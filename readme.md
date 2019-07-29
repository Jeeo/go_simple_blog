## Simple blog

It's a quite simple blog application that I've used to learn the basics of Golang web developemnt

It's does not properly handle with errors 

## Running

#### Pre-requisites
You must have the postgres installed to run this. the database configs are in data/data.go.

All tables schemas are in schemas path

Please use all.sql, but if you need of individual schema, you can use
the individuals inside of schemas path

run inside of project folder: 

```sh
$ go build 
```
so just starts with

```sh
$ ./blog 
```

the default port that this application is running is 3000

so go to the browser and access http://localhost:3000