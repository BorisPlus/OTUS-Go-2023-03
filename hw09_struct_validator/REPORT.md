# Домашнее задание №9 «Валидатор структур»

Описание [задания](./README.md).

## Реализации

```shell
go doc -all ./ > go_doc_-all.txt
```

```text
package hw09structvalidator // import "github.com/BorisPlus/hw09_struct_validator"


VARIABLES

var ValidatorFunctionsMap = map[string]ValidatorRuleFunctionSignature{
    "in":     validateIn,
    "len":    validateLen,
    "max":    validateMax,
    "min":    validateMin,
    "regexp": validateRegexp,
}
    ValidatorFunctionsMap - special public validator-map. This is map of
    validator functions with special signature.


FUNCTIONS

func Validate(v interface{}) error
    Validate - validate structure at field-by-field order.

    It use `validator` struct field tag with validation rule grammatic like:
      - "function:params".
      - "function1:params1|function2:params2".

    Examples:

        type App struct {
          Version string   `validate:"len:5"`
        }

        type Response struct {
          Code int         `validate:"in:200,404,500"`
          Body string      `json:"omitempty"`
        }

        type User struct {
          ID     string    `json:"id" validate:"len:36"`
          Name   string
          Age    int       `validate:"min:18|max:50"`
          Email  string    `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
          Role   UserRole  `validate:"in:admin,stuff"`
          Phones []string  `validate:"len:11"`
        }

        type LevelMonitoring struct {
          Samples []int    `validate:"min:20|max:40"`
        }


TYPES

type ValidationError struct {
    Field string
    Err   error
}
    ValidationError - struct of validation error.

func ValidateFieldByRule(vRule ValidationRule) (validationError ValidationError)
    ValidateFieldByRule - validate structure field by rule.

type ValidationErrors []ValidationError
    ValidationErrors - stack trace of of validation errors.

func ValidateField(field reflect.StructField, value reflect.Value, validationsRules string) ValidationErrors
    ValidateField - validate structure field at rule-by-rule order. It use
    validation rule grammatic "function:params". Params:
      - field - is struct field object;
      - value - is field value;
      - validationsRules - is string representation of field limitation,
        like "in:200,404,500".

func (v ValidationErrors) Error() string
    func (v ValidationErrors) - stack trace representation.

type ValidationRule struct {
    Field  reflect.StructField
    Value  reflect.Value
    Func   string
    Params string
}
    ValidationRule - struct of field and its params for validate.

type ValidatorRuleFunctionSignature func(vRule ValidationRule) error
    ValidatorRuleFunctionSignature - signature of function, witch validate field
    by rule. It is used in special public validator-map.


```

## Тестирование

### Тестирование на структурах с валидными значениями

```shell
go test -v -run TestValidatePositive ./ > TestValidatePositive.txt
```

```text
=== RUN   TestValidatePositive
=== RUN   TestValidatePositive/User_positive_validation
Let's validate field    No.[0]:
 - ID (string) with value '012345678901234567890123456789012345'
 - It must satisfy validator 'len:36'.
Let's skip field    No.[1]:
 - Name (string) with value ''.
 - No any validator.
Let's validate field    No.[2]:
 - Age (int) with value '27'
 - It must satisfy validator 'min:18|max:50'.
Let's validate field    No.[3]:
 - Email (string) with value 'stuff@go.dev'
 - It must satisfy validator 'regexp:^\w+@\w+\.\w+$'.
Let's validate field    No.[4]:
 - Role (UserRole) with value 'stuff'
 - It must satisfy validator 'in:admin,stuff'.
Let's validate field    No.[5]:
 - Phones () with value '[+1012345678 +7701234568]'
 - It must satisfy validator 'len:11'.
Let's skip field    No.[6]:
 - meta (RawMessage) with value '[]'.
 - No any validator.
=== RUN   TestValidatePositive/App_positive_validation
Let's validate field    No.[0]:
 - Version (string) with value 'v.1.9'
 - It must satisfy validator 'len:5'.
=== RUN   TestValidatePositive/Token_positive_validation
Let's skip field    No.[0]:
 - Header () with value '[]'.
 - No any validator.
Let's skip field    No.[1]:
 - Payload () with value '[]'.
 - No any validator.
Let's skip field    No.[2]:
 - Signature () with value '[]'.
 - No any validator.
=== RUN   TestValidatePositive/Response_positive_validation
Let's validate field    No.[0]:
 - Code (int) with value '200'
 - It must satisfy validator 'in:200,404,500'.
Let's skip field    No.[1]:
 - Body (string) with value 'OK HTTP/2'.
 - No any validator.
=== RUN   TestValidatePositive/LevelMonitoring_positive_validation
Let's validate field    No.[0]:
 - Samples () with value '[20 35 27 40]'
 - It must satisfy validator 'min:20|max:40'.
