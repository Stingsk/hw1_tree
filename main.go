package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const tabSymbol = `	`
const verticalSlash = `|`
const newLine = `
`
const firstSymbol = `───`
const lastSymbol = `└───`

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, files bool) error {
	var text = ""
	filesFromDisk, _ := os.ReadDir(path)
	getDir(path, filesFromDisk, files, 0, 0, &text)

	fmt.Fprintln(out, text)

	return nil
}

func getDir(path string, filesFromDisk []os.DirEntry, files bool, countTab int, countTabLast int, text *string) {
	lenFiles := len(filesFromDisk)
	for i := 0; i < lenFiles-1; i++ {
		if filesFromDisk[i].IsDir() {
			getDirName(countTab, countTabLast, filesFromDisk[i], text)
			newPath := path + `\` + filesFromDisk[i].Name()
			newFilesFromDisk, _ := os.ReadDir(newPath)
			getDir(newPath, newFilesFromDisk, files, countTab+1, countTabLast, text)
		} else if files {
			getFileName(countTab, countTabLast, filesFromDisk[i], text)
		}
	}

	if lenFiles > 0 && filesFromDisk[lenFiles-1].IsDir() {
		newPath := path + `\` + filesFromDisk[lenFiles-1].Name()
		newFilesFromDisk, _ := os.ReadDir(newPath)
		getDirName(countTab, countTabLast, filesFromDisk[lenFiles-1], text)
		getDir(newPath, newFilesFromDisk, files, countTab, countTabLast+1, text)
	} else if files {
		getFileName(countTab, countTabLast, filesFromDisk[lenFiles-1], text)
	}
}

func getFileName(countTab int, countTabLast int, file os.DirEntry, text *string) {
	fileInfo, _ := file.Info()
	size := int(fileInfo.Size())
	sizeText := ""
	if size > 0 {
		sizeText = " (" + strconv.Itoa(int(fileInfo.Size())) + "b)"
	} else {
		sizeText = " (empty)"
	}
	*text += getName(countTab, countTabLast, file) + sizeText + newLine
}

func getDirName(countTab int, countTabLast int, file os.DirEntry, text *string) {
	*text += getName(countTab, countTabLast, file) + newLine
}

func getName(countTab int, countTabLast int, file os.DirEntry) string {
	symbol := verticalSlash + firstSymbol
	if countTabLast > 0 {
		symbol = lastSymbol
	}
	return getTab(countTab) + getTabLast(countTabLast) + symbol + file.Name()
}

func getTab(count int) string {
	text := ""
	for i := 0; i < count; i++ {
		text += verticalSlash + tabSymbol
	}
	return text
}

func getTabLast(count int) string {
	text := ""
	for i := 0; i < count; i++ {
		text += tabSymbol
	}
	return text
}
