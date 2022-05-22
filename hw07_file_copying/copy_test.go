package main

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inputFilePath     = "testdata/input.txt"
	resultFileInfo, _ = os.Stat(inputFilePath)
	inputFileSize     = resultFileInfo.Size()
)

func TestCountRealLimit(t *testing.T) {
	halfFileSize := int64(math.Round(float64(inputFileSize / 2)))
	smallNumber := int64(5)

	type testCase struct {
		testName      string
		fileName      string
		inputLimit    int64
		offset        int64
		limitExpected int64
		errorExpected error
	}

	testCases := []testCase{
		{testName: "Default empty flags", limitExpected: inputFileSize},
		{testName: "Small limit", inputLimit: smallNumber, limitExpected: smallNumber},
		{testName: "Limit without offset", inputLimit: halfFileSize, limitExpected: halfFileSize},
		{testName: "Limit more than file size", inputLimit: inputFileSize + halfFileSize, limitExpected: inputFileSize},
		{testName: "Middle of file", inputLimit: halfFileSize, offset: smallNumber, limitExpected: halfFileSize},
		{
			testName:      "From half to EOF",
			inputLimit:    halfFileSize,
			offset:        halfFileSize + smallNumber,
			limitExpected: inputFileSize - (halfFileSize + smallNumber),
		},
		{testName: "Skip half of file without limit", offset: halfFileSize, limitExpected: inputFileSize - halfFileSize},
		// Tests with expected error.
		{testName: "Illegal offset", offset: inputFileSize + smallNumber, errorExpected: ErrOffsetExceedsFileSize},
		{testName: "Read error", fileName: "/tmp/fileNotExist.name", errorExpected: os.ErrNotExist},
		{testName: "Unsupported file", fileName: "/dev/urandom", errorExpected: ErrUnsupportedFile},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			if testCase.fileName == "" {
				testCase.fileName = inputFilePath
			}
			limit, err := countRealLimit(testCase.fileName, testCase.inputLimit, testCase.offset)
			require.Equal(t, testCase.limitExpected, limit)

			if testCase.errorExpected != nil {
				require.ErrorIs(t, err, testCase.errorExpected)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "")
	toPath := tempFile.Name()
	defer os.Remove(toPath)

	t.Run("Full copy", func(t *testing.T) {
		err := Copy(inputFilePath, toPath, 0, 0)
		require.NoError(t, err)

		resultFileInfo, _ = os.Stat(toPath)
		require.Equal(t, inputFileSize, resultFileInfo.Size())
	})

	type testCase struct {
		testName       string
		offset         int64
		limit          int
		expectedResult string
	}

	var (
		headPhrase   = "Go\nDocuments"
		middlePhrase = "Documents"
		tailPhrase   = "\nSupported by Google\n"
	)

	testCases := []testCase{
		{"Head", 0, len(headPhrase), headPhrase},
		{"Middle", 3, len(middlePhrase), middlePhrase},
		{"Tail", inputFileSize - int64(len(tailPhrase)), len(tailPhrase) * 2, tailPhrase},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := Copy(inputFilePath, toPath, testCase.offset, int64(testCase.limit))
			require.NoError(t, err)

			result, err := os.ReadFile(toPath)
			require.NoError(t, err)
			require.Equal(t, testCase.expectedResult, string(result))
		})
	}
}
