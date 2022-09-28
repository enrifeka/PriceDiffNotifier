package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

func getPrices(endpoints []Endpoint, requestTimeout time.Duration) []PostRes {
	resChan := make(chan PostRes)
	for _, e := range endpoints {
		go func(endpoint Endpoint, reqTimeout time.Duration) {
			resChan <- getPrice(endpoint, reqTimeout)
		}(e, requestTimeout)
	}
	var postRes []PostRes
	for range endpoints {
		res := <-resChan
		postRes = append(postRes, res)
	}
	return postRes
}

func getPrice(e Endpoint, requestTimeout time.Duration) PostRes {
	client := http.Client{
		Timeout: requestTimeout,
	}
	res, err := client.Get(e.PriceAPI)
	if err != nil {
		return PostRes{Endpoint: e, Err: err}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return PostRes{Endpoint: e, Err: err}
	}
	value := gjson.Get(string(body), e.PriceTagPath)
	if !value.Exists() {
		return PostRes{Endpoint: e, Err: errors.New("price not found in path: " + e.PriceTagPath)}
	}
	no, err := strconv.ParseFloat(value.String(), 64)
	if err != nil {
		return PostRes{Endpoint: e, Err: err}
	}
	e.ParsedPrice = no
	return PostRes{Endpoint: e}
}

type PostRes struct {
	Endpoint Endpoint
	Err      error
}
