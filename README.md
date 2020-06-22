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
You can find a Docker image in `https://hub.docker.com/r/cquintana92/shoppinglist`.

## 2. How to use

### 2.1. Binary
Once you have your binary, you are ready to go!

In order to explore the different options, you can run `$ shopping-list --help` to inspect them.

```
$ shopping-list --help
GLOBAL OPTIONS:
   --loglevel value  Log Level [TRACE, DEBUG, INFO, WARN, ERROR] (default: "INFO") [$LOGLEVEL]
   --dbPath value    (default: "./shopping.sqlite") [$DB_PATH]
   --port value      (default: 5454) [$PORT]
   --help, -h        show help
   --version, -v     print the version
```

As you can see, the binary has sensible defaults, but you can easily override them via command line flags or environment variables. The default sqlite path is `CWD/shopping.sqlite`. In case the database does not exist, it will be created and initialized.

### 2.2. Docker image
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
      DB_PATH: "/data/shopping.sqlite"
      PORT: 5454
    restart: always
```

