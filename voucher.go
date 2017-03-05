package main

import (
	"crypto/rand"
	"encoding/base64"
)

func genVoucher() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(b)
	return s, nil
}
