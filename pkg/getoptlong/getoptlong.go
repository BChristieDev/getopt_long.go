/*
	@file      pkg/getoptlong/getoptlong.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

package getoptlong

import (
	"fmt"
	"os"
	"strings"

	"github.com/BChristieDev/getopt_long.go/internal/common"
)

type Option struct {
	Name   string
	HasArg int
	Flag   *int
	Val    int
}

const NoArgument, RequiredArgument, OptionalArgument int = 0, 1, 2

var OptArg string
var OptInd, OptErr, OptOpt int
var nextchar int

func init() {
	OptArg = ""
	OptInd = 1
	OptErr = 1
	OptOpt = 0
	nextchar = 0
}

func errInvalidOpt(opt int, errMsg string) int {
	OptInd++
	OptOpt = opt
	nextchar = 0

	if OptErr != 0 {
		fmt.Fprintln(os.Stderr, errMsg)
	}

	return '?'
}

func errRequiresArg(opt int, errMsg string) int {
	OptOpt = opt

	if OptErr != 0 {
		fmt.Fprintln(os.Stderr, errMsg)
		return '?'
	}

	return ':'
}

func parseArg(argc int, argv []string, isOptional bool, isRequired bool, opt int, optargind int, errMsg string) int {
	if (isOptional && optargind > 0) || isRequired {
		if isRequired && optargind <= 0 && OptInd >= argc {
			return errRequiresArg(opt, errMsg)
		}

		OptArg = argv[OptInd][optargind:]
		OptInd++
		nextchar = 0
	} else {
		OptArg = ""
	}

	return 0
}

func parseLongOpt(argc int, argv []string, longopts []Option, indexptr *int) int {
	eq := common.IndexOf(argv[OptInd], "=", 3)
	var opt string

	if eq == -1 {
		opt = argv[OptInd][2:]
	} else {
		opt = argv[OptInd][2:eq]
	}

	optarrIndex := common.FindIndex(longopts, func(longopt Option) bool { return longopt.Name == opt })

	if optarrIndex == -1 {
		return errInvalidOpt(0, fmt.Sprintf("unrecognized option '--%s'", opt))
	}

	if longopts[optarrIndex].HasArg <= NoArgument || longopts[optarrIndex].HasArg > OptionalArgument || eq == -1 {
		OptInd++
	}

	*indexptr = optarrIndex

	isOptional := longopts[optarrIndex].HasArg == OptionalArgument
	isRequired := longopts[optarrIndex].HasArg == RequiredArgument

	err := parseArg(argc, argv, isOptional, isRequired, 0, eq+1, fmt.Sprintf("option '--%s' requires an argument", opt))

	if err != 0 {
		return err
	}

	if longopts[optarrIndex].Flag == nil {
		return longopts[optarrIndex].Val
	}

	*longopts[optarrIndex].Flag = longopts[optarrIndex].Val

	return 0
}

func parseShortOpt(argc int, argv []string, shortopts string) int {
	opt := int(argv[OptInd][nextchar])
	optstrIndex := strings.Index(shortopts, string(rune(opt)))

	if optstrIndex == -1 {
		return errInvalidOpt(opt, fmt.Sprintf("invalid option -- %c", opt))
	}

	nextchar++

	if nextchar == len(argv[OptInd]) {
		OptInd++
		nextchar = 0
	}

	isOptional := common.CharAt(shortopts, optstrIndex+1) == ":" && common.CharAt(shortopts, optstrIndex+2) == ":"
	isRequired := common.CharAt(shortopts, optstrIndex+1) == ":" && common.CharAt(shortopts, optstrIndex+2) != ":"

	err := parseArg(argc, argv, isOptional, isRequired, opt, nextchar, fmt.Sprintf("option requires an argument -- %c", opt))

	if err != 0 {
		return err
	}

	return opt
}

func GetoptLong(argc int, argv []string, shortopts string, longopts []Option, indexptr *int) int {
	if OptInd >= argc {
		return -1
	}

	if nextchar == 0 {
		if argv[OptInd][0] != '-' || argv[OptInd] == "-" {
			return -1
		}

		if argv[OptInd] == "--" {
			OptInd++
			return -1
		}

		if argv[OptInd][1] == '-' {
			return parseLongOpt(argc, argv, longopts, indexptr)
		}

		nextchar++
	}

	return parseShortOpt(argc, argv, shortopts)
}
