package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"
)

const (
	fromBig                 = "testdata/alice29.text"
	fromInputTxt            = "testdata/input.txt"
	toMyFileNameTemplate    = "testdata/out_offset%d_limit%d_test_copy.txt"
	outputIdentedTemplate   = "testdata/output.%d.txt"
	ethalonFileNameTemplate = "testdata/out_offset%d_limit%d.txt"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	}
	for _, testCase := range testCases {
		toMyFileName := fmt.Sprintf(toMyFileNameTemplate, testCase.offset, testCase.limit)
		ethalonFileName := fmt.Sprintf(ethalonFileNameTemplate, testCase.offset, testCase.limit)
		defer func(file string) {
			err := os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}(toMyFileName)
		err := Copy(
			fromInputTxt,
			toMyFileName,
			testCase.offset,
			testCase.limit, false, false,
		)
		if err != nil {
			fmt.Println("Copy Error", err)
		}
		myFile, errOpenMy := os.Open(toMyFileName)
		if errOpenMy != nil {
			t.Errorf("problem with checked file %q\n", toMyFileName)
		}
		defer myFile.Close()
		ethalonFile, errOpenEthalon := os.Open(ethalonFileName)
		if errOpenEthalon != nil {
			t.Errorf("problem with ethalon file %q\n", ethalonFileName)
		}
		defer ethalonFile.Close()

		hashMyFile := md5.New()
		_, err = io.Copy(hashMyFile, myFile)
		if err != nil {
			panic(err)
		}
		hashEthalonFile := md5.New()
		_, err = io.Copy(hashEthalonFile, ethalonFile)
		if err != nil {
			panic(err)
		}

		if fmt.Sprintf("%v", hashMyFile.Sum(nil)) != fmt.Sprintf("%v", hashEthalonFile.Sum(nil)) {
			t.Errorf("files %s and %s are not the same\n", toMyFileName, ethalonFileName)
		}
		fmt.Printf("OK. Результат соотвествует эталону: %s\n", toMyFileName)
	}
}

func TestCopyFast(t *testing.T) {
	testCases := []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	}
	for _, testCase := range testCases {
		toMyFastFileNameTemplate := "testdata/out_offset%d_limit%d_test_copy_fast.txt"
		toMyFastFileName := fmt.Sprintf(toMyFastFileNameTemplate, testCase.offset, testCase.limit)
		ethalonFileName := fmt.Sprintf(ethalonFileNameTemplate, testCase.offset, testCase.limit)
		defer func(file string) {
			err := os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}(toMyFastFileName)
		err := CopyFast(
			fromInputTxt,
			toMyFastFileName,
			testCase.offset,
			testCase.limit,
		)
		if err != nil {
			fmt.Println("CopyFast Error", err)
		}
		myFile, errOpenMy := os.Open(toMyFastFileName)
		if errOpenMy != nil {
			panic(errOpenMy)
		}
		defer myFile.Close()
		ethalonFile, errOpenEthalon := os.Open(ethalonFileName)
		if errOpenEthalon != nil {
			panic(errOpenEthalon)
		}
		defer ethalonFile.Close()

		hashMyFile := md5.New()
		_, err = io.Copy(hashMyFile, myFile)
		if err != nil {
			panic(err)
		}
		hashEthalonFile := md5.New()
		_, err = io.Copy(hashEthalonFile, ethalonFile)
		if err != nil {
			panic(err)
		}

		if fmt.Sprintf("%v", hashMyFile.Sum(nil)) != fmt.Sprintf("%v", hashEthalonFile.Sum(nil)) {
			t.Errorf("files %s and %s are not the same\n", toMyFastFileName, ethalonFileName)
		}
		fmt.Printf("OK. Результат соотвествует эталону: %s\n", toMyFastFileName)
	}
}

