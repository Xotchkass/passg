package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
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

func generate_password(length uint, pool []byte) string {
	pass_arr := make([]byte, length)

	for i := uint(0); i < length; i++ {
		pool_i, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		if err != nil {
			panic(err)
		}
		pass_arr[i] = pool[pool_i.Uint64()]
	}
	return string(pass_arr)
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
		clipboard.WriteAll(generate_password(args.Length, character_pool))
	} else {
		result := make([]string, 0, args.Number)
		for i := 0; i < args.Number; i++ {
			result = append(result, generate_password(args.Length, character_pool))
		}
		fmt.Println(strings.Join(result, "\n"))
	}
}
