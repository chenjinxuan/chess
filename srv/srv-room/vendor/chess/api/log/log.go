package log

import (
	"github.com/Sirupsen/logrus"
)

var (
	Log = logrus.New()
)

type WrapLog struct{}

func (l *WrapLog) Output(calldepth int, s string) error {
	Log.Debugf("Mgo Out(depth:%d,s:(%s)", calldepth, s)
	return nil
}
