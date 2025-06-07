package main

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpointer"
)

func main() {
	fmt.Println("=== JSON Pointer Go Implementation Examples ===")
	fmt.Println()

	// Sample JSON document for demonstrations
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"id":   1,
				"name": "Alice Smith",
				"profile": map[string]any{
					"email": "alice@example.com",
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": true,
							"sms":   false,
						},
						"privacy": "public",
					},
				},
			},
			map[string]any{
				"id":   2,
				"name": "Bob Johnson",
				"profile": map[string]any{
					"email": "bob@example.com",
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": false,
							"sms":   true,
						},
						"privacy": "private",
					},
				},
			},
		},
		"metadata": map[string]any{
			"version":     "1.0.0",
			"created":     "2023-01-01",
			"description": "User management system",
		},
		"special~chars": map[string]any{
			"foo/bar": "value with special characters",
			"tilde~":  "tilde value",
		},
	}

	fmt.Printf("Sample document: %+v\n", doc)
	fmt.Println()

	// Example 1: Basic Get Operations (never fail)
	fmt.Println("=== 1. Get Operations (never fail, returns nil if not found) ===")

	// Get root
	result := jsonpointer.Get(doc, jsonpointer.Path{})
	fmt.Printf("Get root: %+v\n", result)

	// Get simple property
	result = jsonpointer.Get(doc, jsonpointer.Path{"metadata", "version"})
	fmt.Printf("Get /metadata/version: %v\n", result)

	// Get array element
	result = jsonpointer.Get(doc, jsonpointer.Path{"users", 0, "name"})
	fmt.Printf("Get /users/0/name: %v\n", result)

	// Get deep nested property
	result = jsonpointer.Get(doc, jsonpointer.Path{"users", 1, "profile", "settings", "notifications", "sms"})
	fmt.Printf("Get /users/1/profile/settings/notifications/sms: %v\n", result)

	// Get non-existent property (returns nil)
	result = jsonpointer.Get(doc, jsonpointer.Path{"nonexistent", "property"})
	fmt.Printf("Get non-existent property: %v\n", result)

	fmt.Println()

	// Example 2: Find Operations (return errors for invalid paths)
	fmt.Println("=== 2. Find Operations (return errors for invalid paths) ===")

	// Find successful
	ref, err := jsonpointer.Find(doc, jsonpointer.Path{"users", 0, "profile", "email"})
	if err != nil {
		log.Printf("Find error: %v", err)
	} else {
		fmt.Printf("Find result - Val: %v, Key: %v\n", ref.Val, ref.Key)
	}

	// Find with array end marker (creates reference to append position)
	ref, err = jsonpointer.Find(doc, jsonpointer.Path{"users", "-"})
	if err != nil {
		log.Printf("Find error: %v", err)
	} else {
		fmt.Printf("Find array end - Val: %v, Key: %v\n", ref.Val, ref.Key)
	}

	// Find with invalid array index (returns error)
	_, err = jsonpointer.Find(doc, jsonpointer.Path{"users", "invalid_index"})
	if err != nil {
		fmt.Printf("Find invalid array index error: %v\n", err)
	}

	// Find non-existent path (returns error)
	_, err = jsonpointer.Find(doc, jsonpointer.Path{"nonexistent", "deeply", "nested"})
	if err != nil {
		fmt.Printf("Find non-existent path error: %v\n", err)
	}

	fmt.Println()

	// Example 3: FindByPointer - Optimized string-based operations
	fmt.Println("=== 3. FindByPointer - Optimized string operations ===")

	// FindByPointer successful
	ref, err = jsonpointer.FindByPointer("/users/0/profile/settings/privacy", doc)
	if err != nil {
		log.Printf("FindByPointer error: %v", err)
	} else {
		fmt.Printf("FindByPointer result - Val: %v, Key: %v\n", ref.Val, ref.Key)
	}

	// FindByPointer with escaped characters
	ref, err = jsonpointer.FindByPointer("/special~0chars/foo~1bar", doc)
	if err != nil {
		log.Printf("FindByPointer error: %v", err)
	} else {
		fmt.Printf("FindByPointer escaped chars - Val: %v, Key: %v\n", ref.Val, ref.Key)
	}

	// FindByPointer with array end marker
	ref, err = jsonpointer.FindByPointer("/users/-", doc)
	if err != nil {
		log.Printf("FindByPointer error: %v", err)
	} else {
		fmt.Printf("FindByPointer array end - Val: %v, Key: %v\n", ref.Val, ref.Key)
	}

	fmt.Println()

	// Example 4: JSON Pointer Parsing and Formatting
	fmt.Println("=== 4. JSON Pointer Parsing and Formatting ===")

	// Parse JSON pointers to paths
	pointers := []string{
		"",
		"/users",
		"/users/0/name",
		"/metadata/version",
		"/special~0chars/foo~1bar",
		"/users/-",
	}

	for _, pointer := range pointers {
		path := jsonpointer.ParseJsonPointer(pointer)
		formatted := jsonpointer.FormatJsonPointer(path)
		fmt.Printf("Parse '%s' -> %+v -> Format '%s'\n", pointer, path, formatted)
	}

	fmt.Println()

	// Example 5: Component Escaping and Unescaping
	fmt.Println("=== 5. Component Escaping and Unescaping ===")

	components := []string{
		"simple",
		"foo/bar",
		"foo~bar",
		"foo~bar/baz",
		"~tilde~/slash/",
	}

	for _, component := range components {
		escaped := jsonpointer.EscapeComponent(component)
		unescaped := jsonpointer.UnescapeComponent(escaped)
		fmt.Printf("Component '%s' -> Escape '%s' -> Unescape '%s'\n", component, escaped, unescaped)
	}

	fmt.Println()

	// Example 6: Utility Functions
	fmt.Println("=== 6. Utility Functions ===")

	// Path utilities
	path1 := jsonpointer.Path{"users", 0}
	path2 := jsonpointer.Path{"users", 0, "profile"}
	path3 := jsonpointer.Path{"metadata"}

	fmt.Printf("IsRoot(%+v): %v\n", jsonpointer.Path{}, jsonpointer.IsRoot(jsonpointer.Path{}))
	fmt.Printf("IsRoot(%+v): %v\n", path1, jsonpointer.IsRoot(path1))

	fmt.Printf("IsChild(%+v, %+v): %v\n", path1, path2, jsonpointer.IsChild(path1, path2))
	fmt.Printf("IsChild(%+v, %+v): %v\n", path2, path1, jsonpointer.IsChild(path2, path1))

	fmt.Printf("IsPathEqual(%+v, %+v): %v\n", path1, path1, jsonpointer.IsPathEqual(path1, path1))
	fmt.Printf("IsPathEqual(%+v, %+v): %v\n", path1, path3, jsonpointer.IsPathEqual(path1, path3))

	// Parent operations
	parent, err := jsonpointer.Parent(path2)
	if err != nil {
		log.Printf("Parent error: %v", err)
	} else {
		fmt.Printf("Parent of %+v: %+v\n", path2, parent)
	}

	// Parent of root (error)
	_, err = jsonpointer.Parent(jsonpointer.Path{})
	if err != nil {
		fmt.Printf("Parent of root error: %v\n", err)
	}

	// Array index validation
	fmt.Printf("IsValidIndex('0'): %v\n", jsonpointer.IsValidIndex("0"))
	fmt.Printf("IsValidIndex('123'): %v\n", jsonpointer.IsValidIndex("123"))
	fmt.Printf("IsValidIndex('01'): %v\n", jsonpointer.IsValidIndex("01"))   // false - leading zero
	fmt.Printf("IsValidIndex('-'): %v\n", jsonpointer.IsValidIndex("-"))     // false - special marker
	fmt.Printf("IsValidIndex('abc'): %v\n", jsonpointer.IsValidIndex("abc")) // false - not a number
	fmt.Printf("IsValidIndex(42): %v\n", jsonpointer.IsValidIndex(42))       // true - number

	// Integer checking
	fmt.Printf("IsInteger('123'): %v\n", jsonpointer.IsInteger("123"))
	fmt.Printf("IsInteger('01'): %v\n", jsonpointer.IsInteger("01")) // true - has digits
	fmt.Printf("IsInteger('abc'): %v\n", jsonpointer.IsInteger("abc"))

	fmt.Println()

	// Example 7: Type Guards and Reference Types
	fmt.Println("=== 7. Type Guards and Reference Analysis ===")

	// Get references to different types
	arrayRef, _ := jsonpointer.Find(doc, jsonpointer.Path{"users", 0})
	objectRef, _ := jsonpointer.Find(doc, jsonpointer.Path{"metadata", "version"})

	fmt.Printf("Array reference: %+v\n", arrayRef)
	fmt.Printf("IsArrayReference: %v\n", jsonpointer.IsArrayReference(*arrayRef))
	fmt.Printf("IsObjectReference: %v\n", jsonpointer.IsObjectReference(*arrayRef))

	fmt.Printf("Object reference: %+v\n", objectRef)
	fmt.Printf("IsArrayReference: %v\n", jsonpointer.IsArrayReference(*objectRef))
	fmt.Printf("IsObjectReference: %v\n", jsonpointer.IsObjectReference(*objectRef))

	fmt.Println()

	// Example 8: Validation
	fmt.Println("=== 8. Validation ===")

	// Valid pointers
	validPointers := []string{
		"",
		"/users",
		"/users/0/name",
		"/special~0chars/foo~1bar",
	}

	// Invalid pointers
	invalidPointers := []string{
		"users",           // missing leading slash
		"/users/",         // trailing slash
		"/invalid~escape", // invalid escape sequence
	}

	fmt.Println("Valid pointers:")
	for _, pointer := range validPointers {
		err := jsonpointer.ValidateJsonPointer(pointer)
		fmt.Printf("  '%s': %v\n", pointer, err)
	}

	fmt.Println("Invalid pointers:")
	for _, pointer := range invalidPointers {
		err := jsonpointer.ValidateJsonPointer(pointer)
		fmt.Printf("  '%s': %v\n", pointer, err)
	}

	// Path validation
	validPath := jsonpointer.Path{"users", 0, "name"}
	err = jsonpointer.ValidatePath(validPath)
	fmt.Printf("Valid path %+v: %v\n", validPath, err)

	// Create a very long path (over limit)
	longPath := make(jsonpointer.Path, 300) // Over 256 limit
	for i := range longPath {
		longPath[i] = fmt.Sprintf("step%d", i)
	}
	err = jsonpointer.ValidatePath(longPath)
	fmt.Printf("Long path (300 steps): %v\n", err)

	fmt.Println()

	// Example 9: ToPath Utility
	fmt.Println("=== 9. ToPath Utility - Convert string or path to normalized path ===")

	// Convert string pointer to path
	stringPointer := "/users/0/profile/email"
	pathFromString := jsonpointer.ToPath(stringPointer)
	fmt.Printf("ToPath('%s'): %+v\n", stringPointer, pathFromString)

	// Convert existing path (no-op)
	existingPath := jsonpointer.Path{"metadata", "version"}
	pathFromPath := jsonpointer.ToPath(existingPath)
	fmt.Printf("ToPath(%+v): %+v\n", existingPath, pathFromPath)

	fmt.Println()
	fmt.Println("=== All examples completed successfully! ===")
}
