package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/amar-jay/minicache/cache"
	"github.com/amar-jay/minicache/client"
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
	cacher     cache.Cache
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
		log.Printf("exit requested, shutting down", "signal: ", sig)
		s.shutdown(1)
	}()

	return s.start()

}

func (s *Server) start() (err error) {
	l, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return logger.Errorf(ServerError + err.Error())
	}

	logger.Printf("mini-cacher server on: %s\n", s.ListenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Warnf(ServerError, err.Error())
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
			logger.Errorf(ServerError, "parsing error, ", err.Error())
			break
		}

		go func() {
			switch t := cmd.(type) {
			case *proto.SetCmd:
				s.handleSet(conn, t)
			case *proto.GetCmd:
				s.handleGet(conn, t)
			case *proto.JoinCmd:
				log.Errorf(CommandError, t)
			default:
				logger.Errorf(CommandError, t)
			}
		}(conn, cmd)

	}
}

func (s *Server) handleSet(conn net.Conn, cmd *proto.SetCmd) {
	logger.Println("new member joined: ", conn.RemoteAddr().String())
	s.Lock()
	defer s.Unlock()
	s.clients[client.NewWithConn(conn)] = struct{}{}
}
