package filesize

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// this package aims to get the file sizes

// File describes a file on the system
type File struct {
	Parent string // parent directory which this file is rooted in
	os.FileInfo
}

// Directory describes a directory on the system
type Directory struct {
	subDirs []string // sub-directories
	Name    string   // path to this directory
	Largest string   // name of the largest file in this directory
	Size    int64    // size of this directory
}

type DirMap struct {
}

// crawls the directory rooted at s - s must be passed in as a FULLY QUALIFIED DIR NAME
// returns d and its associated information
func dirCrawl(s string) (d Directory, err error) {
	f, err := os.Open(s)
	if err != nil {
		return d, errors.New("Opening error:" + f.Name()) // need to return relevant information for d
	}
	dir, err := f.Readdir(-1)
	if err != nil {
		return d, err // don't return partially read dir, let caller decide what to do
	}
	// now i have a flat slice describing my directory
	var max int64
	var largest string
	var cumSum int64
	subDirs := make([]string, 0)
	for _, file := range dir {
		info, err := os.Lstat(filepath.Join(s, file.Name()))
		if err != nil {
			return d, errors.New("we can't find yo file: :" + filepath.Join(s, info.Name()))
		}
		if info.IsDir() {
			subName := filepath.Join(s, info.Name())
			subDirs = append(subDirs, subName)
		}
		temp := info.Size()
		if temp > max && file.Name() != f.Name() {
			max = temp
			largest = info.Name()
		}
		cumSum += temp
	}
	return Directory{subDirs: subDirs, Name: f.Name(), Largest: largest, Size: cumSum}, nil
}

// DirCrawl is a handler for the dirCrawl method - this handles the recursive calls
func DirCrawl(s string, workList chan<- []string, resultList chan<- Directory, wg *sync.WaitGroup) {
	result, err := dirCrawl(s)
	defer fmt.Println("done")
	defer wg.Done()
	if err != nil {
		fmt.Println(err)
		return
	}

	toSend := make([]string, 0)
	for _, v := range result.subDirs {
		fmt.Println("sending: " + v)
		toSend = append(toSend, v)
	}
	workList <- toSend
	resultList <- result
}
