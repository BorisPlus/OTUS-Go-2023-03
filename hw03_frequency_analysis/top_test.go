package hw03frequencyanalysis

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestStringSpliter(t *testing.T) {
	testString := "qwe asd zxc qwe"
	result := StringSpliter(testString)
	expected := []string{"qwe", "asd", "zxc", "qwe"}
	if len(expected) != len(result) {
		t.Errorf("StringSpliter(%q) = %q, but expected %q.", testString, result, expected)
	}
	for idx := range expected {
		if expected[idx] != result[idx] {
			t.Errorf("StringSpliter(%q) = %q, but expected %q.", testString, result, expected)
		}
	}

	fmt.Printf("OK. StringSpliter(%q) return expected %q.\n", testString, expected)
}

func TestInitDistributionOfMappedWord(t *testing.T) {

	testCases := []struct {
		stringWords         []string
		expectedMappedWords map[string]uint
	}{
		{
			stringWords:         []string{},
			expectedMappedWords: map[string]uint{},
		},
		{
			stringWords:         []string{"x"},
			expectedMappedWords: map[string]uint{"x": 1},
		},
		{
			stringWords:         []string{"x", "x"},
			expectedMappedWords: map[string]uint{"x": 2},
		},
		{
			stringWords:         []string{"y", "x", "x"},
			expectedMappedWords: map[string]uint{"x": 2, "y": 1},
		},
	}

	for _, testCase := range testCases {
		resultMappedWords := InitDistributionOfMappedWord(testCase.stringWords)
		if len(resultMappedWords.Distribution) != len(testCase.expectedMappedWords) {
			t.Errorf("InitDistributionOfMappedWord(%v) = %v, but expected %v.",
				testCase.stringWords, resultMappedWords, testCase.expectedMappedWords)
		}
		for v := range testCase.expectedMappedWords {
			if testCase.expectedMappedWords[v] != resultMappedWords.Distribution[v] {
				t.Errorf("InitDistributionOfMappedWord(%v) = %v, but expected %v.",
					testCase.stringWords, resultMappedWords, testCase.expectedMappedWords)
			}
		}
		fmt.Printf("OK. InitDistributionOfMappedWord(%v) return expected %v.\n", testCase.stringWords, testCase.expectedMappedWords)
	}
}

func TestSortedStructWords(t *testing.T) {
	testCases := []struct {
		mappedWords map[string]uint
		expected    []StructWord
	}{
		{
			mappedWords: map[string]uint{},
			expected:    []StructWord{},
		},
		{
			mappedWords: map[string]uint{"x": 1},
			expected:    []StructWord{{word: "x", count: 1}},
		},
		{
			mappedWords: map[string]uint{"x": 2},
			expected:    []StructWord{{word: "x", count: 2}},
		},
		{
			mappedWords: map[string]uint{"x": 2, "y": 1, "z": 3, "a": 3},
			expected:    []StructWord{{"a", 3}, {"z", 3}, {"x", 2}, {"y", 1}},
		},
	}

	for _, testCase := range testCases {
		dMappedWords := DistributionOfMappedWord{}
		dMappedWords.Distribution = testCase.mappedWords
		structWords := dMappedWords.GetAsSortedStructWords()
		if len(testCase.expected) != len(structWords) {
			t.Errorf("(%v).GetAsSortedStructWords = %v, but expected %v.",
				testCase.mappedWords, structWords, testCase.expected)
		}
		for idx := range structWords {
			if testCase.expected[idx] != structWords[idx] {
				t.Errorf("(%v).GetAsSortedStructWords = %v, but expected %v.",
					testCase.mappedWords, structWords, testCase.expected)
			}
		}
		fmt.Printf("OK. (%v).GetAsSortedStructWords() return expected %v.\n", testCase.mappedWords, testCase.expected)
	}
}

func TestGetTopStructWords(t *testing.T) {
	testCases := []struct {
		structWords []StructWord
		limit       uint
		expected    []StructWord
	}{
		{
			structWords: []StructWord{},
			limit:       10,
			expected:    []StructWord{},
		},
		{
			structWords: []StructWord{{"d", 3}, {"e", 3}, {"v", 2}, {"h", 1}},
			limit:       0,
			expected:    []StructWord{},
		},
		{
			structWords: []StructWord{{"a", 3}, {"z", 3}, {"x", 2}, {"y", 1}},
			limit:       5,
			expected:    []StructWord{{"a", 3}, {"z", 3}, {"x", 2}, {"y", 1}},
		},
		{
			structWords: []StructWord{{"a", 3}, {"z", 3}, {"x", 2}, {"y", 1}},
			limit:       2,
			expected:    []StructWord{{"a", 3}, {"z", 3}},
		},
		{
			structWords: []StructWord{{"a", 3}, {"z", 3}, {"x", 2}, {"y", 1}},
			limit:       1,
			expected:    []StructWord{{"a", 3}},
		},
	}

	for _, testCase := range testCases {

		top := GetTopStructWords(testCase.structWords, testCase.limit)
		if len(top) != len(testCase.expected) {
			t.Errorf("GetTopStructWords(%v, %v) = %v, but expected %v.",
				testCase.structWords, testCase.limit, top, testCase.expected)
		}

		for idx := range testCase.expected {
			if testCase.expected[idx] != top[idx] {
				t.Errorf("GetTopStructWords(%v, %v) = %v, but expected %v.",
					testCase.structWords, testCase.limit, top, testCase.expected)
			}
		}
		fmt.Printf("OK. GetTopStructWords(%v, %v) return expected %v.\n", testCase.structWords, testCase.limit, testCase.expected)
	}
}

func TestWordStructToWordStrings(t *testing.T) {
	testCases := []struct {
		structWords []StructWord
		expected    []string
	}{
		{
			structWords: []StructWord{},
			expected:    []string{},
		},
		{
			structWords: []StructWord{{"ad", 3}, {"ae", 3}, {"vj", 2}, {"h", 1}},
			expected:    []string{"ad", "ae", "vj", "h"},
		},
		{
			structWords: []StructWord{{"ccc", 10}, {"aaa", 3}, {"b", 1}},
			expected:    []string{"ccc", "aaa", "b"},
		},
	}

	for _, testCase := range testCases {
		result := WordStructToWordStrings(testCase.structWords)

		if len(result) != len(testCase.expected) {
			t.Errorf("WordStructToWordStrings(%v) = %v, but expected %v.",
				testCase.structWords, result, testCase.expected)
		}

		for idx := range testCase.expected {
			if testCase.expected[idx] != result[idx] {
				t.Errorf("WordStructToWordStrings(%v) = %v, but expected %v.",
					testCase.structWords, result, testCase.expected)
			}
		}
		fmt.Printf("OK. WordStructToWordStrings(%v) return expected %v.\n", testCase.structWords, testCase.expected)

	}
}

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test at HW example", func(t *testing.T) {
		expected := []string{
			"and",     // 2
			"one",     // 2
			"cat",     // 1
			"cats",    // 1
			"dog,",    // 1
			"dog,two", // 1
			"man",     // 1
		}
		require.Equal(t, expected, Top10("cat and dog, one dog,two cats and one man"))
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})

	t.Run("negative test", func(t *testing.T) {
		if !taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.NotEqual(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.NotEqual(t, expected, Top10(text))
		}
	})
}
