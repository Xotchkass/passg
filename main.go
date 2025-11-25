package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
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

func remove(s []byte, i int) []byte {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func main() {
	Length := flag.Uint("l", 15, "Length of the password")
	Number := flag.Int("n", 1, "Number of passwords to generate")
	CharacterGroups := flag.String("g", "ULD", "Character groups to include: U=uppercase [A-Z], L=lowercase [a-z], D=digits [0-9], S=symbols")
	Include := flag.String("i", "-_!@$&/?\\", "Additional characters to include in the password")
	Exclude := flag.String("e", "", "Characters to exclude from the password")
	Clipboard := flag.Bool("c", false, "Copy password to clipboard instead of printing (ignored if -n > 1)")
	flag.Parse()

	character_pool := []byte{}
	for _, char := range *CharacterGroups {
		switch char {
		case 'U', 'u':
			character_pool = append(character_pool, UPPER_LATIN...)
		case 'L', 'l':
			character_pool = append(character_pool, LOWER_LATIN...)
		case 'D', 'd':
			character_pool = append(character_pool, DIGITS...)
		case 'S', 's':
			character_pool = append(character_pool, SYMBOLS...)
		default:
			fmt.Fprintf(os.Stderr, "Error: Wrong character group %c\n.", char)
			flag.Usage()
			os.Exit(1)
		}
	}
	for _, c := range *Include {
		if c > unicode.MaxASCII {
			fmt.Fprintf(os.Stderr, "Error: non-ASCII characters not supported. Got '%c' in -i parameter.", c)
			os.Exit(1)
		}
		char := byte(c)
		if bytes.IndexByte(character_pool, char) == -1 {
			character_pool = append(character_pool, char)
		}
	}
	for _, c := range *Exclude {
		if c > unicode.MaxASCII {
			continue
		}
		char := byte(c)
		i := bytes.IndexByte(character_pool, char)
		if i >= 0 {
			character_pool = remove(character_pool, i)
		}
	}
	password := make([]byte, *Length)

	if *Clipboard && *Number == 1 {
		generatePassword(character_pool, password)
		clipboard.WriteAll(string(password))
	} else {
		result := strings.Builder{}
		result.Grow(*Number * (int(*Length) + 1))
		for range *Number {
			generatePassword(character_pool, password)
			_, err := result.Write(password)
			if err != nil {
				panic(err)
			}
			result.WriteByte('\n')
		}
		fmt.Print(result.String())
	}
}
