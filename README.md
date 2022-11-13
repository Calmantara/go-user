# Go-User

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [maintainer](#maintainer)

## About <a name = "about"></a>

This is project to handle relation between 
user <-> photos <-> credit card

## Getting Started <a name = "getting_started"></a>

To running this project in local, you need to install `docker compose`.
Set your `ENV` variable as `dev` to local or development purpose. Set as `prod` to running in production system.

You can set your configuration under `./manifest` directory.
```
postgresRead:
  host: 'postgres'
  port: 5432
  database: 'user'
  username: 'postgres'
  password: 'postgresAdmin'
  timeZone: 'GMT'
  autoMigrate: true
  mode: 'read'
  enableLog: false
  maxConnection: 10
  # amount
  maxIdleConnection: 1
  # in minute
  maxIdleConnectionTtl: 5 
```

After all containers up and running, run migration sql script under `./migration` directory under `user` database.

### Prerequisites

This project is under `docker-compose` management.
See [docker-compose](https://docs.docker.com/compose/install/) for installation


### Installing

To run the project, execute command below

```
make up
```

To shutting down the project, execute command below

```
make down
```

All endpoints can be found in [postman.json](./postman.json) file

## Maintainer <a name = "maintainer"></a>

Calmantara
