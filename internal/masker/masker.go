package masker

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// Masker add * mask on AES password for password hint in main TUI menu
func Masker(passwordPhrase string) string {

	var maskerString string
	pwdLen := len(passwordPhrase)
	switch pwdLen {
	case 1:
		maskerString = replaceAtIndex(passwordPhrase, '*', 0)
	case 2:
		maskerString = replaceAtIndex(passwordPhrase, '*', 1)
	case 3:
		maskerString = replaceAtIndex(passwordPhrase, '*', 1)
		maskerString = replaceAtIndex(maskerString, '*', 2)
	case 4:
		maskerString = replaceAtIndex(passwordPhrase, '*', 2)
		maskerString = replaceAtIndex(maskerString, '*', 3)
	default:
		for i := pwdLen / 2; i < pwdLen; i++ {
			maskerString = replaceAtIndex(passwordPhrase, '*', i)
			passwordPhrase = maskerString
		}
	}
	return maskerString
}
