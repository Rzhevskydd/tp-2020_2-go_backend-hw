package uniq

import (
	"errors"
	"strings"
)

type OptionParams struct {
	C	bool
	D 	bool
	U 	bool
	F 	int
	S 	int
	I 	bool
}

func (op * OptionParams) IsValid() bool {
	incorrect1 := op.C && (op.D || op.U)
	incorrect2 := op.D && (op.C || op.U)
	incorrect3 := op.U && (op.C || op.D)
	if incorrect1 || incorrect2 || incorrect3 {
		return false
	}
	return true
}

func UniqUtility(text []string, params *OptionParams) ([]string, error) {
	if !params.IsValid() {
		return []string{}, errors.New("Incorrect using of flags, follow this:\n" +
			"uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
	}

	var counter = Counter{}

	order := 0
	var lineInfo LineInfo
	var exists bool

	for _, sourceLine := range text {
		line := sourceLine

		if params.F > 0 {
			splited := strings.Fields(sourceLine)

			if params.F < len(splited) {
				line = strings.Join(splited[params.F:], " ")
			}
		}

		if params.S > 0 && params.S < len(line) - 1{
			line = line[params.S:]
		}

		if params.I {
			line = strings.ToLower(line)
		}

		lineInfo, exists = counter[line]

		if !exists {
			lineInfo.Value = sourceLine
			lineInfo.Order = order
			order++
		}
		lineInfo.Count++

		counter[line] = lineInfo
	}

	var lines []string

	if params.C {
		lines = counter.GetCountUnique()
	} else if params.D {
		lines = counter.GetNotUnique()
	} else if params.U {
		lines = counter.GetNotRepeated()
	} else{
		lines = counter.GetUnique()
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = line
	}


	return result, nil
}