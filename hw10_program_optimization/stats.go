package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
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
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	pattern, err := regexp.Compile("\\." + domain + "$")
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		var user User
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, err
		}

		matched := pattern.MatchString(user.Email)

		if matched {
			emailParts := strings.SplitN(user.Email, "@", 2)

			if len(emailParts) != 2 {
				return nil, fmt.Errorf("invalid email: %s", user.Email)
			}

			num := result[strings.ToLower(emailParts[1])]
			num++
			result[strings.ToLower(emailParts[1])] = num
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
