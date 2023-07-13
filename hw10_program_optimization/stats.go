package hw10programoptimization


import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	// "errors"
	"bufio"
)

type DomainStat map[string]int

func RowParser(
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	workersCount := 50000 // Enviroment - could be better

	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp := regexp.MustCompile(domainAtEmailRegexp)

	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	rowsChannel := make(chan []byte)

	domainStat := make(DomainStat)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go RowParser(rowsChannel, compiledRegexp, &wg, &mtx, &domainStat)
	}

	// b := make([]byte, 1)
	// chunk := make([]byte, 0)
	// for {
	// 	_, err := r.Read(b)
	// 	if err != nil {
	// 		if errors.Is(err, io.EOF) {
	// 			rowsChannel <- chunk
	// 			break
	// 		}
	// 		fmt.Println(err)
	// 		break
	// 	}
	// 	if b[0] == '\n' {
	// 		rowsChannel <- chunk
	// 		chunk = make([]byte, 0)
	// 		continue
	// 	}
	// 	chunk = append(chunk, b...)
	// }

	scanner := bufio.NewScanner(r)
	const maxCapacity int = 5000000 // It's over 64K !!! Byt why 5Mln
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// bb := bytes.NewReader(r)
	// scanners := bufio.NewReaderSize(bb, 10000000)
	for scanner.Scan() {
		rowsChannel <- scanner.Bytes()
	}

	close(rowsChannel)
	wg.Wait()

	// syncMap.Range(func(key, value interface{}) bool {
	// 	syncMapDomainStat[key.(string)] = value.(int)
	// 	return true
	// })

	return domainStat, nil
}
