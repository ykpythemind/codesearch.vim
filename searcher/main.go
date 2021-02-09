package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)

	if err := run(wrt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(out io.Writer) error {

	// ParseSearchQueryFromInput

	q := SearchQuery{Pattern: "ma"}

	rgargs, err := getRgArgs(q)
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
