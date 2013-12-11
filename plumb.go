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
)

// A plumber message
type PlumbMsg struct {
	Src   string
	Dst   string
	Wdir  string
	Type  string
	Attr  PlumbAttr
	Ndata int
	Data  []byte
}

// Pack a message into a byte slice
func (pm PlumbMsg) Pack() []byte {
	out := new(bytes.Buffer)

	fmt.Fprintf(out, "%s\n", pm.Src)
	fmt.Fprintf(out, "%s\n", pm.Dst)
	fmt.Fprintf(out, "%s\n", pm.Wdir)
	fmt.Fprintf(out, "%s\n", pm.Type)
	fmt.Fprintf(out, "%s\n", pm.Attr.Pack())

	if pm.Ndata <= 0 {
		pm.Ndata = len(pm.Data)
	}
	fmt.Fprintf(out, "%d\n", pm.Ndata)

	fmt.Fprintf(out, "%s", pm.Data[:pm.Ndata])
	return out.Bytes()
}

// Attributes attached to a plumber message
type PlumbAttr map[string]string

// Pack attributes into a byte slice
func (pa PlumbAttr) Pack() []byte {
	var space bool
	out := new(bytes.Buffer)

	for k, v := range pa {
		if space {
			out.WriteRune(' ')
		}
		fmt.Fprintf(out, "%s=%s", k, quote(v))
	}

	return out.Bytes()
}

// A plumber connection
type Plumber struct {
	// fs mount
	fsys *client.Fsys
	// read fid
	fid *client.Fid
}

// Open a plumber port. mode should be OREAD for ports, and OWRITE for "send".
func PlumbOpen(name string, mode uint8) (*Plumber, error) {
	var err error
	p := &Plumber{}

	p.fsys, err = client.MountService("plumb")
	if err != nil {
		return nil, fmt.Errorf("mount plumb: %s", err)
	}

	if err := p.fsys.Access("send", plan9.AWRITE); err != nil {
		return nil, err
	}

	p.fid, err = p.fsys.Open(name, mode)
	if err != nil {
		return nil, fmt.Errorf("open %s: %s", name, err)

		/* try create */
		/*
			p.fid, err = p.fsys.Create(name, mode, 0600)
			if err != nil {
				return nil, fmt.Errorf("create %s: %s", name, err)
			}
		*/
	}

	return p, nil
}

// Read one plumber message. Plumber port must be opened with mode OREAD.
func (p *Plumber) Recv() (*PlumbMsg, error) {
	msg := &PlumbMsg{}
	indata := make([]byte, 8192)

	n, err := p.fid.Read(indata)
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

	ndata, _ := rd.ReadString('\n')
	msg.Ndata, _ = strconv.Atoi(strings.TrimSpace(ndata))

	data := new(bytes.Buffer)
	io.Copy(data, rd)

	msg.Data = data.Bytes()
	msg.Data = msg.Data[:msg.Ndata]

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Send a plumber message. Plumber must be opened with "send" and mode OWRITE.
func (p *Plumber) Send(m *PlumbMsg) error {
	msg := m.Pack()

	if _, err := p.fid.Write(msg); err != nil {
		return err
	}

	return nil
}
