/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : fallback.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Fallback destination resolution for unauthenticated traffic (e.g. samsung.com).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package security

import (
	"fmt"
	"strconv"
)

// ResolveFallbackDest returns a "host:port" string from config dest (string or number).
func ResolveFallbackDest(dest interface{}, defaultHost string) string {
	if dest == nil {
		return defaultHost + ":443"
	}
	switch v := dest.(type) {
	case string:
		if v == "" {
			return defaultHost + ":443"
		}
		return v
	case float64:
		return defaultHost + ":" + strconv.Itoa(int(v))
	case int:
		return defaultHost + ":" + strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
