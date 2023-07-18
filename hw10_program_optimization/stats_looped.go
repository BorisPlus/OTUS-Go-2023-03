package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"bufio"
	"strings"
)


func GetDomainStatLooped(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			continue
		}
		if strings.HasSuffix(user.Email, fmt.Sprintf(".%s", domain)) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
