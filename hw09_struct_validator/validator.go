package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	st := reflect.ValueOf(v)

	for i := 0; i < st.NumFiled(); i++ {
		field := st.Field(i)
		if validateTag, ok := field.Tag.Lookup("validate"); ok {
			if validateTag != "" {
				fmt.Println(validateTag)
			}
		}
	}

	return nil
}
