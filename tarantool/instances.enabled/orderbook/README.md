# Connectors

A sample application used to demonstrate how to connect to a database using connectors for different languages and execute requests for manipulating the data.

## Running

Start the application by executing the following command in the [connectors](../../../connectors) directory:

```console
$ tt start sample_db
```

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

## Spaces
1. User (id=integer, wallet_balance=number, )
2. positions (user_id=integer, market=string, position_size=number, avg_entry_price=number, side=integer, status=string)
3. market_price (market=string, price=number)
4. order_book (price=number, market=string, side=integer,userid,entry_price, position_size)
5. match_history..

#side is 1 (long) and -1 (short)