package jsonpointer

func validatePointerString(pointer string) error {
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
