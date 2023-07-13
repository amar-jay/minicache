package proto

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/amar-jay/minicache/errors"
	"github.com/amar-jay/minicache/logger"
)

type Cmd byte

const (
	Nonce Cmd = iota
	SetCmd
	GetCmd
	DelCmd
	JoinCmd
)

type Command interface {
	Bytes() ([]byte, error)
	String() (string, error)
}

type Set struct {
	Key []byte
	Val []byte
	TTL int64
}
type Get struct {
	Key []byte
}

type Join struct{}

var _ Command = (*Set)(nil)
var _ Command = (*Get)(nil)
var _ Command = (*Join)(nil)

func ParseCmd(r io.Reader) (Command, error) {
	var cmd Cmd

	if err := binary.Read(r, binary.LittleEndian, cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case SetCmd:
		return parseSetCmd(r)
	case GetCmd:
		return parseGetCmd(r)
	case JoinCmd:
		return parseJoinCmd(r)
	default:
		return nil, logger.Errorf(errors.CommandError, cmd)
	}
}

func parseSetCmd(r io.Reader) (*Set, error) {
	var keySize, valSize, ttl uint16

	if err := binary.Read(r, binary.LittleEndian, keySize); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, valSize); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, ttl); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	var key, val = make([]byte, keySize), make([]byte, valSize)

	if err := binary.Read(r, binary.LittleEndian, &key); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	if ttl == 0 {
		return nil, logger.Errorf(errors.CommandError, "ttl cannot be 0")
	}
	var cmd = Set{
		Key: key,
		Val: val,
		TTL: int64(ttl),
	}

	// TODO: find a better way to handle unset ttl
	return &cmd, nil
}

func parseGetCmd(r io.Reader) (*Get, error) {
	var keySize uint16

	if err := binary.Read(r, binary.LittleEndian, keySize); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}
	var key = make([]byte, keySize)

	if err := binary.Read(r, binary.LittleEndian, &key); err != nil {
		return nil, logger.Errorf(errors.InvalidParam, err.Error())
	}

	if keySize == 0 {
		return nil, logger.Errorf(errors.CommandError, "key cannot be empty")
	}

	return &Get{
		Key: key,
	}, nil
}

func parseJoinCmd(r io.Reader) (*Join, error) {
	logger.Panicf(errors.UnimplementedError)
	return &Join{}, nil
}

func (c *Set) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, SetCmd)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, int32(len(c.Val)))
	binary.Write(buf, binary.LittleEndian, c.TTL)
	binary.Write(buf, binary.LittleEndian, c.Key)
	binary.Write(buf, binary.LittleEndian, c.Val)

	if buf.Len() == 0 {
		return nil, logger.Errorf(errors.ParseError, "key or value cannot be empty")
	}

	return buf.Bytes(), nil
}

func (c *Get) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, GetCmd)

	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	if buf.Len() == 0 {
		return nil, logger.Errorf(errors.ParseError, "key cannot be empty")
	}
	return buf.Bytes(), nil
}

func (*Join) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, JoinCmd); err != nil {
		return nil, logger.Errorf(errors.ParseError, err.Error())
	}

	return buf.Bytes(), nil
}

func (s *Set) String() (string, error) {
	if s, err := s.Bytes(); err != nil {
		return "", err
	} else {
		return string(s), nil
	}
}

func (s *Get) String() (string, error) {
	if s, err := s.Bytes(); err != nil {
		return "", err
	} else {
		return string(s), nil
	}
}

func (s *Join) String() (string, error) {
	if s, err := s.Bytes(); err != nil {
		return "", err
	} else {
		return string(s), nil
	}
}
