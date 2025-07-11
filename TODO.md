# TODO

- Add command to clear selection

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
