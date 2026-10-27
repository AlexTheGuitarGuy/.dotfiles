Any write or modification to a system visible to other people (Jira tickets, Confluence pages, GitHub/PR comments, Slack messages, emails, or similar) requires explicit approval from the user before posting, every time, no standing permission.

Before posting, the content must:
1. Read as authored by the user, not by an AI. No AI-tell phrasing, no signature or disclosure that it was AI-written.
2. Follow the no-em-dash rule (see no-em-dash.md), since this is user-facing/external-facing text.
3. Contain no commit hashes, run or comment ids, or other precise magic values (exact limits,
   line numbers, byte counts) that a person would not bother to type. Say "fixed" and describe
   the change in plain words instead; that kind of precision reads as AI-written.

Read-only lookups (viewing tickets, pages, threads) need no approval. Any create/update/comment/transition action always needs a review-and-approve step first, drafted in plain prose without dashes or AI-sounding phrasing.

Why: this kind of content is visible to teammates and represents the user's own voice/authorship, unlike local file edits or memory notes.
