package structs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Job struct {
	SessionName   string `toml:"session_name"`
	Name          string `toml:"name"`
	RequestBody   string `toml:"request_body"`
	URL           string `toml:"url"`
	RequestMethod string `toml:"request_method"`

	StatusCode   int    `toml:"status_code"`
	ResponseBody string `toml:"response_body"`
	ResponseType string `toml:"response_type"`
}

func (j *Job) Run() error {
	jar, err := GlobalSessionStorage.GetJar(j.SessionName)
	if err != nil {
		return err
	}

	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest(
		j.RequestMethod,
		j.URL,
		strings.NewReader(j.RequestBody),
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != j.StatusCode {
		return errors.New("status code is not match")
	}

	switch j.ResponseType {
	case "json":
		a := map[string]interface{}{}
		b := map[string]interface{}{}

		ar := strings.NewReader(j.RequestBody)
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
