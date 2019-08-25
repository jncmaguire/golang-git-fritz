# golang-git-fritz
<!-- badges -->

Fritz is a series of helpers (mainly hooks), written in Golang. The helpers are meant to make it easier to follow:

- [Conventional Commits v1.0.0-beta.4](./adr/0003-follow-conventional-commits.md)

## Features Overview

* **Contained.** Cross-compiled, easy to install, easy to remove.
* **Configurable.** Takes a TOML to personalize your settings, like specifying the list of scopes.
* **As much as you want it to be.** Install or remove as few hooks as you'd like for a specific repo. Hook options:
  - prepare-commit-msg
    * `commit-prep`: Build a commit template based on conventional commits standards.
  - commit-msg
    * `commit-validate`: Validate (and fail) a commit based on conventional commits standards

## Get Started

### Requirements (Execution)

* git
* bash

### Install

#### User Profile

#### Project

1. Check into the root of your project, on the same level as the .git folder.
2. Run `fritz setup <space-delimited list of hooks>` to install the hooks you want.

### Remove

#### Project

1. Check into the root of your project, on the same level as the .git folder.
2. Run `fritz cleanup` to remove the hooks fritz created and your fritz configuration file.

#### User Profile

## Support

* As a one person experimental project, there isn't any. Use at your own peril.

## Misc.

* [docs](/docs)

<!-- References -->