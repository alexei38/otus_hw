package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

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

	Product struct {
		App App `validate:"nested"`
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

type Tests []struct {
	in          interface{}
	expectedErr error
}

func runTests(t *testing.T, tests Tests) {
	t.Helper()
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errs := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, errs)
			} else {
				require.Error(t, errs)
				verrs, ok := errs.(ValidationErrors) //nolint:errorlint
				require.True(t, ok)
				for _, verr := range verrs {
					require.NotNil(t, verr.Err)
					require.ErrorIs(t, verr.Err, tt.expectedErr)
				}
			}
		})
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestValidate(t *testing.T) {
	tests := Tests{
		{App{"12345"}, nil},
		{App{"123"}, ErrLengthMustEqual},
		{App{"123456"}, ErrLengthMustEqual},
		{Product{App{"12345"}}, nil},
		{Product{App{"123"}}, ErrLengthMustEqual},
		{Product{App{"123456"}}, ErrLengthMustEqual},
		{Response{100, `{key: "value"}`}, ErrNotContainsIn},
		{Response{200, `{key: "value"}`}, nil},
		{Response{404, `{key: "value"}`}, nil},
		{Response{500, `{key: "value"}`}, nil},
		{Response{600, `{key: "value"}`}, ErrNotContainsIn},
		{
			User{
				ID:    randString(36),
				Name:  randString(10),
				Age:   18,
				Email: "test@example.com",
				Role:  UserRole("admin"),
				meta:  nil,
			},
			nil,
		},
		{
			User{
				ID:    "",
				Name:  randString(10),
				Age:   18,
				Email: "test@example.com",
				Role:  UserRole("admin"),
			},
			ErrLengthMustEqual,
		},
		{
			User{
				ID:    randString(36),
				Name:  randString(10),
				Age:   0,
				Email: "test@example.com",
				Role:  "admin",
			},
			ErrMustMoreThan,
		},
		{
			User{
				ID:    randString(36),
				Name:  randString(10),
				Age:   55,
				Email: "test@example.com",
				Role:  "admin",
			},
			ErrMustLessThan,
		},
		{
			User{
				ID:    randString(36),
				Name:  randString(10),
				Age:   18,
				Email: randString(10),
				Role:  "admin",
			},
			ErrValidateRegex,
		},
		{
			User{
				ID:     randString(36),
				Name:   randString(10),
				Age:    18,
				Email:  "test@example.com",
				Role:   "stuff",
				Phones: []string{randString(11), randString(11)},
			},
			nil,
		},
		{
			User{
				ID:     randString(36),
				Name:   randString(10),
				Age:    18,
				Email:  "test@example.com",
				Role:   "",
				Phones: []string{randString(11), randString(11)},
			},
			ErrNotContainsIn,
		},
		{
			User{
				ID:     randString(36),
				Name:   randString(10),
				Age:    18,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{randString(21), randString(11)},
			},
			ErrLengthMustEqual,
		},
		{
			User{
				ID:     randString(36),
				Name:   randString(10),
				Age:    18,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{randString(11), randString(21)},
			},
			ErrLengthMustEqual,
		},
	}
	runTests(t, tests)
}

func TestStringValidators(t *testing.T) {
	tests := Tests{
		{
			struct {
				Key string `validate:"len:string"`
			}{},
			ErrMustInt,
		},
		{
			struct {
				Key string `validate:"regexp:!@#$%!@#$%^&*IJHGFDRET"`
			}{Key: "!"},
			ErrValidateRegex,
		},
		{
			struct {
				Key string `validate:"regexp:\\d+"`
			}{Key: "123456"},
			nil,
		},
		{
			struct {
				Key string `validate:"unknown:validator"`
			}{},
			ErrUnsuportedValidator,
		},
		{
			struct {
				Key string `validate:"in:admin,support"`
			}{Key: "admin"},
			nil,
		},
		{
			struct {
				Key string `validate:"in:admin,support"`
			}{Key: "support"},
			nil,
		},
		{
			struct {
				Key string `validate:"in:admin,support"`
			}{Key: "stuff"},
			ErrNotContainsIn,
		},
		{
			struct {
				Key string `validate:"in"`
			}{Key: ""},
			ErrUnsuportedValidator,
		},
	}
	runTests(t, tests)
}

