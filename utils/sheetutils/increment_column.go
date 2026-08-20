package sheetutils

func IncrementColumn(currentColumn string) string {
	if len(currentColumn) == 1 {
		currentChar := currentColumn[0]
		if currentChar == 'Z' {
			return "AA"
		}
		return string(currentChar + 1)
	}
	lastChar := currentColumn[len(currentColumn)-1]
	if lastChar == 'Z' {
		return IncrementColumn(currentColumn[:len(currentColumn)-1]) + "A"
	}
	return currentColumn[:len(currentColumn)-1] + string(lastChar+1)
}

func IncrementColumnN(currentColumn string, n int) string {
	for range n {
		currentColumn = IncrementColumn(currentColumn)
	}
	return currentColumn
}