func TestCopySegmented(t *testing.T) {
	testCases := []struct {
		offset       int64
		limit        int64
		segmentSize  int64
		writersCount int
	}{
		{offset: 0, limit: 0, segmentSize: 1, writersCount: 5},
		{offset: 0, limit: 10, segmentSize: 10, writersCount: 2},
		{offset: 0, limit: 1000, segmentSize: 100, writersCount: 1},
		{offset: 0, limit: 10000, segmentSize: 10, writersCount: 5},
		{offset: 100, limit: 1000, segmentSize: 1, writersCount: 1},
		{offset: 6000, limit: 1000, segmentSize: 3, writersCount: 20},
	}
	for _, testCase := range testCases {
		ethalonFileName := fmt.Sprintf(ethalonFileNameTemplate, testCase.offset, testCase.limit)
		toMyFileName := fmt.Sprintf(toMyFileNameTemplate, testCase.offset, testCase.limit)
		defer func(file string) {
			err := os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}(toMyFileName)
		params := CopySegmentedParams{
			fromInputTxt,
			toMyFileName,
			testCase.offset,
			testCase.limit,
			testCase.segmentSize,
			testCase.writersCount,
			false,
			false,
		}
		err := CopySegmented(params)
		if err != nil {
			fmt.Println("Copy Error", err)
		}

		myFile, errOpenMy := os.Open(toMyFileName)
		if errOpenMy != nil {
			panic(errOpenMy)
		}
		defer myFile.Close()

		ethalonFile, errOpenEthalon := os.Open(ethalonFileName)
		if errOpenEthalon != nil {
			panic(errOpenEthalon)
		}
		defer ethalonFile.Close()

		hashMyFile := md5.New()
		_, err = io.Copy(hashMyFile, myFile)
		if err != nil {
			panic(err)
		}
		hashEthalonFile := md5.New()
		_, err = io.Copy(hashEthalonFile, ethalonFile)
		if err != nil {
			panic(err)
		}

		if fmt.Sprintf("%v", hashMyFile.Sum(nil)) != fmt.Sprintf("%v", hashEthalonFile.Sum(nil)) {
			t.Errorf("files %s and %s are not the same\n", toMyFileName, ethalonFileName)
		}
		fmt.Printf("OK. Результат соотвествует эталону: %s\n", toMyFileName)
	}
}

func TestCopySegmentedCustomParams(t *testing.T) {
	testCases := []struct {
		offset       int64
		limit        int64
		segmentSize  int64
		writersCount int
	}{
		{offset: 0, limit: 10000, segmentSize: 1, writersCount: 1},
		{offset: 0, limit: 10000, segmentSize: 3, writersCount: 20},
		{offset: 0, limit: 10000, segmentSize: 10, writersCount: 5},
	}
	ethalonFileName := "testdata/out_offset0_limit10000.txt"
	for id, testCase := range testCases {
		toMyFileSegmentedNameTemplate := "testdata/out_offset%d_limit%d_test_copy_segmented.%d.txt"
		toMyFileSegmentedName := fmt.Sprintf(toMyFileSegmentedNameTemplate, testCase.offset, testCase.limit, id)
		defer func(file string) {
			err := os.Remove(file)
			if err != nil {
				fmt.Println(err)
			}
		}(toMyFileSegmentedName)
		params := CopySegmentedParams{
			fromInputTxt,
			toMyFileSegmentedName,
			testCase.offset,
			testCase.limit,
			testCase.segmentSize,
			testCase.writersCount,
			false,
			false,
		}
		err := CopySegmented(params)
		if err != nil {
			fmt.Println("Copy Error", err)
		}

		myFile, errOpenMy := os.Open(toMyFileSegmentedName)
		if errOpenMy != nil {
			panic(errOpenMy)
		}
		defer myFile.Close()

		ethalonFile, errOpenEthalon := os.Open(ethalonFileName)
		if errOpenEthalon != nil {
			panic(errOpenEthalon)
		}
		defer ethalonFile.Close()

		hashMyFile := md5.New()
		_, err = io.Copy(hashMyFile, myFile)
		if err != nil {
			panic(err)
		}
		hashEthalonFile := md5.New()
		_, err = io.Copy(hashEthalonFile, ethalonFile)
		if err != nil {
			panic(err)
		}

		if fmt.Sprintf("%v", hashMyFile.Sum(nil)) != fmt.Sprintf("%v", hashEthalonFile.Sum(nil)) {
			t.Errorf("files %s and %s are not the same\n", toMyFileSegmentedName, ethalonFileName)
		}
		fmt.Printf("OK. Результат соотвествует эталону: %s\n", toMyFileSegmentedName)
	}
}

