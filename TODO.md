# TODO

- Expose Model.bottmLine in lua runtime
  - status.set("Message")
  - message = status.set()
  - decouple Runtime from Model with interfaces

- Make colors configurable in runtime
  - colors.bg = 235
  - colors.error = 196
  - decouple Runtime from Model with interfaces

- Add notifications in UI
  - show as floating windows in the upper right corner
  - expire after a timeout
  - have levels with different styling and timeout
  - exposed to lua runtime as events

- Add logging service

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

- change list component:
  - expose movement to runtime as events
  - render line as {date} {sender}→{recipient} {subject} {tags}
  - find a way to update Marked state that does not require a re-fill of the list

- Add error display in UI
