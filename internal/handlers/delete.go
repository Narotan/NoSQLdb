package handlers

import (
	"fmt"
	"nosql_db/internal/api"
	"nosql_db/internal/storage"
)

func handleDelete(coll *storage.Collection, req api.Request) api.Response {
	// Сначала находим документы для удаления через FullScan
	candidates := findFullScan(coll, req.Query)
	deletedCount := 0

	for _, doc := range candidates {
		id, ok := doc["_id"].(string)
		if ok {
			if coll.Delete(id) {
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		if err := coll.Save(); err != nil {
			return api.Response{Status: api.StatusError, Message: "failed to save changes"}
		}
		if err := coll.RebuildAllIndexes(); err != nil {
			return api.Response{Status: api.StatusError, Message: "failed to rebuild indexes"}
		}
	}

	return api.Response{
		Status:  api.StatusSuccess,
		Message: fmt.Sprintf("Deleted %d document(s)", deletedCount),
		Count:   deletedCount,
	}
}
