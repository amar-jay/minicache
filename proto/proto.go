package proto

import (
	"encoding/binary"
	"io"
)

type Cmd struct{}
type SetCmd string
type GetCmd string
type JoinCmd string

func ParseCmd[T JoinCmd | SetCmd | GetCmd](r io.Reader) (T, error) {
	cmd := Cmd{}

	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
	}
}
