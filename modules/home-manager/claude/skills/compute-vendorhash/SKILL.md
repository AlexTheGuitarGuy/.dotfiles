---
name: compute-vendorhash
description: >
  Compute vendorHash, cargoHash, npmDepsHash, and other fixed-output derivation
  hashes for Nix flakes. Use when: (1) packaging a Go module with buildGoModule
  and vendorHash is unknown, (2) packaging a Rust crate with buildRustPackage and
  cargoHash is unknown, (3) packaging a Node.js project with buildNpmPackage and
  npmDepsHash is unknown, (4) any Nix derivation uses a placeholder hash like
  sha256-AAAAAAAA... or lib.fakeSha256 and needs the real hash extracted from a
  build failure, (5) the user says "hash mismatch", "vendorHash", "cargoHash",
  "npmDepsHash", "fakeSha256", or "compute the hash".
---

# Compute VendorHash

## Quick Start

For any derivation with an unknown fixed-output hash, the canonical Nix workflow is:

1. **Set a placeholder hash** in the derivation
2. **Run the build**, it will fail with a hash mismatch
3. **Extract the correct hash** from the error's `got:` line
4. **Replace the placeholder** with the real hash

## Placeholder Hashes

Use these exact placeholders so the build fails predictably:

| Hash Type | Placeholder |
|---|---|
| SRI (sha256) | `sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=` |
| Nix base32 | `lib.fakeSha256` |
| Nix base32 (SHA512) | `lib.fakeSha512` |

Always prefer the SRI placeholder (`sha256-AAAAAAAA...`) for new code, it is
self-describing and produces clearer error messages.

## Extraction Pattern

When the build fails, stderr contains exactly one line like:

```
error: hash mismatch in fixed-output derivation '...-go-modules':
         specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
            got:    sha256-5wMMzSI1WtGW0uu7q6eEHWiU1dXcC+2CMGBDGY8A08k=
```

Copy the **full `got:` value** (including the `sha256-` prefix and `=` suffix)
and paste it into the derivation.

## Derivation-Specific Notes

### buildGoModule (Go)

```nix
pkgs.buildGoModule {
  pname = "mytool";
  version = "0.1.0";
  src = ./scripts/mytool;
  vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
}
```

- The hash covers the `vendor/` directory (Go modules).
- Run: `nix build .#mytool` or `nix build --impure --expr '...'`
- Extract from: `hash mismatch in fixed-output derivation '...-go-modules'`

### buildRustPackage (Rust)

```nix
pkgs.buildRustPackage {
  pname = "mytool";
  version = "0.1.0";
  src = ./scripts/mytool;
  cargoHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
}
```

- The hash covers the `cargo-vendor` directory.
- Same extraction pattern as Go.

### buildNpmPackage (Node.js)

```nix
pkgs.buildNpmPackage {
  pname = "mytool";
  version = "0.1.0";
  src = ./scripts/mytool;
  npmDepsHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
}
```

- The hash covers the `node_modules` tarball.
- Same extraction pattern.

### fetchFromGitHub / fetchurl

```nix
pkgs.fetchFromGitHub {
  owner = "...";
  repo = "...";
  rev = "...";
  hash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
}
```

- The hash covers the downloaded tarball.
- Same extraction pattern.

## Automation Script

Use the bundled script to automate the extract-and-replace workflow:

```bash
nix-shell -p nodejs --run "node ~/.claude/skills/compute-vendorhash/scripts/compute-vendorhash.js <path-to-nix-file> <attribute-name>"
```

The script:
1. Finds the placeholder hash in the Nix file
2. Runs `nix build` on the attribute
3. Extracts the correct hash from stderr
4. Replaces the placeholder in-place

See `references/nix-hashes.md` for a complete catalog of hash types and
derivation-specific quirks.
