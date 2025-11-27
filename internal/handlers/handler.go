package handlers

import (
	"fmt"
	"nosql_db/internal/api"
	"nosql_db/internal/storage"
)

// HandleRequest — точка входа для обработки запросов
func HandleRequest(req api.Request) api.Response {
	if req.Database == "" {
		return api.Response{Status: api.StatusError, Message: "database name is required"}
	}

	coll, err := storage.LoadCollection(req.Database)
	if err != nil {
		return api.Response{Status: api.StatusError, Message: fmt.Sprintf("failed to load database: %v", err)}
	}

	_ = coll.LoadAllIndexes()

	switch req.Command {
	case api.CmdInsert:
		return handleInsert(coll, req)
	case api.CmdFind:
		return handleFind(coll, req)
	case api.CmdDelete:
		return handleDelete(coll, req)
	case api.CmdCreateIndex:
		return handleCreateIndex(coll, req)
	default:
		return api.Response{Status: api.StatusError, Message: fmt.Sprintf("unknown command: %s", req.Command)}
	}
}
