## PR Review Instructions

When performing a Pull Request review, do your typical PR analysis, and:

### 1. Risk Assessment
Provide a **Risk Rating** at the top of the review summary:
- **HIGH:** Changes to core logic, fundamental architectural patterns, or critical shared utilities.
- **MEDIUM:** Significant feature additions or modifications to established business logic.
- **LOW:** Documentation, styling, minor bug fixes in non-critical paths, or boilerplate.

### 2. Targeted Review Areas
Identify specific code blocks that require **scrupulous human review**. Focus on:
- Complex conditional logic or concurrency-prone areas.
- Potential breaking changes in internal or external APIs.
- Logic that lacks sufficient unit test coverage within the PR.

### 3. Reviewer Recommendations
Analyze the `CODEOWNERS` file and the git history (recent editors) to suggest the most qualified reviewers.
- Prioritize individuals who have made significant recent contributions to the specific files modified.
- Cross-reference these contributors with the defined `CODEOWNERS` for the directory.

### 4. Code Style
Give style advice based on the following guides, in order of priority.
1. [Effective Go](https://go.dev/doc/effective_go)
2. [Google Code Review Comments](https://go.dev/wiki/CodeReviewComments)
3. [Google Style Guide](https://google.github.io/styleguide/go/)

Style exceptions are acceptable when aligning with pre-existing "local" style from the same file or package, but they should still be noted.
