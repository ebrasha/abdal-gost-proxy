/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : health_check.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Periodic health check and auto re-dial for tunnel stability.
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
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
	"github.com/ebrasha/abdal-gost-proxy/core/models"
)

// HealthChecker runs periodic checks and triggers re-dial on failure.
type HealthChecker struct {
	cfg    *models.ClientConfig
	runner *XrayClientRunner
}

// NewHealthChecker creates a health checker (does not start).
func NewHealthChecker(cfg *models.ClientConfig, runner *XrayClientRunner) *HealthChecker {
	return &HealthChecker{cfg: cfg, runner: runner}
}

// Run starts the periodic health check loop; blocks until ctx is done.
func (h *HealthChecker) Run(ctx context.Context) {
	if !h.cfg.HealthCheck.Enabled {
		return
	}
	interval := time.Duration(h.cfg.HealthCheck.IntervalSeconds) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	timeout := time.Duration(h.cfg.HealthCheck.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	maxRetries := h.cfg.HealthCheck.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	checkURL := h.cfg.HealthCheck.CheckURL
	if checkURL == "" {
		checkURL = "http://www.google.com/generate_204"
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	failCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok := h.checkOnce(ctx, timeout, checkURL)
			if ok {
				failCount = 0
				continue
			}
			failCount++
			fmt.Print(colors.Yellow(fmt.Sprintf("[Abdal Gost Proxy client] health check failed (%d/%d)\n", failCount, maxRetries)))
			if failCount >= maxRetries {
				fmt.Print(colors.Magenta("[Abdal Gost Proxy client] triggering re-dial (restart tunnel)\n"))
				if err := h.runner.Restart(); err != nil {
					fmt.Print(colors.Red(fmt.Sprintf("[Abdal Gost Proxy client] re-dial error: %v\n", err)))
				} else {
					failCount = 0
				}
			}
		}
	}
}

// checkOnce performs one health check via local SOCKS5 proxy.
func (h *HealthChecker) checkOnce(ctx context.Context, timeout time.Duration, checkURL string) bool {
	u, err := url.Parse(checkURL)
	if err != nil {
		return false
	}
	host := u.Host
	if u.Port() == "" {
		if u.Scheme == "https" {
			host = u.Hostname() + ":443"
		} else {
			host = u.Hostname() + ":80"
		}
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(timeout)
	}
	dialCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	conn, err := dialSocks5(dialCtx, "127.0.0.1", h.cfg.LocalPort, host)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// dialSocks5 connects to target via local SOCKS5 proxy (simple CONNECT).
func dialSocks5(ctx context.Context, proxyHost string, proxyPort int, targetAddr string) (net.Conn, error) {
	var d net.Dialer
	proxyAddr := net.JoinHostPort(proxyHost, strconv.Itoa(proxyPort))
	conn, err := d.DialContext(ctx, "tcp", proxyAddr)
	if err != nil {
		return nil, err
	}
	if err := socks5Connect(conn, targetAddr); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

func socks5Connect(conn net.Conn, targetAddr string) error {
	if _, err := conn.Write([]byte{5, 1, 0}); err != nil {
		return err
	}
	buf := make([]byte, 256)
	if _, err := conn.Read(buf[:2]); err != nil {
		return err
	}
	if buf[0] != 5 || buf[1] != 0 {
		return net.ErrClosed
	}
	host, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return err
	}
	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	req := []byte{5, 1, 0, 3, byte(len(host))}
	req = append(req, host...)
	req = append(req, byte(portNum>>8), byte(portNum))
	if _, err := conn.Write(req); err != nil {
		return err
	}
	n, err := conn.Read(buf)
	if err != nil || n < 4 {
		return err
	}
	if buf[1] != 0 {
		return net.ErrClosed
	}
	return nil
}
