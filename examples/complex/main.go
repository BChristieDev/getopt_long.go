/*
	@file      examples/complex/main.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

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
