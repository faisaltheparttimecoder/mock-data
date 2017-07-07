package core

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/icrowley/fake"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Random text generator based on the length of string needed
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
		return 0, errors.New("Max value is greater or equal to Min value, cannot generate data within these ranges")
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + max, nil
}

// Random Float generator based on precision specified
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RandomFloat(min, max, precision int) (float64, error) {
	output := math.Pow(10, float64(precision))
	randNumber, err := RandomInt(min, max)
	if err != nil {
		return 0.0, err
	}
	return float64(round(float64(randNumber)/rand.Float64()*output)) / output, nil
}

// Random date generator
func RandomDate(fromyear, toyear int) (string, error) {

	if fromyear > toyear {
		return time.Now().Format("2006-01-02 15:04:05.000000"), errors.New("Number of years behind is greater than number of years in future")
	}

	min := time.Now().AddDate(fromyear, 0, 0).Unix()
	max := time.Now().AddDate(toyear, 0, 0).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05.000000"), nil
}

// Random bool generator
func RandomBoolean() bool {
	number, _ := RandomInt(1, 9999)
	if number%2 == 0 {
		return true
	} else {
		return false
	}
}

// Random Paragraphs
func RandomParagraphs() string {
	return fake.ParagraphsN(50)
}

// Random IPv6 & IPv4 Address
func RandomIP() string {
	number, _ := RandomInt(1, 9999)
	if number%2 == 0 {
		return fake.IPv4()
	} else {
		return fake.IPv6()
	}
}

// Random bit
func RandomBit(max int) string {
	var bitValue string
	for i := 0; i < max; i++ {
		if RandomBoolean() {
			bitValue = bitValue + "1"
		} else {
			bitValue = bitValue + "0"
		}
	}
	return bitValue
}

// Random UUID
func RandomUUID() (string, error) {
	// To generate random UUID, we will use unix tool "uuidgen" (unix utility)
	uuidString, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", fmt.Errorf("Unable to run uuidgen to generate UUID data: %v", err)
	}
	return strings.TrimSpace(string(uuidString)), nil
}

// Random Mac Address
func RandomMacAddress() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1))
}
