package model

import "fmt"

type ValidationErrors struct {
	Errors []error `json:"errors"`
}

func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("Multiple errors %s", e.Errors)
}

func (e *ValidationErrors) AddError(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

type ValidationError struct {
	Field string `json:"field"`
	Message  string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation Error %s, %s",e.Field, e.Message)
}

type TokenExpiredError struct {
	Type 		string `json:"error"`
	Description	string `json:"description"`
}

func (e *TokenExpiredError) Error() string {
	return fmt.Sprintf("Token Expired Error %s, %s",e.Type, e.Description)
}

type TokenMissingError struct {
	Type 		string `json:"error"`
	Description	string `json:"description"`
}

func (e *TokenMissingError) Error() string {
	return fmt.Sprintf("Token Missing Error %s, %s",e.Type, e.Description)
}