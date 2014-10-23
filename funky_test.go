package main

import "testing"
import "io/ioutil"
import "os"

func TestFindFunksInStr(t *testing.T) {
	s := []byte("function hello() {} def hello(): func hello() {}")
	if len(findFunksInStr(s, "test")) != 3 {
		t.Error("Fail")
	}
}

func TestFindInFile(t *testing.T) {
	generateFiles(t)
	defer cleanup()
	findFuncNameInFiles("function funcStyle1", "/tmp/testdir")
}

func TestTraversal(t *testing.T) {
	generateFiles(t)
	traverseDir("/tmp/testdir") // TODO: This is incomplete
    // we cant defer clean up cause of the go routines spawned in traverseDir. cleanup() will get called before anything is done
}

func TestOpenAndReadFile(t *testing.T) {
	generateFiles(t)
	defer cleanup()
	openAndSearchFile("/tmp/testdir/funcDEC.code")
	if len(FuncRefs) < 3 {
        t.Error("Fail")
    }
}

// generate some files for testing

func generateFiles(t *testing.T) {
	// create file that declares functions
	os.MkdirAll("/tmp/testdir", 0777)
	text := []byte("function funcStyle1() {}\ndef func_style2():\nfunc funcStyle3() {}")
	err := ioutil.WriteFile("/tmp/testdir/funcDEC.code", text, 0644)
	if err != nil {
		t.Error(err)
	}
	// create a file that calls functions
	text = []byte("funcStyle1(); func_style2() funcStyle3()")
	err = ioutil.WriteFile("/tmp/testdir/funcUSE.code", text, 0644)
	if err != nil {
		t.Error(err)
	}
}

// Clean up test garbage files
func cleanup() {
	os.RemoveAll("/tmp/testdir")
}
