package benchmarks

import (
	"testing"

	"github.com/go-json-experiment/json"

	jp "github.com/kaptinlin/jsonpointer"
)

// Standard JSON parsing produces generic types - most common real-world scenario
func BenchmarkStandardJSON(b *testing.B) {
	jsonData := `{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}]}`
	var data any
	json.Unmarshal([]byte(jsonData), &data)

	for b.Loop() {
		jp.Get(data, "users", "0", "name")
	}
}

// ===== Slice type performance comparison =====

// Specialized type: []string
func BenchmarkStringSlice_Specialized(b *testing.B) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// Generic type: []any (for comparison)
func BenchmarkStringSlice_Generic(b *testing.B) {
	data := []any{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// Specialized type: []int
func BenchmarkIntSlice_Specialized(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// Generic type: []any (for comparison)
func BenchmarkIntSlice_Generic(b *testing.B) {
	data := []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// Specialized type: []float64
func BenchmarkFloat64Slice_Specialized(b *testing.B) {
	data := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9, 10.0}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// Generic type: []any (for comparison)
func BenchmarkFloat64Slice_Generic(b *testing.B) {
	data := []any{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9, 10.0}

	for b.Loop() {
		jp.Get(data, "5")
	}
}

// ===== Map type performance comparison =====

// Specialized type: map[string]string
func BenchmarkMapStringString_Specialized(b *testing.B) {
	data := map[string]string{
		"a": "value1", "b": "value2", "c": "value3",
		"d": "value4", "e": "value5", "f": "value6",
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// Generic type: map[string]any
func BenchmarkMapStringString_Generic(b *testing.B) {
	data := map[string]any{
		"a": "value1", "b": "value2", "c": "value3",
		"d": "value4", "e": "value5", "f": "value6",
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// Specialized type: map[string]int
func BenchmarkMapStringInt_Specialized(b *testing.B) {
	data := map[string]int{
		"a": 1, "b": 2, "c": 3,
		"d": 4, "e": 5, "f": 6,
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// Generic type: map[string]any
func BenchmarkMapStringInt_Generic(b *testing.B) {
	data := map[string]any{
		"a": 1, "b": 2, "c": 3,
		"d": 4, "e": 5, "f": 6,
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// Specialized type: map[string]float64
func BenchmarkMapStringFloat64_Specialized(b *testing.B) {
	data := map[string]float64{
		"a": 1.1, "b": 2.2, "c": 3.3,
		"d": 4.4, "e": 5.5, "f": 6.6,
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// Generic type: map[string]any
func BenchmarkMapStringFloat64_Generic(b *testing.B) {
	data := map[string]any{
		"a": 1.1, "b": 2.2, "c": 3.3,
		"d": 4.4, "e": 5.5, "f": 6.6,
	}

	for b.Loop() {
		jp.Get(data, "d")
	}
}

// ===== Nested access performance comparison =====

// Specialized nested type: map[string][]string
func BenchmarkNested_Specialized(b *testing.B) {
	data := map[string][]string{
		"users": {"Alice", "Bob", "Charlie", "David", "Eve"},
	}

	for b.Loop() {
		jp.Get(data, "users", "2")
	}
}

// Generic nested type: map[string]any + []any
func BenchmarkNested_Generic(b *testing.B) {
	data := map[string]any{
		"users": []any{"Alice", "Bob", "Charlie", "David", "Eve"},
	}

	for b.Loop() {
		jp.Get(data, "users", "2")
	}
}

// ===== Find operation performance comparison =====

// Find with specialized types
func BenchmarkFind_Specialized(b *testing.B) {
	data := map[string][]int{
		"scores": {100, 90, 85, 95, 88},
	}

	for b.Loop() {
		jp.Find(data, "scores", "3")
	}
}

// Find with generic types
func BenchmarkFind_Generic(b *testing.B) {
	data := map[string]any{
		"scores": []any{100, 90, 85, 95, 88},
	}

	for b.Loop() {
		jp.Find(data, "scores", "3")
	}
}
