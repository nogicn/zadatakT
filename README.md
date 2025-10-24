# Project backendT

Backend and mock implementation for the T company hackathon.
Didn't have time to make the frontend (applied for the backend part anyway), and the backend is done.

I tried using copilot/chat as little as possible and mostly used it to help me debug problems.
Most of the problems that made me lose time were connected to me not being familliar with how Echo handles bindings/parsing of data and setting up swagger with Echo and it took me around 4 hours just to set everything up(swagger, db, echo parsing).
I usually use either net/http with gorilla/mux or gofiber, challenged myself to try Echo in a hackathon without using it heavily before.

The base template that has the folder structure and a basic Echo server was made with 
```url
https://github.com/Melkeydev/go-blueprint
```
while the rest (database layer, models, endpoints, logging middleware, T app integration) was made by me.

## Notable features

The backend is written in Echo
```url
https://echo.labstack.com/docs
```
and it does two things.
It uses two models to simulate a small CRUD app (/posts and /users endpoint)
and it also logs all requests it has processed (/logs endpoint) including requests for logs themselves.

Its possible to use multiple filters on the logs endpoint so the frontend doesnt need to do too much work. Example curl requests can be found in routes.go in case I don't have enough time to finish the frontend or you can run them in swagger.

Swagger runs on /swagger/index.html 
```url
https://waps.website/swagger/index.html
```
or if run locally
```url
http://localhost:8080/swagger/index.html
```

By default, there is a single account and post made so you can test the application with swagger right away.

The application can also be connected to the T app with a simple wrapper I made and it can be seen in routes.go.
You only need to add the variables to the env as per instructions on the T company dashboard and the rest will work like magic!

The application uses sqlite for the local database.
Because the go sqlite implementation doesnt pair well with multiple writers, two differerent connections are made to the db.
One is Read-only and the other one is Read-Write but limited to one (1) writer because using multiple connections to write will severely throttle the sqlite implementation. 

For the management of the database and migrations, the project uses Goose.
```url
https://pressly.github.io/goose/
```

This project doesnt use an orm and instead uses Sqlc as the driving force behind the project.
```url 
https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#
```

The reason for choosing Sqlc is the ability to get an experience very close to a full fledged ORM, but with the benefit of seeing and writing your own queries and knowing exactly what they will make/do.
Sqlc automatically creates functions that can be used for scanning data from the database and structures needed for arguments/results, and they can be seen in the repository folder with the name format *.sql.go, while all the queries for every table can be seen in the queries folder.

This backend is deployed on ```https://api.waps.website```

The application has integration tests for endpoints and unit tests for the database repository layer.



## Running the project

### You can run the application using docker or using the Makefile just remember to rename the .env.example file to .env.
The application defaults to port 8080 but can be changed in the .env file.


Below are all possible commands with the application.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create docker container
```bash
make docker-run
```

Shutdown docker Container
```bash
make docker-down
```

Live reload the application using air:
```bash
make air
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
