package darwin_test

import (
	"testing"

	// Frameworks
	mac "github.com/djthorpe/gopi-hw/darwin"
)

var (
	TEST_STRINGS = []string{
		"",
		"abcd",
		"condé nast",
		"paul schütze",
		"☺☻☹",
		"日a本b語ç日ð本Ê語þ日¥本¼語i日©",
		"日a本b語ç日ð本Ê語þ日¥本¼語i日©日a本b語ç日ð本Ê語þ日¥本¼語i日©日a本b語ç日ð本Ê語þ日¥本¼語i日©",
	}
)

func Test_CFString000(t *testing.T) {
	cfstr := mac.NewCFString("")
	defer cfstr.Free()
	if cfstr.String() != "" {
		t.Error("Expected empty string")
	}
}

func Test_CFString001(t *testing.T) {
	for _, str := range TEST_STRINGS {
		cfstr := mac.NewCFString(str)
		defer cfstr.Free()
		if cfstr.String() != str {
			t.Errorf("Expected string: '%v' got: '%v'", str, cfstr.String())
		}
	}
}
