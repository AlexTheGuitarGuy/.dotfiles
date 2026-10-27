Never cancel a running job. This covers GitHub Actions runs (`gh run cancel`), Kubernetes jobs and
cronjob runs, deploys, pipelines, and any other CI or cluster job, including ones I started myself
and ones that look stale or superseded.

When a running job is outdated (for example a deploy building an older commit), let it finish and
start the new run after it, or queue it behind it. If waiting genuinely seems wrong, ask the user
first and say why, instead of cancelling.

Why: I cancelled a running ENT deploy because a newer commit had been pushed. The user does not want
jobs cancelled at all.
