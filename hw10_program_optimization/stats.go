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


type DomainStat map[string]int

func RowParser(
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	domainStat *DomainStat,
	domains <-chan string,

) {
	defer wg.Done()
	for domain := range domains {
		mtx.Lock()
		(*domainStat)[domain]++
		mtx.Unlock()
	}
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	workersCount := 1 
	count, exists := os.LookupEnv("WORKERS_COUNT")
	if exists {
		intVar, err := strconv.Atoi(count)
		if err == nil {
			workersCount = intVar
		}
	}
	// fmt.Println("workersCount =", workersCount)

	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	domainsChannel := make(chan string)

	domainStat := make(DomainStat)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go RowParser(&wg, &mtx, &domainStat, domainsChannel)
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
	// fmt.Println("maxCapacity =", maxCapacity)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		matches := compiledRegexp.FindAllSubmatch(scanner.Bytes(), -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
			domainsChannel <- domainAtLowercase
		}
	}

	close(domainsChannel)
	wg.Wait()

	return domainStat, nil
}
