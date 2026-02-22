/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : enable_other.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-15 12:00:00
 * Description : No-op for ANSI enable on non-Windows (Unix already supports ANSI).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

//go:build !windows

package term

// EnableANSI is a no-op on non-Windows; ANSI is normally supported.
func EnableANSI() {}
