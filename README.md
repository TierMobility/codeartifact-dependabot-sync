# codeartifact-dependabot-sync

Many people are using private registries like [AWS CodeArtifact](https://aws.amazon.com/codeartifact/) to store critical code and distribute it within a controlled domain. [Dependabot](https://github.com/dependabot) is a GitHub integrated tool that allows for dependency analysis with automatic pull requests and alerts for repositories. As the name suggests, private registries are only allowed to be accessed by people and systems that have access.

Until recently, Dependabot's only option to access private registries was to add secrets through the UI. Now they offer additional API operations for [programmatically adding secrets to Dependabot](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/managing-encrypted-secrets-for-dependabot). 

This project aims to become a tool for people who are using AWS CodeArtifact and want to use Dependabot with it. The codeartifact-dependabot-sync enables you to automatically update your secret every 10 hours.

# Getting started

The following instructions show how to setup the environment to run this code within a given environment.

## Prerequisites 

-   A fresh install of Golang 1.17. Please follow [these instructions from the official documentation](https://go.dev/dl/)
    ```console
    foo@bar:~$ go version
    go version go1.17.6 darwin/amd64
    ```

- A GitHub App that has access to Dependabot Secrets for your Repo or Organization. [Official docs](https://docs.github.com/en/developers/apps/getting-started-with-apps/about-apps)

## Installation

1. clone the repo

    ```Bash
    git clone https://github.com/TierMobility/codeartifact-dependabot-sync

    cd codeartifact-dependabot-sync
    ```

1. (optional) get all modules

    ```bash
    export GO111MODULE=on
    go get .
    ```

1. Build it
    ```Bash
    GO111MODULE=on go build . -o /codeartifact-dependabot-sync
    ```

## How to use

> the tool runs indefinitely until the process is killed. This can be dissabled by setting the `DAEMON` option to `false`.

- <a name="setup"></a>Setup the following data:

    | Key  | Description  |
    |---|---|
    | CODEARTIFACT_DOMAIN_OWNER  | Owner (AWS acc) for the AWS CodeArtifact domain. Also used when [using CodeArtifact with AWS Cli](https://docs.aws.amazon.com/cli/latest/reference/codeartifact/login.html)  |
    | CODEARTIFACT_DURATION  | Duration of the AWS CodeArtifact authToken.  |
    | CODEARTIFACT_DOMAIN  | AWS CodeArtifact Domain for which access is required. Also used when [using CodeArtifact with AWS Cli](https://docs.aws.amazon.com/cli/latest/reference/codeartifact/login.html)  |
    | DEPENDABOT_ORG  | The GitHub organization for which the secret should be created  |
    | GITHUB_PRIVATE_KEY  | GitHub secret for GitHub App authentication  |
    | GITHUB_APP_ID  | The ID of the GitHub App used for authentication  |
    | GITHUB_APP_TOKEN  | GitHub App token used for encrypting secrets |

- Using env variables
    1. Setup environment variables regarding [point 1 from installation](#setup)

    2. 
        ```bash
        ./codeartifact-dependabot-sync
        ```

- Using flags

    1. The flags for the tool are the same as demonstrated in [point 1 from installation](#setup). 

        ```Bash
        # Get all the flags and their descriptions:
        ./codeartifact-dependabot-sync -h

        # run it with flag data
        ./codeartifact-dependabot-sync -DEPENDABOT-ORG=exampleOrg  ...
        ```


