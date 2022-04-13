package types

import (
	"fmt"
	"reflect"
	"strings"
)

func Validate(object interface{}) error {
	objectType := reflect.TypeOf(object)
	objectValue := reflect.ValueOf(object)

	errorMessage := ""

	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)

		if !strings.Contains(string(field.Tag), "required") {
			continue
		}

		if objectValue.Field(i).String() == "" {
			errorMessage += fmt.Sprintf("field %s is required but no value was set\n", field.Name)
		}
	}

	if errorMessage == "" {
		return nil
	}

	return fmt.Errorf("error validating %s:\n%s", objectType.Name(), errorMessage)
}
