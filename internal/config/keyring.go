package config

import (
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"
)

// go-keyring wrappers for interacting with the system keyring

const serviceName = "how"

func SetAPIKeyInKeyring(provider, apiKey string) error {
	return keyring.Set(serviceName, provider, apiKey)
}

func GetAPIKeyFromKeyring(provider string) (string, error) {
	key, err := keyring.Get(serviceName, provider)
	if err == keyring.ErrNotFound {
		return "", nil // Return empty string instead of error for not found
	}
	return key, err
}

func DeleteAPIKeyFromKeyring(provider string) error {
	err := keyring.Delete(serviceName, provider)
	if err == keyring.ErrNotFound {
		return nil // Not an error if key doesn't exist
	}
	return err
}

func ListProvidersWithKeys() ([]string, error) {
	providers := GetProviders()
	var providersWithKeys []string

	for _, provider := range providers {
		key, err := keyring.Get(serviceName, provider)
		if err != nil && err != keyring.ErrNotFound {
			return nil, fmt.Errorf("error checking key for %s: %w", provider, err)
		}
		if key != "" {
			providersWithKeys = append(providersWithKeys, provider)
		}
	}

	return providersWithKeys, nil
}

func HasAPIKeyInKeyring(provider string) (bool, error) {
	_, err := keyring.Get(serviceName, provider)
	if err == keyring.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Mask all but the last 8 characters of a key
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		half := len(key) / 2
		return strings.Repeat("*", 8-half) + key[half:]
	}
   return "****" + key[len(key)-8:]
}