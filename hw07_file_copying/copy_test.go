package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	os.Mkdir("tmp", 0755)
	defer os.RemoveAll("tmp")

	// 1. запускаем без параметров
	runTestCase(t, 0, 0)

	// 2. проверяем offset
	runTestCase(t, 10, 0)

	// 3. проверяем limit
	runTestCase(t, 0, 10)

	// 4. проверяем offset и limit
	runTestCase(t, 100, 1000)
}

func runTestCase(t *testing.T, offset, limit int64) {
	testOutFilePath := fmt.Sprintf("tmp/test_out_%d_%d.txt", offset, limit)
	correctOutFilePath := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", offset, limit)
	err := Copy(
		"testdata/input.txt",
		testOutFilePath,
		offset,
		limit,
	)
	require.NoError(t, err)
	correctOut := getMD5SumString(correctOutFilePath)
	testOut := getMD5SumString(testOutFilePath)
	require.Equal(t, correctOut, testOut)
}

func getMD5SumString(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}

	file1Sum := md5.New()

	_, err = io.Copy(file1Sum, f)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%X", file1Sum.Sum(nil))
}
