# Pull Request guidelines

This document addresses how we should create PRs, give and receive reviews. The motivation is to have better code, reduce the time from creation to merge while sharing knowledge and insights that help everyone becoming better developers.

Note that none of this is a hard rule, but suggestions / guidelines. Although everyone is encouraged to stick to this points as much as possible. Use your common sense if some of this do not apply well on a particular PR

## How to create a good PR

- Follow the template, unless for some reason it doesn't fit the content of the PR
- Try hard on doing small PRs (> ~400 lines), in general is better to have 2 small PRs rather than a big one
- Indicate clearly who should review it, ideally 2 team mates
- Author of the PR is responsible for merging. Never do it until you have the approval of the specified reviewers unless you have their explicit permission
- Introduce the purpose of the PR, for example: `Fixes the handle of ...`
- Give brief context on why this is being done and link it to any relevant issue
- Feel free to ask to specific team mates to review specific parts of the PR

## How to do a good review

- In general it's hard to set a quality threshold for changes. A good measure for when to approve is to accept changes once the overall quality of the code has been improved (compared to the code base before the PR)
- Try hard to avoid taking things personally. For instance avoid using `I`, `you`, `I (don't) like`, ...
- Ask, don’t tell. ("What about trying...?" rather than "Don’t do...")
- Try to use positive language. You can even use emoji to clarify tone.
- Be super clear on how confident you are when requesting changes. One way to do it is by starting the message like this:
  - `Opinion: ...` this way you're indicating a some how personal preference to the PR author. It's great to share opinions that may not be based on evidence, The author should understand this as `if you agree with me, you could do this change`
  - `Suggestion: ...` similar to opinion, but should have some back up / evidence / reasoning on why the suggestion is better than what done by the author. It should be read as `I think this could be done better this way, unless you have arguments to defend the original solution you should do it`
  - `Request: ...` indicates that the PR won't be approved unless the request is applied. This should always include arguments that make obvious why the changes being requested are needed
- Avoid doing code reviews for too consecutive time, try to don't do reviews for more than one hour non-stop

## How to receive feedback

- Accept that many programming decisions are opinions. Discuss tradeoffs, which you prefer, and reach a resolution quickly.
- Ask for clarification if needed. ("I don’t understand, can you clarify?")
- Offer clarification, explain the decisions you made to reach a solution in question.
- Try to respond to every comment.
- If there is growing confusion or debate, ask yourself if the written word is still the best form of communication. Talk (virtually) face-to-face, then mutually consider posting a follow-up to summarize any offline discussion (useful for others who be following along, now or later).
- If consensus is still not reached, involve someone else in the discussion. As a last resource the lead of the project could take the decision

## Links and credits

This guide is based on the following content:

- https://smartbear.com/learn/code-review/best-practices-for-peer-code-review/
- https://github.blog/2015-01-21-how-to-write-the-perfect-pull-request/
