package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_getRootCA(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "default root CA",
			args:    args{
				"rootCA.pem",
			},
			want:    "rootCA.pem",
			wantErr: false,
		},
		{
			name:    "missing rootCA",
			args:    args{
				"missing.pem",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRootCA(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRootCA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				want := readFile(t, tt.want)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("getRootCA() got = %v, want %v", got, tt.want)
				}
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