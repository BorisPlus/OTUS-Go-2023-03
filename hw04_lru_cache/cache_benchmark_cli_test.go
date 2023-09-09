package hw04_lru_cache

import (
	"fmt"
	"testing"
	"time"
)

// go test -bench=. -benchmem list.go cache.go cache_test_data.go cache_benchmark_cli_test.go \
// > cache_benchmark_cli_test.txt

// Операции, необходимые для подсчета статистики времени добавления
// x2 - возведение в степень 2.
func x2(x float64) float64 {
	return x * x
}

// sum - сумма.
func sum(arr []float64) float64 {
	var sum float64
	for _, value := range arr {
		sum += value
	}
	return sum
}

// average - среднее значение (времени проведения операции Set/Get).
func average(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// dispersion - дисперсия (времени проведения операции Set/Get).
func dispersion(data []float64) float64 {
	dataAverage := average(data)
	x2up := make([]float64, 0, len(data))
	for _, value := range data {
		x2up = append(x2up, x2(value-dataAverage))
	}
	return sum(x2up) / float64(len(data))
}

func BenchmarkSet(b *testing.B) {
	for _, testDataset := range TestDatasets {
		b.Run(fmt.Sprintf("%d", len(testDataset.data)), func(b *testing.B) {
			// чтоб не было операций вымещения capasity=valuesCount
			cache := NewCache(len(testDataset.data))
			durationsSet, durationsGet := []float64{}, []float64{}
			b.ResetTimer()
			// Собираем данные для статистики в отношении метода Set
			for _, keyValue := range testDataset.data {
				start := time.Now()
				b.StartTimer()

				cache.Set(keyValue.key, keyValue.value)

				b.StopTimer()
				duration := time.Since(start)

				durationsSet = append(durationsSet, float64(duration.Microseconds()))
			}

			// Собираем данные для статистики в отношении метода Get
			for _, keyValue := range testDataset.data {
				start := time.Now()
				b.StartTimer()

				cache.Get(keyValue.key)

				b.StopTimer()
				duration := time.Since(start)

				// durationsGet = append(durationsGet, float64(duration.Microseconds()))
				durationsGet = append(durationsGet, float64(duration.Nanoseconds()))
			}

			// Среднее времени добавления в LRU
			b.ReportMetric(average(durationsSet), "med_t/set")
			// Дисперсия времени добавления в LRU
			b.ReportMetric(dispersion(durationsSet), "disp_t/set")
			// Среднее времени взятия из LRU
			b.ReportMetric(average(durationsGet), "med_t/get")
			// Дисперсия времени взятия из LRU
			b.ReportMetric(dispersion(durationsGet), "disp_t/get")
		})
	}
}
