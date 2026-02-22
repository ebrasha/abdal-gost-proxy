/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : client_config.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:16:06
 * Description : Client-side configuration models for SOCKS5, VLESS, and health check.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package models

// HealthCheckConfig holds periodic health check and stability options.
type HealthCheckConfig struct {
	Enabled        bool   `json:"enabled"`
	IntervalSeconds int   `json:"interval_seconds"`
	TimeoutSeconds  int   `json:"timeout_seconds"`
	MaxRetries     int   `json:"max_retries"`
	CheckURL       string `json:"check_url"`
}

// ClientConfig is the root client configuration loaded from abdal-gost-proxy-client.json.
type ClientConfig struct {
	LocalPort           int                `json:"local_port"`
	ServerAddr          string             `json:"server_addr"`
	ServerPort          int                `json:"server_port"`
	UUID                string             `json:"uuid"`
	RealityPublicKey    string             `json:"reality_public_key"`
	ShortID             string             `json:"short_id"`
	SNI                 string             `json:"sni"`
	Fingerprint        string             `json:"fingerprint"`
	Transport           string             `json:"transport"`
	ServiceName         string             `json:"service_name"`
	HealthCheck         HealthCheckConfig  `json:"health_check"`
}
