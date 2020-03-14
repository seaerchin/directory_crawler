package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/seaerchin/directory_crawler/filesize"
)

type seenList struct {
	sync.RWMutex
	cache map[string]filesize.Directory
}

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

	seenList := seenList{sync.RWMutex{}, make(map[string]filesize.Directory)}
	workList := make(chan []string)
	resultList := make(chan filesize.Directory)
	numTokens := make(chan struct{}, *numCrawler)

	var wg sync.WaitGroup

	fmt.Println(os.Args)
	go func() {
		if len(flag.Args()) == 0 {
			workList <- []string{"."}
		} else {
			workList <- flag.Args()
		}
	}()

	go func() {
		for result := range resultList {
			seenList.Lock()
			seenList.cache[result.Name] = result
			seenList.Unlock()
		}
	}()

	go func() {
		for jobList := range workList {
			for _, job := range jobList {
				seenList.RLock()
				_, ok := seenList.cache[job]
				seenList.RUnlock()
				if !ok {
					numTokens <- struct{}{} // acquire token
					wg.Add(1)
					go filesize.DirCrawl(job, workList, resultList, *force, &wg)
					<-numTokens // release token
				}
			}
		}
	}()

	// no guarantees of wait unless sleep - find a way to do some work in background
	time.Sleep(1 * time.Second)
	wg.Wait()
	close(workList)
	close(resultList)

	for key, value := range seenList.cache {
		fmt.Println("key: " + key + " size: " + strconv.FormatInt(value.Size, 10) + " largest file: " + value.Largest)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("this is it")
	_, _ = reader.ReadByte()
}