func TestIntValidators(t *testing.T) {
	tests := Tests{
		{
			struct {
				Key int `validate:"len:string"`
			}{},
			ErrUnsuportedValidator,
		},
		{
			struct {
				Key int `validate:"min:string"`
			}{},
			ErrMustInt,
		},
		{
			struct {
				Key int `validate:"min:5"`
			}{0},
			ErrMustMoreThan,
		},
		{
			struct {
				Key int `validate:"min:0"`
			}{5},
			nil,
		},
		{
			struct {
				Key int `validate:"max:string"`
			}{},
			ErrMustInt,
		},
		{
			struct {
				Key int `validate:"max:0"`
			}{1},
			ErrMustLessThan,
		},
		{
			struct {
				Key int `validate:"max:5"`
			}{0},
			nil,
		},
		{
			struct {
				Key int `validate:"in"`
			}{},
			ErrUnsuportedValidator,
		},
		{
			struct {
				Key int `validate:"in:"`
			}{},
			ErrUnsuportedValidator,
		},
		{
			struct {
				Key int `validate:"in:string"`
			}{},
			ErrMustInt,
		},
		{
			struct {
				Key int `validate:"in:100"`
			}{100},
			nil,
		},
		{
			struct {
				Key int `validate:"in:100,200,300"`
			}{200},
			nil,
		},
		{
			struct {
				Key int `validate:"in:100,200,300"`
			}{300},
			nil,
		},
		{
			struct {
				Key int `validate:"in:100,200,300"`
			}{150},
			ErrNotContainsIn,
		},
		{
			struct {
				Key int `validate:"min:200|in:200,300,400|max:400"`
			}{
				200,
			},
			nil,
		},
		{
			struct {
				Key int `validate:"min:300|in:200,300,400|max:400"`
			}{
				200,
			},
			ErrMustMoreThan,
		},
		{
			struct {
				Key int `validate:"min:200|in:200,300,400|max:250"`
			}{
				300,
			},
			ErrMustLessThan,
		},
	}
	runTests(t, tests)
}

func TestSliceValidators(t *testing.T) {
	tests := Tests{
		{
			struct {
				Key []string `validate:"len:string"`
			}{
				[]string{"test"},
			},
			ErrMustInt,
		},
		{
			struct {
				Key []string `validate:"len:10"`
			}{
				[]string{randString(10), randString(10), randString(10)},
			},
			nil,
		},
		{
			struct {
				Key []string `validate:"len:10"`
			}{
				[]string{randString(10), randString(30), randString(10)},
			},
			ErrLengthMustEqual,
		},
		{
			struct {
				Key []string `validate:"len:10"`
			}{
				[]string{randString(10), randString(10), randString(30)},
			},
			ErrLengthMustEqual,
		},
		{
			struct {
				Key []string `validate:"in:foo,bar"`
			}{
				[]string{randString(10), randString(10), randString(30)},
			},
			ErrNotContainsIn,
		},
		{
			struct {
				Key []string `validate:"in:foo,bar"`
			}{
				[]string{"foo", randString(10), randString(30)},
			},
			ErrNotContainsIn,
		},
		{
			struct {
				Key []string `validate:"in:foo,bar"`
			}{
				[]string{"foo", "bar", "foo"},
			},
			nil,
		},
		{
			struct {
				Key []int `validate:"min:100"`
			}{
				[]int{200, 300},
			},
			nil,
		},
		{
			struct {
				Key []int `validate:"min:201"`
			}{
				[]int{200, 300},
			},
			ErrMustMoreThan,
		},
		{
			struct {
				Key []int `validate:"max:250"`
			}{
				[]int{200, 300},
			},
			ErrMustLessThan,
		},
		{
			struct {
				Key []int `validate:"max:300"`
			}{
				[]int{200, 300},
			},
			nil,
		},
		{
			struct {
				Key []int `validate:"in:200,300,400"`
			}{
				[]int{200, 300},
			},
			nil,
		},
		{
			struct {
				Key []int `validate:"in:200,300,400"`
			}{
				[]int{200, 300, 500},
			},
			ErrNotContainsIn,
		},
		{
			struct {
				Key []byte `validate:"in:200,300,400"`
			}{
				[]byte{200},
			},
			ErrUnsuportedType,
		},
		{
			struct {
				Key []int `validate:"min:200|in:200,300,400|max:400"`
			}{
				[]int{200, 300},
			},
			nil,
		},
		{
			struct {
				Key []int `validate:"min:300|in:200,300,400|max:400"`
			}{
				[]int{200, 300},
			},
			ErrMustMoreThan,
		},
		{
			struct {
				Key []int `validate:"min:200|in:200,300,400|max:250"`
			}{
				[]int{200, 300},
			},
			ErrMustLessThan,
		},
	}
	runTests(t, tests)
}

func TestStructValidators(t *testing.T) {
	tests := Tests{
		{
			struct {
				Data byte `validate:"len:100"`
			}{},
			ErrUnsuportedType,
		},
		{
			"42",
			ErrUnsuportedType,
		},
		{
			42,
			ErrUnsuportedType,
		},
		{
			struct {
				Nested string `validate:"nested"`
			}{"test"},
			ErrUnsuportedValidator,
		},
		{
			struct {
				Nested struct {
					ID string `validate:"len:30"`
				} `validate:"nested"`
			}{struct {
				ID string `validate:"len:30"`
			}(struct{ ID string }{ID: randString(30)})},
			nil,
		},
		{
			struct {
				Nested struct {
					ID string `validate:"len:30"`
				} `validate:"nested"`
			}{struct {
				ID string `validate:"len:30"`
			}(struct{ ID string }{ID: randString(10)})},
			ErrLengthMustEqual,
		},
		{
			struct {
				Nested struct {
					ID string `validate:"len:30"`
				} `validate:"unknownValidator"`
			}{struct {
				ID string `validate:"len:30"`
			}(struct{ ID string }{ID: randString(10)})},
			ErrUnsuportedValidator,
		},
	}
	runTests(t, tests)
}
