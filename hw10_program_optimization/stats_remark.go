package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

func rowParserRemark(
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	rows <-chan string,
	compiledRegexp regexp.Regexp,
	domainStat DomainStat,
) {
	defer wg.Done()
	for row := range rows {
		matches := compiledRegexp.FindAllStringSubmatch(row, -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(matches[matcheIndex][1])
			mtx.Lock()
			domainStat[domainAtLowercase]++
			mtx.Unlock()
		}
	}
}

func GetDomainStatRemark(r io.Reader, domain string) (DomainStat, error) {
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
		go rowParserRemark(&wg, &mtx, dataChannel, *compiledRegexp, domainStat)
	}
	scanner := bufio.NewScanner(r)
	maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 239)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		dataChannel <- scanner.Text()
	}
	close(dataChannel)
	wg.Wait()
	return domainStat, nil
}
