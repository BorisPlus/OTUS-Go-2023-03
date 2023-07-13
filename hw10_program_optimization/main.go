package hw10programoptimization // change to main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// For Example2: internal mutex in DomainStat

type InternalMutexedDomainStat struct {
	mx sync.RWMutex
	m  DomainStat
}

func (c *InternalMutexedDomainStat) DomainStat() DomainStat {
	return c.m
}

func (c *InternalMutexedDomainStat) Load(key string) (int, bool) {
	c.mx.RLock()
	val, ok := c.m[key]
	c.mx.RUnlock()
	return val, ok
}

func (c *InternalMutexedDomainStat) Store(key string, value int) {
	c.mx.Lock()
	c.m[key] = value
	c.mx.Unlock()
}

func (c *InternalMutexedDomainStat) Increment(key string) {
	c.mx.Lock()
	c.m[key]++
	c.mx.Unlock()
}

func NewInternalMutexedDomainStat() *InternalMutexedDomainStat {
	return &InternalMutexedDomainStat{
		m: make(DomainStat),
	}
}

func MainRowParser(
	row <-chan []byte,
	compiledRegexp *regexp.Regexp,
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	domainStat *DomainStat, // Example1: mutexed DomainStat
	internalMutexedDomainStat *InternalMutexedDomainStat, // Example2: internal mutex in DomainStat
	syncMap *sync.Map, // Example3: sync.Map fo DomainStat
) {
	defer wg.Done()
	for rowData := range row {
		matches := compiledRegexp.FindAllSubmatch(rowData, -1)
		for matcheIndex := range matches {
			domainAtLowercase := strings.ToLower(string(matches[matcheIndex][1]))

			// Example1: mutexed DomainStat

			mtx.Lock()
			(*domainStat)[domainAtLowercase]++
			mtx.Unlock()

			// Example2: internal mutex in DomainStat

			internalMutexedDomainStat.Increment(domainAtLowercase)

			// Example3: sync.Map fo DomainStat

			v, ok := syncMap.LoadOrStore(domainAtLowercase, 1)
			if ok {
				syncMap.Store(domainAtLowercase, v.(int)+1)
			}
		}
	}
}

func main() {

	workersCount := 2 // If 1 - It's OK result
	
	testFile, _ := zip.OpenReader("testdata/users.dat.zip")
	defer testFile.Close()
	data, err := testFile.File[0].Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	domain := "biz"
	domainAtEmailRegexp := fmt.Sprintf(`@(\w+\.%s)`, domain)
	compiledRegexp := regexp.MustCompile(domainAtEmailRegexp)

	// Goroutines block

	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	rowsChannel := make(chan []byte)

	// Example1: mutexed DomainStat
	domainStat := make(DomainStat)
	// Example2: internal mutex in DomainStat
	internalMutexedDomainStat := NewInternalMutexedDomainStat()
	// Example3: sync.Map fo DomainStat
	var syncMap sync.Map

	// Parsing at goroutines

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go MainRowParser(rowsChannel, compiledRegexp, &wg, &mtx, &domainStat, internalMutexedDomainStat, &syncMap)
	}

	// Read file row-by-row to channel

	scanner := bufio.NewScanner(data)
	const maxCapacity int = 100000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		rowsChannel <- scanner.Bytes()
	}

	close(rowsChannel)
	wg.Wait()

	// Example1: mutexed DomainStat

	fileExample1, err := os.Create("testdata/DomainStat_Raw.out")
	if err != nil {
		fmt.Printf("Couldn't open file")
	}
	defer fileExample1.Close()

	for i := range domainStat {
		result := fmt.Sprintf("DOMAIN: %s COUNT: %d\n", i, domainStat[i])
		fileExample1.WriteString(result)
	}

	// Example2: internal mutex in DomainStat

	fileExample2, err := os.Create("testdata/DomainStat_InternalMutexed.out")
	if err != nil {
		fmt.Printf("Couldn't open file")
	}
	defer fileExample2.Close()

	domainStatAtMutex := internalMutexedDomainStat.DomainStat()
	for i := range domainStatAtMutex {
		result := fmt.Sprintf("DOMAIN: %s COUNT: %d\n", i, domainStatAtMutex[i])
		fileExample2.WriteString(result)
	}

	// Example3: sync.Map fo DomainStat

	fileExample3, err := os.Create("testdata/DomainStat_SyncMap.out")
	if err != nil {
		fmt.Printf("Couldn't open file")
	}
	defer fileExample3.Close()

	syncMapDomainStat := make(DomainStat)

	syncMap.Range(func(key, value interface{}) bool {
		syncMapDomainStat[key.(string)] = value.(int)
		return true
	})

	for i := range syncMapDomainStat {
		result := fmt.Sprintf("DOMAIN: %s COUNT: %d\n", i, syncMapDomainStat[i])
		fileExample3.WriteString(result)
	}
}
