[![Actions Status](https://github.com/qwertmax/golang-basic-authentication-jwt/workflows/Go/badge.svg)](https://github.com/qwertmax/golang-basic-authentication-jwt/actions)

# golang-basic-authentication-jwt

### ENV variables

default values for ENV are:


`DB_HOST=0.0.0.0`

`DB_PORT=5432`

`DB_USER=postgres`

`DB_PASSWORD=my-secret-password`

`DB_NAME=golang`



How to use an image of docker container

before you needs to auth via github username 

`docker login docker.pkg.github.com --username <github username>``

`docker pull docker.pkg.github.com/qwertmax/golang-basic-authentication-jwt/golang-basic-authentication-jwt:latest`

or Docker Compose

```yml
version: '3'

services:

  backend:
    image: docker.pkg.github.com/qwertmax/golang-basic-authentication-jwt/golang-basic-authentication-jwt:latest
    environment:
      - DB_HOST=0.0.0.0
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=my-secret-password
      - DB_NAME=golang
    ports:
      - "8000:8000"

  db:
    image: postgres:latest
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=1
      - POSTGRES_DB=golang
    volumes:
      - ./postgres-data:/var/lib/postgresql/data

```
