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

func RowParser(
	row <-chan []byte,
	compiledRegexp *regexp.Regexp,
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	syncMap *sync.Map,
) {
	defer wg.Done()
	for rowData := range row {
		matches := compiledRegexp.FindAllSubmatch(rowData, -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))
			// mtx.Lock()
			// mtx.Unlock()
			_ = mtx
			v, exsist := syncMap.LoadOrStore(domainAtLowercase, 1)
			if exsist {
				syncMap.Store(domainAtLowercase, v.(int)+1)
			}
		}
	}
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	workersCount := 100 // Enviroment - could be better

	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp := regexp.MustCompile(domainAtEmailRegexp)
	
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	rowsChannel := make(chan []byte)

	var syncMap sync.Map

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go RowParser(rowsChannel, compiledRegexp, &wg, &mtx, &syncMap)
	}

	scanner := bufio.NewScanner(r)
	const maxCapacity int = 10000000 // It's over 64K !!!
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		rowsChannel <- scanner.Bytes()
	}

	close(rowsChannel)
	wg.Wait()

	syncMapDomainStat := make(DomainStat)

	syncMap.Range(func(key, value interface{}) bool {
		syncMapDomainStat[key.(string)] = value.(int)
		return true
	})

	return syncMapDomainStat, nil
}
