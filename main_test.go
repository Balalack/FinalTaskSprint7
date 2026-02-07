package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			http.Error(res, "Не удалось иницировать res-req", http.StatusBadRequest)
			return
		}

		assert.Equal(t, v.status, res.Code)
		assert.Equal(t, v.message, strings.TrimSpace(res.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}

	for _, v := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	}
}

func TestCafeCount(t *testing.T) {
	var town = "moscow"

	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request  string
		response int
	}{
		{fmt.Sprintf("/cafe?count=0&city=%s", town), 0},
		{fmt.Sprintf("/cafe?count=1&city=%s", town), 1},
		{fmt.Sprintf("/cafe?count=2&city=%s", town), 2},
		{fmt.Sprintf("/cafe?count=100&city=%s", town), min(len(cafeList[town]), 100)},
	}

	for _, v := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			http.Error(res, "Не удалось иницировать res-req", http.StatusBadRequest)
			return
		}
		result := res.Body.String()

		if v.response != 0 {
			assert.Len(t, strings.Split(result, ","), v.response)
			continue
		}
		assert.Equal(t, "", result)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request  string
		response int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		resultRequest := fmt.Sprintf("/cafe?search=%s&city=moscow", v.request)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", resultRequest, nil)

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			http.Error(res, "Не удалось иницировать res-req", http.StatusBadRequest)
			return
		}
		result := strings.Split(res.Body.String(), ",")

		if v.response == 0 {
			assert.Equal(t, []string{""}, result)
			continue
		}

		assert.Len(t, result, v.response)
		for _, s := range result {
			assert.Contains(t, strings.ToLower(s), v.request)
		}
	}
}
