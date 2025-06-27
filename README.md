<div align="center">
    <h1>getopt_long.go</h1>
    <h4>Go option parser inspired by getopt_long(3)</h4>
    <p>
        <a href="https://github.com/BChristieDev/getopt_long.go/actions/workflows/ci.yml"><img src="https://github.com/BChristieDev/getopt_long.go/actions/workflows/ci.yml/badge.svg"></a>
        <a href="https://pkg.go.dev/github.com/BChristieDev/getopt_long.go"><img src="https://pkg.go.dev/badge/github.com/BChristieDev/getopt_long.go.svg"></a>
    </p>
    <p>
        <a href="#install">Install</a> •
        <a href="#examples">Examples</a> •
        <a href="#maintainers">Maintainers</a> •
        <a href="#contributing">Contributing</a> •
        <a href="#license">License</a>
    </p>
</div>

## Install

```sh
$ go get github.com/BChristieDev/getopt_long.go
```

## Examples

### Simple

```go
package main

import (
	"fmt"
	"os"

	. "github.com/BChristieDev/getopt_long.go/pkg/getoptlong"
)

func main() {
	longopts := []Option{
		{Name: "foo", HasArg: RequiredArgument, Flag: nil, Val: 0},
	}

	var longindex, opt int

	for {
		opt = GetoptLong(len(os.Args), os.Args, "a:", longopts, &longindex)

		if opt == -1 {
			break
		}

		switch opt {
		case 0:
			fmt.Printf("option '%s' has argument '%s'", longopts[longindex].Name, OptArg)
		case 'a':
			fmt.Printf("option '%c' has argument '%s'", opt, OptArg)
		}
	}
}
```

```sh
$ ./simple --foo bar --foo=baz --foo
option 'foo' has argument 'bar'
option 'foo' has argument 'baz'
option '--foo' requires an argument

$ ./simple -a foo -abar -a
option 'a' has argument 'foo'
option 'a' has argument 'bar'
option requires an argument -- a
```

### Complex

```go
package main

import (
	"fmt"
	"os"

	. "github.com/BChristieDev/getopt_long.go/pkg/getoptlong"
)

func main() {
	frob_state := struct {
		unset int
		off   int
		on    int
	}{
		-1,
		0,
		1,
	}

	frob_flag := frob_state.unset

	longopts := []Option{
		{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 'a'},
		{Name: "bar", HasArg: OptionalArgument, Flag: nil, Val: 'b'},
		{Name: "baz", HasArg: RequiredArgument, Flag: nil, Val: 'c'},
		{Name: "on", HasArg: NoArgument, Flag: &frob_flag, Val: frob_state.on},
		{Name: "off", HasArg: NoArgument, Flag: &frob_flag, Val: frob_state.off},
		{Name: "silent", HasArg: NoArgument, Flag: nil, Val: 's'},
	}

	var longindex, opt int

	for {
		opt = GetoptLong(len(os.Args), os.Args, "ab::c:", longopts, &longindex)

		if opt == -1 {
			break
		}

		switch opt {
		case 0:
			fmt.Printf("option '%s' changed frob state to '%d'\n", longopts[longindex].Name, frob_flag)
		case 'a':
			fallthrough
		case 'b':
			fallthrough
		case 'c':
			fmt.Printf("option '%c' has argument '%s'\n", opt, OptArg)
		case 's':
			OptErr = 0
		}
	}

	if OptInd < len(os.Args) {
		fmt.Printf("positional arguments: ")

		for OptInd < len(os.Args) {
			fmt.Printf("%s ", os.Args[OptInd])
			OptInd++
		}

		fmt.Printf("\n")
	}
}
```

```sh
$ ./complex --foo --bar=bar --baz=baz --baz qux -- --quux quux
option 'a' has argument ''
option 'b' has argument 'bar'
option 'c' has argument 'baz'
option 'c' has argument 'qux'
positional arguments: --quux quux

$ ./complex --aa -bb -cc -c d -- -e e
option 'a' has argument ''
option 'a' has argument ''
option 'b' has argument 'b'
option 'c' has argument 'c'
option 'c' has argument 'd'
positional arguments: -e e

$ ./complex --on --off
option 'on' changed frob state to '1'
option 'off' changed frob state to '0'

$ ./complex --silent --foo --qux --bar
option 'a' has argument ''
option 'b' has argument ''
```

## Maintainers

[@BChristieDev](https://github.com/BChristieDev)

## Contributing

PRs accepted.

## License

[MIT](LICENSE) © Brandon Christie
