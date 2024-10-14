---@diagnostic disable: lowercase-global
-- Require the box module --
local box = require('box')

-- Create a space --


-- Create Local function --
local function print_log(value)
    print(value)
    io.flush()
end


-- Create function

function match_orders()
    print_log("Match Orders")
end

box.schema.func.create('match_orders', {if_not_exists = true})
