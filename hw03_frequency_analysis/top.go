// Проект с домашеней работой №3 курса OTUS-Go-2023-03.
package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// StringSpliter - разбивает строку на слова.
//
// Например:
//
//   - StringSpliter("qwe asd zxc qwe") = ["qwe" "asd" "zxc" "qwe"].
func StringSpliter(str string) []string {
	return strings.Fields(str)
}

// StructWord - структура, описывающая частоту встречаемости слова.
type StructWord struct {
	word  string // Слово
	count uint   // Частота встречаемости
}

// InitDistributionOfMappedWord - метод подсчета частоты встречаемости слов с результатом в виде MAP-значения:
// "слово № 1" : частота слова № 1 , ... , "слово № k" : частота слова № k.
//
// Например:
//
//   - InitDistributionOfMappedWord([]) = map[].
//   - InitDistributionOfMappedWord([x]) = map[x:1].
//   - InitDistributionOfMappedWord([x x]) = map[x:2].
//   - InitDistributionOfMappedWord([y x x]) = map[x:2 y:1].
func InitDistributionOfMappedWord(words []string) map[string]uint {
	distributionOfMappedWord := make(map[string]uint)
	for _, word := range words {
		if word == "" {
			continue
		}
		distributionOfMappedWord[word]++
	}
	return distributionOfMappedWord
}

// convertToDistributionOfStructWord - метод представления MAP-значения частоты слова в виде StructWord-слайса.
func convertToDistributionOfStructWord(distributionOfMappedWord map[string]uint) []StructWord {
	structWords := []StructWord{}
	for stringWord, count := range distributionOfMappedWord {
		if stringWord == "" {
			continue
		}
		structWords = append(structWords, StructWord{stringWord, count})
	}
	return structWords
}

// GetAsSortedStructWords - метод представления MAP-значения частоты слова в виде
// лексикографически упорядоченного StructWord-слайса.
//
// Например:
//
//   - GetAsSortedStructWords(map[]) = [].
//   - GetAsSortedStructWords(map[x:1]) = [{x 1}].
//   - GetAsSortedStructWords(map[x:2]) = [{x 2}].
//   - GetAsSortedStructWords(map[a:3 x:2 y:1 z:3]) = [{a 3}{z 3}{x 2}{y 1}].
func GetAsSortedStructWords(distributionOfMappedWord map[string]uint) []StructWord {
	// Сортировка в условии задачи: `ORDER BY count DESC, lexical ASC`
	structWords := convertToDistributionOfStructWord(distributionOfMappedWord)
	sort.Slice(
		structWords,
		func(i, j int) bool {
			return structWords[i].count > structWords[j].count || (structWords[i].count == structWords[j].count &&
				structWords[i].word < structWords[j].word)
		})
	return structWords
}

// GetTopStructWords - метод выборки первых N по очереди элементов с защитой от `slice bounds out of range`.
//
// Например:
//
//   - GetTopStructWords([], 10) = [].
//   - GetTopStructWords([{d 3}{e 3}{v 2}{h 1}], 0) = [].
//   - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 5) = [{a 3}{z 3}{x 2}{y 1}].
//   - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 2) = [{a 3}{z 3}].
//   - GetTopStructWords([{a 3}{z 3}{x 2}{y 1}], 1) = [{a 3}].
func GetTopStructWords(structWords []StructWord, limit uint) []StructWord {
	if uint(len(structWords)) < limit {
		limit = uint(len(structWords))
	}
	return structWords[:limit]
}

// WordStructToWordStrings - метод конвертации StructWord-слайса в string-слайс.
//
// Например:
//
//   - WordStructToWordStrings([]) = [].
//   - WordStructToWordStrings([{ad 3} {ae 3} {vj 2} {h 1}]) = [ad ae vj h].
//   - WordStructToWordStrings([{ccc 10} {aaa 3} {b 1}]) = [ccc aaa b].
func WordStructToWordStrings(structWords []StructWord) []string {
	stringWords := []string{}
	for _, structWord := range structWords {
		stringWords = append(stringWords, structWord.word)
	}
	return stringWords
}

// Top10 - функция с заданной в домашнем задании сигнатурой.
//
// Реализована в виде последовательного вызова разработанных методов,
// максимально декомпозирующих исходную задачу на отдельные атомарные этапы,
// без дублирования кода.
// По моему мнению не содержит в себе ничего лишнего, необходимого сверх решения поставленной задачи.
func Top10(s string) []string {
	stringWords := StringSpliter(s)
	mappedWords := InitDistributionOfMappedWord(stringWords)
	structWords := GetAsSortedStructWords(mappedWords)
	top10StructWords := GetTopStructWords(structWords, 10)
	top10StringWords := WordStructToWordStrings(top10StructWords)
	return top10StringWords
}
