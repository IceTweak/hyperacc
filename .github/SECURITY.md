# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < Latest | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

- **Email**: [Add your security email here]
- **GitHub Security Advisory**: Use the "Report a vulnerability" button on the [Security tab](https://github.com/IceTweak/hyperacc/security) of this repository

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

## Disclosure Policy

When the security team receives a security bug report, they will assign it to a primary handler. This person will coordinate the fix and release process, involving the following steps:

1. Confirm the problem and determine the affected versions
2. Audit code to find any potential similar problems
3. Prepare fixes for all releases still under maintenance
4. Publish a security advisory and release the fixes

We aim to respond to security reports within 48 hours and provide regular updates on the progress of the fix.

## Security Best Practices

When using this library, please follow these security best practices:

- Always validate and sanitize input data
- Use the latest version of the library
- Review access control rules regularly
- Follow HyperLedger Fabric security guidelines
- Keep your Go dependencies up to date

## Acknowledgments

We appreciate the security community's help in keeping this project and its users safe. Security researchers who responsibly disclose vulnerabilities will be acknowledged in our security advisories (unless they prefer to remain anonymous).
