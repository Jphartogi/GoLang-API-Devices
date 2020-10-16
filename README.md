# GoLang-API-Devices

CRUD App for registering smart devices and integration with username and password. All data is stored in NO-SQL database ( mongodb )

# How to run

First build the image for the api-app
```bash
docker build -t api-app:1.0 .
```

Then to run the services using the compose command

```bash
docker-compose up -d
```

# How to check the services is up and running

The grpc server is running on port 5000 tcp

And also provided a http listener on port 8092 to check if the api is running and return simple json

To run 
```
curl http://localhost:8092
```

or simply open your browser and type http://localhost:8092

# To run manually ( without docker ) and the test script

## Setting up database
First you need to make sure mongodb is up and running in your system, easily just use 
```
docker run -d -p 27017-27019:27017-27019 --name mongodb mongo
```

## Run the main script
```
go run main.go
```

## Then to run the test script

```
go test -run main_test.go
```
