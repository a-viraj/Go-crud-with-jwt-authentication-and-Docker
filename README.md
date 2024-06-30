
# Golang Back-end Server

This is a Back-end server in golang, using the Gin framework and Go ORM with MongoDB database.


## Authors

- [@a-viraj](https://www.github.com/a-viraj)


## Implementations

* HTTP Request
* Authorization and Authentication
* Middleware (Auth, Logger)
* MONGODB
* Containerization
## How to run using Docker Container

Ensure you have docker and Docker compose installed
 ```
 docker compose up -d
 
 ```

 ## How to run on a local environment

 * Ensure you have mongodb Compass installed
 * create a .env file 
 * Update the DB connection string in the .env file
 * Install packages `go mod tidy` 
 * You can run the project in the terminal use `go run main.go`
 
## Required package

Reference the `go.mod` and `go.sum` file
