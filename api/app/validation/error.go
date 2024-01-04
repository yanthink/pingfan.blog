package validation

type Error struct {
	Message string
	Errors  map[string][]string
	Err     error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	for _, v := range e.Errors {
		if len(v) > 0 {
			return v[0]
		}
	}

	return e.Err.Error()
}
