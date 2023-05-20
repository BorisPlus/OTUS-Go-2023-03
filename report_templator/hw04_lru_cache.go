package reporttemplator

// go run templator.go hw04_lru_cache.go

func TestHw04() {
	substitutions := make(map[string]string)

	listGo := Template{}
	listGo.loadFromFile("../hw04_lru_cache/list.go")
	substitutions["list.go"] = listGo.render(true)

	cacheGo := Template{}
	cacheGo.loadFromFile("../hw04_lru_cache/cache.go")
	substitutions["cache.go"] = cacheGo.render(true)

	docTxt := Template{}
	docTxt.loadFromFile("../hw04_lru_cache/doc.txt")
	substitutions["doc.txt"] = docTxt.render(true)

	listTestTxt := Template{}
	listTestTxt.loadFromFile("../hw04_lru_cache/list_test.txt")
	substitutions["list_test.txt"] = listTestTxt.render(true)

	cacheTestTxt := Template{}
	cacheTestTxt.loadFromFile("../hw04_lru_cache/cache_test.txt")
	substitutions["cache_test.txt"] = cacheTestTxt.render(true)

	cacheTestDataGo := Template{}
	cacheTestDataGo.loadFromFile("../hw04_lru_cache/cache_test_data.go")
	substitutions["cache_test_data.go"] = cacheTestDataGo.render(true)

	cacheBenchmarkCliTestGo := Template{}
	cacheBenchmarkCliTestGo.loadFromFile("../hw04_lru_cache/cache_benchmark_cli_test.go")
	substitutions["cache_benchmark_cli_test.go"] = cacheBenchmarkCliTestGo.render(true)

	cacheBenchmarkCliTestTxt := Template{}
	cacheBenchmarkCliTestTxt.loadFromFile("../hw04_lru_cache/cache_benchmark_cli_test.txt")
	substitutions["cache_benchmark_cli_test.txt"] = cacheBenchmarkCliTestTxt.render(false)

	cacheBenchmarkNocliTestGo := Template{}
	cacheBenchmarkNocliTestGo.loadFromFile("../hw04_lru_cache/cache_benchmark_nocli_test.go")
	substitutions["cache_benchmark_nocli_test.go"] = cacheBenchmarkNocliTestGo.render(true)

	cacheBenchmarkNocliTestTxt := Template{}
	cacheBenchmarkNocliTestTxt.loadFromFile("../hw04_lru_cache/cache_benchmark_nocli_test.txt")
	substitutions["cache_benchmark_nocli_test.txt"] = cacheBenchmarkNocliTestTxt.render(true)

	report := Template{}
	report.loadFromFile("../hw04_lru_cache/REPORT.template.md")
	report.substitutions = substitutions
	report.renderToFile("../hw04_lru_cache/REPORT.md")
}
