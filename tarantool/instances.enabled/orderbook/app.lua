---@diagnostic disable: lowercase-global
-- Require the box module --
local box = require('box')

-- Create a space and Indexes--
-- Create market price spaces --
box.schema.space.create('market_price', { 
    if_not_exists = true,
    format = {
        { name = 'market', type = 'string' }, -- Market name (BTC or ETC)
        { name = 'price', type = 'number' }   -- Market price
    }
})
box.space.market_price:create_index('primary', { 
    type = 'hash',
    parts = { 1, 'string'},  -- Market name
    if_not_exists = true 
})


-- Create Users spaces --
box.schema.space.create('users', { 
    if_not_exists = true,
    format = {
        {name = 'id', type = 'string'},           -- User ID
        {name = 'wallet_balance', type = 'number'} -- Wallet balance
    }
})
box.space.users:create_index('primary', {
    type = 'hash',
    parts = {1, 'string'},  -- User ID
    if_not_exists = true
})


-- -- Create Positions spaces --
-- box.schema.create_space('positions', {
--     if_not_exists = true,
--     format ={
--         {name = 'user_id', type = 'integer'},          -- User ID
--         {name = 'market', type = 'string'},            -- Market (e.g., 'BTC', 'ETH')
--         {name = 'position_size', type = 'number'},     -- Position size
--         {name = 'status', type = 'string'},            -- (open/closed)
--         {name = 'avg_entry_price', type = 'number'},   -- Average entry price
--         {name = 'side', type = 'integer'}              -- Side (1 for long, -1 for short)
--     }
-- })
-- box.space.positions:create_index('user_market_index', {
--     type = 'hash',
--     parts = {1, 'integer', 2, 'string'},  -- user_id, market
--     if_not_exists = true
-- })


-- Create order book spaces --
box.schema.create_space('order_book', {
    if_not_exists = true,
    format =
    {
        {name = 'id', type = 'string'}, -- uniqueKey (userid:price:side) entry price
        {name = 'price', type = 'number'}, -- entry price
        {name = 'market', type = 'string'}, -- BTC or ETH
        {name = 'side', type = 'string'},  -- 1 for buy, -1 for sell
        {name = 'user_id', type = 'string'}, -- user id
        {name = 'position_size', type = 'number'}, -- position holds
        {name = 'created_at', type = 'unsigned'} -- Use a timestamp or an integer for ordering
    }
})
box.space.order_book:create_index('primary', {
    type = 'hash',
    parts = {1, 'string'},  -- uniqueKey (userid:price:side) entry price
    if_not_exists = true
})
box.space.order_book:create_index('market_side_price_timestamp_index', {
    type = 'TREE',
    parts = {
        3, 'string',  -- market
        4, 'string',   -- side
        2, 'number',   -- entry price
    },
    if_not_exists = true,
    unique=false
})
-- box.space.order_book:create_index("user_price_index", {
--     type = 'hash',
--     parts = {5, 'string', 2, 'number'},  -- user_id, price
--     if_not_exists = true
-- })

-- box.space.order_book:create_index('market_side_index', {
--     type = 'hash',
--     parts = {2, 'string', 3, 'integer'},  -- price, market, side
--     if_not_exists = true
-- })
-- box.space.order_book:create_index('market_index', {
--     type = 'hash',
--     parts = {2, 'string'},  -- price, market, side
--     if_not_exists = true
-- })
-- box.space.order_book:create_index('side_index', {
--     type = 'hash',
--     parts = {3, 'integer'},  -- price, market, side
--     if_not_exists = true
-- })





-- Create function --

-- (Market Price) --
-- Function to get market price --
function get_market_price(key)
    local result = box.space.market_price:select({key})
    if #result > 0 then
        print(result[1][2])
        io.flush()
        return result[1][2]
    else
        return nil
    end
end

box.schema.func.create('get_market_price', {if_not_exists = true})


-- (Users) --
-- Function to get user wallet balance --
function get_user_wallet_balance(key)
    local result = box.space.users:select({key})
    if #result > 0 then
        return result[1][2]
    else
        return nil
    end
end

box.schema.func.create('get_user_wallet_balance', {if_not_exists = true})

-- Function to update user wallet balance --
function update_user_wallet_balance(user_id, new_balance)
    box.space.users:update(user_id, {{'=', 2, new_balance}})
end

box.schema.func.create('update_user_wallet_balance', {if_not_exists = true})

-- Function to create new user --
function create_user_wallet_balance(user_id, wallet_balance)
    print(wallet_balance)
    box.space.users:insert({user_id, wallet_balance})
end

box.schema.func.create('create_user_wallet_balance', {if_not_exists = true})

-- -- Function to delete user --
-- function delete_user(user_id)
--     box.space.users:delete(user_id)
-- end

-- box.schema.func.create('delete_user', {if_not_exists = true})

-- -- Function to get all users --
-- function get_all_users()
--     return box.space.users:select()
-- end

-- box.schema.func.create('get_all_users', {if_not_exists = true})





-- Creating function for orderbooks --

