package structs

import (
	"net/http/cookiejar"
)

var GlobalSessionStorage = SessionStorage{
	cookies:         map[string]*cookiejar.Jar{},
	responseStorage: map[string]map[string]interface{}{},
}

type SessionStorage struct {
	cookies         map[string]*cookiejar.Jar
	responseStorage map[string]map[string]interface{}
}

func (s SessionStorage) StoreResponse(sessionName, key string, v interface{}) {
	_, ok := s.responseStorage[sessionName]
	if !ok {
		s.responseStorage[sessionName] = map[string]interface{}{}
	}
	s.responseStorage[sessionName][key] = v
}

func (s SessionStorage) LoadResponse(sessionName, key string) (interface{}, bool) {
	sessionStore, ok := s.responseStorage[sessionName]
	if !ok {
		return nil, false
	}
	ret, ok := sessionStore[key]
	return ret, ok
}

func (s SessionStorage) Keys(sessionName string) []string {
	keys := []string{}
	_, ok := s.responseStorage[sessionName]
	if !ok {
		return keys
	}
	for k := range s.responseStorage[sessionName] {
		keys = append(keys, k)
	}
	return keys
}

func (s SessionStorage) GetJar(name string) (*cookiejar.Jar, error) {
	jar, ok := s.cookies[name]
	if !ok {
		j, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		s.cookies[name] = j

		return j, nil
	}

	return jar, nil
}
