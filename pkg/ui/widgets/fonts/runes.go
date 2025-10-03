package fonts

var (
	Numbers = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	Alpha   = make([]rune, 52)
)

func init() {
	for a := 'A'; a < 'Z'; a++ {
		Alpha[a-'A'] = a
	}
	for a := 'a'; a < 'z'; a++ {
		Alpha[a-'a'+26] = a
	}
}
