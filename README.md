Ontrack CLI
===========

[Ontrack](https://github.com/nemerosa/ontrack) is an application which store all events which happen in your CI/CD environment: branches, builds, validations, promotions, labels, commits, etc. It allows your delivery chains to reach new levels by driving your pipelines using real-time data.

The Ontrack CLI is a Command Line Interface tool, available on many platforms, which allows you to feed information into Ontrack from any shell platform.

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

### Build setup

Then, you can create a build entry the same way:

```bash
# Setup of the build
ontrack-cli build setup --project <project> --branch <branch> --build <build>
```

where `<build>` is a unique identifier for your build (typically a build number).

## TODO

- Setup of the project for Git
- Setup of the branch for Git
- Setup of the validation stamps
- Setup of the promotions
- Setup of the auto promotions
- Validations
