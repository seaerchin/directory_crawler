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

func main() {
	seenList := make(map[string]filesize.Directory)
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
			seenList[result.Name] = result
		}
	}()

	go func() {
		for i := range workList {
			for _, job := range i {
				if _, ok := seenList[job]; !ok {
					wg.Add(1)
					go filesize.DirCrawl(job, workList, resultList, &wg)
				}
			}
		}
	}()

	time.Sleep(1 * time.Second)

	wg.Wait()
	close(workList)
	close(resultList)

	for key, value := range seenList {
		fmt.Println("key: " + key + " size: " + strconv.FormatInt(value.Size, 10))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("this is it")
	_, _ = reader.ReadByte()
}
