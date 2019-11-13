package log

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

func init() {
	logger.Flags()
}

// Panic calls the logger's Panic function with the provided args
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf calls the logger's Panicf function with the provided args
func Panicf(args ...interface{}) {
	logger.Panicf(args[0].(string), args[1:]...)
}

// Panicln calls the logger's Panicln function with the provided args
func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

// Fatal calls the logger's Fatal function with the provided args
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf calls the logger's Fatalf function with the provided args
func Fatalf(args ...interface{}) {
	logger.Fatalf(args[0].(string), args[1:]...)
}

// Fatalln calls the logger's Fatalln function with the provided args
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

// Print calls the logger's Print function with the provided args
func Print(args ...interface{}) {
	logger.Print(args...)
}

// Printf calls the logger's Printf function with the provided args
func Printf(args ...interface{}) {
	logger.Printf(args[0].(string), args[1:]...)
}

// Println calls the logger's Println function with the provided args
func Println(args ...interface{}) {
	logger.Println(args...)
}
