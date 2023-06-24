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

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(from, to string, offset, limit int64, perc, v bool) error {
	if !v {
		log.SetOutput(io.Discard)
	}

	input, err := os.Open(from)
	if err != nil {
		defer input.Close()
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
	if repairLimit > limit {
		repairLimit = limit
	}

	input.Seek(offset, 0)

	reader := bufio.NewReader(input)
	buf := make([]byte, 1)
	var processed int64
	for {
		processed++
		if limit != 0 && processed > limit {
			break
		}

		_, err := reader.Read(buf)
		if err != nil {
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Println(err)
				return err
			}
		}

		log.Printf("%s", hex.Dump(buf))
		n2, err := output.Write(buf)
		if err != nil {
			return err
		}
		if !v && perc {
			log.Printf("wrote %d bytes\n", n2)
			suffix := "\r"
			if processed == repairLimit {
				suffix = "\n"
			}
			fmt.Fprintf(os.Stdout, "\rCopy ...%.2f%%%s", float32(processed)/float32(repairLimit)*100, suffix)
		}
	}
	return nil
}
