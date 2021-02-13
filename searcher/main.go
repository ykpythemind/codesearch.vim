package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logf, err := os.OpenFile(filepath.Join(home, "codesearch.vim.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logf.Close()

	wrt := io.MultiWriter(os.Stdout, logf)
	log.SetOutput(logf)

	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, "missing file")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	cwd := flag.String("cwd", "", "")
	flag.Parse()

	if err := run(*cwd, f, wrt); err != nil {
		log.Println(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

const metaMarker = "▿"

type SearchQuery struct {
	Pattern  string
	Includes string
	Excludes string
	Option   queryOption
}

type QueryParser struct {
	input          string
	patternReading bool
	pattern        string
}

type queryOption struct {
	useRegexp       bool
	caseSensitivity caseSensitivity
}

type caseSensitivity string

var (
	smartCase     caseSensitivity = "smartCase"
	ignoreCase    caseSensitivity = "ignoreCase"
	caseSensitive caseSensitivity = "caseSensitive"
)

func NewQueryParser(input string) *QueryParser {
	return &QueryParser{input: input}
}

func (p *QueryParser) Parse() (*SearchQuery, error) {
	scanner := bufio.NewScanner(strings.NewReader(p.input))
	p.patternReading = true

	includes := ""
	excludes := ""
	options := queryOption{}
	_ = options

	for scanner.Scan() {
		t := scanner.Text()
	scanned:

		if strings.HasPrefix(t, metaMarker+" includes") {
			p.endPatternReading()
			if scanner.Scan() {
				t = scanner.Text()
				if p.ignoreLine(t) {
					goto scanned
				} else {
					includes = strings.TrimSpace(t)
				}
			}
		} else if strings.HasPrefix(t, metaMarker+" excludes") {
			p.endPatternReading()
			if scanner.Scan() {
				t = scanner.Text()
				if p.ignoreLine(t) {
					goto scanned
				} else {
					excludes = strings.TrimSpace(t)
				}
			}
		} else if strings.HasPrefix(t, metaMarker+" options") {
			p.endPatternReading()
			if scanner.Scan() {
				t = scanner.Text()
				if p.ignoreLine(t) {
					goto scanned
				} else {
					opt, err := parseOptions(strings.TrimSpace(t))
					if err == nil {
						options = *opt
					}
				}
			}
		} else {
			// read pattern
			if p.patternReading {
				p.pattern = p.pattern + t
			}
		}
	}

	// 最後まで読んだ ---

	return &SearchQuery{Pattern: p.pattern, Includes: includes, Excludes: excludes, Option: options}, nil
}

func parseOptions(optStr string) (*queryOption, error) {
	optionLineRegexp := regexp.MustCompile(`caseOption:(.*)\|(\s*)useRegexp:(.*)\|`)
	matched := optionLineRegexp.FindStringSubmatch(optStr)

	// log.Printf("%+v\n", matched)
	if len(matched) != 4 {
		return nil, errors.New("format is wrong")
	}

	casestr := strings.TrimSpace(strings.ToLower(matched[1]))
	var casesence caseSensitivity
	if casestr == "smartcase" {
		casesence = smartCase
	} else if casestr == "ignorecase" {
		casesence = ignoreCase
	} else if casestr == "casesensitive" {
		casesence = caseSensitive
	}
	useregexp := strings.TrimSpace(matched[3])

	return &queryOption{useRegexp: useregexp == "true", caseSensitivity: casesence}, nil
}

func (p *QueryParser) ignoreLine(line string) bool {
	return strings.TrimSpace(line) == "" || strings.HasPrefix(line, metaMarker)
}

func (p *QueryParser) endPatternReading() {
	if !p.patternReading {
		return
	}
	log.Println(p.pattern)
	p.pattern = strings.TrimRight(p.pattern, "\n") // 最後の行の改行は削除する
	p.patternReading = false
}

func run(cwd string, in io.Reader, out io.Writer) error {
	if cwd != "" {
		current, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Chdir(cwd)
		if err != nil {
			return err
		}
		defer os.Chdir(current)
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	parser := NewQueryParser(string(b))
	query, err := parser.Parse()
	if err != nil {
		return err
	}

	rgargs, err := getRgArgs(*query)
	if err != nil {
		return err
	}

	var escapedArgs []string
	for _, arg := range rgargs {
		if !argRegexp.MatchString(arg) {
			// escape
			arg = fmt.Sprintf("'%s'", arg)
		}

		escapedArgs = append(escapedArgs, arg)
	}

	log.Printf("rgargs: %v\n", rgargs)

	// joinedArgs := strings.Join(rgargs, " ")

	cmd := exec.Command("rg", rgargs...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

type RgArgs []string

func (args *RgArgs) Append(otherArg ...string) {
	slice := *args
	slice = append(slice, otherArg...)
	*args = slice
}

var argRegexp = regexp.MustCompile("^-")

// see https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L378
func getRgArgs(query SearchQuery) (RgArgs, error) {
	args := RgArgs{}

	args.Append("--hidden")

	if query.Option.caseSensitivity == smartCase {
		args.Append("--smart-case")
	} else if query.Option.caseSensitivity == ignoreCase {
		args.Append("--ignore-case")
	} else if query.Option.caseSensitivity == caseSensitive {
		args.Append("--case-sensitive")
	} else {
		return nil, fmt.Errorf("%s is unknown", query.Option.caseSensitivity)
	}

	if query.Option.useRegexp {
		return nil, errors.New("regexp is not implemented")
	}

	var doublestarIncludes, otherIncludes []string
	_ = doublestarIncludes

	otherIncludes = strings.Split(query.Includes, ",")
	if otherIncludes[0] != "" {
		// todo: unique

		args.Append("-g", "!*")
		for _, in := range otherIncludes {
			// fixme .から始まるものは拡張子もマッチするという挙動が再現できない
			if strings.HasPrefix(in, ".") {
				args.Append("-g", fmt.Sprintf("*%s", in))
				continue
			}

			for _, glob := range spreadGlobComponents(in) {
				glob = anchorGlob(glob)
				args.Append("-g", glob)
			}
		}
	}

	// doubleStarIncludes

	// Allow $ to match /r/n
	args.Append("--crlf")

	args.Append("--vimgrep")

	// after double dashes

	var searchPatternAfterDoubleDashes string

	// do some parse, use regexp
	searchPatternAfterDoubleDashes = query.Pattern
	args.Append("--fixed-strings")

	args.Append("--")

	if searchPatternAfterDoubleDashes != "" {
		// Put the query after --, in case the query starts with a dash
		args.Append(searchPatternAfterDoubleDashes)
	}

	args.Append(".")

	return args, nil
}

// `"foo/*bar/something"` -> `["foo", "foo/*bar", "foo/*bar/something", "foo/*bar/something/**"]`
func spreadGlobComponents(globArg string) []string {
	components := splitGlobAware(globArg, '/')
	var ret []string

	l := len(components)
	_ = l
	for i := range components {
		r := components[0 : i+1]
		s := strings.Join(r, "/")
		ret = append(ret, s)
		// これ足りない？
		if i == l-1 && !strings.HasSuffix(s, "*") {
			s += "/**"
			ret = append(ret, s)
		}
	}

	return ret
}

func anchorGlob(glob string) string {
	if strings.HasPrefix(glob, "**") || strings.HasPrefix(glob, "/") {
		return glob
	} else {
		return "/" + glob
	}
}

func splitGlobAware(pattern string, splitChar rune) (segments []string) {
	if pattern == "" {
		return
	}
	inBraces := false
	inBrackets := false

	val := ""

	for _, char := range pattern {
		switch char {
		case splitChar:
			if !inBraces && !inBrackets {
				segments = append(segments, val)
				val = ""
				continue
			}
			break
		case '{':
			inBraces = true
			break
		case '}':
			inBraces = false
			break
		case '[':
			inBrackets = true
			break
		case ']':
			inBrackets = false
			break
		}

		val += string(char)
	}

	if val != "" {
		segments = append(segments, val)
	}

	return
}
