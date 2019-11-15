package verifier

import (
	"io/ioutil"
	"testing"
)

func TestVerifyTimestamp(t *testing.T) {
	tests := []struct{
		name string
		container timestampContainer
		wantErr bool
	}{
		{
			name: "valid single timestamp",
			container: timestampContainer{
				data: []byte("hello world\n"),
				timestamps: []*Timestamped{
					{
						Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
						LtvTimestamp: map[string]*LTV{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nested timestamp",
			container: timestampContainer{
				data:       []byte("hello world\n"),
				timestamps: []*Timestamped{
					{
						Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
						LtvTimestamp: map[string]*LTV{},
					},
					{
						Rfc3161Timestamp: readFile(t, "hello_world_response.tsr.data_response.tsr"),
						LtvTimestamp: map[string]*LTV{},
					},
				},
			},
		},
		{
			name: "hash mismatch",
			container: timestampContainer{
				data: []byte("hello world"),
				timestamps: []*Timestamped{
					{
						Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
						LtvTimestamp: map[string]*LTV{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no timestamps",
			container:timestampContainer{
				data:       []byte("hello world"),
				timestamps: []*Timestamped{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.container.Verify(); err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		t.Errorf("Could not read file: %s", err)
	}
	return data
}