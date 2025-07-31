# TODO

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
