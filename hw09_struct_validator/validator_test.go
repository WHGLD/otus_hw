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
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{Version: "1234"},
			ValidationErrors{
				ValidationError{"Version", errors.New("value len is not correct")},
			},
		},
		{App{Version: "12345"}, nil},

		{
			Response{Code: 200, Body: ""},
			ValidationErrors{
				ValidationError{"Code", errors.New("in rules is not valid")},
			},
		},

		{Token{}, nil},

		{
			User{
				ID:     "12345",
				Name:   "",
				Age:    111,
				Email:  "email@ya.ru",
				Role:   "admin",
				Phones: []string{"12345678912"},
			},
			ValidationErrors{
				ValidationError{"ID", errors.New("value len is not correct")},
				ValidationError{"Age", errors.New("value > max")},
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
				require.ErrorAs(t, err, &ValidationErrors{})
				require.EqualError(t, err, tt.expectedErr.Error())
			}

			_ = tt
		})
	}

	t.Run("case - not correct in", func(t *testing.T) {
		in := "not struct"
		errOut := NewProgramError("the input is not a struct, app stopped")

		err := Validate(in)
		if err == nil {
			require.NoError(t, err)
		} else {
			require.ErrorAs(t, err, &ProgramError{})
			require.EqualError(t, err, errOut.Error())
		}
	})
}
