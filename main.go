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
	getDir(path, filesFromDisk, files, 0, &text)

	fmt.Fprintln(out, text)

	return nil
}

func getDir(path string, filesFromDisk []os.DirEntry, files bool, countTab int, text *string) {
	lenFiles := len(filesFromDisk)
	for i := 0; i < lenFiles - 1; i++  {
		if filesFromDisk[i].IsDir() {
			*text +=  newLine + getTab(countTab) + verticalSlash + firstSymbol + filesFromDisk[i].Name()
			newPath := path+`\`+filesFromDisk[i].Name()
			newFilesFromDisk, _ := os.ReadDir(newPath)
			getDir(newPath, newFilesFromDisk, files, countTab+1, text)
		} else if files {
			getFile(countTab + 1, filesFromDisk[i], text, firstSymbol)
		}
	}

	if lenFiles > 0 && filesFromDisk[lenFiles - 1].IsDir() {
		newPath := path+`\`+filesFromDisk[lenFiles - 1].Name()
		newFilesFromDisk, _ := os.ReadDir(newPath)
		if len(newFilesFromDisk) > 1 {
			*text +=  newLine + getTabLast(countTab)  + lastSymbol + filesFromDisk[lenFiles - 1].Name()
		} else  {
			*text +=  newLine + getTab(countTab)  + lastSymbol + filesFromDisk[lenFiles - 1].Name()
		}
		getDir(newPath, newFilesFromDisk, files, countTab+1, text)
	} else if files {
		getFileLast(countTab + 1, filesFromDisk[lenFiles - 1], text, lastSymbol)
	}
}

func getFile(path int, file os.DirEntry, text *string, symbol string) {
	fileInfo, _ := file.Info()
	size := strconv.Itoa(int(fileInfo.Size()))
	*text += newLine + getTab(path) + tabSymbol + symbol + file.Name() + " (" + size +"b)" + newLine
}

func getFileLast(path int, file os.DirEntry, text *string, symbol string) {
	fileInfo, _ := file.Info()
	size := strconv.Itoa(int(fileInfo.Size()))
	*text += newLine + getTabLast(path) + tabSymbol + symbol + file.Name() + " (" + size +"b)" + newLine
}

func getFileName(countTab int, countTabLast int, file os.DirEntry, text *string) {
	fileInfo, _ := file.Info()
	size := int(fileInfo.Size())
	symbol := verticalSlash+firstSymbol
	if countTabLast > 0 {
		symbol = lastSymbol
	}
	sizeText := ""
	if size > 0 {
		sizeText = " (" + strconv.Itoa(int(fileInfo.Size())) +"b)"
	} else {

		sizeText = " (empty)"
	}
	*text += newLine + getTab(countTab) +  getTabLast(countTabLast) + tabSymbol + symbol + file.Name() + sizeText + newLine
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