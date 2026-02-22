/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : socks5_gate.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : SOCKS5 gate: listens on local_port and forwards to Xray SOCKS5 (internal port).
 *                Acts as the user-facing gate; can be replaced by Gost for chain/re-dial features.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package client

import (
	"context"
	"io"
	"net"
	"strconv"
	"sync"
)

// SOCKS5Gate listens on localPort and relays TCP to internal SOCKS5 (Xray).
type SOCKS5Gate struct {
	localPort   int
	internalPort int
	listener   net.Listener
	mu         sync.Mutex
}

// NewSOCKS5Gate creates and starts the gate (listen on localPort, forward to internalPort).
func NewSOCKS5Gate(localPort, internalPort int) (*SOCKS5Gate, error) {
	if internalPort <= 0 {
		internalPort = InternalSocksPort
	}
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(localPort))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	g := &SOCKS5Gate{localPort: localPort, internalPort: internalPort, listener: ln}
	go g.serve()
	return g, nil
}

func (g *SOCKS5Gate) serve() {
	for {
		conn, err := g.listener.Accept()
		if err != nil {
			return
		}
		go g.handle(conn)
	}
}

func (g *SOCKS5Gate) handle(in net.Conn) {
	defer in.Close()
	target := net.JoinHostPort("127.0.0.1", strconv.Itoa(g.internalPort))
	out, err := net.Dial("tcp", target)
	if err != nil {
		return
	}
	defer out.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { io.Copy(out, in); wg.Done() }()
	go func() { io.Copy(in, out); wg.Done() }()
	wg.Wait()
}

// Run blocks until ctx is done; then closes the gate.
func (g *SOCKS5Gate) Run(ctx context.Context) error {
	<-ctx.Done()
	return g.Close()
}

// Close stops the listener.
func (g *SOCKS5Gate) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.listener == nil {
		return nil
	}
	err := g.listener.Close()
	g.listener = nil
	return err
}
