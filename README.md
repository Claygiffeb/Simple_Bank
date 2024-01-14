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
![Alt text](image.png)

## 3. Database Transaction

The most important transaction in this project is tranfering money. Although it's sound simple, but in order to provide reliable and consistent, it need to follow the ACID property.

For example: The transaction transfer 10 USD form account A to account B will compose below steps
- Step 1: Create a record of the transaction with amount = 10
- Step 2: Create an entry account for A with amount = -10
- Step 3: Create an entry account for B with amount = 10
- Step 4: Subtract 10 from the balance of A
- Step 5: Add 10 to the balance of B

Note the Step 4 and Step will require locking protocol