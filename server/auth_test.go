package server

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func Test_auth_AuthenticateContext(t *testing.T) {
	type fields struct {
		sKey string
		aKey string
	}

	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx: metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdSI6ImF1dGhrZXkifQ.SdJFhs6LJsOErMVQ_6s5PAVShN3wx5KFU9Tc9w7-jQs"}}),
		},
		{
			name: "Fail - No such key",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     context.Background(),
			wantErr: ErrMissingKey,
		},
		{
			name: "Fail - Missing token",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{}),
			wantErr: ErrMissingKey,
		},
		{
			name: "Fail - Missing token",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {""}}),
			wantErr: ErrMissingKey,
		},
		{
			name: "Fail - Broken token",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdSI6ImF1dGhrZXkifQ"}}),
			wantErr: ErrInvalidTokenFormat,
		},
		{
			name: "Fail - No claims",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.we6lzGXF28AJ3rjDLfntt-lp3VZEJwjNCRBWeGKBQoM"}}),
			wantErr: ErrInvalidTokenFormat,
		},
		{
			name: "Fail - Wrong auth key",
			fields: fields{
				aKey: DEFAULT_AUTH,
				sKey: DEFAULT_SIGN,
			},
			ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdSI6InNvbWUgcmFuZG9tIHN0dWZmIn0.pBAkfQS_LY4IhjkqAoOs_p-G03ICxeV01G5PbPeYCks"}}),
			wantErr: ErrInvalidAuthKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &auth{
				sKey: tt.fields.sKey,
				aKey: tt.fields.aKey,
			}
			if err := a.AuthenticateContext(tt.ctx); err != tt.wantErr {
				t.Errorf("auth.AuthenticateContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
