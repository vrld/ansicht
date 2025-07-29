# TODO

- Make colors configurable in runtime

      ansicht.set_theme{
        background = 0,
        muted = 8,
        foreground = 7,
        highlight = 15,
        accent = 3,
        secondary = 4,
        tertiary = 6,
        accent_bright = 11,
        secondary_bright = 12,
        tertiary_bright = 14,
      }

  - decouple Runtime from ui Model using the existing interface pattern

- Add `ansicht.notify{message, severity="info"|"warning"|"error", timeout=number}` to lua runtime
  - overwrite status with given message
  - show original status message after timeout
  - timeout defaults to different values based on severity (default "info")
  - different styles for severities
  - multiple competing notifications are resolved by severity, then recency (newest wins)
  - example:

        ansicht.notify{message="hello", timeout=20}
        ansicht.notify{message="WARNING", severity="warning", timeout=10}
        ansicht.notify{message="world", severity="info", timeout=15}

    - severity defaults to info
    - will show "WARNING" for 10 seconds (severity trumps other messages)
    - will then show "world" for 5 seconds (more recent)
    - will then show "hello" for 5 seconds (longest timeout)

- Refactor history service
  - selection should not be part of the service
    - just hold history of inputs sorted by prompt
    - `Get(prompt, index)`
    - `Count(prompt)`
    - `Add(prompt, input)`
    - `Remove(prompt, index[, len])`
    - `Clear(prompt)`

  - externalize selection as own object, bound to prompt
    - `selection := service.InputHistory().GetSelection(prompt)`
    - `selection.Next()`, ...
    - `selection.Get()`

- Expose history service in runtime

      #ansicht.history[prompt]  # Count(prompt)
      ansicht.history[prompt][i]
      ansicht.history[prompt][i + 1] = "foo"

- Save input history to file

- Tab completion on input
  - tab completes to the most recent input with the current input as prefix

- Responsive list item renderer
  - Should omit recipient if line is too narrow
  - May break over multiple lines

- Refactor use of list item component
  - expose movement to runtime as events
  - find a way to update Marked state that does not require a re-fill of the list
