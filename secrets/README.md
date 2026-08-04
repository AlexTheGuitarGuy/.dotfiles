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

## Tools reference

| Command | Purpose |
|---|---|
| `age-keygen -o ~/.config/sops/age/keys.txt` | Generate a new age keypair |
| `sops secrets/secrets.yaml` | Edit encrypted secrets |
| `sops -d secrets/secrets.yaml` | Decrypt to stdout |
| `sops updatekeys secrets/secrets.yaml` | Re-encrypt for all recipients in `.sops.yaml` |
