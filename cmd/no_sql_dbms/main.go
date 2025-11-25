package main

import (
	"encoding/json"
	"fmt"
	"nosql_db/internal/operators"
	"nosql_db/internal/query"
	"nosql_db/internal/storage"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	dbName := os.Args[1]
	command := os.Args[2]

	var jsonQuery string
	if len(os.Args) >= 4 {
		jsonQuery = os.Args[3]
	}

	if err := executeCommand(dbName, command, jsonQuery); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeCommand(dbName, command, jsonQuery string) error {
	switch command {
	case "insert":
		return cmdInsert(dbName, jsonQuery)
	case "find":
		return cmdFind(dbName, jsonQuery)
	case "delete":
		return cmdDelete(dbName, jsonQuery)
	case "create_index":
		return cmdCreateIndex(dbName, jsonQuery)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func cmdInsert(dbName, jsonDoc string) error {
	if jsonDoc == "" {
		return fmt.Errorf("insert requires a JSON document")
	}

	// Парсим документ
	doc, err := query.ParseDocument(jsonDoc)
	if err != nil {
		return err
	}

	// Загружаем коллекцию
	coll, err := storage.LoadCollection(dbName)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	// Вставляем документ
	id, err := coll.Insert(doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	// Сохраняем коллекцию
	if err := coll.Save(); err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}

	fmt.Printf("Document inserted successfully. ID: %s\n", id)
	return nil
}

func cmdFind(dbName, jsonQuery string) error {
	// Парсим запрос
	q, err := query.Parse(jsonQuery)
	if err != nil {
		return err
	}

	// Загружаем коллекцию
	coll, err := storage.LoadCollection(dbName)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	// Получаем все документы и фильтруем
	allDocs := coll.All()
	results := []map[string]any{}

	for _, doc := range allDocs {
		if operators.MatchDocument(doc, q.Conditions) {
			results = append(results, doc)
		}
	}

	// Выводим результаты в JSON
	if len(results) == 0 {
		fmt.Println("[]")
		return nil
	}

	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func cmdDelete(dbName, jsonQuery string) error {
	// Парсим запрос
	q, err := query.Parse(jsonQuery)
	if err != nil {
		return err
	}

	// Загружаем коллекцию
	coll, err := storage.LoadCollection(dbName)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	// Находим документы для удаления
	allDocs := coll.All()
	deletedCount := 0

	for _, doc := range allDocs {
		if operators.MatchDocument(doc, q.Conditions) {
			// Получаем _id и удаляем
			if id, ok := doc["_id"].(string); ok {
				if coll.Delete(id) {
					deletedCount++
				}
			}
		}
	}

	// Сохраняем коллекцию
	if err := coll.Save(); err != nil {
		return fmt.Errorf("failed to save collection: %w", err)
	}

	fmt.Printf("Deleted %d document(s).\n", deletedCount)
	return nil
}

func cmdCreateIndex(dbName, fieldName string) error {
	// Задание со звездочкой - пока заглушка
	fmt.Printf("Index creation on field '%s' for collection '%s' is not yet implemented (bonus task).\n", fieldName, dbName)
	return nilx
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  no_sql_dbms <database_name> insert '<json_document>'")
	fmt.Println("  no_sql_dbms <database_name> find '<json_query>'")
	fmt.Println("  no_sql_dbms <database_name> delete '<json_query>'")
	fmt.Println("  no_sql_dbms <database_name> create_index <field_name>")
	fmt.Println("\nExamples:")
	fmt.Println(`  no_sql_dbms my_database insert '{"name": "Alice", "age": 25, "city": "London"}'`)
	fmt.Println(`  no_sql_dbms my_database find '{"age": 25}'`)
	fmt.Println(`  no_sql_dbms my_database find '{"age": {"$gt": 20}}'`)
	fmt.Println(`  no_sql_dbms my_database delete '{"name": {"$like": "A%"}}'`)
}
