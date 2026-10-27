Never take over the keyboard, mouse or screen without asking the user right before, every time.
This covers anything that opens or focuses windows, types or clicks on the desktop: `term-shot`,
`k9s-shot`, `xdotool`, `wmctrl`, driving the Bruno app, fullscreen windows. Say what will open,
what gets typed and roughly how long it takes, then wait for an explicit go. An earlier approval
does not carry over to the next run, a retake or a related step.

The Playwright browser MCP is not covered: it runs in its own browser and does not touch the
user's input.

Why: screenshot tooling took over the desktop twice while the user was working.
