# passgen-go
just another CLI password generator

## Install
`go install github.com/Xotchkass/passgen-go@latest`

## Usage 
```
Usage: main [--length LENGTH] [--number NUMBER] [--charactergroups CHARACTERGROUPS] [--include INCLUDE] [--exclude EXCLUDE] [--clipboard]

Options:
  --length LENGTH, -l LENGTH
                         Length of generated password [default: 15]
  --number NUMBER, -n NUMBER
                         Number of generated passwords [default: 1]
  --charactergroups CHARACTERGROUPS, -g CHARACTERGROUPS
                         Which group of characters include in password. Options: u - upper case latin letters [A-Z]. l - lower case latin letters [a-z]. d - digits [0-9]. s - symbols [~!@#$%^&*()_-+={[}]|\:;"'<,>.?/] [default: ULD]
  --include INCLUDE, -i INCLUDE
                         additional characters to include
  --exclude EXCLUDE, -e EXCLUDE
                         characters to exclude
  --clipboard, -c        if set - writes generated password in clipboard instead of stdin. ignored if '-n' > 1
  --help, -h             display this help and exit
```
