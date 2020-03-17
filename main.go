package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/disiqueira/gotree"
	"github.com/seaerchin/directory_crawler/filesize"
)

// TODO: impl verbose
// TODO: add tests

var numCrawler = flag.Int("n", 20, "set the maximum number of crawlers here!")
var verbose = flag.Bool("v", false, "set to true for constant updates on the progress")
var force = flag.Bool("f", false, "set to true to allow deletion of the largest files in ALL subdirectories")

func main() {
	flag.Parse()
	// defensive check
	if *numCrawler < 1 {
		*numCrawler = 1
	}

	var wg sync.WaitGroup
	var forest []gotree.Tree

	if len(flag.Args()) > 0 {
		for _, job := range flag.Args() {
			root := gotree.New(fmt.Sprintf("%s - files: %d", job, filesize.GetSize(job)))
			r := filesize.Request{root, job, make(chan interface{})}
			forest = append(forest, r.Root)
			wg.Add(1)
			go filesize.DirCrawl(r, *force, &wg)
		}
	} else {
		// repeated code but wtvr
		root := gotree.New(fmt.Sprintf("%s - files: %d", ".", filesize.GetSize(".")))
		r := filesize.Request{root, ".", make(chan interface{})}
		forest = append(forest, r.Root)
		wg.Add(1)
		go filesize.DirCrawl(r, *force, &wg)
	}

	// no guarantees of wait unless sleep - find a way to do some work in background
	wg.Wait()

	for _, tree := range forest {
		fmt.Println("here's your tree")
		fmt.Println(tree.Print())
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("this is it")
	_, _ = reader.ReadByte()
}
