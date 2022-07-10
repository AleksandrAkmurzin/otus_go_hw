package hw09structvalidator

import (
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
		m      Meta     `validate:"nested"`
	}

	Meta struct {
		app App `validate:"nested"`
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
		m:      Meta{App{"1.2.3"}},
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
			func(user User) User { user.m.app.Version = "invalid"; return user }(validUser),
			ValidationErrors{
				{"m.app.Version", ErrInvalidStringLength},
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
					require.Equal(t, tt.expectedErr, vErrs)
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
		name        string
		value       interface{}
		errContains string
	}{
		{"OnlyStructs", "no struct", "struct"},
		{"No bool validation", struct {
			b bool `validate:"in:true"`
		}{}, "bool"},
		{"UnsupportedIntRule", struct {
			i int `validate:"sign:unsigned"`
		}{}, "sign"},
		{"UnsupportedStringRule", struct {
			s string `validate:"maxLength:1024"`
		}{}, "maxLength"},
		{"InvalidRuleSyntax", struct {
			invalidInt int `validate:"min:a"`
		}{}, "invalidInt"},
		{"NestedField", struct {
			n struct {
				i int `validate:"unknown:"`
			} `validate:"nested"`
		}{}, "n.i"},
		{"Valid", struct {
			s string `validate:"len:3"`
		}{"abc"}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.value)
			require.Error(t, err)
			require.ErrorContains(t, err, test.errContains)

			var vErrs ValidationErrors
			if errors.As(err, &vErrs) {
				require.Zero(t, len(vErrs))
			}
		})
	}
}
