package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func parseArgs(s string, left string, i int) (float64, float64, int, error) {
	var next  = len(s[i + 1:])
	var right string
	for j, char := range s[i + 1:] {
		next = j + i + 1
		symb := fmt.Sprintf("%c", char)
		if !(symb >= "0" && symb <= "9") {
			next--
			break
		}
		right += fmt.Sprintf("%c", char)
	}
	next++

	num1, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	num2, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return num1, num2, next, nil
}

func evaluateExp(s string) (res float64) {
	var left string
	for i, c := range s {
		switch c {
		case '+':
			num, err := strconv.ParseFloat(left, 64)
			if err != nil && left != "" {
				panic(err.Error())
			}
			res = num + evaluateExp(s[i + 1:])
			return res
		case '-':
			num, err := strconv.ParseFloat(left, 64)
			if err != nil {
				panic(err.Error())
			}
			res = num - evaluateExp(s[i + 1:])
			return res
		case '*':
			num1, num2, next, err := parseArgs(s, left, i)
			if err != nil {
				panic(err.Error())
			}

			res = num1 * num2
			s = strings.Replace(s, s[:next], fmt.Sprintf("%f", res), 1)

			return evaluateExp(s)

		case '/':
			num1, num2, next, err := parseArgs(s, left, i)
			if err != nil {
				panic(err.Error())
			}

			res = num1 / num2
			s = strings.Replace(s, s[:next], fmt.Sprintf("%f", res), 1)

			return evaluateExp(s)

		case '(':
			var closeBrace int
			for j, char := range s {
				symb := fmt.Sprintf("%c", char)
				if symb == ")" {
					closeBrace = j
					break
				}
			}

			tmp := evaluateExp(s[i+1 : closeBrace])

			s = strings.Replace(s, s[:closeBrace+1], fmt.Sprintf("%f", tmp), 1)
			return evaluateExp(s)

		default:
			left += fmt.Sprintf("%c", c)
		}
	}

	res, err := strconv.ParseFloat(left, 64)
	if err != nil {
		panic(err.Error())
	}
	return res
}

func main() {
	var expression string

	if !(len(os.Args) == 2) {
		byteStr, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err.Error())
		}
		expression = strings.TrimSpace(string(byteStr))
	} else {
		expression = strings.TrimSpace(os.Args[1])
	}

	res := evaluateExp(expression)
	fmt.Println(res)
}
