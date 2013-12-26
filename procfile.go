package main

import (
	"bufio"
	"errors"
	"io"
)

var (
	InvalidProcessName = errors.New("Invalid characters for process type")
)

type Procfile struct {
	Entries []ProcfileEntry
	data    io.Reader
}

type ProcfileEntry struct {
	Type    string
	Command string
}

func NewProcfile(r io.Reader) *Procfile {
	procfile := new(Procfile)
	procfile.data = r
	return procfile
}

func ParseProcfile(r io.Reader) (*Procfile, error) {
	procfile := NewProcfile(r)
	err := procfile.Parse()
	return procfile, err
}

func (p *Procfile) Parse() error {
	r := bufio.NewReader(p.data)

  var ƒ parseStateFunc
  ƒ = new(startLine)
	var k []byte
	var v []byte

	for {
		c, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				addEntry(&k, &v, p)
				return nil
			}
			return err
		}
    ƒ = ƒ.parse(c, &k, &v, p)
	}
}

type parseStateFunc interface {
  parse(c byte, k *[]byte, v *[]byte, p *Procfile) parseStateFunc
}

type startLine struct {}
type readType struct {}
type endType struct {}
type readProc struct {}

func (s *startLine) parse(c byte, k *[]byte, v *[]byte, p *Procfile) parseStateFunc {
  switch {
  case isWhitespace(c):
    return new(startLine)
  case isValid(c):
    *k = append(*k, c)
    return new(readType)
  default:
    panic(InvalidProcessName)
  }
}

func (s *readType) parse(c byte, k *[]byte, v *[]byte, p *Procfile) parseStateFunc {
  switch {
  case isValid(c):
    *k = append(*k, c)
    return new(readType)
  case c == ':':
    return new(readProc)
  case isIgnored(c):
    return new(endType)
  default:
    panic(InvalidProcessName)
  }
}

func (s *endType) parse(c byte, k *[]byte, v *[]byte, p *Procfile) parseStateFunc {
  switch {
  case isIgnored(c):
    return new(endType)
  case c == ':':
    return new(readProc)
  default:
    panic(InvalidProcessName)
  }
}

func (s *readProc) parse(c byte, k *[]byte, v *[]byte, p *Procfile) parseStateFunc {
  switch c {
  case '\n':
    addEntry(k, v, p)
    return new(startLine)
  default:
    *v = append(*v, c)
    return new(readProc)
  }
}

func isWhitespace(c byte) bool {
  return c == ' ' || c == '\t' || c == '\n'
}

func isIgnored(c byte) bool {
	return c == ' ' || c == '\t'
}

func isValid(c byte) bool {
	return (c > '0' && c < '9') ||
		(c > 'A' && c < 'Z') ||
		(c > 'a' && c < 'z') ||
		(c == '_')
}

func addEntry(k *[]byte, v *[]byte, pf *Procfile) {
	if len(*k) > 0 && len(*v) > 0 {
		pe := &ProcfileEntry{
			Type:    string(*k),
			Command: string(*v),
		}
		pf.Entries = append(pf.Entries, *pe)
		*k = nil
		*v = nil
	}
}
