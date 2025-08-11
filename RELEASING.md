# Releasing addlicense

Releases of `addlicense` are done on an *ad hoc* basis when significant changes
have been made. If you are unable to run addlicense at head (e.g. it's part of a
GitHub workflow) and need a recent feature, please file an issue requesting a
new release.

## Process

The release workflow is configured at
[.github/workflows/release.yaml](./.github/workflows/release.yaml). A new
release will be automatically created whenever any tag is pushed to the
repository. (TODO: consider only matching `v[0-9]*`.) Start a release by adding
a tag, then push that tag to GitHub:

```sh
git tag -a v1.2.3-rc0
git push origin master --tags
```

The release workflow will start a Go build and create a
[GitHub release](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository).
You can check the status of the release workflow and investigate any errors
[on the Actions page](https://github.com/google/addlicense/actions/).

When the workflow finishes, check the
[Releases page](https://github.com/google/addlicense/releases) and find the
release for the new tag. This will have an automatic changelog list. Edit the
release to highlight any major new features or other items of note. If this is a
release candidate (`-rc0` suffix on the tag), check the "Set as a pre-release"
box. Then publish the release! If you are creating a new release after a
pre-release step, make sure the "Previous tag" is set to the prior main release
when generating the release notes, e.g. create a list of changes between
`v1.2.0` and `v1.3.0` rather than between `v1.3.0-rc1` and `v1.3.0`.
