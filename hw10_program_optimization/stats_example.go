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

// TODO: with "reflect" - func LoadOrDefault[T any](name string, asDefault T) T {}.
func loadEnviromentOrDefault(name string, asDefault int) int {
	value, exists := os.LookupEnv(name)
	if exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return asDefault
}

func rowParserExample(
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	rows <-chan []byte,
	compiledRegexp regexp.Regexp,
	domainStat DomainStat,
) {
	defer wg.Done()
	for row := range rows {
		matches := compiledRegexp.FindAllSubmatch(row, -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
			mtx.Lock()
			domainStat[domainAtLowercase]++
			mtx.Unlock()
		}
	}
}

func GetDomainStatExample(r io.Reader, domain string) (DomainStat, error) {
	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	dataChannel := make(chan []byte)
	domainStat := make(DomainStat)
	workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 1)
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go rowParserExample(&wg, &mtx, dataChannel, *compiledRegexp, domainStat)
	}
	scanner := bufio.NewScanner(r)
	maxCapacity := loadEnviromentOrDefault("MAX_CAPACITY", 2_000_000) // Magick big value
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		dataChannel <- scanner.Bytes()
	}
	close(dataChannel)
	wg.Wait()
	return domainStat, nil
}
