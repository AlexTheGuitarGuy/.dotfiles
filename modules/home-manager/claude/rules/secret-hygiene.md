Never print secret values. Do not run unfiltered `env`, `printenv`, `set`, `kubectl get secret -o yaml`
or `-o json`, full `docker inspect`, or `az keyvault secret show`, and never enable shell tracing
(`set -x`). Never put a secret in a command argument. Redact logs and response bodies before quoting
them. Read key or variable names only, never values.

Why: logs and transcripts persist, and one stray env dump leaks every credential in the pod.
