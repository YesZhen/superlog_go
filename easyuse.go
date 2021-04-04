package superlog_go

import (
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

func Read() superlog.Msg {
	return superlog.Read(r, 1)[0]
} 