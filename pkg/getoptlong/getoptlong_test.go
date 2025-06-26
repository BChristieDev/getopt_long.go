/*
	@file      pkg/getoptlong/getoptlong_test.go
	@author    Brandon Christie <bchristie.dev@gmail.com>
*/

package getoptlong_test

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	. "github.com/BChristieDev/getopt_long.go/pkg/getoptlong"
)

func cleanup(t *testing.T) {
	t.Helper()

	OptArg = ""
	OptInd = 1
	OptErr = 1
	OptOpt = 0
}

func TestLongOptions(t *testing.T) {
	t.Run("End of options delimiter", func(t *testing.T) {
		args := []string{"", "--foo", "--", "--bar"}
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 'f'},
			{Name: "bar", HasArg: NoArgument, Flag: nil, Val: 'b'},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'f' {
				t.Errorf("opt is '%c'. Expected 'f'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "--bar" {
			t.Errorf("positional argument is '%s'. Expected '--bar'.\n", args[OptInd])
		}
	})

	t.Run("Flag", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		var foo int
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: &foo, Val: 1},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%d'. Expected '0'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if foo != 1 {
			t.Errorf("flag 'foo' is '%d'. Expected '1'.\n", foo)
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Index pointer", func(t *testing.T) {
		args := []string{"", "--bar", "qux"}
		var indexptr int
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 0},
			{Name: "bar", HasArg: NoArgument, Flag: nil, Val: 0},
			{Name: "baz", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, &indexptr)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%d'. Expected '0'.\n", opt)
			}

			if longopts[indexptr].Name != "bar" {
				t.Errorf("indexptr is '%s'. Expected 'bar'.\n", longopts[indexptr].Name)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "qux" {
			t.Errorf("positional argument is '%s'. Expected 'qux'.\n", args[OptInd])
		}
	})

	t.Run("Invalid option silent", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w
		OptErr = 0

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "" {
				t.Errorf("stderr is '%s'. Expected ''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Invalid option", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "unrecognized option '--foo'" {
				t.Errorf("stderr is '%s'. Expected 'unrecognized option '--foo''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Expects no argument passed optional argument", func(t *testing.T) {
		args := []string{"", "--foo=bar", "baz"}
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "baz" {
			t.Errorf("positional argument is '%s'. Expected 'baz'.\n", args[OptInd])
		}
	})

	t.Run("Expects no argument", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Expects optional argument passed no argument", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		longopts := []Option{
			{Name: "foo", HasArg: OptionalArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Expects optional argument passed optional argument", func(t *testing.T) {
		args := []string{"", "--foo=bar", "baz"}
		longopts := []Option{
			{Name: "foo", HasArg: OptionalArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "bar" {
				t.Errorf("optarg is '%s'. Expected 'bar'.\n", OptArg)
			}
		}

		if args[OptInd] != "baz" {
			t.Errorf("positional argument is '%s'. Expected 'baz'.\n", args[OptInd])
		}
	})

	t.Run("Expects required argument passed no argument silent", func(t *testing.T) {
		args := []string{"", "--foo"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "foo", HasArg: RequiredArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w
		OptErr = 0

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != ':' {
				t.Errorf("opt is '%c'. Expected ':'.\n", opt)
			}

			if string(stderr) != "" {
				t.Errorf("stderr is '%s'. Expected ''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd-1] != "--foo" {
			t.Errorf("last argument is '%s'. Expected '--foo'.\n", args[OptInd-1])
		}
	})

	t.Run("Expects required argument passed no argument", func(t *testing.T) {
		args := []string{"", "--foo"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "foo", HasArg: RequiredArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "option '--foo' requires an argument" {
				t.Errorf("stderr is '%s'. Expected ''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd-1] != "--foo" {
			t.Errorf("last argument is '%s'. Expected '--foo'.\n", args[OptInd-1])
		}
	})

	t.Run("Expects required argument passed optional argument", func(t *testing.T) {
		args := []string{"", "--foo=bar", "baz"}
		longopts := []Option{
			{Name: "foo", HasArg: RequiredArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "bar" {
				t.Errorf("optarg is '%s'. Expected 'bar'.\n", OptArg)
			}
		}

		if args[OptInd] != "baz" {
			t.Errorf("positional argument is '%s'. Expected 'baz'.\n", args[OptInd])
		}
	})

	t.Run("Expects required argument passed required argument", func(t *testing.T) {
		args := []string{"", "--foo", "bar", "baz"}
		longopts := []Option{
			{Name: "foo", HasArg: RequiredArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 0 {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "bar" {
				t.Errorf("optarg is '%s'. Expected 'bar'.\n", OptArg)
			}
		}

		if args[OptInd] != "baz" {
			t.Errorf("positional argument is '%s'. Expected 'baz'.\n", args[OptInd])
		}
	})

	t.Run("Val", func(t *testing.T) {
		args := []string{"", "--foo", "bar"}
		longopts := []Option{
			{Name: "foo", HasArg: NoArgument, Flag: nil, Val: 'f'},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'f' {
				t.Errorf("opt is '%c'. Expected 'f'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})
}

func TestShortOptions(t *testing.T) {
	t.Run("End of options delimiter", func(t *testing.T) {
		args := []string{"", "-a", "--", "-b"}
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "ab", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'a' {
				t.Errorf("opt is '%c'. Expected '0'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "-b" {
			t.Errorf("positional argument is '%s'. Expected '-b'.\n", args[OptInd])
		}
	})

	t.Run("Invalid option silent", func(t *testing.T) {
		args := []string{"", "-a", "foo"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, ":", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "" {
				t.Errorf("stderr is '%s'. Expected ''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "foo" {
			t.Errorf("positional argument is '%s'. Expected 'foo'.\n", args[OptInd])
		}
	})

	t.Run("Invalid option", func(t *testing.T) {
		args := []string{"", "-a", "foo"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, "", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "invalid option -- a" {
				t.Errorf("stderr is '%s'. Expected 'invalid option -- a'.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "foo" {
			t.Errorf("positional argument is '%s'. Expected 'foo'.\n", args[OptInd])
		}
	})

	t.Run("Expects no argument", func(t *testing.T) {
		args := []string{"", "-abc", "foo"}
		optstring := "abc"
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, optstring, longopts, nil)

			if opt == -1 {
				break
			}

			if !strings.ContainsRune(optstring, rune(opt)) {
				t.Errorf("opt is '%t'. Expected 'true'.\n", strings.ContainsRune(optstring, rune(opt)))
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "foo" {
			t.Errorf("positional argument is '%s'. Expected 'foo'.\n", args[OptInd])
		}
	})

	t.Run("Expects optional argument passed no argument", func(t *testing.T) {
		args := []string{"", "-a", "foo"}
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "a::", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'a' {
				t.Errorf("opt is '%c'. Expected 'a'.\n", opt)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd] != "foo" {
			t.Errorf("positional argument is '%s'. Expected 'foo'.\n", args[OptInd])
		}
	})

	t.Run("Expects optional argument passed optional argument", func(t *testing.T) {
		args := []string{"", "-afoo", "bar"}
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "a::", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'a' {
				t.Errorf("opt is '%c'. Expected 'a'.\n", opt)
			}

			if OptArg != "foo" {
				t.Errorf("optarg is '%s'. Expected 'foo'.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Expects required argument passed no argument silent", func(t *testing.T) {
		args := []string{"", "-a"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, ":a:", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != ':' {
				t.Errorf("opt is '%c'. Expected ':'.\n", opt)
			}

			if string(stderr) != "" {
				t.Errorf("stderr is '%s'. Expected ''.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd-1] != "-a" {
			t.Errorf("last argument is '%s'. Expected '-a'.\n", args[OptInd-1])
		}
	})

	t.Run("Expects required argument passed no argument", func(t *testing.T) {
		args := []string{"", "-a"}
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		os.Stderr = w

		for {
			opt = GetoptLong(len(args), args, "a:", longopts, nil)

			if opt == -1 {
				break
			}

			w.Close()
			regex := regexp.MustCompile("[\r\n]")
			stderr, _ := io.ReadAll(r)
			stderr = regex.ReplaceAll(stderr, nil)
			os.Stderr = oldStderr

			if opt != '?' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if string(stderr) != "option requires an argument -- a" {
				t.Errorf("stderr is '%s'. Expected 'option requires an argument -- a'.\n", stderr)
			}

			if OptArg != "" {
				t.Errorf("optarg is '%s'. Expected ''.\n", OptArg)
			}
		}

		if args[OptInd-1] != "-a" {
			t.Errorf("last argument is '%s'. Expected '-a'.\n", args[OptInd-1])
		}
	})

	t.Run("Expects required argument passed optional argument", func(t *testing.T) {
		args := []string{"", "-afoo", "bar"}
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "a:", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'a' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if OptArg != "foo" {
				t.Errorf("optarg is '%s'. Expected 'foo'.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})

	t.Run("Expects required argument passed required argument", func(t *testing.T) {
		args := []string{"", "-a", "foo", "bar"}
		longopts := []Option{
			{Name: "", HasArg: NoArgument, Flag: nil, Val: 0},
		}
		var opt int

		t.Cleanup(func() { cleanup(t) })

		for {
			opt = GetoptLong(len(args), args, "a:", longopts, nil)

			if opt == -1 {
				break
			}

			if opt != 'a' {
				t.Errorf("opt is '%c'. Expected '?'.\n", opt)
			}

			if OptArg != "foo" {
				t.Errorf("optarg is '%s'. Expected 'foo'.\n", OptArg)
			}
		}

		if args[OptInd] != "bar" {
			t.Errorf("positional argument is '%s'. Expected 'bar'.\n", args[OptInd])
		}
	})
}
