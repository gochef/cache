package cache

import "fmt"

type (
	// Driver represents a cache driver instance
	Driver interface {
		Get(key string) (interface{}, bool)
		Put(key string, data interface{}, expire int64)
		Remove(key string)
		Clear()
	}

	// Config is the cache instance configuration
	Config struct {
		Driver   string
		MaxSize  int
		Address  string
		Username string
		Password string
		Use bool
	}

	// Cache represents a cache instance
	Cache struct {
		config *Config
		driver Driver
	}
)

var (
	drivers = map[string]Driver{
		"memory": NewMemoryCache(),
	}
)

// New returns a cache instance with provided config
func New(cfg *Config) *Cache {
	driver, ok := drivers[cfg.Driver]
	if !ok {
		errStr := fmt.Sprintf("cache: cache provider %s is not registered", cfg.Driver)
		panic(errStr)
	}

	return &Cache{
		config: cfg,
		driver: driver,
	}
}

// RegisterDriver registers a driver
// panics if driver is already registered
func RegisterDriver(name string, driver Driver) {
	if driver == nil {
		panic("cache: driver is nil")
	}

	if _, ok := drivers[name]; ok {
		errStr := fmt.Sprintf("cache: driver %s is already registered", name)
		panic(errStr)
	}

	drivers[name] = driver
}

// Get fetches an item from session store by key,
// returns an empty interface and false if it doesnt exist
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.driver.Get(key)
}

// GetString returns a string item from session store
func (c *Cache) GetString(key string) (string, bool) {
	data, ok := c.Get(key)
	if !ok {
		return "", false
	}

	str, ok := data.(string)
	return str, ok
}

// GetInt returns an integer item from session store
func (c *Cache) GetInt(key string) (int, bool) {
	data, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	str, ok := data.(int)
	return str, ok
}

// Put adds an item to cache for the specified duration
// identified by provided key
func (c *Cache) Put(key string, data interface{}, duration int64) {
	c.driver.Put(key, data, duration)
}

// PutForever adds an item to the cache forever
func (c *Cache) PutForever(key string, data interface{}) {
	c.driver.Put(key, data, 0)
}

// Remove deletes an item from session store by provided key
func (c *Cache) Remove(key string) {
	c.driver.Remove(key)
}

// Pull gets an item from session store and deletes the item from session
func (c *Cache) Pull(key string) (interface{}, bool) {
	data, ok := c.driver.Get(key)
	c.driver.Remove(key)

	return data, ok
}

// PullString gets a string item from session store and deletes the item from session
func (c *Cache) PullString(key string) (string, bool) {
	data, ok := c.GetString(key)
	c.driver.Remove(key)

	return data, ok
}

// PullInt gets an integer item from session store and deletes the item from session
func (c *Cache) PullInt(key string) (int, bool) {
	data, ok := c.GetInt(key)
	c.driver.Remove(key)

	return data, ok
}

// Remember fetches an item from the cache, if the item does not exist,
// passed callback is executed, the data from the callback is stored in the cache
// for the passed duration and returned to the caller
func (c *Cache) Remember(key string, duration int64, cb func() interface{}) interface{} {
	data, ok := c.Get(key)
	if ok {
		return data
	}

	data = cb()
	c.Put(key, data, duration)
	return data
}

// RememberForever does the same as Remember except, the data is stored forever
func (c *Cache) RememberForever(key string, cb func() interface{}) interface{} {
	return c.Remember(key, 0, cb)
}

// Clear empties the session store
func (c *Cache) Clear() {
	c.driver.Clear()
}
