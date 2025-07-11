-- keep lsp errors contained here
local key = key
local event = event
local messages = messages
local notmuch = notmuch

-- map keys to events
-- events are available through the key table
key.r = event.refresh()

-- keys can also map to functions that return an event
key.q = event.quit
-- multiple bindings can map to the same event
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

-- create and navigate queries
key["/"] = event.query.new()
key.left = event.query.prev()
key.right = event.query.next()

-- mark messages for tagging
key[" "] = event.marks.toggle()
key.i = event.marks.invert()

key.enter = function()
  local message = messages.selected() -- this gives the currently highlighted/selected message
  -- TODO: async 'spawn(command, arg, arg, arg)'
  local command = table.concat {
    "/home/matthias/Projekte/übersicht.mail/einsicht/build/bin/einsicht",
    " ",
    message.filename,
    ">/dev/null",
    " ",
    "2>&1"
  }
  if pcall(os.execute, command) then
    notmuch.tag(message, "-unread")
  end
  return event.refresh { message }
end

-- wrapper function that returns a function that tags selected messages
-- with the given tags and returns a refresh event
local function tag_selected_messages(tags)
  return function()
    local selected = messages.selected()
    local messages_of_interest = { selected }
    -- messages.marked() gives a table of all messages marked with event.marks.*
    for _, message in pairs(messages.marked()) do
      if selected ~= message then
        messages_of_interest[#messages_of_interest + 1] = message
      end
    end
    -- notmuch.tag({msg1, msg2}, "+tag1", "-tag2", "+tag3")
    -- equivalent to notmuch tag +tag1 -tag2 +tag3 id:... id:...
    notmuch.tag(messages_of_interest, table.unpack(tags))

    return event.refresh(messages_of_interest)
  end
end

key.d = tag_selected_messages { "+deleted", "-unread", "-inbox" }
key.a = tag_selected_messages { "+archive", "-inbox" }
key.u = tag_selected_messages { "+unread" }
