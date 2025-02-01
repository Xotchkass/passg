package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode"

	"github.com/alexflint/go-arg"
	"github.com/atotto/clipboard"
)

var args struct {
	Length          uint   `arg:"-l" default:"15" help:"Length of the password"`
	Number          int    `arg:"-n" default:"1" help:"Number of passwords to generate"`
	CharacterGroups string `arg:"-g" default:"ULD" help:"Character groups to include: U/u=uppercase [A-Z], L/l=lowercase [a-z], D/d=digits [0-9], S/s=symbols"`
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

func generatePassword(length uint, pool []byte) string {
	password := make([]byte, length)

	pool_i_buff := [8]byte{0}
	for i := uint(0); i < length; i++ {
		_, err := rand.Read(pool_i_buff[:])
		if err != nil {
			panic(err)
		}
		pool_i := binary.LittleEndian.Uint64(pool_i_buff[:]) % uint64(len(pool))
		password[i] = pool[pool_i]
	}
	return string(password)
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

	if args.Clipboard && args.Number == 1 {
		clipboard.WriteAll(generatePassword(args.Length, character_pool))
	} else {
		result := make([]string, 0, args.Number)
		for i := 0; i < args.Number; i++ {
			result = append(result, generatePassword(args.Length, character_pool))
		}
		fmt.Println(strings.Join(result, "\n"))
	}
}
