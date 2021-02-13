package main

import (
	"reflect"
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

func TestParseOptions(t *testing.T) {
	type test struct {
		input string
		want  *queryOption
	}

	tc := []test{
		{input: "caseOption: smartcase | useRegexp: false | useIgnoreSettingFile: false", want: &queryOption{caseSensitivity: smartCase, useRegexp: false}},
		{input: "caseOption: ignorecase | useRegexp: true | useIgnoreSettingFile: false", want: &queryOption{caseSensitivity: ignoreCase, useRegexp: true}},
		{input: "caseOption:ignorecase|useRegexp: true| useIgnoreSettingFile: false", want: &queryOption{caseSensitivity: ignoreCase, useRegexp: true}},
		{input: "caseOption: ignorecase | useRegexp:", want: nil},
		{input: "caseOption: ignorecase | useRegexp: false | useIgnoreSettingFile: true", want: &queryOption{caseSensitivity: ignoreCase, useIgnoreSettingFile: true}},
	}

	for _, tc := range tc {
		tc := tc

		got, err := parseOptions(tc.input)
		if err != nil {
			if tc.want != nil {
				t.Error(err)
			}
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("input: %s, want: %v, got: %v", tc.input, tc.want, got)
		}
	}

}
