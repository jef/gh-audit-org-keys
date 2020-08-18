# audit-org-keys [![ci](https://github.com/jef/audit-org-keys/workflows/ci/badge.svg)](https://github.com/jef/audit-org-keys/actions?query=workflow%3Aci+branch%3Amaster) [![codeql](https://github.com/jef/audit-org-keys/workflows/codeql/badge.svg)](https://github.com/jef/audit-org-keys/actions?query=workflow%3Acodeql+branch%3Amaster)

The point of this project is to help demonstrate that users of GitHub could potentially fall victim to getting their private SSH key cracked. This based on the size and complexity of the key the user generates.

Programs like `ssh2john` from **John the Ripper** can best demonstrate how fast an SSH private key can be solved from a _not so_ complex algorithm with low key lengths (think RSA < 1024 bits).

## Getting started

### Releases

| Tag | Description | 
|:---:|---|
| `latest` | Built against tagged releases; stable
| `nightly` | Built against HEAD; generally considered stable, but could have problems |

```
GITHUB_ORGANIZATION=actions
GITHUB_PAT=mysecrettoken

docker run --rm -it \
    --env "GITHUB_ORGANIZATION=$GITHUB_ORGANIZATION" \
    --env "GITHUB_PAT=$GITHUB_PAT" \
    "docker.pkg.github.com/jef/audit-org-keys/audit-org-keys:<tag>"
```

> :point_right: View [Available arguments](#available-arguments) and [Available environment variables](#available-environment-variables) below if you'd like to customize input and output

### Development

#### Requirements

- Go 1.14+ or Docker

#### Running

```sh
GITHUB_ORGANIZATION=actions
GITHUB_PAT=mysecrettoken

# Golang
go build
./audit-org-keys

# show users with multiple keys
./audit-org-keys -show-users=multiple

# Docker
docker build -t audit-org-keys:localhost .

docker run --rm -it \
    --env "GITHUB_ORGANIZATION=$GITHUB_ORGANIZATION" \
    --env "GITHUB_PAT=$GITHUB_PAT" \
    audit-org-keys:localhost

# show users without keys
docker run --rm -it \
    --env "GITHUB_ORGANIZATION=$GITHUB_ORGANIZATION" \
    --env "GITHUB_PAT=$GITHUB_PAT" \
    audit-org-keys:localhost -show-users=without
```

##### Available arguments

- `-show-users=<filter>`: display users with filter (`all`, `with`, `without`, `multiple`)

##### Available environment variables

- `GITHUB_ORGANIZATION`*: The organization under audit
- `GITHUB_PAT`*: GitHub Personal Access Token
    - [Create a PAT](https://github.com/settings/tokens) with `read:org` scope
    - Some organizations have SSO; if yours does, then you also need to enable it
- `LOG_LEVEL`: Sets zap log level

> :point_right: Required denoted by `*`

### Acknowledgments

- [Auditing GitHub users’ SSH key quality](https://blog.benjojo.co.uk/post/auditing-github-users-keys)
- [Openwall - John the Ripper](https://www.openwall.com/john/)
    - [magnumripper/JohnTheRipper](https://github.com/magnumripper/JohnTheRipper)
