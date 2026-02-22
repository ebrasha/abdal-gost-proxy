/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : xray_client_runner.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Runs Xray-core client (SOCKS5 inbound + VLESS Reality gRPC outbound); supports restart for re-dial.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package client

import (
	"bytes"
	"context"
	"sync"

	"github.com/ebrasha/abdal-gost-proxy/core/models"
	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/xtls/xray-core/infra/conf/serial"
)

// XrayClientRunner runs Xray client and supports restart (re-dial).
type XrayClientRunner struct {
	config   *models.ClientConfig
	instance *core.Instance
	mu       sync.Mutex
}

// NewXrayClientRunner creates runner and starts Xray client (SOCKS5 on local_port).
func NewXrayClientRunner(cfg *models.ClientConfig) (*XrayClientRunner, error) {
	r := &XrayClientRunner{config: cfg}
	return r, r.start()
}

func (r *XrayClientRunner) start() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Use config's local_port so Xray listens on the user-chosen port (no separate gate).
	jsonBytes, err := BuildXrayClientJSON(r.config, r.config.LocalPort)
	if err != nil {
		return err
	}
	xrayConfig, err := serial.LoadJSONConfig(bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	instance, err := core.New(xrayConfig)
	if err != nil {
		return err
	}
	if err := instance.Start(); err != nil {
		return err
	}
	r.instance = instance
	return nil
}

// Restart closes current instance and starts a new one (re-dial).
func (r *XrayClientRunner) Restart() error {
	r.mu.Lock()
	if r.instance != nil {
		_ = r.instance.Close()
		r.instance = nil
	}
	r.mu.Unlock()
	return r.start()
}

// Close stops the Xray client instance.
func (r *XrayClientRunner) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.instance == nil {
		return nil
	}
	err := r.instance.Close()
	r.instance = nil
	return err
}

// Run keeps the runner alive until context is done.
func (r *XrayClientRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return r.Close()
}
