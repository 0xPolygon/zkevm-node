# ok-to-test workflow

By default, when proposing a PR from a fork, the secrets won't be available in
the actions executed for that PR. This can be problematic if the tests require
secrets to run, in that case the reviewers can't know if the changes proposed
will break the tests or not.

## How it works
In order to solve this issue we integrated the [ok-to-test workflow]. It works
as follows:
* The Github actions executed on PRs are split in 3 groups:
  * actions that don't require secrets.
  * trusted actions, they require secrets and are executed from branches in the
  `zkevm-node` repo, secrets are always available to them
  * from-fork actions, they require secrets and are executed from forks, secrets
  are available after approval from users with write access to the repo.
* When a PR is created from a branch in the `zkevm-node` repo, the actions that
don't require secrets and trusted actions are executed.
* When a PR is created from a fork, the actions that don't require secrets are
executed, and the from-fork actions are executed after a user with write access
comments in the PR `/ok-to-test sha=<commit sha>` with the sha of the commit over
which the actions should run. Before adding this comment the user should have
reviewed the code and verified that it won't suppose a security threat (for
instance, the reviewer should verify that the PR is not adding new actions, or is
not trying to disclose the existing secrets).

## Requirements
Our setup relies on the existence of a repo secret called `PERSONAL_ACCESS_TOKEN`
with the value of a personal access token with repo access scope.

## How to add the ok-to-test functionality to an existing workflow
In order to transform an existing workflow into one that uses the ok-to-test
functionality it should be changed like this:
* Add the `repository_dispatch` entry like here https://github.com/0xPolygonHermez/zkevm-bridge-service/pull/148/files#diff-107e910e9f2ebfb9a741fa10b2aa7100cc1fc4f5f3aca2dfe78b905cbd73c0d2R9-R10
* Duplicate the job, if it is called `build`, copy it to `from-fork-build` and
rename `build` to `trusted-build`.
* In `trusted-build` add this `if` as the first item:
```
if: github.event_name == 'pull_request' && github.event.pull_request.head.repo.full_name == github.repository
```
* In `from-fork-build`:
  * Add this `if` as the first item:
  ```
  if:
      github.event_name == 'repository_dispatch' &&
      github.event.client_payload.slash_command.sha != '' &&
      contains(github.event.client_payload.pull_request.head.sha, github.event.client_payload.slash_command.sha)
  ```
  * If it has a checkout action, replace it with:
  ```
  - name: Fork based /ok-to-test checkout
    uses: actions/checkout@v3
    with:
      ref: 'refs/pull/${{ github.event.client_payload.pull_request.number }}/merge'
  ```
  * Add this code as the last step https://github.com/0xPolygonHermez/zkevm-bridge-service/pull/148/files#diff-107e910e9f2ebfb9a741fa10b2aa7100cc1fc4f5f3aca2dfe78b905cbd73c0d2R60-R88

[ok-to-test workflow]: https://github.com/imjohnbo/ok-to-test
