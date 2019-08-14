package log

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_UdpSyncer(t *testing.T) {
	Convey("Write udp stream", t, func() {
		writer := NewAsyncUdpWriter("no where")
		l, err := writer.Write([]byte("test"))
		So(l, ShouldEqual, 4)
		So(err, ShouldEqual, nil)
	})
}
