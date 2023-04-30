package main

// go run templator.go hw04_lru_cache.go
func main() {

	substitutions := make(map[string]string)
	//
	list_go := Template{}
	list_go.loadFromFile("../hw04_lru_cache/list.go")
	substitutions["list.go"] = list_go.render(true)
	//
	cache_go := Template{}
	cache_go.loadFromFile("../hw04_lru_cache/cache.go")
	substitutions["cache.go"] = cache_go.render(true)
	//
	doc_txt := Template{}
	doc_txt.loadFromFile("../hw04_lru_cache/doc.txt")
	substitutions["doc.txt"] = doc_txt.render(true)
	//
	list_test_txt := Template{}
	list_test_txt.loadFromFile("../hw04_lru_cache/list_test.txt")
	substitutions["list_test.txt"] = list_test_txt.render(true)
	//
	cache_test_txt := Template{}
	cache_test_txt.loadFromFile("../hw04_lru_cache/cache_test.txt")
	substitutions["cache_test.txt"] = cache_test_txt.render(true)
	// 
	cache_test_data_go := Template{}
	cache_test_data_go.loadFromFile("../hw04_lru_cache/cache_test_data.go")
	substitutions["cache_test_data.go"] = cache_test_data_go.render(true)
	//
	cache_benchmark_cli_test_go := Template{}
	cache_benchmark_cli_test_go.loadFromFile("../hw04_lru_cache/cache_benchmark_cli_test.go")
	substitutions["cache_benchmark_cli_test.go"] = cache_benchmark_cli_test_go.render(true)
	//
	cache_benchmark_cli_test_txt := Template{}
	cache_benchmark_cli_test_txt.loadFromFile("../hw04_lru_cache/cache_benchmark_cli_test.txt")
	substitutions["cache_benchmark_cli_test.txt"] = cache_benchmark_cli_test_txt.render(false)
	//
	cache_benchmark_nocli_test_go := Template{}
	cache_benchmark_nocli_test_go.loadFromFile("../hw04_lru_cache/cache_benchmark_nocli_test.go")
	substitutions["cache_benchmark_nocli_test.go"] = cache_benchmark_nocli_test_go.render(true)
	//
	cache_benchmark_nocli_test_txt := Template{}
	cache_benchmark_nocli_test_txt.loadFromFile("../hw04_lru_cache/cache_benchmark_nocli_test.txt")
	substitutions["cache_benchmark_nocli_test.txt"] = cache_benchmark_nocli_test_txt.render(true)
	//
	report := Template{}
	report.loadFromFile("../hw04_lru_cache/REPORT.template.md")
	report.substitutions = substitutions
	report.renderToFile("../hw04_lru_cache/REPORT.md")
}
