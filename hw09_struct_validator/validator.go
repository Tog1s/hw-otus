package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}
type ValidationErrors []ValidationError

type ValidatorMap map[string]string

type ValidatorFn func(reflect.Value, ValidatorMap) error

var (
	ErrorValNotStruct                = errors.New("value is not a struct")
	ErrorValIncorrectRule            = errors.New("unable to parse rule")
	ErrorValIncorrectStringLength    = errors.New("incorrect string length")
	ErrorValRegexString              = errors.New("string do not satisfy regex")
	ErrorValNoMatchingElementInSlice = errors.New("there is no matching element")
	ErrorInputMinLimit               = errors.New("provided value is less than min")
	ErrorInputMaxLimit               = errors.New("provided value is more than max")
	ErrorValUnsupportedType          = errors.New("provided type is not supported")
	ErrorValTagValueShouldBeInteger  = errors.New("validation tag value should be integer")
	ErrorValIncorrectTagRegexPattern = errors.New("provided regexp pattern is invalid")
)

func (v ValidationErrors) Error() string {
	err := strings.Builder{}

	for _, e := range v {
		errMessage := e.Err.Error()
		if err.Len() != 0 {
			errMessage += "\n\t"
		}
		err.WriteString(errMessage)
	}

	return err.String()
}

func Validate(v interface{}) error {
	var errs ValidationErrors

	reflectType := reflect.TypeOf(v)
	reflectValue := reflect.ValueOf(v)

	if reflectType.Kind() != reflect.Struct {
		return ErrorValNotStruct
	}

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectType.Field(i)
		value := reflectValue.Field(i)

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		validator, err := parseTag(validateTag)
		if err != nil {
			return append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
		errs = append(errs, validate(field, value, validator)...)
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func parseTag(tag string) (ValidatorMap, error) {
	validatorMap := make(ValidatorMap)

	for _, rule := range strings.Split(tag, "|") {
		fieldRule := strings.Split(rule, ":")

		if len(fieldRule) != 2 || (fieldRule[0] == "") || fieldRule[1] == "" {
			return nil, fmt.Errorf("%w: rule: %s", ErrorValIncorrectRule, tag)
		}

		validatorMap[fieldRule[0]] = fieldRule[1]
	}

	return validatorMap, nil
}

func validate(field reflect.StructField, value reflect.Value, rules ValidatorMap) ValidationErrors {
	var errs ValidationErrors
	addError := func(field string, err error) {
		errs = append(errs, ValidationError{
			Field: field,
			Err:   err,
		})
	}

	kind := field.Type.Kind()
	switch kind { //nolint: exhaustive
	case reflect.Int:
		if err := validateInt(value, rules); err != nil {
			addError(field.Name, err)
		}
	case reflect.String:
		if err := validateString(value, rules); err != nil {
			addError(field.Name, err)
		}
	case reflect.Slice:
		elementKind := field.Type.Elem().Kind()
		switch elementKind { //nolint: exhaustive
		case reflect.Int:
			for _, err := range validateSlice(validateInt, value, rules) {
				addError(field.Name, err)
			}
		case reflect.String:
			for _, err := range validateSlice(validateString, value, rules) {
				addError(field.Name, err)
			}
		default:
			addError(field.Name, ErrorValUnsupportedType)
		}
	default:
		addError(field.Name, ErrorValUnsupportedType)
	}
	return errs
}

func validateInt(value reflect.Value, rules ValidatorMap) error {
	for rule, val := range rules {
		switch rule {
		case "max":
		case "min":
			expected, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid value %s: %s", ErrorValTagValueShouldBeInteger, rule, val)
			}

			if value.Int() < expected {
				return ErrorInputMinLimit
			}

			if rule == "max" && value.Int() > expected {
				return ErrorInputMaxLimit
			}
		case "in":
			for _, integer := range strings.Split(val, ",") {
				expected, err := strconv.ParseInt(integer, 10, 64)
				if err != nil {
					return fmt.Errorf("%w: invalid slice value: %s", ErrorValTagValueShouldBeInteger, integer)
				}

				if value.Int() == expected {
					return nil
				}
			}

			return ErrorValNoMatchingElementInSlice
		}
	}

	return nil
}

func validateString(value reflect.Value, rules ValidatorMap) error {
	for rule, val := range rules {
		switch rule {
		case "len":
			expected, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("%w: invalid value length: %s", err, val)
			}

			if value.Len() != expected {
				return ErrorValIncorrectStringLength
			}
		case "regexp":
			rx, err := regexp.Compile(val)
			if err != nil {
				return fmt.Errorf("%w; invalid regex: %s", ErrorValIncorrectTagRegexPattern, val)
			}

			if !rx.MatchString(value.String()) {
				return ErrorValRegexString
			}
		case "in":
			for _, word := range strings.Split(val, ",") {
				if value.String() == word {
					return nil
				}
			}

			return ErrorValNoMatchingElementInSlice
		}
	}

	return nil
}

func validateSlice(fn ValidatorFn, values reflect.Value, rules ValidatorMap) []error {
	var errs []error

	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)

		if err := fn(value, rules); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (v ValidationErrors) Is(target error) bool {
	var targetErr ValidationErrors

	if !errors.As(target, &targetErr) {
		return false
	}

	if len(v) != len(targetErr) {
		return false
	}

	for i := 0; i < len(targetErr); i++ {
		areFieldsEqual := v[i].Field == targetErr[i].Field
		areErrsEqual := errors.Is(v[i].Err, targetErr[i].Err)
		if !areFieldsEqual || !areErrsEqual {
			return false
		}
	}

	return true
}
