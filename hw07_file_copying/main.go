package main

import (
	"flag"
	"fmt"
)

var (
	from, to               string
	limit, offset, segment int64
	v, perc                bool
	writers                int
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "byte-limit to copy (default: 0)")
	flag.Int64Var(&offset, "offset", 0, "byte-offset of input file (default: 0)")
	//
	flag.BoolVar(&v, "v", false, "verbose log output (default: false)")
	flag.BoolVar(&perc, "perc", false, "indicate percent processing (default: false)")
	//
	flag.Int64Var(&segment, "segment", 1, "byte-segmentation count (default: 1)")
	//
	flag.IntVar(&writers, "writers", 1, "parallel writers count(default: 1)")
	//
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage \"copier.goc\":\n")
		order := []string{
			"from", "to", "limit", "offset", // Сopy
			"segment", "writers", // СopySegmented
			"perc", "v", // Сopy
		}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("-%s\n", flag.Name)
			fmt.Printf("  %s\n", flag.Usage)
		}
		fmt.Println()
		fmt.Printf("Example:\n")

		fmt.Printf(
			`  ./copier.goc -from=%q \
               -to=%q \
               -limit=%d \
               -offset=%d \
               -segment=%d \
               -writers=%d \
               -perc=%s \
               -v=%s
			   `,
			"./testdata/input.txt",
			"./testdata/output.txt",
			1000,
			0,
			200,
			2,
			"true",
			"false",
		)
		fmt.Println()
	}
}

func main() {
	flag.Parse()

	// Copy - UNCOMMENT ME FOR THIS
	// err := Copy(from, to, offset, limit, perc, v)
	err := CopyFast(from, to, offset, limit)
	// params := CopySegmentedParams{
	// 	from,
	// 	to,
	// 	offset,
	// 	limit,
	// 	segment,
	// 	writers,
	// 	perc,
	// 	v,
	// }
	// err := CopySegmented(params)
	if err != nil {
		fmt.Println("Error", err)
	}
}
