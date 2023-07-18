package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

func GetDomainStatGoroutinedFastjson(r io.Reader, domain string) (DomainStat, error) {
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	dataChannel := make(chan string)
	domainStat := make(DomainStat)
	workersCount := loadEnviromentOrDefault("WORKERS_COUNT", 100)
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
