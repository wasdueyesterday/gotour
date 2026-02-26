package gotour


func LongestNorepeatSubStr(s string) int {
	l := 0
	longest := 0
	sett := make(map[byte]struct{})

	// r is right index, relative to l
	// when valid, move r to expand the window
	// when invalid, move l to shrink the window
	for r := 0; r < len(s); r++ {
		for {
			// v is rune, the unicode point, since we're dealing only with ASCII, it's a byte (uint8)
			if _, exist := sett[s[r]]; !exist {
				break
			}
			delete(sett, s[l])
			l++
		}
		// calc the window len btw l, r pointers
		var w int = r - l + 1
		longest = max(longest, w)
		// in go, hash map mimics set, struct{}{} takes no memory
		sett[s[r]] = struct{}{}
	}
	return longest
}