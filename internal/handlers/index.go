package handlers

import (
	"fmt"
	"nosql_db/internal/storage"
)

// cmdCreateIndex обрабатывает команду создания индекса
func cmdCreateIndex(dbName, fieldName string) error {
	coll, err := storage.LoadCollection(dbName)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}
	if err := coll.CreateIndex(fieldName, 64); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	fmt.Printf("Index created successfully on field '%s'.\n", fieldName)
	return nil
}
