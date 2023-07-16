package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

type DomainStat map[string]int

func rowParser(
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
	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	dataChannel := make(chan string)
	domainStat := make(DomainStat)
	workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go rowParser(&wg, &mtx, &domainStat, dataChannel)
	}

	scanner := bufio.NewScanner(r)
	maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 239)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		matches := compiledRegexp.FindAllSubmatch(scanner.Bytes(), -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
			dataChannel <- domainAtLowercase
		}
	}

	close(dataChannel)
	wg.Wait()

	return domainStat, nil
}
