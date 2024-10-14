---@diagnostic disable: lowercase-global
-- Require the box module --
local box = require('box')

-- Create a space --
box.schema.space.create('market_price', { 
    if_not_exists = true,
    format = {
        { name = 'asset', type = 'string' },
        { name = 'price', type = 'number' }
    }
})

-- Create Indexes --
box.space.market_price:create_index('primary', { parts = { 1 }, if_not_exists = true })


-- Create function --
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