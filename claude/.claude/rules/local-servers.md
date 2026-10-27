Stop a local server by its port, never by name. Find the pid with `ss -ltnp | grep :<port>` and
kill that pid. Never use `pkill -f` or a grep over `ps` output: the pattern also matches the
agent's own shell and misses servers started under another command line.

Before testing against a local server, check that the process on the port is the one just
started: its start time (`ps -o lstart -p <pid>`) and, when config matters, the env key names
it runs with (names only, never values). After stopping, check that the port is free.

Why: an old API instance kept holding the port after a "stop", and a dry run ran against it
instead of the fresh one; earlier a `pkill -f` pattern matched the agent's own shell.
