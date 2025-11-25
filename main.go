package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	cb "github.com/atotto/clipboard"
)

const (
	UPPER_LATIN = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LOWER_LATIN = "abcdefghijklmnopqrstuvwxyz"
	SYMBOLS     = "~!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
	DIGITS      = "0123456789"
)

func generatePassword(pool []byte, password []byte) {
	_, err := rand.Read(password)
	if err != nil {
		panic(err)
	}

	for i := range len(password) {
		password[i] = pool[int(password[i])%len(pool)]
	}

}

func main() {
	length := flag.Uint("l", 15, "Length of the password")
	number := flag.Int("n", 1, "Number of passwords to generate")
	characterGroups := flag.String("g", "ULD", "Character groups to include: U=uppercase [A-Z], L=lowercase [a-z], D=digits [0-9], S=symbols")
	include := flag.String("i", "-_!@$&/?\\", "Additional characters to include in the password")
	exclude := flag.String("e", "", "Characters to exclude from the password")
	clipboard := flag.Bool("c", false, "Copy password to clipboard instead of printing (ignored if -n > 1)")
	flag.Parse()

	characterPool := []byte{}
	for _, char := range *characterGroups {
		switch char {
		case 'U', 'u':
			characterPool = append(characterPool, UPPER_LATIN...)
		case 'L', 'l':
			characterPool = append(characterPool, LOWER_LATIN...)
		case 'D', 'd':
			characterPool = append(characterPool, DIGITS...)
		case 'S', 's':
			characterPool = append(characterPool, SYMBOLS...)
		default:
			fmt.Fprintf(os.Stderr, "Error: Wrong character group %c\n.", char)
			flag.Usage()
			os.Exit(1)
		}
	}
	for _, c := range *include {
		if c > unicode.MaxASCII {
			fmt.Fprintf(os.Stderr, "Error: non-ASCII characters not supported. Got '%c' in -i parameter.", c)
			os.Exit(1)
		}
		char := byte(c)
		if bytes.IndexByte(characterPool, char) == -1 {
			characterPool = append(characterPool, char)
		}
	}
	{
		finalPool := []byte{}
		for _, b := range characterPool {
			if !strings.ContainsRune(*exclude, rune(b)) {
				finalPool = append(finalPool, b)
			}
		}
		characterPool = finalPool
	}
	password := make([]byte, *length)

	if *clipboard && *number == 1 {
		generatePassword(characterPool, password)
		cb.WriteAll(string(password))
	} else {
		result := strings.Builder{}
		result.Grow(*number * (int(*length) + 1))
		for range *number {
			generatePassword(characterPool, password)
			_, err := result.Write(password)
			if err != nil {
				panic(err)
			}
			result.WriteByte('\n')
		}
		fmt.Print(result.String())
	}
}
