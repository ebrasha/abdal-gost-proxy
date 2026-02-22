/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal Gost Proxy
 * File Name : keys.go
 * Author : Ebrahim Shafiei (EbraSha)
 * Email : Prof.Shafiei@Gmail.com
 * Created On : 2026-02-14 22:30:03
 * Description : Curve25519 key pair, UUID v4, and short_id (hex) generation for XTLS-Reality / VLESS.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// GenerateRealityKeys returns private and public key for Reality (X25519), base64 RawURL encoded.
// Uses only standard library (crypto/ecdh, Go 1.20+) so no extra module dependency is needed.
func GenerateRealityKeys() (privateKey, publicKey string, err error) {
	privKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	privBytes := privKey.Bytes()
	pubBytes := privKey.PublicKey().Bytes()
	privateKey = base64.RawURLEncoding.EncodeToString(privBytes)
	publicKey = base64.RawURLEncoding.EncodeToString(pubBytes)
	return privateKey, publicKey, nil
}

// GenerateUserUUID generates a UUID v4 for VLESS user authentication (standard library only).
func GenerateUserUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// UUID v4: set version (4) and variant bits (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16])), nil
}

// GenerateShortIDs returns hex strings for Reality short_ids. Rule: 2–16 chars each. Count = number of IDs.
func GenerateShortIDs(count int) ([]string, error) {
	if count <= 0 {
		count = 2
	}
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		out = append(out, hex.EncodeToString(b))
	}
	return out, nil
}
