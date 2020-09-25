package structs

import (
	"net/http/cookiejar"
)

var GlobalSessionStorage = SessionStorage{}

type SessionStorage map[string]*cookiejar.Jar

func (s SessionStorage) GetJar(name string) (*cookiejar.Jar, error) {
	jar, ok := s[name]
	if !ok {
		j, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		s[name] = j

		jar = j
	}

	return jar, nil
}
