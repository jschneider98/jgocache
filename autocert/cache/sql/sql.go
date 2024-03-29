package sql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/crypto/acme/autocert"

	"github.com/jschneider98/jgocache/autocert/cache/sql/mysql"
	"github.com/jschneider98/jgocache/autocert/cache/sql/postgresql"
)

var _ autocert.Cache = &Cache{}

// Cache implements autocert.Cache with a SQL database
type Cache struct {
	driver autocert.Cache
}

// NewCache returns Cache instance
func NewCache(db *sql.DB) (*Cache, error) {
	var driver autocert.Cache

	switch fmt.Sprintf("Driver: %s", reflect.TypeOf(db.Driver())) {
	case "Driver: *mysql.MySQLDriver":
		driver = &mysql.MySQL{
			DB: db,
		}

	case "Driver: *pq.Driver":
		driver = &postgresql.PostgreSQL{
			DB: db,
		}

	default:
		return nil, errors.New("Unsupported driver")
	}

	return &Cache{
		driver: driver,
	}, nil
}

// Get retrieves certificate data from cache
func (m *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := m.driver.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	data, err = m.decode(data)
	if err != nil {
		return nil, err
	}

	return data, err
}

// Put stores certificate data to cache
func (m *Cache) Put(ctx context.Context, key string, data []byte) error {
	data = m.encode(data)
	err := m.driver.Put(ctx, key, data)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes certificate data from cache
func (m *Cache) Delete(ctx context.Context, key string) error {
	return m.driver.Delete(ctx, key)
}

func (m *Cache) decode(input []byte) ([]byte, error) {
	ct := make([]byte, base64.StdEncoding.DecodedLen(len(input)))
	l, err := base64.StdEncoding.Decode(ct, input)
	if err != nil {
		return nil, err
	}

	return ct[:l], nil
}

func (m *Cache) encode(input []byte) []byte {
	ct := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(ct, input)

	return ct
}
