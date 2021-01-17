package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"

	"github.com/aau-network-security/haaukins-exercises/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_AUTH = "authkey"
	DEFAULT_SIGN = "signkey"
)

type Server struct {
	store store.Store
	auth  Authenticator
	tls   bool
}

func NewServer(conf *Config) (*Server, error) {

	st, err := store.NewStore(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Pass)
	if err != nil {
		return nil, err
	}

	s := &Server{
		store: st,
		auth:  NewAuthenticator(conf.SigninKey, conf.AuthKey),
		tls:   conf.TLS.Enabled,
	}
	return s, nil
}

func (s *Server) NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {

	streamInterceptor := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := s.auth.AuthenticateContext(stream.Context()); err != nil {
			return err
		}
		return handler(srv, stream)
	}

	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := s.auth.AuthenticateContext(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}

	opts = append([]grpc.ServerOption{
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	}, opts...)
	return grpc.NewServer(opts...)
}

type certificate struct {
	cPath    string
	cKeyPath string
}

func (s *Server) GrpcOpts(conf *Config) ([]grpc.ServerOption, error) {

	if conf.TLS.Enabled {
		creds, err := GetCreds(conf)

		if err != nil {
			return []grpc.ServerOption{}, errors.New("Error on retrieving certificates: " + err.Error())
		}
		log.Printf("INFO server is running in secure mode !")
		return []grpc.ServerOption{grpc.Creds(creds)}, nil
	}
	return []grpc.ServerOption{}, nil
}

func GetCreds(conf *Config) (credentials.TransportCredentials, error) {
	log.Printf("INFO preparing credentials for RPC")

	certificateProps := certificate{
		cPath:    conf.TLS.CertFile,
		cKeyPath: conf.TLS.CertKey,
	}

	creds, err := credentials.NewServerTLSFromFile(certificateProps.cPath, certificateProps.cKeyPath)
	if err != nil {
		return nil, err
	}
	return creds, nil
}

type Config struct {
	Host      string `yaml:"host"`
	Port      uint   `yaml:"port"` //gRPC endpoint
	AuthKey   string `yaml:"auth-key"`
	SigninKey string `yaml:"signin-key"`
	DB        struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"db"`
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"certfile"`
		CertKey  string `yaml:"certkey"`
		CAFile   string `yaml:"cafile"`
	} `tls:"tls,omitempty"`
}

func NewConfigFromFile(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, err
	}

	if c.Host == "" {
		log.Println("DBG host not provided in the configuration file")
		c.Host = "localhost"
	}

	if c.Port == 0 {
		log.Println("DBGpPort not provided in the configuration file")
		c.Port = 50095
	}

	if c.SigninKey == "" {
		log.Println("DBG signinKey not provided in the configuration file")
		c.SigninKey = DEFAULT_SIGN
	}

	if c.AuthKey == "" {
		log.Println("DBG authKey not provided in the configuration file")
		c.AuthKey = DEFAULT_AUTH
	}

	if c.DB.Host == "" || c.DB.User == "" || c.DB.Pass == "" {
		return nil, errors.New("DB parameters missing in the configuration file")
	}

	if c.DB.Port == 0 {
		c.DB.Port = 27017 //default port of Mongo DB
	}

	if c.TLS.Enabled {
		if c.TLS.CAFile == "" || c.TLS.CertKey == "" || c.TLS.CertFile == "" {
			return nil, errors.New("certificates parameters missing in the configuration file")
		}
	}

	return &c, nil
}
