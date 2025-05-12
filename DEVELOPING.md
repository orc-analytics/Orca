## Developer Guide

### Release process

Releases are managed through git tags. To perform a release, following this flow:

1. Checkout a branch with the release pattern, `release/v*.*.*`. This states the version of the intended release
2. Perform all neccessary dev work and testing on this release branch. There should not be major work required as
   the code should be close to release ready. Ensure that all pipelines pass.
3. Perform a pre-release of this branch, by creating a tag locally and pushing it to the origin: `git tag v0.1.0+test`.
   Test the artifacts
4. Once happy, and all documentation is up to date, merge the release branch into main.
5. Create a release tag locally, from main, and push to the origin: `git tag v0.1.0 && push origin v0.1.0`
