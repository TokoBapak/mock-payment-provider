package presentation

type RequestValidationError struct {
	Err error
}

func (r *RequestValidationError) Error() string {
	return r.Err.Error()
}
