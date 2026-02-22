/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : banner.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-15 03:21:57
 * Description : Startup banner for server and client (phosphorescent/cyberpunk style).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package display

import (
	"fmt"

	"github.com/ebrasha/abdal-gost-proxy/core/colors"
)

const (
	appName    = "Abdal Gost Proxy"
	programmer = "Ebrahim Shafiei (EbraSha)"
	email      = "Prof.Shafiei@Gmail.com"
	separator  = "  ─────────────────────────────────────────"
)

// PrintBanner prints app name, mode (Server/Client), programmer info and separator. No box.
func PrintBanner(mode string) {
	fmt.Print("  " + colors.Cyan(appName))
	fmt.Println(colors.Magenta(" — " + mode))
	fmt.Println()
	fmt.Println("  " + colors.Green("Programmer: "+programmer))
	fmt.Println("  " + colors.Cyan("Email: "+email))
	fmt.Println(colors.Magenta(separator))
	fmt.Println()
}
