# Secrets Management

Secrets are encrypted at rest with [sops](https://github.com/getsops/sops) and
[age](https://github.com/FiloSottile/age) and committed to git as ciphertext.
They are never exported into the shell environment, so ordinary processes
(install scripts, dlx tools, agents) cannot read them. Each command that needs a
secret decrypts only the keys it uses, for that one run. See `zsh/.zshrc` and
`nushell/.config/nushell/config.nu` for the `npm`/`pnpm`/`npx` wrappers.

## Architecture

```
secrets/secrets.yaml        <-- encrypted file, committed to git
.sops.yaml                  <-- creation rules mapping age pubkeys to files
~/.config/sops/age/keys.txt <-- private key, one per machine, never in git
```

sops encrypts value-by-value: key names stay readable in the committed file, only
values are ciphertext. The private key is the only thing that must never be
committed and must be backed up outside git, losing it makes every secret
encrypted with it permanently unrecoverable.

## First-time setup on a new machine

```bash
mkdir -p ~/.config/sops/age
age-keygen -o ~/.config/sops/age/keys.txt
chmod 600 ~/.config/sops/age/keys.txt
```

Note the printed public key (`age1...`), then add it to `.sops.yaml`:

```yaml
creation_rules:
  - path_regex: secrets/secrets\.yaml$
    key_groups:
      - age:
          - age1existing...
          - age1new...   # <-- add the new machine's pubkey here
```

Re-encrypt so every listed recipient can decrypt:

```bash
sops updatekeys secrets/secrets.yaml
```

Commit both files. Back up the new machine's private key
(`~/.config/sops/age/keys.txt`) somewhere outside git.

## Editing secrets

```bash
sops secrets/secrets.yaml
```

Opens the decrypted file in `$EDITOR`; saving re-encrypts automatically.

## Adding a new secret

Add a key via `sops secrets/secrets.yaml` (or
`sops set secrets/secrets.yaml '["NAME"]' '"value"'`), save. Nothing loads it
automatically. Give it to the command that needs it, one of:

- one-off: `sops exec-env secrets/secrets.yaml '<command>'`
- a tool you run often: add a wrapper next to the `npm` one in `config.nu`
  (`with-env (secrets NAME) { ^tool ...$rest }`) and `.zshrc`
  (`NAME=$(_sec NAME) command tool "$@"`)
- an MCP server: a `headersHelper` like context7's below

## Claude Code MCP servers (bootstrap on a new machine)

Claude Code's user-scoped MCP server config lives in `~/.claude.json`, which is
machine-local state and not tracked in this repo. There's no settings.json key
for it (`mcpServers` is not part of the settings schema, confirmed by testing).
On a new machine, run these once (the sops age key must be in place first):

```bash
claude mcp add-json --scope user context7 '{"type":"http","url":"https://mcp.context7.com/mcp","headersHelper":"sops -d --extract '\''[\"CONTEXT7_API_KEY\"]'\'' /home/alex/.dotfiles/secrets/secrets.yaml | jq -Rc '\''{CONTEXT7_API_KEY: .}'\''"}'

claude mcp add --scope user apiportal-mcp -- pnpm dlx @dvag/apiportal-mcp@1.1.0

claude mcp add --transport http --scope user atlassian https://mcp.atlassian.com/v1/mcp/authv2
```

Verify with `claude mcp list`. context7 uses a `headersHelper` that decrypts
`CONTEXT7_API_KEY` from sops at connect time, so no key is stored in
`~/.claude.json` and rotating it in sops needs no re-registration. A static
`--header "...: ${CONTEXT7_API_KEY}"` does not work at user scope: Claude Code
only expands `${VAR}` in project `.mcp.json`, so the literal string is sent.

## Tools reference

| Command | Purpose |
|---|---|
| `age-keygen -o ~/.config/sops/age/keys.txt` | Generate a new age keypair |
| `sops secrets/secrets.yaml` | Edit encrypted secrets |
| `sops -d secrets/secrets.yaml` | Decrypt to stdout |
| `sops updatekeys secrets/secrets.yaml` | Re-encrypt for all recipients in `.sops.yaml` |
