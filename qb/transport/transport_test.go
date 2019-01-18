package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"mooncamp.com/dgraphtools/gql"
	"mooncamp.com/dgraphtools/qb/endpoint"
)

func Test_parse_template(t *testing.T) {
	handler := NewHTTPHandler(endpoint.NewEndpointSet(), "/api/v1")
	server := httptest.NewServer(handler)

	u, _ := url.Parse(server.URL)
	u.Path = "/api/v1/parse"

	body := bytes.NewBuffer(nil)

	parseRequest := struct {
		Query     string            `json:"query"`
		Variables map[string]string `json:"variables"`
	}{
		Query: `
{
  bladerunner(func: uid(0x107b2c)) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
`,
		Variables: map[string]string{},
	}

	if err := json.NewEncoder(body).Encode(parseRequest); err != nil {
		t.Fatalf("encode: %v", err)
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		t.Fatalf("new req: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do req: %v", err)
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("parse: %s", string(buf))
	}

	var queries []gql.GraphQuery
	_ = json.NewDecoder(bytes.NewReader(buf)).Decode(&queries)

	templateRequest := struct {
		Queries   []gql.GraphQuery  `json:"queries"`
		Alias     string            `json:"alias"`
		Variables map[string]string `json:"variables"`
	}{
		Queries:   queries,
		Alias:     "",
		Variables: map[string]string{},
	}

	body.Truncate(0)
	if err := json.NewEncoder(body).Encode(templateRequest); err != nil {
		t.Fatalf("encode: %v", err)
	}

	u.Path = "/api/v1/template"
	req, err = http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		t.Fatalf("new req: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("do req: %v, %s", err, string(buf))
	}

	_, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("template: %s", string(buf))
	}
}
