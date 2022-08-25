package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
  "github.com/atotto/clipboard"
	"github.com/alexflint/go-arg"
)

var args struct {
  Length uint `arg:"-l" help:"Length of generated password" default:"15"`
  Number int `arg:"-n" help:"Number of generated passwords" default:"1"`
  CharacterGroups string `arg:"-g" default:"ULD" help:"Which group of characters include in password. Options: u - upper case latin letters [A-Z]. l - lower case latin letters [a-z]. d - digits [0-9]. s - symbols [~!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/]"`
  Include string `arg:"-i" help:"additional characters to include"`
  Exclude string `arg:"-e" help:"characters to exclude"`
  Clipboard bool `arg:"-c" help:"if set - writes generated password in clipboard instead of stdin. ignored if '-n' > 1"`
}

const (
    UPPER_LATIN = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    LOWER_LATIN = "abcdefghijklmnopqrstuvwxyz"
    SYMBOLS = "~!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
    DIGITS = "0123456789"
    )

func generate_password(length uint, pool []byte) string {
  pass_arr  := []byte{}
  
  for i := uint(0); i<length; i++ {
    pool_i, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
    if err != nil{
      panic("Unable to generate random number. For some reason.")
    }
    pass_arr = append(pass_arr, pool[pool_i.Uint64()])
  }
  return string(pass_arr)
}

func remove(s []byte, i int) []byte {
   return append(s[:i], s[i+1:]...)
}
    
func main() {
  arg.MustParse(&args)
 
  character_pool := []byte{}
  for _,char := range args.CharacterGroups{
    switch char{
      case 'U': 
        character_pool = append(character_pool, UPPER_LATIN...)
      case 'L': 
        character_pool = append(character_pool, LOWER_LATIN...)
      case 'D': 
        character_pool = append(character_pool, DIGITS...)
      case 'S': 
        character_pool = append(character_pool, SYMBOLS...)
      default: 
        panic(fmt.Errorf("Wrong character group %c", char))
    }
  }
  for _, c := range args.Include {
    char := byte(c)
    if !bytes.Contains(character_pool, []byte{char}){
      character_pool = append(character_pool, char)
    }
  }
  for _, c := range args.Exclude {
    char := byte(c)
    i := bytes.Index(character_pool, []byte{char}); if i >= 0{
      character_pool = remove(character_pool,i)
    }
  }
  if args.Clipboard && args.Number == 1{
    clipboard.WriteAll(generate_password(args.Length, character_pool))
  } else { 
    for i:=0; i<args.Number;i++ {
    fmt.Println(generate_password(args.Length, character_pool))
    }
  }
}
