package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

var LIMIT = 10

// Частота встречаемости слова
// в виде структуры.
type WordObject struct {
	word  string // Слово
	count uint   // Частота встречаемости
}

// Частота встречаемости слов
// в виде map.
type DistributionMap struct {
	mappedStatValue map[string]uint
}

// Метод получения Мap из строки.
func (distMap *DistributionMap) Calculate(values []string) {
	for _, word := range values {
		if word == "" {
			continue
		}
		distMap.mappedStatValue[word]++
	}
}

// Метод представления Мap в виде []WordObject.
func (distMap *DistributionMap) AsWordObjectStructs() []WordObject {
	ws := []WordObject{}
	for word, count := range distMap.mappedStatValue {
		if word == "" {
			continue
		}
		ws = append(ws, WordObject{word, count})
	}
	return ws
}

func Top10(s string) []string {
	// Все слова без иных пробельных символов
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\f", " ")
	s = strings.ReplaceAll(s, "\v", " ")
	stringWords := strings.Split(s, " ")
	// fmt.Println("stringWords", stringWords)

	// Распределение в формате MAP key-value: слово -> его частота
	distMap := DistributionMap{make(map[string]uint)}
	distMap.Calculate(stringWords)
	// fmt.Println("distribution", distMap)
	wordsObjects := distMap.AsWordObjectStructs()
	// fmt.Println("wordsObjects", wordsObjects)

	// Сортировка в условии задачи: `ORDER BY count DESC, lexical ASC`
	sort.Slice(
		wordsObjects,
		func(i, j int) bool {
			return wordsObjects[i].count > wordsObjects[j].count || (wordsObjects[i].count == wordsObjects[j].count &&
				wordsObjects[i].word < wordsObjects[j].word)
		})
	// Синтаксический сахар
	orderedWordsObjects := wordsObjects
	// fmt.Println("wordsObjects", wordsObjects)

	wordsLen := len(orderedWordsObjects)
	iter := 0
	returnedWords := []string{}
	// Для отладки: returnedWordsObjects := []WordObject{}
	for iter < LIMIT && iter < wordsLen {
		returnedWords = append(returnedWords, wordsObjects[iter].word)
		// Для отладки: returnedWordsObjects = append(returnedWordsObjects, wordsObjects[iter])
		iter++
	}

	// fmt.Println("returnedWords", returnedWords)
	// fmt.Println("returnedWordsObjects", returnedWordsObjects)
	return returnedWords
}
