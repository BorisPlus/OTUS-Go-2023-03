# Домашнее задание №9 «Валидатор структур»

Описание [задания](./README.md).

## Реализации

```shell
go doc -all ./ > go_doc_-all.txt
```

```text
{{ go_doc_-all.txt }}
```

## Тестирование

### Тестирование на структурах с валидными значениями

```shell
go test -v -run TestValidatePositive ./ > TestValidatePositive.txt
```

```text
{{ TestValidatePositive.txt }}
```

### Тестирование на структурах с невалидными значениями

```shell
go test -v -run TestValidateNegative ./ > TestValidateNegative.txt
```

```text
{{ TestValidateNegative.txt }}
```

### Тестирование не на структуре

```shell
go test -v -run TestValidateNotStructObject ./ > TestValidateNotStructObject.txt
```

```text
{{ TestValidateNotStructObject.txt }}
```

### Тестирование на структуре с нереализованными валидаторами

```shell
go test -v -run TestValidateNotImplemented ./ > TestValidateNotImplemented.txt
```

```text
{{ TestValidateNotImplemented.txt }}
```

### Ожидаемый стек ошибок на примере нереализованных валидаторов

```shell
go test -v -run TestValidateExpectedNotImplemented ./ > TestValidateExpectedNotImplemented.txt
```

```text
{{ TestValidateExpectedNotImplemented.txt }}
```

### Отсутствие неожидаемого стека ошибок на примере нереализованных валидаторов

> Важен порядок формирования стека ошибок валидации. В сравнении с прошлым примером он просто изменен в ожидаемом выводе.

```shell
go test -v -run TestValidateUnxpectedNotImplemented ./ > TestValidateUnxpectedNotImplemented.txt
```

```text
{{ TestValidateUnxpectedNotImplemented.txt }}
```

## Вывод

Реализован тег упрощенного параметризуемого валидатора произвольного поля структуры.

## Заметка для себя (составление отчета)

```shell
golangci-lint run --out-format=github-actions ./

cd ../hw09_struct_validator
go doc -all ./ > go_doc_-all.txt
go test -v -run TestValidatePositive ./ > TestValidatePositive.txt
go test -v -run TestValidateNegative ./ > TestValidateNegative.txt
go test -v -run TestValidateNotStructObject ./ > TestValidateNotStructObject.txt
go test -v -run TestValidateNotImplemented ./ > TestValidateNotImplemented.txt
go test -v -run TestValidateExpectedNotImplemented ./ > TestValidateExpectedNotImplemented.txt
go test -v -run TestValidateUnxpectedNotImplemented ./ > TestValidateUnxpectedNotImplemented.txt
cd ../report_templator
go test templator.go hw09_struct_validator_test.go
cd ../hw09_struct_validator

cat ./REPORT.md | grep FAIL
```
