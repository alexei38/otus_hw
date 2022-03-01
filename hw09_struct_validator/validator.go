package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnsuportedValidator = errors.New("unsupported validator")
	ErrUnsuportedType      = errors.New("unsupported type")
	ErrMustInt             = errors.New("must be int")
	ErrMustMoreThan        = errors.New("must be more than")
	ErrMustLessThan        = errors.New("must be less than")
	ErrNotContainsIn       = errors.New("not contains in")
	ErrValidateRegex       = errors.New("regex validate failed")
	ErrLengthMustEqual     = errors.New("length must be")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := strings.Builder{}
	for _, err := range v {
		if err.Err != nil {
			s.WriteString(err.Err.Error())
			s.WriteString("\n")
		}
	}
	return s.String()
}

func validateStr(str string, validators map[string]string) error {
	for k, v := range validators {
		switch k {
		case "len":
			i, err := strconv.Atoi(v)
			if err != nil {
				return ErrMustInt
			}
			if len(str) != i {
				return fmt.Errorf("%w %d", ErrLengthMustEqual, i)
			}
		case "regexp":
			compile, err := regexp.Compile(v)
			if err != nil {
				return err
			}
			if !compile.MatchString(str) {
				return ErrValidateRegex
			}
		case "in":
			subs := strings.Split(v, ",")
			subs = deleteEmpty(subs)
			if len(subs) == 0 {
				return ErrUnsuportedValidator
			}
			found := contains(strings.Split(v, ","), str)
			if !found {
				return fmt.Errorf("'%s' %w %s", str, ErrNotContainsIn, v)
			}
		default:
			return ErrUnsuportedValidator
		}
	}
	return nil
}

func validateInt(val int, validators map[string]string) error {
	for k, v := range validators {
		switch k {
		case "min":
			i, err := strconv.Atoi(v)
			if err != nil {
				return ErrMustInt
			}
			if val < i {
				return fmt.Errorf("%w %d", ErrMustMoreThan, i)
			}
		case "max":
			i, err := strconv.Atoi(v)
			if err != nil {
				return ErrMustInt
			}
			if val > i {
				return fmt.Errorf("%w %d", ErrMustLessThan, i)
			}
		case "in":
			subs := deleteEmpty(strings.Split(v, ","))
			if len(subs) == 0 {
				return ErrUnsuportedValidator
			}
			for _, i := range strings.Split(v, ",") {
				_, err := strconv.Atoi(i)
				if err != nil {
					return ErrMustInt
				}
			}
			found := contains(strings.Split(v, ","), strconv.Itoa(val))
			if !found {
				return fmt.Errorf("'%d' %w %s", val, ErrNotContainsIn, v)
			}
		default:
			return ErrUnsuportedValidator
		}
	}
	return nil
}

func validateSlice(name string, value reflect.Value, validators map[string]string) ValidationErrors {
	var errors ValidationErrors
	for i := 0; i < value.Len(); i++ {
		v := value.Index(i)
		switch value.Type().Elem().Kind() { //nolint:exhaustive
		case reflect.String:
			err := validateStr(v.String(), validators)
			if err != nil {
				errors = append(errors, ValidationError{name, err})
			}
		case reflect.Int:
			err := validateInt(int(v.Int()), validators)
			if err != nil {
				errors = append(errors, ValidationError{name, err})
			}
		default:
			errors = append(errors, ValidationError{name, ErrUnsuportedType})
		}
	}
	return errors
}

func validateStruct(value reflect.Value) ValidationErrors {
	var errors ValidationErrors
	rt := value.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fv := value.Field(i)
		validators, ok := readTag(field)
		if !ok {
			continue
		}
		errors = append(errors, validateField(field, fv, validators)...)
	}
	return errors
}

func validateField(field reflect.StructField, value reflect.Value, validators map[string]string) ValidationErrors {
	var errors ValidationErrors
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		err := validateStr(value.String(), validators)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.Int:
		err := validateInt(int(value.Int()), validators)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.Slice:
		errors = append(errors, validateSlice(field.Name, value, validators)...)
	case reflect.Struct:
		for k := range validators {
			switch k {
			case "nested":
				errors = append(errors, validateStruct(value)...)
			default:
				errors = append(errors, ValidationError{
					Field: field.Name,
					Err:   ErrUnsuportedValidator,
				})
			}
		}
	default:
		errors = append(errors, ValidationError{
			Field: field.Name,
			Err:   ErrUnsuportedType,
		})
	}
	return errors
}

func Validate(v interface{}) error {
	var errors ValidationErrors
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Struct {
		errors = append(errors, ValidationError{
			Field: rv.Type().Name(),
			Err:   ErrUnsuportedType,
		})
		return errors
	}
	errors = append(errors, validateStruct(rv)...)

	if len(errors) > 0 {
		return errors
	}
	return nil
}
