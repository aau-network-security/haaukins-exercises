haaukins-exercises [![Build Status](https://github.com/aau-network-security/haaukins-exercises/workflows/test/badge.svg)](https://github.com/aau-network-security/haaukins-exercises/actions)

Haaukins exercises is internally used for managing exercises provided by [Haaukins](https://github.com/aau-network-security/haaukins) platform.
This microservices is basically a store (using MongoDB) where users, through gRPC communication, can read and insert exercises.
Exercises structure can be found in `model` folder.

## Production usage

Docker image of haaukins exercises could be used in any docker compose file if environment variables provided correctly.
When using in production, you can specify image address instead of building it from source code.

Haaukins exercises image with recent changes will be available at docker hub, with released tag.
Therefore just make sure that `.env`  and `config.yml` files are set correctly in order to run this microservice.

Steps to run it in production:

 - Make sure you have configured `.env` and `config.yml` according to the instructions in [configuration](#configuration) section.
 - `curl -o docker-compose.yml https://raw.githubusercontent.com/aau-network-security/haaukins-exercises/main/docker-compose.yml`
 - Change image version to the latest one

## Configuration
Haaukins exercises uses two crucial configuration files that have to be set: `.env` and `config.yml`. 

Here is the information which should be included into `.env` file: 

```text
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=toor

CONFIG_PATH=/scratch/configsconfig.yml
MONGO_DATA_PATH=/scratch/configs/data
CERTS_PATH=/scratch/configs/certs

ME_CONFIG_MONGODB_ADMINUSERNAME=root
ME_CONFIG_MONGODB_ADMINPASSWORD=toor
ME_CONFIG_BASICAUTH_USERNAME=admin
ME_CONFIG_BASICAUTH_PASSWORD=secret
```

- `MONGO_INITDB_ROOT_USERNAME`: admin user for the DB (it has to be the same with `config.yml`)
- `MONGO_INITDB_ROOT_PASSWORD`: password for admin user for the DB (it has to be the same with `config.yml`)
- `CERTS_PATH`: Should be provided if TLS is enabled  and certificates should be valid for provided host in `config.yml` file. 
- `CONFIG_PATH`: Path of your `config.yml` file which is mount in `docker-compose.yml` file
- `ME_CONFIG_MONGODB_ADMINUSERNAME`: (Optional) 
- `ME_CONFIG_MONGODB_ADMINPASSWORD`: (Optional)
- `ME_CONFIG_BASICAUTH_USERNAME`: (Optional)
- `ME_CONFIG_BASICAUTH_PASSWORD`: (Optional)

Here is the information which should be included into `config.yml` file: 

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

Haaukins exercises could be run by:

```bash
docker-compose up -d --build
```
