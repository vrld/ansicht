-- keep lsp errors contained here
local key = key
local ansicht = ansicht

-- map keys to functions
key.r = ansicht.refresh

-- multiple bindings can map to the same event
key.q = ansicht.quit
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

-- create and navigate queries
key["/"] = function()
  ansicht.input{
    placeholder = "tag:unread",
    prompt = "notmuch search ",
    with_input = function(query)
      ansicht.query.new(query)
      ansicht.status.set("")
    end,
  }
end
key.left = ansicht.query.prev
key.right = ansicht.query.next

-- mark messages for tagging
key[" "] = ansicht.marks.toggle
key.i = ansicht.marks.invert
key.x = ansicht.marks.clear

key.enter = function()
  local message = ansicht.messages.selected() -- this gives the currently highlighted/selected message
  ansicht.spawn{
    "/home/matthias/Projekte/übersicht.mail/einsicht/result/bin/einsicht",
    message.filename,
    next=function()
      ansicht.tag(message, "-unread")
      ansicht.refresh { message }
    end
  }
end

-- wrapper function that returns a function that tags selected messages
-- with the given tags and returns a refresh event
local function tag_selected_messages(tags)
  local selected = ansicht.messages.selected()
  local messages_of_interest = { selected }
  -- messages.marked() gives a table of all messages marked with event.marks.*
  for _, message in pairs(ansicht.messages.marked()) do
    if selected ~= message then
      messages_of_interest[#messages_of_interest + 1] = message
    end
  end
  -- notmuch.tag({msg1, msg2}, "+tag1", "-tag2", "+tag3")
  -- equivalent to notmuch tag +tag1 -tag2 +tag3 id:... id:...
  ansicht.tag(messages_of_interest, table.unpack(tags))

  ansicht.status.set("Tagged " .. #messages_of_interest .. " messages: " .. table.concat(tags, " "))
  ansicht.refresh(messages_of_interest)
end

key.d = function() tag_selected_messages { "+deleted", "-unread", "-inbox" } end
key.a = function() tag_selected_messages { "+archive", "-inbox" } end
key.u = function() tag_selected_messages { "+unread" } end

key.t = function ()
  ansicht.input {
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
      tag_selected_messages(tags)
    end,
  }
end

-- test status functionality
function Startup()
  return ansicht.status.set("ansicht")
end
