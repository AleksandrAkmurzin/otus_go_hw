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

func (vErr ValidationError) Error() string {
	b := strings.Builder{}
	b.WriteString("error on field [")
	b.WriteString(vErr.Field)
	b.WriteString("]: ")
	b.WriteString(vErr.Err.Error())
	return b.String()
}

func (v ValidationErrors) Error() string {
	childErrors := make([]string, 0, len(v))
	for i := range v {
		childErrors = append(childErrors, v[i].Error())
	}

	return "validation errors: " + strings.Join(childErrors, ", ")
}

var (
	ErrValidation = errors.New("validation error")

	ErrInvalidInt      = fmt.Errorf("%w: invalid int", ErrValidation)
	ErrIntTooLow       = fmt.Errorf("%w: value less then allowed min", ErrInvalidInt)
	ErrIntTooLarge     = fmt.Errorf("%w: value more then allowed max", ErrInvalidInt)
	ErrIntNotInAllowed = fmt.Errorf("%w: value is not present in allowed list", ErrInvalidInt)

	ErrInvalidString        = fmt.Errorf("%w: invalid string", ErrValidation)
	ErrInvalidStringLength  = fmt.Errorf("%w: invalid string length", ErrInvalidString)
	ErrStringNotMatchRegexp = fmt.Errorf("%w: string not match given regexp", ErrInvalidString)
	ErrStringNotInAllowed   = fmt.Errorf("%w: value is not present in allowed list", ErrInvalidString)
)

func Validate(v interface{}) error {
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct {
		return errors.New("unsupported input value: must be struct")
	}

	rv := reflect.ValueOf(v)

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < rt.NumField(); i++ {
		vErrs, err := validateStructField(rt.Field(i), rv.Field(i), "")
		if err != nil {
			return err
		}

		validationErrors = append(validationErrors, vErrs...)
	}

	return validationErrors
}

func validateStructField(
	f reflect.StructField,
	v reflect.Value,
	parentFieldName string,
) (vErrs ValidationErrors, err error) {
	tag, ok := f.Tag.Lookup("validate")
	if !ok {
		return
	}

	if tag == "nested" {
		if v.Kind() != reflect.Struct {
			err = errors.New("tag validate:nested applicable only for struct")
			return
		}

		for i := 0; i < f.Type.NumField(); i++ {
			nestedFieldValue := v.Field(i)
			nestedParentFieldName := f.Name
			if parentFieldName != "" {
				nestedParentFieldName = parentFieldName + "." + nestedParentFieldName
			}
			nestedVErrs, err := validateStructField(
				f.Type.Field(i),
				nestedFieldValue,
				nestedParentFieldName,
			)
			if err != nil {
				return vErrs, err
			}

			vErrs = append(vErrs, nestedVErrs...)
		}

		return
	}

	fieldName := f.Name
	if parentFieldName != "" {
		fieldName = parentFieldName + "." + fieldName
	}

	conditions := strings.Split(tag, "|")
	for i := range conditions {
		parts := strings.Split(conditions[i], ":")
		if len(parts) != 2 {
			err = fmt.Errorf("wrong validation rule format for field %s", fieldName)
			return
		}

		conditionVErrs, err := validateCondition(parts[0], parts[1], fieldName, v)
		if err != nil {
			return vErrs, err
		}

		vErrs = append(vErrs, conditionVErrs...)
	}

	return
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

			vErrs = append(vErrs, sliceValidationErrors...)
		}

	default:
		return ValidationErrors{}, fmt.Errorf("validation of type %v is not supported", kind)
	}

	if err != nil {
		err = fmt.Errorf("%w for field %s", err, fieldName)
	}

	return
}

func validateString(value string, ruleName string, condition string) error {
	stringValidationRules := map[string]func(condition string, value string) error{
		"len": func(exactLen string, testValue string) error {
			targetValueLength, err := strconv.Atoi(exactLen)
			if err != nil {
				return err
			}
			if len(testValue) != targetValueLength {
				return ErrInvalidStringLength
			}

			return nil
		},
		"regexp": func(re string, testValue string) error {
			expr, err := regexp.Compile(re)
			if err != nil {
				return err
			}

			if !expr.MatchString(testValue) {
				return ErrStringNotMatchRegexp
			}

			return nil
		},
		"in": func(allowedStrings string, testValue string) error {
			for _, allowedString := range strings.Split(allowedStrings, ",") {
				if testValue == allowedString {
					return nil
				}
			}

			return ErrStringNotInAllowed
		},
	}

	validator, ok := stringValidationRules[ruleName]
	if !ok {
		return errors.New("unsupported string validation rule: " + ruleName)
	}

	return validator(condition, value)
}

func validateInt(value int, ruleName string, condition string) error {
	intValidationRules := map[string]func(condition string, testValue int) error{
		"min": func(allowedMinValue string, testValue int) error {
			minValue, err := strconv.Atoi(allowedMinValue)
			if err != nil {
				return err
			}
			if testValue < minValue {
				return ErrIntTooLow
			}

			return nil
		},
		"max": func(allowedMaxValue string, testValue int) error {
			maxValue, err := strconv.Atoi(allowedMaxValue)
			if err != nil {
				return err
			}
			if testValue > maxValue {
				return ErrIntTooLarge
			}

			return nil
		},
		"in": func(allowedIntegers string, testValue int) error {
			for _, allowedIntS := range strings.Split(allowedIntegers, ",") {
				allowedInt, err := strconv.Atoi(allowedIntS)
				if err != nil {
					return err
				}

				if testValue == allowedInt {
					return nil
				}
			}

			return ErrIntNotInAllowed
		},
	}

	validator, ok := intValidationRules[ruleName]
	if !ok {
		return errors.New("unsupported int validation rule: " + ruleName)
	}

	return validator(condition, value)
}
