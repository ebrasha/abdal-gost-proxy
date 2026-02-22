/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : proxy_dialer.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Network abstraction for proxy dialing (used by client services).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package api

import (
	"context"
	"net"
)

// ProxyDialer abstracts dialing through the VLESS proxy (e.g. for health checks).
type ProxyDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}
