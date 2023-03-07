package server

import (
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Config
		wantErr bool
	}{
		{
			name: "Success",
			path: "../testdata/success.yml",
			want: &Config{
				Host:      "localhost",
				Port:      50095,
				AuthKey:   "test",
				SigninKey: "test",
				DB: Remote{
					Host: "localhost",
					Port: 27017,
					User: "haaukins",
					Pass: "haaukins",
				},
				TLS: TLSConf{
					Enabled:  true,
					CertFile: "path-to-cert",
					CertKey:  "path-to-key",
				},
			},
		},
		{
			name: "Success - default values",
			path: "../testdata/success-defaults.yml",
			want: &Config{
				Host:      "localhost",
				Port:      50095,
				AuthKey:   "authkey",
				SigninKey: "signkey",
				DB: Remote{
					Host: "localhost",
					Port: 27017,
					User: "haaukins",
					Pass: "haaukins",
				},
				TLS: TLSConf{
					Enabled:  true,
					CertFile: "path-to-cert",
					CertKey:  "path-to-key",
				},
			},
		},
		{
			name:    "Fail - Database",
			path:    "../testdata/fail-databaseconf.yml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Fail - tls",
			path:    "../testdata/fail-tls.yml",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
