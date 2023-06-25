package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	ErrNotPositiveWritersCount = errors.New("writers count must be positive")
	ErrReadSegment             = errors.New("segment read error")
)

type Segment struct {
	data   []byte
	offset int64
}

func percenteger(part, full int64) {
	suffix := "\r"
	if part == full {
		suffix = "\n"
		// fmt.Fprintf(os.Stdout, "\rCopySegmented ...%.2f%%%s", 100., suffix)
		// return
	}
	fmt.Fprintf(os.Stdout, "\rCopySegmented ...%.2f%%%s", float32(part)/float32(full)*100, suffix)
	// fmt.Println()
}

func percentbarRunner(waitGroup *sync.WaitGroup, percent100value int64, equitiesChan <-chan int64, percentaging bool) {
	defer waitGroup.Done()
	var equitysum int64
	for equity := range equitiesChan {
		if percentaging {
			equitysum += equity
			percenteger(equitysum, percent100value)
		}
	}
}

type CopySegmentedParams struct {
	from, to               string
	offset, limit, segment int64
	writers                int
	perc, v                bool
}

func writer(number int, waitGroup *sync.WaitGroup,
	// mutex *sync.RWMutex,
	segmentsForWriteAt <-chan Segment,
	writedPercentages chan<- int64, target *os.File,
) {
	defer func() {
		waitGroup.Done()
	}()
	for segment := range segmentsForWriteAt {
		log.Printf("WRITER (No. %d): offset = %d\n", number, segment.offset)
		log.Printf("WRITER (No. %d): data\n%s\n", number, hex.Dump(segment.data))
		log.Printf("WRITER (No. %d): must be len() = %d\n", number, len(segment.data))

		// TODO: IT IS NO NEED - mutex - OS RULEZ
		// mutex.Lock()
		// mutex.Unlock()
		count, err := target.WriteAt(segment.data, segment.offset)
		// TODO: IT IS NO NEED
		// target.Sync()
		if err != nil {
			log.Println(err)
			log.Printf("Error %q\n", err)
		}

		log.Printf("WRITER (No. %d): wrotet len()=%d\n", number, count)

		writedPercentages <- int64(count)
	}
}

