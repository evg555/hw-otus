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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: nil,
		},
		{
			in:          UserRole("test role"),
			expectedErr: ProgramError{Err: errors.New("input 'test role': not a struct")},
		},
		{
			in: Token{
				Payload: []byte("test payload"),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "5a4e751d-0f3b-4975-9884-aeceec16494d",
				Age:    20,
				Email:  "test@test.com",
				Role:   UserRole("admin"),
				Phones: []string{"+7911111111"},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "test",
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Field string `validate:"len-5"`
			}{Field: "test1"},
			expectedErr: ProgramError{Err: errors.New("len-5: invalid rule")},
		},
		{
			in: struct {
				Field string `validate:"len:test"`
			}{Field: "test1"},
			expectedErr: ProgramError{Err: errors.New("value 'test': not a number")},
		},
		{
			in: struct {
				Field bool `validate:"len:4"`
			}{Field: true},
			expectedErr: ProgramError{Err: errors.New("type 'bool': unsupported field type")},
		},
		{
			in: struct {
				Field string `validate:"regexp:]["`
			}{Field: "test1"},
			expectedErr: ProgramError{Err: errors.New("rule value '][': invalid regexp")},
		},
		{
			in: struct {
				Field string `validate:"test:test"`
			}{Field: "test1"},
			expectedErr: ProgramError{Err: errors.New("rule name 'test': unsupported rule name")},
		},
		{
			in: struct {
				Field int `validate:"min:test"`
			}{Field: 1},
			expectedErr: ProgramError{Err: errors.New("value 'test': not a number")},
		},
		{
			in: struct {
				Field int `validate:"max:test"`
			}{Field: 1},
			expectedErr: ProgramError{Err: errors.New("value 'test': not a number")},
		},
		{
			in: struct {
				Field int `validate:"in:test"`
			}{Field: 1},
			expectedErr: ProgramError{Err: errors.New("value 'test': not a number")},
		},
		{
			in: struct {
				Field int `validate:"test:test"`
			}{Field: 1},
			expectedErr: ProgramError{Err: errors.New("rule name 'test': unsupported rule name")},
		},
		{
			in: struct {
				Field []int `validate:"test:test"`
			}{Field: []int{1, 2, 3}},
			expectedErr: ProgramError{Err: errors.New("rule name 'test': unsupported rule name")},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, tt.expectedErr, err.Error())
				require.ErrorAs(t, err, &tt.expectedErr)
			}

			_ = tt
		})
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: nil,
		},
		{
			in: App{Version: "12"},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   errors.New("length must be 5"),
				},
			},
		},
		{
			in: Response{
				Code: 300,
				Body: "test",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   errors.New("must be one of [200,404,500]"),
				},
			},
		},
		{
			in: User{},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   errors.New("length must be 36"),
				},
				{
					Field: "Age",
					Err:   errors.New("must be >= 18"),
				},
				{
					Field: "Email",
					Err:   errors.New("must match regexp ^\\w+@\\w+\\.\\w+$"),
				},
				{
					Field: "Role",
					Err:   errors.New("must be one of [admin,stuff]"),
				},
			},
		},
		{
			in: User{
				ID:     "5a4e751d-0f3b-4975-9884-aeceec16494d",
				Email:  "test@test.com",
				Role:   UserRole("admin"),
				Age:    100,
				Phones: []string{"1111"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   errors.New("must be <= 50"),
				},
				{
					Field: "Phones",
					Err:   errors.New("element 0: length must be 11"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, tt.expectedErr, err.Error())
				require.ErrorAs(t, err, &tt.expectedErr)
			}

			_ = tt
		})
	}
}
