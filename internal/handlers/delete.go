package handlers

import (
	"fmt"
	"nosql_db/internal/operators"
	"nosql_db/internal/query"
	"nosql_db/internal/storage"
)

// cmdDelete обрабатывает команду удаления документов
func cmdDelete(dbName, jsonQuery string) error {
	q, err := query.Parse(jsonQuery)
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
	allDocs := coll.All()
	deletedCount := 0
	for _, doc := range allDocs {
		if operators.MatchDocument(doc, q.Conditions) {
			if id, ok := doc["_id"].(string); ok {
				if coll.Delete(id) {
					deletedCount++
				}
			}
		}
	}
	if err := coll.Save(); err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}

	if err := coll.SaveAllIndexes(); err != nil {
		return fmt.Errorf("failed to save indexes: %w", err)
	}

	fmt.Printf("Deleted %d document(s).\n", deletedCount)
	return nil
}
