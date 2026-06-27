# CI/CD Rules (PART 28)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use Makefile in CI/CD — CI runs steps directly
- ❌ Pin third-party actions to a tag — full commit SHA required
- ❌ Use pull_request_target for untrusted code
- ❌ Use gitleaks for secret scanning — use truffleHog
- ❌ Create build-toolchain.yml for Go projects — casjaysdev/go:latest needs no custom toolchain
- ❌ Create ci.yml or release.yml before code is complete and tests pass
- ❌ Skip workflow policy gates
- ❌ Reference secrets in untrusted code paths

## CRITICAL - ALWAYS DO
- ✅ Implement CI/CD on all 5 providers
- ✅ Pin all third-party actions to full commit SHA (never a tag)
- ✅ Run secret scan with truffleHog
- ✅ Publish checksums, SBOM, and provenance with every release
- ✅ Run lint+test+secret-scan+workflow-policy+vuln-scan in parallel before build
- ✅ Follow workflow creation order (security → ci → release)

## REQUIRED CI/CD PROVIDERS (all 5)
| Provider | File Location |
|----------|--------------|
| GitHub Actions | .github/workflows/ |
| GitLab CI | .gitlab-ci.yml |
| Gitea Actions | .gitea/workflows/ |
| Forgejo Actions | .forgejo/workflows/ |
| Jenkins | Jenkinsfile |

Same quality gates must be applied on all providers.

## GITHUB ACTIONS JOB ORDER
```
[parallel gate]
  lint + test + secret-scan + workflow-policy + vuln-scan
        ↓
     build (all 8 platforms)
        ↓
  [parallel post-build]
  coverage + image-scan + upload-artifacts
```

## WORKFLOW CREATION ORDER (non-negotiable sequence)
1. **Security workflows** — secret scan, SHA/digest policy, dependency audit
   (no build dependency; safe to add first)
2. **ci.yml** — add LAST, only after code is complete and tests pass
3. **release.yml** — add LAST, only after ci.yml is proven green

Go projects NEVER get build-toolchain.yml.

## RELEASE ARTIFACTS (required for every release)
| Artifact | Description |
|----------|-------------|
| cvedex-{os}-{arch}[.exe] | All 8 platform binaries |
| checksums.txt | SHA-256 of all binaries |
| SBOM | Software Bill of Materials |
| provenance | SLSA build provenance |

## SECRET SCANNING
- Tool: truffleHog (NEVER gitleaks)
- Run on every PR and push to main
- Block merge on any finding

## THIRD-PARTY ACTION PINNING
```yaml
# WRONG — tag can be moved
- uses: actions/checkout@v4

# CORRECT — immutable SHA
- uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
```

---
For complete details, see AI.md PART 28
