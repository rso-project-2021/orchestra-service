# orchestra-service
![Build](https://github.com/rso-project-2021/orchestra-service/actions/workflows/build.yml/badge.svg)
![Deploy](https://github.com/rso-project-2021/orchestra-service/actions/workflows/deploy.yml/badge.svg)  
Microservice used for orchestrating operations of other microservices. It is used for making quick station reservations.

## Environment file
In root of your local repository add `config.json` file.
```
{
    "grpc_address": "localhost:9000",
    "server_address": "0.0.0.0:8080",
    "gin_mode": "debug"
}
```