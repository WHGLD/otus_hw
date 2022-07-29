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
	//u := User{}
	//err := Validate(u)
	//if err != nil {
	//	return
	//}

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
		//
		//{Response{Code: 200, Body: ""}, nil},
		//{Response{Code: 201, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		//
		//{Response{Code: 403, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},
		//{Response{Code: 404, Body: ""}, nil},
		//
		//{Response{Code: 500, Body: ""}, nil},
		//{Response{Code: 501, Body: ""}, ValidationErrors{ValidationError{Field: "Code", Err: ErrIn}}},

		{User{
			ID:     "12345",
			Name:   "",
			Age:    99,
			Email:  "email@ya.ru",
			Role:   "admin",
			Phones: []string{"12345678912"},
		}, nil},
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
}
