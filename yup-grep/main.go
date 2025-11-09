package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	gloo "github.com/gloo-foo/framework"
	. "github.com/yupsh/grep"
)

const (
	flagIgnoreCase       = "ignore-case"
	flagLineNumber       = "line-number"
	flagCount            = "count"
	flagInvertMatch      = "invert-match"
	flagWordRegexp       = "word-regexp"
	flagFixedStrings     = "fixed-strings"
	flagRecursive        = "recursive"
	flagFilesWithMatches = "files-with-matches"
	flagQuiet            = "quiet"
)

func main() {
	app := &cli.App{
		Name:  "grep",
		Usage: "search for patterns in files",
		UsageText: `grep [OPTIONS] PATTERN [FILE...]

   Search for PATTERN in each FILE.
   When no FILE is specified, read from standard input.
   PATTERN is a regular expression by default.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    flagIgnoreCase,
				Aliases: []string{"i"},
				Usage:   "ignore case distinctions in patterns and data",
			},
			&cli.BoolFlag{
				Name:    flagLineNumber,
				Aliases: []string{"n"},
				Usage:   "prefix each line of output with the line number",
			},
			&cli.BoolFlag{
				Name:    flagCount,
				Aliases: []string{"c"},
				Usage:   "only print a count of matching lines per FILE",
			},
			&cli.BoolFlag{
				Name:    flagInvertMatch,
				Aliases: []string{"v"},
				Usage:   "select non-matching lines",
			},
			&cli.BoolFlag{
				Name:    flagWordRegexp,
				Aliases: []string{"w"},
				Usage:   "match only whole words",
			},
			&cli.BoolFlag{
				Name:    flagFixedStrings,
				Aliases: []string{"F"},
				Usage:   "interpret PATTERN as a fixed string, not a regular expression",
			},
			&cli.BoolFlag{
				Name:    flagRecursive,
				Aliases: []string{"r", "R"},
				Usage:   "read all files under each directory, recursively",
			},
			&cli.BoolFlag{
				Name:    flagFilesWithMatches,
				Aliases: []string{"l"},
				Usage:   "print only names of FILEs with selected lines",
			},
			&cli.BoolFlag{
				Name:    flagQuiet,
				Aliases: []string{"q", "silent"},
				Usage:   "suppress all normal output",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "grep: %v\n", err)
		os.Exit(2)
	}
}

func action(c *cli.Context) error {
	// Require at least the pattern
	if c.NArg() < 1 {
		return fmt.Errorf("missing pattern argument")
	}

	// Extract pattern
	pattern := Pattern(c.Args().Get(0))

	// Build parameters array with files and flags
	var params []any

	// Add file arguments (or none for stdin)
	for i := 1; i < c.NArg(); i++ {
		params = append(params, gloo.File(c.Args().Get(i)))
	}

	// Add flags based on CLI options
	if c.Bool(flagIgnoreCase) {
		params = append(params, IgnoreCase)
	}
	if c.Bool(flagLineNumber) {
		params = append(params, LineNumber)
	}
	if c.Bool(flagCount) {
		params = append(params, Count)
	}
	if c.Bool(flagInvertMatch) {
		params = append(params, Invert)
	}
	if c.Bool(flagWordRegexp) {
		params = append(params, WholeWord)
	}
	if c.Bool(flagFixedStrings) {
		params = append(params, FixedStrings)
	}
	if c.Bool(flagRecursive) {
		params = append(params, Recursive)
	}
	if c.Bool(flagFilesWithMatches) {
		params = append(params, FilesOnly)
	}
	if c.Bool(flagQuiet) {
		params = append(params, Quiet)
	}

	// Create and execute the grep command
	cmd := Grep(pattern, params...)
	return gloo.Run(cmd)
}
