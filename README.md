# audit-org-keys

The point of this project is to help demonstrate that users of GitHub could potentially fall victim to getting their private SSH key cracked. This based on the size and complexity of the key the user generates.

Programs like `ssh2john` from **John the Ripper** can best demonstrate how fast an SSH private key can be solved from a _not so_ complex algorithm with low key lengths (think RSA < 1024 bits).

## Getting started

### Requirements

- Go 1.14+ or Docker
- GitHub Personal Access Token
- GitHub Organization that you can read
    - Example: [actions](https://github.com/actions)

### Running

#### Golang
```sh
export GITHUB_ORGANIZATION=actions
export GITHUB_PAT=mysecrettoken

# native
make run

# Docker
make run-docker
```

### Acknowledgments

- [Auditing GitHub users’ SSH key quality](https://blog.benjojo.co.uk/post/auditing-github-users-keys)
- [Openwall - John the Ripper](https://www.openwall.com/john/)
    - [magnumripper/JohnTheRipper](https://github.com/magnumripper/JohnTheRipper)
