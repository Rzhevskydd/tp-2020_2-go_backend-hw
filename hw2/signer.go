package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type channels map[string]chan interface{}

type ch chan interface{}

func ExecutePipeline(tasks ...job) {

	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))

	worker := func(waiter *sync.WaitGroup, task job, in ch) ch {
		out := make(ch, 1)
		go func() {
			defer func() {
				close(out)
				waiter.Done()
			}()
			task(in, out)
		}()
		return out
	}

	c := make(ch, 1)
	for i, task := range tasks {
		if i == 0 {
			close(c)
		}
		c = worker(wg, task, c)
	}
	wg.Wait()
}

func produceCrc32(out chan string, in chan string) {
	out <- DataSignerCrc32(<-in)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for input := range in {
		dataStr := fmt.Sprintf("%v", input)

		wg.Add(1)
		go func(out ch, waiter *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()

			firstCrc32 := make(chan string, 1)
			md5 := make(chan string, 1)
			secondCrc32 := make(chan string, 1)


			in := make(chan string, 1)
			in <- dataStr

			go produceCrc32(firstCrc32, in)
			go func(c chan string) {
				mu.Lock()
				c <- DataSignerMd5(dataStr)
				mu.Unlock()
			}(md5)
			go produceCrc32(secondCrc32, md5)

			out <- <-firstCrc32 + "~" + <-secondCrc32
		}(out, wg, mu)
	}
	wg.Wait()
}

type inter struct {
	i   int
	str string
}

func produceMultiHash(out ch, waiter *sync.WaitGroup, dataStr string) {
	defer waiter.Done()

	c := make(chan inter)
	for i := 0; i < 6; i++ {
		go func(i int, c chan inter) {
			c <- inter{i, DataSignerCrc32(strconv.FormatInt(int64(i), 10) + dataStr)}
		}(i, c)
	}

	hashesToConcat := make([]string, 6, 6)
	for i := 0; i < 6; i++ {
		chunk := <-c
		hashesToConcat[chunk.i] = chunk.str
	}
	concated := strings.Join(hashesToConcat, "")

	out <- concated
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for input := range in {
		switch dataStr := input.(type) {
		case string:
			wg.Add(1)
			go produceMultiHash(out, wg, dataStr)
		}
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var lst []string
	for input := range in {
		hash := fmt.Sprintf("%v", input)
		lst = append(lst, hash)
	}
	sort.Strings(lst)

	res := strings.Join(lst, "_")
	out <- res
}
