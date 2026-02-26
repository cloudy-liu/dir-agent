# Release Process

## Version reset policy

The project release line was reset to `v0.1.0`.

- Old tags/releases (`v0.5`, `v0.6`) must be removed before publishing `v0.1.0`.
- New public releases start from `v0.1.0` and continue with semantic version tags.

## Trigger

Release is fully automated by GitHub Actions.

- Trigger: push a version tag matching `v*`
- Workflow: `.github/workflows/release.yml`
- Release author: `github-actions[bot]`

## One-command flow

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions will run tests, build binaries, package release bundles, and publish release assets automatically.

## Release assets format

Each platform gets a single zip file:

- `diragent_<tag>_windows_amd64.zip`
- `diragent_<tag>_windows_arm64.zip`
- `diragent_<tag>_darwin_amd64.zip`
- `diragent_<tag>_darwin_arm64.zip`
- `diragent_<tag>_linux_amd64.zip`
- `diragent_<tag>_linux_arm64.zip`

Each zip contains:

- `diragent` or `diragent.exe`
- one-click entrypoints: `install` and `uninstall`
- `scripts/install.*`
- `scripts/uninstall.*`
- `README.quickstart.md`

## Cleanup old tags/releases

Use GitHub CLI:

```bash
gh release delete v0.5 --yes
gh release delete v0.6 --yes
git tag -d v0.5 v0.6
git push origin :refs/tags/v0.5 :refs/tags/v0.6
```
