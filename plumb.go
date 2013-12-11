// See http://man.cat-v.org/plan_9/6/plumb for details
package plumb

import (
	"bufio"
	"bytes"
	"code.google.com/p/goplan9/plan9"
	"code.google.com/p/goplan9/plan9/client"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once
var fsys *client.Fsys
var fsysErr error

func mountPlumb() {
	fsys, fsysErr = client.MountService("plumb")
}

// Attributes attached to a plumber message
type Attr map[string]string

// Pack attributes into a byte slice
func (pa Attr) pack() []byte {
	var space bool
	out := new(bytes.Buffer)

	for k, v := range pa {
		if space {
			out.WriteRune(' ')
		}
		fmt.Fprintf(out, "%s=%s", k, quote(v))

		space = true
	}

	return out.Bytes()
}

// A plumber message
type Msg struct {
	Src  string
	Dst  string
	Wdir string
	Type string
	Attr Attr
	Data []byte
}

// Pack a message into a byte slice
func (pm Msg) pack() []byte {
	out := new(bytes.Buffer)

	fmt.Fprintf(out, "%s\n%s\n%s\n%s\n%s\n", pm.Src, pm.Dst, pm.Wdir, pm.Type, pm.Attr.pack())
	ln := len(pm.Data)
	fmt.Fprintf(out, "%d\n%s", ln, pm.Data[:ln])

	return out.Bytes()
}

func (pm Msg) String() string {
	return fmt.Sprintf("Dst: %s Src: %s Wdir: %s Type: %s Attr: %s Ndata: %d Data: %q",
		pm.Src, pm.Dst, pm.Wdir, pm.Type, pm.Attr.pack(), len(pm.Data), pm.Data)
}

// A plumber connection
type Port client.Fid

// Open a plumber port. mode should be OREAD for ports, and OWRITE for "send".
func Open(name string, mode uint8) (*Port, error) {
	once.Do(mountPlumb)
	if fsysErr != nil {
		return nil, fsysErr
	}

	if err := fsys.Access("send", plan9.AWRITE); err != nil {
		return nil, err
	}

	if fid, err := fsys.Open(name, mode); err != nil {
		return nil, fmt.Errorf("open %s: %s", name, err)
	} else {
		return (*Port)(fid), nil
	}
}

// Read one plumber message. Plumber port must be opened with mode OREAD.
func (p *Port) Recv() (*Msg, error) {
	msg := &Msg{}
	indata := make([]byte, 8192)

	fid := (*client.Fid)(p)
	n, err := fid.Read(indata)
	if n <= 0 {
		return nil, err
	}

	buf := bytes.NewBuffer(indata)
	rd := bufio.NewReader(buf)

	msg.Src, _ = rd.ReadString('\n')
	msg.Src = strings.TrimSpace(msg.Src)

	msg.Dst, _ = rd.ReadString('\n')
	msg.Dst = strings.TrimSpace(msg.Dst)

	msg.Wdir, _ = rd.ReadString('\n')
	msg.Wdir = strings.TrimSpace(msg.Wdir)

	msg.Type, _ = rd.ReadString('\n')
	msg.Type = strings.TrimSpace(msg.Type)

	attr, _ := rd.ReadString('\n')
	msg.Attr, _ = ParseAttr(strings.TrimSpace(attr))

	ndatastr, _ := rd.ReadString('\n')
	ndata, _ := strconv.Atoi(strings.TrimSpace(ndatastr))

	data := new(bytes.Buffer)
	io.Copy(data, rd)

	msg.Data = data.Bytes()
	msg.Data = msg.Data[:ndata]

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Send a plumber message. Plumber must be opened with "send" and mode OWRITE.
func (p *Port) Send(m *Msg) error {
	msg := m.pack()

	fid := (*client.Fid)(p)
	if _, err := fid.Write(msg); err != nil {
		return err
	}

	return nil
}
