/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : client.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Client orchestration: SOCKS5 gate (Xray), health check, and auto re-dial.
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
	"sync"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
	"github.com/ebrasha/abdal-gost-proxy/core/models"
)

// Run starts the Abdal Gost Proxy client: Xray SOCKS5 on local_port (from config) -> VLESS+Reality+gRPC; health check and re-dial.
func Run(ctx context.Context, cfg *models.ClientConfig) error {
	runner, err := NewXrayClientRunner(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = runner.Close() }()

	var wg sync.WaitGroup
	fmt.Print(colors.Cyan(fmt.Sprintf("[Abdal Gost Proxy client] SOCKS5 on 127.0.0.1:%d -> %s:%d (VLESS+Reality+gRPC)\n", cfg.LocalPort, cfg.ServerAddr, cfg.ServerPort)))

	health := NewHealthChecker(cfg, runner)
	if cfg.HealthCheck.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			health.Run(ctx)
		}()
	}

	// Block until context cancelled.
	<-ctx.Done()
	wg.Wait()
	return nil
}
