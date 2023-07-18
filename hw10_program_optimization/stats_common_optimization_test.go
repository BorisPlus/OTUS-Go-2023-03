//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"io"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetDomainFuncSignature func(r io.Reader, domain string) (DomainStat, error)

// go test -v -count=1 -timeout=30s -tags bench .
func TemplateCommonGetDomainStat_Time_And_Memory(t *testing.T, Func GetDomainFuncSignature) {
	bench := func(b *testing.B) {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, 1, len(r.File))

		iodata, err := r.File[0].Open()
		require.NoError(t, err)

		b.StartTimer()
		stat, err := Func(iodata, "biz")
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

func TestCommonGetDomainStatInitial_Time_And_Memory(t *testing.T) {
	fmt.Printf("\nTest with GetDomainStat: Initial\n")
	TemplateCommonGetDomainStat_Time_And_Memory(t, GetDomainStatInitial)
}

func TestCommonGetDomainStatAlternate_Time_And_Memory(t *testing.T) {
	fmt.Printf("\nTest with GetDomainStat: Alternate\n")
	TemplateCommonGetDomainStat_Time_And_Memory(t, GetDomainStatAlternate)
}

func TestCommonGetDomainStatGoroutinedFastJson_Time_And_Memory(t *testing.T) {
	fmt.Printf("\nTest with GetDomainStat: Goroutined + FastJson\n")
	TemplateCommonGetDomainStat_Time_And_Memory(t, GetDomainStatGoroutinedFastJson)
}

