package handlers

import (
	"fmt"
	"nosql_db/internal/api"
	"nosql_db/internal/storage"
)

func handleCreateIndex(coll *storage.Collection, req api.Request) api.Response {
	fieldName := ""
	for k := range req.Query {
		fieldName = k
		break
	}

	if fieldName == "" {
		return api.Response{Status: api.StatusError, Message: "field name required in query"}
	}

	if err := coll.CreateIndex(fieldName, 64); err != nil {
		return api.Response{Status: api.StatusError, Message: fmt.Sprintf("failed to create index: %v", err)}
	}

	return api.Response{
		Status:  api.StatusSuccess,
		Message: fmt.Sprintf("Index created on field '%s'", fieldName),
	}
}
