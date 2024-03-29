# Autocert/cache

`autocert.Cache` implementations
Allow the use of precaching for realy fast reads and encryption for better security.

# Configuration

## Dir Cache backend

Simply wraps [`golang.org/x/crypto/acme/autocert.DirCache`](https://godoc.org/golang.org/x/crypto/acme/autocert#DirCache)

Required options:
 * `path` - path to the directory to store cache. Directory may exist or will be
   created automatically.

## Redis Cache backend

Redis cache

Required options:
* `addr` - Redis database address

Additional options:
* `password` - database password
* `db` - database (0 - for default database)

## SQL Cache backend

Backend allows to use relational database as a storage.

Required options:
 * `driver` - database driver to use. Allowed options: `mysql`, `postgres`
 * `dsn` - Data source name refering database in the form supported by Go's drivers.
   Examples:
    - MySQL: `root@tcp(127.0.0.1:3306)/db_name?parseTime=true`
    - PostgreSQL: `postgres://postgres@localhost/db_name?sslmode=disable`

   More examples can be found in tests.

# Benchmarks

```
› make benchmark
go test -bench=. -cpu=1,2,3,4
goos: darwin
goarch: amd64
pkg: jgocache/autocert/cache
BenchmarkCacheGetSQLMySQLWithoutEncryptionAndPrecaching              	    2000	    916999 ns/op
BenchmarkCacheGetSQLMySQLWithoutEncryptionAndPrecaching-2            	    2000	    930788 ns/op
BenchmarkCacheGetSQLMySQLWithoutEncryptionAndPrecaching-3            	    2000	    926944 ns/op
BenchmarkCacheGetSQLMySQLWithoutEncryptionAndPrecaching-4            	    2000	    928512 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndWithoutPrecaching          	    2000	    927463 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndWithoutPrecaching-2        	    2000	    944562 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndWithoutPrecaching-3        	    2000	    942038 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndWithoutPrecaching-4        	    2000	    943408 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndPrecaching                 	30000000	        39.0 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndPrecaching-2               	30000000	        38.0 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndPrecaching-3               	30000000	        38.4 ns/op
BenchmarkCacheGetSQLMySQLWithEncryptionAndPrecaching-4               	30000000	        37.9 ns/op
BenchmarkCacheGetSQLPostgreSQLWithoutEncryptionAndPrecaching         	    2000	    832391 ns/op
BenchmarkCacheGetSQLPostgreSQLWithoutEncryptionAndPrecaching-2       	    2000	    837306 ns/op
BenchmarkCacheGetSQLPostgreSQLWithoutEncryptionAndPrecaching-3       	    2000	    851568 ns/op
BenchmarkCacheGetSQLPostgreSQLWithoutEncryptionAndPrecaching-4       	    2000	    945629 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndWithoutPrecaching     	    2000	   1037167 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndWithoutPrecaching-2   	    2000	    860692 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndWithoutPrecaching-3   	    2000	    934574 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndWithoutPrecaching-4   	    2000	    932242 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndPrecaching            	30000000	        39.6 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndPrecaching-2          	30000000	        41.8 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndPrecaching-3          	30000000	        38.6 ns/op
BenchmarkCacheGetSQLPostgreSQLWithEncryptionAndPrecaching-4          	30000000	        38.7 ns/op
BenchmarkCacheGetDirWithoutEncryptionAndPrecaching                   	   50000	     25212 ns/op
BenchmarkCacheGetDirWithoutEncryptionAndPrecaching-2                 	   50000	     28623 ns/op
BenchmarkCacheGetDirWithoutEncryptionAndPrecaching-3                 	   50000	     28575 ns/op
BenchmarkCacheGetDirWithoutEncryptionAndPrecaching-4                 	   50000	     28373 ns/op
BenchmarkCacheGetDirWithEncryptionAndWithoutPrecaching               	   50000	     26619 ns/op
BenchmarkCacheGetDirWithEncryptionAndWithoutPrecaching-2             	   50000	     30195 ns/op
BenchmarkCacheGetDirWithEncryptionAndWithoutPrecaching-3             	   50000	     30278 ns/op
BenchmarkCacheGetDirWithEncryptionAndWithoutPrecaching-4             	   50000	     30136 ns/op
BenchmarkCacheGetDirWithEncryptionAndPrecaching                      	30000000	        40.8 ns/op
BenchmarkCacheGetDirWithEncryptionAndPrecaching-2                    	30000000	        39.3 ns/op
BenchmarkCacheGetDirWithEncryptionAndPrecaching-3                    	30000000	        39.4 ns/op
BenchmarkCacheGetDirWithEncryptionAndPrecaching-4                    	30000000	        39.4 ns/op
PASS
ok  	jgocache/autocert/cache	60.091s
```
