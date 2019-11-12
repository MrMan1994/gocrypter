package log

import(
	"log"
	"os"
)

var(
	logger = log.New(os.Stdout, "", 0)
)

func init() {
	logger.Flags()
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Panicf(args ...interface{}) {
	logger.Panicf(args[0].(string), args[1:]...)
}

func Panicln(args ...interface{}) {
	logger.Panicln(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Fatalf(args ...interface{}) {
	logger.Fatalf(args[0].(string), args[1:]...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args)
}

func Print(args ...interface{}) {
	logger.Print(args)
}

func Printf(args ...interface{}) {
	logger.Printf(args[0].(string), args[1:]...)
}

func Println(args ...interface{}) {
	logger.Println(args)
}