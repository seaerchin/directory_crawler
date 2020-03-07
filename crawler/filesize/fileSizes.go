package filesize

import (
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
	Info    os.FileInfo
	Name    string // path to this directory
	Largest string // name of the largest file in this directory
}

type DirMap map 

// crawls the directory rooted at s - s must be passed in as a FULLY QUALIFIED DIR NAME
func dirCrawl(s string) {

}
