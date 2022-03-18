# codeartifact-dependabot-sync

This project is a little go tool that pushes a refreshed AWS CodeArtifact token to an entire Organizations Dependabot. 

# How to run

1. building the binary/tool

    ```Bash
    go build .
    ```

2. using tool with run flags or env variables:

- CODEARTIFACT_DOMAIN 

    AWS CodeArtifact Domain for which access is required

- CODEARTIFACT_DOMAIN_OWNER 
    
    owner (AWS acc) for the AWS CodeArtifact domain

- CODEARTIFACT_DURATION

    duration of the AWS CodeArtifact authToken

- DEPENDABOT_ORGA
    
    the GitHub organization for which the secret should be created

- DEPENDABOT_OWNER
    
    owner of the GitHub organization

- GITHUB_APP_ID
    
    the ID of the GitHub App used for authentication

- GITHUB_SECRET
    
    GitHub secret for GitHub App authentication

