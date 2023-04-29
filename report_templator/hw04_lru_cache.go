package main

// go run templator.go hw04_lru_cache.go
func main() {

	substitutions := make(map[string]string)
	//
	list_go := Template{}
	list_go.loadFromFile("../hw04_lru_cache/list.go", false)
	substitutions["list.go"] = list_go.render()
	//
	cache_go := Template{}
	cache_go.loadFromFile("../hw04_lru_cache/cache.go", false)
	substitutions["cache.go"] = cache_go.render()
	//
	doc_txt := Template{}
	doc_txt.loadFromFile("../hw04_lru_cache/doc.txt", false)
	substitutions["doc.txt"] = doc_txt.render()
	//
	O__txt := Template{}
	O__txt.loadFromFile("../hw04_lru_cache/O?.txt", false)
	substitutions["O?.txt"] = O__txt.render()
	//
	list_testing_txt := Template{}
	list_testing_txt.loadFromFile("../hw04_lru_cache/list_testing.txt", false)
	substitutions["list_testing.txt"] = list_testing_txt.render()
	//
	cache_testing_txt := Template{}
	cache_testing_txt.loadFromFile("../hw04_lru_cache/cache_testing.txt", false)
	substitutions["cache_testing.txt"] = cache_testing_txt.render()
	//
	report := Template{}
	report.loadFromFile("../hw04_lru_cache/REPORT.template.md", false)
	report.substitutions = substitutions
	report.renderToFile("../hw04_lru_cache/REPORT.md")
}
