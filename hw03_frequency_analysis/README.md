# Домашнее задание №3 «Частотный анализ»

Необходимо написать Go функцию, принимающую на вход строку с текстом и
возвращающую слайс с 10-ю наиболее часто встречаемыми в тексте словами.

Если слова имеют одинаковую частоту, то должны быть отсортированы **лексикографически**.

* Словом считается набор символов, разделенных пробельными символами.

* Если есть более 10 самых частотых слов (например 15 разных слов встречаются ровно 133 раза,
остальные < 100), то следует вернуть 10 лексикографически первых слов.

* Словоформы не учитываем: "нога", "ногу", "ноги" - это разные слова.

* Слово с большой и маленькой буквы считать за разные слова. "Нога" и "нога" - это разные слова.

* Знаки препинания считать "буквами" слова или отдельными словами.
"-" (тире) - это отдельное слово. "нога," и "нога" - это разные слова.

## Пример

```text
cat and dog, one dog,two cats and one man
```

Топ 7:

* `and`     (2)
* `one`     (2)
* `cat`     (1)
* `cats`    (1)
* `dog,`    (1)
* `dog,two` (1)
* `man`     (1)

При необходимости можно выделять дополнительные функции / ошибки.

## (*) Дополнительное задание 

* учитывать большие/маленькие буквы и знаки препинания
* "Нога" и "нога" - это одинаковые слова, "нога!", "нога", "нога," и " 'нога' " - это одинаковые слова;
* "какой-то" и "какойто" - это разные слова, "-" (тире) - это не слово.

## Критерии оценки

* Пайплайн зелёный - 4 балла
* Добавлены новые юнит-тесты - до 4 баллов
* Понятность и чистота кода - до 2 баллов
* Дополнительное задание на баллы не влияет

Зачёт от 7 баллов

## Подсказки

* `regexp.MustCompile`
* `strings.Split`
* `strings.Fields`
* `sort.Slice`

## Частые ошибки

* `regexp.MustCompile` используется в функции, а не уровне пакета - это плохо по следующим причинам:
  * производительность: нет смысла компилировать регулярку каждый раз при вызове функции;
  * функция не должна паниковать!
* При выполнении задания со звёздочкой забывают, что тире не должно являться словом.

## Демонстрация работоспособности

Стремление к покрытию кода тестами

```shell
go test -v -cover top.go top_test.go
```

```text
=== RUN   TestStringSpliter
OK. StringSpliter("qwe asd zxc qwe") return expected ["qwe" "asd" "zxc" "qwe"].
--- PASS: TestStringSpliter (0.00s)
=== RUN   TestInitDistributionOfMappedWord
OK. InitDistributionOfMappedWord([]) return expected map[].
OK. InitDistributionOfMappedWord([x]) return expected map[x:1].
OK. InitDistributionOfMappedWord([x x]) return expected map[x:2].
OK. InitDistributionOfMappedWord([y x x]) return expected map[x:2 y:1].
--- PASS: TestInitDistributionOfMappedWord (0.00s)
=== RUN   TestSortedStructWords
OK. (map[]).GetAsSortedStructWords() return expected [].
OK. (map[x:1]).GetAsSortedStructWords() return expected [{x 1}].
OK. (map[x:2]).GetAsSortedStructWords() return expected [{x 2}].
OK. (map[a:3 x:2 y:1 z:3]).GetAsSortedStructWords() return expected [{a 3} {z 3} {x 2} {y 1}].
--- PASS: TestSortedStructWords (0.00s)
=== RUN   TestGetTopStructWords
OK. GetTopStructWords([], 10) return expected [].
OK. GetTopStructWords([{d 3} {e 3} {v 2} {h 1}], 0) return expected [].
OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 5) return expected [{a 3} {z 3} {x 2} {y 1}].
OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 2) return expected [{a 3} {z 3}].
OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 1) return expected [{a 3}].
--- PASS: TestGetTopStructWords (0.00s)
=== RUN   TestWordStructToWordStrings
OK. WordStructToWordStrings([]) return expected [].
OK. WordStructToWordStrings([{ad 3} {ae 3} {vj 2} {h 1}]) return expected [ad ae vj h].
OK. WordStructToWordStrings([{ccc 10} {aaa 3} {b 1}]) return expected [ccc aaa b].
--- PASS: TestWordStructToWordStrings (0.00s)
=== RUN   TestTop10
=== RUN   TestTop10/no_words_in_empty_string
=== RUN   TestTop10/positive_test_at_HW_example
=== RUN   TestTop10/positive_test
=== RUN   TestTop10/negative_test
--- PASS: TestTop10 (0.00s)
    --- PASS: TestTop10/no_words_in_empty_string (0.00s)
    --- PASS: TestTop10/positive_test_at_HW_example (0.00s)
    --- PASS: TestTop10/positive_test (0.00s)
    --- PASS: TestTop10/negative_test (0.00s)
PASS
coverage: 97.2% of statements
ok      command-line-arguments  0.007s  coverage: 97.2% of statements
```

