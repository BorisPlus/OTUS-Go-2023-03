package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func CopyFast(from, to string, offset, limit int64) error {
	// if !v {
	log.SetOutput(io.Discard)
	// }

	input, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("\"from\" - %w", err)
	}
	defer input.Close()

	fileInfo, err := input.Stat()
	if err != nil {
		return err
	}
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	output, err := os.Create(to)
	if err != nil {
		defer output.Close()
		return fmt.Errorf("\"to\" - %w", err)
	}
	defer output.Close()

	remainder := fileInfo.Size() - offset
	repairLimit := remainder
	if limit != 0 && repairLimit > limit {
		repairLimit = limit
	}
	input.Seek(offset, 0)

	reader := bufio.NewReaderSize(input, int(repairLimit))
	buf := make([]byte, repairLimit)
	log.Println("repairLimit", repairLimit)
	var processed int
	for {
		processed++
		if limit != 0 && processed > int(repairLimit) {
			break
		}

		count, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Println("1")
			log.Println(err)
			return err
		}
		processed += count
		log.Printf("%s", hex.Dump(buf))
		n2, err := output.Write(buf)
		if err != nil {
			return err
		}
		if !v && perc {
			log.Printf("wrote %d bytes\n", n2)
			suffix := "\r"
			if processed == int(repairLimit) {
				suffix = "\n"
			}
			fmt.Fprintf(os.Stdout, "\rCopy ...%.2f%%%s", float32(processed)/float32(repairLimit)*100, suffix)
		}
	}
	return nil
}
