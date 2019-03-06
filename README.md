# JWT authentication system

A very simple authentication system based on JWT (JavaWebTokens).

# Installation

This application uses a MySQL (or MariaDB) database engine. In order to test the application you have to create a database and then import the `database.sql` script located in the application directory. You can do it from many programs, such as phpMyAdmin, Adminer, etc. Or you can do it directly from the command line. For example:

```bash
$ mysql u -root -p

mysql> create database mydb;
mysql> use mydb
mysql> source dbschema.sql
```

Then rename or copy `example.env` to `.env` and change the Token Database sections. After that you can now compile and execute the application:

```bash
# compile and execute the application
$ go build
$ ./jwtserver
2019/03/06 21:57:24 Server started at port localhost:8080
```

You can now test this application from Postman:  
https://documenter.getpostman.com/view/412470/S11PpFqg

And that's all Folks!
