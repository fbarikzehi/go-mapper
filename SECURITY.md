# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please follow these steps:

1. **Do not** create a public GitHub issue
2. Email the maintainers at security@example.com
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and provide a timeline for fixes.

## Security Best Practices

When using gomap:

1. **Validate input data** before mapping
2. **Use depth limits** to prevent stack overflow
3. **Enable circular reference detection** for untrusted data
4. **Sanitize data** from external sources
5. **Review custom converters** for security implications

## Disclosure Policy

- Security vulnerabilities will be disclosed publicly after a fix is released
- We follow a 90-day disclosure timeline
- Critical vulnerabilities may be disclosed sooner if actively exploited

Thank you for helping keep gomap secure!

---
