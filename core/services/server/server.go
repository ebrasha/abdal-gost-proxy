/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : server.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Server orchestration: loads config, starts Xray (VLESS+Reality+gRPC).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package server

import (
	"context"
	"fmt"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
	"github.com/ebrasha/abdal-gost-proxy/core/models"
)

// Run starts the Abdal Gost Proxy server (VLESS + Reality + gRPC on listen_port).
func Run(ctx context.Context, cfg *models.ServerConfig) error {
	runner, err := NewXrayRunner(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := runner.Close(); err != nil {
			fmt.Print(colors.Magenta(fmt.Sprintf("[Abdal Gost Proxy server] close: %v\n", err)))
		}
	}()
	fmt.Print(colors.Cyan(fmt.Sprintf("[Abdal Gost Proxy server] listening on %s:%d (VLESS+gRPC+Reality)\n", cfg.ListenAddress, cfg.ListenPort)))
	return runner.Start(ctx)
}