func TestCopySegmentedBigFile(t *testing.T) {
	verbose := false
	percentaging := false

	testCases := []struct {
		message      string
		offset       int64
		limit        int64
		segmentSize  int64
		writersCount int
	}{
		{
			message: "Побайтное копирование 50000 байт с 1 врайтером.",
			offset:  0, limit: 50000, segmentSize: 1, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 1 врайтером.",
			offset:  0, limit: 50000, segmentSize: 256, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 4 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 256, writersCount: 4,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 10 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 256, writersCount: 10,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 1 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 500, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 10 врайтером.",
			offset:  0, limit: 50000, segmentSize: 500, writersCount: 10,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 100 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 500, writersCount: 100,
		},
		{
			message: "Копирование 50000 байт с буфером 1000-байт с 5 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 1000, writersCount: 5,
		},
		{
			message: "Копирование 50000 байт с буфером 1000-байт с 50 врайтерами.",
			offset:  0, limit: 50000, segmentSize: 1000, writersCount: 50,
		},
		{
			message: "Копирую 50000 байт в 1 врайтер :)",
			offset:  0, limit: 50000, segmentSize: -1, writersCount: 1,
		},
	}

	ethalonFileName := "./testdata/alice29.ethalon.text"

	for id, testCase := range testCases {
		fmt.Printf("Run [ID %d]: %v\n", id+1, testCase.message)

		const outputIdentedTemplate = "testdata/output.%d.txt"
		toMyFileName := fmt.Sprintf(outputIdentedTemplate, id+1)
		t.Run(
			testCase.message,
			func(t *testing.T) {
				defer func(fileName string) {
					err := os.Remove(fileName)
					if err != nil {
						fmt.Println(err)
					}
				}(toMyFileName)
				params := CopySegmentedParams{
					fromBig,
					toMyFileName,
					testCase.offset,
					testCase.limit,
					testCase.segmentSize,
					testCase.writersCount,
					verbose,
					percentaging,
				}
				err := CopySegmented(params)
				if err != nil {
					fmt.Println("СopySegmented Error", err)
				}

				toMyFile, errOpenMy := os.Open(toMyFileName)
				if errOpenMy != nil {
					panic(errOpenMy)
				}
				defer toMyFile.Close()

				hashMyFile := md5.New()
				_, err = io.Copy(hashMyFile, toMyFile)
				if err != nil {
					panic(err)
				}

				ethalonFile, errOpenEthalon := os.Open(ethalonFileName)
				if errOpenEthalon != nil {
					panic(errOpenEthalon)
				}
				defer ethalonFile.Close()

				hashEthalonFile := md5.New()
				_, err = io.Copy(hashEthalonFile, ethalonFile)
				if err != nil {
					panic(err)
				}

				if fmt.Sprintf("%v", hashMyFile.Sum(nil)) != fmt.Sprintf("%v", hashEthalonFile.Sum(nil)) {
					fmt.Printf("hashMyFile  %v\n", hashMyFile.Sum(nil))
					fmt.Printf("ethalonHash %v\n", hashEthalonFile.Sum(nil))
					t.Errorf("Run [ID %d]: file %s has bad hash-sum\n", id+1, toMyFileName)
				}
				fmt.Printf("ethalonHash %v\n", hashEthalonFile.Sum(nil))
				fmt.Printf("hashMyFile  %v\n", hashMyFile.Sum(nil))
				fmt.Printf("OK. Результат %s соотвествует эталону: %s\n", toMyFileName, ethalonFileName)
			},
		)
	}
}

