# ShoppingList (Server)

This repository contains the server for the ShoppingList.

Link to the Android client repository: https://gitlab.com/cquintana92/shoppinglist-android

## 1. How to get

### 1.1. Download a precompiled binary
You can go to the [releases section of this project](https://gitlab.com/cquintana92/shoppinglist-server/-/releases) and download the binary for your system (only Mac and Linux are supported (64 bit architectures) due to the sqlite3 dependency, which makes it hard to cross-compile).

### 1.2. Build from source.
This project has a Makefile that provides an easy way to run standard commands. You can run `make help` to see which options are available.
In order to build it, you can just run `make build`.

### 1.3. Docker image
You can find a Docker image in [Docker hub](https://hub.docker.com/r/cquintana92/shoppinglist).

## 2. How to use

### 2.1. Binary
Once you have your binary, you are ready to go!

In order to explore the different options, you can run `$ shopping-list --help` to inspect them.

```
$ shopping-list --help
GLOBAL OPTIONS:
   --loglevel value        Log Level [TRACE, DEBUG, INFO, WARN, ERROR] (default: "INFO") [$LOGLEVEL]
   --dbUrl value           Database URL connection (either sqlite3://PATH_TO_DB or PostgreSQL connection string) [$DB_PATH]
   --secretEndpoint value  Secret endpoint without authorization (useful for 3rd party integrations) [$SECRET_ENDPOINT]
   --secretBearer value    Secret bearer authorization in order to secure your server (Header: Authorization: Bearer YOUR_SECRET [$SECRET_BEARER]
   --port value            Port where the server will listen (default: 5454) [$PORT]
   --help, -h              show help
   --version, -v           print the version
```

As you can see, the binary has sensible defaults, but you can easily override them via command line flags or environment variables.

If you want to use an sqlite database, you will need to set `--dbUrl sqlite3://PATH_TO_THE_DATABASE`. In case the database does not exist, it will be created and initialized.

### 2.2. Securing your server

The server can receive a `secretBearer` parameter which will be used in order to prevent undesired access to your shopping list.

It must be passed to the server in the HTTP requests by setting the `Authorization` header to the following value:

```
Authorization: Bearer YOUR_SECRET_BEARER
```

That means, if your secret bearer is "ILikeTrains", your HTTP requests must contain the header:

```
Authorization: Bearer ILikeTrains
```

However, there are some third party integrations that do not support setting headers to your requests (such as IFTTT). In order to support these integrations, you can define a secret creation endpoint which does not need any authorization. It can be configured with the `secretEndpoint` parameter. 

That means, if you have set your secret bearer to `ILikeTrains` and your secret endpoint to `iliketrains`, you will be able to create new items by either:

- Sending a POST request with the `Authorization: Bearer ILikeTrains` header (`POST /`).
- Sending a POST request without authorization (`POST /iliketrains`).
 

### 2.3. Docker image
You can use a docker-compose like the following:

```yaml
version: '3'
services:
  shopping:
    image: cquintana92/shoppinglist:latest
    ports:
      - "5454:5454"
    volumes:
      - "./data:/data"
    environment:
      LOGLEVEL: INFO
      DB_URL: "sqlite3:///data/shopping.sqlite"
      PORT: 5454
      SECRET_ENDPOINT: super_secret_endpoint
      SECRET_BEARER: my_super_secret
    restart: always
```

Or for deploying with postgresql:

```yaml
version: '3'
services:
  # PostgreSQL database
  postgresql:
    image: postgres:12-alpine
    environment:
      POSTGRES_USER: postgresuser
      POSTGRES_PASSWORD: postgrespassword
      POSTGRES_DB: shoppinglistdb
    volumes:
      - "postgres_data:/var/lib/postgresql"
    restart: always

  # Shopping list server
  shopping:
    image: cquintana92/shoppinglist:latest
    ports:
      - "5454:5454"
    volumes:
      - "./data:/data"
    environment:
      LOGLEVEL: INFO
      DB_URL: "postgresql://postgresuser:postgrespassword@postgresql:5432/shoppinglistdb?sslmode=disable"
      PORT: 5454
      SECRET_ENDPOINT: super_secret_endpoint
      SECRET_BEARER: my_super_secret
    restart: always

volumes:
  postgres_data:
```

