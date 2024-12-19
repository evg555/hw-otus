package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTagName      = "validate"
	rulesDelimiter       = "|"
	ruleNameValDelimiter = ":"

	ruleLen    = "len"
	ruleIn     = "in"
	ruleMax    = "max"
	ruleMin    = "min"
	ruleRegexp = "regexp"
)

var (
	ErrNotStruct            = errors.New("not a struct")
	ErrInvalidRule          = errors.New("invalid rule")
	ErrUnsupportedFieldType = errors.New("unsupported field type")
	ErrInvalidRegexp        = errors.New("invalid regexp")
	ErrUnsupportedRuleName  = errors.New("unsupported rule name")
	ErrNotNumber            = errors.New("not a number")
)

type ProgramError struct {
	Err error
}

func (p ProgramError) Error() string {
	return p.Err.Error()
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))

	for i, err := range v {
		errs[i] = fmt.Sprintf("%s: %s", err.Field, err.Err)
	}

	return "validation errors: " + strings.Join(errs, "; ")
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	refVal := reflect.ValueOf(v)

	if refVal.Kind() != reflect.Struct {
		return ProgramError{Err: fmt.Errorf("input '%v': %w", v, ErrNotStruct)}
	}

	var validationErrors ValidationErrors

	refType := refVal.Type()

	for i := 0; i < refVal.NumField(); i++ {
		var fieldErrors ValidationErrors

		field := refVal.Field(i)
		fieldType := refType.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		validateTag := fieldType.Tag.Get(validateTagName)
		if validateTag == "" {
			continue
		}

		err := validateField(field, validateTag, fieldType.Name)
		if err != nil {
			if errors.As(err, &fieldErrors) {
				if len(fieldErrors) > 0 {
					validationErrors = append(validationErrors, fieldErrors...)
				}
			} else {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(field reflect.Value, tag string, name string) error {
	var errs ValidationErrors

	rules := strings.Split(tag, rulesDelimiter)

	for _, rule := range rules {
		err := applyRule(field, rule)
		if err != nil {
			if errors.As(err, &ProgramError{}) {
				return err
			}

			errs = append(errs, ValidationError{
				Field: name,
				Err:   err,
			})
		}
	}

	return errs
}

func applyRule(field reflect.Value, rule string) error {
	parts := strings.SplitN(rule, ruleNameValDelimiter, 2)
	if len(parts) != 2 {
		return ProgramError{Err: fmt.Errorf("%s: %w", rule, ErrInvalidRule)}
	}

	ruleName, ruleValue := parts[0], parts[1]

	switch field.Kind() { //nolint:exhaustive
	case reflect.String:
		return validateString(field.String(), ruleName, ruleValue)
	case reflect.Int:
		return validateInt(int(field.Int()), ruleName, ruleValue)
	case reflect.Slice:
		return validateSlice(field, ruleName, ruleValue)
	default:
		return ProgramError{Err: fmt.Errorf("type '%s': %w", field.Type().Name(), ErrUnsupportedFieldType)}
	}
}

func validateString(value string, ruleName, ruleValue string) error {
	switch ruleName {
	case ruleLen:
		expectedLen, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ProgramError{Err: fmt.Errorf("value '%s': %w", ruleValue, ErrNotNumber)}
		}

		if len(value) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}
	case ruleRegexp:
		re, err := regexp.Compile(ruleValue)
		if err != nil {
			return ProgramError{Err: fmt.Errorf("rule value '%s': %w", ruleValue, ErrInvalidRegexp)}
		}

		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp %s", ruleValue)
		}
	case ruleIn:
		options := strings.Split(ruleValue, ",")

		for _, option := range options {
			if value == option {
				return nil
			}
		}

		return fmt.Errorf("must be one of [%s]", ruleValue)
	default:
		return ProgramError{Err: fmt.Errorf("rule name '%s': %w", ruleValue, ErrUnsupportedRuleName)}
	}

	return nil
}

func validateInt(value int, ruleName, ruleValue string) error {
	switch ruleName {
	case ruleMin:
		expectedMin, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ProgramError{Err: fmt.Errorf("value '%s': %w", ruleValue, ErrNotNumber)}
		}

		if value < expectedMin {
			return fmt.Errorf("must be >= %d", expectedMin)
		}
	case ruleMax:
		expectedMax, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ProgramError{Err: fmt.Errorf("value '%s': %w", ruleValue, ErrNotNumber)}
		}

		if value > expectedMax {
			return fmt.Errorf("must be <= %d", expectedMax)
		}
	case ruleIn:
		options := strings.Split(ruleValue, ",")
		for _, option := range options {
			opt, err := strconv.Atoi(option)
			if err != nil {
				return ProgramError{Err: fmt.Errorf("value '%s': %w", ruleValue, ErrNotNumber)}
			}

			if value == opt {
				return nil
			}
		}

		return fmt.Errorf("must be one of [%s]", ruleValue)
	default:
		return ProgramError{Err: fmt.Errorf("rule name '%s': %w", ruleValue, ErrUnsupportedRuleName)}
	}

	return nil
}

func validateSlice(field reflect.Value, ruleName, ruleValue string) error {
	for i := 0; i < field.Len(); i++ {
		if err := applyRule(field.Index(i), fmt.Sprintf("%s:%s", ruleName, ruleValue)); err != nil {
			if errors.As(err, &ProgramError{}) {
				return err
			}

			return fmt.Errorf("element %d: %w", i, err)
		}
	}

	return nil
}
