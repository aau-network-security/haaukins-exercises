package server

import (
	"context"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/aau-network-security/haaukins-exercises/proto"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(ctx context.Context) (proto.ExerciseStoreClient, func()) {
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

	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	serv, err := NewServer(Config{
		Host:      "localhost",
		Port:      50095,
		AuthKey:   "test",
		SigninKey: "test",
		DB: Remote{
			Host: "localhost",
			Port: uint(port),
			User: "haaukins",
			Pass: "haaukins",
		},
		TLS: TLSConf{
			Enabled:  true,
			CertFile: "path-to-cert",
			CertKey:  "path-to-key",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("error creating testserver")
	}
	proto.RegisterExerciseStoreServer(baseServer, serv)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("error serving grpc")
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error().Err(err).Msg("error connecting testserver")
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing test server")
		}
		baseServer.Stop()
		if err := pool.Purge(resource); err != nil {
			log.Error().Err(err).Msg("error closing mongo")
		}
	}

	client := proto.NewExerciseStoreClient(conn)

	return client, closer
}

func TestServer_AddExercises(t *testing.T) {
	client, closer := setupTestServer(context.Background())

	type args struct {
		ctx     context.Context
		request *proto.AddExercisesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.Empty
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:     metadata.NewIncomingContext(context.Background(), map[string][]string{"token": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdSI6ImF1dGhrZXkifQ.SdJFhs6LJsOErMVQ_6s5PAVShN3wx5KFU9Tc9w7-jQs"}}),
				request: &proto.AddExercisesRequest{Exercises: []*proto.Exercise{}},
			},
			want: &proto.Empty{},
		},
	}
	for _, tt := range tests {
		t.Cleanup(closer)
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.AddExercises(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.AddExercises() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Server.AddExercises() = %v, want %v", got, tt.want)
			}
		})
	}
}
