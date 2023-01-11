package schema

// Error is a struct that holds the common error data
// to be returned
type Error struct {
	StatusCode    int    `json:"status_code,string"`
	StatusMessage string `json:"status_message"`
}

type ValidationIssue struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationError is a specific struct that holds the request validation error data
// to be returned
type ValidationError struct {
	Error
	Issues []ValidationIssue `json:"issues"`
}
