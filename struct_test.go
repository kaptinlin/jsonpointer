package jsonpointer

import (
	"errors"
	"reflect"
	"testing"
)

// Test structs
type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string
	private string // private field, should be ignored
	Ignored string `json:"-"` // ignored field
}

type Profile struct {
	User     User   `json:"user"`
	Location string `json:"location"`
}

type EmptyTagName struct {
	Name string `json:",omitempty"`
}

type EmbeddedName struct {
	EmbeddedName string `json:"embedded_name"`
}

type EmbeddedUser struct {
	EmbeddedName
	Name string `json:"name"`
}

func TestStructField(t *testing.T) {
	t.Parallel()

	user := User{
		Name:    "Alice",
		Age:     30,
		Email:   "alice@example.com",
		private: "secret",
		Ignored: "ignored",
	}

	tests := []struct {
		name     string
		field    string
		expected any
		found    bool
	}{
		{"JSON tag field", "name", "Alice", true},
		{"JSON tag field age", "age", 30, true},
		{"Regular field", "Email", "alice@example.com", true},
		{"Private field", "private", nil, false},
		{"Ignored field", "Ignored", nil, false},
		{"Nonexistent field", "nonexistent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := reflect.ValueOf(user)
			found := structField(tt.field, &value)

			if found != tt.found {
				t.Errorf("structField() found = %v, want %v", found, tt.found)
			}

			if found && value.Interface() != tt.expected {
				t.Errorf("structField() value = %v, want %v", value.Interface(), tt.expected)
			}
		})
	}
}

func TestStructFieldWithPointer(t *testing.T) {
	t.Parallel()

	user := &User{
		Name:  "Bob",
		Age:   25,
		Email: "bob@example.com",
	}

	value := reflect.ValueOf(user)
	found := structField("name", &value)

	if !found {
		t.Error("structField() should be able to find field in pointer to struct")
	}

	if value.Interface() != "Bob" {
		t.Errorf("structField() value = %v, want %v", value.Interface(), "Bob")
	}
}

func TestStructFieldWithEmptyJSONTagName(t *testing.T) {
	t.Parallel()

	value := reflect.ValueOf(EmptyTagName{Name: "fallback"})
	found := structField("Name", &value)

	if !found {
		t.Fatal("structField() should fall back to the field name when the JSON tag name is empty")
	}

	if value.Interface() != "fallback" {
		t.Errorf("structField() value = %v, want %v", value.Interface(), "fallback")
	}
}

func TestStructFieldWithEmbeddedFieldOrder(t *testing.T) {
	t.Parallel()

	value := reflect.ValueOf(EmbeddedUser{
		EmbeddedName: EmbeddedName{EmbeddedName: "embedded"},
		Name:         "top-level",
	})
	found := structField("name", &value)

	if !found {
		t.Fatal("structField() should find the top-level tagged field")
	}

	if value.Interface() != "top-level" {
		t.Errorf("structField() value = %v, want %v", value.Interface(), "top-level")
	}
}

func TestGetWithEmptyJSONTagName(t *testing.T) {
	t.Parallel()

	result, err := Get(EmptyTagName{Name: "fallback"}, "Name")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result != "fallback" {
		t.Errorf("Get() = %v, want %v", result, "fallback")
	}
}

