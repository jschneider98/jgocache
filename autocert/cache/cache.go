package cache

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/pbkdf2"

	"github.com/go-redis/redis"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	rediscache "github.com/jschneider98/jgocache/autocert/cache/redis"
	sqlcache "github.com/jschneider98/jgocache/autocert/cache/sql"
)

var _ autocert.Cache = &Cache{}

const (
	// CacheBackendDir specifies autocert's native DirCache caching backend
	// Which uses local directory for storing certificates
	CacheBackendDir string = "dir"

	// CacheBackendRedis specifies redis caching driver
	CacheBackendRedis string = "redis"

	// CacheBackendSQL specifies database/sql family of caching drivers(MySQL, PostgreSQL, etc.)
	CacheBackendSQL string = "sql"
)

// Cache struct to store high level cache instance
type Cache struct {
	backend       autocert.Cache
	encryptionKey []byte
	usePrecaching bool
	precache      sync.Map
}

// Get calls backend's Get to retrieve data from cache by key.
// If precaching is used Get will attempt to retrieve data from cache first,
// and if this attempt fails - will ask backend's Get to fill cache and return
// the data.
// If encryption is turned on Get will try to decrypt data retrieved from
// backend's Get before filling cache and returning the data.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {

	if c.usePrecaching {
		data, ok := c.precache.Load(key)

		if ok {
			return data.([]byte), nil
		}
	}

	data, err := c.backend.Get(ctx, key)

	if err != nil {
		return nil, err
	}

	if c.encryptionKey == nil {
		c.precache.Store(key, data)
		return data, nil
	}

	data, err = c.decrypt(data)

	if err != nil {
		return nil, err
	}

	if c.usePrecaching {
		c.precache.Store(key, data)
	}

	return data, nil
}

// Put calls backend's Put to store the data to cache.
// If precache is used it will also fill cache if no errors returned by backend's
// Put.
// If encryption is turned on it will encrypt data before calling backend's Put,
// but precache will be filled with plaintext data.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	var resultData []byte

	if c.encryptionKey != nil {
		var err error
		resultData, err = c.encrypt(data)
		if err != nil {
			return err
		}
	} else {
		resultData = data
	}

	err := c.backend.Put(ctx, key, resultData)
	if err != nil {
		return err
	}

	if c.usePrecaching {
		c.precache.Store(key, data)
	}

	return nil
}

// Delete calls backend's Delete method.
// If precache is used it will first delete from cache no matter if backend's
// Put completed successfully or failed.
func (c *Cache) Delete(ctx context.Context, key string) error {

	if c.usePrecaching {
		c.precache.Delete(key)
	}

	return c.backend.Delete(ctx, key)
}

func (c *Cache) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)

	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("Ciphertext is too short. Probably corrupted data")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, err
}

func (c *Cache) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)

	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// NewCacheFactory returns Cache instance
func NewCacheFactory(options map[string]string) (autocert.Cache, error) {
	var err error

	backend, ok := options["backend"]

	if !ok {
		return nil, fmt.Errorf("Missing required 'backend' parameter")
	}

	var usePrecaching bool
	usePrecachingString, ok := options["usePrecaching"]

	if ok {
		usePrecaching, _ = strconv.ParseBool(usePrecachingString)
	}

	var encKey []byte
	encryptionKey, ok := options["encryptionKey"]

	if ok && encryptionKey != "" {
		h := sha256.New()
		h.Write([]byte(encryptionKey))
		key := h.Sum(nil)
		encKey = pbkdf2.Key(key[:15], key[16:32], 1048, 32, sha256.New)
	}

	var b autocert.Cache

	switch backend {
		case CacheBackendSQL:
			b, err = newSQLCacheBackend(options)

			if err != nil {
				return nil, err
			}
		case CacheBackendDir:
			b, err = newDirCacheBackend(options)

			if err != nil {
				return nil, err
			}
		case CacheBackendRedis:
			b, err = newRedisCacheBackend(options)

			if err != nil {
				return nil, err
			}
	}

	var c autocert.Cache = &Cache{
		encryptionKey: encKey,
		usePrecaching: usePrecaching,
		precache:      sync.Map{},
		backend:       b,
	}

	return c, err
}

func newDirCacheBackend(options map[string]string) (autocert.Cache, error) {
	dir, ok := options["path"]

	if !ok {
		return nil, fmt.Errorf("Option 'path' is required for 'dir' cache backend")
	}

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, err
		}
	}

	return autocert.DirCache(dir), nil
}

func newSQLCacheBackend(options map[string]string) (autocert.Cache, error) {
	driver, ok := options["driver"]

	if !ok {
		return nil, fmt.Errorf("No driver specified")
	}

	dsn, ok := options["dsn"]

	if !ok {
		return nil, fmt.Errorf("dsn option to backend is required")
	}

	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("Error establishing database connection: %s", err)
	}

	if e := db.Ping(); e != nil {
		return nil, fmt.Errorf("Error contacting database: %s", e)
	}

	return sqlcache.NewCache(db)
}

func newRedisCacheBackend(options map[string]string) (autocert.Cache, error) {
	addr, ok := options["addr"]

	if !ok || addr == "" {
		return nil, fmt.Errorf("No 'addr' specified for Redis cache")
	}

	password, _ := options["password"]

	dbS, ok := options["db"]

	if !ok {
		dbS = "0"
	}

	db, err := strconv.Atoi(dbS)

	if err != nil {
		return nil, fmt.Errorf("Error parsing db field: %s", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return rediscache.NewCache(client)
}
