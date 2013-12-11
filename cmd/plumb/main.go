package main

import (
	"bytes"
	"code.google.com/p/goplan9/plan9"
	"flag"
	"fmt"
	"github.com/mischief/plumb"
	"io"
	"os"
)

var p = flag.String("p", "send", "plumbfile")
var a = flag.String("a", "", "attributes")
var s = flag.String("s", "goplumb", "source")
var d = flag.String("d", "", "destination port")
var t = flag.String("t", "text", "type")
var w = flag.String("w", "", "working directory")
var i = flag.Bool("i", false, "read from stdin")

func main() {
	flag.Parse()

	port, err := plumb.Open(*p, plan9.OWRITE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg := &plumb.Msg{}

	msg.Src = *s
	msg.Dst = *d
	msg.Type = *t
	msg.Attr, _ = plumb.ParseAttr(*a)

	if *w == "" {
		wd, _ := os.Getwd()
		msg.Wdir = wd
	} else {
		msg.Wdir = *w
	}

	if *i == true {
		if _, ok := msg.Attr["action"]; !ok {
			msg.Attr["action"] = "showdata"
		}

		data := new(bytes.Buffer)
		io.Copy(data, os.Stdin)
		msg.Data = data.Bytes()
	} else {
		for _, s := range flag.Args() {
			msg.Data = []byte(s)
			port.Send(msg)
		}
		os.Exit(0)
	}

	port.Send(msg)
}
