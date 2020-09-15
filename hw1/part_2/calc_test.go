package main

import (
	"testing"
)

func TestCalc(t *testing.T) {
	var testCases = []struct {
		input string
		result float64
	}{
		{
			input: "5+10",
			result: 15,
		},
		{
			input: "2+2*2",
			result: 6,
		},
		{
			input: "1000+2*(25+100/4*(2*(2+2)))",
			result: 1450,
		},
		{
			input: "(1000+5)",
			result: 1005,
		},
		{
			input: "(1+2)-3",
			result: 0,
		},
		{
			input: "(1+2)*3",
			result: 9,
		},
		{
			input: "100/4/5",
			result: 5,
		},
	}

	for num, test := range testCases {
		res := evaluateExp(test.input)
		if res != test.result {
			t.Errorf("%.2f != %.2f\n Test number: %d", res, test.result, num)
		}
	}
}