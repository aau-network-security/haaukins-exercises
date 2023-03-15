package mongodb

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/aau-network-security/haaukins-exercises/proto"
	"github.com/aau-network-security/haaukins-exercises/testdata"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	ctx := context.TODO()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create connection to docker")
	}

	resource, err := pool.Run("mongo", "5", []string{"MONGO_INITDB_ROOT_USERNAME=haaukins", "MONGO_INITDB_ROOT_PASSWORD=haaukins"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create docker container")
	}

	port, err := strconv.ParseUint(resource.GetPort("27017/tcp"), 10, 32)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get port for service")
	}

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	type args struct {
		host string
		port uint
		user string
		pass string
	}
	tests := []struct {
		name    string
		args    args
		want    *store
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				host: "localhost",
				port: uint(port),
				user: "haaukins",
				pass: "haaukins",
			},
			wantErr: false,
			want:    &store{},
		},
		{
			name: "Fail",
			args: args{
				host: "example.com",
				port: 0,
				user: "haaukins",
				pass: "haaukins",
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStore(ctx, tt.args.host, tt.args.port, tt.args.user, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlePrivacyUniverse(t *testing.T) {

	tests := []struct {
		name string
		in   []*proto.Exercise
		want []*proto.Exercise
	}{
		{
			name: "Success - Jobspace only",
			in:   testdata.InSuccessJobspaceOnly,
			want: testdata.OutSuccessJobspaceOnly,
		},
		{
			name: "success - no Privacy universe",
			in:   testdata.InSuccessNoPU,
			want: testdata.OutSuccessNoPU,
		},
		{
			name: "Success - mixed",
			in:   testdata.InMixSuccess,
			want: testdata.OutMixSuccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlePrivacyUniverse(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handlePrivacyUniverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
