package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}
	for _, err := range v {
		strBuilder.WriteString("Field \"" + err.Field + "\": " + err.Err.Error() + "\n")
	}
	return strBuilder.String()
}

func (v ValidationErrors) Add(errList []ValidationError) []ValidationError {
	return append(v, errList...)
}

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return v.Field + ":" + v.Err.Error()
}

func NewValidationError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

type ProgramError struct {
	Err error
}

func (p ProgramError) Error() string {
	return "program error: " + p.Err.Error()
}

func NewProgramError(msg string) ProgramError {
	return ProgramError{
		Err: errors.New(msg),
	}
}

func Validate(v interface{}) error {
	errorList := make(ValidationErrors, 0)

	valueType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		// program error
		return NewProgramError("the input is not a struct, app stopped")
	}

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		if isPrivateField(field.Name) {
			continue
		}

		fieldValue := value.Field(i)
		rules := field.Tag.Get("validate")
		if rules == "" {
			continue
		}
		rulesSlice := strings.Split(rules, "|")

		switch field.Type.Kind() { //nolint
		case reflect.Int:
			//
			errList := validateInt(fieldValue.Int(), field.Name, rulesSlice)
			if len(errList) != 0 {
				// validation error
				errorList = errorList.Add(errList)
			}

		case reflect.String:
			//
			errList := validateString(fieldValue.String(), field.Name, rulesSlice)
			if len(errList) != 0 {
				// validation error
				errorList = errorList.Add(errList)
			}

		case reflect.Slice:
			//
			sliceErrorList := validateSlice(fieldValue.Interface(), field.Name, rulesSlice)
			if len(sliceErrorList) > 0 {
				// validation error
				errorList = errorList.Add(sliceErrorList)
			}

		default:
			// validation error
			errorList = append(errorList, NewValidationError(
				field.Name,
				errors.New("the struct field type is not possible to validate"),
			))
		}
	}

	if len(errorList) == 0 {
		return nil
	}

	return errorList
}

func validateSlice(fieldValue interface{}, filedName string, rulesSlice []string) []ValidationError {
	var sliceErrorList []ValidationError

	switch sliceValues := fieldValue.(type) {
	case []int:
		for _, item := range sliceValues {
			errIntList := validateInt(int64(item), filedName, rulesSlice)
			if len(errIntList) != 0 {
				// validation error
				sliceErrorList = append(sliceErrorList, errIntList...)
			}
		}
	case []string:
		for _, item := range sliceValues {
			errStrList := validateString(item, filedName, rulesSlice)
			if len(errStrList) != 0 {
				// validation error
				sliceErrorList = append(sliceErrorList, errStrList...)
			}
		}
	}

	return sliceErrorList
}

func validateInt(fieldValue int64, filedName string, rulesSlice []string) []ValidationError {
	var errList []ValidationError

	for _, rule := range rulesSlice {
		parsedRules, err := parseRules(rule)
		if err != nil {
			continue
		}

		if ruleValue, ok := parsedRules["max"]; ok {
			fieldValueInt := int(fieldValue)
			ruleValueInt, errInt := strconv.Atoi(ruleValue)
			if errInt != nil {
				continue
			}
			if fieldValueInt > ruleValueInt {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value > max"),
					),
				)
				continue
			}
		}

		if ruleValue, ok := parsedRules["min"]; ok {
			fieldValueInt := int(fieldValue)
			ruleValueInt, errInt := strconv.Atoi(ruleValue)
			if errInt != nil {
				continue
			}
			if fieldValueInt < ruleValueInt {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value < min"),
					),
				)
				continue
			}
		}

		if ruleValue, ok := parsedRules["in"]; ok {
			fieldValueInt := int(fieldValue)

			inValues := strings.Split(ruleValue, ",")

			if len(inValues) != 2 {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("in rules is not valid"),
					),
				)
				continue
			}

			min, _ := strconv.Atoi(inValues[0])
			max, _ := strconv.Atoi(inValues[1])
			if fieldValueInt < min || fieldValueInt > max {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value is not in the range"),
					),
				)
				continue
			}
		}
	}

	return errList
}

func validateString(fieldValue, filedName string, rulesSlice []string) []ValidationError {
	var errList []ValidationError

	for _, rule := range rulesSlice {
		parsedRules, err := parseRules(rule)
		if err != nil {
			continue
		}

		if ruleValue, ok := parsedRules["len"]; ok {
			ruleValueInt, errStr := strconv.Atoi(ruleValue)
			if errStr != nil {
				continue
			}
			if len(fieldValue) != ruleValueInt {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value len is not correct"),
					),
				)
				continue
			}
		}

		if ruleValue, ok := parsedRules["regexp"]; ok {
			match, errRegEx := regexp.MatchString(ruleValue, fieldValue)
			if !match || errRegEx != nil {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value is not matched by regex"),
					),
				)
				continue
			}
		}

		if ruleValue, ok := parsedRules["in"]; ok {
			inValues := strings.Split(ruleValue, ",")
			if len(inValues) != 2 {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("in rules is not valid"),
					),
				)
				continue
			}

			if !contains(inValues, fieldValue) {
				// validation error
				errList = append(
					errList,
					NewValidationError(
						filedName,
						errors.New("value not in valid list"),
					),
				)
				continue
			}
		}
	}

	return errList
}

func parseRules(source string) (paramsMap map[string]string, err error) {
	paramsMap = make(map[string]string)

	parsedSource := strings.Split(source, ":")
	ruleKey := parsedSource[0]
	ruleParam := parsedSource[1]
	paramsMap[ruleKey] = ruleParam

	if paramsMap == nil {
		return paramsMap, NewProgramError("not acceptable validate rule")
	}

	return paramsMap, nil
}

func isPrivateField(name string) bool {
	firstRune := []rune(name)[0]
	return unicode.IsLower(firstRune)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
