package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexMatch(t *testing.T) {
	regex := "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}"

	assert.True(t, regexMatch("test@acme.com", regex))
	assert.True(t, regexMatch("Test@Acme.Com", regex))
	assert.False(t, regexMatch("test", regex))
	assert.False(t, regexMatch("test.com", regex))
	assert.False(t, regexMatch("test@acme", regex))
}

func TestMustRegexMatch(t *testing.T) {
	regex := "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}"

	o, err := mustRegexMatch("test@acme.com", regex)
	assert.True(t, o)
	assert.Nil(t, err)

	o, err = mustRegexMatch("Test@Acme.Com", regex)
	assert.True(t, o)
	assert.Nil(t, err)

	o, err = mustRegexMatch("test", regex)
	assert.False(t, o)
	assert.Nil(t, err)

	o, err = mustRegexMatch("test.com", regex)
	assert.False(t, o)
	assert.Nil(t, err)

	o, err = mustRegexMatch("test@acme", regex)
	assert.False(t, o)
	assert.Nil(t, err)
}

func TestRegexFindAll(t *testing.T) {
	regex := "a{2}"
	assert.Equal(t, 1, len(regexFindAll("aa", regex, -1)))
	assert.Equal(t, 1, len(regexFindAll("aaaaaaaa", regex, 1)))
	assert.Equal(t, 2, len(regexFindAll("aaaa", regex, -1)))
	assert.Equal(t, 0, len(regexFindAll("none", regex, -1)))
}

func TestMustRegexFindAll(t *testing.T) {
	type args struct {
		regex, s string
		n        int
	}
	cases := []struct {
		expected int
		args     args
	}{
		{1, args{"a{2}", "aa", -1}},
		{1, args{"a{2}", "aaaaaaaa", 1}},
		{2, args{"a{2}", "aaaa", -1}},
		{0, args{"a{2}", "none", -1}},
	}

	for _, c := range cases {
		res, err := mustRegexFindAll(c.args.s, c.args.regex, c.args.n)
		if err != nil {
			t.Errorf("regexFindAll test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, len(res), "case %#v", c.args)
	}
}

func TestRegexFindl(t *testing.T) {
	regex := "fo.?"
	assert.Equal(t, "foo", regexFind("foorbar", regex))
	assert.Equal(t, "foo", regexFind("foo foe fome", regex))
	assert.Equal(t, "", regexFind("none", regex))
}

func TestMustRegexFindl(t *testing.T) {
	type args struct{ regex, s string }
	cases := []struct {
		expected string
		args     args
	}{
		{"foo", args{"fo.?", "foorbar"}},
		{"foo", args{"fo.?", "foo foe fome"}},
		{"", args{"fo.?", "none"}},
	}

	for _, c := range cases {
		res, err := mustRegexFind(c.args.s, c.args.regex)
		if err != nil {
			t.Errorf("regexFind test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexReplaceAll(t *testing.T) {
	regex := "a(x*)b"
	assert.Equal(t, "-T-T-", regexReplaceAll("-ab-axxb-", regex, "T"))
	assert.Equal(t, "--xx-", regexReplaceAll("-ab-axxb-", regex, "$1"))
	assert.Equal(t, "---", regexReplaceAll("-ab-axxb-", regex, "$1W"))
	assert.Equal(t, "-W-xxW-", regexReplaceAll("-ab-axxb-", regex, "${1}W"))
}

func TestMustRegexReplaceAll(t *testing.T) {
	type args struct{ regex, s, repl string }
	cases := []struct {
		expected string
		args     args
	}{
		{"-T-T-", args{"a(x*)b", "-ab-axxb-", "T"}},
		{"--xx-", args{"a(x*)b", "-ab-axxb-", "$1"}},
		{"---", args{"a(x*)b", "-ab-axxb-", "$1W"}},
		{"-W-xxW-", args{"a(x*)b", "-ab-axxb-", "${1}W"}},
	}

	for _, c := range cases {
		res, err := mustRegexReplaceAll(c.args.s, c.args.regex, c.args.repl)
		if err != nil {
			t.Errorf("regexReplaceAll test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexReplaceAllLiteral(t *testing.T) {
	regex := "a(x*)b"
	assert.Equal(t, "-T-T-", regexReplaceAllLiteral("-ab-axxb-", regex, "T"))
	assert.Equal(t, "-$1-$1-", regexReplaceAllLiteral("-ab-axxb-", regex, "$1"))
	assert.Equal(t, "-${1}-${1}-", regexReplaceAllLiteral("-ab-axxb-", regex, "${1}"))
}

func TestMustRegexReplaceAllLiteral(t *testing.T) {
	type args struct{ regex, s, repl string }
	cases := []struct {
		expected string
		args     args
	}{
		{"-T-T-", args{"a(x*)b", "-ab-axxb-", "T"}},
		{"-$1-$1-", args{"a(x*)b", "-ab-axxb-", "$1"}},
		{"-${1}-${1}-", args{"a(x*)b", "-ab-axxb-", "${1}"}},
	}

	for _, c := range cases {
		res, err := mustRegexReplaceAllLiteral(c.args.s, c.args.regex, c.args.repl)
		if err != nil {
			t.Errorf("regexReplaceAllLiteral test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, res, "case %#v", c.args)
	}
}

func TestRegexSplit(t *testing.T) {
	regex := "a"
	assert.Equal(t, 4, len(regexSplit("banana", regex, -1)))
	assert.Equal(t, 0, len(regexSplit("banana", regex, 0)))
	assert.Equal(t, 1, len(regexSplit("banana", regex, 1)))
	assert.Equal(t, 2, len(regexSplit("banana", regex, 2)))

	regex = "z+"
	assert.Equal(t, 2, len(regexSplit("pizza", regex, -1)))
	assert.Equal(t, 0, len(regexSplit("pizza", regex, 0)))
	assert.Equal(t, 1, len(regexSplit("pizza", regex, 1)))
	assert.Equal(t, 2, len(regexSplit("pizza", regex, 2)))
}

func TestMustRegexSplit(t *testing.T) {
	type args struct {
		regex, s string
		n        int
	}
	cases := []struct {
		expected int
		args     args
	}{
		{4, args{"a", "banana", -1}},
		{0, args{"a", "banana", 0}},
		{1, args{"a", "banana", 1}},
		{2, args{"a", "banana", 2}},
		{2, args{"z+", "pizza", -1}},
		{0, args{"z+", "pizza", 0}},
		{1, args{"z+", "pizza", 1}},
		{2, args{"z+", "pizza", 2}},
	}

	for _, c := range cases {
		res, err := mustRegexSplit(c.args.s, c.args.regex, c.args.n)
		if err != nil {
			t.Errorf("regexSplit test case %v failed with err %s", c, err)
		}
		assert.Equal(t, c.expected, len(res), "case %#v", c.args)
	}
}

func TestRegexQuoteMeta(t *testing.T) {
	assert.Equal(t, "1\\.2\\.3", regexQuoteMeta("1.2.3"))
	assert.Equal(t, "pretzel", regexQuoteMeta("pretzel"))
}
