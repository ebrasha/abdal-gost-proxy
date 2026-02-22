/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : reality.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : XTLS-Reality handshake and settings handling for server inbound.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package security

import "github.com/ebrasha/abdal-gost-proxy/core/models"

// RealityParams holds validated Reality parameters for building Xray config.
type RealityParams struct {
	Dest       string
	ServerName string
	ServerNames []string
	PrivateKey string
	ShortIDs   []string
}

// FromServerReality builds RealityParams from server config.
func FromServerReality(rs *models.RealitySettings) RealityParams {
	p := RealityParams{
		Dest:        rs.Dest,
		PrivateKey:  rs.PrivateKey,
		ShortIDs:    append([]string{}, rs.ShortIDs...),
		ServerNames: append([]string{}, rs.ServerNames...),
	}
	if len(p.ServerNames) > 0 {
		p.ServerName = p.ServerNames[0]
	}
	return p
}
