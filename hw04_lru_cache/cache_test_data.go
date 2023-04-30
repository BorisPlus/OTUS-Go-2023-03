package hw04lrucache

import (
	"fmt"
	"math/rand"
)

func generateData(count int) []KeyValue {
	keyValues := []KeyValue{}

	for key, value := range rand.Perm(count) {
		keyValues = append(keyValues, KeyValue{key: Key(fmt.Sprint(key)), value: value})
	}
	return keyValues
}

var TestDatasets = []struct {
	data []KeyValue
}{
	{data: generateData(10)},
	{data: generateData(100)},
	{data: generateData(10000)},
}
