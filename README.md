# Abdal Gost Proxy

**English** | [فارسی (README.fa.md)](README.fa.md)

---

![dgram.jpg](dgram.jpg)

---

## Why This Software Exists

**Abdal Gost Proxy** is a client/server proxy system that helps you:

- **Bypass censorship and DPI** by carrying traffic over VLESS + gRPC over **XTLS-Reality**, so it looks like normal HTTPS to the target site (e.g. google.com or google.com).
- **Avoid dedicated SSL certificates** on your server: Reality “borrows” the appearance of a real site, so the server does not need its own TLS cert.
- **Use a local SOCKS5 proxy** on your PC; you point your browser or apps to it, and traffic is sent through the tunnel to your server.

It combines **[Gost](https://github.com/go-gost/gost)**-style tunnel management with **[Xray-core](https://github.com/XTLS/Xray-core)** for the VLESS protocol and Reality security.

---

## Features

| Feature | Description |
|--------|-------------|
| **VLESS** | Lightweight protocol for proxy traffic. |
| **XTLS-Reality** | Hides the proxy as a visit to a real site (e.g. google.com) without your own SSL cert. |
| **gRPC (HTTP/2)** | Transport layer that helps traffic pass through DPI. |
| **uTLS fingerprint** | Client can mimic browser TLS (e.g. Chrome) to reduce fingerprinting. |
| **SOCKS5 gate** | Local SOCKS5 listener; apps use it as a normal proxy. |
| **Health check** | Optional periodic checks; configurable interval, timeout, retries, and URL. |
| **Auto re-dial** | After repeated health-check failures, the client restarts the tunnel. |
| **Fallback** | On the server, unauthenticated traffic can be sent to a fallback site (e.g. google.com). |

---

## How It Works

![dgram.jpg](dgram.jpg)

```
[Browser/App] → SOCKS5 (127.0.0.1:10808) → [Gate] → [Abdal Gost Proxy client]
                                                          ↓
                                              VLESS + gRPC + Reality
                                                          ↓
[Abdal Gost Proxy server] (0.0.0.0:443) ← Reality handshake ← [Internet]
```

- **Server:** Listens on 443. Accepts VLESS over gRPC over Reality. If the connection is not valid Reality/VLESS, it can fallback to a real site (e.g. show google.com).
- **Client:** Listens on a local port (e.g. 10808) as SOCKS5. Forwards to Xray, which wraps traffic in VLESS and sends it over gRPC + Reality to the server. Optional health check and auto re-dial keep the tunnel up.

---

## Requirements

- **Go 1.21+** (for building).
- **Server:** A VPS or machine with a public IP, port 443 open.
- **Client:** Windows or Linux; config file next to the executable.

---

## Build

From the project root:

```bash
# All platforms (Windows + Linux) and configs into dist/
build-dist.bat
```

Output:

- `dist/windows/` — Server, client, and reality-keygen (`.exe`) + config samples.
- `dist/linux/` — Same, without `.exe`.

Or build manually:

```bash
go mod tidy
go build -o abdal-gost-proxy-server.exe main.go
go build -o abdal-gost-proxy-client.exe client_main.go
go build -o reality-keygen.exe ./tools/reality-keygen
```

---

## Configuration

### 1. Generate keys and UUID

Run once:

```bash
go run ./tools/reality-keygen
```

You get:

- **uuid** — Use the same value for server `users[].id` and client `uuid`.
- **private_key** — Server only; put in `abdal-gost-proxy-server.json` → `reality_settings.private_key`.
- **reality_public_key** — Client only; put in `abdal-gost-proxy-client.json` → `reality_public_key`.

Generate **short_id** (e.g. hex):

```bash
openssl rand -hex 8
```

Use one or more in server `reality_settings.short_ids` and pick one for client `short_id`.

---

### 2. Server config: `abdal-gost-proxy-server.json`

| Field | Description |
|-------|-------------|
| `listen_address` | Bind address (e.g. `0.0.0.0`). |
| `listen_port` | Usually `443`. |
| `protocol` | `vless`. |
| `users` | List of VLESS users; each has `id` (UUID), `email`, `flow` (leave empty or omit for gRPC). |
| `reality_settings.dest` | Fallback site:port when connection is not valid (e.g. `www.google.com:443`). |
| `reality_settings.server_names` | SNI list (e.g. `["www.google.com","google.com"]`). |
| `reality_settings.private_key` | From reality-keygen (keep secret). |
| `reality_settings.short_ids` | List of short IDs (e.g. from `openssl rand -hex 8`). |
| `transport.type` | `grpc`. |
| `transport.service_name` | Must match client (e.g. `abdal-grpc-stream`). |
| `transport.multi_mode` | Optional gRPC multi-mode. |
| `fallback.dest` | Fallback port or host:port if needed. |
| `fallback.xver` | Proxy protocol version (e.g. `0`). |

**Supported options (server):**

| Option | Allowed values / notes |
|--------|------------------------|
| `listen_address` | `0.0.0.0` (all interfaces) or a specific IP, e.g. `192.168.1.1`. |
| `listen_port` | Any free port; typically `443`. |
| `protocol` | `vless` (only protocol used in this system). |
| `users[].flow` | `""` (empty) when using gRPC — **required**. Do **not** use `xtls-rprx-vision` with gRPC. |
| `reality_settings.dest` | String `"host:port"`, e.g. `www.google.com:443`, `www.google.com:443`. Must be a real TLS site. |
| `reality_settings.server_names` | Array of SNI strings; first usually matches `dest` hostname. |
| `reality_settings.short_ids` | Array of hex strings (e.g. from `openssl rand -hex 8`); 2–16 chars each. |
| `transport.type` | `grpc` (this system uses gRPC only). |
| `transport.service_name` | Any string; must match client. Avoid default names; e.g. `abdal-grpc-stream`. |
| `transport.multi_mode` | `true` or `false`. |
| `fallback.dest` | Number (port only, e.g. `80`) or string `"host:port"`. |
| `fallback.xver` | `0` (off), `1`, or `2` (Proxy Protocol). |

Example (minimal):

```json
{
  "listen_address": "0.0.0.0",
  "listen_port": 443,
  "protocol": "vless",
  "users": [
    {
      "id": "YOUR-UUID-FROM-REALITY-KEYGEN",
      "email": "admin@abdal",
      "flow": ""
    }
  ],
  "reality_settings": {
    "enabled": true,
    "dest": "www.google.com:443",
    "server_names": ["www.google.com", "google.com"],
    "private_key": "YOUR_PRIVATE_KEY",
    "short_ids": ["1a2b3c4d5e6f"]
  },
  "transport": {
    "type": "grpc",
    "service_name": "abdal-grpc-stream",
    "multi_mode": true
  },
  "fallback": { "dest": 80, "xver": 0 }
}
```

**Note:** For gRPC transport, leave `flow` empty (or omit it). Do not use `xtls-rprx-vision` with gRPC.

---

### 3. Client config: `abdal-gost-proxy-client.json`

| Field | Description |
|-------|-------------|
| `local_port` | Local SOCKS5 port (e.g. `10808`). |
| `server_addr` | Server IP or domain. |
| `server_port` | Usually `443`. |
| `uuid` | Same as server user `id`. |
| `reality_public_key` | From reality-keygen (public key of server). |
| `short_id` | One of server’s `short_ids`. |
| `sni` | Must match Reality site (e.g. `www.google.com`). |
| `fingerprint` | uTLS fingerprint: e.g. `chrome`, `firefox`. |
| `transport` | `grpc`. |
| `service_name` | Must match server (e.g. `abdal-grpc-stream`). |
| `health_check` | Optional: `enabled`, `interval_seconds`, `timeout_seconds`, `max_retries`, `check_url`. |

**Supported options (client):**

| Option | Allowed values / notes |
|--------|------------------------|
| `local_port` | Any free port 1–65535; typical `10808`. |
| `server_addr` | Server IP or domain. |
| `server_port` | Usually `443`. |
| `sni` | Must match the Reality site (same as server `reality_settings.dest` hostname), e.g. `www.google.com`, `www.google.com`. |
| `fingerprint` | uTLS fingerprint; one of: `chrome`, `firefox`, `safari`, `ios`, `android`, `edge`, `360`, `qq`, `random`. Default if empty: `chrome`. |
| `transport` | `grpc` (this system uses gRPC only). |
| `service_name` | Must match server exactly. |
| `health_check.enabled` | `true` or `false`. |
| `health_check.interval_seconds` | Positive number; interval between checks (e.g. `5`). |
| `health_check.timeout_seconds` | Positive number; timeout per check (e.g. `3`). |
| `health_check.max_retries` | Positive number; retries before re-dial (e.g. `3`). |
| `health_check.check_url` | Any HTTP(S) URL used to test connectivity via the proxy (e.g. `http://www.google.com/generate_204`). |

Example:

```json
{
  "local_port": 10808,
  "server_addr": "YOUR_SERVER_IP",
  "server_port": 443,
  "uuid": "YOUR-UUID",
  "reality_public_key": "YOUR_PUBLIC_KEY",
  "short_id": "1a2b3c4d5e6f",
  "sni": "www.google.com",
  "fingerprint": "chrome",
  "transport": "grpc",
  "service_name": "abdal-grpc-stream",
  "health_check": {
    "enabled": true,
    "interval_seconds": 5,
    "timeout_seconds": 3,
    "max_retries": 3,
    "check_url": "http://www.google.com/generate_204"
  }
}
```

---

## How to Use

### Server

1. Put `abdal-gost-proxy-server.json` next to the server binary.
2. Run:
   - Windows: `abdal-gost-proxy-server.exe`
   - Linux: `./abdal-gost-proxy-server`
3. Ensure port 443 is open in the firewall.

### Client

1. Put `abdal-gost-proxy-client.json` next to the client binary.
2. Run:
   - Windows: `abdal-gost-proxy-client.exe`
   - Linux: `./abdal-gost-proxy-client`
3. Set your browser or system proxy to **SOCKS5**, host **127.0.0.1**, port **10808** (or your `local_port`).

---


## 🐛 Reporting Issues

If you encounter any issues or have configuration problems, please reach out via email at **Prof.Shafiei@Gmail.com**. You can also report issues on GitLab or GitHub.

---

## ❤️ Donation

If you find this project helpful and would like to support further development, please consider making a donation:

- [Donate Here](https://alphajet.ir/abdal-donation)

---

## 🤵 Programmer

Handcrafted with passion by **Ebrahim Shafiei (EbraSha)**

- **E-Mail:** Prof.Shafiei@Gmail.com  
- **Telegram:** [@ProfShafiei](https://t.me/ProfShafiei)  
- **GitHub:** [ebrasha](https://github.com/ebrasha)  
- **LinkedIn:** [profshafiei](https://www.linkedin.com/in/profshafiei/)

---

## 📜 License

This project is licensed under the **GPLv2 or later** License.
