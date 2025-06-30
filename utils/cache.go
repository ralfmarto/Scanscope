package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Cache handles scanning cache.
type Cache struct {
	Entries map[string]bool `json:"entries"`
}

// LoadCache reads cache from file.
func LoadCache(path string) (*Cache, error) {
	c := &Cache{Entries: make(map[string]bool)}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return nil, err
	}
	json.Unmarshal(data, &c.Entries)
	return c, nil
}

// Has returns true if hash exists.
func (c *Cache) Has(hash string) bool {
	return c.Entries[hash]
}

// Add adds hash to cache.
func (c *Cache) Add(hash string) {
	c.Entries[hash] = true
}

// Save writes cache to file.
func (c *Cache) Save(path string) error {
	data, err := json.MarshalIndent(c.Entries, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
