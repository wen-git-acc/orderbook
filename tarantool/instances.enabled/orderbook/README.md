# Connectors

Tarantool DB as in-memory database for simple order book


## Running
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

## Connect to db for data Viewing
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


## Spaces
Total 4 Spaces creates:
1. users, hold userid and wallet balance
2. positions, hold current opening positions
3. market_price, hold market price for the market (eth, btc etc.)
4. order_book, hold all orders data.

