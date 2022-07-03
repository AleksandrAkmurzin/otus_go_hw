package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
)

func TestValidate(t *testing.T) {
	validUser := User{
		ID:     "123456789a123456789b123456789c123456",
		Name:   "No matter",
		Age:    25,
		Email:  "mail@example.com",
		Role:   "admin",
		Phones: []string{"12345678901", "1234567890a"},
		meta:   nil,
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			validUser,
			nil,
		},
		{
			func(user User) User { user.ID = "short"; return user }(validUser),
			ValidationErrors{{"ID", ErrInvalidStringLength}},
		},
		{
			func(user User) User { user.Age = 16; return user }(validUser),
			ValidationErrors{{"Age", ErrIntTooLow}},
		},
		{
			func(user User) User { user.Age = 60; return user }(validUser),
			ValidationErrors{{"Age", ErrIntTooLarge}},
		},
		{
			func(user User) User { user.Email = "invalid"; return user }(validUser),
			ValidationErrors{{"Email", ErrStringNotMatchRegexp}},
		},
		{
			func(user User) User { user.Role = "hacker"; return user }(validUser),
			ValidationErrors{{"Role", ErrStringNotInAllowed}},
		},
		{
			func(user User) User { user.Phones = []string{"valid_11_dg", "invalid"}; return user }(validUser),
			ValidationErrors{{"Phones[1]", ErrInvalidStringLength}},
		},
		{
			func(user User) User { user.Phones = []string{"invalid", "invalidToo"}; return user }(validUser),
			ValidationErrors{
				{"Phones[0]", ErrInvalidStringLength},
				{"Phones[1]", ErrInvalidStringLength},
			},
		},
		{
			Response{200, ""},
			nil,
		},
		{
			Response{301, "Redirect"},
			ValidationErrors{{"Code", ErrIntNotInAllowed}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var vErrs ValidationErrors
			if errors.As(err, &vErrs) {
				if tt.expectedErr == nil {
					require.Zero(t, len(vErrs))
				} else {
					require.Equal(t, vErrs, tt.expectedErr)
				}
				return
			}

			require.NoError(t, err)
			_ = tt
		})
	}
}

func TestErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"Only structs", "no struct"},
		{"No bool validation", struct {
			b bool `validate:"unknown"`
		}{}},
		{"UnsupportedIntRule", struct {
			i int `validate:"sign:unsigned"`
		}{}},
		{"UnsupportedStringRule", struct {
			s string `validate:"maxLength:1024"`
		}{}},
		{"InvalidRuleSyntax", struct {
			i int `validate:"min:a"`
		}{}},
		{"Valid", struct {
			s string `validate:"len:3"`
		}{"abc"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.value)
			require.Error(t, err)

			var vErrs ValidationErrors
			if errors.As(err, &vErrs) {
				require.Zero(t, len(vErrs))
			}
		})
	}
}
