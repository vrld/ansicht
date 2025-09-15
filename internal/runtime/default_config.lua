local ansicht = ...

-- global (!) table that maps key presses to functions
key = {
  R = ansicht.refresh,  -- ansicht.refresh() reloads the query

  -- you can bind any function
  q = function()
    ansicht.history.save()  -- save input history between sessions
    ansicht.quit()
  end,

  ["/"] = function()
    -- switches to input mode
    ansicht.input {
      prompt = "notmuch search ",
      placeholder = "tag:unread",
      -- called with the input text when the input is committed
      with_input = ansicht.query.new,
    }
  end,
  left = ansicht.query.prev,
  right = ansicht.query.next,

  -- mark messages to fill ansicht.messages.marked()
  [" "] = ansicht.marks.toggle,
  i = ansicht.marks.invert,
  x = ansicht.marks.clear,

  enter = function()
    -- get the highlighted message with fields:
    --   id,
    --   thread_id,
    --   date,
    --   filename,
    --   from,
    --   to,
    --   subject.
    local message = ansicht.messages.selected()

    -- run external commands
    ansicht.exec {
      "xdg-open", message.filename,
      next = function(err)
        if err == nil then
          ansicht.tag(message, "-unread")
          ansicht.refresh{ message }
        end
      end,
    }
  end,

  -- get reply templates from notmuch
  r = function()
    local template = ansicht.reply(ansicht.messages.selected())
    ansicht.exec { "vim", stdin = template }
  end,

  g = function()
    local template = ansicht.reply_group(ansicht.messages.selected())
    ansicht.exec { "vim", stdin = template }
  end,
}

-- define aliases like so
key["ctrl+c"] = key.q
key["ctrl+d"] = key.q

-- use full lua scripting

-- wrapper function that returns a function that tags selected messages
-- with the given tags and refreshes the messages
local function tag_selected_messages(tags)
  local selected = ansicht.messages.selected()
  local messages_of_interest = { selected }

  -- messages.marked() gives a table of all messages marked with `ansicht.marks.*` (see above)
  for _, message in pairs(ansicht.messages.marked()) do
    if selected ~= message then
      messages_of_interest[#messages_of_interest + 1] = message
    end
  end

  -- ansicht.tag({msg1, msg2}, "+tag1", "-tag2", "+tag3")
  -- equivalent to notmuch tag +tag1 -tag2 +tag3 id:... id:...
  ansicht.tag(messages_of_interest, table.unpack(tags))

  ansicht.status.set("Tagged " .. #messages_of_interest .. " messages: " .. table.concat(tags, " "))
  ansicht.refresh(messages_of_interest)
end

key.t = function()
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

key.d = function() tag_selected_messages { "+deleted", "-unread", "-inbox" } end
key.a = function() tag_selected_messages { "+archive", "-inbox" } end
key.u = function() tag_selected_messages { "+unread" } end

-- run code that depends on the UI here
function Startup()
  ansicht.status.set("ansicht")
  ansicht.history.load()  -- load input history from previous session
end
