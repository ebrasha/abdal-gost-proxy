/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : colors.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-15 12:00:00
 * Description : ANSI color codes and helpers for all console output (neon palette, Windows-compatible).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package colors

// ANSI Color codes for neon colors (ANSI256; compatible with Windows 10 CMD/PowerShell when ANSI is enabled).
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[38;5;196m"  // Neon Red
	ColorGreen   = "\033[38;5;46m"   // Neon Green
	ColorYellow  = "\033[38;5;226m"  // Neon Yellow
	ColorBlue    = "\033[38;5;51m"   // Neon Blue/Cyan
	ColorPurple  = "\033[38;5;129m"   // Neon Purple
	ColorPink    = "\033[38;5;201m"  // Neon Pink
	ColorOrange  = "\033[38;5;208m"  // Neon Orange
	ColorWhite   = "\033[38;5;15m"   // Bright White
	ColorCyan    = "\033[38;5;87m"   // Bright Cyan
	ColorMagenta = "\033[38;5;165m"  // Bright Magenta
)

// Red returns s wrapped with Neon Red and reset.
func Red(s string) string { return ColorRed + s + ColorReset }

// Green returns s wrapped with Neon Green and reset.
func Green(s string) string { return ColorGreen + s + ColorReset }

// Yellow returns s wrapped with Neon Yellow and reset.
func Yellow(s string) string { return ColorYellow + s + ColorReset }

// Blue returns s wrapped with Neon Blue/Cyan and reset.
func Blue(s string) string { return ColorBlue + s + ColorReset }

// Purple returns s wrapped with Neon Purple and reset.
func Purple(s string) string { return ColorPurple + s + ColorReset }

// Pink returns s wrapped with Neon Pink and reset.
func Pink(s string) string { return ColorPink + s + ColorReset }

// Orange returns s wrapped with Neon Orange and reset.
func Orange(s string) string { return ColorOrange + s + ColorReset }

// White returns s wrapped with Bright White and reset.
func White(s string) string { return ColorWhite + s + ColorReset }

// Cyan returns s wrapped with Bright Cyan and reset.
func Cyan(s string) string { return ColorCyan + s + ColorReset }

// Magenta returns s wrapped with Bright Magenta and reset.
func Magenta(s string) string { return ColorMagenta + s + ColorReset }
