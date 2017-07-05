package main

import (
	redo "github.com/gonet2/libs/nsq-redo"
	"testing"
	"time"
)

type testdoc struct {
	Name string
	Age  int
}

func TestRedo(t *testing.T) {
	doc := testdoc{}
	// subdoc
	r := redo.NewRedoRecord(1, "test1", ts())
	doc.Name = "name1"
	doc.Age = 18
	r.AddChange("test", "xxx", doc)
	redo.Publish(r)

	r = redo.NewRedoRecord(2, "test2", ts())
	doc.Name = "name2"
	doc.Age = 22
	r.AddChange("test", "", doc)
	redo.Publish(r)
}

const TS_MASK = 0x1FFFFFFFFFF

func ts() uint64 {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return (uint64(t) & TS_MASK) << 22
}
