package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"unicode"

	"github.com/alexflint/go-arg"
	"github.com/atotto/clipboard"
)

var args struct {
	Length          uint   `arg:"-l" default:"15" help:"Length of the password"`
	Number          int    `arg:"-n" default:"1" help:"Number of passwords to generate"`
	CharacterGroups string `arg:"-g" default:"ULD" help:"Character groups to include: U=uppercase [A-Z], L=lowercase [a-z], D=digits [0-9], S=symbols"`
	Include         string `arg:"-i" help:"Additional characters to include in the password" default:"-_!@$&/?\\"`
	Exclude         string `arg:"-e" help:"Characters to exclude from the password"`
	Clipboard       bool   `arg:"-c" help:"Copy password to clipboard instead of printing (ignored if -n > 1)"`
}

const (
	UPPER_LATIN = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LOWER_LATIN = "abcdefghijklmnopqrstuvwxyz"
	SYMBOLS     = "~!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
	DIGITS      = "0123456789"
)

func generatePassword(pool []byte, password []byte){
	_, err := rand.Read(password)
	if err != nil {
		panic(err)
	}

	for i:= range len(password) {
		password[i] = pool[int(password[i]) % len(pool)]
	}

}

func remove(s []byte, i int) []byte {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func main() {
	arg.MustParse(&args)

	character_pool := []byte{}
	for _, char := range args.CharacterGroups {
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
			panic(fmt.Errorf("wrong character group %c", char))
		}
	}
	for _, c := range args.Include {
		if c > unicode.MaxASCII {
			panic(fmt.Errorf("non-ASCII characters not supported. Got '%c'", c))
		}
		char := byte(c)
		if bytes.IndexByte(character_pool, char) == -1 {
			character_pool = append(character_pool, char)
		}
	}
	for _, c := range args.Exclude {
		if c > unicode.MaxASCII {
			continue
		}
		char := byte(c)
		i := bytes.IndexByte(character_pool, char)
		if i >= 0 {
			character_pool = remove(character_pool, i)
		}
	}
	password := make([]byte, args.Length)

	if args.Clipboard && args.Number == 1 {
		generatePassword(character_pool, password)
		clipboard.WriteAll(string(password))
	} else {
		result := strings.Builder{}
		result.Grow(args.Number * (int(args.Length) + 1))
		for range args.Number {
			generatePassword(character_pool, password)
			_, err := result.Write(password)
			if err != nil {
				panic(err)
			}
			err = result.WriteByte(byte('\n'))
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(result.String())
	}
}