--- PASS: TestValidatePositive (0.00s)
    --- PASS: TestValidatePositive/User_positive_validation (0.00s)
    --- PASS: TestValidatePositive/App_positive_validation (0.00s)
    --- PASS: TestValidatePositive/Token_positive_validation (0.00s)
    --- PASS: TestValidatePositive/Response_positive_validation (0.00s)
    --- PASS: TestValidatePositive/LevelMonitoring_positive_validation (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

### Тестирование на структурах с невалидными значениями

```shell
go test -v -run TestValidateNegative ./ > TestValidateNegative.txt
```

```text
=== RUN   TestValidateNegative
=== RUN   TestValidateNegative/User_negative_validation
Let's validate field    No.[0]:
 - ID (string) with value 'I have short id'
 - It must satisfy validator 'len:36'.
Let's skip field    No.[1]:
 - Name (string) with value ''.
 - No any validator.
Let's validate field    No.[2]:
 - Age (int) with value '14'
 - It must satisfy validator 'min:18|max:50'.
Let's validate field    No.[3]:
 - Email (string) with value 'stuff@go.dev.ru'
 - It must satisfy validator 'regexp:^\w+@\w+\.\w+$'.
Let's validate field    No.[4]:
 - Role (UserRole) with value 'guest'
 - It must satisfy validator 'in:admin,stuff'.
Let's validate field    No.[5]:
 - Phones () with value '[+10123456789 +77012345689]'
 - It must satisfy validator 'len:11'.
Let's skip field    No.[6]:
 - meta (RawMessage) with value '[]'.
 - No any validator.
Test catch expected error. It's OK.
Validation Errors stack trace:
 - ID: field value len("I have short id") return 15, but expected 36
 - Age: value 14 is less than infimum 18
 - Email: value "stuff@go.dev.ru" is not mutch to template "^\\w+@\\w+\\.\\w+$"
 - Role: value guest not in [admin stuff]
 - Phones: field value len("+10123456789") return 12, but expected 11
 - Phones: field value len("+77012345689") return 12, but expected 11
=== RUN   TestValidateNegative/LevelMonitoring_negative_validation
Let's validate field    No.[0]:
 - Samples () with value '[20 45 27 15]'
 - It must satisfy validator 'min:20|max:40'.
Test catch expected error. It's OK.
Validation Errors stack trace:
 - Samples: value 15 is less than infimum 20
 - Samples: value 45 is more than supremum 40
--- PASS: TestValidateNegative (0.00s)
    --- PASS: TestValidateNegative/User_negative_validation (0.00s)
    --- PASS: TestValidateNegative/LevelMonitoring_negative_validation (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

### Тестирование не на структуре

```shell
go test -v -run TestValidateNotStructObject ./ > TestValidateNotStructObject.txt
```

```text
=== RUN   TestValidateNotStructObject
Test catch expected error. It's OK.
object type "string" can not be validate, because it is not a Struct
--- PASS: TestValidateNotStructObject (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

### Тестирование на структуре с нереализованными валидаторами

```shell
go test -v -run TestValidateNotImplemented ./ > TestValidateNotImplemented.txt
```

```text
=== RUN   TestValidateNotImplemented
Let's validate field    No.[0]:
 - A (int) with value '1'
 - It must satisfy validator 'len:1'.
Let's validate field    No.[1]:
 - B (string) with value 'fieldB'
 - It must satisfy validator 'length:1'.
Test catch expected error. It's OK.
Validation Errors stack trace:
 - A: validate function "len" for field type "int" is not implemented
 - B: validate function "length" for field type "string" is not implemented
--- PASS: TestValidateNotImplemented (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

### Ожидаемый стек ошибок на примере нереализованных валидаторов

```shell
go test -v -run TestValidateExpectedNotImplemented ./ > TestValidateExpectedNotImplemented.txt
```

```text
=== RUN   TestValidateExpectedNotImplemented
Let's validate field    No.[0]:
 - A (int) with value '1'
 - It must satisfy validator 'len:1'.
Let's validate field    No.[1]:
 - B (string) with value 'fieldB'
 - It must satisfy validator 'length:1'.
Test catch expected error. It's OK.
Validation Errors stack trace:
 - A: validate function "len" for field type "int" is not implemented
 - B: validate function "length" for field type "string" is not implemented

--- PASS: TestValidateExpectedNotImplemented (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

### Отсутствие неожидаемого стека ошибок на примере нереализованных валидаторов

> Важен порядок формирования стека ошибок валидации. В сравнении с прошлым примером он просто изменен в ожидаемом выводе.

```shell
go test -v -run TestValidateUnxpectedNotImplemented ./ > TestValidateUnxpectedNotImplemented.txt
```

```text
=== RUN   TestValidateUnxpectedNotImplemented
Let's validate field    No.[0]:
 - A (int) with value '1'
 - It must satisfy validator 'len:1'.
Let's validate field    No.[1]:
 - B (string) with value 'fieldB'
 - It must satisfy validator 'length:1'.
Test catch expected fake error. It's OK.

Get error. It's OK.
Validation Errors stack trace:
 - A: validate function "len" for field type "int" is not implemented
 - B: validate function "length" for field type "string" is not implemented

Expected fake error. It's OK.
Validation Errors stack trace:
 - B: validate function "length" for field type "string" is not implemented
 - A: validate function "len" for field type "int" is not implemented

--- PASS: TestValidateUnxpectedNotImplemented (0.00s)
PASS
ok      github.com/BorisPlus/hw09_struct_validator    (cached)

```

## Вывод

Реализован тег упрощенного параметризуемого валидатора.

## Для составления отчета (для себя)

```shell
golangci-lint run --out-format=github-actions ./

cd ../hw09_struct_validator
go doc -all ./ > go_doc_-all.txt &&
go test -v -run TestValidatePositive ./ > TestValidatePositive.txt &&
go test -v -run TestValidateNegative ./ > TestValidateNegative.txt &&
go test -v -run TestValidateNotStructObject ./ > TestValidateNotStructObject.txt &&
go test -v -run TestValidateNotImplemented ./ > TestValidateNotImplemented.txt &&
go test -v -run TestValidateExpectedNotImplemented ./ > TestValidateExpectedNotImplemented.txt &&
go test -v -run TestValidateUnxpectedNotImplemented ./ > TestValidateUnxpectedNotImplemented.txt &&
cd ../report_templator &&
go test templator.go hw09_struct_validator_test.go &&
cd ../hw09_struct_validator
```
