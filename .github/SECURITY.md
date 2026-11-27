# Security policy

## Supported versions

The latest stable release is supported. Security fixes are developed against
`main`; beta and development builds are provided for evaluation and are not
supported production releases.

## Report a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/webappsgo/cvedex/security/advisories/new)
or choose **Security → Report a vulnerability** in the repository. Do not file
vulnerabilities as public issues. Mirrors should direct reporters to the
canonical GitHub repository.

Include the affected release or commit, deployment environment, impact, and
reproduction details. Do not include real credentials or private user data.
The maintainers will acknowledge the report, investigate it privately, and
coordinate a fix and disclosure timeline with the reporter. Please allow time
for a fix to reach users before public disclosure.

The application publishes `/.well-known/security.txt`, operator security
information at `/server/security`, a disclosure policy at
`/server/security/policy`, and general contact at `/server/contact`. These
instance pages do not replace private GitHub reporting for repository
vulnerabilities.
