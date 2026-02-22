/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : xray_config_builder.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Builds Xray-core JSON config from Abdal server config (VLESS, gRPC, Reality, Fallback).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package server

import (
	"bytes"
	"encoding/json"

	"github.com/ebrasha/abdal-gost-proxy/core/models"
	"github.com/ebrasha/abdal-gost-proxy/core/security"
)

// xrayInbound represents one inbound in Xray JSON format.
type xrayInbound struct {
	Listen   string          `json:"listen"`
	Port     int             `json:"port"`
	Protocol string          `json:"protocol"`
	Settings *xrayVLESSSet    `json:"settings"`
	StreamSettings *xrayStream `json:"streamSettings"`
	Sniffing *xraySniffing   `json:"sniffing,omitempty"`
}

type xrayVLESSSet struct {
	Clients      []xrayClient `json:"clients"`
	Decryption   string       `json:"decryption"`
}

type xrayClient struct {
	ID    string `json:"id"`
	Email string `json:"email,omitempty"`
	Flow  string `json:"flow,omitempty"`
}

type xrayStream struct {
	Network   string            `json:"network"`
	Security  string            `json:"security"`
	RealitySettings *xrayReality `json:"realitySettings,omitempty"`
	GRPCSettings   *xrayGRPC   `json:"grpcSettings,omitempty"`
}

type xrayReality struct {
	Show        bool     `json:"show"`
	Dest        string   `json:"dest"`
	Xver        int      `json:"xver"`
	ServerNames []string `json:"serverNames"`
	PrivateKey  string   `json:"privateKey"`
	ShortIDs    []string `json:"shortIds"`
}

type xrayGRPC struct {
	ServiceName string `json:"serviceName"`
	MultiMode   bool   `json:"multiMode"`
}

type xraySniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

type xrayOutbound struct {
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
}

type xrayConfig struct {
	Log       *xrayLog       `json:"log,omitempty"`
	Inbounds  []xrayInbound  `json:"inbounds"`
	Outbounds []xrayOutbound `json:"outbounds"`
}

type xrayLog struct {
	Loglevel string `json:"loglevel"`
}

// BuildXrayJSON converts ServerConfig to Xray-compatible JSON (VLESS + gRPC + Reality).
func BuildXrayJSON(cfg *models.ServerConfig) ([]byte, error) {
	realityParams := security.FromServerReality(&cfg.RealitySettings)
	defaultHost := cfg.Fallback.DefaultHost
	if defaultHost == "" {
		defaultHost = "www.samsung.com"
	}
	fallbackDest := security.ResolveFallbackDest(cfg.Fallback.Dest, defaultHost)
	if realityParams.Dest != "" {
		fallbackDest = realityParams.Dest
	}

	protocol := cfg.Protocol
	if protocol == "" {
		protocol = "vless"
	}
	network := cfg.Transport.Type
	if network == "" {
		network = "grpc"
	}

	// For gRPC transport, flow must be empty (XTLS/Vision only for direct TCP+TLS/Reality)
	clients := make([]xrayClient, 0, len(cfg.Users))
	for _, u := range cfg.Users {
		flow := u.Flow
		if network == "grpc" {
			flow = ""
		}
		clients = append(clients, xrayClient{ID: u.ID, Email: u.Email, Flow: flow})
	}

	serviceName := cfg.Transport.ServiceName
	if serviceName == "" {
		serviceName = "abdal-grpc-stream"
	}

	streamSettings := &xrayStream{
		Network:  network,
		Security: "reality",
		RealitySettings: &xrayReality{
			Show:        false,
			Dest:        fallbackDest,
			Xver:        cfg.Fallback.XVer,
			ServerNames: realityParams.ServerNames,
			PrivateKey:  realityParams.PrivateKey,
			ShortIDs:    realityParams.ShortIDs,
		},
		GRPCSettings: nil,
	}
	if network == "grpc" {
		streamSettings.GRPCSettings = &xrayGRPC{
			ServiceName: serviceName,
			MultiMode:   cfg.Transport.MultiMode,
		}
	}

	xcfg := xrayConfig{
		Log: &xrayLog{Loglevel: "warning"},
		Inbounds: []xrayInbound{{
			Listen:   cfg.ListenAddress,
			Port:     cfg.ListenPort,
			Protocol: protocol,
			Settings: &xrayVLESSSet{
				Clients:    clients,
				Decryption: "none",
			},
			StreamSettings: streamSettings,
			Sniffing: &xraySniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls", "quic"},
			},
		}},
		Outbounds: []xrayOutbound{
			{Protocol: "freedom", Tag: "direct"},
			{Protocol: "blackhole", Tag: "block"},
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
