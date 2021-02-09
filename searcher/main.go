package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)

	if err := run(os.Stdin, wrt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func ParseSearchQuery(str string) (*SearchQuery, error) {
	scanner := bufio.NewScanner(strings.NewReader(str))

	pattern := ""
	patternReading := true
	matchedEndpattern := false

	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "---===---") {
			matchedEndpattern = true
			patternReading = false
			// remove before line break
			pattern = strings.TrimRight(pattern, "\n")
		} else {
			// read pattern
			if patternReading {
				pattern += t
			}
		}
	}

	// 最後まで読んだ ---

	if !matchedEndpattern {
		pattern = strings.TrimRight(pattern, "\n")
	}

	return &SearchQuery{Pattern: pattern}, nil
}

func run(in io.Reader, out io.Writer) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	query, err := ParseSearchQuery(string(b))
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

type SearchQuery struct {
	Pattern string
}

type RgArgs []string

func (args *RgArgs) Append(otherArg string) {
	slice := *args
	slice = append(slice, otherArg)
	*args = slice
}

// see https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L378
func getRgArgs(query SearchQuery) (RgArgs, error) {
	args := RgArgs{}

	args.Append("--hidden")

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