func CopySegmented(params CopySegmentedParams) error {
	if !params.v {
		log.SetOutput(io.Discard)
	}
	if params.writers < 1 {
		return ErrNotPositiveWritersCount
	}

	log.Printf("CopySegmented")
	log.Printf("from = %s\n", params.from)
	input, err := os.Open(params.from)
	if err != nil {
		defer input.Close()
		log.Println(err)
		// return fmt.Errorf("\"to\" - %q", err.Error())
		return err
	}
	defer input.Close()

	fileInfo, err := input.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	if params.offset > fileInfo.Size() {
		log.Println("ErrOffsetExceedsFileSize")
		return ErrOffsetExceedsFileSize
	}

	repairLimit := fileInfo.Size() - params.offset
	log.Printf("offset = %d\n", params.offset)
	log.Printf("limit = %d\n", params.limit)

	if params.limit == 0 {
		params.limit = repairLimit
	}

	if params.limit != 0 && repairLimit > params.limit {
		repairLimit = params.limit
	}
	log.Printf("fileInfo.Size() = %d\n", fileInfo.Size())
	log.Printf("repairLimit = %d\n", repairLimit)

	if params.segment < 1 {
		params.segment = repairLimit
		log.Printf("segmentSize = %d\n", params.segment)
	}

	output, err := os.Create(params.to)
	if err != nil {
		defer output.Close()
		log.Println(err)
		// return fmt.Errorf("\"to\" - %q", err.Error())
		return err
	}
	if err := output.Truncate(repairLimit); err != nil {
		log.Println(err)
		return err
	}
	output.Sync()

	defer output.Close()

	var n int64 = 1

	segments := make(chan Segment)
	percentages := make(chan int64)
	wgWriters := sync.WaitGroup{}
	// wgMutex := sync.RWMutex{}

	for i := 0; i < params.writers; i++ {
		wgWriters.Add(1)
		// go func(number int, waitGroup *sync.WaitGroup,
		// 	mutex *sync.RWMutex, segmentsForWriteAt <-chan Segment,
		// 	writedPercentages chan<- int64, target *os.File,
		// ) {
		// 	defer func() {
		// 		waitGroup.Done()
		// 	}()
		// 	for segment := range segmentsForWriteAt {
		// 		log.Printf("WRITER (No. %d): offset = %d\n", number, segment.offset)
		// 		log.Printf("WRITER (No. %d): data\n%s\n", number, hex.Dump(segment.data))
		// 		log.Printf("WRITER (No. %d): must be len() = %d\n", number, len(segment.data))

		// 		// TODO: IT IS NO NEED - mutex - OS RULEZ
		// 		// mutex.Lock()
		// 		// mutex.Unlock()
		// 		_ = mutex
		// 		count, err := target.WriteAt(segment.data, segment.offset)
		// 		// TODO: IT IS NO NEED
		// 		// target.Sync()
		// 		if err != nil {
		// 			log.Println(err)
		// 			log.Printf("Error %q\n", err)
		// 		}

		// 		log.Printf("WRITER (No. %d): wrotet len()=%d\n", number, count)

		// 		writedPercentages <- int64(count)
		// 	}
		// }(i, &wgWriters, &wgMutex, segments, percentages, output)
		go writer(i, &wgWriters, segments, percentages, output)
	}

	wgPercenter := sync.WaitGroup{}
	wgPercenter.Add(1)
	// go func(waitGroup *sync.WaitGroup, percent100value int64, equitiesChan <-chan int64, percentaging bool) {
	// 	defer waitGroup.Done()
	// 	var equitysum int64
	// 	for equity := range equitiesChan {
	// 		if percentaging {
	// 			equitysum += equity
	// 			percenteger(equitysum, percent100value)
	// 		}
	// 	}
	// }(&wgPercenter, repairLimit, percentages, !params.v && params.perc)

	go percentbarRunner(&wgPercenter, repairLimit, percentages, !params.v && params.perc)

	input.Seek(params.offset, 0)

	var prevSegmentsSizesSum int64
	partition := params.segment

	// Ужс
	reader := bufio.NewReaderSize(input, int(repairLimit))

	number := 0
	for {
		number++
		log.Println()
		log.Printf("READER (No. %d): initial segment size = %d\n", number, params.segment)

		log.Printf("STEP READ n = %d\n", n)
		log.Printf("STEP READ n = %d - repairLimit = %d\n", n, repairLimit)
		log.Printf("STEP READ n = %d - inisegmentSize = %d\n", n, params.segment)
		log.Printf("STEP READ n = %d - prevSegmentsSizesSum = %d\n", n, prevSegmentsSizesSum)

		if (prevSegmentsSizesSum + params.segment) > repairLimit {
			partition = repairLimit - prevSegmentsSizesSum
		}
		log.Printf("READER (No. %d): segment size = %d\n", number, partition)
		buf := make([]byte, partition)
		log.Printf("READER (No. %d): expected read = %d\n", number, len(buf))
		count, err := reader.Read(buf)
		log.Printf("READER (No. %d): result read = %d\n", number, count)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Println(err)
			return err
		}
		if count == 0 {
			log.Printf("ErrReadSegment")
			return ErrReadSegment
		}
		newSegment := Segment{data: buf, offset: prevSegmentsSizesSum}
		segments <- newSegment
		log.Printf("READER (No. %d): put in channel-segment with len() %d, at offset %d with data\n%s.",
			number, len(newSegment.data), prevSegmentsSizesSum, hex.Dump(newSegment.data))

		prevSegmentsSizesSum += int64(count)

		if prevSegmentsSizesSum == repairLimit {
			close(segments)
			break
		}
		n++
	}

	log.Println("all writers Wait()")
	wgWriters.Wait()
	log.Println("close(percents)")
	close(percentages)
	log.Println("percenter Wait()")
	wgPercenter.Wait()
	return nil
}
