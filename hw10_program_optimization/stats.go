package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

// TODO: with "reflect" - func LoadOrDefault[T any](name string, asDefault T) T {}.
func loadEnviromentOrDefault(name string, asDefault int) int { //nolint:unparam
	value, exists := os.LookupEnv(name)
	if exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return asDefault
}

func domainStatCalc(
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
	domains <-chan string,
	domainStat DomainStat,
) {
	defer wg.Done()
	for domain := range domains {
		mtx.Lock()
		domainStat[domain]++
		mtx.Unlock()
	}
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	dataChannel := make(chan string)
	domainStat := make(DomainStat)
	workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 1000)
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go domainStatCalc(&wg, &mtx, dataChannel, domainStat)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		if strings.HasSuffix(email, fmt.Sprintf(".%s", domain)) {
			dataChannel <- strings.ToLower(strings.SplitN(email, "@", 2)[1])
		}
	}
	close(dataChannel)
	wg.Wait()
	return domainStat, nil
}
