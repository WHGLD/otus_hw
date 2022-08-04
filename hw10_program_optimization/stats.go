//go:generate easyjson -no_std_marshalers stats.go
package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	index := 0
	for scanner.Scan() {
		var user User
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}
		result[index] = user
		index++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	regEx, _ := regexp.Compile("\\." + domain)

	for _, user := range u {
		if matched := regEx.MatchString(user.Email); matched {
			domainEmail := strings.SplitN(user.Email, "@", 2)[1]
			num := result[strings.ToLower(domainEmail)]
			num++
			result[strings.ToLower(domainEmail)] = num
		}
	}
	return result, nil
}
