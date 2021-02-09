package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
}

type QueryParser struct {
	input          string
	patternReading bool
	pattern        string
}

func NewQueryParser(input string) *QueryParser {
	return &QueryParser{input: input}
}

func (p *QueryParser) Parse() (*SearchQuery, error) {
	scanner := bufio.NewScanner(strings.NewReader(p.input))
	p.patternReading = true

	includes := ""
	excludes := ""

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
					includes = trim(t)
				}
			}
		} else if strings.HasPrefix(t, metaMarker+" excludes") {
			p.endPatternReading()
			if scanner.Scan() {
				t = scanner.Text()
				if p.ignoreLine(t) {
					goto scanned
				} else {
					excludes = trim(t)
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

	return &SearchQuery{Pattern: p.pattern, Includes: includes, Excludes: excludes}, nil
}

func trim(str string) string {
	return strings.TrimSpace(str)
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

	/*
		const cwd = options.folder.fsPath;

		const escapedArgs = rgArgs
			.map(arg => arg.match(/^-/) ? arg : `'${arg}'`)
			.join(' ');
		this.outputChannel.appendLine(`${rgDiskPath} ${escapedArgs}\n - cwd: ${cwd}`);
	*/

	log.Printf("rgargs: %v\n", rgargs)

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

// see https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L378
func getRgArgs(query SearchQuery) (RgArgs, error) {
	args := RgArgs{}

	args.Append("--hidden")

	var doublestarIncludes, otherIncludes []string
	_ = doublestarIncludes

	otherIncludes = strings.Split(query.Includes, ",")
	if otherIncludes[0] != "" {
		// todo: unique

		args.Append("-g", "!*")
		for _, in := range otherIncludes {
			// want this logic https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L393
			globArg := in
			args.Append("-g", globArg)
		}
	}

	// Allow $ to match /r/n
	args.Append("--crlf")

	args.Append("--vimgrep")

	// after double dashes

	var searchPatternAfterDoubleDashes string

	// do some parse
	searchPatternAfterDoubleDashes = query.Pattern

	// これで区別が必要
	args.Append("--")

	if searchPatternAfterDoubleDashes != "" {
		// Put the query after --, in case the query starts with a dash
		args.Append(searchPatternAfterDoubleDashes)
	}

	args.Append(".")

	return args, nil
}
