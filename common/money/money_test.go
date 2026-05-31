package money

import "testing"

func TestFormatYuan(t *testing.T) {
	cases := map[int64]string{0: "0.00", 5: "0.05", 1800: "18.00", 1805: "18.05", -1250: "-12.50"}
	for in, want := range cases {
		if got := FormatYuan(in); got != want {
			t.Errorf("FormatYuan(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestParseYuan(t *testing.T) {
	cases := map[string]int64{"18.00": 1800, "18": 1800, "18.5": 1850, "0.05": 5, "-12.50": -1250}
	for in, want := range cases {
		got, err := ParseYuan(in)
		if err != nil {
			t.Errorf("ParseYuan(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseYuan(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestParseYuan_Invalid(t *testing.T) {
	for _, in := range []string{"", "abc", "1.234"} {
		if _, err := ParseYuan(in); err == nil {
			t.Errorf("ParseYuan(%q) expected error", in)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	c, err := ParseYuan("123.45")
	if err != nil {
		t.Fatal(err)
	}
	if got := FormatYuan(c); got != "123.45" {
		t.Errorf("roundtrip = %q", got)
	}
}
