# TODO

- Add async `spawn` method to execute shell commands in the background
  - report result (return code, output) as event
    - will be used later
  - use `spawn` in key.enter binding:

        key.enter = function()
          local message = messages.selected()
          spawn("einsicht", message.filename)
          notmuch.tag(message, "-unread")
        end

- Tab completion on input
  - tab completes to the most recent input with the current input as prefix

- Expose history service in runtime

- Investigate the use of the repeated function calling in OnKey:

        function decorate(fun)
            some_stuff()
            event = fun()
            return modify(event)
        end

        key.enter = decorate(function() ... end)

- Add error event

      key.e = event.error("test error")

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
