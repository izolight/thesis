package verifier

import (
	"crypto"
	"encoding/hex"
	"testing"
)

func TestVerifyHash(t *testing.T) {
	tests := []struct{
		name string
		data []byte
		hash string
		algorithm crypto.Hash
		wantErr bool
	}{
		{
			name: "valid sha256",
			data: []byte("Hello world"),
			hash: "64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c",
			algorithm: crypto.SHA256,
			wantErr: false,
		},
		{
			name: "wrong sha256",
			data: []byte("Hello world"),
			hash: "1894a19c85ba153acbf743ac4e43fc004c891604b26f8c69e1e83ea2afc7c48f",
			algorithm: crypto.SHA256,
			wantErr: true,
		},
		{
			name: "sha256 hash with sha384 algorithm",
			data: []byte("Hello world"),
			hash: "64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c",
			algorithm: crypto.SHA384,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, _ := hex.DecodeString(tt.hash)
			if err := verifyHash(tt.data, hash, tt.algorithm); err != nil != tt.wantErr {
				t.Errorf("VerifyHashes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