-- Getting orders for long orders --
-- function get_orders_by_price(market, side, price)
--     local results = {}

--     -- Select orders using the composite index
--     -- Using the market, side and sorting by price and created_at
--     local orders = box.space.order_book:index('market_side_price_timestamp_index'):select({market, side})
--     print("orders retrieve")
--     print(orders)
--     io.flush()
--     -- Build the result array
--     for _, order in ipairs(orders) do
--         local price = order[1]  -- price
--         local user_info = {order[4], order[5]}  -- user_id, position_size, entry_price

--         -- Check if there's already a price entry in results
--         if not results[price] then
--             results[price] = {price, {user_info}}  -- Create a new entry if it doesn't exist
--         else
--             table.insert(results[price][2], user_info)  -- Append user info to existing price entry
--         end
--     end

--     -- Convert results from table to array
--     local result_array = {}
--     for _, value in pairs(results) do
--         table.insert(result_array, value)
--     end


--     -- table.sort(result, function(a, b) return a[1] < b[1] end)


--     --     -- Sort the result array in descending order by price
--     --     table.sort(result_array, function(a, b)
--     --         return a[1] > b[1]  -- Sort by price in descending order
--     --     end)



--     return result_array
-- end

-- Getting orders for short orders --





function get_orders_by_market_side_and_price(market, side, price)
    local results = {}

    -- Determine the iterator based on the comparator
    local iterator
    if side == -1 then
        -- Short will retrieved bid orders
        iterator = box.index.LE
    else
        -- Long will retrieved ask orders
        iterator = box.index.GE
    end
 

    -- Select orders using the composite index with the specified price condition
    local orders = box.space.order_book.index.market_side_price_timestamp_index:select(
        {market, side, price},  -- Include the price in the selection criteria
        {iterator = iterator, limit = 100}  -- Use the selected iterator
    )

    -- Build the result array
    for _, order in ipairs(orders) do
        local order_price = order[1]  -- price
        local user_info = {order[4], order[5], order[6]}  -- user_id, position_size, entry_price
        print("user_info")
        if not results[order_price] then
            results[order_price] = {order_price, {user_info}}  -- Create a new entry if it doesn't exist
        else
            table.insert(results[order_price][2], user_info)  -- Append user info to existing price entry
        end

        -- -- Check if the order price meets the threshold
        -- if (comparator == "LE" and order_price <= price) or (comparator == "GE" and order_price >= price) then
        --     -- Check if there's already a price entry in results
        --     if not results[order_price] then
        --         results[order_price] = {order_price, {user_info}}  -- Create a new entry if it doesn't exist
        --     else
        --         table.insert(results[order_price][2], user_info)  -- Append user info to existing price entry
        --     end
        -- end
    end


    -- Convert results from table to array
    local result_array = {}
    for _, value in pairs(results) do
        table.insert(result_array, value)
    end

    print(result_array)

    -- -- Sort the result array in descending order by price
    -- table.sort(result_array, function(a, b)
    --     return a[1] > b[1]  -- Sort by price in descending order
    -- end)

    return result_array
end
box.schema.func.create('get_orders_by_market_side_and_price', {if_not_exists = true})



--- Insert Order Data to order book --
function insert_order_data(primary_key,price, market, side, user_id, position_size)
    local existing_order = get_order_by_primary_key(primary_key)
    local created_at = os.time()
    local position_size = tonumber(position_size)

    if existing_order then
        -- Extract the current position size and increase it
        local current_position_size = existing_order[6]
        position_size = current_position_size + position_size

        -- Delete the current order
        box.space.order_book.index.primary:delete({primary_key})
    end

    box.space.order_book:insert({primary_key,price, market, side, user_id, position_size, created_at})
end
box.schema.func.create('insert_order_data', {if_not_exists = true})


-- Get Order by Price and User ID --
function get_order_by_price_and_user_id(user_id, price)
    local result = box.space.order_book.index.user_price_index:select({user_id, price})
    if #result > 0 then
        return result[1]
    else
        return nil
    end
end
box.schema.func.create('get_order_by_price_and_user_id', {if_not_exists = true})


function get_order_by_primary_key(primaryKey)
    local result = box.space.order_book.index.primary:select({primaryKey})
    if #result > 0 then
        return result[1]
    else
        return nil
    end
end
box.schema.func.create('get_order_by_primary_key', {if_not_exists = true})



 -- box.space.order_book.index.market_side_price_timestamp_index:select({"BTC","1",100},{iterator="GE"}) --
 -- For ask (short) orders, the price should be greater than or equal to the specified price (from orderbook)
 -- the result array will return in order of lowest buy to highest buy with incremental timestamp, we can build ordermap sequentially, and tranverse from behind to match"

 --  box.space.order_book.index.market_side_price_timestamp_index:select({"BTC","-1",300},{iterator="LE"}) --
 -- This is for (bid) long orders, the price should be less than or equal to the specified price (from orderbook)
 -- the result array is arraneging from highest sell to lowest sell with decremental timestamp, we can build ordermap from bottom and sort it from high price to low price"