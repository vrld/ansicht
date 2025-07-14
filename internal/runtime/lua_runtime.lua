-- keep lsp errors contained here
local event = event

-- make event.input more ergonomical
-- luaPushInput registers the callback in the registry and when
-- the event is consumed, the registry is cleaned. this means
-- that we need a fresh event on every key press, otherwise the
-- callback will only be executed on the first invocation.
-- OnKey will execute functions until it reaches the userdata
event.input_raw = event.input
event.input = function(config)
  return function() return event.input_raw(config) end
end
