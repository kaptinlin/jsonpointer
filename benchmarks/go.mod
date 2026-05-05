module jsonpointer-compare

go 1.26.2

// Local library.
replace github.com/kaptinlin/jsonpointer => ../

require (
	github.com/bragdond/jsonpointer-go v1.0.0
	github.com/dolmen-go/jsonptr v0.0.0-20240328010033-38530b85cd9c
	github.com/go-json-experiment/json v0.0.0-20260504200034-64a0a05799db
	github.com/go-openapi/jsonpointer v0.23.1
	github.com/kaptinlin/jsonpointer v0.4.18
	github.com/woodsbury/jsonpointer v0.7.2
)

require github.com/go-openapi/swag/jsonname v0.26.0 // indirect
