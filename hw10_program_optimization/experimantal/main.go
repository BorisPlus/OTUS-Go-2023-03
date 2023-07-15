package main

import (
	"archive/zip"
	"fmt"
	"os"

	hw10 "github.com/BorisPlus/OTUS-Go-2023-03/hw10_program_optimization"
)

func main() {

	fmt.Println("WORKERS_COUNT", "=", os.Getenv("WORKERS_COUNT"))
	fmt.Println("MAX_CAPACITY ", "=", os.Getenv("MAX_CAPACITY"))

	r, err := zip.OpenReader("../testdata/users.dat.zip")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	data, err := r.File[0].Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	stat, err := hw10.GetDomainStat(data, "biz")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("I get GetDomainStat")
	fmt.Println("Let's check it with ethalon")

	// t := testing.T{}
	// require.Equal(&t, hw10.ExpectedBizStat, stat)

	// LEFT OUTER JOIN
	for key := range stat {
		if stat[key] != hw10.ExpectedBizStat[key] {
			fmt.Println("FAIL")
			print(key)
			return
		}
	}
	// RIGHT OUTER JOIN
	for key := range hw10.ExpectedBizStat {
		if stat[key] != hw10.ExpectedBizStat[key] {
			fmt.Println("FAIL")
			return
		}
	}
	fmt.Println("OK")
}
