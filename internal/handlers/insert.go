package handlers

import (
	"fmt"
	"nosql_db/internal/api"
	"nosql_db/internal/storage"
)

func handleInsert(coll *storage.Collection, req api.Request) api.Response {
	if len(req.Data) == 0 {
		return api.Response{Status: api.StatusError, Message: "no data provided for insert"}
	}

	count := 0
	for _, doc := range req.Data {
		_, err := coll.Insert(doc)
		if err != nil {
			return api.Response{Status: api.StatusError, Message: fmt.Sprintf("insert error: %v", err)}
		}
		count++
	}

	if err := coll.Save(); err != nil {
		return api.Response{Status: api.StatusError, Message: "failed to save data"}
	}

	if err := coll.SaveAllIndexes(); err != nil {
		return api.Response{Status: api.StatusError, Message: "failed to save indexes"}
	}

	return api.Response{
		Status:  api.StatusSuccess,
		Message: fmt.Sprintf("Inserted %d document(s)", count),
		Count:   count,
	}
}
