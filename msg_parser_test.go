package link

import (
	"bytes"
	"github.com/playnb/util"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

var mp = NewMsgParser()

func TestMsgParser(t *testing.T) {
	convey.Convey("TestMsgParser", t, func() {
		buf := bytes.NewBuffer(nil)

		payLoad := "123123123AA"
		src, err := mp.WriteAll([]byte(payLoad))
		convey.So(err, convey.ShouldBeNil)

		_, err = buf.Write(src.GetPayload())
		convey.So(err, convey.ShouldBeNil)

		dst, err := mp.Read(buf)
		convey.So(err, convey.ShouldBeNil)

		convey.So([]byte(payLoad), util.ShouldByteSilceEqual, dst.GetPayload())
	})
}
