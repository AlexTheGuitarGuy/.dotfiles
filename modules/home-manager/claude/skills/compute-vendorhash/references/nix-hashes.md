# Nix Hash Types Reference

## SRI Hashes (Preferred)

SRI (Subresource Integrity) hashes are self-describing and include the algorithm
prefix. Nix supports:

- `sha256-...`: Most common
- `sha512-...`: For larger inputs or when sha256 collides

Format: `<algorithm>-<base64>`

Example: `sha256-5wMMzSI1WtGW0uu7q6eEHWiU1dXcC+2CMGBDGY8A08k=`

## Nix Base32 Hashes (Legacy)

Older Nix code uses base32-encoded hashes without an algorithm prefix. These are
identified by length:

- SHA-256: 52 characters
- SHA-512: 103 characters

Helper functions:
- `lib.fakeSha256`: Placeholder for SHA-256
- `lib.fakeSha512`: Placeholder for SHA-512

## Derivation-Specific Hash Attributes

| Builder | Hash Attribute | What It Covers |
|---|---|---|
| `buildGoModule` | `vendorHash` | `vendor/` directory (Go modules) |
| `buildRustPackage` | `cargoHash` | `cargo-vendor/` directory |
| `buildNpmPackage` | `npmDepsHash` | `node_modules` tarball |
| `fetchFromGitHub` | `hash` | Downloaded source tarball |
| `fetchurl` / `fetchgit` | `sha256` or `hash` | Downloaded file/archive |
| `rustPlatform.buildRustPackage` | `cargoHash` | Same as `buildRustPackage` |
| `buildGradlePackage` | `gradleDepsHash` | Gradle dependencies |
| `buildPythonPackage` with `format = "pyproject"` | `dependenciesHash` | Poetry/pdm lock |

## Common Error Patterns

### buildGoModule

```
error: hash mismatch in fixed-output derivation '...-go-modules':
         specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
            got:    sha256-5wMMzSI1WtGW0uu7q6eEHWiU1dXcC+2CMGBDGY8A08k=
```

### buildRustPackage

```
error: hash mismatch in fixed-output derivation '...-cargo-deps':
         specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
            got:    sha256-...
```

### fetchFromGitHub

```
error: hash mismatch in fixed-output derivation '...-source':
         specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
            got:    sha256-...
```

## Tools

### nix-prefetch-url

```bash
nix-prefetch-url https://example.com/file.tar.gz
```

Outputs the SRI hash directly.

### nix-prefetch-git

```bash
nix-prefetch-git --url https://github.com/owner/repo --rev v1.0.0
```

Outputs JSON with `hash` field.

### nix-prefetch-github

```bash
nix-prefetch-github owner repo --rev v1.0.0
```

Same as above, shorthand for GitHub.
