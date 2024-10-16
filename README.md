# [ORDERBOOK]

Orderbook - RabbitX

## Project Assumptions
- Currency pass into the service is assumed all the same
- Assuming with limit, market order that identified through maker-taker model.
- Assuming free to use, no transaction fee charge :).
- Assuming users have diamond hand, baghold no matter what position, no closed position endpoint provided :).
- Assumimng simplified order book, long when user have short position does not close the short position and 


## Projet Overview

This project is built using the Gin framework, which handles concurrent requests by default. Hence, the trade (write) order operation is exposed as rest api endpoint.

It supports a **maker-taker model** for order execution, as this is a simple orderbook some assumptions is given.

In this model:

- Maker: A maker is a trader who provides liquidity to the market by placing limit orders. These orders are set at a specific price and sit on the order book until matched with a market(taker) order. Makers contribute to the market's depth, and they often benefit from lower fees due to their role in enhancing liquidity.

- Taker: A taker, on the other hand, is a trader who places market orders that are executed immediately at the best available price. Takers remove liquidity from the market, as they fulfill existing orders rather than adding new ones. They typically pay higher fees than makers because their trades can create volatility.

- Self match is not allowed, it will skip to next best available price.

## How does the matching engine logic works (high level)?:

1. Matching engine is triggered on every order submitted. Every order requested via the `POST /orderbook/orders/insert` endpoint will first pass into the matching engine after going through the first level of simple validation below.
    - Retrieve current's position to check the current state, proceed to calculate net position size if open position is at opposite direction.
        -- With that say, placing reverse direction order will closed your current position and realized PNL.
        -- If the order position amount is larger than open position at reverse direction, all open position PNL will be realized and proceeding to fulfiled the net position. 
    - Retrieve the current user's wallet balance and validate if the wallet balance is equal to or more than the submitted request amount.
2. The order is then passed into the matching engine and sent to the respective matching engine based on the order's side (1 is long, -1 is short).
    - Two matching engine logics are created. The logic is similar with slight differences that are hard to realize; hence, two functions are created to improve readability and future modification.
    - A panic handler is added to handle any unforeseen scenarios (technically).
3. In the matching engine (`/server/pkg/tarantool_pkh/matching_engine.go`), the order book for the respective side is retrieved.
    - For long orders, the ask order book will be retrieved according to the market (eth/btc/etc).
    - For short orders, the bid order book will be retrieved according to the market (eth/btc/etc).
    - The orderbook data is extracted and rebuilt into an **array of arrays**, with ascending price as the first layer and ordered timestamp as the second layer. Please refer to the example below:
    ```
    e.g.
    [[order1($12.1),order2($12.1)],[order1($13.1),order3($13.1)]]
    ```
    - With this data structure:
        - Long orders: iterate the array through ascending order to match the lowest available best ask price.
        - Short orders: iterate the array in reverse direction (last element) to match the highest available bid price.
    - **(Caveat)** Due to simple implementation, the order book is stored in an in-memory DB (Tarantool) in the form of each order being a row, holding details like side, user ID, quantity, price, market, etc. As mentioned, this will then be extracted and rebuilt into an array.
    - **(Improvement)** Since Tarantool is an in-memory database, a future improvement could involve converting this storage method into a hashmap for better and more efficient data retrieval and price identification (faster matching engine).
4. In the matching engine, the following steps occur:
    - The current order will be identified as either a **taker/maker** order by checking the **current order** price against the **best available price**. It will switch to a maker order (limit order) depending on the scenario.
        - Long order: current `order price < lowest asking price` **(it's a maker order)**, insert into the order book, deduct the amount from the user's wallet balance, end of the engine.
        - Short order: current `order price > highest bidding price` **(it's a maker order)**, insert into the order book, deduct the amount from the user's wallet balance, end of the engine.
    - If this is a **taker order** (matching engine does not end at the step above), iterate over the orderbook array mentioned in Step 3.
    - If this is a **taker order** (matching engine does not end at the step above), the matching engine will try to fulfill all the requested position size as long as the current order price and maker price conditions are valid.  `The engine will stop when`:
        - All positions in the current order are fulfilled.
        - Long order: the engine stops when the current `order price < lowest asking price`.
        - Short order: the engine stops when the current `order price > highest bidding price`.
    - The execution price will always be based on the maker order (limit) from the order book.
    - In every fulfillment cycle:
        - Maker: 
            - If the open order position is fully filled, remove it from the order book and open active positions in the position list.
            - If the open order is partially filled, update the position size in the current order list and the timestamp remains unchanged (FIFO).
            - The order will only requeue (update timestamp) if the same order is being modified from the maker's end.
        - Taker:
            - If partially fulfilled, open the position in the position list, deduct the position size and wallet balance, and use it for the next iteration.
            - If fully filled, the engine ends with the open position in the position list, deduct the position size, and the wallet balance.
        - Note:
            - The state of position list and orderbook will be updated every iteration when there is matched.
            - The logic will update the current position by recomputing the average price and increasing the position size holdings if the position is already present in the position list with the same `side and market`.
            - Account margin (>10%) is checked every cycle to ensure there is enough to execute (match) the order. `When the account margin check fails:`
                - Taker: revert the order, not executing and stop the engine.
                - Maker: revert all the current order in order book, and refund the initial deducted amount when writing into orderbook (taker does not need refund because balance is not deducted yet).
                - Note: When calculating the account margin, it is based on the condition after the order is executed, meaning it includes the updated wallet balance and position holdings.
            - User wallet is checked together with account margin to ensure there is enough balance to cover the current execution price.
    - The market price of the particular instrument/market (eth, btc, etc.) will be updated based on the latest traded price.

Other Note:
1. When order cancelled through cancelling endpoint, user gets refunded. Likewise, refund also happen when margin less than 10% and all the order book orders is automtically cancelled.
2. Please refer to other [endpoints](#api-endpoints) that mimic actual trading, hence, you can perform action like inserting order, deposit money to test the orderbook engine.
3. Please refer to tarantool setup to understand how to reset the db as some of the edge cases endpoint might not be provide for you to delete data etc [tarantool docs](#tarantool-db), alternatively, you can refer to [tarantool docs](#tarantool-db) the cli command below to check how to connect into db and modified the data through cli.

## Technical Views
- How is the project structre
- when use imemorydb
- stateless 
-atomicity
- precision issue.

## API Endpoints

This project provides the following REST API endpoints, please refer to `orderbook_postman_collection.json` on the endpoint contract details like (request body):

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

