Test data on the DFS environments (households, persons and anything hanging off them) may be
created and changed by the AI without asking, within these limits.

Environments: only `d01`-`d05`, `entwicklung` and `integration`. Never write to `abnahme` or
`produktion`, whatever the request says; the scripts refuse them too.

Tooling lives in `~/projects/api-collections/dvag-api-collections/scripts/ai/`. Every run goes
through sops so credentials come from `~/.dotfiles/secrets/secrets.yaml` without ever being
printed:

    sops exec-env ~/.dotfiles/secrets/secrets.yaml 'node ~/projects/api-collections/dvag-api-collections/scripts/ai/<script> ...'

Always this exact form (it is what the auto mode allow rule matches).

- `create-household.mjs --env <e> [--persons n] [--ticket T]`: creates a household with a main
  person (and n extra persons), last names end in `AI<timestamp>`, ids are appended to the
  sops registry `~/.dotfiles/secrets/ai-test-data.yaml`.
- `write.mjs --env <e> --method <M> --service <slug> --path <p> --owner-id <id> [--body f]`:
  any other write. Refuses unless `--owner-id` is in the registry for that env and appears in
  `--path`.
- `run-ticket.mjs --env <e> --ticket <T> [--household <id>]`: runs the ticket's Bruno folder
  against the newest AI household (or the given owned one), JSON report in the tmp dir.

Ownership: modify only data listed in the registry. Data the AI did not create (a teammate's
household, the shared `advisor_sophie` data, the default `householdId`/`personId` in the Bruno
environments) is read-only unless the user names that exact id in the current conversation;
only then pass `--allow-foreign <id>`. Every write goes through `write.mjs` or
`create-household.mjs`, never a raw curl or an ad hoc Bruno request with a write method.

Advisors cannot be created through DFS APIs. Use the advisor Sophie from sops:
`DVAG_ADVISOR_SOPHIE_USERNAME` on every env, `DVAG_ADVISOR_SOPHIE_ENTWICKLUNG_PASSWORD` on
entwicklung and `DVAG_ADVISOR_SOPHIE_PASSWORD` on integration and d01-d05.

Where services run: asset-service only on entwicklung, every other service only on
integration and d01-d05; the scripts refuse other combinations. Households can only be
created on integration (and d0x). On entwicklung use the advisor Sophie household, the
`householdId`/`personId` from `environments/entwicklung.yml`, registered as granted by the user. If a needed key is missing, stop and ask the user to add it; never borrow another
login.

Secrets: never run `sops -d` on `secrets.yaml`, never echo or log a token or password, never
put one in a report, commit, Bruno file or PR comment. Refer to accounts by env var name only.

Side effects to keep in mind: creating a person also pushes wish goals to Salesforce, and
there is no DELETE for persons or households, so created data stays. Reuse the newest AI
household for the env before creating another one.
