Scope: only when preparing your own PR for review (the pr-testing-workflow SHOT steps). During
investigation, debugging, reproducing a bug or reviewing someone else's PR, take no screenshots and do not drive the Bruno
or browser UI: use the CLI (`run-ticket.mjs`, `write.mjs`, kubectl) and report the output in chat.

Proof screenshots (evidence for PRs, tickets, test runs) are always taken of the real tool,
never of text rendered into an image. A rendered PNG proves nothing.

Pick the source by what is being proven:

1. CLI output: a real terminal, via `term-shot` (`~/.local/bin/term-shot`, stowed from
   `~/.dotfiles/local-bin`):

       term-shot <out.png> <wait-seconds-per-command> '<cmd1>' '<cmd2>' ...

   It opens Ghostty with nushell, makes it fullscreen and checks the size before typing
   anything, then types each command at the prompt in turn and takes one screenshot of the
   window at the end. Run it with `timeout`.
   - No titles or echo banners: the typed commands are the titles.
   - Commands must not wrap. Shorten with an `alias k = kubectl --context .. -n ..`, a
     `let` for long values, or a `cd` into the evidence folder first.
   - Output must not wrap or scroll off: nu tables with `select` for the needed columns,
     `to text` for log lines, `slice` and a second screenshot when it does not fit.
   - Prefer live sources (kubectl, az, an exec into a pod) over opening saved logs. Saved
     logs are fine for runs that are over.
   - Never show a secret on screen.
   - Look at every image before using it; redo any that wraps, is cut off, or errors.

2. API calls: the Bruno desktop app (`/opt/Bruno/bruno`), driven with `xdotool`: open the
   collection and the ticket folder, select the environment, send the request, and capture
   the window (`import -window <id>`) with request and response both visible.
   `run-ticket.mjs` stays for running the folder and reporting results; its `--evidence`
   PNGs are not proof.

3. Browser flows: the Playwright MCP, a screenshot plus a video per step.

4. Mails: MailDev in the Playwright browser. Open only the mail sent to the AI-owned
   address and screenshot it with subject, recipient, date and body visible. Never
   screenshot the inbox list or anyone else's mail. The user logs the profile in once;
   never type the MailDev password into the browser.

`shot-text.mjs` and `mail-shot.mjs` are no longer used for proof.
