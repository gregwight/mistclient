# Contributing to mistclient

We welcome contributions! Please follow these guidelines to contribute.

## Development Process

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix: `git checkout -b my-new-feature`.
3.  Make your changes.
4.  Add or update tests for your changes. All contributions should be tested.
5.  Ensure the test suite passes: `go test -v ./...`.
6.  Format your code: `go fmt ./...`.
7.  Commit your changes following the Commit Message Guidelines below.
8.  Push to the branch: `git push origin my-new-feature`.
9.  Create a new Pull Request.

## Commit Message Guidelines

To maintain a clean and readable commit history that facilitates automatic changelog generation, this project follows the Conventional Commits specification.

Your commit messages should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Common types:**

-   **feat**: A new feature (correlates with a `MINOR` version bump).
-   **fix**: A bug fix (correlates with a `PATCH` version bump).
-   **docs**: Documentation only changes.
-   **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc).
-   **refactor**: A code change that neither fixes a bug nor adds a feature.
-   **perf**: A code change that improves performance.
-   **test**: Adding missing tests or correcting existing tests.
-   **build**: Changes that affect the build system or external dependencies (e.g., `go.mod`).
-   **ci**: Changes to our CI configuration files and scripts.

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.