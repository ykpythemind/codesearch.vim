package main

import (
	"testing"
)

func TestParseSearchQuery(t *testing.T) {
	type test struct {
		input string
		want  SearchQuery
	}

	tc := []test{
		{input: `hoge
`, want: SearchQuery{Pattern: "hoge"}},
		{input: `fuga ▿`, want: SearchQuery{Pattern: "fuga ▿"}},
		{input: `piyo

fuga
`, want: SearchQuery{Pattern: "piyo\n\nfuga"}},
		{input: `piyo
▿ includes
app/,*.jpg`, want: SearchQuery{Pattern: "piyo", Includes: "app/,*.jpg"}},
		{input: `piyo

▿ includes
app/,*.jpg

▿ excludes
hoge
`, want: SearchQuery{Pattern: "piyo\n", Includes: "app/,*.jpg", Excludes: "hoge"}},
	}

	for _, tc := range tc {
		tc := tc

		parser := NewQueryParser(tc.input)
		q, err := parser.Parse()
		if err != nil {
			t.Error(err)
		}

		if *q != tc.want {
			t.Errorf("bad: want %+v got %+v", tc.want, *q)
		}
	}
}
