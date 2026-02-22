/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : client_main.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Minimal client entry: loads config and runs core client (SOCKS5 + health check + re-dial).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
	"github.com/ebrasha/abdal-gost-proxy/core/display"
	"github.com/ebrasha/abdal-gost-proxy/core/models"
	"github.com/ebrasha/abdal-gost-proxy/core/services/client"
	"github.com/ebrasha/abdal-gost-proxy/core/term"
)

func main() {
	term.EnableANSI()
	display.PrintBanner("Client")
	cfgPath := chooseClientConfig()
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("read config %s: %v", cfgPath, err)
	}
	var cfg models.ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := client.Run(ctx, &cfg); err != nil && ctx.Err() == nil {
		log.Fatalf("client: %v", err)
	}
}

// chooseClientConfig returns absolute path to the JSON config: from args, or auto if single .json, or prompt if multiple.
// Config files are listed and read from the executable's directory so the selected file is always the one next to the exe.
func chooseClientConfig() string {
	if len(os.Args) > 1 {
		argPath := os.Args[1]
		if filepath.IsAbs(argPath) {
			return argPath
		}
		if abs, err := filepath.Abs(argPath); err == nil {
			return abs
		}
		return argPath
	}
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	dir := filepath.Dir(exePath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("list config dir %s: %v", dir, err)
	}
	var jsonFiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(strings.ToLower(name), ".json") {
			jsonFiles = append(jsonFiles, name)
		}
	}
	sort.Strings(jsonFiles)
	switch len(jsonFiles) {
	case 0:
		log.Fatalf("no .json config file found in %s", dir)
	case 1:
		cfgPath := filepath.Join(dir, jsonFiles[0])
		fmt.Print(colors.Green("Using profile: " + strings.TrimSuffix(jsonFiles[0], filepath.Ext(jsonFiles[0])) + "\n\n"))
		return cfgPath
	}
	fmt.Println(colors.Cyan("Available profiles (.json):"))
	for i, name := range jsonFiles {
		profileName := strings.TrimSuffix(name, filepath.Ext(name))
		fmt.Printf("  %s\n", colors.Magenta(fmt.Sprintf("%d. %s", i+1, profileName)))
	}
	fmt.Print(colors.Yellow("Select profile (number or name): "))
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("read input: %v", err)
	}
	choice := strings.TrimSpace(line)
	if choice == "" {
		log.Fatalf("no profile selected")
	}
	if num, err := strconv.Atoi(choice); err == nil && num >= 1 && num <= len(jsonFiles) {
		return filepath.Join(dir, jsonFiles[num-1])
	}
	choiceNoExt := strings.TrimSuffix(choice, ".json")
	for _, name := range jsonFiles {
		profileName := strings.TrimSuffix(name, filepath.Ext(name))
		if profileName == choice || profileName == choiceNoExt || name == choice {
			return filepath.Join(dir, name)
		}
	}
	log.Fatalf("invalid profile: %q", choice)
	return ""
}
