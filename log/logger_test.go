package log

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func Test_Output(t *testing.T) {
	Convey("Test loadSyncer", t, func() {
		So(loadSyncer(&Config{}), ShouldEqual, os.Stdout)
	})
}

func Test_Fields(t *testing.T) {
	Stdout()
	Convey("Verfing fields output", t, func() {
		Convey("log.String", func() {
			So(func() { V(0).With(String("test", "test")).Info("test") }, ShouldNotPanic)
		})
		Convey("log.Int", func() {
			So(func() { V(0).With(Int("test", 1)).Info("test") }, ShouldNotPanic)
		})
		Convey("log.Type", func() {
			So(func() { V(0).With(Type("test")).Info("test") }, ShouldNotPanic)
		})
	})
}

func Test_Levels(t *testing.T) {
	Stdout()
	Convey("Verfing log level functions", t, func() {
		Convey("log.Info", func() {
			So(func() {
				Info("test")
				Infof("test %v", "infof")
			}, ShouldNotPanic)
		})
		Convey("log.Warn", func() {
			So(func() {
				Warn("test")
				Warnf("test %v", "warnf")
			}, ShouldNotPanic)
		})
		Convey("log.Error", func() {
			So(func() {
				Error("test")
				Errorf("test %v", "errorf")
			}, ShouldNotPanic)
		})
		Convey("log.Debug", func() {
			So(func() {
				Debug("test")
				Debugf("test %v", "debugf")
			}, ShouldNotPanic)
		})
		Convey("log.Verbose", func() {
			So(func() {
				V(0).Info("verbose 0")
				V(1).Info("verbose 1")
				V(2).Info("verbose 2")
			}, ShouldNotPanic)
		})
	})
}

func Test_Disable(t *testing.T) {
	Initialize(&Config{Disable: true})
	Convey("Verfing log level functions", t, func() {
		Convey("log.Info", func() {
			So(func() {
				Info("test")
				Infof("test %v", "infof")
			}, ShouldNotPanic)
		})
		Convey("log.Warn", func() {
			So(func() {
				Warn("test")
				Warnf("test %v", "warnf")
			}, ShouldNotPanic)
		})
		Convey("log.Error", func() {
			So(func() {
				Error("test")
				Errorf("test %v", "errorf")
			}, ShouldNotPanic)
		})
		Convey("log.Debug", func() {
			So(func() {
				Debug("test")
				Debugf("test %v", "debugf")
			}, ShouldNotPanic)
		})
		Convey("log.Verbose", func() {
			So(func() {
				V(0).Info("verbose 0")
				V(1).Info("verbose 1")
				V(2).Info("verbose 2")
			}, ShouldNotPanic)
		})
	})
}
