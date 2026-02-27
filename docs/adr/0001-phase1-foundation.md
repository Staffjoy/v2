# ADR 0001: Phase 1 modernization foundation

## Status
Accepted

## Context
The inherited Staffjoy v2 repository is difficult to run and evolve quickly due to historical tooling and fragmented setup conventions.

## Decision
As the first executable modernization step, introduce a minimal, tool-agnostic local foundation:

1. Standard `.env.example` for local runtime defaults.
2. `docker-compose.yml` for Postgres + MailHog development dependencies.
3. `scripts/check-env.sh` for required environment validation.
4. `Makefile` shortcuts for repeatable setup and compose actions.
5. GitHub Action to verify environment and compose configuration on PRs.

## Consequences

### Positive
- Faster onboarding for contributors.
- Reduced setup drift between local environments.
- Simple CI gate that catches config regressions early.

### Tradeoffs
- This does not yet migrate services off legacy build/runtime architecture.
- Compose currently provides dependencies, not a full service mesh replacement.

## Follow-up
- Add service containers incrementally as each service is modernized.
- Migrate CI checks to include service-level tests per milestone.
