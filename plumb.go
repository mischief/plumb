// See http://man.cat-v.org/plan_9/6/plumb for details
package plumb

import (
	"code.google.com/p/goplan9/plan9/client"
  "fmt"
  "bytes"
  "bufio"
  "strconv"
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
	fid  *client.Fid
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
		/* try create */
    return nil, fmt.Errorf("open %s: %s", name, err)
		p.fid, err = p.fsys.Create(name, mode, 0600)
		if err != nil {
      return nil, fmt.Errorf("create %s: %s", name, err)
		}
	}

	return p, nil
}

// Read one plumber message
func (p *Plumber) Recv() (*PlumbMsg, error) {
	msg := &PlumbMsg{}
	data := make([]byte, 8192)

	n, err := p.fid.Read(data)
  if n <= 0 {
    return nil, err
  }

  buf := bytes.NewBuffer(data)
  scan := bufio.NewScanner(buf)

  if scan.Scan() {
    msg.Src = scan.Text()
  }

  if scan.Scan() {
    msg.Dst = scan.Text()
  }

  if scan.Scan() {
    msg.Wdir = scan.Text()
  }

  if scan.Scan() {
    msg.Type = scan.Text()
  }

  if scan.Scan() {
    attr := scan.Text()
    msg.Attr, _ = ParseAttr(attr)
  }

  if scan.Scan() {
    ndata := scan.Text()
    msg.Ndata, _ = strconv.Atoi(ndata)
  }

  for scan.Scan() {
    msg.Data += scan.Text() + "\n"
  }

  if err = scan.Err(); err != nil {
    return nil, err
  }

	return msg, nil
}
