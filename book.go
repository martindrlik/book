package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	L Log

	voucher = make(map[string]string)
	ch      = make(chan struct{ k, v string }, 100)
)

func main() {
	go func() {
		for d := range ch {
			voucher[d.k] = d.v
			L.Debug("book", d)
		}
	}()
	http.HandleFunc("/book", book)
	http.HandleFunc("/verify", verify)
	panic(http.ListenAndServe(":8085", nil))
}

func book(w http.ResponseWriter, r *http.Request) {
	d := struct {
		Token string `json:"token"`
	}{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		L.Warn("book decode", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := voucher[d.Token]; ok {
		L.Debug("book already booked", d)
		http.Error(w, `{"error":"already booked"}`, http.StatusFound)
		return
	}
	v, err := genVoucher()
	if err != nil {
		L.Error("book voucher", err)
		http.Error(w, "Not your fault. Try book it later.", http.StatusInternalServerError)
		return
	}
	ch <- struct{ k, v string }{d.Token, v}
	fmt.Fprintf(w, `{"voucher":%q}`, v)
}

func verify(w http.ResponseWriter, r *http.Request) {
	d := struct {
		Token   string `json:"token"`
		Voucher string `json:"voucher"`
	}{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&d)
	if err != nil {
		L.Warn("verify decode", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, ok := voucher[d.Token]
	if !ok {
		L.Debug("verify never booked", d)
		http.Error(w, `{"error":"never booked"}`, http.StatusFound)
		return
	}
	if v != d.Voucher {
		L.Debug("verify never booked", d)
		http.Error(w, `{"error":"bad voucher"}`, http.StatusFound)
		return
	}
	fmt.Fprint(w, `{"message":"all seems fine"}`)
}
