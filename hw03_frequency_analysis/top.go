package hw03frequencyanalysis

import (
	"sort"
	"strings"
)


// Разбиение строки на слова
// Tests:
// OK. StringSpliter("qwe asd zxc qwe") return expected ["qwe" "asd" "zxc" "qwe"].
func StringSpliter(str string) []string {
	// Все слова без иных пробельных символов
	whitespace := " "
	for _, spacevalue := range []string{"\n", "\t", "\r", "\f", "\v"} {
		str = strings.ReplaceAll(str, spacevalue, whitespace)
	}
	return strings.Split(str, whitespace)
}

// Частота встречаемости слов
// в виде структуры.
type StructWord struct {
	word  string // Слово
	count uint   // Частота встречаемости
}

// Частота встречаемости слов
// в виде map.
type DistributionOfMappedWord struct {
	Distribution map[string]uint
}

// Метод получения распределения слов в виде Мap.
// Tests:
// OK. InitDistributionOfMappedWord([]) return expected map[].
// OK. InitDistributionOfMappedWord([x]) return expected map[x:1].
// OK. InitDistributionOfMappedWord([x x]) return expected map[x:2].
// OK. InitDistributionOfMappedWord([y x x]) return expected map[x:2 y:1].
func InitDistributionOfMappedWord(words []string) DistributionOfMappedWord {
	distributionOfMappedWord := DistributionOfMappedWord{make(map[string]uint)}
	for _, word := range words {
		if word == "" {
			continue
		}
		distributionOfMappedWord.Distribution[word]++
	}
	return distributionOfMappedWord
}

// Метод представления Мap в виде []WordObject.
func (distributionOfMappedWord *DistributionOfMappedWord) convertToDistributionOfStructWord() []StructWord {
	structWords := []StructWord{}
	for stringWord, count := range distributionOfMappedWord.Distribution {
		if stringWord == "" {
			continue
		}
		structWords = append(structWords, StructWord{stringWord, count})
	}
	return structWords
}

// Tests:
// OK. (map[]).GetAsSortedStructWords() return expected [].
// OK. (map[x:1]).GetAsSortedStructWords() return expected [{x 1}].
// OK. (map[x:2]).GetAsSortedStructWords() return expected [{x 2}].
// OK. (map[a:3 x:2 y:1 z:3]).GetAsSortedStructWords() return expected [{a 3} {z 3} {x 2} {y 1}].
func (distributionOfMappedWord *DistributionOfMappedWord) GetAsSortedStructWords() []StructWord {
	// Сортировка в условии задачи: `ORDER BY count DESC, lexical ASC`
	structWords := distributionOfMappedWord.convertToDistributionOfStructWord()
	sort.Slice(
		structWords,
		func(i, j int) bool {
			return structWords[i].count > structWords[j].count || (structWords[i].count == structWords[j].count &&
				structWords[i].word < structWords[j].word)
		})
	return structWords
}

// Tests:
// OK. GetTopStructWords([], 10) return expected [].
// OK. GetTopStructWords([{d 3} {e 3} {v 2} {h 1}], 0) return expected [].
// OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 5) return expected [{a 3} {z 3} {x 2} {y 1}].
// OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 2) return expected [{a 3} {z 3}].
// OK. GetTopStructWords([{a 3} {z 3} {x 2} {y 1}], 1) return expected [{a 3}].
func GetTopStructWords(structWords []StructWord, limit uint) []StructWord {
	min := func(x, y uint) uint {
		if x < y {
			return x
		}
		return y
	}
	wordStructsCount := uint(len(structWords))
	return structWords[:min(wordStructsCount, limit)]
}

// Tests:
// OK. WordStructToWordStrings([]) return expected [].
// OK. WordStructToWordStrings([{ad 3} {ae 3} {vj 2} {h 1}]) return expected [ad ae vj h].
// OK. WordStructToWordStrings([{ccc 10} {aaa 3} {b 1}]) return expected [ccc aaa b].
func WordStructToWordStrings(structWords []StructWord) []string {
	stringWords := []string{}
	for _, structWord := range structWords {
		stringWords = append(stringWords, structWord.word)
	}
	return stringWords
}

func Top10(s string) []string {
	stringWords := StringSpliter(s)
	mappedWords := InitDistributionOfMappedWord(stringWords)
	structWords := mappedWords.GetAsSortedStructWords()
	top10StructWords := GetTopStructWords(structWords, 10)
	top10StringWords := WordStructToWordStrings(top10StructWords)
	return top10StringWords
}
