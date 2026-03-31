# Security Policy

## Scope

tagaudit is a static analysis tool that reads Go source code and reports struct tag issues. It does not:

- Execute analyzed code
- Make network requests
- Handle credentials or secrets
- Process untrusted input beyond Go source files loaded via `golang.org/x/tools/go/packages`

Security concerns for this project are limited to:

- Bugs in the analysis pipeline that could cause unexpected file writes (via `--fix`)
- Dependency vulnerabilities in third-party modules
- Path traversal or injection in the CLI or config file handling

## Supported Versions

Only the latest release is supported with security fixes.

## Reporting a Vulnerability

If you find a security issue, please report it privately rather than opening a public issue.

**Email:** Open a [private security advisory](https://github.com/emm5317/tagaudit/security/advisories/new) on GitHub.

You should expect an initial response within 7 days. If the issue is confirmed, a fix will be released as soon as practical.

## Dependencies

tagaudit depends on:

- `golang.org/x/tools` — Go analysis framework
- `github.com/fatih/structtag` — Struct tag parsing
- `github.com/spf13/cobra` — CLI framework
- `gopkg.in/yaml.v3` — YAML config parsing

Dependency vulnerabilities are monitored via Dependabot. If you notice a vulnerable dependency, please open an issue.
