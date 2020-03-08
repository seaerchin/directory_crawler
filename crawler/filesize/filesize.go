package filesize

import (
	"fmt"
	"os"
)

// this package aims to get the file sizes

// File describes a file on the system
type File struct {
	Parent string // parent directory which this file is rooted in
	os.FileInfo
}

// Directory describes a directory on the system
type Directory struct {
	Info    os.FileInfo // fileinfo of this directory
	Name    string      // path to this directory
	Largest string      // name of the largest file in this directory
	Size    int64       // size of this directory
}

type DirMap struct {
}

// crawls the directory rooted at s - s must be passed in as a FULLY QUALIFIED DIR NAME
// returns d and its associated information
func dirCrawl(s string) (d Directory, err error) {
	f, err := os.Open(s)
	if err != nil {
		return d, err // need to return relevant information for d
	}
	dir, err := f.Readdir(-1)
	if err != nil {
		return d, err // don't return partially read dir, let caller decide what to do
	}
	// now i have a flat slice describing my directory\
	var max int64
	var largest string
	var cumSum int64
	for _, file := range dir {
		info, _ := os.Lstat(file.Name())
		temp := info.Size()
		if temp > max && file.Name() != f.Name() {
			max = temp
			largest = info.Name()
		}
		cumSum += temp
	}
	// largest filesize known
	this, err := os.Lstat(f.Name())
	if err != nil {
		return d, err
	}
	return Directory{Info: this, Name: f.Name(), Largest: largest, Size: cumSum}, nil
}

func DirCrawl(s string) {
	res, _ := dirCrawl(s)
	fmt.Printf("%+v", res)
}
