package log

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_StatsLibLogger(t *testing.T) {
	sl := NewStatslibLogger()
	Convey("Printf", t, func() {
		So(func() {
			sl.Printf("test %v", "statslib")
		}, ShouldNotPanic)
	})
}
