package logger

import (
	"fmt"
	"log"
)

func Printf(msg string, args ...interface{}) {
	fmt.Printf("[LOG] "+msg, args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}

func Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	fmt.Printf("[WARN] "+msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	log.Panicf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}
