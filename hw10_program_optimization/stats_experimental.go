package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

func GetDomainStatExperimental(r io.Reader, domain string) (DomainStat, error) {
	domainAtEmailRegexp := fmt.Sprintf(`@\w+\.%s`, domain)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	dataChannel := make(chan string)
	domainStat := make(DomainStat)
	workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 1)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go rowParser(&wg, &mtx, &domainStat, dataChannel)
	}

	scanner := bufio.NewScanner(r)
	maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 64_000)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		matches := compiledRegexp.FindAll(scanner.Bytes(), -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(strings.SplitN(string(matches[matcheIndex]), "@", 2)[1])
			dataChannel <- domainAtLowercase
		}
	}

	close(dataChannel)
	wg.Wait()

	return domainStat, nil
}
