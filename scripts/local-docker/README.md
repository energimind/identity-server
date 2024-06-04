# About

## docker-compose.yaml

This YAML file defines the Docker Compose services for the Identity Server. It includes services for MongoDB, Redis,
and the Identity Web application. Each service has a profiles attribute that determines which profiles the service
belongs to. For example, the mongo and redis services belong to the default and web profiles, while the identity-web
service only belongs to the web profile.

## start.sh

This is a shell script used to start the Docker Compose services. It uses the docker compose command to start the
services defined in the docker-compose.yaml file.

Usage:

```bash
./start.sh
```

or to start with a specific profile:

```bash
IS_PROFILE=web ./start.sh
```

Additional docker compose parameters can be passed to the script after the ``start.sh``.

## stop.sh

This is a shell script used to stop the Docker Compose services. It uses the docker compose command to stop the
services defined in the docker-compose.yaml file.

Usage:

```bash
./stop.sh
```

Additional docker compose parameters can be passed to the script after the ``stop.sh``.

## logs.sh

This is a shell script used to view the logs of the Docker Compose services. It uses the docker compose command to
view the logs of the services defined in the docker-compose.yaml file.

Usage:

```bash
./logs.sh
```

Additional docker compose parameters can be passed to the script after the ``logs.sh``.

## mongo-cli.sh

This is a shell script used to run the MongoDB CLI. It uses the docker compose command to run the MongoDB CLI.

Usage:

```bash
./mongo-cli.sh
```

## redis-cli.sh

This is a shell script used to run the Redis CLI. It uses the docker compose command to run the Redis CLI.

Usage:

```bash
./redis-cli.sh
```
