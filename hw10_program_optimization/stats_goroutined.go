package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

func GetDomainStatGoroutined(r io.Reader, domain string) (DomainStat, error) {
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
	var user User
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			continue
		}
		if strings.HasSuffix(user.Email, fmt.Sprintf(".%s", domain)) {
			dataChannel <- strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		}
	}
	close(dataChannel)
	wg.Wait()
	return domainStat, nil
}
