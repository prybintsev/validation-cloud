# Ethereum API

## Prerequisites

Running this service requires SQLite: https://www.sqlite.org/
It is usually included in MacOS, but the service run on a different OS, please make sure it is installed

JWT_SECRET_KEY environment variable should contain the JWT secrete key. Any string will work, although in real life
it should be a long randomly generated string. This key is used to generate JWT tokens.

When the service starts, it creates the database file (if it does not exist already) in data directory under the root
project directory.

## Setup
