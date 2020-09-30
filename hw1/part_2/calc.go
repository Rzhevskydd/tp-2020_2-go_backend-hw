package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


func parseBraces(startBrace int, s string) (float64, int) {
	var closeBrace int
	var nestedBracesCount int
	for j, char := range s[startBrace + 1:] {
		symb := fmt.Sprintf("%c", char)
		if symb == ")" {
			if !(nestedBracesCount > 0) {
				closeBrace = j
				break
			}
			nestedBracesCount--

		} else if symb == "(" {
			nestedBracesCount++
		}
	}

	cut := s[startBrace + 1 : startBrace + closeBrace + 1]
	res := evaluateExp(cut)
	return res, startBrace + closeBrace + 1
}

func extractRightOperand(s string, left string, cutFrom int) (float64, float64, int, error) {
	var idxToReplace = len(s[cutFrom+ 1:])
	var right string

	cutStr := s[cutFrom+ 1:]
	for i, char := range cutStr {
		if char == '(' {
			replace, closeBrace := parseBraces(cutFrom + 1, s);
			s = strings.Replace(s, s[cutFrom + 1 : closeBrace + 1], fmt.Sprintf("%.4f", replace), 1)
			right = s[cutFrom+1:]
			idxToReplace = closeBrace
			break
		}

		idxToReplace = i + cutFrom + 1
		symb := fmt.Sprintf("%c", char)
		if !(symb >= "0" && symb <= "9") {
			idxToReplace--
			break
		}
		right += fmt.Sprintf("%c", char)
	}
	idxToReplace++

	num1, err := strconv.ParseFloat(left, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	num2, err := strconv.ParseFloat(right, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return num1, num2, idxToReplace, nil
}

func evaluateExp(s string) (res float64) {
	var leftOperand string
	for i, c := range s {
		switch c {
		case '+':
			num, err := strconv.ParseFloat(leftOperand, 64)
			if err != nil && leftOperand != "" {
				panic(err.Error())
			}
			res = num + evaluateExp(s[i + 1:])
			return res
		case '-':
			// negative number
			if leftOperand == "" {
				leftOperand += fmt.Sprintf("%c", c)
				continue
			}

			//num, err := strconv.ParseFloat(leftOperand, 64)
			num1, num2, idxToReplace, err := extractRightOperand(s, leftOperand, i)
			if err != nil {
				panic(err.Error())
			}
			//res = num - evaluateExp(s[i + 1:])
			res = num1 - num2
			s = strings.Replace(s, s[:idxToReplace], fmt.Sprintf("%.4f", res), 1)
			return evaluateExp(s)
		case '*':
			num1, num2, idxToReplace, err := extractRightOperand(s, leftOperand, i)
			if err != nil {
				panic(err.Error())
			}

			res = num1 * num2
			s = strings.Replace(s, s[:idxToReplace], fmt.Sprintf("%.4f", res), 1)

			return evaluateExp(s)

		case '/':
			num1, num2, idxToReplace, err := extractRightOperand(s, leftOperand, i)
			if err != nil {
				panic(err.Error())
			}

			res = num1 / num2
			s = strings.Replace(s, s[:idxToReplace], fmt.Sprintf("%.4f", res), 1)

			return evaluateExp(s)

		case '(':
			replace, closeBrace := parseBraces(i, s);

			s = strings.Replace(s, s[:closeBrace+1], fmt.Sprintf("%.4f", replace), 1)
			return evaluateExp(s)

		default:
			leftOperand += fmt.Sprintf("%c", c)
		}
	}

	res, err := strconv.ParseFloat(leftOperand, 64)
	if err != nil {
		panic(err.Error())
	}
	return res
}

func main() {
	var expression string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Please, enter correct expression\n" +
				"calc.go supports '*', '/', '+', '-' and '(...) operators")
		}
	}()

	if !(len(os.Args) == 2) {
		byteStr, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err.Error())
		}
		expression = strings.TrimSpace(string(byteStr))
	} else {
		expression = strings.TrimSpace(os.Args[1])
	}

	// remove all spaces
	expression = strings.Replace(expression, " ", "", -1)
	
	res := evaluateExp(expression)
	fmt.Println(res)
}
