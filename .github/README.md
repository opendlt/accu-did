# No CI Here

**All builds, tests, and releases are local-only.**

This project uses a complete local development and release workflow instead of GitHub Actions. See the [root README](../README.md) for details on:

- Container-first development with `make dev-shell`
- Complete local CI pipeline with `make ci-local`
- Local release workflow with `make release-local`
- Cross-platform builds and multi-arch Docker images
- Security scanning and SBOM generation

**Why Local-Only?**

- **Faster feedback:** No queue times or runner limitations
- **Consistent environment:** Same container across all developers
- **Complete autonomy:** No dependency on external CI services
- **Cost effective:** No CI/CD service costs or complexity
- **Security focused:** Private builds with local artifact control

**Getting Started:**

```bash
# Launch interactive development environment
make dev-shell

# Run complete CI validation
make ci-local

# Create production release
make release-local
```

For detailed instructions, see [README.md](../README.md) and [docs/ops/OPERATIONS.md](../docs/ops/OPERATIONS.md).