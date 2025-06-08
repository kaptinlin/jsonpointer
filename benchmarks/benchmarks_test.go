package benchmarks

import (
	"encoding/json"
	"testing"

	// Our implementation
	ourjp "github.com/kaptinlin/jsonpointer"

	// Comparison library 1: go-openapi/jsonpointer
	gojp "github.com/go-openapi/jsonpointer"

	// Comparison library 2: BragdonD/jsonpointer-go
	bragdonjp "github.com/bragdond/jsonpointer-go"

	// Comparison library 3: woodsbury/jsonpointer
	woodsjp "github.com/woodsbury/jsonpointer"

	// Comparison library 4: dolmen-go/jsonptr
	dolmenjp "github.com/dolmen-go/jsonptr"
)

// Test data
var (
	smallData = map[string]any{
		"name": "Alice",
		"age":  30,
		"profile": map[string]any{
			"email": "alice@example.com",
			"settings": map[string]any{
				"theme":         "dark",
				"notifications": true,
			},
		},
		"hobbies": []any{"reading", "coding", "music"},
		"scores":  []any{95, 87, 92},
	}

	mediumData = generateMediumData()

	// Struct test data
	structData = generateStructData()
	mapData    = generateMapData()
)

func generateMediumData() map[string]any {
	users := make([]any, 100)
	for i := 0; i < 100; i++ {
		users[i] = map[string]any{
			"id":   i,
			"name": "User " + string(rune(i)),
			"profile": map[string]any{
				"email": "user" + string(rune(i)) + "@example.com",
				"age":   20 + (i % 50),
			},
		}
	}
	return map[string]any{
		"users": users,
		"total": 100,
	}
}

// Struct definitions for benchmark testing
type BenchUser struct {
	Name    string       `json:"name"`
	Age     int          `json:"age"`
	Profile BenchProfile `json:"profile"`
	Hobbies []string     `json:"hobbies"`
	Scores  []int        `json:"scores"`
}

type BenchProfile struct {
	Email    string        `json:"email"`
	Settings BenchSettings `json:"settings"`
}

type BenchSettings struct {
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
}

func generateStructData() BenchUser {
	return BenchUser{
		Name: "Alice",
		Age:  30,
		Profile: BenchProfile{
			Email: "alice@example.com",
			Settings: BenchSettings{
				Theme:         "dark",
				Notifications: true,
			},
		},
		Hobbies: []string{"reading", "coding", "music"},
		Scores:  []int{95, 87, 92},
	}
}

func generateMapData() map[string]any {
	return map[string]any{
		"name": "Alice",
		"age":  30,
		"profile": map[string]any{
			"email": "alice@example.com",
			"settings": map[string]any{
				"theme":         "dark",
				"notifications": true,
			},
		},
		"hobbies": []any{"reading", "coding", "music"},
		"scores":  []any{95, 87, 92},
	}
}

// ===== Our library benchmarks =====

func BenchmarkOur_Find_Root(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("", smallData)
	}
}

func BenchmarkOur_Get_Root(b *testing.B) {
	path := ourjp.ParseJsonPointer("")
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(smallData, path)
	}
}

func BenchmarkOur_Find_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/name", smallData)
	}
}

func BenchmarkOur_Get_Shallow(b *testing.B) {
	path := ourjp.ParseJsonPointer("/name")
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(smallData, path)
	}
}

func BenchmarkOur_Find_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/profile/settings/theme", smallData)
	}
}

func BenchmarkOur_Get_Deep(b *testing.B) {
	path := ourjp.ParseJsonPointer("/profile/settings/theme")
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(smallData, path)
	}
}

func BenchmarkOur_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.ParseJsonPointer("/profile/settings/theme")
	}
}

func BenchmarkOur_Medium_FindUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/users/50/name", mediumData)
	}
}

func BenchmarkOur_NotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/nonexistent", smallData)
	}
}

func BenchmarkOur_PrecompiledPath(b *testing.B) {
	path := ourjp.ParseJsonPointer("/profile/settings/theme")
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.Find(smallData, path)
	}
}

// ===== go-openapi/jsonpointer benchmarks =====

func BenchmarkGoOpenAPI_Get_Root(b *testing.B) {
	jsonData, _ := json.Marshal(smallData)
	for i := 0; i < b.N; i++ {
		ptr, _ := gojp.New("")
		_, _, _ = ptr.Get(jsonData)
	}
}

func BenchmarkGoOpenAPI_Get_Shallow(b *testing.B) {
	jsonData, _ := json.Marshal(smallData)
	for i := 0; i < b.N; i++ {
		ptr, _ := gojp.New("/name")
		_, _, _ = ptr.Get(jsonData)
	}
}

func BenchmarkGoOpenAPI_Get_Deep(b *testing.B) {
	jsonData, _ := json.Marshal(smallData)
	for i := 0; i < b.N; i++ {
		ptr, _ := gojp.New("/profile/settings/theme")
		_, _, _ = ptr.Get(jsonData)
	}
}

func BenchmarkGoOpenAPI_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = gojp.New("/profile/settings/theme")
	}
}

func BenchmarkGoOpenAPI_Medium_GetUser(b *testing.B) {
	jsonData, _ := json.Marshal(mediumData)
	for i := 0; i < b.N; i++ {
		ptr, _ := gojp.New("/users/50/name")
		_, _, _ = ptr.Get(jsonData)
	}
}

