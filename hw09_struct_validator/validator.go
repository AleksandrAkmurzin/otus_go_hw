package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (vErr ValidationError) Error() string {
	childErr := strings.Builder{}
	childErr.WriteString("error on field [")
	childErr.WriteString(vErr.Field)
	childErr.WriteString("]: ")
	childErr.WriteString(vErr.Err.Error())
	return childErr.String()
}

func (v ValidationErrors) Error() string {
	childErrors := make([]string, 0, len(v))
	for i := range v {
		childErrors = append(childErrors, v[i].Error())
	}

	return "validation errors: " + strings.Join(childErrors, ", ")
}

func Validate(v interface{}) error {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct {
		return errors.New("unsupported input value: must be struct")
	}

	rv := reflect.ValueOf(v)

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag, ok := f.Tag.Lookup("validate")
		if !ok {
			continue
		}

		value := rv.Field(i)

		conditions := strings.Split(tag, "|")
		for i := range conditions {
			parts := strings.Split(conditions[i], ":")
			if len(parts) != 2 {
				return fmt.Errorf("wrong validation rule format for field %s", f.Name)
			}
			vErrs, err := validateCondition(parts[0], parts[1], f.Name, value)
			if err != nil {
				return err
			}

			for errIndex := range vErrs {
				validationErrors = append(validationErrors, vErrs[errIndex])
			}
		}
	}

	return validationErrors
}

func validateCondition(
	ruleName string,
	ruleCondition string,
	fieldName string,
	v reflect.Value,
) (vErrs ValidationErrors, err error) {
	kind := v.Kind()
	// nolint:exhaustive
	switch kind {
	case reflect.String:
		err = validateString(v.String(), ruleName, ruleCondition)
		if errors.Is(err, ErrValidation) {
			vErrs = append(vErrs, ValidationError{fieldName, err})
			err = nil
		}

	case reflect.Int:
		err = validateInt(int(v.Int()), ruleName, ruleCondition)
		if errors.Is(err, ErrValidation) {
			vErrs = append(vErrs, ValidationError{fieldName, err})
			err = nil
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			sliceValidationErrors, err := validateCondition(
				ruleName,
				ruleCondition,
				fmt.Sprintf("%s[%d]", fieldName, i),
				v.Index(i),
			)
			if err != nil {
				return vErrs, err
			}
			for errIndex := range sliceValidationErrors {
				vErrs = append(vErrs, sliceValidationErrors[errIndex])
			}
		}

	default:
		return ValidationErrors{}, fmt.Errorf("validation of type %v is not supported", kind)
	}

	return
}

func validateString(value string, ruleName string, condition string) error {
	validator, ok := stringValidationRules[ruleName]
	if !ok {
		return errors.New("unsupported string validation rule: " + ruleName)
	}

	return validator(condition, value)
}

func validateInt(value int, ruleName string, condition string) error {
	validator, ok := intValidationRules[ruleName]
	if !ok {
		return errors.New("unsupported int validation rule: " + ruleName)
	}

	return validator(condition, value)
}
