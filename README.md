# Simple Bank Project

## 1. Introduction
This project provides APIs for:
- Create a new bank account and manage it
- Record the history of transactions
- Perform a transaction (money transfer) between two accounts

Technologies
- Programming languages: Golang
- Postgres, Gin, Docker, CI/CD, MockDB

## 2. Services
The service that we’re going to build is a simple bank. It will provide APIs for the frontend to do following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.


## 3. Setup local development

### Install tools

- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    sudo apt install golang-migrate
    ```


- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    sudo apt install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### How to generate code


- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```



### How to run

- Run server:

    ```bash
    make server
    ```

## 4. Future features
- Change User Password API
- Email Verification services
