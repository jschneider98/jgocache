package cache

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/acme/autocert"
)

type CacheTestSuite struct {
	suite.Suite
}

func (s *CacheTestSuite) TestInitializeSQLCacheWithEncryption() {
	options := map[string]string{
		"backend":       "sql",
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/test_db?parseTime=true",
		"encryptionKey": "testkey",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	// Put data sample
	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	// Get data sample back
	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func (s *CacheTestSuite) TestInitializeSQLCacheNoEncryption() {
	options := map[string]string{
		"backend": "sql",
		"driver":  "mysql",
		"dsn":     "root@tcp(127.0.0.1:3306)/test_db?parseTime=true",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func (s *CacheTestSuite) TestInitializeSQLCacheWithPrecaching() {
	options := map[string]string{
		"backend":       "sql",
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/test_db?parseTime=true",
		"usePrecaching": "true",
		"encryptionKey": "",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	// Put data sample
	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	// Get data sample back
	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func (s *CacheTestSuite) TestInitializeSQLCacheWithEncryptionAndPrecaching() {
	options := map[string]string{
		"backend":       "sql",
		"driver":        "mysql",
		"dsn":           "root@tcp(127.0.0.1:3306)/test_db?parseTime=true",
		"encryptionKey": "testkey",
		"usePrecaching": "true",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	// Put data sample
	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	// Get data sample back
	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func (s *CacheTestSuite) TestInitializeDirCacheWithEncryptionAndPrecaching() {
	dir, err := ioutil.TempDir("", "dircache_test")
	s.Require().NoError(err)

	options := map[string]string{
		"backend":       "dir",
		"path":          dir,
		"encryptionKey": "testkey",
		"usePrecaching": "true",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	// Put data sample
	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	// Get data sample back
	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func (s *CacheTestSuite) TestInitializeRedisCacheNoEncryption() {
	options := map[string]string{
		"backend": "redis",
		"addr": "127.0.0.1:6379",
	}
	c, err := NewCacheFactory(options)
	s.Require().NoError(err)

	err = c.Put(context.Background(), "test-key", []byte("test-data"))
	s.Require().NoError(err)

	data, err := c.Get(context.Background(), "test-key")
	s.Require().NoError(err)
	s.Equal([]byte("test-data"), data)

	// Delete sample
	err = c.Delete(context.Background(), "test-key")
	s.Require().NoError(err)

	// Test if sample is still present
	data, err = c.Get(context.Background(), "test-key")
	s.Require().Equal(autocert.ErrCacheMiss, err)
	s.Equal([]byte(nil), data)
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
