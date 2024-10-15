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


-- Create Positions spaces --
box.schema.create_space('positions', {
    if_not_exists = true,
    format ={
        {name = 'user_id', type = 'string'},          -- User ID
        {name = 'market', type = 'string'},            -- Market (e.g., 'BTC', 'ETH')
        {name = 'position_size', type = 'number'},     -- Position size
        {name = 'avg_entry_price', type = 'number'},   -- Average entry price
        {name = 'side', type = 'string'}              -- Side (1 for long, -1 for short)
    }
})
box.space.positions:create_index('user_market_side_index', {
    type = 'hash',
    parts = {1, 'string', 2, 'string', 5, 'string'},  -- user_id, market
    if_not_exists = true
})


-- Create order book spaces --
box.schema.create_space('order_book', {
    if_not_exists = true,
    format =
    {
        {name = 'id', type = 'string'}, -- uniqueKey (userid:price:side:market) entry price
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
box.space.order_book:create_index('market_side_index',{
    parts = {
        3, 'string',  -- market
        4, 'string',   -- entry price
    },
    if_not_exists = true,
    unique=false
})
box.space.order_book:create_index('userid_index',{
    parts = {
        5, 'string',  -- user_id
    },
    if_not_exists = true,
    unique=false
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


-- Create Procedure --
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
-- Procedure to get user wallet balance --
function get_user_wallet_balance(key)
    local result = box.space.users:select({key})
    if #result > 0 then
        return result[1][2]
    else
        return nil
    end
end
box.schema.func.create('get_user_wallet_balance', {if_not_exists = true})

-- Procedure to update user wallet balance --
function update_user_wallet_balance(user_id, new_balance)
    box.space.users:update(user_id, {{'=', 2, new_balance}})
end
box.schema.func.create('update_user_wallet_balance', {if_not_exists = true})

-- Procedure to create new user --
function create_user_wallet_balance(user_id, wallet_balance)
    print(wallet_balance)
    box.space.users:insert({user_id, wallet_balance})
end
box.schema.func.create('create_user_wallet_balance', {if_not_exists = true})



-- (Orderbook) --
--- Insert Order Data to order book --
function insert_order_data(primary_key,price, market, side, user_id, position_size, created_at)
    local existing_order = get_order_by_primary_key(primary_key)
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

-- Get order by primary key --
function get_order_by_primary_key(primaryKey)
    local result = box.space.order_book.index.primary:select({primaryKey})
    if #result > 0 then
        return result[1]
    else
        return nil
    end
end
box.schema.func.create('get_order_by_primary_key', {if_not_exists = true})

