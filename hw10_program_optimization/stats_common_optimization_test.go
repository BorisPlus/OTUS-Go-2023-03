//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetDomainFuncSignature func(r io.Reader, domain string) (DomainStat, error)

func TemplateCommonGetDomainStatTimeAndMemory(t *testing.T, function GetDomainFuncSignature) {
	t.Helper()
	bench := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := function(data, "biz")
		b.StopTimer()
		require.NoError(t, err)

		require.Equal(t, expectedBizStat, stat)
	}

	result := testing.Benchmark(bench)
	mem := result.MemBytes
	t.Logf("time used: %s / %s", result.T, timeLimit)
	t.Logf("memory used: %dMb / %dMb", mem/mb, memoryLimit/mb)

	require.Less(t, int64(result.T), int64(timeLimit), "the program is too slow")
	require.Less(t, mem, memoryLimit, "the program is too greedy")
}

// func TestCommonGetDomainStatInitial_Time_And_Memory(t *testing.T) {
// 	fmt.Printf("\nTest with GetDomainStat: Initial\n")
// 	TemplateCommonGetDomainStatTimeAndMemory(t, GetDomainStatInitial)
// }

func TestCommonGetDomainStatAlternate_Time_And_Memory(t *testing.T) {
	fmt.Printf("\nTest with GetDomainStat: Alternate\n")
	TemplateCommonGetDomainStatTimeAndMemory(t, GetDomainStatAlternate)
}

func TestCommonGetDomainStatGoroutinedFastJson_Time_And_Memory(t *testing.T) {
	fmt.Printf("\nTest with GetDomainStat: Goroutined + FastJson\n")
	TemplateCommonGetDomainStatTimeAndMemory(t, GetDomainStatGoroutinedFastjson)
}
