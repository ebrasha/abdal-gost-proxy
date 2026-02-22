/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : enable_windows.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-15 12:00:00
 * Description : Enables Windows ANSI console for colored output (Windows build).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

//go:build windows

package term

import "github.com/muesli/termenv"

// EnableANSI enables virtual terminal processing on Windows so ANSI colors work in CMD/PowerShell.
// Call this at the start of main for Windows 10 compatibility.
func EnableANSI() {
	_, _ = termenv.EnableWindowsANSIConsole()
}
