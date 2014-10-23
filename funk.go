package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Declare function definers
var FUNC_DECS []string = []string{
	"function",
	"def",
	"func",
}

var FuncRefs []FileRef = []FileRef{}

type FileRef struct {
	FilePath string
	Text     string
	LineNum  string
}

func (f FileRef) String() string {
	return fmt.Sprintf("File: %v LineNum: %v\n", f.FilePath, f.LineNum)
}

func main() {
	dirpath := flag.String("d", ".", "Path to directory to search")
	flag.Parse()
	dir := filepath.Dir(*dirpath + "/")
	fmt.Println(dir)
	traverseDir(dir)
	// Now we have collected all the func names
	for i := range FuncRefs {
		ref := FuncRefs[i]
		fmt.Print(fmt.Sprintf("\nSearching for: %v\nDeclared in: %v, Line: %v", ref.Text[:len(ref.Text)-1], ref.FilePath, ref.LineNum)) // TODO: Need to display what file/line func is declared
		fmt.Print("\n-----------------------------------------------------\n")
		results := findFuncNameInFiles(ref.Text, dir)
		fmt.Print(fmt.Sprintf("%d Occurances\n", len(results)))
		fmt.Print(results)
	}

}

// Traverse dir path
func traverseDir(dirpath string) {
	err := filepath.Walk(dirpath, checkFile)
	if err != nil {
		fmt.Printf("traverseDir returned %v\n", err)
	}
}

func checkFile(path string, f os.FileInfo, err error) error {
	if f != nil {
        // OMIT DIRS and swp files(CAUSE EVERYONE USES VIM RIGHT????) TODO:Later we'll filter out thigns properly
		if !f.IsDir() && !strings.Contains(f.Name(), ".swp") { 
			openAndSearchFile(path)
		}
	}
	return nil
}

func openAndSearchFile(filepath string) {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	results := findFunksInStr(fileContent, filepath)
	FuncRefs = append(FuncRefs, results...) // TODO: THis global thing is gnarly
}

// Find all functions in string and return a slice of func names
func findFunksInStr(str []byte, sourceFilePath string) []FileRef {
	// TODO split str into lines
	lines := bytes.Split(str, []byte("\n"))
	matches := []FileRef{}
	if len(lines) != 0 {
		for i := range lines {
			line := string(lines[i])
			rx := "("
			for i := range FUNC_DECS {
				rx += fmt.Sprintf("\\b%v\\b \\w+\\s?\\(|", FUNC_DECS[i])
			}
			rx = rx[:len(rx)-1] + ")"
			re := regexp.MustCompile(rx)
			rawMatches := re.FindAllString(line, -1)
			for j := range rawMatches {
				matches = append(matches, FileRef{sourceFilePath, rawMatches[j], fmt.Sprintf("%d", i+1)})
			}
		}
	}

	return matches
}

// find a function name in a file at a line
func findFuncNameInFiles(name string, dirpath string) []FileRef {
	// clean the name of the function so we dont have to search with the function prefix
	re := regexp.MustCompile("(\\bfunction\\b|\\bfunc\\b|\\bdef\\b|\\s)") // TODO: Build this string out of FUNC_DECS
	cleanName := re.ReplaceAllString(name, "")
	return grepWrap(cleanName, dirpath)
}

func grepWrap(str string, dirpath string) []FileRef {
	results := []FileRef{}
	out, err := exec.Command("grep", "-rinwo", str, dirpath).Output()
	strB := []byte(str)
	if err != nil {
		//fmt.Println(err, "<--- error") // TODO: What should happen if grep errors out? Maybe only check for grep error code 2?
		return results
	}

	s := bytes.Split(out, []byte("\n"))
	for i := range s {
		ln := s[i]
		if bytes.Contains(ln, strB) {
			// parse the grep output. It looks like this ```/tmp/funcDEC.code:1:funcStyle1```
			parts := bytes.Split(ln, []byte(":"))
			if len(parts) > 1 {
				results = append(results, FileRef{string(parts[0]), string(parts[2]), string(parts[1])})
			}
		}
	}

	return results

}