func BenchmarkGoOpenAPI_NotFound(b *testing.B) {
	jsonData, _ := json.Marshal(smallData)
	for i := 0; i < b.N; i++ {
		ptr, _ := gojp.New("/nonexistent")
		_, _, _ = ptr.Get(jsonData)
	}
}

func BenchmarkGoOpenAPI_PrecompiledPath(b *testing.B) {
	jsonData, _ := json.Marshal(smallData)
	ptr, _ := gojp.New("/profile/settings/theme")
	for i := 0; i < b.N; i++ {
		_, _, _ = ptr.Get(jsonData)
	}
}

// ===== BragdonD/jsonpointer-go benchmarks =====

func BenchmarkBragdonD_Get_Root(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr, _ := bragdonjp.NewJSONPointer("")
		_, _ = ptr.GetValue(smallData)
	}
}

func BenchmarkBragdonD_Get_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr, _ := bragdonjp.NewJSONPointer("/name")
		_, _ = ptr.GetValue(smallData)
	}
}

func BenchmarkBragdonD_Get_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr, _ := bragdonjp.NewJSONPointer("/profile/settings/theme")
		_, _ = ptr.GetValue(smallData)
	}
}

func BenchmarkBragdonD_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = bragdonjp.NewJSONPointer("/profile/settings/theme")
	}
}

func BenchmarkBragdonD_Medium_GetUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr, _ := bragdonjp.NewJSONPointer("/users/50/name")
		_, _ = ptr.GetValue(mediumData)
	}
}

func BenchmarkBragdonD_NotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr, _ := bragdonjp.NewJSONPointer("/nonexistent")
		_, _ = ptr.GetValue(smallData)
	}
}

func BenchmarkBragdonD_Precompiled(b *testing.B) {
	ptr, _ := bragdonjp.NewJSONPointer("/profile/settings/theme")
	for i := 0; i < b.N; i++ {
		_, _ = ptr.GetValue(smallData)
	}
}

// ===== woodsbury/jsonpointer benchmarks =====

func BenchmarkWoodsbury_Get_Root(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = woodsjp.Get("", smallData)
	}
}

func BenchmarkWoodsbury_Get_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = woodsjp.Get("/name", smallData)
	}
}

func BenchmarkWoodsbury_Get_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = woodsjp.Get("/profile/settings/theme", smallData)
	}
}

func BenchmarkWoodsbury_Medium_GetUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = woodsjp.Get("/users/50/name", mediumData)
	}
}

func BenchmarkWoodsbury_NotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = woodsjp.Get("/nonexistent", smallData)
	}
}

// ===== dolmen-go/jsonptr benchmarks =====

func BenchmarkDolmen_Get_Root(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dolmenjp.Get(smallData, "")
	}
}

func BenchmarkDolmen_Get_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dolmenjp.Get(smallData, "/name")
	}
}

func BenchmarkDolmen_Get_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dolmenjp.Get(smallData, "/profile/settings/theme")
	}
}

func BenchmarkDolmen_Medium_GetUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dolmenjp.Get(mediumData, "/users/50/name")
	}
}

func BenchmarkDolmen_NotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = dolmenjp.Get(smallData, "/nonexistent")
	}
}

// ===== Struct vs Map performance comparisons =====

// Struct access benchmarks
func BenchmarkOur_Struct_Get_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(structData, ourjp.Path{"name"})
	}
}

func BenchmarkOur_Struct_Get_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(structData, ourjp.Path{"profile", "settings", "theme"})
	}
}

func BenchmarkOur_Struct_Find_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/name", structData)
	}
}

func BenchmarkOur_Struct_Find_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/profile/settings/theme", structData)
	}
}

// Map access benchmarks (for comparison)
func BenchmarkOur_Map_Get_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(mapData, ourjp.Path{"name"})
	}
}

func BenchmarkOur_Map_Get_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(mapData, ourjp.Path{"profile", "settings", "theme"})
	}
}

func BenchmarkOur_Map_Find_Shallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/name", mapData)
	}
}

func BenchmarkOur_Map_Find_Deep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ourjp.FindByPointer("/profile/settings/theme", mapData)
	}
}

// Mixed data benchmark (struct containing maps)
func BenchmarkOur_Mixed_StructWithMap(b *testing.B) {
	mixed := map[string]any{
		"user": structData,
		"meta": map[string]any{"version": "1.0"},
	}
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(mixed, ourjp.Path{"user", "profile", "settings", "theme"})
	}
}

// Test struct field caching effectiveness
func BenchmarkOur_Struct_FieldCaching(b *testing.B) {
	// This benchmark tests if repeated access to the same struct type
	// benefits from field caching
	users := make([]BenchUser, 10)
	for i := range users {
		users[i] = generateStructData()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range users {
			_ = ourjp.Get(users[j], ourjp.Path{"name"})
		}
	}
}

// JSON tag vs field name access
func BenchmarkOur_Struct_JSONTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(structData, ourjp.Path{"name"}) // Uses JSON tag
	}
}

func BenchmarkOur_Struct_FieldName(b *testing.B) {
	type TestStruct struct {
		Name string // No JSON tag, uses field name
	}
	data := TestStruct{Name: "Alice"}

	for i := 0; i < b.N; i++ {
		_ = ourjp.Get(data, ourjp.Path{"Name"}) // Uses field name
	}
}
