/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : main.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Minimal server entry: loads config and runs core server.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebrasha/abdal-gost-proxy/core/display"
	"github.com/ebrasha/abdal-gost-proxy/core/models"
	"github.com/ebrasha/abdal-gost-proxy/core/services/server"
	"github.com/ebrasha/abdal-gost-proxy/core/term"
)

const defaultServerConfigPath = "abdal-gost-proxy-server.json"

func main() {
	term.EnableANSI()
	display.PrintBanner("Server")
	cfgPath := defaultServerConfigPath
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("read config %s: %v", cfgPath, err)
	}
	var cfg models.ServerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := server.Run(ctx, &cfg); err != nil && ctx.Err() == nil {
		log.Fatalf("server: %v", err)
	}
}
