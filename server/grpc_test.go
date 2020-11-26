package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"testing"

	pb "github.com/aau-network-security/haaukins-exercises/proto"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	NExercises  = []string{"ftp", "xxs", "xxe", "sql", "mitm", "crypto", "shad", "rand", "ccs"}
	NCategories = []string{"forensics", "binary"}
)

func createTestClientConn() (*grpc.ClientConn, error) {

	tokenCorret := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		AUTH_KEY: DEFAULT_AUTH,
	})

	tokenString, err := tokenCorret.SignedString([]byte(DEFAULT_SIGN))
	if err != nil {
		return nil, err
	}

	authCreds := Creds{Token: tokenString}

	// Load the client certificates from disk
	certificate, err := tls.LoadX509KeyPair(testCertPath, testCertKeyPath)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(testCAPath)
	if err != nil {
		return nil, err
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
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

	// Create a connection with the TLS credentials
	conn, err := grpc.Dial(HOST, dialOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TestServer_GetExerciseByTags(t *testing.T) {
	conn, err := createTestClientConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewExerciseStoreClient(conn)

	tt := []struct {
		name     string
		tags     []string
		expected int
		err      bool
	}{
		{name: "Normal Get exercises by tags not empty", tags: NExercises[:4], expected: 4},
		{name: "Normal Get exercises by tags empty", tags: []string{}, expected: 0},
		{name: "Invalid tags", tags: []string{"randomex"}, expected: 0, err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.GetExerciseByTags(context.Background(), &pb.GetExerciseByTagsRequest{Tag: tc.tags})
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("Error get exercises: %v", err)
			}
			if tc.err {
				t.Fatal("Error expected")
			}
			if len(resp.Exercises) != tc.expected {
				t.Fatalf("Expected number of challenges %d, got %d", tc.expected, len(resp.Exercises))
			}
		})
	}
}

func TestServer_GetCategories(t *testing.T) {
	conn, err := createTestClientConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewExerciseStoreClient(conn)

	resp, err := c.GetCategories(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Categories) != len(NCategories) {
		t.Fatalf("Expected number of category %d, got %d", len(NCategories), len(resp.Categories))
	}
}

func TestServer_AddCategory(t *testing.T) {
	conn, err := createTestClientConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewExerciseStoreClient(conn)

	tt := []struct {
		name  string
		categ string
		err   bool
	}{
		{name: "Normal category", categ: "randomcategory"},
		{name: "Already existing category", categ: "randomcategory", err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.AddCategory(context.Background(), &pb.AddCategoryRequest{
				Tag:  tc.categ,
				Name: tc.categ,
			})
			if err != nil {
				if tc.err {
					return
				}
				t.Fatalf("Error insert category: %v", err)
			}
			if tc.err {
				t.Fatal("Error expected")
			}
		})
	}
}
