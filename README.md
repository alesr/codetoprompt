# [CTP] codetoprompt

[CTP] codetoprompt is a Go package that loads files from a specified directory and outputs their contents in a format that can be used as a prompt for chat-GPT questions. The package can either display the output in the terminal or write it to a file.

The package accepts flags to specify the directory to load files from, the file to output the results to (if any), and which directories or files to exclude. Additionally, the package can optionally include or exclude blank lines in the output.

## Installation

If you have Go installed, you can install the package by running the following command:

```bash
$ go install github.com/alesr/codetoprompt
```

Alternatively, you can download the binary for your platform from the artifacts built by the CI pipeline. The latest version can be found [here](https://github.com/alesr/codetoprompt/releases/tag/v1.0.0)


As a suggestion is to rename the binary to `ctp` and add it to your path.

## Flags

The following flags are available:

```
  -dir string
        The directory to load files from (default ".")
  -exclude string
        A comma-separated list of directories or files to exclude
  -blanklines
        Whether to include blank lines in the output
  -out string
        The file to output the results to
```


## Usage

To use the package, run the following command:

```shell
$ ctp -dir . -out out.txt -exclude go.mod,go.sum,.git,LICENSE,.gitignore,README.md,.github
```

The above command will load all files from the current directory, excluding the files specified in the exclude flag, and output the results to the file specified in the out flag.

### Output Example

The output of the above command will be in the following format:

```
---
Filename: codetoprompt.go

package codetoprompt
import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)
const (
	fileOverwritePrompt string = "The output file already exists. Do you want to overwrite it? (y/n)"
)
var (
	dir               string
	outputPath        string
	exclude           string
	excludeBlankLines bool
	errOutputFileAlreadyExists = errors.New("output file already exists")
)
type file struct {
	name    string
	path    string
	content []byte
}
func parseFlags() {
	flag.StringVar(&dir, "dir", "", "the root directory you want to load files from")
	flag.StringVar(&outputPath, "out", "", "the filepath to the output file, if not provided, the output will be displayed in the terminal")
	flag.StringVar(&exclude, "exclude", "", "exclude a directory or a file from the output file")
	flag.BoolVar(&excludeBlankLines, "blanklines", true, "include blank lines in the output file")
	flag.Parse()
}
func Run() error {
	parseFlags()
	if err := validateFlagsInput(); err != nil {
		if !errors.Is(err, errOutputFileAlreadyExists) {
			return fmt.Errorf("error validating flags: %w", err)
		}
		if !overwriteFile(outputPath) {
			fmt.Println("Aborting...")
			os.Exit(0)
		}
	}
	ignoreList := strings.Split(exclude, ",")
	if outputPath != "" {
		// Adding the output file to the ignore list to avoid writing it to itself
		ignoreList = append(ignoreList, outputPath)
	}
	files, err := loadFiles(dir, ignoreList)
	if err != nil {
		return fmt.Errorf("error loading files: %w", err)
	}
	if outputPath != "" {
		fmt.Println("Writing files to output file...")
		outputFile, err := createFile(outputPath)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		if err := writeFiles(files, outputFile, excludeBlankLines); err != nil {
			return fmt.Errorf("error writing files: %w", err)
		}
		os.Exit(0)
	}
	for _, file := range files {
		fmt.Println("Filename: ", file.name)
		fmt.Println(string(file.content))
		fmt.Println("---")
	}
	return nil
}
func validateFlagsInput() error {
	if dir == "" {
		return errors.New("root path not provided")
	}
	if outputPath != "" {
		alreadyExists, err := fileAlreadyExists(outputPath)
		if err != nil {
			return fmt.Errorf("error checking if output file exists: %w", err)
		}
		if alreadyExists {
			return errOutputFileAlreadyExists
		}
	}
	return nil
}
func fileAlreadyExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if file exists: %w", err)
	}
	return true, nil
}
func overwriteFile(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		var input string
		fmt.Println(fileOverwritePrompt)
		fmt.Scanln(&input)
		switch strings.ToLower(input) {
		case "y", "yes":
			return true
		}
	}
	return false
}
func loadFiles(rootPath string, ignoreList []string) ([]file, error) {
	var files []file
	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not access path %q: %v", path, err)
		}
		if info.IsDir() {
			for _, item := range ignoreList {
				if info.Name() == item {
					return filepath.SkipDir
				}
			}
		}
		if !info.IsDir() {
			for _, item := range ignoreList {
				if info.Name() == item {
					return nil
				}
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read file %q: %v", path, err)
			}
			files = append(files, file{
				name:    info.Name(),
				path:    path,
				content: content,
			})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not walk the path %q: %v", rootPath, err)
	}
	return files, nil
}
func createFile(outputPath string) (*os.File, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("could not create file: %w", err)
	}
	return file, nil
}
func writeFiles(files []file, output io.WriteCloser, includeBlanklines bool) error {
	defer output.Close()
	if _, err := output.Write([]byte("---\n")); err != nil {
		return fmt.Errorf("could not write opening code separator for file: %w", err)
	}
	for _, file := range files {
		var strBulider strings.Builder
		strBulider.WriteString("Filename: ")
		strBulider.WriteString(file.name)
		strBulider.WriteString("\n\n")
		if includeBlanklines {
			for _, line := range strings.Split(string(file.content), "\n") {
				if line == "" {
					continue
				}
				strBulider.WriteString(line)
				strBulider.WriteString("\n")
			}
		} else {
			strBulider.WriteString(string(file.content))
		}
		strBulider.WriteString("---\n")
		if _, err := output.Write([]byte(strBulider.String())); err != nil {
			return fmt.Errorf("could not write file %q: %v", file.name, err)
		}
	}
	return nil
}
---
Filename: main.go

package main
import (
	"log"
	"github.com/alesr/codetoprompt/internal/codetoprompt"
)
func main() {
	if err := codetoprompt.Run(); err != nil {
		log.Fatalln(err)
	}
}
---
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Authors

- [@alesr](https://www.github.com/alesr)




