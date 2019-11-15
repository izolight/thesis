package verifier

import "testing"

func TestVerifyIDToken(t *testing.T) {
	tests := []struct{
		name string
		verifier Verifier
		wantErr bool
	}{
		{
			name: "valid id token",
			verifier: idTokenVerifier{

			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verifier.Verify(); err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}