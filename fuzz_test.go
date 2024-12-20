// Copyright 2024 Kristopher Rahim Afful-Brown. All Rights Reserved.
//
// Distributed under MIT license.
// See file LICENSE for detail or copy at https://opensource.org/licenses/MIT

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Fuzz_handleCalcHigh(f *testing.F) {
	s := newServer(f)

	type testcase struct {
		s string
	}
	tt := []testcase{
		{`{"values":[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`},
		{`{"values":[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]}`},
		{`{"values":[-50, -9, -8, -7, -6, -5, -4, -3, -2, -1]}`},
		{`{"values":[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]}`},
		{`{"values":[10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200]}`},
	}
	for _, tc := range tt {
		f.Add([]byte(tc.s))
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		if !json.Valid(a) {
			t.Skip("invalid JSON")
		}
		var b struct {
			Values []int `json:"values"`
		}
		err := json.Unmarshal(a, &b)
		if err != nil {
			t.Skip("only correct requests are interestring")
		}

		resp, err := s.Client().Post(s.URL+"/calc", "application/json", bytes.NewReader(a))
		if err != nil {
			t.Errorf("Client.Post: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d; got %d", http.StatusOK, resp.StatusCode)
		}
		// unmarshal
	})
}

func Fuzz_handleListOffset(f *testing.F) {
	s := newServer(f)

	type testcase struct {
		s string
	}
	tt := []testcase{
		{`{"limit":-10,"offset":-10}`},
		{`{"limit":0,"offset":0}`},
		{`{"limit":100,"offset":100}`},
		{`{"limit":200,"offset":200}`},
	}
	for _, tc := range tt {
		f.Add([]byte(tc.s))
	}

	f.Fuzz(func(t *testing.T, a []byte) {
		if !json.Valid(a) {
			t.Skip("invalid JSON")
		}
		var b struct {
			Values []int `json:"values"`
		}
		err := json.Unmarshal(a, &b)
		if err != nil {
			t.Skip("only correct requests are interestring")
		}

		resp, err := s.Client().Post(s.URL+"/list", "application/json", bytes.NewBuffer(a))
		if err != nil {
			t.Errorf("Client.Post: %v", err)
		}
		t.Cleanup(func() { resp.Body.Close() })

		if resp.StatusCode == http.StatusBadRequest {
			t.Skip("invalid json")
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d; got %d", http.StatusOK, resp.StatusCode)
		}
		// unmarshal
	})
}

func newServer(t testing.TB) *httptest.Server {
	ts := httptest.NewServer(Handler())
	t.Cleanup(func() { ts.Close() })
	return ts
}
