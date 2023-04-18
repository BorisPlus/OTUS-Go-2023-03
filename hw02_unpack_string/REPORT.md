# Отчет о подходе к решению задачи "распаковки" строки

## Идея

Задача обработки строки в соотвествии с изложенными правилами "распаковки" подпадает удовлетворяет подходу [`порождающей грамматики`](https://ru.wikipedia.org/wiki/%D0%9F%D0%BE%D1%80%D0%BE%D0%B6%D0%B4%D0%B0%D1%8E%D1%89%D0%B0%D1%8F_%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0) и сводится к составлению последовательной цепочки (`предложения`) блоков слов (`лексем`).

`Лексема` представляет собой структуру из двух полей:

- Руна (перечень ограничивается прилагаемым в коде списком).
- Число повторений Руны. При этом число представлено только одной цифрой. Если число повторений не задано, то оно равно 1.

"Распаковка" `лексемы` приводит к "порождению" `слова`.

`Блок предложения` задается совокупностью из изначальной `лексемы` и ссылки на следующий `блок предложения`.
"Распаковка" `блока предложения` приводит к "распаковке" `лексемы`.

"Распаковка" `предложения` заключается в последовательной распаковке `блоков` по цепочке, начиная с первого в цепочке.
Если ссылка на следующий `блок предложения` нулевая (`nil`), то всё `предложение` распаковано.

### Лексема

[lexeme.go](./OTUS-Go-2023-03/hw02_unpack_string/lexeme.go)

```go
type Lexeme struct {
    _rune        rune
    _runeWasSet  bool // необходимо для отсеживания факта присвоения
    _count       uint
    _countWasSet bool // необходимо для отсеживания факта присвоения, так как
                      // если число повторений не задано, то _count равно 1
}
```

### Блок предложения

```go
type StatementBlock struct {
    BlockLexeme Lexeme     // текущая лексема
    NextBlock   *Statement // ссылка на следующий блок
}
```

### Предложение

Задается стратовым блоком `FirstStatementBlock` всей цепочки, схематично:

```text
FirstStatementBlock --> *NextStatementBlock --> ... --> *NextStatementBlock --> nil
```

### Распаковка

Алгоритм заключается в постороении вышеуказанной цепочки `блоков предложения`:

```go
func Unpack(inputString string) (string, error) {
    ...
    var currentBlockTmp StatementBlock = StatementBlock{BlockLexeme: lexeme, NextBlock: currentBlock}
    currentBlock = &currentBlockTmp
    lexeme = Lexeme{}
    ...
}
```

__Замечание__: Особенность в том, что обход исходной строки проводится в обратном порядке - с конца в начало. Проще строить цепочку, определяя указатель на следующий `блок предложения`, так как именно он был распознан на этапе итерации. Кроме того, результирующий программный код при реверсивном обходе становится менее громоздким, а значит более читабельным и наглядным.

## Тестирование реализации

```shell
cd hw02_unpack_string/
go test -v lexeme.go statement.go unpack.go lexeme_test.go statement_test.go unpack_test.go 
```

Рассмотрим тесты поподробнее.

### Тестирование структуры `Лексема`

```shell
go test -v lexeme.go lexeme_test.go 
```

Вывод:

```text
=== RUN   TestLexeme             --,
Lexeme 'a'*0 is valid. It's OK.    |
Lexeme 'a'*1 is valid. It's OK.    |--> Проверка допустимых значений
Lexeme 'a'*9 is valid. It's OK.    |    Рун и Числа их повторений в Лексеме
--- PASS: TestLexeme (0.00s)     --`
=== RUN   TestLexemePanic)          --,
Lexeme 'a'*10 is not valid. It's OK.  |--> Проверка на недопустимость определенных 
Lexeme 'B'*1 is not valid. It's OK.   |    Рун и Числа их повторений в Лексеме
--- PASS: TestLexemePanic (0.00s)   --`
=== RUN   TestLexemeUnpack                     --,
Lexeme 'a'*0 unpacked to: "".                    |
Lexeme 'a'*1 unpacked to: "a".                   |--> Проверка результатов распаковки Лексемы
Lexeme '\t'*9 unpacked to: "\t\t\t\t\t\t\t\t\t". |
--- PASS: TestLexemeUnpack (0.00s)             --`
PASS
ok      command-line-arguments  0.007s
```

### Тестирование структуры `Блок Предложения`

```shell
go test -v lexeme.go statement.go statement_test.go
```

Вывод:

```text
Iterate test for "a1b2c9de0" unpacked.
Lexemes for Statements blocks:                 --,
Lexeme 'a'*1 unpacked to: "a".                   |
Lexeme 'b'*2 unpacked to: "bb".                  |--> подготовка лексем в TestMain
Lexeme 'c'*9 unpacked to: "ccccccccc".           |
Lexeme 'd' without any count unpacked to: "d".   |
Lexeme 'e'*0 unpacked to: "".                  --`
=== RUN   TestStatementUnpack
Reverse realization.                                --,
Statement "e0" unpacked to: "".                       |
Statement "de0" unpacked to: "d".                     |
Statement "c9de0" unpacked to: "cccccccccd".          |--> Реверсивный обход строки "a1b2c9de0"
Statement "b2c9de0" unpacked to: "bbcccccccccd".      |    с визуализацией промежуточных итогов
Statement "a1b2c9de0" unpacked to: "abbcccccccccd".   |
--- PASS: TestStatementUnpack (0.00s)               --`
PASS
ok      command-line-arguments  0.004s
```

### Тестирование алгоритма распаковки строки

```shell
go test -v lexeme.go statement.go unpack.go unpack_test.go 
```

Вывод:

```text
=== RUN   TestUnpack                            --,
=== RUN   TestUnpack/a4bc2d5e                     |
String "a4bc2d5e" unpacked to: "aaaabccddddde".   |
=== RUN   TestUnpack/abccd                        |
String "abccd" unpacked to: "abccd".              |
=== RUN   TestUnpack/#00                          |
String "" unpacked to: "".                        |--> проверка результатов распаковки 
=== RUN   TestUnpack/aaa0b                        |    на допустимых строках
String "aaa0b" unpacked to: "aab".                |
=== RUN   TestUnpack/d_5abc  (!)тут: "d\n5abc"    |
String "d\n5abc" unpacked to: "d\n\n\n\n\nabc"    |
--- PASS: TestUnpack (0.00s)                      |
    --- PASS: TestUnpack/a4bc2d5e (0.00s)         |
    --- PASS: TestUnpack/abccd (0.00s)            |
    --- PASS: TestUnpack/#00 (0.00s)              |
    --- PASS: TestUnpack/aaa0b (0.00s)            |
    --- PASS: TestUnpack/d_5abc (0.00s)         --`
=== RUN   TestUnpackInvalidString
=== RUN   TestUnpackInvalidString/3abc
String "3abc" unpacked to: "", with error: "Not valid symbol '3' in position `0`.".       --,
=== RUN   TestUnpackInvalidString/45                                                        |
String "45" unpacked to: "", with error: "Not valid count-symbol '5' in position `1`.".     |
=== RUN   TestUnpackInvalidString/aaa10b                                                    |
String "aaa10b" unpacked to: "", with error: "Not valid count-symbol '0' in position `4`.". |
--- PASS: TestUnpackInvalidString (0.00s)                                                   /--> проверка ошибки  
    --- PASS: TestUnpackInvalidString/3abc (0.00s)                                         /     распаковки на 
    --- PASS: TestUnpackInvalidString/45 (0.00s)                                          /      недопустимых строках
    --- PASS: TestUnpackInvalidString/aaa10b (0.00s)                                   --`
PASS
ok      command-line-arguments  0.007s
```

## Вывод

Функция распаковки реализована в соотвествии с интерфейсом исходной задачи:

```go
func Unpack(_ string) (string, error) {
    // Place your code here.
    return "", nil
}
```

её тест также приведен в изначальном варианте, за исключением отказа от обобщенного вида ошибки

```go
var ErrInvalidString = errors.New("invalid string")
```

в угоду более конкретизированного варианта, в частности:

```text
... with error: "Not valid symbol '3' in position `0`.". 
```

Предлагаемый подход ориентирован на изначальную проверку корректности исходного светрочного выражения.
Только в случае допустимости всех `лексем` в постороенной цепочке `блоков предложения` будет произведена их последовательная распаковка.

Теоретически, в случае ошибки в исходном светрочном выражении вычислительных ресурсов будет затрачено меньше, чем при вариаенте "поточной" распаковки на лету (необходимы замеры времени исполнения обоих вариантов).
