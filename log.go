package main

import (
	"fmt"
)

type Log struct{}

func (Log) Error(a ...interface{}) {
	fmt.Print("ERROR: ")
	fmt.Println(a...)
}

func (Log) Warn(a ...interface{}) {
	fmt.Print("WARN: ")
	fmt.Println(a...)
}

func (Log) Debug(a ...interface{}) {
	fmt.Print("DEBUG: ")
	fmt.Println(a...)
}
