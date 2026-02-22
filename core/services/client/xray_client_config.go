/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : xray_client_config.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Builds Xray-core client JSON config (SOCKS5 inbound + VLESS Reality gRPC outbound).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package client

import (
	"bytes"
	"encoding/json"

	"github.com/ebrasha/abdal-gost-proxy/core/models"
)

// clientInbound represents SOCKS5 inbound in Xray client config.
type clientInbound struct {
	Listen   string `json:"listen"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

// clientOutboundVless represents VLESS outbound in Xray client config.
type clientOutboundVless struct {
	Protocol string              `json:"protocol"`
	Tag      string              `json:"tag"`
	Settings *clientVlessSettings `json:"settings"`
	StreamSettings *clientStream `json:"streamSettings"`
}

type clientVlessSettings struct {
	Vnext []clientVnext `json:"vnext"`
}

type clientVnext struct {
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Users   []clientUser  `json:"users"`
}

type clientUser struct {
	ID         string `json:"id"`
	Encryption string `json:"encryption"`
	Flow       string `json:"flow,omitempty"`
}

type clientStream struct {
	Network   string           `json:"network"`
	Security  string           `json:"security"`
	RealitySettings *clientReality `json:"realitySettings,omitempty"`
	GRPCSettings   *clientGRPC   `json:"grpcSettings,omitempty"`
}

type clientReality struct {
	Show          bool   `json:"show"`
	Fingerprint   string `json:"fingerprint"`
	ServerName    string `json:"serverName"`
	PublicKey     string `json:"publicKey"`
	ShortID       string `json:"shortId"`
}

type clientGRPC struct {
	ServiceName string `json:"serviceName"`
}

type clientConfig struct {
	Log       *xrayLogClient   `json:"log,omitempty"`
	Inbounds  []clientInbound   `json:"inbounds"`
	Outbounds []interface{}    `json:"outbounds"`
}

type xrayLogClient struct {
	Loglevel string `json:"loglevel"`
}

// InternalSocksPort is used only as fallback when localPort from config is missing or <= 0.
const InternalSocksPort = 10809

// BuildXrayClientJSON produces Xray client config JSON (SOCKS5 inbound on localPort from config, VLESS+Reality+gRPC out).
func BuildXrayClientJSON(cfg *models.ClientConfig, localPort int) ([]byte, error) {
	if localPort <= 0 {
		localPort = InternalSocksPort
	}
	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "abdal-grpc-stream"
	}
	fingerprint := cfg.Fingerprint
	if fingerprint == "" {
		fingerprint = "chrome"
	}
	network := cfg.Transport
	if network == "" {
		network = "grpc"
	}

	streamSettings := &clientStream{
		Network:  network,
		Security: "reality",
		RealitySettings: &clientReality{
			Show:        false,
			Fingerprint: fingerprint,
			ServerName:  cfg.SNI,
			PublicKey:   cfg.RealityPublicKey,
			ShortID:     cfg.ShortID,
		},
		GRPCSettings: nil,
	}
	if network == "grpc" {
		streamSettings.GRPCSettings = &clientGRPC{ServiceName: serviceName}
	}

	xcfg := clientConfig{
		Log: &xrayLogClient{Loglevel: "warning"},
		Inbounds: []clientInbound{{
			Listen:   "127.0.0.1",
			Port:     localPort,
			Protocol: "socks",
			Tag:      "socks-in",
		}},
		Outbounds: []interface{}{
			clientOutboundVless{
				Protocol: "vless",
				Tag:      "proxy",
				Settings: &clientVlessSettings{
					Vnext: []clientVnext{{
						Address: cfg.ServerAddr,
						Port:    cfg.ServerPort,
						Users: []clientUser{{
							ID:         cfg.UUID,
							Encryption: "none",
							Flow:       "", // empty for gRPC; "xtls-rprx-vision" only for direct TCP+TLS/Reality
						}},
					}},
				},
				StreamSettings: streamSettings,
			},
			struct {
				Protocol string `json:"protocol"`
				Tag      string `json:"tag"`
			}{Protocol: "freedom", Tag: "direct"},
			struct {
				Protocol string `json:"protocol"`
				Tag      string `json:"tag"`
			}{Protocol: "blackhole", Tag: "block"},
		},
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(xcfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
