package jsonpointer

import (
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

func TestStructField(t *testing.T) {
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

func TestGetWithStruct(t *testing.T) {
	user := User{
		Name:  "Charlie",
		Age:   35,
		Email: "charlie@example.com",
	}

	tests := []struct {
		name        string
		path        Path
		expected    any
		expectError bool
	}{
		{"Get name via JSON tag", Path{"name"}, "Charlie", false},
		{"Get age via JSON tag", Path{"age"}, 35, false},
		{"Get email via field name", Path{"Email"}, "charlie@example.com", false},
		{"Get nonexistent field", Path{"nonexistent"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Get(user, tt.path...)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for nonexistent field, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindWithStruct(t *testing.T) {
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
			result, _ := Get(profile, tt.path...)
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMixedMapAndStruct(t *testing.T) {
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
			result, _ := Get(data, tt.path...)
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindByPointerWithStruct(t *testing.T) {
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
	user := &User{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	// Test Get with pointer to struct
	t.Run("Get with pointer to struct", func(t *testing.T) {
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
				result, _ := Get(user, tt.path...)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	// Test Find with pointer to struct
	t.Run("Find with pointer to struct", func(t *testing.T) {
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
			result, _ := Get(profile, tt.path...)
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
	user := &User{
		Name:  "Charlie",
		Age:   35,
		Email: "charlie@example.com",
	}

	// Pointer to pointer to struct
	userPtr := &user

	// This should still work by dereferencing all pointers
	name, _ := Get(userPtr, "name")
	if name != "Charlie" {
		t.Errorf("Get() with double pointer = %v, want %v", name, "Charlie")
	}
}

// Test comprehensive mixed struct and map data scenarios
func TestMixedStructMapComprehensive(t *testing.T) {
	// Complex nested structure with mixed types
	type Company struct {
		Name      string                 `json:"name"`
		Founded   int                    `json:"founded"`
		Employees []User                 `json:"employees"`
		Metadata  map[string]any         `json:"metadata"`
		Locations map[string]interface{} `json:"locations"`
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
		Locations: map[string]interface{}{
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
				result, _ := Get(company, tt.path...)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}

		// Test array access separately (can't compare slices directly)
		t.Run("Employee array access", func(t *testing.T) {
			employees, _ := Get(company, "employees")
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
				result, _ := Get(company, tt.path...)
				if result != tt.expected {
					t.Errorf("Get() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	// Test with FindByPointer
	t.Run("FindByPointer with complex mixed data", func(t *testing.T) {
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
			result, _ := Get(data, tt.path...)
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
			result, _ := Get(data, tt.path...)
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test nil pointer handling
func TestNilPointerHandling(t *testing.T) {
	var user *User = nil

	// Should return error when trying to access field of nil pointer
	result, err := Get(user, "name")
	if result != nil {
		t.Errorf("Get() with nil pointer = %v, want nil", result)
	}
	// Now we expect an error when trying to access fields of nil pointer
	if err == nil {
		t.Error("Get() with nil pointer should return error")
	}

	// FindByPointer should also return error for nil pointer field access
	ref, err := FindByPointer(user, "/name")
	if err == nil {
		t.Error("FindByPointer() with nil pointer should return error")
	}
	if ref != nil {
		t.Error("FindByPointer() should return nil reference for nil pointer")
	}
}

// Test Get function behavior with missing fields (should return error)
func TestGetMissingFieldBehavior(t *testing.T) {
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
		name          string
		data          any
		path          []string
		expectedValue any
		expectedError bool
		description   string
	}{
		{
			name:          "Missing field at end of path",
			data:          user,
			path:          []string{"nonexistent"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing struct field",
		},
		{
			name:          "Missing field in middle of path",
			data:          user,
			path:          []string{"nonexistent", "nested"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error when missing field is in middle of path",
		},
		{
			name:          "Missing field in deeper nesting",
			data:          user,
			path:          []string{"nonexistent", "very", "deep", "path"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for deeply nested missing fields",
		},
		{
			name:          "Missing nested field in struct",
			data:          profile,
			path:          []string{"user", "nonexistent"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing field in nested struct",
		},
		{
			name:          "Missing field with more path after",
			data:          profile,
			path:          []string{"user", "nonexistent", "more", "path"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error when missing field has more path after it",
		},
		{
			name:          "Valid field should still work",
			data:          user,
			path:          []string{"name"},
			expectedValue: "Alice",
			expectedError: false,
			description:   "Valid fields should continue to work normally",
		},
		{
			name:          "Valid nested field should still work",
			data:          profile,
			path:          []string{"user", "name"},
			expectedValue: "Alice",
			expectedError: false,
			description:   "Valid nested fields should continue to work normally",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Get(tt.data, tt.path...)

			// Check error expectation
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none. %s", tt.description)
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v. %s", err, tt.description)
			}

			// Check result
			if result != tt.expectedValue {
				t.Errorf("Get() = %v, want %v. %s", result, tt.expectedValue, tt.description)
			}
		})
	}
}

// Test Get function with maps to ensure consistent behavior
func TestGetMissingMapKeyBehavior(t *testing.T) {
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
		name          string
		path          []string
		expectedValue any
		expectedError bool
		description   string
	}{
		{
			name:          "Missing top-level key",
			path:          []string{"missing"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing top-level key",
		},
		{
			name:          "Missing key in middle of path",
			path:          []string{"missing", "nested"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error when missing key is in middle of path",
		},
		{
			name:          "Missing nested key",
			path:          []string{"user", "missing"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing nested key",
		},
		{
			name:          "Missing nested key with more path",
			path:          []string{"user", "missing", "more"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error when missing nested key has more path",
		},
		{
			name:          "Valid nested access should work",
			path:          []string{"user", "name"},
			expectedValue: "Bob",
			expectedError: false,
			description:   "Valid nested access should continue to work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Get(data, tt.path...)

			// Check error expectation
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none. %s", tt.description)
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v. %s", err, tt.description)
			}

			// Check result
			if result != tt.expectedValue {
				t.Errorf("Get() = %v, want %v. %s", result, tt.expectedValue, tt.description)
			}
		})
	}
}

// Test mixed struct and map missing field behavior
func TestGetMissingFieldMixedData(t *testing.T) {
	user := User{Name: "Charlie", Age: 35}
	data := map[string]any{
		"user":   user,
		"config": map[string]any{"debug": true},
	}

	tests := []struct {
		name          string
		path          []string
		expectedValue any
		expectedError bool
		description   string
	}{
		{
			name:          "Missing field in struct within map",
			path:          []string{"user", "missing"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing field in struct within map",
		},
		{
			name:          "Missing field in struct with more path",
			path:          []string{"user", "missing", "deep"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing field in struct with more path",
		},
		{
			name:          "Missing key in map within map",
			path:          []string{"config", "missing"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing key in nested map",
		},
		{
			name:          "Missing key in map with more path",
			path:          []string{"config", "missing", "deep"},
			expectedValue: nil,
			expectedError: true,
			description:   "Should return error for missing key in map with more path",
		},
		{
			name:          "Valid access should work",
			path:          []string{"user", "name"},
			expectedValue: "Charlie",
			expectedError: false,
			description:   "Valid access should continue to work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Get(data, tt.path...)

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none. %s", tt.description)
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error but got: %v. %s", err, tt.description)
			}

			if result != tt.expectedValue {
				t.Errorf("Get() = %v, want %v. %s", result, tt.expectedValue, tt.description)
			}
		})
	}
}
