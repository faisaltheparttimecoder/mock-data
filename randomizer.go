package main

import (
	"errors"
	"math/rand"
	"time"
)

// Random text generator based on the length of string needed
var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Random String generator
func RandomString(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// Random Number generator based on the min and max specified
func RandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, errors.New("max value is greater or equal to Min value, cannot generate data within these ranges")
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min, nil
}
