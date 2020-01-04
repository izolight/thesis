package api_test

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/api"
	"testing"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		caFile string
		wantErr bool
	}{
		{
			name: "normal router",
			caFile: "rootCA.pem",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := testLogger{
				t:       t,
				wantErr: tt.wantErr,
			}
			caFile := readFile(t, "rootCA.pem")
			api.NewRouter(logger, caFile)
		})
	}
}

type testLogger struct {
	t *testing.T
	wantErr bool
}

func (t testLogger) Print(...interface{}) {
	panic("implement me")
}

func (t testLogger) Printf(string, ...interface{}) {
	panic("implement me")
}

func (t testLogger) Println(...interface{}) {
	panic("implement me")
}

func (t testLogger) Fatal(...interface{}) {
	panic("implement me")
}

func (t testLogger) Fatalf(string, ...interface{}) {
	panic("implement me")
}

func (t testLogger) Fatalln(args ...interface{}) {
	if !t.wantErr {
		t.t.Error(args...)
	}
	panic("implement me")
}

func (t testLogger) Panic(...interface{}) {
	panic("implement me")
}

func (t testLogger) Panicf(string, ...interface{}) {
	panic("implement me")
}

func (t testLogger) Panicln(...interface{}) {
	panic("implement me")
}
