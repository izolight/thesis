package verifier_test

import "testing"

func TestVerifySignatureData(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
	}{
		{
			name:     "valid signature",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verifier.Verify(); err != nil != tt.wantErr {
				t.Fatalf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
