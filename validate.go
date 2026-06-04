package jsonpointer

const (
	// MaxPointerLength is the maximum allowed length for JSON Pointer strings.
	MaxPointerLength = 1024

	// MaxPathLength is the maximum allowed number of pointer tokens.
	MaxPathLength = 256
)

func validatePointerString(pointer string) error {
	if len(pointer) > MaxPointerLength {
		return ErrPointerTooLong
	}
	if pointer == "" {
		return nil
	}
	if pointer[0] != '/' {
		return ErrInvalidPointer
	}
	for i := 0; i < len(pointer); i++ {
		if pointer[i] != '~' {
			continue
		}
		if i+1 >= len(pointer) {
			return ErrInvalidPointer
		}
		next := pointer[i+1]
		if next != '0' && next != '1' {
			return ErrInvalidPointer
		}
		i++
	}
	return nil
}

func validateTokens(tokens []string) error {
	if len(tokens) > MaxPathLength {
		return ErrPathTooLong
	}
	if pointerLength(tokens) > MaxPointerLength {
		return ErrPointerTooLong
	}
	return nil
}
