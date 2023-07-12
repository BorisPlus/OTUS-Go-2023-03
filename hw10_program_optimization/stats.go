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

func worker(waitGroup *sync.WaitGroup, mutex *sync.Mutex, bytesSlices <-chan []byte, re regexp.Regexp, domainStat *DomainStat) {
	defer waitGroup.Done()
	for bytesSlice := range bytesSlices {
		submatches := re.FindAllSubmatch(bytesSlice, -1)
		for matcheIndex := range submatches {
			domainAtLowercase := strings.ToLower(string(submatches[matcheIndex][1]))
			mutex.Lock()
			(*domainStat)[domainAtLowercase]++
			mutex.Unlock()
		}
	}
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {

	domainAtEmailRegexp := fmt.Sprintf(`@(\w*\.%s)`, domain)
	fmt.Println(domainAtEmailRegexp)
	compiledRegexp, err := regexp.Compile(domainAtEmailRegexp)

	if err != nil {
		return nil, err
	}

	domainStat := make(DomainStat)
	bytesSlicesChannel := make(chan []byte)
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(&wg, &mtx, bytesSlicesChannel, *compiledRegexp, &domainStat)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		row := scanner.Bytes()
		// fmt.Println(string(row))
		bytesSlicesChannel <- row
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	close(bytesSlicesChannel)
	wg.Wait()

	return domainStat, nil
}
