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

		if isPrivateField(field.Name) == true {
			continue
		}

		fieldValue := value.Field(i)
		rules := field.Tag.Get("validate")
		if rules == "" {
			continue
		}
		rulesSlice := strings.Split(rules, "|")

		switch field.Type.Kind() {
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
			var sliceErrorList []ValidationError

			switch fieldValue.Interface().(type) {
			case []int:
				sliceValues := fieldValue.Interface().([]int)
				for _, item := range sliceValues {
					errIntList := validateInt(int64(item), field.Name, rulesSlice)
					if len(errIntList) != 0 {
						// validation error
						sliceErrorList = append(sliceErrorList, errIntList...)
					}
				}
			case []string:
				sliceValues := fieldValue.Interface().([]string)
				for _, item := range sliceValues {
					errStrList := validateString(item, field.Name, rulesSlice)
					if len(errStrList) != 0 {
						// validation error
						sliceErrorList = append(sliceErrorList, errStrList...)
					}
				}
			}

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

func validateInt(fieldValue int64, filedName string, rulesSlice []string) []ValidationError {
	var errList []ValidationError

	for _, rule := range rulesSlice {
		parsedRules, err := parseRules(rule)
		if err != nil {
			continue
		}

		if ruleValue, ok := parsedRules["max"]; ok {
			fieldValueInt := int(fieldValue)
			ruleValueInt, _ := strconv.Atoi(ruleValue)
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
			ruleValueInt, _ := strconv.Atoi(ruleValue)
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
			ruleValueInt, _ := strconv.Atoi(ruleValue)
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
			match, _ := regexp.MatchString(ruleValue, fieldValue)
			if match == false {
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
			first := inValues[0]
			second := inValues[1]
			if fieldValue != first || fieldValue != second {
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
	regexps := map[string]string{
		"max":    `(?P<max>:\d{1,})`,
		"min":    `(?P<min>:\d{1,})`,
		"in":     `(?P<in>:.+)`,
		"regexp": `(?P<regexp>:.+)`,
		"len":    `(?P<len>:[0-9]{1,})`,
	}

	var regEx string

	for rule, regex := range regexps {
		match, _ := regexp.MatchString(rule, source)
		if match {
			regEx = regex
			break
		}
	}

	paramsMap = make(map[string]string)

	if regEx == "" {
		return paramsMap, NewProgramError("not acceptable validate rule")
	}

	compRegEx := regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(source)

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			value := strings.ReplaceAll(match[i], ":", "")
			paramsMap[name] = value
		}
	}

	return paramsMap, nil
}

func isPrivateField(name string) bool {
	firstRune := []rune(name)[0]
	return unicode.IsLower(firstRune)
}
