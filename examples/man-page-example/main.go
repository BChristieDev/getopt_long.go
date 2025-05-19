/*
	@file      examples/man-page-example/main.go
	@author    Brandon Christie <bchristie.dev@gmail.com>

	@see       https://linux.die.net/man/3/getopt_long
*/

package main

import (
	"fmt"
	"os"

	"github.com/BChristieDev/getopt_long.go/pkg/getoptlong"
)

func main() {
	var c int
	digit_optind := 0

	for {
		var this_option_optind int
		option_index := 0

		if getoptlong.OptInd > 0 {
			this_option_optind = getoptlong.OptInd
		} else {
			this_option_optind = 1
		}

		long_options := []getoptlong.Option{
			{Name: "add", HasArg: getoptlong.RequiredArgument, Flag: nil, Val: 0},
			{Name: "append", HasArg: getoptlong.NoArgument, Flag: nil, Val: 0},
			{Name: "delete", HasArg: getoptlong.RequiredArgument, Flag: nil, Val: 0},
			{Name: "verbose", HasArg: getoptlong.NoArgument, Flag: nil, Val: 0},
			{Name: "create", HasArg: getoptlong.RequiredArgument, Flag: nil, Val: 'c'},
			{Name: "file", HasArg: getoptlong.RequiredArgument, Flag: nil, Val: 0},
		}

		c = getoptlong.GetoptLong(len(os.Args), os.Args, "abc:d:012", long_options, &option_index)

		if c == -1 {
			break
		}

		switch c {
		case 0:
			fmt.Printf("options %s", long_options[option_index].Name)
			if getoptlong.OptArg != "" {
				fmt.Printf(" with arg %s", getoptlong.OptArg)
			}
			fmt.Printf("\n")
		case '0':
			fallthrough
		case '1':
			fallthrough
		case '2':
			if digit_optind != 0 && digit_optind != this_option_optind {
				fmt.Printf("digits occur in two different argv-elements.\n")
			}
			digit_optind = this_option_optind
			fmt.Printf("option %c\n", c)
		case 'a':
			fmt.Printf("option a\n")
		case 'b':
			fmt.Printf("option b\n")
		case 'c':
			fmt.Printf("option c with value '%s'\n", getoptlong.OptArg)
		case 'd':
			fmt.Printf("option d with value '%s'\n", getoptlong.OptArg)
		case '?':
		default:
			fmt.Printf("?? getopt returned character code 0%o ??\n", c)
		}
	}

	if getoptlong.OptInd < len(os.Args) {
		fmt.Printf("non-option ARGV-elements: ")
		for getoptlong.OptInd < len(os.Args) {
			fmt.Printf("%s ", os.Args[getoptlong.OptInd])
			getoptlong.OptInd++
		}
		fmt.Printf("\n")
	}
}
