package uniq

import "sort"

type LineInfo struct {
	Value string
	Count int
	Order int
}

var order = func(l1, l2 * LineInfo) bool {
	return l1.Order < l2.Order
}
type By func(u1, u2 * LineInfo) bool

func (by By) Sort(lines []LineInfo) {
	ls := &LineSorter{
		lines: lines,
		by: by,
	}
	sort.Sort(ls)
}

type LineSorter struct {
	lines []LineInfo
	by func(l1, l2 * LineInfo) bool
}

func (ls * LineSorter) Len() int {
	return len(ls.lines)
}

func (ls * LineSorter) Swap(i, j int) {
	ls.lines[i], ls.lines[j] = ls.lines[j], ls.lines[i]
}

func (ls * LineSorter) Less(i, j int) bool {
	return ls.by(&ls.lines[i], &ls.lines[j])
}