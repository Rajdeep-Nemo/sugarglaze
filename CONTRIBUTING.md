# Contributing

We appreciate your interest in contributing. Please read through these guidelines to help keep the project consistent and the review process smooth.

---

## Branching Strategy

| Branch | Purpose |
|--------|---------|
| `main` | Contains production-ready code that passes all tests. Do not branch from here. |
| `dev`  | Active development happens here. Branch from `dev` and open a PR back into it. |

---

## Workflow

```
Clone → Branch from dev → Commit → Pull request
```

Keep changes small and focused. One logical change per branch makes review faster and easier.

---

## Commit & PR Prefixes

All commits and PR titles must begin with one of the following prefixes:

| Prefix | Use for |
|--------|---------|
| `bug:` | Defect fixes |
| `feature:` | New functionality |
| `docs:` | Documentation updates |
| `refactor:` | Code restructuring |

Use short, descriptive messages. Clearly describe what changed and why in the PR body.

---

## Code Style

- Follow standard Go naming conventions.
- Match the style of surrounding code.
- Prefer simple, readable code over clever solutions.
- Use only the Go standard library — this project is intentionally self-contained and does not use external dependencies.

---

## Tests & Docs

- Add or update tests for every change.
- Update relevant documentation alongside code changes.

---

## Before Submitting

- Verify the project builds cleanly.
- Avoid introducing new dependencies.
- Clearly describe what changed and why in the PR.
- Format the code by running `go fmt ./...` from main directory.

---

## Contributors

Once your PR is merged, add your name to [CONTRIBUTORS.md](CONTRIBUTORS.md).