func BenchmarkCopy(b *testing.B) {
	testCases := []struct {
		message string
		offset  int64
		limit   int64
	}{
		{
			message: "Побайтное копирование 50000 байт [0:5000].",
			offset:  0, limit: 50000,
		},
		{
			message: "Побайтное копирование 50000 байт [100:5100].",
			offset:  100, limit: 50000,
		},
		{
			message: "Побайтное копирование 50000 байт [1000:6000].",
			offset:  1000, limit: 50000,
		},
	}

	for id, testCase := range testCases {
		fmt.Printf("Run [ID %d]: %v\n", id+1, testCase.message)

		toMyFileName := fmt.Sprintf(outputIdentedTemplate, id+1)
		b.Run(
			testCase.message,
			func(b *testing.B) {
				defer func(fileName string) {
					err := os.Remove(fileName)
					if err != nil {
						fmt.Println(err)
					}
				}(toMyFileName)
				err := Copy(
					fromBig,
					toMyFileName,
					testCase.offset,
					testCase.limit,
					false,
					false)
				if err != nil {
					fmt.Println("Copy Error", err)
				}
			},
		)
	}
}

func BenchmarkCopySegmented(b *testing.B) {
	verbose := false
	percentaging := false

	testCases := []struct {
		message      string
		offset       int64
		limit        int64
		segmentSize  int64
		writersCount int
	}{
		{
			message: "Побайтное копирование 50000 байт с 1 врайтером.",
			offset:  100, limit: 50000, segmentSize: 1, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 1 врайтером.",
			offset:  100, limit: 50000, segmentSize: 256, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 4 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 256, writersCount: 4,
		},
		{
			message: "Копирование 50000 байт с буфером 256-байт с 10 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 256, writersCount: 10,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 1 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 500, writersCount: 1,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 10 врайтером.",
			offset:  100, limit: 50000, segmentSize: 500, writersCount: 10,
		},
		{
			message: "Копирование 50000 байт с буфером 500-байт с 100 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 500, writersCount: 100,
		},
		{
			message: "Копирование 50000 байт с буфером 1000-байт с 5 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 1000, writersCount: 5,
		},
		{
			message: "Копирование 50000 байт с буфером 1000-байт с 50 врайтерами.",
			offset:  100, limit: 50000, segmentSize: 1000, writersCount: 50,
		},
		{
			message: "Смотрите как быстро копирую 50000 байт в 1 врайтер :)",
			offset:  100, limit: 50000, segmentSize: -1, writersCount: 1,
		},
	}

	for id, testCase := range testCases {
		fmt.Printf("Run [ID %d]: %v\n", id+1, testCase.message)

		toMyFileName := fmt.Sprintf(outputIdentedTemplate, id+1)
		b.Run(
			testCase.message,
			func(b *testing.B) {
				defer func(fileName string) {
					err := os.Remove(fileName)
					if err != nil {
						fmt.Println(err)
					}
				}(toMyFileName)
				params := CopySegmentedParams{
					fromBig,
					toMyFileName,
					testCase.offset,
					testCase.limit,
					testCase.segmentSize,
					testCase.writersCount,
					verbose,
					percentaging,
				}
				err := CopySegmented(params)
				if err != nil {
					fmt.Println("СopySegmented Error", err)
				}
			},
		)
	}
}

func BenchmarkCopyFast(b *testing.B) {
	testCases := []struct {
		message       string
		offset, limit int64
	}{
		{
			message: "Копирование 50000 байт с отступом 0.",
			offset:  0, limit: 50000,
		},
		{
			message: "Копирование 50000 байт с отступом 100.",
			offset:  100, limit: 50000,
		},
		{
			message: "Копирование 50000 байт с отступом 10000.",
			offset:  10000, limit: 50000,
		},
	}

	for id, testCase := range testCases {
		fmt.Printf("Run [ID %d]: %v\n", id+1, testCase.message)
		toMyFileName := fmt.Sprintf(outputIdentedTemplate, id+1)
		b.Run(
			testCase.message,
			func(b *testing.B) {
				defer func(fileName string) {
					err := os.Remove(fileName)
					if err != nil {
						fmt.Println(err)
					}
				}(toMyFileName)
				err := CopyFast(
					fromBig,
					toMyFileName,
					testCase.offset,
					testCase.limit,
				)
				if err != nil {
					fmt.Println("СopyFast Error", err)
				}
			},
		)
	}
}
