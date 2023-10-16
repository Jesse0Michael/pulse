# Pulse - AI Empowered Insights

Use AI to get insights from different services. 
Pulse pull activity from a specified service and use AI to summarize that activity. 
Add parameters to narrow the scope of the data.

## Data Sources

### Github
Pulse uses the [Github API](https://docs.github.com/en/rest) to pull activity from the [ListEventsPerformedByUser](https://docs.github.com/en/rest/activity/events#list-events-for-the-authenticated-user) endpoint.
#### Events:
- PushEvent
- PullRequestEvent

These events are parsed and formatted into the content of our summary generation prompt.

## Generative AI

### OpenAI
Pulse uses the [OpenAI API](https://platform.openai.com/docs/api-reference) to generate summaries from the [CreateChatCompletions](https://platform.openai.com/docs/api-reference/chat/create) endpoint using the  data pulled from the data sources.


## Pulse CLI

### Installation
```go
go install github.com/jesse0michael/pulse/cmd/pulse@latest
```

### Usage
```
generate a summary of a github user's activity

Usage:
   github [username] [flags]

Flags:
  -h, --help   help for github

Environment:
  GITHUB_URL         the url for accessing the GitHub API
  GITHUB_TOKEN       the authentication token to use with the GitHub API
  OPENAI_URL         the url for accessing the OpenAI API
  OPENAI_TOKEN       the authentication token to use with the OpenAI API
```

### Example
``` bash
pulse github jesse0michael
```
>The user "jesse0michael" has been active on multiple repositories on GitHub. Here is a summary of their recent activity:
>
>1. Repository: Jesse0Michael/pulse
>   - Pushed commits with the following commit messages:
>     - Fix: Use correct directory to build CLI
>     - Chore: Remove debug line
>   - Merged branch 'main' into the repository
>
>2. Repository: Jesse0Michael/pulse
>   - Pushed commits with the following commit messages:
>     - Chore: Add Environment configuration to CLI usage
>     - Chore: Rename pulse cmd directory
>   - Merged branch 'main' into the repository
>
>3. Repository: Jesse0Michael/pulse
>   - Pushed commits with the following commit messages:
>     - Feat: Update API route to use GitHub and username
>     - CI: Use standard GITHUB_TOKEN
>
>4. Repository: Jesse0Michael/pulse
>   - Pushed commits with the following commit messages:
>     - Fix: Refactor code into Pulser service struct to use in both the CLI and API
>     - Fix: openai service summary output formatting
>     - Fix: Clear default GitHub service URL
>     - Test: Test openai service
>
>5. Repository: Jesse0Michael/pulse
>   - Pushed commits with the following commit messages:
>     - Fix: Make GitHub service testable
>
>6. Repository: Jesse0Michael/go-rest-assured
>   - Pushed commits with the following commit messages:
>     - CI: Release from main
>     - Chore: Change default branch to main
>     - Build: Move docker labels
>     - Feat: BREAKING CHANGE upgrade rest assured to v4
>     - Feat: Export Serve method to start HTTP listener
>     - Remove client context and error channel
>     - Rely on the caller of the package to appropriately call Serve
>     - Chore: Upgrade to Go 1.21
>     - Feat: Move to log/slog for logging
>     - Fix: Use google/uuid package
>     - Test: Serve rest assured client in tests
>     - Chore: Update license
>     - Feat: Add NewClientServe function
>     - Build: Add docker labels
>     - CI: Update Go version
>   - Opened a pull request in the repository with the title "feat: BREAKING CHANGE upgrade rest assured to v4" and a description containing links to related commits.
>
>7. Repository: Jesse0Michael/fetcher
>   - Pushed commits with the following commit message:
>     - Fix: Disable Instagram feed
