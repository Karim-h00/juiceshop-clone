package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"strings"
)

func getAssetPath(mediaType string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	id := base64.RawURLEncoding.EncodeToString(base)

	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", id, ext)
}

func (cfg config) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func (cfg config) getAssetURL(assetPath string) string {
	return fmt.Sprintf("%s/assets/%s", cfg.baseURL, assetPath)
}

func slugToName(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
		return ip
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func validateRegisterInput(email, username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func validateLoginInput(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) > 254 {
		return errors.New("email is too long")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}

var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

func validateUsername(username string) error {

	if !validUsername.MatchString(username) {
		return errors.New("username must be 3-32 characters, letters/numbers/'_' only")
	}
	return nil
}

func validatePassword(password string) error {
	length := []rune(password)
	if len(length) < 8 || len(length) > 128 {
		return errors.New("password length must be 8-128 characters")
	}
	return nil
}
