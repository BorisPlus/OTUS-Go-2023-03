package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		return
	}
	log.Println("os.Args[0]", os.Args[0])
	log.Println("os.Args[1]", os.Args[1])
	log.Println("os.Args[2]", os.Args[2])
	environment, err := ReadEnvDir(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(environment)
	_ = RunCmd(
		os.Args[2:],
		environment,
	)
}
