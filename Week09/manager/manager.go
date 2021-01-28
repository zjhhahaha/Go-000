package manager

import (
	"bufio"
	"net"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
	mu     sync.Mutex
	conns  map[string]*conn
	logger log.Logger
}

func (m *Manager) Close() {
	for _, conn := range m.conns {
		conn.closeOnce.Do(func() {
			conn.handleClose(errors.New("Close cause server exit"))
		})
	}
}

func New(logger log.Logger) *Manager {
	return &Manager{
		logger: logger,
		conns:  make(map[string]*conn),
	}
}

func (m *Manager) Connect(c net.Conn) {
	reader := bufio.NewReaderSize(c, 1024)
	writer := bufio.NewWriterSize(c, 1024)
	data, _, err := reader.ReadLine()
	if err != nil {
		level.Error(m.logger).Log("msg", "connection init: read error", "err", err)
		return
	}
	name := string(data)
	conn := &conn{
		conn:      c,
		logger:    log.WithPrefix(m.logger, "name", name),
		writeChan: make(chan []byte),
	}
	m.add(name, conn)
	g := errgroup.Group{}
	g.Go(func() error {
		return conn.handleRead(reader)
	})
	g.Go(func() error {
		return conn.handleWrite(writer)
	})
	err = g.Wait()
	if err != nil {
		conn.closeOnce.Do(func() {
			conn.handleClose(err)
		})
	}
	m.mu.Lock()
	delete(m.conns, name)
	m.mu.Unlock()
}

func (m *Manager) add(name string, c *conn) {
	m.mu.Lock()
	m.conns[name] = c
	m.mu.Unlock()
}

type conn struct {
	logger    log.Logger
	id        int
	conn      net.Conn
	closeOnce sync.Once
	writeChan chan []byte
}

func (c *conn) handleWrite(writer *bufio.Writer) error {
	for msg := range c.writeChan {
		_, err := writer.Write(msg)
		if err != nil {
			return errors.Wrap(err, "write message error")
		}
		err = writer.Flush()
		if err != nil {
			return errors.Wrap(err, "flush message error")
		}
	}
	return nil
}

func (c *conn) handleRead(reader *bufio.Reader) error {
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			return errors.Wrap(err, "read message error")
		}
		level.Info(c.logger).Log("msg", "accept message", "detail", string(buf))
		c.writeChan <- []byte("accept")
	}
}

func (c *conn) handleClose(err error) {
	level.Error(c.logger).Log("msg", "connection close", "err", err)
	close(c.writeChan)
}
