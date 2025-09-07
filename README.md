Ontrack CLI
===========

[![Build](https://github.com/nemerosa/ontrack-cli/actions/workflows/go.yml/badge.svg)](https://github.com/nemerosa/ontrack-cli/actions/workflows/go.yml)

[Ontrack](https://github.com/nemerosa/ontrack) is an application which store all events which happen in your CI/CD environment: branches, builds, validations, promotions, labels, commits, etc. It allows your delivery chains to reach new levels by driving your pipelines using real-time data.

The Ontrack CLI is a Command Line Interface tool, available on many platforms, which allows you to feed information into Ontrack from any shell platform.

> The Ontrack CLI works only with the version 4 of [Ontrack](https://github.com/nemerosa/ontrack).

# Installation

Download the latest version for your platform from the [releases](https://github.com/nemerosa/ontrack-cli/releases) page.

No further installation step is needed; the CLI is coded in Golang and does not need any dependency.

# Setup

You need to register a configuration:

```bash
ontrack-cli config create prod https://ontrack.example.com --token <token>
```

This registers an installation called `prod`, located at https://ontrack.example.com, using an authentication token.

The configuration is stored on disk, in `~/.ontrack-cli-config.yaml` and the `config create` needs to be done only once.

> The Ontrack CLI supports only version 4.x and beyond of Ontrack.

# Usage

After the configuration has been set, injection of data into Ontrack from a CI pipeline can be typically done this way.

## Branch setup

We make sure the branch managed by the pipeline is registered into Ontrack:

```bash
# Setup of the branch
ontrack-cli branch setup --project <project> --branch <branch>
```

Here, `<project>` is the name of your project or repository, and `<branch>` is typically the Git branch name
or the PR name (like `PR-123`). The `branch setup` operation is idempotent.

> Run `ontrack-cli branch setup --help` for additional options.

## Validation stamps setup

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

## Promotions and auto promotion

Promotions can be created using:

```bash
ontrack-cli promotion-level setup --project <project> --branch <branch> --promotion <promotion>
```

Their auto promotion can be set using:

```bash
ontrack-cli promotion-level setup --project <project> --branch <branch> --promotion <promotion> \
   --validation <stamp1> \
   --validation <stamp2> \
   --depends-on <other-promotion-1> \
   --depends-on <other-promotion-2>
```

The validation stamps and promotions this command depends on will be created if they don't exist already.

## Build setup

Then, you can create a build entry the same way:

```bash
# Setup of the build
ontrack-cli build setup --project <project> --branch <branch> --build <build>
```

where `<build>` is a unique identifier for your build (typically a build number).

If you need to associated a release label to your build, you can use the `--release` option:

```bash
ontrack-cli build setup --project <project> --branch <branch> --build <build> --release <label>
```

The same way, you can associate a Git commit property to the build with the `--commit` option:

```bash
ontrack-cli build setup --project <project> --branch <branch> --build <build> --commit <commit>
```

## Git integration

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

> Note that the Git commit property on a build can be set directly using the `build setup` command and the `--commit` option:

```bash
ontrack-cli build setup --project <project> --branch <branch> --build <build> --commit <commit>
```

# Validation

One of the most important point of Ontrack is to record _validations_:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> --status <status>
```

where `<status>` is an Ontrack validation run status like `PASSED`, `WARNING` or `FAILED`.

## Data validation

Additionally, a validation run can be created with some
[data](https://static.nemerosa.net/ontrack/release/latest/docs/doc/index.html#validation-stamps-data). For example, to create a test summary validation:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    --data-type net.nemerosa.ontrack.extension.general.validation.TestSummaryValidationDataType \
    --data {passed: 1, skipped: 2, failed: 3}
```

The later syntax is pretty cumbersome and the CLI provides dedicated commands for the most used data types:

* for CHML data type:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    chml \
        --critical 0 \
        --high 2 \
        --medium 25 \
        --low 1214
```

* for test summary data type:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    tests \
        --passed 20 \
        --skipped 2 \
        --failed 1
```

* for percentage data type:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    percentage \
        --value 87
```

* for metrics data type:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    metrics \
        --metric speed=1.5 \
        --metric acceleration=0.25 \
        --metrics weight=145,height=185.1
```

## Run info

The `validate` commands accept additional flags to set the run info on a validation (source & trigger, duration):

* `--run-time` - duration of the validation in seconds
* `--source-type` - type of source for the validation, for example, the CI name like `jenkins`
* `--source-uri` - the URI to the source of the validation, for example, the URL to a Jenkins job
* `--trigger-type` - how the validation was triggered (for example: `scm`)
* `--trigger-data` - data associated with the trigger (for example, a Git commit)

For example, to set the validation duration on a test summary validation:

```bash
ontrack-cli validate --project <project> --branch <branch> --build <build> --validation <validation> \
    --run-time 80 \
    tests \
        --passed 20 \
        --skipped 2 \
        --failed 1
```

# Auto-versioning

The Ontrack CLI can be used to set up the auto-versioning configuration for a branch.

Given the `auto-versioning.yaml` file containing the configuration, the call looks like:

```bash
ontrack-cli branch --project <project> --branch <branch> auto-versioning \
    --yaml auto-versioning.yaml
```

The `auto-versioning.yaml` file looks like:

```yaml
dependencies:
  - sourceProject: my-library
    sourceBranch: release-1.3
    sourcePromotion: IRON
    targetPath: gradle.properties
    targetProperty: my-version
    postProcessing: jenkins
    postProcessingConfig:
      dockerImage  : openjdk:8
      dockerCommand: ./gradlew clean
```

> The format of this file is fully described in the Ontrack documentation at
> https://static.nemerosa.net/ontrack/release/latest/docs/doc/index.html#auto-versioning-config

In a parent repository, you can use the auto-versioning check to automatically create the dependency links.

```yaml
ontrack-cli build auto-versioning-check \
  --project <project> \
  --branch <branch> \
  --build <build>
```

# Misc

## Direct GraphQL calls

The Ontrack CLI uses the GraphQL API of Ontrack for its communication. The `graphql` command allows to run raw GraphQL queries.

For example:

```bash
ontrack-cli graphql \
    --query 'query ProjectList($name: String!) { projects(name: $name) { id name branches { name } } }' \
    --var name=ontrack-cli
```

## General options

The `--graphqh-log` flag is available for all commands, to enable some tracing on the console for the GraphQL requests and responses.

# Integrations

While the Ontrack CLI can be used directly, there are direct integrations in some environments.

## Jenkins

The [`ontrack-jenkins-cli-pipeline`](https://github.com/nemerosa/ontrack-jenkins-cli-pipeline/) Jenkins pipeline library allows an easy integration between your `Jenkinsfile` pipelines and Ontrack.

## GitHub actions

* [`ontrack-github-actions-cli-setup`](https://github.com/nemerosa/ontrack-github-actions-cli-setup) - _installation of the CLI and simplified GitHub/Git setup_
* [`ontrack-github-actions-cli-validation`](https://github.com/nemerosa/ontrack-github-actions-cli-validation) - _creation of validation runs based on GitHub workflow information_

# Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
