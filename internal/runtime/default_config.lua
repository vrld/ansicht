-- keep lsp errors contained here
local key = key
local event = event
local messages = messages
local notmuch = notmuch
local spawn = spawn
local status = status

-- map keys to events
-- events are available through the key table
key.r = event.refresh()

-- keys can also map to functions that return an event
key.q = event.quit
-- multiple bindings can map to the same event
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

-- create and navigate queries
key["/"] = event.input {
  placeholder = "tag:unread",
  prompt = "notmuch search ",
  with_input = event.query.new
}
key.left = event.query.prev()
key.right = event.query.next()

-- mark messages for tagging
key[" "] = event.marks.toggle()
key.i = event.marks.invert()

key.enter = function()
  local message = messages.selected() -- this gives the currently highlighted/selected message
  spawn{
    "/home/matthias/Projekte/übersicht.mail/einsicht/build/bin/einsicht",
    message.filename,
    next=function()
      notmuch.tag(message, "-unread")
      return event.refresh { message }
    end
  }
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

key.x = event.marks.clear
key.t = event.input {
  placeholder = "-unread +act",
  prompt = "notmuch tag ",
  with_input = function(tags_str)
    local idx_space = tags_str:find(" ")
    local tags = { tags_str:sub(1, idx_space) }
    while idx_space ~= nil do
      local next_word_idx = idx_space + 1
      idx_space = tags_str:find(" ", next_word_idx)
      tags[#tags + 1] = tags_str:sub(next_word_idx, idx_space)
    end
    return tag_selected_messages(tags)
  end,
}

-- test status functionality
key.s = function()
  return event.status("Hello from Lua! Current status: " .. status.get())
end

function Startup()
  return event.status("ansicht")
end
