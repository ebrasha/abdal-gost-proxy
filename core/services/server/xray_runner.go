/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : xray_runner.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Runs Xray-core as library with built config (VLESS, Reality, gRPC).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package server

import (
	"bytes"
	"context"

	"github.com/ebrasha/abdal-gost-proxy/core/models"
	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/xtls/xray-core/infra/conf/serial"
)

// XrayRunner holds Xray instance and config for lifecycle management.
type XrayRunner struct {
	instance *core.Instance
	config   *models.ServerConfig
}

// NewXrayRunner builds Xray config from ServerConfig and creates runner (does not start).
func NewXrayRunner(cfg *models.ServerConfig) (*XrayRunner, error) {
	jsonBytes, err := BuildXrayJSON(cfg)
	if err != nil {
		return nil, err
	}
	xrayConfig, err := serial.LoadJSONConfig(bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	instance, err := core.New(xrayConfig)
	if err != nil {
		return nil, err
	}
	return &XrayRunner{instance: instance, config: cfg}, nil
}

// Start starts the Xray instance (blocking until context is cancelled).
func (r *XrayRunner) Start(ctx context.Context) error {
	if err := r.instance.Start(); err != nil {
		return err
	}
	<-ctx.Done()
	return r.instance.Close()
}

// Close stops the Xray instance.
func (r *XrayRunner) Close() error {
	if r.instance == nil {
		return nil
	}
	return r.instance.Close()
}
