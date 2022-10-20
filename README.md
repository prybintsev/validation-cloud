# Ethereum API

## Prerequisites

Running this service requires SQLite: https://www.sqlite.org/
It is usually included in MacOS, but the service run on a different OS, please make sure it is installed

JWT_SECRET_KEY environment variable should contain the JWT secrete key. Any string will work, although in real life
it should be a long randomly generated string. This key is used to generate JWT tokens.

When the service starts, it creates the database file (if it does not exist already) in data directory under the root
project directory.

## Running the application
To run the application, execute the following command in terminal

The application includes a web server an a job that collects samples once a minute. In order to be able to retrieve the
average growth rate, the application needs to run for one minute or more so at least 2 samples are collected.

```
JWT_SECRET_KEY="my_secret_key" make run
```
## Examples

### Creating a user

To create a user, run the following command:

```shell
curl --location --request POST 'http://localhost:8080/auth/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user-name": "myusername",
    "password": "mypass"
}'
```

You should receive the following response:

```shell
{
    "code": 200,
    "message": "ok"
}
```

### Generating JWT token

To generate a JWT token, run the following command using the credentials from the previous command

```shell
curl --location --request POST 'http://localhost:8080/auth/generate-token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user-name": "myusername",
    "password": "mypass"
}'
```

You should get a response with the generated token:

```shell
{
    "code": 200,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjYyNDE2NzEsInVzZXIiOiJteXVzZXJuYW1lIn0.wC5Hgeg91Tlaoxgpoa4SrOKq7nIeOTAlF01KIEIQEFQ"
}
```

The token is valid for a period of 10 minutes after which a new token needs to be generated.

### Retrieving the average blockchain growth

To get the average growth of the blockchain per minute, run the following command using the token from the token generation response:

```shell
curl --location --request GET 'http://localhost:8080/geth/avg-growth' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjYyNDE2NzEsInVzZXIiOiJteXVzZXJuYW1lIn0.wC5Hgeg91Tlaoxgpoa4SrOKq7nIeOTAlF01KIEIQEFQ'
```

You should get a response with the average number of block generated per minute
```shell
{
    "AverageGrowth": 5.00288133919318
}
```
