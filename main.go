package main

import (
	"bufio"
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

func main() {
	seenList := seenList{sync.RWMutex{}, make(map[string]filesize.Directory)}
	workList := make(chan []string)
	resultList := make(chan filesize.Directory)
	var wg sync.WaitGroup

	// TODO: check if only 1 args or refactor into accepting multiple args
	// TODO: user to set number of crawlers

	go func() {
		if len(os.Args[1:]) == 0 {
			workList <- []string{"."}
		} else {
			workList <- []string{os.Args[1:][0]}
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
					wg.Add(1)
					go filesize.DirCrawl(job, workList, resultList, &wg)
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

	// add a sort here using sort.interface
	// add some formatting

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("this is it")
	_, _ = reader.ReadByte()
}
