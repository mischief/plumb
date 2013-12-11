package plumb

import (
	"bytes"
	"code.google.com/p/goplan9/plan9"
	"testing"
)

var (
	tmsg = &Msg{
		Src:  "plumb_test",
		Dst:  "edit",
		Wdir: "/tmp",
		Type: "text",
		Attr: map[string]string{
			"quux": "foo",
			"foo":  "bar baz",
		},
		Data: []byte("foo.c")}
)

func TestPlumbSendRecv(t *testing.T) {
	pls, err := Open("send", plan9.OWRITE)
	if err != nil {
		t.Fatal(err)
	}

	pl, err := Open("edit", plan9.OREAD)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("omsg %s", tmsg)
	pls.Send(tmsg)

	msg, err := pl.Recv()

	t.Logf("imsg %s", msg)

	if err != nil {
		t.Fatalf("recv: %s", err)
	}

	if msg.Src != "plumb_test" {
		t.Fatal("recv src")
	}
	if msg.Dst != "edit" {
		t.Fatal("recv dst")
	}
	if msg.Wdir != "/tmp" {
		t.Fatal("recv wdir")
	}
	if msg.Type != "text" {
		t.Fatal("recv text")
	}
	if msg.Attr["foo"] != "bar baz" {
		t.Fatal("recv attr")
	}
	if !bytes.Equal(msg.Data, []byte("foo.c")) {
		t.Logf("%q %X", msg.Data, msg.Data)
		t.Fatal("recv data")
	}
}

func BenchmarkPlumbSendRecv(b *testing.B) {
	pls, err := Open("send", plan9.OWRITE)
	if err != nil {
		b.Fatal(err)
	}

	pl, err := Open("edit", plan9.OREAD)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		pls.Send(tmsg)

		if _, err := pl.Recv(); err != nil {
			b.Fatal(err)
		}
	}
}
