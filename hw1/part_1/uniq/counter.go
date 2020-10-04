package uniq

import "strconv"

type Counter map[string]LineInfo

func (c *Counter) GetUnique() []string {
	result := make([]string, len(*c))
	for _, lineInfo := range *c {
		result[lineInfo.Order] = lineInfo.Value
	}
	return result
}

func (c *Counter) GetCountUnique() []string {
	result := make([]string, len(*c))
	for _, lineInfo := range *c {
		result[lineInfo.Order] = strconv.FormatInt(int64(lineInfo.Count), 10) + " " + lineInfo.Value
	}
	return result
}

func (c *Counter) GetNotUnique() []string {
	var notUniq []LineInfo
	for _, lineInfo := range *c {
		if lineInfo.Count > 1 {
			notUniq = append(notUniq, lineInfo)
		}
	}

	By(order).Sort(notUniq)

	result := make([]string, len(notUniq))
	for i, line := range notUniq {
		result[i] = line.Value
	}

	return result
}

func (c *Counter) GetNotRepeated() []string {
	var notRepeated []LineInfo
	for _, lineInfo := range *c {
		if lineInfo.Count == 1 {
			notRepeated = append(notRepeated, lineInfo)
		}
	}

	By(order).Sort(notRepeated)

	result := make([]string, len(notRepeated))
	for i, line := range notRepeated {
		result[i] = line.Value
	}

	return result
}
