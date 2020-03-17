package filesize

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/disiqueira/gotree"
)

// this package aims to get the file sizes

// File describes a file on the system
type File struct {
	Parent string // parent directory which this file is rooted in
	os.FileInfo
}

// Directory describes a directory on the system
type directory struct {
	subDirs []string // sub-directories - deprecate  but kept
	name    string   // path to this directory
	largest string   // name of the largest file in this directory
	size    int64    // size of files in this directory
}

// crawls the directory rooted at s - s must be passed in as a FULLY QUALIFIED DIR NAME
// returns d and its associated information
func dirCrawl(s string, canDelete bool) (d directory, err error) {
	f, err := os.Open(s)
	if err != nil {
		return d, errors.New("Opening error: " + s) // need to return relevant information for d
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
	if canDelete {
		os.Remove(filepath.Join(s, largest))
	}

	return directory{subDirs: subDirs, name: f.Name(), largest: largest, size: cumSum}, nil
}

// DirCrawl is a handler for the dirCrawl method - this handles the recursive calls
// function will walk the directory rooted at s using filepath.walk
func DirCrawl(r Request, canDelete bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(r.Result)
	// d, err := dirCrawl(r.Job, canDelete) // double work done here

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// crawler gets called for every file in job
	crawler := func(paths string, f os.FileInfo, err error) error {
		// qualified directory
		// fmt.Println(paths)
		if f.IsDir() && paths != r.Job {
			info, err := dirCrawl(paths, canDelete)
			if err != nil {
				return err
			}
			stringRep := fmt.Sprintf("%s - files: %d", f.Name(), info.size)
			root := r.Root.Add(stringRep)
			newRequest := Request{root, paths, make(chan interface{})}
			wg.Add(1)
			go DirCrawl(newRequest, canDelete, wg) // refactor to use requests so sub directories return a request; calling function appends subtrees to parent
			return filepath.SkipDir
		}
		return nil
	}

	// close(r.SubDirs) // ordering maybe wrong
	// for range d.subDirs {
	// 	root.AddTree((<-r.SubDirs).(gotree.Tree))
	// }
	// r.Result <- r.Root
	filepath.Walk(r.Job, filepath.WalkFunc(crawler))
}

// Request is a request to traverse the directory rooted at s; results are returned through the result channel
type Request struct {
	Root   gotree.Tree
	Job    string
	Result chan interface{} // communicate subroutine done
}

func GetSize(s string) int64 {
	info, err := dirCrawl(s, false)
	if err != nil {
		panic("dude call this on a qualified dir man, wtf")
	}
	return info.size
}
