# Project backendT

Backend and mock implementation for the T company hackathon.

## Notable features

The backend does two things.
It uses two models to simulate a small CRUD app (/posts and /users endpoint)
and it also logs all requests it has processed (/logs endpoint) including logs requests themselves (inception :D).

Its possible to use multiple filters on the logs endpoint so the frontend doesnt need to do too much work. Example curl requests can be found in routes.go in case I don't have enough time to finish the frontend or you can run them in swagger.

The specific details of the api can be seen on 
```https://waps.website/swagger/index.html```

The application can also be connected to the T app with a simple wrapper as seen in routes.go.
You only need to add the variables to the env file as per instructions on the T company dashboard and the rest will work like magic!

The application uses sqlite for the local database.
Because the go sqlite implementation doesnt pair well with multiple writers, two differerent connections are made to the db.
One is Read-only and the other one is Read-Write but limited to one (1) writer because using multiple connections to write will severely throttle the sqlite implementation. 

For the management of the database and migrations, the project uses Goose.

This project doesnt use an orm and instead uses Sqlc as the driving force behind the project.

The reason for choosing Sqlc is the ability to get an experience very close to a full fledged ORM, but with the benefit of seeing and writing your own queries and knowing exactly what they will make/do when run.
Sqlc automatically creates functions that can be used for scanning data from the database and structures needed for arguments/results, and they can be seen in the repository folder with the name format *.sql.go, while all the queries for every table can be seen in the queries folder.

This backend is deployed on ```https://api.waps.website```,
and the frontend will be deployed on ```https://waps.website```




## Running the project

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

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
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
