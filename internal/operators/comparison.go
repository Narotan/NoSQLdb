package operators

import (
	"fmt"
	"reflect"
)

// CompareEq проверяет равенство значений
func CompareEq(fieldValue, queryValue any) bool {
	return reflect.DeepEqual(fieldValue, queryValue)
}

// CompareGt проверяет больше (>)
func CompareGt(fieldValue, queryValue any) bool {
	return compareNumeric(fieldValue, queryValue, func(a, b float64) bool { return a > b })
}

// CompareLt проверяет меньше (<)
func CompareLt(fieldValue, queryValue any) bool {
	return compareNumeric(fieldValue, queryValue, func(a, b float64) bool { return a < b })
}

// CompareLike проверяет строку с wildcard-паттерном (% - любая строка, _ - один символ)
func CompareLike(fieldValue, pattern any) bool {
	fieldStr, ok1 := fieldValue.(string)
	patternStr, ok2 := pattern.(string)
	if !ok1 || !ok2 {
		return false
	}

	// Конвертируем SQL LIKE паттерн в regex-подобную проверку
	// % -> любая последовательность символов
	// _ -> один символ
	return matchLikePattern(fieldStr, patternStr)
}

// CompareIn проверяет принадлежность к массиву значений
func CompareIn(fieldValue any, values any) bool {
	// values должен быть массивом
	valuesSlice, ok := values.([]any)
	if !ok {
		return false
	}

	for _, v := range valuesSlice {
		if reflect.DeepEqual(fieldValue, v) {
			return true
		}
	}
	return false
}

// compareNumeric - вспомогательная функция для сравнения числовых значений
func compareNumeric(a, b any, cmp func(float64, float64) bool) bool {
	aNum, err1 := toFloat64(a)
	bNum, err2 := toFloat64(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return cmp(aNum, bNum)
}

// toFloat64 конвертирует различные числовые типы в float64
func toFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", val)
	}
}

// matchLikePattern проверяет соответствие строки LIKE-паттерну
func matchLikePattern(str, pattern string) bool {
	return matchLikeHelper(str, pattern, 0, 0)
}

func matchLikeHelper(str, pattern string, strIdx, patIdx int) bool {
	// Если оба индекса достигли конца - совпадение
	if strIdx == len(str) && patIdx == len(pattern) {
		return true
	}

	// Если паттерн закончился, а строка нет - не совпадение
	if patIdx == len(pattern) {
		return false
	}

	// Обработка %
	if pattern[patIdx] == '%' {
		// Пропускаем последовательные %
		for patIdx < len(pattern) && pattern[patIdx] == '%' {
			patIdx++
		}
		// Если % в конце паттерна - всегда совпадение
		if patIdx == len(pattern) {
			return true
		}
		// Пробуем сопоставить остаток паттерна с различными позициями строки
		for i := strIdx; i <= len(str); i++ {
			if matchLikeHelper(str, pattern, i, patIdx) {
				return true
			}
		}
		return false
	}

	// Если строка закончилась - не совпадение
	if strIdx == len(str) {
		return false
	}

	// Обработка _
	if pattern[patIdx] == '_' {
		return matchLikeHelper(str, pattern, strIdx+1, patIdx+1)
	}

	// Обычный символ
	if str[strIdx] == pattern[patIdx] {
		return matchLikeHelper(str, pattern, strIdx+1, patIdx+1)
	}

	return false
}
