# [ORDERBOOK]

Orderbook - RabbitX

## Projet Overview

This project is built using the Gin framework, which handles concurrent requests by default. It supports a maker-taker model for order execution. In this model:

Taker Orders: When an aggressive taker order is placed, it will continue to fulfill the position until the best price condition fails. For short positions, the execution stops when the bid price falls below the short price. For long positions, it stops when the ask price exceeds the long price.
Order Book Management: If an instant match is not found, the order is recorded in the order book. Once an order is placed in the order book, the wallet balance is deducted accordingly.
Matching Charges: For taker orders, a fee is charged each time a match is found.

Note:
1. All the inserted market name and userid is assume to be small letter.
2. Assume all the user have permission to access and start using trading endpoints.
3. When order cancelled, user gets refunded. Likewise, refund also happen when margin less than 10% and all the order book orders is automtically cancelled.

## API Endpoints

This project provides the following REST API endpoints:

- **Health Check**
  - `GET /health_check`: Check the health status of the service.

- **Order Management**
  - `POST /orderbook/orders/insert`: Insert a maker or taker order.
  - `POST /orderbook/orders/cancel`: Cancel an order from the order book.

- **User Management**
  - `POST /orderbook/user/deposit`: Make a deposit as an existing or new user.
  - `GET /orderbook/user/wallet/:userId`: Get the current wallet balance of a user.
  - `GET /orderbook/user/:userId/positions`: View the current open positions (matched) of a user.

- **Market Data**
  - `GET /orderbook/:market`: Get the order book for a specific market (e.g., eth, btc).
  - `GET /orderbook/market-price/:market`: Get the current market price of a specific market (e.g., eth, btc).

- **Position Viewing**
  - `GET /orderbook/view/positions`: View all open positions regardless of the user.


## Project Navigation
- navigation to entry point /server/api/controllers/orderbook_controller.go 
- navigation to main.go /server/cmd/gin/main.go

## Set Up Environment
1. Copy and paste the launch.json schema below, and run go cli (triggers /server/cmd/cli/main.go) to pre-poluate user with wallet balance and orderbook for local test.
2. Alternatively, you can run the script file by navigating to the mentioned directory with 
```sh
go run .
```
3. You can run the script file as many time you wish to increase position size in each order, deposit will remain unchanged unless the value is reconfigured.
4. Please refer guide for tarantool below if you are looking to clean/reconfigure the database.
5. Lastly please import the postman collection for easy set up if needed :).

## Test Command

Navigate to ./server/tests and run command below

```sh
go test ./...
```

## Project structure

    ├── server                 - Application source code
    │   ├── api                - api related code (controleers, handlers, dto, middleware...)
    │   ├── cmd                - Go appliction entry point
    │   ├── config             - Application config from env.
    │   ├── pkg                - Packages or Dependencies
    │   ├── tests              - Test files.  
    ├── scripts                - Folder dedicated for bash script or others.
    ├── tarantool              - Tarantool environment config files.
    ├── Dockerfile             - Dockerfile for building the image
    └── README.md              - Current view.

## Debugging code in VS Code

Create `launch.json` under the `.vscode` folder in the root directory of the project. Add the following configurations:

This allowed you to run the project and debug the code in Visual Studio Code.

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.1.0",
    "configurations": [
      {
        "name": "Launch Go",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "env": {
          "DEBUG": "true",
          "MODE": "development",
          "SECRET_POSTGRES_DB_PASSWORD": "postgres_password",
          "TARANTOOL_PASSWORD": "123456"
        },
        "program": "${workspaceFolder}/server/cmd/gin"
      },
      {
        "name": "Launch Go Cli",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "env": {
          "DEBUG": "true",
          "MODE": "development",
          "SECRET_POSTGRES_DB_PASSWORD": "postgres_password",
          "TARANTOOL_PASSWORD": "123456"
        },
        "program": "${workspaceFolder}/server/cmd/cli"
      }
    ]
}
```
## Environment Configuration

Staging and Production might not be useful as we could be using configmaps and secret for k8s.

## Setup local postgresql database

### Start postgresql in your local

Please change the launch.json and environment variable respectively if you have set up user and password for your db.

```shell
$ brew install postgresql

$ brew services start postgresql
```

## Tarantool DB
### Connectors

Tarantool DB as in-memory database for simple order book

### Installation
```sh
brew install tarantool
```

```sh
$ brew install tt
```
### Running
Navigate to ./tarantool from root.

Interactive mode:
```sh
tt start orderbook -i
```

Normal mode:
```sh
tt start orderbook
```

Clean data for fresh new instances
```sh
tt clean orderbook
```

### Connect to db for data Viewing
1. 
```sh
tt start orderbook
```

2.
```sh
tt connect orderbook
```

3.
```sh
box.space.users:select{}
```
```sh
box.space.market_price:select{}
```
```sh
box.space.positions:select{}
```
```sh
box.space.order_book:select{}
```

### Spaces
Total 4 Spaces created:
1. users, hold userid and wallet balance
2. positions, hold current opening positions
3. market_price, hold market price for the market (eth, btc etc.)
4. order_book, hold all orders data.

