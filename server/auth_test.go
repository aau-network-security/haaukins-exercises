package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"google.golang.org/grpc/credentials"

	pb "github.com/aau-network-security/haaukins-exercises/proto"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	HOST = "localhost:50095" //default in the server.go -> NewConfigFromFile
)

var (
	testCertPath    = os.Getenv("CERT")
	testCertKeyPath = os.Getenv("CERT_KEY")
	testCAPath      = os.Getenv("CA")
)

type Creds struct {
	Token    string
	Insecure bool
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

//Check authentication with the server
//Certificate created from here https://gist.github.com/cecilemuller/9492b848eb8fe46d462abeb26656c4f8
func TestStoreConnection(t *testing.T) {

	tokenCorret := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: DEFAULT_AUTH,
	})

	tokenError := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: "wrong-token",
	})

	tt := []struct {
		name  string
		token *jwt.Token
		err   string
	}{
		{name: "Test Normal Authentication", token: tokenCorret},
		{name: "Test Unauthorized", token: tokenError, err: "Invalid Authentication Key"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tokenString, err := tc.token.SignedString([]byte(DEFAULT_SIGN))
			if err != nil {
				t.Fatalf("Error creating the token")
			}

			authCreds := Creds{Token: tokenString}

			// Load the client certificates from disk
			certificate, err := tls.LoadX509KeyPair(testCertPath, testCertKeyPath)
			if err != nil {
				t.Fatalf("could not load client key pair: %s", err)
			}

			// Create a certificate pool from the certificate authority
			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(testCAPath)
			if err != nil {
				t.Fatalf("could not read ca certificate: %s", err)
			}

			// Append the certificates from the CA
			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				t.Fatalf("failed to append ca certs")
			}

			creds := credentials.NewTLS(&tls.Config{
				ServerName:   "localhost",
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})

			dialOpts := []grpc.DialOption{
				grpc.WithTransportCredentials(creds),
				grpc.WithPerRPCCredentials(authCreds),
			}

			conn, err := grpc.Dial(HOST, dialOpts...)
			if err != nil {
				t.Fatalf("Connection error: %v", err)
			}
			defer conn.Close()

			c := pb.NewExerciseStoreClient(conn)

			_, err = c.GetExercises(context.Background(), &pb.Empty{})

			if err != nil {
				st, ok := status.FromError(err)
				if ok {
					err = fmt.Errorf(st.Message())
				}

				if tc.err != "" {
					if tc.err != err.Error() {
						t.Fatalf("unexpected error (expected: %s) received: %s", tc.err, err.Error())
					}
					return
				}
				t.Fatalf("expected no error, but received: %s", err)
			}

			if tc.err != "" {
				t.Fatalf("expected error, but received none")
			}
		})
	}
}
