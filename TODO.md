# TODO

- Refactor search as general input, move behavior to query

      keys["/"] = cmd.input(function(query)
        return cmd.query.new(query)
      end)

  - `cmd.input` enters interactive mode (ref `m.focusSearch`)
  - takes a function as argument
  - once finished, the function will be called with the input text and the selected messages
  - result will give the next command
  - `cmd.query` signals to search for the given query (new event)

- Implement tag function

      keys["t"] = cmd.input(function(tags)
        notmuch.tag(...)
        return cmd.refresh
      end)

- Add async spawn method to execute shell commands in the background
  - use spawn in key.enter binding

- Tweak UI Layout:
  > {query} {%d marked}/{%d total}
  > {list of messages}
  > {query tabs}
  > {input line}

- change list component:
  - expose movement to runtime as events
  - render line as {date} {sender}→{recipient} {subject} {tags}
  - find a way to update Marked state that does not require a re-fill of the list

- Add error display in UI
