# haaukins-exercises [![Build Status](https://github.com/aau-network-security/haaukins-exercises/workflows/test/badge.svg)](https://github.com/aau-network-security/haaukins-exercises/actions)

Haaukins microservice that store all the CTF challenges

### Configuration
Make sure to have a `config.yml` with these parameters in order to run the app

```yaml
host: localhost
port: 50095
auth-key: random
signin-key: random
db:
  host: mongo
  user: root
  pass: toor
  port: 27017
tls:
  enabled: false
  certfile: /certs/localhost.crt
  certkey: /certs/localhost.key
  cafile: /certs/haaukins-store.com.crt
```

### Run

```bash
docker-compose up -d --build
```
