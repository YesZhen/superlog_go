package superlog

import (
	"testing"
	"time"
)

// func TestInitDefault(t *testing.T) {
// 	r, err := InitDefault()
// 	if err != nil || r == nil {
// 		t.Errorf("InitDefault() error, %v", err)
// 	}
// }

func TestReadWrite(t *testing.T) {
	r, err := InitDefault()
	if err != nil || r == nil {
		t.Errorf("InitDefault() error, %v", err)
	}


	SetName("RunC")

	t3 := time.Now().UnixNano()

	Log(r, "Hello!!!")

	t4 := time.Now().UnixNano()

	Log(r, "Yes!!!")

	t5 := time.Now().UnixNano()

	m := Read(r, 1)
	m = Read(r, 1)

	t6 := time.Now().UnixNano()

	if m == nil || len(m) <= 0 {
		t.Error("read error")
	}

	t.Log(m[0].Timestamp)
	t.Log(m[0].Pid)
	t.Log(m[0].Tid)
	t.Log(m[0].Data)
	t.Log(m[0].Name)

	t.Log(t3, t4, t5, t6)
}