package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadCollection загружает коллекцию из базы данных
func LoadCollection(name string) (*Collection, error) {
	path := filepath.Join("data", name+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewCollection(name), nil
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]any
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	hmap := NewHashMap()
	for k, v := range raw {
		hmap.Put(k, v)
	}
	coll := NewCollection(name)
	coll.Data = hmap
	return coll, nil
}

// Save сохраняет данные в json в базе данных
func (c *Collection) Save() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	items := c.Data.Items()
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}
	path := filepath.Join("data", c.Name+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}
