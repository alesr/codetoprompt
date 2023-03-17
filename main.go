package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultOutputFile string = "output.txt"

type file struct {
	name    string
	path    string
	content []byte
}

func main() {
	rootPtr := flag.String("root", "", "the filepath to the root directory you want to load files from")
	outputPtr := flag.String("output", defaultOutputFile, "the filepath to the output file")
	includeBlanklines := flag.Bool("blanklines", false, "include blank lines in the output file")
	expludeDirs := flag.String("excludedirs", "", "exclude a directory from the output file")

	flag.Parse()

	if rootPtr == nil || *rootPtr == "" {
		log.Fatal("root flag is required")
	}

	var dirsToIgnore []string
	if expludeDirs != nil && *expludeDirs != "" {
		dirsToIgnore = strings.Split(*expludeDirs, ",")
	}

	var files []file

	if *outputPtr != defaultOutputFile {
		if _, err := os.Stat(*outputPtr); err == nil {
			var input string
			fmt.Println("The output file already exists. Do you want to overwrite it? (y/n)")
			fmt.Scanln(&input)

			switch strings.ToLower(input) {
			case "y", "yes":
				fmt.Println("Overwriting file...")
			default:
				fmt.Println("Exiting...")
				os.Exit(0)
			}
		}
	}

	files, err := loadFiles(*rootPtr, dirsToIgnore)
	if err != nil {
		log.Fatalf("failed to load files: %s", err)
	}

	if len(files) == 0 {
		fmt.Println("no files found")
		os.Exit(0)
	}

	output, err := os.Create(*outputPtr)
	if err != nil {
		log.Fatalf("failed to create output file: %s", err)
	}

	if err := writeFiles(files, output, *includeBlanklines); err != nil {
		log.Fatalf("failed to write files to output file: %s", err)
	}

	outputContent, err := os.ReadFile(*outputPtr)
	if err != nil {
		log.Fatalf("failed to read output file: %s", err)
	}
	fmt.Println(string(outputContent))
}

func loadFiles(root string, dirsToIgnore []string) ([]file, error) {
	var files []file

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not access path %q: %v", path, err)
		}

		if info.IsDir() {
			for _, dir := range dirsToIgnore {
				if info.Name() == dir {
					return filepath.SkipDir
				}
			}
		}

		if !info.IsDir() {
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
		return nil, fmt.Errorf("could not walk the path %q: %v", root, err)
	}
	return files, nil
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

		if !includeBlanklines {
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
