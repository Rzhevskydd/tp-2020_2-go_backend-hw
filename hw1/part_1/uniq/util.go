package uniq

func removeSpaces(str string) string {
	//TODO
	byteStr := []rune(str)
	keys := make(map[rune]bool)
	var res []rune

	for _, entry := range byteStr {
		if _, val := keys[entry]; !val {
			keys[entry] = true
			res = append(res, entry)
		}
	}
	return string(res)
}
