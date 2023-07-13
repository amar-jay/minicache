package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/amar-jay/minicache/cache"
	"github.com/amar-jay/minicache/client"
	"github.com/amar-jay/minicache/errors"
	"github.com/amar-jay/minicache/logger"
	"github.com/amar-jay/minicache/proto"
	"github.com/urfave/cli/v2"
)

// this is the storage of the cache
type Members struct {
	clients map[*client.Client]any // all clients connected to server
	sync.Mutex
}
type Server struct {
	Members
	ListenAddr string
	cache      cache.Cache
}

// new server with urfave cli context
func NewWithCtx(c *cli.Context) error {
	s := &Server{
		ListenAddr: c.String("listen"),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		logger.Printf("exit requested, shutting down signal: %v", sig)
		s.shutdown(1)
	}()

	return s.start()

}

func (s *Server) start() (err error) {
	l, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return logger.Errorf(errors.ServerError, err.Error())
	}

	logger.Printf("mini-cacher server on: %s\n", s.ListenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Warnf(errors.ServerError, err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) shutdown(code int) error {

	os.Exit(code)
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, err := proto.ParseCmd(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Errorf(errors.ServerError, err.Error())
			break
		}

		go func(conn net.Conn, cmd proto.Command) {
			switch t := cmd.(type) {
			case *proto.Set:
				s.handleSet(conn, t)
			case *proto.Get:
				s.handleGet(conn, t)
			case *proto.Join:
				s.handleJoin(conn, t)
			default:
				logger.Errorf(errors.CommandError, t)
			}
		}(conn, cmd)

	}
}

func (s *Server) handleSet(conn net.Conn, cmd *proto.Set) error {
	logger.Printf("SET %s to %s", cmd.Key, cmd.Val)

	go func() {
		for member := range s.clients {
			err := member.Set(context.TODO(), cmd.Key, cmd.Val, cmd.TTL)
			if err != nil {
				logger.Println("forward to member error:", err)
			}
		}
	}()

	resp := proto.Response{}
	if err := s.cache.Set(cmd.Key, cmd.Val, time.Duration(cmd.TTL)); err != nil {
		resp.Status = http.StatusInternalServerError
		b, _ := resp.Bytes()
		_, err := conn.Write(b)
		return err
	}

	resp.Status = http.StatusOK
	b, _ := resp.Bytes()
	_, err := conn.Write(b)

	return err
}

func (s *Server) handleGet(conn net.Conn, cmd *proto.Get) (err error) {
	logger.Println("get command received: ", conn.RemoteAddr().String())
	res := new(proto.Response)
	res.Key = cmd.Key

	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		res.Status = http.StatusNotFound
		res.Value = []byte{}
		b, _ := res.Bytes()
		if _, err = conn.Write(b); err != nil {
			return logger.Errorf(errors.ServerError, err.Error())
		}
		return
	}
	res.Status = http.StatusOK
	res.Value = val

	b, _ := res.Bytes()
	_, err = conn.Write(b)
	return err

}

func (s *Server) handleJoin(conn net.Conn, cmd *proto.Join) {
	logger.Println("new member joined: ", conn.RemoteAddr().String())
	s.Lock()
	defer s.Unlock()
	s.clients[client.NewWithConn(conn)] = struct{}{}
}
