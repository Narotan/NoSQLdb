package handlers

import (
	"fmt"
	"nosql_db/internal/query"
	"nosql_db/internal/storage"
)

// cmdInsert обрабатывает команду вставки документа
func cmdInsert(dbName, jsonDoc string) error {
	if jsonDoc == "" {
		return fmt.Errorf("insert requires a JSON document")
	}
	doc, err := query.ParseDocument(jsonDoc)
	if err != nil {
		return err
	}
	coll, err := storage.LoadCollection(dbName)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}
	if err := coll.LoadAllIndexes(); err != nil {
		return fmt.Errorf("failed to load indexes: %w", err)
	}
	id, err := coll.Insert(doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}
	if err := coll.Save(); err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}
	if err := coll.SaveAllIndexes(); err != nil {
		return fmt.Errorf("failed to save indexes: %w", err)
	}
	fmt.Printf("Document inserted successfully. ID: %s\n", id)
	return nil
}
