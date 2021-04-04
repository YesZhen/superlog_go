package superlog_go

import "testing"

func TestLog(t *testing.T) {
	Log("Yep")
}

func TestRead(t *testing.T) {
	m := Read()
	t.Log(m.Name)
	t.Log(m.Timestamp)
}