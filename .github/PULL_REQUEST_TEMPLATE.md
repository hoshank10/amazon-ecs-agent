<!--
Please make sure you've read and understood our contributing guidelines;
https://github.com/aws/amazon-ecs-agent/blob/master/CONTRIBUTING.md

Please provide the following information:
-->

### Summary
<!-- What does this pull request do? -->

### Implementation details
<!-- How are the changes implemented? -->

### Testing
<!-- How was this tested? -->
<!--
Note for external contributors:
`make test` and `make run-integ-tests` can run in a Linux development
environment like your laptop.  `go test -timeout=30s ./agent/...` and
`.\scripts\run-integ.tests.ps1` can run in a Windows development environment
like your laptop.  Please ensure unit and integration tests pass (on at least
one platform) before opening the pull request.
Once you open the pull request, there will be 14 automatic test checks on the bottom
of the pull request, please make sure they all pass before you merge it. You can
use `bot/test` label to rerun the automatic tests multiple times.
-->

New tests cover the changes: <!-- yes|no -->

### Description for the changelog
<!--
Write a short (one line) summary that describes the changes in this
pull request for inclusion in the changelog.
You can see our changelog entry style here:
https://github.com/aws/amazon-ecs-agent/commit/c9aefebc2b3007f09468f651f6308136bd7b384f
-->

### Additional Information

**Does this PR include breaking model changes? If so, Have you added transformation functions?**
<!-- If yes, next release should have a upgraded minor version -->  

**Does this PR include the addition of new environment variables in the README?**
<!-- 
If it is a sensitive variable, add it to this blocklist in ecs-logs-collector here: https://github.com/aws/amazon-ecs-logs-collector/blob/b0958c2aa424c6dc578d5a8def4422c51791a076/ecs-logs-collector.sh#L63
If it is not a sensitive variable, add it to the allowlist in ecs-logs-collector here: https://github.com/aws/amazon-ecs-logs-collector/blob/b0958c2aa424c6dc578d5a8def4422c51791a076/ecs-logs-collector.sh#L66
-->

### Licensing

By submitting this pull request, I confirm that my contribution is made under the terms of the Apache 2.0 license.
