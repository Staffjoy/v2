# Schedulerr Modernization Roadmap

This document defines a practical path to evolve this Staffjoy v2 fork into **Schedulerr**: a production-ready, multi-tenant workforce scheduling SaaS.

## Product scope

Schedulerr targets **workforce scheduling** across multiple verticals (restaurants, retail, warehouses, healthcare, and similar shift-based teams). It is intentionally not appointment booking software.

## Guiding principles

- Preserve working business logic where possible; replace only where modernization meaningfully reduces long-term complexity.
- Deliver in incremental phases that keep the system deployable at each stage.
- Build for SaaS fundamentals from day one: tenancy boundaries, observability, billing, and secure defaults.
- Keep architecture AI-ready by modeling schedule constraints in a way that can be consumed by optimization engines later.

## Target end state

- Multi-tenant SaaS platform with organization/location/worker hierarchy.
- PostgreSQL-backed data model with explicit tenant isolation strategy.
- API-first backend with stable auth, permissions, and auditability.
- Stripe-based subscription billing with webhook-driven entitlement updates.
- Production delivery pipeline (CI/CD), containerized services, and runtime observability.

## Phase 0: Architecture and code audit (immediate)

1. Inventory services, dependencies, and ownership boundaries in the monorepo.
2. Classify components into:
   - Keep and modernize
   - Keep temporarily for compatibility
   - Remove/replace
3. Capture external contracts that must remain stable during migration.
4. Produce a short risk register (legacy build tooling, proto/gRPC coupling, auth flows, migration hazards).

**Deliverable:** architecture decision record (ADR) baseline and migration backlog.

## Phase 1: Build and runtime simplification

1. Replace Bazel-centric local development workflow with straightforward service-level builds.
2. Standardize service execution with Docker and docker-compose for local development.
3. Add a unified `.env.example` strategy and runtime config loader conventions.
4. Introduce modern linting/formatting and reproducible CI checks.

**Goal:** make the repo easy to run, test, and change without deep legacy tooling knowledge.

## Phase 2: SaaS domain foundation

1. Add explicit SaaS domain entities:
   - Tenants (organizations)
   - Locations
   - Teams / roles
   - Memberships and permissions
2. Implement RBAC with scoped permissions (org-level and location-level).
3. Add audit logging for sensitive actions (publishing schedules, permission changes, billing-impacting actions).
4. Formalize API versioning and backward-compatible schema evolution rules.

**Goal:** establish durable multi-tenant boundaries before adding monetization and advanced scheduling.

## Phase 3: Billing and entitlements (Stripe)

1. Define subscription plans and feature matrix (free/trial/pro tiers).
2. Integrate Stripe Checkout + Customer Portal.
3. Add webhook handlers for:
   - Subscription created/updated/canceled
   - Invoice payment success/failure
   - Trial status transitions
4. Implement entitlement resolution in backend middleware to enforce plan limits.
5. Add billing admin screens and invoice history endpoints.

**Goal:** monetize cleanly with reliable, webhook-driven state.

## Phase 4: Production hardening

1. Add structured logging, request IDs, and correlation across services.
2. Add metrics + health probes + SLO-aligned alerts.
3. Add error monitoring integration and incident-friendly dashboards.
4. Tighten security posture:
   - secrets management
   - least privilege service accounts
   - secure cookie/session defaults
   - dependency and container scanning
5. Validate backup/restore and migration rollback processes.

**Goal:** production readiness with clear operational controls.

## Phase 5: Scheduling intelligence readiness

1. Model scheduling constraints explicitly (availability, labor targets, compliance constraints, fairness).
2. Introduce a pluggable scheduling engine boundary (internal solver now, Timefold-compatible adapter later).
3. Capture optimization telemetry (solver inputs, output score, accepted/rejected suggestions).
4. Add simulation mode for “what-if” schedule generation without impacting published schedules.

**Goal:** future AI/optimization integration without architectural rework.

## Definition of done for “production-ready” milestone

- Tenant isolation tested and enforced.
- RBAC and audit logging active for critical workflows.
- Stripe subscription lifecycle fully automated with tested webhook processing.
- CI pipeline covers linting, tests, migrations, and image builds.
- Deployment playbook, rollback playbook, and on-call runbook published.
- Core scheduling flows validated under realistic load.

## Suggested execution cadence

- Work in 1–2 week milestones.
- End each milestone with:
  - changelog entry
  - migration notes
  - explicit “safe to deploy” status
- Prefer vertical slices (API + DB + UI + tests) over broad partial rewrites.


## Implemented in this repository (current baseline)

- Added `.env.example` as a standardized local configuration template.
- Added `docker-compose.yml` with PostgreSQL + MailHog to simplify local dependency bootstrapping.
- Added `scripts/check-env.sh` and `Makefile` targets for repeatable environment validation and compose workflows.
- Added CI workflow (`.github/workflows/phase1-foundation.yml`) to validate env shape and compose configuration.
- Added ADR `docs/adr/0001-phase1-foundation.md` to document Phase 1 decisions.

## Immediate next actions

1. Freeze a baseline architecture diagram and service inventory from the current repo.
2. Decide the near-term platform direction:
   - continue in Go with focused modernization, or
   - staged migration to a different stack where it materially simplifies operations.
3. Open tracking issues per phase and tag them by risk/impact.
4. Implement Phase 1 first so every following phase moves faster.
