package hw10programoptimization

import (
	"bufio"
	"github.com/mailru/easyjson"
	"io"
	"strings"
)

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
	var (
		domainSuffix = "." + domain
		stat         = make(DomainStat)

		fullDomain string
		user       User
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if !strings.HasSuffix(user.Email, domainSuffix) {
			continue
		}

		fullDomain = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		stat[fullDomain]++
	}

	return stat, nil
}