coverage: **97.2%** of statements

## Документирование

```text
Top10 - функция с заданной в домашнем задании сигнатурой.

    Реализована в виде последовательного вызова разработанных методов,
    максимально декомпозирующих исходную задачу на отдельные атомарные
    этапы, без дублирования кода. По моему мнению не содержит в себе ничего
    лишнего, необходимого сверх решения поставленной задачи.
```

<details> 
<summary> см. подробнее (получено с использованием go doc -all ./) </summary>

```text
package hw03frequencyanalysis // import "github.com/BorisPlus/OTUS-Go-2023-03/tree/master/hw03_frequency_analysis"

Проект с домашеней работой №3 курса OTUS-Go-2023-03.

FUNCTIONS

func InitDistributionOfMappedWord(words []string) map[string]uint
    InitDistributionOfMappedWord - метод подсчета частоты встречаемости слов
    с результатом в виде MAP-значения: "слово № 1" : частота слова № 1 , ... ,
    "слово № k" : частота слова № k.

    Например:

      - InitDistributionOfMappedWord([]) = map[].
      - InitDistributionOfMappedWord([x]) = map[x:1].
      - InitDistributionOfMappedWord([x x]) = map[x:2].
      - InitDistributionOfMappedWord([y x x]) = map[x:2 y:1].

func StringSpliter(str string) []string
    StringSpliter - разбивает строку на слова.

    Например:

      - StringSpliter("qwe asd zxc qwe") = ["qwe" "asd" "zxc" "qwe"].

func Top10(s string) []string
    Top10 - функция с заданной в домашнем задании сигнатурой.

    Реализована в виде последовательного вызова разработанных методов,
    максимально декомпозирующих исходную задачу на отдельные атомарные этапы,
    без дублирования кода. По моему мнению не содержит в себе ничего лишнего,
    необходимого сверх решения поставленной задачи.

func WordStructToWordStrings(structWords []StructWord) []string
    WordStructToWordStrings - метод конвертации StructWord-слайса в
    string-слайс.

    Например:

      - WordStructToWordStrings([]) = [].
      - WordStructToWordStrings([{ad 3} {ae 3} {vj 2} {h 1}]) = [ad ae vj h].
      - WordStructToWordStrings([{ccc 10} {aaa 3} {b 1}]) = [ccc aaa b].


TYPES

type StructWord struct {
  // Has unexported fields.
}
    StructWord - структура, описывающая частоту встречаемости слова.

func GetAsSortedStructWords(distributionOfMappedWord map[string]uint) []StructWord
    GetAsSortedStructWords - метод представления MAP-значения частоты слова в
    виде лексикографически упорядоченного StructWord-слайса.

    Например:

      - GetAsSortedStructWords(map[]) = [].
      - GetAsSortedStructWords(map[x:1]) = [{x 1}].
      - GetAsSortedStructWords(map[x:2]) = [{x 2}].
      - GetAsSortedStructWords(map[a:3 x:2 y:1 z:3]) = [{a 3}{z 3}{x 2}{y 1}].

func GetTopStructWords(structWords []StructWord, limit uint) []StructWord
    GetTopStructWords - метод выборки первых N по очереди элементов с защитой от
    `slice bounds out of range`.

    Например:

      - GetTopStructWords([], 10) = [].
      - GetTopStructWords([{d 3}{e 3}{v 2}{h 1}], 0) = [].
      - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 5) = [{a 3}{z 3}{x 2}{y 1}].
      - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 2) = [{a 3}{z 3}].
      - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 1) = [{a 3}].

```

</details>
