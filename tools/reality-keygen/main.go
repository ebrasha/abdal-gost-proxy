/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : main.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:30:03
 * Description : CLI tool to generate Reality public/private key pair for Abdal Gost Proxy.
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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
	"github.com/ebrasha/abdal-gost-proxy/core/term"
)

func main() {
	term.EnableANSI()
	printBanner()
	reader := bufio.NewReader(os.Stdin)
	for {
		uuidStr, err := GenerateUserUUID()
		if err != nil {
			fmt.Fprint(os.Stderr, colors.Red(fmt.Sprintf("error: %v\n", err)))
			os.Exit(1)
		}
		privateKey, publicKey, err := GenerateRealityKeys()
		if err != nil {
			fmt.Fprint(os.Stderr, colors.Red(fmt.Sprintf("error: %v\n", err)))
			os.Exit(1)
		}
		shortIDs, err := GenerateShortIDs(2)
		if err != nil {
			fmt.Fprint(os.Stderr, colors.Red(fmt.Sprintf("error: %v\n", err)))
			os.Exit(1)
		}
		fmt.Println(colors.Cyan("uuid (use for VLESS user id on server and client):"))
		fmt.Println(colors.Magenta(uuidStr))
		fmt.Println()
		fmt.Println(colors.Cyan("private_key (use on server, keep secret):"))
		fmt.Println(colors.Magenta(privateKey))
		fmt.Println()
		fmt.Println(colors.Cyan("reality_public_key (use on client):"))
		fmt.Println(colors.Magenta(publicKey))
		fmt.Println()
		fmt.Println(colors.Cyan("short_ids (array of hex 2–16 chars; server: short_ids, client: pick one as short_id):"))
		for _, s := range shortIDs {
			fmt.Println(colors.Magenta(s))
		}
		arr, _ := json.Marshal(shortIDs)
		fmt.Println(colors.Cyan("  Server JSON short_ids:"))
		fmt.Println(colors.Magenta("  " + string(arr)))
		fmt.Println()
		fmt.Print(colors.Yellow("Generate again? (y/n): "))
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "n" || line == "no" || line == "نه" {
			fmt.Println(colors.Green("Bye."))
			break
		}
		if line == "y" || line == "yes" || line == "آره" || line == "بله" {
			continue
		}
		break
	}
}

// printBanner shows software name and programmer info (no box), neon color style.
func printBanner() {
	fmt.Println("  " + colors.Cyan("Abdal Gost Proxy — Reality Keygen"))
	fmt.Println("  " + colors.Magenta("UUID & X25519 key pair for VLESS/Reality"))
	fmt.Println()
	fmt.Println("  " + colors.Green("Programmer: Ebrahim Shafiei (EbraSha)"))
	fmt.Println("  " + colors.Cyan("Email: Prof.Shafiei@Gmail.com"))
	fmt.Println(colors.Magenta("  ─────────────────────────────────────────"))
	fmt.Println()
}
