module jsonpointer-compare

go 1.26.3

// Local library.
replace github.com/kaptinlin/jsonpointer => ../

require (
	github.com/bragdond/jsonpointer-go v1.0.0
	github.com/dolmen-go/jsonptr v0.0.0-20260529085001-d6b11e72da90
	github.com/go-json-experiment/json v0.0.0-20260601182631-00ed12fed2a6
	github.com/go-openapi/jsonpointer v0.23.1
	github.com/kaptinlin/jsonpointer v0.4.18
	github.com/woodsbury/jsonpointer v0.7.2
)

require github.com/go-openapi/swag/jsonname v0.26.0 // indirect
