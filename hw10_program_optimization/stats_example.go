package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)


func RowParserExample(
	row <-chan []byte,
	compiledRegexp *regexp.Regexp,
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	domainStat *DomainStat,
) {
	defer wg.Done()
	for rowData := range row {
		matches := compiledRegexp.FindAllSubmatch(rowData, -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
			mtx.Lock()
			(*domainStat)[domainAtLowercase]++
			mtx.Unlock()
		}
	}
}

func GetDomainStatExample(r io.Reader, domain string) (DomainStat, error) {
	workersCount := 1 
	count, exists := os.LookupEnv("WORKERS_COUNT")
	if exists {
		intVar, err := strconv.Atoi(count)
		if err == nil {
			workersCount = intVar
		}
	}
	fmt.Println("workersCount =", workersCount)

	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	rowsChannel := make(chan []byte)

	domainStat := make(DomainStat)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go RowParserExample(rowsChannel, compiledRegexp, &wg, &mtx, &domainStat)
	}

	scanner := bufio.NewScanner(r)
	maxCapacity := 2_000_000 // It's over 64K !!!
	c, exists := os.LookupEnv("MAX_CAPACITY")
	if exists {
		intVar, err := strconv.Atoi(c)
		if err == nil {
			maxCapacity = intVar
		}
	}
	fmt.Println("maxCapacity =", maxCapacity)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		rowsChannel <- scanner.Bytes()
	}

	close(rowsChannel)
	wg.Wait()

	return domainStat, nil
}
