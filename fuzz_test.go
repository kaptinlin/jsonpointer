package jsonpointer

import "testing"

func FuzzParseRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("/")
	f.Add("/foo")
	f.Add("/foo/bar")
	f.Add("/users/0/name")
	f.Add("/~0~1")
	f.Add("/foo%20bar")
	f.Add("/special/chars!@#$%^&*()")
	f.Add("/array/-")
	f.Add("/with spaces/and\ttabs")

	f.Fuzz(func(t *testing.T, pointer string) {
		p, err := Parse(pointer)
		if err != nil {
			return
		}

		roundtrip, err := Parse(p.String())
		if err != nil {
			t.Fatalf("Parse(%q) after String() returned error: %v", p.String(), err)
		}
		if p.String() != roundtrip.String() {
			t.Fatalf("roundtrip mismatch: %q -> %q", p.String(), roundtrip.String())
		}
	})
}

func FuzzEscapeUnescapeToken(f *testing.F) {
	f.Add("")
	f.Add("simple")
	f.Add("with/slash")
	f.Add("with~tilde")
	f.Add("~/both")
	f.Add("~0~1")
	f.Add("unicode")

	f.Fuzz(func(t *testing.T, token string) {
		escaped := EscapeToken(token)
		got, err := UnescapeToken(escaped)
		if err != nil {
			t.Fatalf("UnescapeToken(EscapeToken(%q)) returned error: %v", token, err)
		}
		if got != token {
			t.Fatalf("escape roundtrip mismatch: %q -> %q -> %q", token, escaped, got)
		}
	})
}
