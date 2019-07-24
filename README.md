# jgocache
jgocache is almost an exact copy of teran's (Igor Shishkin) autocert cache code used in his reverse proxy project (https://github.com/teran/svcproxy). Most of the credit for the code is his.
***
I've made a few modest changes to the code for my needs.

* Isolated the cache code out of the reverse proxy.
* I removed the "migrateMaybe" logic for the DB schema. It didn't seem worth the additional code dependencies and didn't always work for me in my product environments. As such if you're using a SQL backend you'll have to either use your own schema migration code or manually add the table.
* I simplified the cache factory interface.

```
// original interface
func NewCacheFactory(backend CacheBackend, options map[string]string) (autocert.Cache, error) {

// updated interface
func NewCacheFactory(options map[string]string) (autocert.Cache, error) {

```

* I corrected a place or two where a precache conditional was missing.
### Future
***
I may add additional caching beyond autocert.
