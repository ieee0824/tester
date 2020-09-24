package structs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Job struct {
	SessionName   string `yaml:"session_name"`
	Name          string `yaml:"name"`
	RequestBody   string `yaml:"request_body"`
	URL           string `yaml:"url"`
	RequestMethod string `yaml:"request_method"`

	StatusCode   int     `yaml:"status_code"`
	ResponseBody *string `yaml:"response_body"`
	ResponseType string  `yaml:"response_type"`
}

func (j *Job) Run() error {
	jar, err := GlobalSessionStorage.GetJar(j.SessionName)
	if err != nil {
		return err
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
		a := map[string]interface{}{}
		b := map[string]interface{}{}

		ar := strings.NewReader(*j.ResponseBody)
		br := resp.Body

		if err := json.NewDecoder(ar).Decode(&a); err != nil {
			return err
		}
		if err := json.NewDecoder(br).Decode(&b); err != nil {
			return err
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
