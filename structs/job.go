package structs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/itchyny/gojq"
)

type JobOption struct {
	Debug bool
}

type Job struct {
	SessionName   string `yaml:"session_name"`
	Name          string `yaml:"name"`
	RequestBody   string `yaml:"request_body"`
	URL           string `yaml:"url"`
	RequestMethod string `yaml:"request_method"`

	StatusCode         int      `yaml:"status_code"`
	ResponseBody       *string  `yaml:"response_body"`
	IgnoreResponseKeys []string `yaml:"ignore_response_keys"`

	ResponseType string `yaml:"response_type"`
}

func (j *Job) Run(opts ...JobOption) error {
	isDebugMode := false
	if len(opts) != 0 {
		isDebugMode = opts[0].Debug
	}

	jar, err := GlobalSessionStorage.GetJar(j.SessionName)
	if err != nil {
		return err
	}

	if isDebugMode {
		log.Printf("jar: %v\n", jar)
		u, err := url.Parse(j.URL)
		if err != nil {
			panic(err)
		}
		cookies := jar.Cookies(u)

		log.Println("setted cookies")
		for _, c := range cookies {
			log.Println(c)
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Jar:       jar,
		Transport: tr,
	}

	var req *http.Request
	switch j.RequestMethod {
	case "POST":
		r, err := http.NewRequest(
			j.RequestMethod,
			j.URL,
			strings.NewReader(j.RequestBody),
		)
		if err != nil {
			return err
		}
		req = r
	default:
		r, err := http.NewRequest(
			j.RequestMethod,
			j.URL,
			nil,
		)
		if err != nil {
			return err
		}
		req = r
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != j.StatusCode {
		return fmt.Errorf("status code is not match %d, %d", j.StatusCode, resp.StatusCode)
	}

	if j.ResponseBody == nil {
		return nil
	}
	switch j.ResponseType {
	case "json":
		var a interface{}
		var b interface{}

		ar := strings.NewReader(*j.ResponseBody)
		br := resp.Body

		if err := json.NewDecoder(ar).Decode(&a); err != nil {
			return err
		}
		if err := json.NewDecoder(br).Decode(&b); err != nil {
			return err
		}

		for _, ignoreKey := range j.IgnoreResponseKeys {
			query, err := gojq.Parse(fmt.Sprintf("del(%s)", ignoreKey))
			if err != nil {
				return err
			}
			iter := query.Run(b)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return err
				}
				b = v
			}
		}

		if !reflect.DeepEqual(a, b) {
			return fmt.Errorf("response body not mathc: %v, %v", a, b)
		}
	default:
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)

		if j.RequestBody != buf.String() {
			return fmt.Errorf("response body not mathc: %s, %s", j.RequestBody, buf.String())
		}
	}

	return nil
}
