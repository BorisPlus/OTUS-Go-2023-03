package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func sliceAtoi(sa []string) ([]int64, error) {
	sliceInt := make([]int64, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return sliceInt, err
		}
		sliceInt = append(sliceInt, int64(i))
	}
	return sliceInt, nil
}

const validationTagName = "validate"

// ValidationError - struct of validation error.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors - stack trace of of validation errors.
type ValidationErrors []ValidationError

// func (v ValidationErrors) - stack trace representation.
func (v ValidationErrors) Error() string {
	result := ""
	for i, e := range v {
		if i == 0 {
			result = "Validation Errors stack trace:"
		}
		result = fmt.Sprintf("%s\n - %s: %s", result, e.Field, e.Err)
	}
	return result
}

// ValidationRule - struct of field and its params for validate.
type ValidationRule struct {
	Field  reflect.StructField
	Value  reflect.Value
	Func   string
	Params string
}

// ValidatorFunctionsMap - special public validator-map.
// This is map of validator functions with special signature.
var ValidatorFunctionsMap = map[string]ValidatorRuleFunctionSignature{
	"in":     validateIn,
	"len":    validateLen,
	"max":    validateMax,
	"min":    validateMin,
	"regexp": validateRegexp,
}

// ValidateFieldByRule - validate structure field by rule.
func ValidateFieldByRule(vRule ValidationRule) (validationError ValidationError) {
	validationError.Field = vRule.Field.Name
	validatorFunction, ok := ValidatorFunctionsMap[vRule.Func]
	if ok {
		validationError.Err = validatorFunction(vRule)
	} else {
		validationError.Err = tryValidateNotImplementedRule(vRule)
	}
	return validationError
}

// ValidateField - validate structure field at rule-by-rule order.
// It use validation rule grammatic "function:params".
// Params:
//   - field - is struct field object;
//   - value - is field value;
//   - validationsRules - is string representation of field limitation, like "in:200,404,500".
func ValidateField(field reflect.StructField, value reflect.Value, validationsRules string) ValidationErrors {
	fieldValidationErrors := make(ValidationErrors, 0)
	for _, validatorFieldRule := range strings.Split(validationsRules, "|") {
		// fmt.Printf(" > validate rule %q\n", validatorFieldRule)
		validatorFieldRuleFuncParams := strings.SplitN(validatorFieldRule, ":", 2)
		validatorFieldRuleFunc := validatorFieldRuleFuncParams[0]
		validatorFieldRuleParams := validatorFieldRuleFuncParams[1]
		if field.Type.Kind() == reflect.Slice {
			// fmt.Printf("field.Type.Kind() = %s\n", field.Type.Kind())
			// fmt.Printf("r.Slice = %s\n", r.Slice)
			for i := 0; i < value.Len(); i++ {
				iValue := value.Index(i)
				// log.Printf("iValue = %s\n", iValue)
				// fmt.Printf("iValue.Kind() = %s\n", iValue.Kind())
				virtualField := reflect.StructField{
					Name:      field.Name,
					PkgPath:   field.PkgPath,
					Type:      iValue.Type(),
					Tag:       field.Tag,
					Offset:    field.Offset,
					Index:     field.Index,
					Anonymous: field.Anonymous,
				}
				vRule := ValidationRule{virtualField, iValue, validatorFieldRuleFunc, validatorFieldRuleParams}
				fieldRuleValidationError := ValidateFieldByRule(vRule)
				if fieldRuleValidationError.Err != nil {
					fieldValidationErrors = append(fieldValidationErrors, fieldRuleValidationError)
				}
			}
		} else {
			vRule := ValidationRule{field, value, validatorFieldRuleFunc, validatorFieldRuleParams}
			fieldRuleValidationError := ValidateFieldByRule(vRule)
			if fieldRuleValidationError.Err != nil {
				fieldValidationErrors = append(fieldValidationErrors, fieldRuleValidationError)
			}
		}
	}
	return fieldValidationErrors
}

// Validate - validate structure at field-by-field order.
//
// It use `validator` struct field tag with validation rule grammatic like:
//   - "function:params".
//   - "function1:params1|function2:params2".
//
// Examples:
//
//	type App struct {
//	  Version string   `validate:"len:5"`
//	}
//
//	type Response struct {
//	  Code int         `validate:"in:200,404,500"`
//	  Body string      `json:"omitempty"`
//	}
//
//	type User struct {
//	  ID     string    `json:"id" validate:"len:36"`
//	  Name   string
//	  Age    int       `validate:"min:18|max:50"`
//	  Email  string    `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
//	  Role   UserRole  `validate:"in:admin,stuff"`
//	  Phones []string  `validate:"len:11"`
//	}
//
//	type LevelMonitoring struct {
//	  Samples []int    `validate:"min:20|max:40"`
//	}
//
//	type Block struct {
//	  Chain string `validate:"regexp:\\d+|len:20"`
//	}
func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("object type %q can not be validate, because it is not a Struct", t.Name())
	}
	structValidationErrors := make(ValidationErrors, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		mustBeValidateAs := field.Tag.Get(validationTagName)
		if mustBeValidateAs == "" {
			fmt.Printf("Let's skip field\tNo.%d:\n", field.Index)
			// WTF - linter!!!
			// golangci-lint run --out-format=github-actions
			// ...line is 121 characters (with 111 limit)
			fieldValue := reflect.ValueOf(v).FieldByName(field.Name)
			fmt.Printf(" - %v (%v) with value '%v'.\n", field.Name, field.Type.Name(), fieldValue)
			fmt.Printf(" - No any validator.\n")
			continue
		}
		fmt.Printf("Let's validate field\tNo.%d:\n", field.Index)
		fmt.Printf(" - %v (%v) with value '%v'\n", field.Name, field.Type.Name(), reflect.ValueOf(v).FieldByName(field.Name))
		fmt.Printf(" - It must satisfy validator '%v'.\n", mustBeValidateAs)
		value := reflect.ValueOf(v).FieldByName(field.Name)
		fieldValidationErrors := ValidateField(field, value, mustBeValidateAs)
		if len(fieldValidationErrors) > 0 {
			structValidationErrors = append(structValidationErrors, fieldValidationErrors...)
		}
	}
	// return structValidationErrors // Wow, it is normal return because structValidationErrors implenemt Error() interface
	if len(structValidationErrors) == 0 {
		return nil
	}
	return fmt.Errorf(structValidationErrors.Error())
}

