/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : server_config.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Server-side configuration models for VLESS, Reality, and Gost.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package models

// ServerUser represents a VLESS client (user) on the server.
type ServerUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Flow  string `json:"flow"`
}

// RealitySettings holds XTLS-Reality server configuration.
type RealitySettings struct {
	Enabled     bool     `json:"enabled"`
	Dest        string   `json:"dest"`
	ServerNames []string `json:"server_names"`
	PrivateKey  string   `json:"private_key"`
	ShortIDs    []string `json:"short_ids"`
}

// TransportConfig defines gRPC transport options.
type TransportConfig struct {
	Type        string `json:"type"`
	ServiceName string `json:"service_name"`
	MultiMode   bool   `json:"multi_mode"`
}

// FallbackConfig defines fallback destination for unauthenticated traffic (e.g. host:port or port).
type FallbackConfig struct {
	Dest       interface{} `json:"dest"`         // string "host:443" or number 80
	XVer       int         `json:"xver"`
	DefaultHost string     `json:"default_host"` // when dest is a port number, host to use (optional)
}

// GostConfig holds Gost-specific options for chaining and limits.
type GostConfig struct {
	EnableChaining  bool `json:"enable_chaining"`
	MaxConnections  int  `json:"max_connections"`
}

// ServerConfig is the root server configuration loaded from abdal-gost-proxy-server.json.
type ServerConfig struct {
	ListenAddress   string           `json:"listen_address"`
	ListenPort      int              `json:"listen_port"`
	Protocol        string           `json:"protocol"`
	Users           []ServerUser     `json:"users"`
	RealitySettings RealitySettings  `json:"reality_settings"`
	Transport       TransportConfig  `json:"transport"`
	Fallback        FallbackConfig   `json:"fallback"`
	GostConfig      GostConfig       `json:"gost_config"`
}
