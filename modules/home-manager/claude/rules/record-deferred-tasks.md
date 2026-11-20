Record every deferred task the moment it is deferred, in the same turn, without asking first. A task is
deferred when the user says later or not now, when a check can only happen on another system or after a
reboot ("next time you boot Windows..."), or when a step is handed to the user and they have not
confirmed it was done.

Personal tasks about the machines and devices themselves (system config, apps, hardware, backups) go
into `~/.dotfiles/TODO.md`, under the system they belong to (Windows, Linux, Phone). Tasks of one repo
or project (an app, a thesis, a game, a ticket or PR in a work repo) go into that repo's own list: its
tracked TODO file if it has one, otherwise `.excluded/todos.md` in its main checkout (the `main`
worktree in a bare layout), never in a feature worktree, since removing the worktree deletes its
`.excluded/`. The work task list holds only work tasks that span repos (tokens, tooling, rules), plus
one pointer line per repo with open items. Neither goes into `~/.dotfiles/TODO.md`. Write the entry so
it stands on its own: what to do, where the steps or files are, and what is unknown.

Never offer to add a task ("if you want, I can add it to the list") and wait for an answer: add it,
then mention that it was added. An unclear "done" from the user counts as open until it is clear which
item it refers to; ask which one instead of ticking anything off.

Before reporting the full list of open tasks, follow the pointer lines into the repo lists, and check it against the recent sessions and notes for items
that were mentioned but never recorded, and add those first.

Why: Windows LibreWolf setup, an earphone check on Windows, a Windows clock fix, monitor profiles and
old conversation retention were each mentioned in chat but never written down. One session offered to
add them and ended without an answer, another read "windows stuff done" as the LibreWolf setup, and the
incomplete list was later copied as if it were complete.
