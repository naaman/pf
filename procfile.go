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

	var k []byte
	var v []byte
	state := 0

	for {
		c, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				addEntry(&k, &v, p)
				return nil
			}
			return err
		}

		switch c {
		case '\n':
			addEntry(&k, &v, p)
			state = 0
		case ':':
			state = 1
		default:
			switch state {
			case 0:
				if k == nil && isIgnored(c) {
					break
				}
				if isValid(c) {
					k = append(k, c)
				} else {
					return InvalidProcessName
				}
			case 1:
				v = append(v, c)
			}
		}
	}
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