type gString interface {
	string
}

const gStringType = "gString"

type gInt interface {
	int | int8 | int16 | int32 | int64
}
type gOrdered interface {
	gInt | gString
}

// ValidatorRuleFunctionSignature - signature of function, witch validate field by rule.
// It is used in special public validator-map.
type ValidatorRuleFunctionSignature func(vRule ValidationRule) error

//

func tryValidateNotImplementedRule(vRule ValidationRule) error {
	// WTF - linter!!!
	// golangci-lint run --out-format=github-actions
	// ...line is 125 characters (with 111 limit)
	f := vRule.Func
	t := vRule.Field.Type.Kind().String()
	return fmt.Errorf("validate function %q for field type %q is not implemented", f, t)
}

// [T gString | gSlice]

func checkLen[T gString](object T, expectedLen int) error {
	objectLen := len(object)
	if objectLen == expectedLen {
		return nil
	}
	return fmt.Errorf("field value len(\"%v\") return %d, but expected %d", object, objectLen, expectedLen)
}

func validateLen(vRule ValidationRule) error {
	if vRule.Field.Type.Kind() == reflect.String {
		parsedLen, err := strconv.Atoi(vRule.Params)
		if err != nil {
			return err
		}
		return checkLen[string](vRule.Value.String(), parsedLen)
	}
	return tryValidateNotImplementedRule(vRule)
}

// Lets overhead from - infimum INT

func checkMin[T gOrdered](object, infimum T) error {
	if object >= infimum {
		return nil
	}
	if reflect.TypeOf(object).Name() == gStringType { // Just for output formatting
		return fmt.Errorf("value %q is less than infimum %q", object, infimum)
	}
	return fmt.Errorf("value %v is less than infimum %v", object, infimum)
}

func validateMin(vRule ValidationRule) error {
	if vRule.Field.Type.Kind() == reflect.String {
		return checkMin[string](vRule.Value.String(), vRule.Params)
	}
	if vRule.Field.Type.Kind() == reflect.Int {
		parsedMin, err := strconv.Atoi(vRule.Params)
		if err != nil {
			return err
		}
		return checkMin[int64](vRule.Value.Int(), int64(parsedMin))
	}
	return tryValidateNotImplementedRule(vRule)
}

// Lets overhead from - supremum INT

func checkMax[T gOrdered](object, supremum T) error {
	if object <= supremum {
		return nil
	}
	if reflect.TypeOf(object).Name() == gStringType { // Just for output formatting
		return fmt.Errorf("value %q is less than infimum %q", object, supremum)
	}
	return fmt.Errorf("value %v is more than supremum %v", object, supremum)
}

func validateMax(vRule ValidationRule) error {
	if vRule.Field.Type.Kind() == reflect.String {
		return checkMax[string](vRule.Value.String(), vRule.Params)
	}
	if vRule.Field.Type.Kind() == reflect.Int {
		parsedMin, err := strconv.Atoi(vRule.Params)
		if err != nil {
			return err
		}
		return checkMax[int64](vRule.Value.Int(), int64(parsedMin))
	}
	return tryValidateNotImplementedRule(vRule)
}

//

func checkIn[T comparable](object T, collection []T) error {
	for _, v := range collection {
		if v == object {
			return nil
		}
	}
	if reflect.TypeOf(object).Name() == gStringType { // Just for output formatting
		return fmt.Errorf("value \"%v\" not in %v", object, collection)
	}
	return fmt.Errorf("value %v not in %v", object, collection)
}

func validateIn(vRule ValidationRule) error {
	if vRule.Field.Type.Kind() == reflect.String {
		return checkIn[string](vRule.Value.String(), strings.Split(vRule.Params, ","))
	}
	if vRule.Field.Type.Kind() == reflect.Int {
		intArray, err := sliceAtoi(strings.Split(vRule.Params, ","))
		if err != nil {
			return err
		}
		return checkIn[int64](vRule.Value.Int(), intArray) // TODO:
	}
	return tryValidateNotImplementedRule(vRule)
}

//

func checkRegexp(s, regexpTemplate string) error {
	re, err := regexp.Compile(regexpTemplate)
	if err != nil {
		return err
	}
	if re.MatchString(s) {
		return nil
	}
	return fmt.Errorf("value %q is not mutch to template %q", s, regexpTemplate)
}

func validateRegexp(vRule ValidationRule) error {
	if vRule.Field.Type.Kind() == reflect.String {
		return checkRegexp(vRule.Value.String(), vRule.Params)
	}
	return tryValidateNotImplementedRule(vRule)
}

//
