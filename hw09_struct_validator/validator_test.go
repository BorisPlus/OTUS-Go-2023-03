package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	LevelMonitoring struct {
		Samples []int `validate:"min:20|max:40"`
	}

	FakeValid struct {
		A int    `validate:"len:1"`
		B string `validate:"length:1"`
	}
)

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"User positive validation",
			User{
				ID:     "012345678901234567890123456789012345",
				Age:    27,
				Email:  "stuff@go.dev",
				Role:   "stuff",
				Phones: []string{"+1012345678", "+7701234568"},
			},
			nil,
		},
		{
			"App positive validation",
			App{
				Version: "v.1.9",
			},
			nil,
		},
		{
			"Token positive validation",
			Token{[]byte{}, []byte{}, []byte{}},
			nil,
		},
		{
			"Response positive validation",
			Response{
				Code: 200,
				Body: "OK HTTP/2",
			},
			nil,
		},
		{
			"LevelMonitoring positive validation",
			LevelMonitoring{
				Samples: []int{20, 35, 27, 40},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			// t.Parallel()
			validationErrors := Validate(tt.in)
			if !errors.Is(validationErrors, tt.expectedErr) {
				fmt.Println(validationErrors)
				t.Errorf("Error in %q case\n", tt.name)
			}
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			"User negative validation",
			User{
				ID:     "I have short id",
				Age:    14,
				Email:  "stuff@go.dev.ru",
				Role:   "guest",
				Phones: []string{"+10123456789", "+77012345689"},
			},
		},
		{
			"LevelMonitoring negative validation",
			LevelMonitoring{
				Samples: []int{20, 45, 27, 15},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			// t.Parallel()
			validationErrors := Validate(tt.in)
			if validationErrors == nil {
				t.Errorf("Error must be in %q case\n", tt.name)
			} else {
				fmt.Println("Test catch expected error. It's OK.")
				fmt.Println(validationErrors)
			}
		})
	}
}

func TestValidateNotStructObject(t *testing.T) {
	a := "I'm not a struct"
	validationErrors := Validate(a)
	if validationErrors != nil {
		fmt.Println("Test catch expected error. It's OK.")
		fmt.Println(validationErrors)
	} else {
		t.Error("Error must be in validation of not struct object")
	}
}

func TestValidateNotImplemented(t *testing.T) {
	fakeValid := FakeValid{
		A: 1,
		B: "fieldB",
	}
	validationErrors := Validate(fakeValid)
	if validationErrors == nil {
		t.Errorf("Error must be in FakeValid struct\n")
	} else {
		fmt.Println("Test catch expected error. It's OK.")
		fmt.Println(validationErrors)
	}
}

func TestValidateExpectedNotImplemented(t *testing.T) {
	fakeValid := FakeValid{
		A: 1,
		B: "fieldB",
	}
	expectedError := make(ValidationErrors, 0)
	// TODO: fmt.Errorf must be generate by pkg
	expectedError = append(
		expectedError,
		ValidationError{
			Field: "A",
			Err:   fmt.Errorf("validate function \"len\" for field type \"int\" is not implemented"),
		})
	expectedError = append(
		expectedError,
		ValidationError{
			Field: "B",
			Err:   fmt.Errorf("validate function \"length\" for field type \"string\" is not implemented"),
		})
	validationErrors := Validate(fakeValid)
	if validationErrors.Error() == expectedError.Error() {
		fmt.Println("Test catch expected error. It's OK.")
		fmt.Println(validationErrors)
		fmt.Println()
	} else {
		t.Error("Test catch unexpected error")
		fmt.Println()
		fmt.Println("Get error.")
		fmt.Println(validationErrors)
		fmt.Println()
		fmt.Println("Expected error.")
		fmt.Println(expectedError)
		fmt.Println()
	}
}

func TestValidateUnxpectedNotImplemented(t *testing.T) {
	fakeValid := FakeValid{
		A: 1,
		B: "fieldB",
	}
	unexpectedError := make(ValidationErrors, 0)
	unexpectedError = append(
		unexpectedError,
		ValidationError{
			Field: "B",
			Err:   fmt.Errorf("validate function \"length\" for field type \"string\" is not implemented"),
		})
	unexpectedError = append(
		unexpectedError,
		ValidationError{
			Field: "A",
			Err:   fmt.Errorf("validate function \"len\" for field type \"int\" is not implemented"),
		})
	validationErrors := Validate(fakeValid)
	if validationErrors.Error() == unexpectedError.Error() {
		t.Error("Error must be unexpected")
	} else {
		fmt.Println("Test catch expected fake error. It's OK.")
		fmt.Println()
		fmt.Println("Get error. It's OK.")
		fmt.Println(validationErrors)
		fmt.Println()
		fmt.Println("Expected fake error. It's OK.")
		fmt.Println(unexpectedError)
		fmt.Println()
	}
}
