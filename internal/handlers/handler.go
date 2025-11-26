package handlers

import "fmt"

// ExecuteCommand главный обработчик команд
func ExecuteCommand(dbName, command, jsonQuery string) error {
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

// PrintUsage выводит справку по использованию
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  server <database_name> insert '<json_document>'")
	fmt.Println("  server <database_name> find '<json_query>'")
	fmt.Println("  server <database_name> delete '<json_query>'")
	fmt.Println("  server <database_name> create_index <field_name>")
	fmt.Println("\nExamples:")
	fmt.Println(`  server my_database insert '{"name": "Alice", "age": 25, "city": "London"}'`)
	fmt.Println(`  server my_database find '{"age": 25}'`)
	fmt.Println(`  server my_database find '{"age": {"$gt": 20}}'`)
	fmt.Println(`  server my_database delete '{"name": {"$like": "A%"}}'`)
}
