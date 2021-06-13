# audit-org-keys [![Release](https://img.shields.io/github/workflow/status/jef/audit-org-keys/Release?color=24292e&label=Release&logo=github&logoColor=white&style=flat-square)](https://github.com/jef/audit-org-keys/actions/workflows/release.yaml) [![Nightly Release](https://img.shields.io/github/workflow/status/jef/audit-org-keys/Nightly%20Release?color=24292e&label=Nightly%20Release&logo=github&logoColor=white&style=flat-square)](https://github.com/jef/audit-org-keys/actions/workflows/nightly-release.yaml)

The point of this project is to help demonstrate that users of GitHub could potentially fall victim to getting their private SSH key cracked. This based on the size and complexity of the key the user generates.

Programs like `ssh2john` from **John the Ripper** can best demonstrate how fast an SSH private key can be solved from a _not so_ complex algorithm with low key lengths (think RSA < 1024 bits).

## Installation

`go get -u github.com/jef/audit-org-keys`

Also available under [GitHub Releases](https://github.com/jef/audit-org-keys/releases) as an executable.

## Usage

It is required that you use a GitHub Personal Access Token (PAT). You can generate one [here](https://github.com/settings/tokens/new). The required scopes are `['read:org']`. Set your PAT to environment variable `GITHUB_TOKEN`. If `GITHUB_TOKEN` isn't set, then you may not get the results you expect.

```shell
Usage of audit_org_keys:
  -o, --organization string   [required] GitHub organization provided to inspect
  -s, --show-users all        display users with filter (all, `with`, `without`, `multiple`)
```

### Examples

- `audit-org-keys --organization="actions"`
- `audit-org-keys --organization="actions" --show-users="all"`

## Releases

| Tag | Description | 
|:---:|---|
| `latest` | Built against tagged releases; stable
| `nightly` | Built against HEAD; generally considered stable, but could have problems |

### Acknowledgments

- [Auditing GitHub users’ SSH key quality](https://blog.benjojo.co.uk/post/auditing-github-users-keys)
- [Openwall - John the Ripper](https://www.openwall.com/john/)
    - [magnumripper/JohnTheRipper](https://github.com/magnumripper/JohnTheRipper)
