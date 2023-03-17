# fileloader

File Loader is a command-line tool that loads the content of all the files in a directory (excluding those in excluded directories) and writes them to a single output file. It includes an option to exclude blank lines in the output file.

It very useful for GPT prompts =]

## Installation
To install File Loader, run the following command:

```go
go get github.com/alesr/fileloader
```

## Usage

```shell
fileloader -root="pathttomyrepo"
```

#### Flags

-root: The filepath to the root directory you want to load files from (required)
-output: The filepath to the output file (default: output.txt)
-blanklines: Include blank lines in the output file (default: false)
-excludedirs: Exclude a directory from the output file (comma-separated list)

##### Example
```shell
Copy code
file-loader -root ./mydirectory -output ./output.txt -blanklines true -excludedirs node_modules,.git
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
