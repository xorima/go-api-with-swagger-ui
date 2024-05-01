# Demo API

A basic API to demo/play with golang http apis. 

## Endpoints

It is known by me that exposing the user id is a bad idea, but this is a demo api 
and I am not going to implement any kind of real security.

`/users`
- Method: GET
  - Response: `{"name": "foo", "id": 1}`
- Method: POST
  - Input: `{"name": "foo"}`
- Method: PUT
  - Input: `{"name": "foo", "id": 1}`
- Method: Delete
    - Input: `{"id": 1}`
`/users/{id}`
    - Method: GET
    - Response: `{"name": "foo", "id": 1}`

Authentication: 
In this example we will use a basic Bearer token authentication, where the token is `1234`.

Observability: 
We will use prometheus to monitor the api.

Logging:
Handled via slog library.

Database: 
Postgres for the database.

## Updating swagger

Remember to run swag init each time. Later put this into the precommit hooks
