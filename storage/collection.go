package storage

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Collection struct {
	Name string
	Data *HashMap
}

func NewCollection(name string) *Collection {
	return &Collection{
		Name: name,
		Data: NewHashMap(),
	}
}

// generateID — простая генерация уникального ID
func generateID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000000))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LoadCollection загружает колллекцию из базы данных
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

	return NewCollection(name), nil
}

// Save сохраняет данные в json в базу данных
func (c *Collection) Save() error {
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

func (c *Collection) Insert(doc map[string]any) (string, error) {
	id := generateID()
	doc["_id"] = id
	c.Data.Put(id, doc)
	return id, nil
}

// GetByID получает документ по _id
func (c *Collection) GetByID(id string) (map[string]any, bool) {
	val, ok := c.Data.Get(id)
	if !ok {
		return nil, false
	}

	doc, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}
	return doc, true
}

// Delete удаляет документ по _id
func (c *Collection) Delete(id string) bool {
	return c.Data.Remove(id)
}

// All возвращает все документы как []map[string]any
func (c *Collection) All() []map[string]any {
	result := []map[string]any{}
	for _, v := range c.Data.Items() {
		if doc, ok := v.(map[string]any); ok {
			result = append(result, doc)
		}
	}
	return result
}
