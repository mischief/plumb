// See http://man.cat-v.org/plan_9/6/plumb for details
package plumb

import (
	"bufio"
	"bytes"
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
	Data  string
}

// Attributes attached to a plumber message
type PlumbAttr map[string]string

// A plumber connection
type Plumber struct {
	// fs mount
	fsys *client.Fsys
	// read fid
	fid *client.Fid
}

// Open a plumber connection to the plumber. mode should be
//0 to read plumber messages.
func PlumbOpen(name string, mode uint8) (*Plumber, error) {
	var err error
	p := &Plumber{}

	p.fsys, err = client.MountService("plumb")
	if err != nil {
		return nil, fmt.Errorf("mount plumb: %s", err)
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

// Read one plumber message
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

	msg.Data = data.String()

	if err != nil {
		return nil, err
	}

	return msg, nil
}
