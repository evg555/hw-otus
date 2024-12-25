package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user User
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, err
		}

		matched := strings.HasSuffix(user.Email, "."+domain)

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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
