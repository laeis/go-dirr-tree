package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// print symbols
const (
	indent        = "\t"
	tShaped       = "├───"
	end           = "└───"
	connectingine = "|"
)

//filterDirectory remove files from slice , if do not need display they
func filterDirectory(slice *[]os.FileInfo, withFiles bool) {
	if !withFiles { // transform the slise only if files don't need
		for i := 0; ; {
			if i >= len(*slice) { // exit if slice don't change
				break
			}
			if !(*slice)[i].IsDir() { //if curent element is file(not directory) remove it
				*slice = append((*slice)[:i], (*slice)[i+1:]...)
			} else { //else go to next element
				i++
			}
		}
	}
}

//dirTree entry point for traversing the directory tree
func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	goalDir, err := filepath.Abs(path) //get absolute directory path
	if err != nil {
		return fmt.Errorf("Directory can be set as absolute")
	}
	//recursive function for traversing the directory tree
	err = recursivePrintNodes(out, goalDir, printFiles, 0, "", false)
	return
}

//getPrefixSymbol return prefix for child folder
func getPrefixSymbol(oldprefix string, level int, isLast bool) (newPrefix string) {
	if level != 0 {
		if !isLast {
			newPrefix = oldprefix + "│\t"
		} else {
			newPrefix = oldprefix + "\t"
		}
	}
	return
}

//getNodeString return ready for print string with symbols, name and size
func getNodeString(fileInfo os.FileInfo, isLastNode bool, beforePrefix string) string {
	var prefix, size string
	prefix += beforePrefix
	if isLastNode {
		prefix += end
	} else {
		prefix += tShaped
	}
	result := prefix + fileInfo.Name()
	if !fileInfo.IsDir() {
		size = formatSize(fileInfo.Size())
	}
	if len(size) > 0 {
		result += " " + size
	}
	result += "\n"
	return result
}

// formatSize return formated string for file size
func formatSize(bite int64) (size string) {
	if bite == 0 {
		size = "(empty)"
	} else {
		size = fmt.Sprintf("(%db)", bite)
	}
	return
}

//recursivePrintNodes traversing the directory tree and print it elements
func recursivePrintNodes(
	output io.Writer,
	path string,
	printFiles bool,
	level int,
	prefix string,
	isLast bool) (err error) {
	//create path for use in a recursive call
	pathPrefix := path + string(os.PathSeparator)
	//get *File from name
	file, err := os.Open(path)
	if err != nil {
		return
	}
	//read current file directory
	directiryList, err := file.Readdir(0)
	file.Close()
	if err != nil {
		return
	}
	//sort directory clice by name
	sort.Slice(directiryList, func(i, j int) bool { return directiryList[i].Name() < directiryList[j].Name() })
	// files are removing if it necessary from directory slice
	filterDirectory(&directiryList, printFiles)
	//get prefix for result string
	prefix = getPrefixSymbol(prefix, level, isLast)
	//get slice length for use
	for pos, item := range directiryList {
		// prefix = getPrefixSymbol(prefix, len(directiryList)-1 == pos, level+1)
		last := len(directiryList)-1 == pos
		//elements print  from from sorted and filtered slice
		fmt.Fprintf(output, getNodeString(item, last, prefix))
		//if current elemnt is directory then call function for it
		if item.IsDir() {
			pathName := pathPrefix + item.Name()
			recursivePrintNodes(output, pathName, printFiles, level+1, prefix, last)
		}

	}
	return
}

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
