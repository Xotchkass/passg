# passg

small CLI password generator just for fun and practice

## Install

`go install github.com/Xotchkass/passg@latest`

## Usage

```text
Usage: passg [--length LENGTH] [--number NUMBER] [--charactergroups CHARACTERGROUPS] [--include INCLUDE] [--exclude EXCLUDE] [--clipboard]

Options:
  --length LENGTH, -l LENGTH
                         Length of the password [default: 15]
  --number NUMBER, -n NUMBER
                         Number of passwords to generate [default: 1]
  --charactergroups CHARACTERGROUPS, -g CHARACTERGROUPS
                         Character groups to include: U/u=uppercase [A-Z], L/l=lowercase [a-z], D/d=digits [0-9], S/s=symbols [default: ULD]
  --include INCLUDE, -i INCLUDE
                         Additional characters to include in the password [default: -_!@$&/?\]
  --exclude EXCLUDE, -e EXCLUDE
                         Characters to exclude from the password
  --clipboard, -c        Copy password to clipboard instead of printing (ignored if -n > 1)
  --help, -h             display this help and exit
```
