# codeartifact-dependabot-sync

This project is a little go tool that pushes a refreshed AWS CodeArtifact token to an entire Organizations Dependabot. 

# How to run

1. building the binary/tool

    ```Bash
    go build .
    ```

2. using tool with run flags or env variables:

-CODEARTIFACT_DOMAIN string
    AWS CodeArtifact Domain for which access is required
-CODEARTIFACT_DOMAIN_OWNER string
    owner (AWS acc) for the AWS CodeArtifact domain
-CODEARTIFACT_DURATION string
    duration of the AWS CodeArtifact authToken
-DEPENDABOT_ORGA string
    the GitHub organization for which the secret should be created
-DEPENDABOT_OWNER string
    owner of the GitHub organization
-GITHUB_APP_ID string
    the ID of the GitHub App used for authentication
-GITHUB_SECRET string
    GitHub secret for GitHub App authentication

