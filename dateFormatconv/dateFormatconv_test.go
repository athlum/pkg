package dateFormatconv

import (
	"fmt"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	s, err := Format("yyyy-MM-dd HH:mm:ss.SSS")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)

	n, err := time.ParseInLocation(s, "2018-07-11 03:44:45.905", time.Local)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(n)
}
