module jsonpointer-compare

go 1.26.2

// 本地库
replace github.com/kaptinlin/jsonpointer => ../

require (
	github.com/bragdond/jsonpointer-go v1.0.0
	github.com/dolmen-go/jsonptr v0.0.0-20240328010033-38530b85cd9c
	github.com/go-json-experiment/json v0.0.0-20251027170946-4849db3c2f7e
	github.com/go-openapi/jsonpointer v0.22.4
	github.com/kaptinlin/jsonpointer v0.4.18
	github.com/woodsbury/jsonpointer v0.7.1
)

require github.com/go-openapi/swag/jsonname v0.25.4 // indirect
