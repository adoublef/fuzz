// Copyright 2024 Kristopher Rahim Afful-Brown. All Rights Reserved.
//
// Distributed under MIT license.
// See file LICENSE for detail or copy at https://opensource.org/licenses/MIT

package main

import (
	"encoding/json"
	"net/http"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/calc", handleCalcHigh())
	mux.Handle("/list", handleListOffset())
	return mux
}

func handleCalcHigh() http.HandlerFunc {
	type payload struct {
		Values []int `json:"values"`
	}
	parse := func(w http.ResponseWriter, r *http.Request) ([]int, error) {
		v, err := decode[payload](w, r)
		if err != nil {
			return nil, err
		}
		return v.Values, nil
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vv, err := parse(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var high int
		for _, v := range vv {
			if v > high {
				high = v
			}
		}
		if high == 50 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleListOffset() http.HandlerFunc {
	type request struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
	parse := func(w http.ResponseWriter, r *http.Request) (limit, offset int, err error) {
		v, err := decode[request](w, r)
		if err != nil {
			return 0, 0, err
		}
		if v.Limit <= 0 {
			v.Limit = 1
		}
		if v.Offset < 0 {
			v.Offset = 0
		}
		return v.Limit, v.Offset, nil
	}

	type response struct {
		Results    []int `json:"items"`
		PagesCount int   `json:"pagesCount"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		l, o, err := parse(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// apply offsert and limit to static data
		all := make([]int, 1000)
		start := o
		end := o + l
		if o > len(all) {
			start = len(all) - 1
		}
		if end > len(all) {
			end = len(all)
		}
		respond(w, r, response{Results: all[start:end], PagesCount: len(all) / l}, http.StatusOK)
	}
}

func decode[V any](w http.ResponseWriter, r *http.Request) (V, error) {
	defer r.Body.Close()
	var v V
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func respond[V any](w http.ResponseWriter, r *http.Request, data V, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(&data)
}
