# gh-audit-org-keys [![Release](https://github.com/jef/gh-audit-org-keys/actions/workflows/release.yaml/badge.svg)](https://github.com/jef/gh-vanity/actions/workflows/release.yaml)

The point of this project is to help demonstrate that users of GitHub could potentially fall victim to getting their private SSH key cracked. This based on the size and complexity of the key the user generates.

Programs like `ssh2john` from **John the Ripper** can best demonstrate how fast an SSH private key can be solved from a _not so_ complex algorithm with low key lengths (think RSA < 1024 bits).

## Installation

1. Install the `gh` cli - see the [installation](https://github.com/cli/cli#installation)

   _Installation requires a minimum version (2.0.0) of the GitHub CLI that supports extensions._

2. Install this extension:

   ```shell
   gh extension install jef/gh-audit-org-keys
   ```

<details>
<summary><strong>Manual Installation</strong></summary>

Requirements: `cli/cli` and `go`.

1. Clone the repository

   ```shell
   # git
   git clone git@github.com:jef/gh-audit-org-keys.git

   # GitHub CLI
   gh repo clone jef/gh-audit-org-keys
   ```

2. `cd` into it

   ```shell
   cd gh-audit-org-keys
   ```

3. Build it

   ```shell
   make build
   ```

4. Install it locally

   ```shell
   gh extension install .
   ```
</details>

## Usage

To run:

```shell
gh audit-org-keys
```

To upgrade:

```sh
gh extension upgrade audit-org-keys
```

### Examples

- `gh audit-org-keys --organization="actions"`
- `gh audit-org-keys --organization="actions" --show-users="all"`

### Acknowledgments

- [Auditing GitHub users’ SSH key quality](https://blog.benjojo.co.uk/post/auditing-github-users-keys)
- [Openwall - John the Ripper](https://www.openwall.com/john/)
    - [magnumripper/JohnTheRipper](https://github.com/magnumripper/JohnTheRipper)
