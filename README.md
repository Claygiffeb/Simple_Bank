# Simple Bank Project

## 1. Introduction
This project provides APIs for:
- Create a new bank account and manage it
- Record the history of transactions
- Perform a transaction (money transfer) between two accounts

Technologies
- Programming languages: Golang
- Postgres, Redis, Gin, gRPC, Docker, Kubernetes, AWS, CI/CD
## 2. Database Design
![Alt text](image-1.png)

- accounts: Table of accounts
- entries: Entries in transaction
- transfers: Transfers between two accounts
- users: owners of the accounts, one user can have many accounts, which have different currencies. To enforce this constraint, we will create a unique composite index. 

## 3. Database Transaction

The most important transaction in this project is tranfering money. Although it's sound simple, but in order to provide reliable and consistent, it need to follow the ACID property.

For example: The transaction transfer 10 USD form account A to account B will compose below steps
- Step 1: Create a record of the transaction with amount = 10
- Step 2: Create an entry account for A with amount = -10
- Step 3: Create an entry account for B with amount = 10
- Step 4: Subtract 10 from the balance of A
- Step 5: Add 10 to the balance of B

Note the Step 4 and Step will require locking protocol

## 4. Testting
Here we use MockDB. Mock databases are used to mimic the behavior of real databases without actually interacting with a live database system. This is advantageous in testing scenarios to ensure that tests are predictable, repeatable, and independent of external dependencies.

We will use reflection mode here

## 5. Validator

Since we don't want to hardcode the validation of requests (in the future, there might be many currency and we need to include them everywhere in the required configuration). 

## 6. Database Migrations
 
 Since we can modify the databases many times in the future, it's not good to change the database migration files, we need to create a new migration file for each migration.

## 7. Password Secure Storage

For security reasons, we will hash the password and store it using BCRYPT HASH (COST,SALT) algorithm

![Alt text](image-2.png)

Then when the user enter the password, we will hash that password using the same algorithm with the same (COST,SALT) and compare two hash value

![Alt text](image-3.png)

Detail implementation in the util package (password.go) and 