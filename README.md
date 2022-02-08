# Gosnel

## Overview

## Quick Start

Install Gosnel:

```shell
$ git clone https://github.com/youngjae-lim/gosnel
$ cd gosnel
$ make install_cli
$ cd ..
$ rm -rf gosnel
```

Create an app and run:

````shell
# Check if gosnel app is in your ```/go/bin/``` path
$ which gosnel

# Create a new app
$ gosnel new testapp

# Go into testapp directory
$ cd testapp

# Run the application
$ make start
````

Go to your favorite webbrowser and type `http://localhost:4000`

You will be able to see a homepage if the app runs successfully.

## Features

- Migrations
- Authentication
- Models
- Handlers
- Sessions
- Mailing
- Remote File Systems

  - S3
  - Minio
  - SFTP
  - WEBDAV

- Maintenance Mode
- Screen Capture of testing webpage

## Configurations

### Postgres database

1 Configure `/config/database.yml`

```yml
development:
  dialect: postgres
  database: gosnel
  user: postgres
  password: password
  host: localhost
  port: 5432
  pool: 5
```

2 Configure .env

```env
# database config - postgres or mysql
# 💁 make sure you stop postgres service on your local machine
# otherwise, it will try to connect to postgre db on your local machine, not on docker container
DATABASE_TYPE=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASS=password
DATABASE_NAME=gosnel
DATABASE_SSL_MODE=disable
```

3. Run `docker-compose`.

```shell
$ docker-compose up -d

# Check if all service containers are running fine
$ docker ps

# Start the application
$ make start

# To stop the application
$ make stop
```

## CLI Application

### Migrations

### Authentication

To make authentication functionality, run the follwoing in the project root:

```shell
$gosnel make auth
```

Then open up `models.go` and update Model struct and New() function as follows:

```go
// models.go
package data

import (
 "database/sql"
 "fmt"
 "os"

 db2 "github.com/upper/db/v4"
 "github.com/upper/db/v4/adapter/mysql"
 "github.com/upper/db/v4/adapter/postgresql"
)

var db *sql.DB
var upper db2.Session

// Models is the wrapper for all database models.
type Models struct {
 // any models inserted here (and in the New function)
 // are easily accessible throughout the entire application

 // For auth, these models are needed
 RememberToken RememberToken
 Tokens        Token
 Users         User
}

// New initializes the models package for use.
func New(databasePool *sql.DB) Models {
 db = databasePool

 switch os.Getenv("DATABASE_TYPE") {
 case "mysql", "mariadb":
  upper, _ = mysql.New(databasePool)
 case "postgres", "postgresql":
  upper, _ = postgresql.New(databasePool)
 default:
  // do nothing
 }

 return Models{
  RememberToken: RememberToken{},
  Users:         User{},
  Tokens:        Token{},
 }
}

// getInsertID returns the integer value of a newly inserted id (using upper).
func getInsertID(i db2.ID) int {
 idType := fmt.Sprintf("%T", i)
 if idType == "int64" {
  return int(i.(int64))
 }

 return i.(int)
}
```

Then open up `routes.go` file and add the routes as following:

```go
// routes.go
package main

import (
 "net/http"

 "github.com/go-chi/chi/v5"
 "github.com/youngjae-lim/gosnel"
)

func (a *application) routes() *chi.Mux {
 // middleware must come before any routes

 // add routes here
 a.get("/", a.Handlers.Home)

 // add these routes after running 'gosnel make auth'
 a.get("/users/login", a.Handlers.UserLogin)
 a.post("/users/login", a.Handlers.PostUserLogin)
 a.get("/users/logout", a.Handlers.Logout)

 a.get("/auth/{provider}", a.Handlers.SocialLogin)
 a.get("/auth/{provider}/callback", a.Handlers.SocialMediaCallBack)

 // static files
 fileServer := http.FileServer(http.Dir("./public"))
 a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

 // routes from gosnel
 a.App.Routes.Mount("/gosnel", gosnel.Routes())
 a.App.Routes.Mount("/api", a.ApiRoutes())

 return a.App.Routes
}
```

#### Social Login via Github and Google

Configure .env file for key and secret for both Github, and Google.

```env
# social login for github
GITHUB_KEY=<your_github_key>
GITHUB_SECRET=<your_github_secret>
GITHUB_CALLBACK=http://localhost:4000/auth/github/callback

# social login for google
GOOGLE_KEY=<your_google_key>
GOOGLE_SECRET=<your_google_secret>
GOOGLE_CALLBACK=http://localhost:4000/auth/google/callback
```

#### Add a Navigation for login and logout

