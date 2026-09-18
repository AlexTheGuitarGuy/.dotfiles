# Secrets Management

Secrets are encrypted at rest with [sops](https://github.com/getsops/sops) and
[age](https://github.com/FiloSottile/age), committed to git as ciphertext, and
decrypted into the shell environment at zsh/nushell startup. See
`zsh/.zshrc` and `nushell/.config/nushell/config.nu` for the decrypt-on-shell-start
wiring.

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

Add a key via `sops secrets/secrets.yaml`, save. It's picked up automatically at
next shell start since both `zsh/.zshrc` and `nushell/config.nu` decrypt and export
every key in the file.

## Claude Code MCP servers (bootstrap on a new machine)

Claude Code's user-scoped MCP server config lives in `~/.claude.json`, which is
machine-local state and not tracked in this repo. There's no settings.json key
for it (`mcpServers` is not part of the settings schema, confirmed by testing).
On a new machine, run these once, after the shell has picked up
`CONTEXT7_API_KEY` from sops (open a new shell first so the variable is set):

```bash
claude mcp add --transport http --scope user context7 https://mcp.context7.com/mcp \
  --header "CONTEXT7_API_KEY: $CONTEXT7_API_KEY"

claude mcp add --scope user apiportal-mcp -- pnpm dlx @dvag/apiportal-mcp@1.1.0

claude mcp add --scope user lens -- /opt/Lens/resources/cli/bin/lens-cli-linux-x64 mcp-server

claude mcp add --transport http --scope user atlassian https://mcp.atlassian.com/v1/mcp/authv2
```

Verify with `claude mcp list`. Note: the `CONTEXT7_API_KEY` currently in sops
is rejected by Context7 as invalid ("API keys should start with 'ctx7sk'
prefix" even though it does) - get a fresh key from context7.com and
`sops secrets/secrets.yaml` to replace it before relying on that server.

## Tools reference

| Command | Purpose |
|---|---|
| `age-keygen -o ~/.config/sops/age/keys.txt` | Generate a new age keypair |
| `sops secrets/secrets.yaml` | Edit encrypted secrets |
| `sops -d secrets/secrets.yaml` | Decrypt to stdout |
| `sops updatekeys secrets/secrets.yaml` | Re-encrypt for all recipients in `.sops.yaml` |
