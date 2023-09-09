package hw04_lru_cache

import (
	"fmt"
	"testing"
)

// go test -v list.go cache.go cache_test_data.go cache_benchmark_nocli_test.go > cache_benchmark_nocli_test.txt

func TestBenchmark(t *testing.T) {
	_ = t
	for _, testDataset := range TestDatasets {
		dataset := testDataset.data
		cache := NewCache(len(dataset))

		resForSet := testing.Benchmark(
			func(b *testing.B) {
				b.Helper()
				// Вынести в отдельную функцию это никак из-за необходимости сигнатуры f(b *B).
				// Вынужденное анонимное замыкание на dataset.
				b.ResetTimer()
				for _, keyValue := range dataset {
					b.StartTimer()
					cache.Set(keyValue.key, keyValue.value)
					b.StopTimer()
				}
			},
		)
		fmt.Printf("--------------------------------------------------------\n")
		fmt.Printf("Operation - Set - with dataset %d values\n", len(dataset))
		fmt.Printf("--------------------------------------------------------\n")
		fmt.Printf("Number of run: %d\n", resForSet.N)
		fmt.Printf("Memory allocations: %d\n", resForSet.MemAllocs)
		fmt.Printf("Memory allocations (AVERAGE): %f\n", float64(resForSet.MemAllocs)/float64(len(dataset)))
		fmt.Printf("Number of bytes allocated: %d\n", resForSet.Bytes)
		fmt.Printf("Number of bytes allocated (AVERAGE): %f\n", float64(resForSet.Bytes)/float64(len(dataset)))
		fmt.Printf("Time taken: %s\n", resForSet.T)
		fmt.Printf("Time taken (AVERAGE, nanosecs.): %f  \n", float64(resForSet.T.Nanoseconds())/float64(len(dataset)))
		fmt.Printf("\n\n")
		res := testing.Benchmark(
			func(b *testing.B) {
				b.Helper()
				// Вынести в отдельную функцию это никак из-за необходимости сигнатуры f(b *B).
				// Вынужденное анонимное замыкание на dataset.
				b.ResetTimer()
				for _, keyValue := range dataset {
					b.StartTimer()
					cache.Get(keyValue.key)
					b.StopTimer()
				}
			},
		)
		fmt.Printf("--------------------------------------------------------\n")
		fmt.Printf("Operation - Get - with dataset %d values\n", len(dataset))
		fmt.Printf("--------------------------------------------------------\n")
		fmt.Printf("Number of run: %d\n", res.N)
		fmt.Printf("Memory allocations: %d\n", res.MemAllocs)
		fmt.Printf("Memory allocations (AVERAGE): %f\n", float64(res.MemAllocs)/float64(len(dataset)))
		fmt.Printf("Number of bytes allocated: %d\n", res.Bytes)
		fmt.Printf("Number of bytes allocated (AVERAGE): %f\n", float64(res.Bytes)/float64(len(dataset)))
		fmt.Printf("Time taken: %s\n", res.T)
		fmt.Printf("Time taken (AVERAGE, nanosecs.): %f  \n", float64(res.T.Nanoseconds())/float64(len(dataset)))
		fmt.Printf("\n\n")
	}
}
