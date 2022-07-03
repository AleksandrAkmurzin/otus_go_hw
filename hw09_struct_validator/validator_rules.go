package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

	intValidationRules = map[string]func(condition string, testValue int) error{
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

	stringValidationRules = map[string]func(condition string, value string) error{
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
)
