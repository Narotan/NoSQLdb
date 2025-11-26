package handlers

import (
	"encoding/json"
	"fmt"
	"nosql_db/internal/index"
	"nosql_db/internal/operators"
	"nosql_db/internal/query"
	"nosql_db/internal/storage"
)

// cmdFind обрабатывает команду поиска документов
func cmdFind(dbName, jsonQuery string) error {
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
	var results []map[string]any
	if len(q.Conditions) == 1 && !hasLogicalOperators(q.Conditions) {
		for field, condition := range q.Conditions {
			if coll.HasIndex(field) {
				results = findWithIndex(coll, field, condition)
				break
			}
		}
	}
	if results == nil {
		results = findFullScan(coll, q)
	}
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

// hasLogicalOperators проверяет наличие $or или $and в условиях
func hasLogicalOperators(conditions map[string]any) bool {
	_, hasOr := conditions["$or"]
	_, hasAnd := conditions["$and"]
	return hasOr || hasAnd
}

// findWithIndex выполняет поиск используя индекс
func findWithIndex(coll *storage.Collection, field string, condition any) []map[string]any {
	btree, ok := coll.GetIndex(field)
	if !ok {
		return nil
	}
	var docIDs []string
	switch v := condition.(type) {
	case float64, int, int64, string, bool:
		key := index.ValueToKey(v)
		values := btree.Search(key)
		docIDs = index.ValuesToStrings(values)
	case map[string]any:
		if gtValue, exists := v["$gt"]; exists {
			key := index.ValueToKey(gtValue)
			values := btree.SearchGreaterThan(key)
			docIDs = index.ValuesToStrings(values)
		} else if ltValue, exists := v["$lt"]; exists {
			key := index.ValueToKey(ltValue)
			values := btree.SearchLessThan(key)
			docIDs = index.ValuesToStrings(values)
		} else if eqValue, exists := v["$eq"]; exists {
			key := index.ValueToKey(eqValue)
			values := btree.Search(key)
			docIDs = index.ValuesToStrings(values)
		} else if inValues, exists := v["$in"]; exists {
			if inArray, ok := inValues.([]any); ok {
				var keys []index.Key
				for _, val := range inArray {
					keys = append(keys, index.ValueToKey(val))
				}
				values := btree.SearchIn(keys)
				docIDs = index.ValuesToStrings(values)
			}
		}
	}
	var results []map[string]any
	for _, id := range docIDs {
		if doc, ok := coll.GetByID(id); ok {
			results = append(results, doc)
		}
	}
	return results
}

// findFullScan выполняет полное сканирование коллекции
func findFullScan(coll *storage.Collection, q *query.Query) []map[string]any {
	var results []map[string]any
	allDocs := coll.All()
	for _, doc := range allDocs {
		if operators.MatchDocument(doc, q.Conditions) {
			results = append(results, doc)
		}
	}
	return results
}
