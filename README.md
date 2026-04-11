# 🐥 Chirpy
Twitter-like app for [boot.dev](https://boot.dev) course on http servers in Go.
I mainly did it to improve my golang skills and general http webserver knowledge

## What it is:
It's basically mimicing a site that allows posting short messages to a social media platform, similar to many others that exist, but much more bare-bones. It supports basic CRUD operations and implements some basic Authentication/Authorization mechanisms via JSON Web Tokens [JWTs](https://en.wikipedia.org/wiki/JSON_Web_Token). Essentially, this means users need to be authenticated before posting "chirps", they cannot edit other users' email/password entries, they cannot delete "chirps" that do not belong to them, etc. They also do not need to pass their password in the request body of each POST request; that information is stored in the HTTP Headers.

## Prerequisites:
* Golang 1.25 or higher
* Postgresql 15
* A .env file to store environment variables

## 🛠️ Installation:
Install with `go get https://github.com/beerded/Chirpy` (yeah I know, it shouldn't be capitalized)

### Installing Go:
On macOS with `brew install go`

### Installing Postgres:
- On macOS with `brew install postgresql@15`
- Create a new database for the project: `CREATE DATABASE chirpy;`
- Run all the migrations under `Chirpy/sql/schema` to get the database initialized properly

## Environment Variables
This project makes use of several environment variables, namely:
- "DB\_URL": Connection string for your database, e.g. `postgres://username:@localhost:5432/chirpy?sslmode=disable`
- "PLATFORM": Set it to "dev" if you want to be able to reset the database. Otherwise that endpoint is disabled.
- "JWT\_SECRET": This is the secret string that is used to identify individual users and create their JWTs.
- "POLKA\_KEY": Used for authenticating the webhook clients for the endpoint that upgrades users to "Chirpy Red" mode

## 🧙 API Documentation

### **Chirps** 
Interact with microblogging posts, known as "Chirps"
##### **GET `/api/chirps`**
Retrieve multiple chirps
###### Parameters
| Name           | Optional (Y/N) | Description                   |
| -------------- | -------------- | ----------------------------- |
| **author_id**<br> string UUID<br> *(query)*        | Y              | Allows filtering by author-id (default is to list all chirps) |
| **sort**<br> string \["asc","desc"\]<br> *(query)* | Y              | Sort results in ascending or descending order (default is asending)|
###### Responses
Content-Type: `application/json`
<table>
  <tr>
    <td>Code</td>
    <td>Description</td>
  </tr>
  <tr>
    <td>200</td>
    <td>Successful Operation</br>Example output:
      
```json
{
  [
    {
      "id": "8bb3f27d-1d5d-4be7-8913-9e792544be39",
      "created_at": "2020-01-01 09:38:21.325756",
      "updated_at": "2020-01-01 09:38:21.325756",
      "body": "Happy New Year everyone!",
      "user_id": "98d1c380-fccc-4834-99cb-ce96a7d26aba"
    },
    {
      "id": "8bb3f27d-1d5d-4be7-8913-9e792544be39"
      "created_at": "2020-01-01 09:43:37.315756",
      "updated_at": "2020-01-01 09:43:37.315756",
      "body": "Happy New Year to you too Steve!",
      "user_id": "7adfc380-aeec-1894-44aa-ce96a7d26aba"
    }
  ]
}
```
   </td>
  </tr>
  <tr>
    <td>404</td>
    <td>Not Found</br>Example output:

```json
{
  "error": "Could not find user"
}
```
  </tr>
</table>
    
### Endpoint: GET `/api/chirps/{id}`
**Retrieve a single chirp by ID**
### Endpoint: POST `/api/chirps`
**Create a chirp**
### Endpoint: DELETE `/api/chirps/{id}`
**Delete a single chirp by ID**

### Endpoint: POST `/api/login`
**Login as a particular user**

### Endpoint: PUT `/api/users`
**Edit your user data**
### Endpoint: POST `/api/users`
**Create a new user**

### Endpoint: POST `/api/refresh`
**Refresh your access token**
### Endpoint: POST `/api/revoke`
**Revoke your refresh token**

### Endpoint: POST `/api/polka/webhooks`
**Send an event to the chirpy backend via webhook**

