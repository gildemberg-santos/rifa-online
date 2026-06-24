package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

func makeRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
