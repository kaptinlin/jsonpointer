package jsonpointer

func find(val any, path Path) (*Reference, error) {
	if len(path) == 0 {
		return &Reference{Val: val}, nil
	}

	var obj any
	var key string
	current := val

	for _, step := range path {
		obj = current
		key = step

		if current == nil {
			return nil, ErrNotFound
		}

		result, handled, err := tryArrayAccess(current, step)
		if err != nil {
			return nil, err
		}
		if handled {
			current = result
			continue
		}

		result, handled, err = tryObjectAccess(current, step)
		if err != nil {
			return nil, err
		}
		if handled {
			current = result
			continue
		}

		return nil, ErrNotFound
	}

	return &Reference{Val: current, Obj: obj, Key: key}, nil
}
