package superlog_go

import (
	"time"
	"fmt"

	"github.com/YesZhen/superlog_go/superlog"
)

var r *superlog.Ring

func init() {
	r, _ = superlog.InitDefault()
	// superlog.SetName("RunC")
}

func SetName(name string) {
	superlog.SetName(name)
}

func Log(data string) {
	superlog.Log(r, data)
}

func Info(data string) {
	Infof(data)
}

func Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if len(s) == 0 {
		return
	}
	
	n := superlog.MAX_DATA_LEN

	times := ((len(s) - 1) / n) + 1

	for i := 0; i < times; i++ {
		if i != times - 1 {
			superlog.Log(r, s[n * i : n * i + n])
		} else {
			superlog.Log(r, s[n * i :])
		}
	}
}

func LogBegin(data string) (string, int64) {
	return data, (superlog.Log(r, data + " begin")).Timestamp
}

func LogEnd(data string, beginT int64) {
	endT := superlog.Log(r, data + " end").Timestamp
	s := fmt.Sprintf("%s cost: %v", data,  time.Duration(endT - beginT) * time.Microsecond)
	superlog.Log(r, s)
}

func Read() superlog.Msg {
	return superlog.Read(r, 1)[0]
} 