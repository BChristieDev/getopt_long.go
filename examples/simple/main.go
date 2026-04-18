/*
	@file      examples/simple/main.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

package main

import (
	"fmt"
	"os"

	"github.com/BChristieDev/getopt_long.go/pkg/getoptlong"
)

func main() {
	longopts := []getoptlong.Option{
		{Name: "foo", HasArg: getoptlong.RequiredArgument, Flag: nil, Val: 0},
	}

	var longindex, opt int

	for {
		opt = getoptlong.Parse(len(os.Args), os.Args, "a:", longopts, &longindex)

		if opt == -1 {
			break
		}

		switch opt {
		case 0:
			fmt.Printf("option '%s' has argument '%s'\n", longopts[longindex].Name, getoptlong.OptArg)
		case 'a':
			fmt.Printf("option '%c' has argument '%s'\n", opt, getoptlong.OptArg)
		}
	}
}
