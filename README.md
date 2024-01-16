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

Detail implementation in the util package (password.go) 

## 8. Authentication: Token-based authentication

In this project, we use JWT v5 for authentication with login user and password, more details on Security is in the Appendix

## 9. Middleware Authentication and Authorization rules

In banking systems, it is crucial that the authorization rules are forced (so that Alice can not tranfer money from Bod to Alice?). This project create a middleware to handle that business logic.
![Alt text](image-20.png)
![Alt text](image-21.png)

We will implement role authorization to enforce four below rules:
![Alt text](image-22.png)


## Apendix: Security - SSL/TLS overview

![Alt text](image.png)

### 1. Symmetric Cryptographic

Suppose Alice want to send a message to Bob, they will use their Shared secret key to encrypt and decrypt the message.

![Alt text](image-7.png)

### 2. Authenticated Encryption
In cases there is a hacker that take the message, he cant decrypt the message, but he can change the encrypted data using Bit-Flipping attacks.

![Alt text](image-8.png)

So, Symmetric Cryptographic is not enough, the idea is we will authenticate the message. The idea is to combine the message with MAC (Message Authentication code)

![Alt text](image-10.png)

Then, Bob can decrypt the message. First he Untag the Encrypted message with MAC to get the MAC and the Encrypted message. Then he run MAC algorithm to generate Tag information from the Shared secret key and compare it with the MAC untagged before.

If the tags are equal, he will encrypt the message to get to message from Alice. Otherwise, he knows that the message has authentication problems.

![Alt text](image-11.png)

### 3. Asymmetric Cryptography

A nature question might come, how Alice and Bob share the secret key before exchange message with out leaking it. The idea is to use public key and private key:
- Diffie-Hellman Ephemeral (DHE)
- Elliptic Curve Diffie-Hellman Ephemeral (ECDHE)

![Alt text](image-12.png)

The idea is kinda like "color mixing", but with some math. 

What about ECDHE? Maybe later

![Alt text](image-13.png)

When Alice want send the message to Bob, she will use the Bob's public key to encrypt the message, since only Bob's private key can decrypt it, we are save (?). 

But what about this scenarios?
![Alt text](image-14.png)

Suppose a hacker can replace the  key to his key, then he can decrypt the message.

We need other ways to authenticate the key sharing too.


### 4. Digital Signature 

The last problem of key replacement: How can we ensure the authentication of the message sender -> The idea is kinda like the real life situation, we use another Source of Trust (like a government database?) to prove that the message was signed by the sender. 

![Alt text](image-15.png)

The steps are simple, we can think of it like signing in real life and go to the government so that they will prove that we did the signature

After that. Bob can send the certificate 

![Alt text](image-16.png)

But the problem is the CA (Certificate Authority) must be trustful! (Like it's not going to send the private key to the Hacker?)

![Alt text](image-17.png)

### 5. Token-based Authentication: Paseto vs JWT

![Alt text](image-4.png)

1. The client requests to login, if it is correctly authenticated, the server send a access token to the Client

2. The client request with the token in the head of request.

Each token has a duration

#### 1. Json Web Token

![Alt text](image-5.png)

It combines 3 fields that are seperated by a dot. Note that the last field contains a digital signature, which is used to verify the client request. 

![Alt text](image-6.png)

But what problem with JWT? 
As we can see, the developer may choose many different signature algorithms, which might lead to problems?

