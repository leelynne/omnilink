package proto

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

const key = "BE-33-DB-F3-50-4D-79-A7-52-CB-51-A4-D0-72-AF-AF"

func TestNew(t *testing.T) {
	k, _ := parseKey(key)
	_, err := NewConnection("192.168.1.85:4369", k)
	if err != nil {
		t.Fatal(err)
	}
}

func parseKey(key string) (StaticKey, error) {
	hexOnly := strings.Replace(key, "-", "", -1)
	keyBytes, err := hex.DecodeString(hexOnly)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 16 {
		return nil, fmt.Errorf("Key %s must be 16 bytes long", key)
	}
	return StaticKey(keyBytes), nil
}