Update `./views/layouts/base.jet` as follows:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta
      name="viewport"
      content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
    />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Gosnel: {{yield browserTitle()}}</title>

    <link
      rel="apple-touch-icon"
      sizes="180x180"
      href="/public/ico/apple-touch-icon.png"
    />
    <link
      rel="icon"
      type="image/png"
      sizes="32x32"
      href="/public/ico/favicon-32x32.png"
    />
    <link
      rel="icon"
      type="image/png"
      sizes="16x16"
      href="/public/ico/favicon-16x16.png"
    />
    <link rel="manifest" href="/public/ico/site.webmanifest" />

    <link
      href="//cdn.jsdelivr.net/npm/bootstrap@5.1.0/dist/css/bootstrap.min.css"
      rel="stylesheet"
      integrity="sha384-KyZXEAg3QhqLMpG8r+8fhAXLRk2vvoC2f3B09zVXn8CA5QIVfZOJ3BCsw2P0p/We"
      crossorigin="anonymous"
    />
    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.7.1/font/bootstrap-icons.css"
    />
    <meta name="csrf-token" content="{{.CSRFToken}}" />

    {{yield css()}}
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container-fluid">
        <a class="navbar-brand" href="#">Navbar</a>
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <a class="nav-link" href="/">Home</a>
            </li>
            <!-- File Systems Dropdown won't be working until you add file systems functionality to your app -->
            <li class="nav-item dropdown">
              <a
                class="nav-link dropdown-toggle"
                href="#"
                id="navbarDropdown"
                role="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                Filesystems
              </a>
              <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
                <li><a class="dropdown-item" href="/list-fs">List files</a></li>
                <li>
                  <a class="dropdown-item" href="/files/upload">Upload File</a>
                </li>
              </ul>
            </li>

            <li class="nav-item dropdown">
              <a
                class="nav-link dropdown-toggle"
                href="#"
                id="navbarDropdown2"
                role="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                Users
              </a>
              <ul class="dropdown-menu" aria-labelledby="navbarDropdown2">
                <li><a class="dropdown-item" href="/users/login">Login</a></li>
                {{ if .IsAuthenticated }}
                <li>
                  <a class="dropdown-item" href="/users/logout">Logout</a>
                </li>
                {{ end }}
              </ul>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container">
      <div class="row">
        <div class="col-md-8 offset-md-2">
          {{if .Flash }}
          <div
            class="alert alert-success alert-dismissible fade show mt-2"
            role="alert"
          >
            {{.Flash}}
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="alert"
              aria-label="Close"
            ></button>
          </div>
          {{end}} {{if .Error }}
          <div
            class="alert alert-success alert-dismissible fade show mt-2"
            role="alert"
          >
            {{.Error}}
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="alert"
              aria-label="Close"
            ></button>
          </div>
          {{end}} {{yield pageContent()}}
        </div>
      </div>
    </div>

    <script
      src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.1/dist/js/bootstrap.bundle.min.js"
      integrity="sha384-/bQdsTh/da6pkI1MST/rWKFNjaCP5gBSY4sEBT38Q/9RBh9AH40zEOg7Hlq2THRZ"
      crossorigin="anonymous"
    ></script>

    {{yield js()}}
  </body>
</html>
```

### Models

### Handlers

### Sessions

### Mailing

### Maintenance Mode

To go into maintenance mode,

```shell
$gosnel down
```

To go back into live mode,

```shell
$gosnel up
```

### Example of .env configuration

```env
# Give your application a unique name (no spaces)
APP_NAME=myapp
APP_URL=http://localhost:4000

# false for production, true for development
DEBUG=true

# the port should we listen on
PORT=4000
RPC_PORT=12345
ALLOWED_URLS="/login,/admin"

# the server name, e.g., www.mysite.com
SERVER_NAME=localhost

# should we use HTTPS?
SECURE=false

# database config - postgres or mysql
# 💁 make sure you stop postgres service on your local machine
# otherwise, it will try to connect to postgre db on your local machine,
# not on docker container
DATABASE_TYPE=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASS=password
DATABASE_NAME=gosnel
DATABASE_SSL_MODE=disable

# redis config
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_PREFIX=myapp

# cache: Redis or badger
CACHE=redis

# cooking seetings
COOKIE_NAME=myapp
COOKIE_LIFETIME=1440
COOKIE_PERSIST=true
COOKIE_SECURE=false
COOKIE_DOMAIN=localhost

# session store: cookie, redis, mysql, or postgres
SESSION_TYPE=redis

# mail settings
SMTP_HOST=
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_PORT=
SMTP_ENCRYPTION=
# domain is not required for sendgrid
MAIL_DOMAIN=
FROM_NAME=
FROM_ADDRESS=

# mail settings for API services
# For MAILER_API, please type smtp, sendgrid, mailgun, or sparkhost
MAILER_API=
MAILER_KEY=
# url is not required for SendGrid
MAILER_URL=

# template engine: go or jet
RENDERER=jet

# the encryption key; must be exactly 32 characters long
KEY=l7HrN539xnhnNvjFfVv5NF9ThfLdl5rX

# remote file system variables for s3, minio, sftp, and webdav

# S3 object storage for linode
S3_SECRET=<your_s3_secret>
S3_KEY=<your_s3_key>
S3_REGION=us-east-1
S3_ENDPOINT=us-east-1.linodeobjects.com
S3_BUCKET=gosnel

# Minio
MINIO_ENDPOINT=127.0.0.1:9000
MINIO_KEY=root
MINIO_SECRET=password
MINIO_USESSL=false
MINIO_REGION=us-east-1
MINIO_BUCKET=testbucket

# SFTP
SFTP_HOST=localhost
SFTP_USER=sftp
SFTP_PASS=password
SFTP_PORT=2022

# WEBDAV
WEBDAV_HOST=http://127.0.0.1:10080
WEBDAV_USER=sftp
WEBDAV_PASS=password

# permitted upload types
ALLOWED_MIMETYPES="image/gif,image/jpeg,image/png,application/pdf"
# max upload size 10mb
MAX_UPLOAD_SIZE=1048576000

# social login for github
GITHUB_KEY=<your_github_key>
GITHUB_SECRET=<your_github_secret>
GITHUB_CALLBACK=http://localhost:4000/auth/github/callback

# social login for google
GOOGLE_KEY=<your_google_key>
GOOGLE_SECRET=<your_google_secret>
GOOGLE_CALLBACK=http://localhost:4000/auth/google/callback
```
