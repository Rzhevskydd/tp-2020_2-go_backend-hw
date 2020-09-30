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
	for _, task := range tasks {
		c = worker(wg, task, c)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for input := range in {
		dataStr := fmt.Sprintf("%v", input)

		wg.Add(1)
		go func(out ch, waiter *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			mu.Lock()
			md5 := DataSignerMd5(dataStr)
			mu.Unlock()

			c_md5 := make(chan string)

			go func() {
				c_md5 <- DataSignerCrc32(md5)
			}()

			out <- DataSignerCrc32(dataStr) +
				"~" + <-c_md5
		}(out, wg, mu)
	}
	wg.Wait()
}

type inter struct {
	i   int
	str string
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for input := range in {
		switch dataStr := input.(type) {
		case string:
			wg.Add(1)
			go func(out ch, waiter *sync.WaitGroup) {
				defer wg.Done()

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
			}(out, wg)
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
	var res string
	for idx, el := range lst {
		res += el
		if idx != len(lst)-1 {
			res += "_"
		}
	}
	out <- res
}
