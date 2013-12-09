package plumb

import (
	"testing"
)

func TestPlumbRecv(t *testing.T) {
	pl, err := PlumbOpen("test", 0)
	if err != nil {
		t.Fatal(err)
	}

	for {
		msg, err := pl.Recv()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%+v", msg)

		if msg.Attr["exit"] == "true" {
			break
		}
	}

}
