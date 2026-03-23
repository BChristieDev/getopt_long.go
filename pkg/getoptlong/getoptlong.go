/*
	@file      pkg/getoptlong/getoptlong.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

package getoptlong

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BChristieDev/getopt_long.go/internal/common"
)

type Option struct {
	/* Name of the long option. */
	Name string
	/*
		NoArgument (or 0) if the option does not take an argument;
		RequiredArgument (or 1) if the option requires an argument; or
		OptionalArgument (or 2) if the option takes an optional argument.
	*/
	HasArg int
	/*
		Specifies how results are returned for a long option. If Flag is not nil, then GetoptLong
		returns 0 and Val will be assigned to the integer Flag is pointing to, otherwise GetoptLong
		returns Val.
	*/
	Flag *int
	/* Value to return, or be assigned to the integer Flag is pointing to. */
	Val int
}

const (
	/* No argument to the option is expected. */
	NoArgument = 0
	/* An argument to the option is required */
	RequiredArgument = 1
	/* An argument to the option may be presented */
	OptionalArgument = 2
)

var (
	/* Stores the argument of an option. */
	OptArg string
	/* Next argument in argv array to process; default 1. */
	OptInd int
	/* Error reporting flag, set to 0 to suppress default error messages; default 1 */
	OptErr int
	/* Stores option that causes an error. */
	OptOpt   int
	nextchar int
)

func init() {
	OptArg = ""
	OptInd = 1
	OptErr = 1
	OptOpt = 0
	nextchar = 0
}

func errInvalidOpt(optopt int, errMsg string) int {
	OptInd++
	OptOpt = optopt
	nextchar = 0

	if OptErr != 0 {
		fmt.Fprintln(os.Stderr, errMsg)
	}

	return '?'
}

func errRequiresArg(optopt int, errMsg string) int {
	OptOpt = optopt

	if OptErr != 0 {
		fmt.Fprintln(os.Stderr, errMsg)
		return '?'
	}

	return ':'
}

func parseArg(argc int, argv []string, hasArg int, optargind int, opt int, errMsg string) int {
	if hasArg == RequiredArgument && optargind <= 0 && OptInd >= argc {
		return errRequiresArg(opt, errMsg)
	}

	if hasArg == RequiredArgument || (hasArg == OptionalArgument && optargind > 0) {
		OptArg = argv[OptInd][optargind:]
		OptInd++
		nextchar = 0
	} else {
		OptArg = ""
	}

	return 0
}

func parseLongOpt(argc int, argv []string, longopts []Option, indexptr *int) int {
	progname := filepath.Base(argv[0])
	eq := common.IndexOf(argv[OptInd], "=", 3)
	var opt string

	if eq == -1 {
		opt = argv[OptInd][2:]
	} else {
		opt = argv[OptInd][2:eq]
	}

	optarrind := common.FindIndex(longopts, func(longopt Option) bool { return longopt.Name == opt })

	if optarrind == -1 {
		return errInvalidOpt(0, fmt.Sprintf("%s: unrecognized option '--%s'", progname, opt))
	}

	if longopts[optarrind].HasArg <= NoArgument || longopts[optarrind].HasArg > OptionalArgument || eq == -1 {
		OptInd++
	}

	if indexptr != nil {
		*indexptr = optarrind
	}

	err := parseArg(argc, argv, longopts[optarrind].HasArg, eq+1, 0, fmt.Sprintf("%s: option '--%s' requires an argument", progname, opt))

	if err != 0 {
		return err
	}

	if longopts[optarrind].Flag == nil {
		return longopts[optarrind].Val
	}

	*longopts[optarrind].Flag = longopts[optarrind].Val

	return 0
}

func parseShortOpt(argc int, argv []string, shortopts string) int {
	progname := filepath.Base(argv[0])
	opt := int(argv[OptInd][nextchar])
	optstrind := strings.Index(shortopts, string(rune(opt)))
	var hasArg int

	if optstrind == -1 {
		return errInvalidOpt(opt, fmt.Sprintf("%s: invalid option -- '%c'", progname, opt))
	}

	nextchar++

	if nextchar == len(argv[OptInd]) {
		OptInd++
		nextchar = 0
	}

	if common.CharAt(shortopts, optstrind+1) == ":" && common.CharAt(shortopts, optstrind+2) == ":" {
		hasArg = RequiredArgument
	} else if common.CharAt(shortopts, optstrind+1) == ":" && common.CharAt(shortopts, optstrind+2) != ":" {
		hasArg = OptionalArgument
	} else {
		hasArg = NoArgument
	}

	err := parseArg(argc, argv, hasArg, nextchar, opt, fmt.Sprintf("%s: option requires an argument -- '%c'", progname, opt))

	if err != 0 {
		return err
	}

	return opt
}

/*
If a short option is recognized the option character is returned. If a long option is recognized
Val is returned if Flag is nil, otherwise 0 is returned and Val is assigned to the integer Flag
is pointing to. If indexptr is not nil, then the index of the long option in longopts is assigned to
the integer indexptr is pointing to.

If an unrecognized option is encountered '?' is returned. If an option with a missing argument is
encountered '?' is returned with OptErr is is non-zero, otherwise ':' is returned.

If all options are parsed -1 is returned.
*/
func GetoptLong(argc int, argv []string, shortopts string, longopts []Option, indexptr *int) int {
	if OptInd >= argc {
		return -1
	}

	if common.CharAt(shortopts, 0) == ":" {
		OptErr = 0
	}

	if nextchar == 0 {
		if common.CharAt(argv[OptInd], 0) != "-" || argv[OptInd] == "-" {
			return -1
		}

		if argv[OptInd] == "--" {
			OptInd++
			return -1
		}

		if common.CharAt(argv[OptInd], 1) == "-" {
			return parseLongOpt(argc, argv, longopts, indexptr)
		}

		nextchar++
	}

	return parseShortOpt(argc, argv, shortopts)
}
