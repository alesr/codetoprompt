package codetoprompt

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFiles(t *testing.T) {
	testCases := []struct {
		name          string
		givenRootPath string
		givenIgnore   []string
		expectedFiles []file
		expectedError error
	}{
		{
			name:          "no root path",
			givenRootPath: "",
			givenIgnore:   []string{},
			expectedFiles: nil,
			expectedError: errNoRootPathProvided,
		},
		{
			name:          "success",
			givenRootPath: "testdata",
			givenIgnore:   []string{"dir_to_ignore_qux", "file_to_ignore_bar.txt"},
			expectedFiles: []file{
				{
					name:    "file_foo.txt",
					path:    "testdata/file_foo.txt",
					content: []byte("line 1\n\nline 3\n"),
				},
				{
					name:    "file_bar.txt",
					path:    "testdata/dir_foo/file_bar.txt",
					content: []byte("line 1\n\nline 3\n\nline 5\nline 6\n"),
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			observedFiles, err := loadFiles(tc.givenRootPath, tc.givenIgnore)

			if tc.expectedError == nil {
				require.NoError(t, err)

				for _, expectedFile := range tc.expectedFiles {
					assert.Contains(t, observedFiles, expectedFile)
				}

			} else {
				assert.True(t, errors.Is(err, tc.expectedError))
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	testCases := []struct {
		name          string
		givenFilePath string
		expectedError error
	}{
		{
			name:          "no file path",
			givenFilePath: "",
			expectedError: errNoFilePathProvided,
		},
		{
			name:          "file is created",
			givenFilePath: "foo",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := createFile(tc.givenFilePath)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tc.expectedError))
			}

			if tc.givenFilePath != "" && tc.expectedError == nil {
				require.NoError(t, os.Remove(tc.givenFilePath))
			}
		})
	}
}

func TestWriteFiles(t *testing.T) {
	testCases := []struct {
		name                   string
		givenFiles             []file
		givenOutput            io.Writer
		givenIncludeBlanklines bool
		expectedError          error
	}{
		{
			name:                   "no files",
			givenFiles:             nil,
			givenOutput:            nil,
			givenIncludeBlanklines: false,
			expectedError:          errNoFilesToWrite,
		},
		{
			name:                   "no output",
			givenFiles:             []file{{name: "test", path: "test", content: []byte("test")}},
			givenOutput:            nil,
			givenIncludeBlanklines: false,
			expectedError:          errNoOutputFileProvided,
		},
		{
			name:                   "no blank lines",
			givenFiles:             []file{{name: "test", path: "test", content: []byte("test")}},
			givenOutput:            bytes.NewBuffer(nil),
			givenIncludeBlanklines: false,
			expectedError:          nil,
		},
		{
			name:                   "with blank lines",
			givenFiles:             []file{{name: "test", path: "test", content: []byte("test")}},
			givenOutput:            bytes.NewBuffer(nil),
			givenIncludeBlanklines: true,
			expectedError:          nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := writeFiles(tc.givenFiles, tc.givenOutput, tc.givenIncludeBlanklines)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tc.expectedError))
			}

			// Check if the output file contains the expected content
			if tc.givenOutput != nil {
				output := tc.givenOutput.(*bytes.Buffer)
				assert.Contains(t, output.String(), "---")
				assert.Contains(t, output.String(), "Filename: test")
				assert.Contains(t, output.String(), "test")
			}
		})
	}
}
