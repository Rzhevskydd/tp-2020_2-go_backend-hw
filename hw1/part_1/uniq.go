package main

import (
	"./uniq"
	"bufio"
	"flag"
	"os"
)
import "fmt"

func main() {
	cFlag := flag.Bool("c", false, "")
	dFlag := flag.Bool("d", false, "")
	uFlag := flag.Bool("u", false, "")

	numFields := flag.Int("f", 0, "")
	numChars := flag.Int("s", 0, "")

	iFlag := flag.Bool("i", false, "")

	flag.Parse()

	args := flag.Args()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if len(args) > 2 {
		panic("uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
	}

	var input, output = os.Stdin, os.Stdout
	var err error

	for i, arg := range args {
		switch i {
		case 0:
			defer input.Close()
			input, err = os.Open(arg)
			if err != nil {
				panic(err.Error())
			}
		case 1:
			defer output.Close()
			output, err = os.Create(arg)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	optional := &uniq.OptionParams{
		C:   *cFlag,
		D:   *dFlag,
		U:   *uFlag,
		F:   *numFields,
		S:   *numChars,
		I:   *iFlag,
	}

	in := bufio.NewScanner(input)
	var text []string
	for in.Scan() {
		text = append(text, in.Text())
	}

	result, err := uniq.UniqUtility(text, optional)
	if err != nil {
		fmt.Printf(err.Error())
	}

	out := bufio.NewWriter(output)
	defer out.Flush()

	for i, line := range result {
		if i != len(result) - 1 {
			_, err = fmt.Fprintln(out, line)
		} else {
			_, err = fmt.Fprint(out, line)
		}
		if err != nil {
			panic(err.Error())
		}
	}

}
