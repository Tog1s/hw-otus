package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
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
		meta   json.RawMessage //nolint:nolintlint
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
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "User struct is valid",
			in: User{
				ID:     "PXZVT49B5osdY70e0tbLruN4eme3pL6j2V0j",
				Name:   "John",
				Age:    21,
				Email:  "some_mail@example.com",
				Role:   UserRole("admin"),
				Phones: []string{"11122233344"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			name:        "Value is not a struct",
			in:          42,
			expectedErr: ErrorValNotStruct,
		},
		{
			name: "Unable to parse rule",
			in: struct {
				Value string `validate:"in"`
			}{
				Value: "for",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrorValIncorrectRule,
				},
			},
		},
		{
			name: "User is not valid",
			in: User{
				ID:     "8429084",
				Name:   "John",
				Age:    17,
				Email:  "some_mail",
				Role:   UserRole("root"),
				Phones: []string{"911"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrorValIncorrectStringLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrorInputMinLimit,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrorValRegexString,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrorValNoMatchingElementInSlice,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrorValIncorrectStringLength,
				},
			},
		},
		{
			name: "No tags, no errors",
			in: Token{
				[]byte{3, 6, 2, 5},
				[]byte{1, 8, 3, 6},
				[]byte{10, 11, 15, 16},
			},
			expectedErr: nil,
		},
		{
			name:        "Empty struct",
			in:          struct{}{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
