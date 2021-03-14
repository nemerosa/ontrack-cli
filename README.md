Ontrack CLI
===========

[![Build](https://github.com/nemerosa/ontrack-cli/actions/workflows/go.yml/badge.svg)](https://github.com/nemerosa/ontrack-cli/actions/workflows/go.yml)

[Ontrack](https://github.com/nemerosa/ontrack) is an application which store all events which happen in your CI/CD environment: branches, builds, validations, promotions, labels, commits, etc. It allows your delivery chains to reach new levels by driving your pipelines using real-time data.

The Ontrack CLI is a Command Line Interface tool, available on many platforms, which allows you to feed information into Ontrack from any shell platform.

## Table of contents

* [Installation](#installation)
* [Setup](#setup)
* [Usage](#usage)
  * [Branch setup](#branch-setup)
  * [Validation stamps setup](#validation-stamps-setup)
  * [Build setup](#build-setup)
  * [Git integration](#git-integration)
* [Integrations](#integrations)
* [TODO](#todo)

## Installation

Download the latest version for your platform from the [releases](https://github.com/nemerosa/ontrack-cli/releases) page.

No further installation step is needed; the CLI is coded in Golang and does not need any dependency.

## Setup

You need to register a configuration:

```bash
ontrack-cli config create prod https://ontrack.example.com --token <token>
```

This registers an installation called `prod`, located at https://ontrack.example.com, using an authentication token.

The configuration is stored on disk, in `~/.ontrack-cli-config.yaml` and the `config create` needs to be done only once.

> The Ontrack CLI supports only version 4.x and beyond of Ontrack.

## Usage

After the configuratio has been set, injection of data into Ontrack from a CI pipeline can be typically done this way.

### Branch setup

We make sure the branch managed by the pipeline is registered into Ontrack:

```bash
# Setup of the branch
ontrack-cli branch setup --project <project> --branch <branch>
```

Here, `<project>` is the name of your project or repository, and `<branch>` is typically the Git branch name
or the PR name (like `PR-123`). The `branch setup` operation is idempotent.

### Validation stamps setup

The CLI can be used to create validation stamps:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation>
```

The `validation-stamp setup` (or `vs setup` for a shortcut) command is idempotent.

Additionally, a validation stamp can be created with a
[data type](https://static.nemerosa.net/ontrack/release/latest/docs/doc/index.html#validation-stamps-data)
and its configuration. For example, to create a CHML validation type:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation> \
    --data-type net.nemerosa.ontrack.extension.general.validation.CHMLValidationDataType \
    --data-config '{warningLevel: {level: "HIGH",value:1},failedLevel:{level:"CRITICAL",value:1}}'
```

The later syntax is pretty cumbersome and the CLI provides dedicated commands for the most used data types:

* for CHML data type:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation> \
    chml \
        --warning HIGH=1 \
        --failed CRITICAL=1
```

* for test summary data type:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation> \
    tests --warning-if-skipped true
```

* for percentage data type:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation> \
    percentage \
        --warning 60 \
        --failure 50 \
        --ok-if-greater false
```

* for metrics data type:

```bash
ontrack-cli validation-stamp setup --project <project> --branch <branch> --validation <validation> \
    metrics
```

### Build setup

Then, you can create a build entry the same way:

```bash
# Setup of the build
ontrack-cli build setup --project <project> --branch <branch> --build <build>
```

where `<build>` is a unique identifier for your build (typically a build number).

### Git integration

Ontrack can leverage SCM information stored in its model, in order to compute change logs or to allow searches based on commits.

For example, to associate a project with a GitHub repository:

```bash
# GitHub setup of the project
ontrack-cli project set-property --project <project> github \
    --configuration github.com \
    --repository nemerosa/ontrack-cli \
    --indexation 30 \
    --issue-service self
```

This command associates the project with the `nemerosa/ontrack-cli` repository, using the credentials defined by the `github.com` GitHub configuration stored in Ontrack. Additionally, Ontrack will index the content of this repository every `30` minutes and the GitHub issues will be used to track issues.

Whenever a branch is created, you associate it with the corresponding Git branch this way:

```bash
# Git setup of the branch
ontrack-cli branch set-property --project <project> --branch <branch> git \
    --git-branch <branch>
```

> Note that pull requests are also supported. In this case, the `--git-branch` must be something like `PR-123`.

Finally, each build can be associated with a Git commit:

```bash
# Git setup of the build
ontrack-cli build set-property --project <project> --branch <branch> --build <build> git-commit \
    --commit <full commit hash>
```

## Integrations

While the Ontrack CLI can be used directly, there are direct integrations in some environments:

* [`ontrack-github-actions-cli-setup` GitHub action](https://github.com/nemerosa/ontrack-github-actions-cli-setup) - _installation of the CLI and simplified GitHub/Git setup_

## TODO

- Setup of the promotions
- Setup of the auto promotions
- Validations
- Validation & build run infos (timings)
