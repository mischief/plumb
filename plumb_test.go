package plumb

import (
	"bytes"
	"code.google.com/p/goplan9/plan9"
	"testing"
)

func TestPlumbSendRecv(t *testing.T) {
	pls, err := PlumbOpen("send", plan9.OWRITE)
	if err != nil {
		t.Fatal(err)
	}

	pl, err := PlumbOpen("edit", plan9.OREAD)
	if err != nil {
		t.Fatal(err)
	}

	pls.Send(&PlumbMsg{Src: "plumb_test", Dst: "edit", Wdir: "/tmp", Type: "text", Data: []byte("foo.c")})

	msg, err := pl.Recv()

	t.Logf("%+v", msg)

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
	if !bytes.Equal(msg.Data, []byte("foo.c")) {
		t.Logf("%q %X", msg.Data, msg.Data)
		t.Fatal("recv data")
	}
}
