package superlog_go

import (
	// "fmt"
	"testing"
	// "github.com/YesZhen/superlog_go/superlog"
)

// func TestLog(t *testing.T) {
// 	Log("Yep")
// }

// func TestRead(t *testing.T) {
// 	m := Read()
// 	t.Log(m.Name)
// 	t.Log(m.Timestamp)
// }

func TestBeginEnd(t *testing.T) {
	testBeginEnd()

	m := Read()
	t.Log(m.Name)
	t.Log(m.Data)
	t.Log(m.Timestamp)

	m = Read()
	t.Log(m.Name)
	t.Log(m.Data)
	t.Log(m.Timestamp)

	m = Read()
	t.Log(m.Name)
	t.Log(m.Data)
	t.Log(m.Timestamp)

}

func testBeginEnd() {
	defer LogEnd(LogBegin("test"))


	d, t := LogBegin("a")
	//...
	LogEnd(d, t)
}

// func TestInfof(t *testing.T) {
// 	Infof("%s, %s, %s, %s", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "cccccccccccccccccccccccccccccc", "dddddddddddddddddddddddddddddd")

// 	m := Read()
// 	t.Log(m.Name)
// 	t.Log(m.Data)
// 	t.Log(m.Timestamp)

// 	m = Read()
// 	t.Log(m.Name)
// 	t.Log(m.Data)
// 	t.Log(m.Timestamp)

// 	m = Read()
// 	t.Log(m.Name)
// 	t.Log(m.Data)
// 	t.Log(m.Timestamp)

// 	m = Read()
// 	t.Log(m.Name)
// 	t.Log(m.Data)
// 	t.Log(m.Timestamp)
// }