func TestGetWithStruct(t *testing.T) {
	t.Parallel()

	user := User{
		Name:  "Charlie",
		Age:   35,
		Email: "charlie@example.com",
	}

	tests := []struct {
		name    string
		path    Path
		want    any
		wantErr error
	}{
		{name: "Get name via JSON tag", path: Path{"name"}, want: "Charlie"},
		{name: "Get age via JSON tag", path: Path{"age"}, want: 35},
		{name: "Get email via field name", path: Path{"Email"}, want: "charlie@example.com"},
		{name: "Get nonexistent field", path: Path{"nonexistent"}, wantErr: ErrFieldNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get(user, tt.path...)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindWithStruct(t *testing.T) {
	t.Parallel()

	user := User{
		Name:  "David",
		Age:   40,
		Email: "david@example.com",
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		{"Find name via JSON tag", Path{"name"}, "David"},
		{"Find age via JSON tag", Path{"age"}, 40},
		{"Find email via field name", Path{"Email"}, "david@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ref, err := Find(user, tt.path...)
			if err != nil {
				t.Errorf("Find() error = %v", err)
				return
			}
			if ref.Val != tt.expected {
				t.Errorf("Find() = %v, want %v", ref.Val, tt.expected)
			}
		})
	}
}

func TestNestedStruct(t *testing.T) {
	t.Parallel()

	profile := Profile{
		User: User{
			Name:  "Eve",
			Age:   28,
			Email: "eve@example.com",
		},
		Location: "New York",
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		{"Get user object", Path{"user"}, profile.User},
		{"Get nested user name", Path{"user", "name"}, "Eve"},
		{"Get nested user age", Path{"user", "age"}, 28},
		{"Get location", Path{"location"}, "New York"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Get(profile, tt.path...)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMixedMapAndStruct(t *testing.T) {
	t.Parallel()

	// Test mixed usage of map and struct
	data := map[string]any{
		"user": User{
			Name:  "Frank",
			Age:   45,
			Email: "frank@example.com",
		},
		"meta": map[string]any{
			"version": "1.0",
		},
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		{"Get user from map", Path{"user"}, data["user"]},
		{"Get user name from struct in map", Path{"user", "name"}, "Frank"},
		{"Get meta version", Path{"meta", "version"}, "1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Get(data, tt.path...)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindByPointerWithStruct(t *testing.T) {
	t.Parallel()

	user := User{
		Name:  "Grace",
		Age:   32,
		Email: "grace@example.com",
	}

	tests := []struct {
		name     string
		pointer  string
		expected any
	}{
		{"Find name via JSON tag", "/name", "Grace"},
		{"Find age via JSON tag", "/age", 32},
		{"Find email via field name", "/Email", "grace@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ref, err := FindByPointer(user, tt.pointer)
			if err != nil {
				t.Errorf("FindByPointer() error = %v", err)
				return
			}
			if ref.Val != tt.expected {
				t.Errorf("FindByPointer() = %v, want %v", ref.Val, tt.expected)
			}
		})
	}
}

func TestFindByPointerNestedStruct(t *testing.T) {
	t.Parallel()

	profile := Profile{
		User: User{
			Name:  "Henry",
			Age:   27,
			Email: "henry@example.com",
		},
		Location: "London",
	}

	tests := []struct {
		name     string
		pointer  string
		expected any
	}{
		{"Find nested user name", "/user/name", "Henry"},
		{"Find nested user age", "/user/age", 27},
		{"Find location", "/location", "London"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ref, err := FindByPointer(profile, tt.pointer)
			if err != nil {
				t.Errorf("FindByPointer() error = %v", err)
				return
			}
			if ref.Val != tt.expected {
				t.Errorf("FindByPointer() = %v, want %v", ref.Val, tt.expected)
			}
		})
	}
}

// Test pointer to struct support across all API functions
func TestPointerToStruct(t *testing.T) {
	t.Parallel()

	user := &User{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Test Get with pointer to struct
	t.Run("Get with pointer to struct", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			path     Path
			expected any
		}{
			{"Get name via JSON tag", Path{"name"}, "Alice"},
			{"Get age via JSON tag", Path{"age"}, 30},
			{"Get email via field name", Path{"Email"}, "alice@example.com"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result, err := Get(user, tt.path...)
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	// Test Find with pointer to struct
	t.Run("Find with pointer to struct", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			path     Path
			expected any
		}{
			{"Find name via JSON tag", Path{"name"}, "Alice"},
			{"Find age via JSON tag", Path{"age"}, 30},
			{"Find email via field name", Path{"Email"}, "alice@example.com"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ref, err := Find(user, tt.path...)
				if err != nil {
					t.Errorf("Find() error = %v", err)
					return
				}
				if ref.Val != tt.expected {
					t.Errorf("Find() = %v, want %v", ref.Val, tt.expected)
				}
			})
		}
	})

	// Test FindByPointer with pointer to struct
	t.Run("FindByPointer with pointer to struct", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			pointer  string
			expected any
		}{
			{"Find name via JSON tag", "/name", "Alice"},
			{"Find age via JSON tag", "/age", 30},
			{"Find email via field name", "/Email", "alice@example.com"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ref, err := FindByPointer(user, tt.pointer)
				if err != nil {
					t.Errorf("FindByPointer() error = %v", err)
					return
				}
				if ref.Val != tt.expected {
					t.Errorf("FindByPointer() = %v, want %v", ref.Val, tt.expected)
				}
			})
		}
	})
}

// Test nested pointer to struct
func TestNestedPointerToStruct(t *testing.T) {
	t.Parallel()

	profile := &Profile{
		User: User{
			Name:  "Bob",
			Age:   25,
			Email: "bob@example.com",
		},
		Location: "Tokyo",
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		{"Get user from pointer to profile", Path{"user"}, profile.User},
		{"Get nested user name", Path{"user", "name"}, "Bob"},
		{"Get location from pointer to profile", Path{"location"}, "Tokyo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Get(profile, tt.path...)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test FindByPointer with nested pointer
	ref, err := FindByPointer(profile, "/user/name")
	if err != nil {
		t.Errorf("FindByPointer() error = %v", err)
		return
	}
	if ref.Val != "Bob" {
		t.Errorf("FindByPointer() = %v, want %v", ref.Val, "Bob")
	}
}

// Test multiple levels of pointers
func TestMultipleLevelsPointers(t *testing.T) {
	t.Parallel()

	user := &User{
		Name:  "Charlie",
		Age:   35,
		Email: "charlie@example.com",
	}

	// Pointer to pointer to struct
	userPtr := &user

	// This should still work by dereferencing all pointers
	name, err := Get(userPtr, "name")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if name != "Charlie" {
		t.Errorf("Get() with double pointer = %v, want %v", name, "Charlie")
	}
}

// Test comprehensive mixed struct and map data scenarios
func TestMixedStructMapComprehensive(t *testing.T) {
	t.Parallel()

	// Complex nested structure with mixed types
	type Company struct {
		Name      string         `json:"name"`
		Founded   int            `json:"founded"`
		Employees []User         `json:"employees"`
		Metadata  map[string]any `json:"metadata"`
		Locations map[string]any `json:"locations"`
	}

	company := Company{
		Name:    "Tech Corp",
		Founded: 2020,
		Employees: []User{
			{Name: "Alice", Age: 30, Email: "alice@techcorp.com"},
			{Name: "Bob", Age: 25, Email: "bob@techcorp.com"},
		},
		Metadata: map[string]any{
			"industry": "Technology",
			"size":     "Medium",
			"public":   true,
		},
		Locations: map[string]any{
			"headquarters": map[string]any{
				"city":    "San Francisco",
				"country": "USA",
			},
			"branch": map[string]any{
				"city":    "New York",
				"country": "USA",
			},
		},
	}

	// Test struct containing arrays of structs
	t.Run("Struct containing arrays of structs", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			path     Path
			expected any
		}{
			{"Company name", Path{"name"}, "Tech Corp"},
			{"First employee name", Path{"employees", "0", "name"}, "Alice"},
			{"Second employee email", Path{"employees", "1", "Email"}, "bob@techcorp.com"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result, err := Get(company, tt.path...)
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}

		// Test array access separately (can't compare slices directly)
		t.Run("Employee array access", func(t *testing.T) {
			t.Parallel()

			employees, err := Get(company, "employees")
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if employees == nil {
				t.Error("Get() employees should not be nil")
				return
			}

			// Verify it's the correct type and length
			if emp, ok := employees.([]User); ok {
				if len(emp) != 2 {
					t.Errorf("Expected 2 employees, got %d", len(emp))
				}
				if emp[0].Name != "Alice" {
					t.Errorf("First employee name = %v, want Alice", emp[0].Name)
				}
			} else {
				t.Errorf("Expected []User, got %T", employees)
			}
		})
	})

	// Test struct containing maps
	t.Run("Struct containing maps", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			path     Path
			expected any
		}{
			{"Metadata industry", Path{"metadata", "industry"}, "Technology"},
			{"Metadata public", Path{"metadata", "public"}, true},
			{"HQ city", Path{"locations", "headquarters", "city"}, "San Francisco"},
			{"Branch country", Path{"locations", "branch", "country"}, "USA"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result, err := Get(company, tt.path...)
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	// Test with FindByPointer
	t.Run("FindByPointer with complex mixed data", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			pointer  string
			expected any
		}{
			{"Employee via JSON pointer", "/employees/0/name", "Alice"},
			{"Metadata via JSON pointer", "/metadata/industry", "Technology"},
			{"Nested map via JSON pointer", "/locations/headquarters/city", "San Francisco"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ref, err := FindByPointer(company, tt.pointer)
				if err != nil {
					t.Errorf("FindByPointer() error = %v", err)
					return
				}
				if ref.Val != tt.expected {
					t.Errorf("FindByPointer() = %v, want %v", ref.Val, tt.expected)
				}
			})
		}
	})
}

// Test maps containing structs and complex nesting
func TestMapContainingStructs(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"users": []User{
			{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
			{Name: "David", Age: 40, Email: "david@example.com"},
		},
		"profiles": map[string]Profile{
			"admin": {
				User:     User{Name: "Admin", Age: 45, Email: "admin@example.com"},
				Location: "Server Room",
			},
			"guest": {
				User:     User{Name: "Guest", Age: 0, Email: "guest@example.com"},
				Location: "Lobby",
			},
		},
		"config": map[string]any{
			"database": map[string]any{
				"host": "localhost",
				"port": 5432,
				"settings": map[string]any{
					"ssl":     true,
					"timeout": 30,
				},
			},
			"features": []string{"auth", "logging", "metrics"},
		},
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		// Array of structs in map
		{"First user name", Path{"users", "0", "name"}, "Charlie"},
		{"Second user age", Path{"users", "1", "age"}, 40},

		// Map of structs
		{"Admin profile name", Path{"profiles", "admin", "user", "name"}, "Admin"},
		{"Guest location", Path{"profiles", "guest", "location"}, "Lobby"},

		// Deeply nested maps
		{"Database host", Path{"config", "database", "host"}, "localhost"},
		{"Database SSL setting", Path{"config", "database", "settings", "ssl"}, true},
		{"First feature", Path{"config", "features", "0"}, "auth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Get(data, tt.path...)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test with FindByPointer
	ref, err := FindByPointer(data, "/profiles/admin/user/Email")
	if err != nil {
		t.Errorf("FindByPointer() error = %v", err)
	} else if ref.Val != "admin@example.com" {
		t.Errorf("FindByPointer() = %v, want %v", ref.Val, "admin@example.com")
	}
}

// Test edge cases with mixed data
func TestMixedDataEdgeCases(t *testing.T) {
	t.Parallel()

	// Empty struct in map
	type Empty struct{}

	data := map[string]any{
		"empty_struct": Empty{},
		"nil_value":    nil,
		"user_ptr":     &User{Name: "Pointer User", Age: 30},
		"nested": map[string]any{
			"array_with_mixed": []any{
				User{Name: "Struct in array", Age: 20},
				map[string]any{"type": "map", "value": 42},
				"plain string",
				123,
			},
		},
	}

	tests := []struct {
		name     string
		path     Path
		expected any
	}{
		{"Empty struct", Path{"empty_struct"}, Empty{}},
		{"Nil value", Path{"nil_value"}, nil},
		{"Pointer to struct name", Path{"user_ptr", "name"}, "Pointer User"},
		{"Struct in mixed array", Path{"nested", "array_with_mixed", "0", "name"}, "Struct in array"},
		{"Map in mixed array", Path{"nested", "array_with_mixed", "1", "type"}, "map"},
		{"String in mixed array", Path{"nested", "array_with_mixed", "2"}, "plain string"},
		{"Number in mixed array", Path{"nested", "array_with_mixed", "3"}, 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Get(data, tt.path...)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test nil pointer handling
func TestNilPointerHandling(t *testing.T) {
	t.Parallel()

	var user *User

	result, err := Get(user, "name")
	if !errors.Is(err, ErrNilPointer) {
		t.Fatalf("Get() error = %v, want %v", err, ErrNilPointer)
	}
	if result != nil {
		t.Errorf("Get() with nil pointer = %v, want nil", result)
	}

	ref, err := FindByPointer(user, "/name")
	if !errors.Is(err, ErrNilPointer) {
		t.Fatalf("FindByPointer() error = %v, want %v", err, ErrNilPointer)
	}
	if ref != nil {
		t.Error("FindByPointer() should return nil reference for nil pointer")
	}
}

// Test Get function behavior with missing fields (should return error)
func TestGetMissingFieldBehavior(t *testing.T) {
	t.Parallel()

	user := User{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	profile := Profile{
		User:     user,
		Location: "New York",
	}

	tests := []struct {
		name    string
		data    any
		path    []string
		want    any
		wantErr error
	}{
		{name: "Missing field at end of path", data: user, path: []string{"nonexistent"}, wantErr: ErrFieldNotFound},
		{name: "Missing field in middle of path", data: user, path: []string{"nonexistent", "nested"}, wantErr: ErrFieldNotFound},
		{name: "Missing field in deeper nesting", data: user, path: []string{"nonexistent", "very", "deep", "path"}, wantErr: ErrFieldNotFound},
		{name: "Missing nested field in struct", data: profile, path: []string{"user", "nonexistent"}, wantErr: ErrFieldNotFound},
		{name: "Missing field with more path after", data: profile, path: []string{"user", "nonexistent", "more", "path"}, wantErr: ErrFieldNotFound},
		{name: "Valid field should still work", data: user, path: []string{"name"}, want: "Alice"},
		{name: "Valid nested field should still work", data: profile, path: []string{"user", "name"}, want: "Alice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get(tt.data, tt.path...)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test Get function with maps to ensure consistent behavior
func TestGetMissingMapKeyBehavior(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"user": map[string]any{
			"name": "Bob",
			"age":  25,
		},
		"config": map[string]any{
			"debug": true,
		},
	}

	tests := []struct {
		name    string
		path    []string
		want    any
		wantErr error
	}{
		{name: "Missing top-level key", path: []string{"missing"}, wantErr: ErrKeyNotFound},
		{name: "Missing key in middle of path", path: []string{"missing", "nested"}, wantErr: ErrKeyNotFound},
		{name: "Missing nested key", path: []string{"user", "missing"}, wantErr: ErrKeyNotFound},
		{name: "Missing nested key with more path", path: []string{"user", "missing", "more"}, wantErr: ErrKeyNotFound},
		{name: "Valid nested access should work", path: []string{"user", "name"}, want: "Bob"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get(data, tt.path...)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test mixed struct and map missing field behavior
func TestGetMissingFieldMixedData(t *testing.T) {
	t.Parallel()

	user := User{Name: "Charlie", Age: 35}
	data := map[string]any{
		"user":   user,
		"config": map[string]any{"debug": true},
	}

	tests := []struct {
		name    string
		path    []string
		want    any
		wantErr error
	}{
		{name: "Missing field in struct within map", path: []string{"user", "missing"}, wantErr: ErrFieldNotFound},
		{name: "Missing field in struct with more path", path: []string{"user", "missing", "deep"}, wantErr: ErrFieldNotFound},
		{name: "Missing key in map within map", path: []string{"config", "missing"}, wantErr: ErrKeyNotFound},
		{name: "Missing key in map with more path", path: []string{"config", "missing", "deep"}, wantErr: ErrKeyNotFound},
		{name: "Valid access should work", path: []string{"user", "name"}, want: "Charlie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get(data, tt.path...)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test array bounds checking behavior
func TestArrayBoundsChecking(t *testing.T) {
	t.Parallel()

	arr := []any{1, 2, 3} // length = 3

	tests := []struct {
		name    string
		index   string
		want    any
		wantErr error
	}{
		{name: "Valid index access", index: "1", want: 2},
		{name: "Array end marker", index: "-", wantErr: ErrIndexOutOfBounds},
		{name: "Index at array length", index: "3", wantErr: ErrIndexOutOfBounds},
		{name: "Index beyond array length", index: "5", wantErr: ErrIndexOutOfBounds},
		{name: "Much larger index", index: "100", wantErr: ErrIndexOutOfBounds},
	}

	t.Run("Get function array bounds", func(t *testing.T) {
		t.Parallel()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := Get(arr, tt.index)
				if tt.wantErr != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("Get() error = %v, want %v", err, tt.wantErr)
					}
					return
				}
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}
				if got != tt.want {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Find function array bounds", func(t *testing.T) {
		t.Parallel()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ref, err := Find(arr, tt.index)
				if tt.wantErr != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("Find() error = %v, want %v", err, tt.wantErr)
					}
					return
				}
				if err != nil {
					t.Fatalf("Find() error = %v", err)
				}
				if ref.Val != tt.want {
					t.Errorf("Find().Val = %v, want %v", ref.Val, tt.want)
				}
			})
		}
	})

	t.Run("FindByPointer function array bounds", func(t *testing.T) {
		t.Parallel()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ref, err := FindByPointer(arr, "/"+tt.index)
				if tt.wantErr != nil {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("FindByPointer() error = %v, want %v", err, tt.wantErr)
					}
					return
				}
				if err != nil {
					t.Fatalf("FindByPointer() error = %v", err)
				}
				if ref.Val != tt.want {
					t.Errorf("FindByPointer().Val = %v, want %v", ref.Val, tt.want)
				}
			})
		}
	})
}
