package main

// ValidationError is ...
type ValidationError struct {
	Field     string      `json:"field"`
	Tag       string      `json:"tag"`
	TagTarget string      `json:"tagTarget"`
	Value     interface{} `json:"value"`
}

// ResponseValidationError is ...
type ResponseValidationError struct {
	Errors []ValidationError `json:"errors"`
}

// ResponseUserID is ...
type ResponseUserID struct {
	ID int64 `json:"id"`
}
