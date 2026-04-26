# Frontend-Backend Integration with Google Session, Pinia, and Real Booking Flows

## TL;DR
> **Summary**: Replace the current mock-driven Nuxt flows with real backend-backed catalog, auth, and booking behavior using a same-origin Google cookie session, Pinia stores, and new frontend-facing read/write contracts. Fix backend schema/API gaps first, then migrate `/login`, `/register`, `/lounge`, `/event`, and `/gaming` onto one shared session/catalog/booking backbone.
> **Deliverables**:
> - Cookie-backed Google auth session with `session/logout/dev-login` contract for frontend and QA
> - Backend public catalog/read-model APIs for home, lounge, event, and gaming
> - Backend booking hardening: session-owned booking creation, overlap/capacity checks, gaming seat/place support
> - Pinia-backed frontend API layer and shared stores for session, catalog, and booking drafts
> - Real data integration for `/login`, `/register`, `/lounge`, `/event`, `/gaming`, and auth-aware header behavior
> - `docs/integration.md` with route/API/store/env/verification matrix and implementation write-up
> **Effort**: XL
> **Parallel**: YES - 4 waves
> **Critical Path**: 1 → 2 → 3 → 6 → 7 → 8 → 9/10/11 → 12 → 13 → 14

## Context
### Original Request
Create a detailed plan to connect everything currently on the frontend to the backend so it works end-to-end, identify what is missing on the backend, add the backend work, store state in Pinia, and document in detail what was done.

### Interview Summary
- All current frontend surfaces in scope: `/`, `/gaming`, `/lounge`, `/event`, `/login`, `/register`.
- Auth model is locked to **Google-only**.
- Public browsing remains allowed; auth is required for booking/profile actions.
- `gaming` stays **in scope in the same overall plan**, not deferred out of scope.
- Documentation deliverable is a new canonical file: `docs/integration.md`.
- Testing strategy is **tests-after**.

### Metis Review (gaps addressed)
- Locked automated auth QA to a deterministic **non-production dev-login/session-seed** path; plan does not rely on manual Google login.
- Locked session contract semantics so acceptance criteria are binary and agent-executable.
- Locked `gaming` to an explicit MVP: browse availability, choose zone/place/config/time window, create booking, show real availability conflicts; no tournaments/admin tooling/payment redesign.
- Locked post-login scope to auth-aware browsing/booking/logout and draft restoration only; **no my-bookings/profile-management UI expansion** in this pass.
- Locked `event` to a **frontend-facing read model over services/zones + metadata**, not a full admin-managed event CMS.

## Work Objectives
### Core Objective
Turn the current mock frontend into a real integrated client over the existing Go backend by adding the missing session, catalog, booking, and gaming capabilities; migrate app state to Pinia; and document the resulting architecture and verification workflow.

### Deliverables
- Backend auth/session contract for browser use
- Backend public read-model endpoints for home/lounge/event/gaming
- Backend booking model hardening and gaming place/config support
- Frontend API client and Pinia bootstrap
- Pinia stores: `session`, `catalog`, `booking`
- Real frontend integration for `/login`, `/register`, `/lounge`, `/event`, `/gaming`, and auth-aware header/shell behavior
- New `docs/integration.md` and targeted cross-links in existing docs

### Definition of Done (verifiable conditions with commands)
- [ ] `cd backend && go test ./...` passes
- [ ] `cd frontend && corepack pnpm run lint` passes
- [ ] `cd frontend && corepack pnpm run typecheck` passes
- [ ] `cd frontend && corepack pnpm run smoke` passes
- [ ] `curl` checks prove public catalog endpoints are guest-accessible, protected booking writes reject unauthenticated callers, and authenticated writes derive booking ownership from session rather than request body
- [ ] Lounge, event, and gaming all load real backend data and create real persisted bookings on success paths
- [ ] `docs/integration.md` exists and documents route→store→endpoint mapping, auth/session contract, env vars, verification commands, and scope decisions

### Must Have
- Google-only cookie-backed session for frontend auth flows
- Session endpoints usable by Nuxt + Pinia (`session`, `logout`, deterministic QA login)
- Public catalog APIs for `/`, `/lounge`, `/event`, `/gaming`
- Server-side booking invariants: overlap, capacity, session-owned actor, gaming place availability
- Pinia-backed frontend state for session/catalog/booking
- Real backend integration on `/login`, `/register`, `/lounge`, `/event`, `/gaming`
- New `docs/integration.md`

### Must NOT Have (guardrails, AI slop patterns, scope boundaries)
- No email/password auth, guest booking, or provider expansion beyond Google
- No admin/seller UI rollout in this pass
- No full event CMS/scheduler backend beyond catalog metadata needed for the existing client surface
- No payment redesign, notification system, analytics, or unrelated UI refresh
- No “manual Google login required” acceptance criteria
- No Pinia sprawl for purely local presentational state that does not need sharing, persistence, or server sync

## Verification Strategy
> ZERO HUMAN INTERVENTION - all verification is agent-executed.
- Test decision: **tests-after**
- Backend verification: `go test ./...` + `curl` contract checks using cookie jar
- Frontend verification: existing Playwright smoke expanded for real auth-aware flows + `lint` + `typecheck`
- Auth QA rule: add a **non-production-only dev session seed endpoint** that produces the same cookie/session shape as the real Google callback so agents can verify protected flows deterministically
- Evidence: `.sisyphus/evidence/task-{N}-{slug}.{ext}`

## Execution Strategy
### Parallel Execution Waves
> Target: 5-8 tasks per wave. Shared auth/catalog/booking foundations are extracted early for maximum parallelism.

Wave 1: backend foundation and contract repair
- Task 1: Browser session contract and auth middleware
- Task 2: Zone/schema contract repair and metadata foundation
- Task 3: Booking contract hardening and invariants
- Task 4: Public lounge/event catalog read models
- Task 5: Gaming inventory/place/config backend model and public read model

Wave 2: frontend infrastructure and shared state backbone
- Task 6: Frontend API client, Nuxt proxy/runtime config, and Pinia bootstrap
- Task 7: Session store plus `/login` `/register` `/header` auth UX
- Task 8: Catalog and booking stores with draft persistence and route/store/API matrix wiring

Wave 3: route migrations onto real data
- Task 9: Integrate `/lounge` with real backend catalog and booking submit
- Task 10: Integrate `/event` with real backend catalog and booking submit
- Task 11: Integrate `/gaming` with real backend catalog, place/config selection, and booking submit
- Task 12: Home route integration and cross-route auth return/draft restoration polish

Wave 4: docs and verification expansion
- Task 13: Write `docs/integration.md` and cross-link existing docs
- Task 14: Expand automated verification for backend contracts and frontend real flows

### Dependency Matrix (full, all tasks)
- **1**: Blocked By: none | Blocks: 3, 6, 7, 8, 9, 10, 11, 12, 14
- **2**: Blocked By: none | Blocks: 4, 5, 8, 9, 10, 11, 13, 14
- **3**: Blocked By: 1, 2 | Blocks: 8, 9, 10, 11, 12, 14
- **4**: Blocked By: 1, 2 | Blocks: 8, 9, 10, 12, 14
- **5**: Blocked By: 1, 2, 3 | Blocks: 8, 11, 12, 14
- **6**: Blocked By: 1 | Blocks: 7, 8, 9, 10, 11, 12, 14
- **7**: Blocked By: 1, 6 | Blocks: 8, 9, 10, 11, 12, 14
- **8**: Blocked By: 1, 2, 3, 4, 5, 6, 7 | Blocks: 9, 10, 11, 12, 14
- **9**: Blocked By: 4, 6, 7, 8 | Blocks: 12, 14
- **10**: Blocked By: 4, 6, 7, 8 | Blocks: 12, 14
- **11**: Blocked By: 5, 6, 7, 8 | Blocks: 12, 14
- **12**: Blocked By: 7, 8, 9, 10, 11 | Blocks: 13, 14
- **13**: Blocked By: 2, 12 | Blocks: F1-F4
- **14**: Blocked By: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 | Blocks: F1-F4

### Agent Dispatch Summary (wave → task count → categories)
- **Wave 1**: 5 tasks → `deep`, `unspecified-high`
- **Wave 2**: 3 tasks → `deep`, `quick`
- **Wave 3**: 4 tasks → `deep`, `visual-engineering`
- **Wave 4**: 2 tasks → `writing`, `unspecified-high`

## TODOs
> Implementation + Test = ONE task. Never separate.
> Every task below includes exact QA scenarios and no executor judgment calls.

- [ ] 1. Implement browser session contract and auth middleware

  **What to do**:
  - Keep Google OAuth as the only auth provider.
  - Replace the current “callback returns user JSON only” backend behavior with a browser session contract:
    - `GET /api/v1/auth/google?return_to=<frontend-path>` starts auth and preserves return path.
    - `GET /api/v1/auth/google/callback` completes Google auth, upserts the user, sets the browser session cookie, and redirects to the frontend `return_to`.
    - `GET /api/v1/auth/session` returns `200` with `{ authenticated: false, user: null }` when signed out and `200` with `{ authenticated: true, user: {...} }` when signed in.
    - `POST /api/v1/auth/logout` clears the session cookie and returns `204`.
    - `POST /api/v1/auth/dev-login` exists only in non-production and seeds the same session/cookie shape for automated QA.
  - Add backend auth middleware that resolves the session user and injects authenticated identity into protected handlers.
  - Keep catalog browsing public; protect booking/profile write routes only.

  **Must NOT do**:
  - Do not add email/password, JWT, OTP, or a second auth mechanism.
  - Do not require auth for browsing `/`, `/gaming`, `/lounge`, `/event`.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: cross-cutting auth/session/middleware/redirect work
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 3, 6, 7, 8, 9, 10, 11, 12, 14 | Blocked By: none

  **References**:
  - Pattern: `backend/internal/http/server/server.go:48-53` - current auth routes exist only as begin/callback and must expand into a browser session contract.
  - Pattern: `backend/internal/http/handlers/auth.go:24-100` - current OAuth begin/callback flow to extend instead of replacing.
  - Pattern: `backend/sql/queries/users.sql:1-15` - current user upsert behavior and client role assignment.
  - Requirement: `docs/third_stage_requirements.md:36-42` - authentication/authorization is a stated system requirement.
  - Requirement: `docs/sixth_stage_interfaces.md:7-13` - client interface explicitly includes OAuth login.
  - Prior plan: `.sisyphus/plans/google-auth-booking-ui-refresh.md:67-75` - same Google-only direction already established in a broader integration plan.

  **Acceptance Criteria**:
  - [ ] `GET /api/v1/auth/session` returns `200` plus `{ authenticated: false, user: null }` without a session cookie.
  - [ ] Successful Google callback or `dev-login` returns/sets a valid browser session cookie that `GET /api/v1/auth/session` recognizes.
  - [ ] `POST /api/v1/auth/logout` clears the cookie and the next `GET /api/v1/auth/session` returns signed-out state.
  - [ ] Protected booking/profile write routes reject unauthenticated requests with `401`.
  - [ ] `return_to` survives the auth round-trip.

  **QA Scenarios**:
  ```text
  Scenario: Deterministic dev login seeds a real browser session
    Tool: Bash
    Steps:
      1. POST /api/v1/auth/dev-login with the required non-production secret header and a client cookie jar.
      2. Reuse the cookie jar on GET /api/v1/auth/session.
      3. POST /api/v1/auth/logout with the same cookie jar.
      4. Re-check GET /api/v1/auth/session.
    Expected: Session endpoint reports authenticated after dev-login and unauthenticated after logout.
    Evidence: .sisyphus/evidence/task-1-session-contract.txt

  Scenario: Real OAuth initiation preserves return path
    Tool: Playwright
    Steps:
      1. Open http://127.0.0.1:3000/lounge.
      2. Attempt a protected booking submit while signed out.
      3. Assert navigation reaches the backend Google auth start URL with a return_to value for /lounge.
    Expected: Auth starts from backend and includes the intended frontend return path.
    Evidence: .sisyphus/evidence/task-1-auth-begin.png
  ```

  **Commit**: NO | Message: `feat(auth): add browser session contract` | Files: `backend/internal/http/handlers/auth.go`, `backend/internal/http/server/server.go`, new middleware/session files

- [ ] 2. Repair zone contracts and add metadata foundation for frontend read models

  **What to do**:
  - Fix the current zone schema/query/DTO mismatch by carrying `zone_tag_id` through zone DTOs, SQL queries, service layer, and HTTP handlers.
  - Remove stale `vip` enum acceptance from zone domain validation and align all code with the actual DB enum set `game|event|lounge|sys`.
  - Add a metadata strategy that can drive current frontend cards without inventing a CMS:
    - add `details_json JSONB NOT NULL DEFAULT '{}'::jsonb` to `zones`
    - add `details_json JSONB NOT NULL DEFAULT '{}'::jsonb` to `services`
  - Expose `zone_tags`, `site_settings`, and `computer_configurations` through backend query/service/handler layers so frontend read models can be assembled without hardcoded data.

  **Must NOT do**:
  - Do not introduce an admin editing UI in this task.
  - Do not create a full event-management subsystem.

  **Recommended Agent Profile**:
  - Category: `unspecified-high` - Reason: schema + sqlc + service + handler alignment
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 4, 5, 8, 9, 10, 11, 13, 14 | Blocked By: none

  **References**:
  - API/Type: `backend/sql/schema.sql:26-56` - `zone_tags`, `zones`, and `services` current shape and missing metadata support.
  - API/Type: `backend/sql/schema.sql:115-126` - existing `computer_configurations` and `site_settings` tables already available in schema.
  - Query: `backend/sql/queries/zones.sql:1-30` - current zone SQL omits `zone_tag_id` entirely.
  - API/Type: `backend/internal/domain/zone.go:7-21` - current DTO still allows removed `vip` type and omits `zone_tag_id`.
  - Requirement: `docs/frontend.md:35-45` - frontend stage ended with static screens ready for backend integration; this task creates the backend-side truth those screens need.

  **Acceptance Criteria**:
  - [ ] Zone create/get/patch contracts include `zone_tag_id` and align with the DB schema.
  - [ ] Zone validation no longer accepts `vip`.
  - [ ] `zones` and `services` can store frontend-facing metadata in `details_json`.
  - [ ] Query/service/handler layers exist for `zone_tags`, `site_settings`, and `computer_configurations`.
  - [ ] `go test ./...` covers the repaired zone contract and metadata read/write behavior.

  **QA Scenarios**:
  ```text
  Scenario: Zone create respects required zone_tag_id
    Tool: Bash
    Steps:
      1. Create a zone tag fixture.
      2. POST a zone with zone_tag_id and valid zone_type.
      3. GET the created zone.
    Expected: Create succeeds, returned payload includes zone_tag_id, and no schema error occurs.
    Evidence: .sisyphus/evidence/task-2-zone-tag-contract.txt

  Scenario: Invalid legacy vip type is rejected
    Tool: Bash
    Steps:
      1. POST a zone payload using type=vip.
    Expected: Request returns 400 with a validation error.
    Evidence: .sisyphus/evidence/task-2-zone-type-reject.txt
  ```

  **Commit**: NO | Message: `fix(zone): align schema and api contracts` | Files: `backend/sql/schema.sql`, migrations, `backend/sql/queries/zones.sql`, `backend/internal/domain/zone.go`, related service/handler files

- [ ] 3. Harden booking write contract and server-side invariants

  **What to do**:
  - Remove client authority over booking ownership:
    - booking create derives user identity from the authenticated session
    - booking patch/delete authorization also uses session identity/role rules instead of trusting payload
  - Extend bookings with the minimum fields required by real integrated flows:
    - `details_json JSONB NOT NULL DEFAULT '{}'::jsonb`
    - nullable `place_id BIGINT` for gaming seat/place bookings
    - immutable contact snapshot fields: `contact_name`, `contact_email`, `contact_phone`
  - Add server-side booking checks:
    - reject unauthenticated writes
    - reject overlapping bookings for the same resource/time window
    - reject gaming bookings when `place_id` is unavailable
    - reject participants that exceed zone/service capacity
    - reject duplicate active submits for same actor + service/place + start time
  - Add protected profile/contact update endpoint used by booking flows when Google profile lacks enough local contact data.

  **Must NOT do**:
  - Do not keep trusting `user_id` from the frontend.
  - Do not add my-bookings/history UI here; keep this task to write contracts and invariants.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: critical business invariants and auth coupling
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 8, 9, 10, 11, 12, 14 | Blocked By: 1, 2

  **References**:
  - API/Type: `backend/internal/domain/booking.go:7-25` - current create/patch DTO trusts client `user_id`.
  - API/Type: `backend/sql/schema.sql:58-70` - current bookings table lacks place/contact/details fields needed by integrated flows.
  - Requirement: `docs/third_stage_requirements.md:7-14` - bookings must check availability/prevent overlaps.
  - Requirement: `docs/sixth_stage_interfaces.md:9-13` - client interface includes booking creation and personal data editing.

  **Acceptance Criteria**:
  - [ ] Unauthenticated booking writes return `401`.
  - [ ] Booking create ignores/forbids client-supplied `user_id` and stores session-owned actor only.
  - [ ] Overlapping bookings for the same zone/service/place/time are rejected with a deterministic conflict response.
  - [ ] Booking writes persist `details_json`, `place_id` when relevant, and contact snapshot fields.
  - [ ] `go test ./...` covers overlap, duplicate-submit, unauthenticated write, and tampered payload cases.

  **QA Scenarios**:
  ```text
  Scenario: Tampered user_id cannot create someone else's booking
    Tool: Bash
    Steps:
      1. Seed an authenticated client session.
      2. POST a booking payload containing a different user_id value.
      3. Read the persisted booking or response payload.
    Expected: Server either rejects the payload or ignores the tampered user_id and assigns the session user.
    Evidence: .sisyphus/evidence/task-3-user-id-tamper.txt

  Scenario: Overlapping booking is rejected
    Tool: Bash
    Steps:
      1. Create one authenticated booking for a fixed resource/time window.
      2. Attempt a second booking on the same resource/time window.
    Expected: Second request returns conflict and does not persist a new booking.
    Evidence: .sisyphus/evidence/task-3-overlap-reject.txt
  ```

  **Commit**: NO | Message: `feat(booking): enforce server-side booking invariants` | Files: schema/migrations, booking domain/service/repo/handler files, auth middleware integration

- [ ] 4. Build public lounge and event catalog read-model APIs

  **What to do**:
  - Add public read endpoints under `/api/v1/public` for the existing client surfaces instead of forcing the frontend to compose raw CRUD responses.
  - Implement:
    - `GET /api/v1/public/home`
    - `GET /api/v1/public/lounge`
    - `GET /api/v1/public/event`
  - Shape the responses around current frontend needs, not raw tables:
    - home: cards/CTA-supporting summaries for gaming/lounge/event
    - lounge: zones, perks, mood, capacity/remaining, premium/accent, slot availability
    - event: event cards, date/time, format, availability, sold-out flag, display metadata
  - Derive these read models from `zones`, `services`, `zone_tags`, `site_settings`, and `details_json` rather than from mock files.

  **Must NOT do**:
  - Do not expose raw internal CRUD payloads directly as the frontend contract.
  - Do not create admin write endpoints for event authoring in this task.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: backend read-model design for multiple client routes
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 8, 9, 10, 12, 14 | Blocked By: 1, 2

  **References**:
  - Pattern: `frontend/app/pages/index.vue:74-100` - home cards that need real backend summaries.
  - Pattern: `frontend/app/pages/lounge.vue:67-268` - lounge catalog and booking form expect zones, slot availability, and display metadata.
  - Pattern: `frontend/app/pages/event.vue:68-246` - event catalog/registration expects event card metadata and sold-out state.
  - Shared types: `frontend/app/utils/experienceData.ts:35-151` - current mock shape to replace with real API DTOs.
  - Backend router: `backend/internal/http/server/server.go:48-93` - current `/api/v1` route placement to extend.

  **Acceptance Criteria**:
  - [ ] Guest callers can read `/api/v1/public/home`, `/api/v1/public/lounge`, and `/api/v1/public/event` with `200`.
  - [ ] Lounge payload contains enough fields to replace `loungeZones`, `loungeTimeSlots`, and associated card metadata.
  - [ ] Event payload contains enough fields to replace `eventItems` and sold-out/remaining logic.
  - [ ] Empty catalog states return stable empty payloads, not `500`s.

  **QA Scenarios**:
  ```text
  Scenario: Public lounge and event catalogs are guest-accessible
    Tool: Bash
    Steps:
      1. GET /api/v1/public/lounge without auth.
      2. GET /api/v1/public/event without auth.
    Expected: Both endpoints return 200 and stable JSON payloads for guest users.
    Evidence: .sisyphus/evidence/task-4-public-catalogs.txt

  Scenario: Empty-state payload is stable
    Tool: Bash
    Steps:
      1. Seed or configure a test fixture with no active lounge or event entries.
      2. GET the corresponding public endpoint.
    Expected: Response is a valid empty payload the frontend can render as empty state.
    Evidence: .sisyphus/evidence/task-4-empty-payloads.txt
  ```

  **Commit**: NO | Message: `feat(api): add public lounge and event read models` | Files: new public handler/service/query files, router registration, related SQL/query changes

- [ ] 5. Add gaming place/config model and public gaming read-model API

  **What to do**:
  - Extend the backend to support the existing gaming UI surface instead of leaving it as static local state.
  - Add the minimum gaming persistence model:
    - `zone_places` table with `zone_id`, `label`, `configuration_id`, `sort_order`, `is_active`
    - use existing `computer_configurations` as reusable config templates or align it into a direct place/config relationship
    - make bookings optionally reference `place_id`
  - Implement `GET /api/v1/public/gaming` that returns:
    - gaming tabs/classes/offers needed by the route
    - zones and places
    - configuration/spec summaries for each place or place template
    - availability grid for the requested date/window
  - Implement any backend helper endpoints needed for validated gaming availability lookup, but keep booking submission on the main booking write contract from Task 3.

  **Must NOT do**:
  - Do not build tournaments, rankings, payments, or admin editing UI.
  - Do not invent a second booking system separate from the shared booking backbone.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: new gaming inventory/place/config modeling with availability derivation
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 8, 11, 12, 14 | Blocked By: 1, 2, 3

  **References**:
  - Pattern: `frontend/app/pages/gaming.vue:146-220` - current gaming booking panel expects date, zone, place, and per-place selection.
  - API/Type: `backend/sql/schema.sql:115-126` - current backend already has `computer_configurations` and `site_settings`, but lacks place inventory.
  - Requirement: `docs/third_stage_requirements.md:9-12` - system must manage game zones/resources and prevent overlapping bookings.
  - Requirement: `docs/sixth_stage_interfaces.md:10-13` - client can view available game zones/services and create bookings.

  **Acceptance Criteria**:
  - [ ] Backend has a persisted concept of gaming places/seats linked to zones/configuration.
  - [ ] `GET /api/v1/public/gaming` returns enough data to replace local `zones`, `hours`, and place/config display logic.
  - [ ] Gaming availability for a selected date/time is computed from persisted bookings and place state.
  - [ ] Conflicting gaming bookings for the same place/time cannot be created.

  **QA Scenarios**:
  ```text
  Scenario: Gaming catalog returns place/config inventory
    Tool: Bash
    Steps:
      1. GET /api/v1/public/gaming.
    Expected: Response includes zones, place labels, configuration/spec information, and availability data for frontend rendering.
    Evidence: .sisyphus/evidence/task-5-gaming-catalog.txt

  Scenario: Conflicting place booking is rejected
    Tool: Bash
    Steps:
      1. Create an authenticated booking for one gaming place/time slot.
      2. Attempt a second booking for the same place/time.
    Expected: Second request is rejected with conflict.
    Evidence: .sisyphus/evidence/task-5-gaming-conflict.txt
  ```

  **Commit**: NO | Message: `feat(gaming): add place inventory and public read model` | Files: schema/migrations, new gaming queries/services/handlers, booking integration changes

- [ ] 6. Add frontend API client, same-origin proxy/runtime config, and Pinia bootstrap

  **What to do**:
  - Add Pinia to the Nuxt frontend and register it in the app.
  - Add one shared API wrapper for frontend calls (`$fetch`-based or equivalent) that centralizes base URL, credentials, typed error handling, and request defaults.
  - Configure frontend/backend communication for same-origin browser behavior in local/dev so cookie auth works without ad hoc per-call hacks.
  - Add typed frontend API modules for `auth/session`, `public catalogs`, and `booking`.
  - Keep the API layer thin: request/response mapping, credentials, and typed errors only.

  **Must NOT do**:
  - Do not scatter direct ad hoc `fetch` calls across page components.
  - Do not put backend business rules into the API client.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: cross-cutting frontend infrastructure and auth/cookie behavior
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 7, 8, 9, 10, 11, 12, 14 | Blocked By: 1

  **References**:
  - Pattern: `frontend/nuxt.config.ts:1-27` - current Nuxt config has no API/runtime integration and must be extended carefully.
  - Pattern: `frontend/package.json:5-13` - current scripts and package setup to extend with Pinia dependency only.
  - Pattern: `frontend/app/app.vue:2-7` - current app shell where shared app bootstrap behavior is introduced.
  - Finding: `frontend` currently has no `defineStore`, no `fetch`, and no API client layer.

  **Acceptance Criteria**:
  - [ ] Pinia is installed and registered in Nuxt.
  - [ ] Frontend has one shared API client layer used by stores instead of pages.
  - [ ] Same-origin/cookie credentials work for session-protected calls in local dev.
  - [ ] `lint`, `typecheck`, and smoke still run after the infrastructure addition.

  **QA Scenarios**:
  ```text
  Scenario: Frontend can call authenticated and public APIs with one client layer
    Tool: Bash
    Steps:
      1. Run frontend lint and typecheck.
      2. Start frontend dev server.
      3. Exercise one public catalog call and one authenticated session call through the browser app.
    Expected: Public calls work signed out and authenticated calls include browser credentials without per-page fetch hacks.
    Evidence: .sisyphus/evidence/task-6-api-layer.txt

  Scenario: Missing backend session is handled as typed signed-out state
    Tool: Playwright
    Steps:
      1. Open /login signed out.
      2. Assert app boot does not crash if session endpoint returns signed-out state.
    Expected: Frontend handles signed-out bootstrap gracefully.
    Evidence: .sisyphus/evidence/task-6-signed-out-bootstrap.png
  ```

  **Commit**: NO | Message: `feat(frontend): add pinia and shared api client` | Files: `frontend/nuxt.config.ts`, `frontend/package.json`, new frontend API/store bootstrap files

- [ ] 7. Implement session store and Google-only auth UX on login/register/header

  **What to do**:
  - Add a Pinia `session` store that owns:
    - `authenticated` flag
    - current user object
    - session bootstrap/loading/error state
    - login start, logout, and return-to handling
  - Convert `/login` and `/register` into Google-only entry surfaces:
    - `/login` copy = sign in with Google
    - `/register` copy = create account with Google
    - both routes call the same backend auth begin endpoint with `return_to`
  - If already authenticated, visiting `/login` or `/register` redirects to `return_to` or `/`.
  - Update `AppHeader` so it reflects signed-out vs signed-in state and exposes logout.

  **Must NOT do**:
  - Do not keep dead local credential form submissions.
  - Do not add profile or my-bookings pages in this task.

  **Recommended Agent Profile**:
  - Category: `quick` - Reason: bounded auth-aware frontend surfaces over the new session store
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 8, 9, 10, 11, 12, 14 | Blocked By: 1, 6

  **References**:
  - Pattern: `frontend/app/pages/login.vue:1-89` - current login surface is visual-only and must become Google-only.
  - Pattern: `frontend/app/pages/register.vue:1-92` - current register surface is visual-only and must become Google-only.
  - Pattern: `frontend/app/components/AppHeader.vue:17-44` - current public header actions to make session-aware.
  - Requirement: `docs/sixth_stage_interfaces.md:7-13` - OAuth login is required by client interface design.

  **Acceptance Criteria**:
  - [ ] `/login` and `/register` each render one Google CTA and no local credential submit flow.
  - [ ] Signed-out header shows Google entry action(s); signed-in header shows user summary and logout.
  - [ ] Already-authenticated visits to `/login` and `/register` redirect away correctly.
  - [ ] Logout clears session state and returns the shell to signed-out mode.

  **QA Scenarios**:
  ```text
  Scenario: Google-only auth surfaces behave correctly
    Tool: Playwright
    Steps:
      1. Open /login signed out.
      2. Assert exactly one primary Google CTA is visible.
      3. Open /register signed out.
      4. Assert exactly one primary Google CTA is visible.
    Expected: No local credential submission path remains.
    Evidence: .sisyphus/evidence/task-7-google-only-auth.png

  Scenario: Authenticated user is redirected away from auth pages
    Tool: Playwright
    Steps:
      1. Seed authenticated session with dev-login.
      2. Open /login and then /register.
      3. Assert browser redirects to / or return_to destination.
    Expected: Auth pages do not remain active for authenticated users.
    Evidence: .sisyphus/evidence/task-7-auth-redirects.png
  ```

  **Commit**: NO | Message: `feat(frontend): add session store and google auth ux` | Files: login/register/header files, new session store, app bootstrap files

- [ ] 8. Implement shared Pinia catalog and booking stores with redirect-safe drafts

  **What to do**:
  - Add a Pinia `catalog` store that loads and caches the public read models for home/lounge/event/gaming.
  - Add a Pinia `booking` store that owns:
    - lounge draft fields
    - event draft fields
    - gaming draft fields
    - shared pending/error/success state
    - return-to and rehydrate logic for booking after auth redirect
  - Persist only the booking draft/intent client-side in `sessionStorage` so auth redirects do not destroy in-progress work.
  - Keep ephemeral purely visual UI state local when it does not need sharing or persistence.

  **Must NOT do**:
  - Do not mirror every presentational ref into Pinia.
  - Do not keep mock dataset files as the production source of truth.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: central application state and auth-return flow logic
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 9, 10, 11, 12, 14 | Blocked By: 1, 2, 3, 4, 5, 6, 7

  **References**:
  - Pattern: `frontend/app/composables/useExperienceMockState.ts:4-30` - current query-state mock mechanism to retire or restrict to dev-only fixtures.
  - Pattern: `frontend/app/composables/useMockSubmission.ts:1-28` - current fake submit path to replace with real store actions.
  - Pattern: `frontend/app/utils/experienceData.ts:35-151` - current mock catalog data to replace with API-driven state.
  - Pattern: `frontend/app/pages/lounge.vue:295-412` - current local lounge draft/submit logic to migrate.
  - Pattern: `frontend/app/pages/event.vue:271-354` - current local event draft/submit logic to migrate.
  - Pattern: `frontend/app/pages/gaming.vue:328-425` - current gaming local state that needs store-backed integration.

  **Acceptance Criteria**:
  - [ ] Catalog data for home/lounge/event/gaming is loaded through the `catalog` store, not local mocks.
  - [ ] Booking drafts survive auth redirects through `sessionStorage` and restore into the correct route.
  - [ ] Booking store exposes deterministic pending, error, and success states used by all three route flows.
  - [ ] Local mock composables are no longer the production data/submit path for integrated routes.

  **QA Scenarios**:
  ```text
  Scenario: Lounge draft survives auth redirect
    Tool: Playwright
    Steps:
      1. Open /lounge signed out.
      2. Fill the booking draft and trigger protected submit.
      3. Complete deterministic auth flow.
      4. Return to /lounge.
    Expected: Draft state is restored from store/sessionStorage and remains ready for submit.
    Evidence: .sisyphus/evidence/task-8-lounge-draft-restore.png

  Scenario: Event and gaming drafts remain isolated
    Tool: Playwright
    Steps:
      1. Create a draft on /event.
      2. Create a different draft on /gaming.
      3. Navigate back and forth between routes.
    Expected: Each route restores only its own draft state.
    Evidence: .sisyphus/evidence/task-8-draft-isolation.png
  ```

  **Commit**: NO | Message: `feat(frontend): add catalog and booking pinia stores` | Files: new stores, API modules, integrated route state files

- [ ] 9. Integrate `/lounge` with real backend catalog and booking flow

  **What to do**:
  - Replace local lounge mocks with `catalog` and `booking` store data/actions.
  - Keep the existing route structure, selectors, and UX rhythm where feasible.
  - Render real lounge zones, perks, mood, availability, slot selection, and booking summary from backend data.
  - On protected submit:
    - if signed out, preserve draft and start auth flow
    - if signed in, submit real booking to backend and surface real success/error messages
  - Keep deterministic loading/empty/error/success UI states, but they must now reflect real API/store states rather than mock query fixtures for the production path.

  **Must NOT do**:
  - Do not keep `experienceData.ts` as the live source for lounge.
  - Do not change the route path or remove the current booking selectors without replacing them in smoke coverage.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: route-level integration over real catalog/store/booking behavior
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: 12, 14 | Blocked By: 4, 6, 7, 8

  **References**:
  - Pattern: `frontend/app/pages/lounge.vue:67-268` - current lounge shell and form layout to preserve.
  - Pattern: `frontend/app/pages/lounge.vue:295-412` - current local draft/submit logic to migrate to stores.
  - Shared UI: `frontend/app/components/ExperienceHero.vue:1-76`, `ExperienceFlowPanel.vue:1-27`, `ExperienceCard.vue:1-91`, `ExperienceStatePanel.vue:1-35`.
  - Smoke: `frontend/tests/smoke.spec.ts:3-70` - current selectors and flow checks to preserve/update.

  **Acceptance Criteria**:
  - [ ] `/lounge` loads from backend-backed store data with no production dependency on hardcoded lounge mocks.
  - [ ] Signed-out protected submit redirects into auth with draft preserved.
  - [ ] Signed-in protected submit creates a real booking and surfaces backend success/conflict/validation errors.
  - [ ] Empty/loading/error lounge states are driven by real request/store state.

  **QA Scenarios**:
  ```text
  Scenario: Authenticated lounge booking succeeds and persists
    Tool: Playwright
    Steps:
      1. Seed authenticated session.
      2. Open /lounge.
      3. Select an available zone/slot, fill required fields, and submit.
      4. Verify success UI.
      5. Verify persisted booking via API or confirmation payload.
    Expected: Real lounge booking is created and reflected in UI/backend.
    Evidence: .sisyphus/evidence/task-9-lounge-real-booking.png

  Scenario: Signed-out lounge submit preserves draft and enters auth flow
    Tool: Playwright
    Steps:
      1. Open /lounge signed out.
      2. Fill draft and submit.
      3. Assert redirect to Google auth start path with return_to.
    Expected: User is not booked while signed out, and draft is preserved for post-auth recovery.
    Evidence: .sisyphus/evidence/task-9-lounge-auth-gate.png
  ```

  **Commit**: NO | Message: `feat(lounge): connect lounge flow to backend` | Files: `frontend/app/pages/lounge.vue`, related stores/API modules/tests

- [ ] 10. Integrate `/event` with real backend catalog and booking flow

  **What to do**:
  - Replace local event mocks with backend-backed event read-model data through the `catalog` store.
  - Preserve the route’s current card/listing/registration structure while switching to real data.
  - Submit real backend bookings/registrations using the shared booking store and booking API contract.
  - Surface sold-out and validation/conflict behavior from the backend rather than from local mock logic.

  **Must NOT do**:
  - Do not introduce a staff/admin event editor.
  - Do not create a separate event booking subsystem apart from shared bookings.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: route-level real data integration with availability and sold-out semantics
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: 12, 14 | Blocked By: 4, 6, 7, 8

  **References**:
  - Pattern: `frontend/app/pages/event.vue:68-246` - current event shell and registration form layout.
  - Pattern: `frontend/app/pages/event.vue:271-354` - current local registration draft/submit logic to migrate.
  - Shared types: `frontend/app/utils/experienceData.ts:107-151` - current event mock shape to replace with backend response.
  - Requirement: `docs/sixth_stage_interfaces.md:9-13` - client interface includes viewing services and creating bookings.

  **Acceptance Criteria**:
  - [ ] `/event` loads from backend-backed event catalog data.
  - [ ] Sold-out/unavailable entries are driven by backend availability data.
  - [ ] Authenticated event submit creates a real booking/registration entry.
  - [ ] Signed-out submit uses the same auth-gated draft-preserving behavior as lounge.

  **QA Scenarios**:
  ```text
  Scenario: Authenticated event booking succeeds
    Tool: Playwright
    Steps:
      1. Seed authenticated session.
      2. Open /event.
      3. Select an available event, fill required fields, and submit.
      4. Verify success UI and persisted booking.
    Expected: Real event registration is created and confirmed through UI/backend.
    Evidence: .sisyphus/evidence/task-10-event-real-booking.png

  Scenario: Sold-out event cannot be booked
    Tool: Playwright
    Steps:
      1. Open /event with backend fixture or seeded data containing a sold-out item.
      2. Attempt to select/submit the sold-out event.
    Expected: UI blocks the action and backend does not create a booking.
    Evidence: .sisyphus/evidence/task-10-event-sold-out.png
  ```

  **Commit**: NO | Message: `feat(event): connect event flow to backend` | Files: `frontend/app/pages/event.vue`, related stores/API modules/tests

- [ ] 11. Integrate `/gaming` with real backend catalog, place/config selection, and booking flow

  **What to do**:
  - Replace local gaming data (`zones`, `hours`, `dates`, place lists) with backend-backed gaming catalog/store state.
  - Keep the current route as a client surface for:
    - choosing gaming class/zone
    - choosing a specific place/computer
    - choosing a date/time window
    - submitting a real booking
  - Map the UI to the shared booking contract while populating `place_id` and gaming-specific metadata in `details_json`.
  - Surface real place conflicts and unavailable slots from backend availability.

  **Must NOT do**:
  - Do not rebuild the whole gaming page visually.
  - Do not ship tournaments/leaderboards/admin controls.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: largest route migration with new backend inventory semantics
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: 12, 14 | Blocked By: 5, 6, 7, 8

  **References**:
  - Pattern: `frontend/app/pages/gaming.vue:146-220` - current booking panel UX and selectors to preserve.
  - Pattern: `frontend/app/pages/gaming.vue:328-425` - local zone/place/hour state that must become backend-backed.
  - Backend foundation: `backend/sql/schema.sql:115-126` - current configuration tables that gaming integration extends.

  **Acceptance Criteria**:
  - [ ] `/gaming` loads real zones/places/configurations/availability from backend-backed store data.
  - [ ] User can select a specific gaming place and time window and submit a real booking.
  - [ ] Availability conflicts from backend are reflected in UI and block invalid booking.
  - [ ] No production dependency on hardcoded local gaming inventory remains.

  **QA Scenarios**:
  ```text
  Scenario: Authenticated gaming place booking succeeds
    Tool: Playwright
    Steps:
      1. Seed authenticated session.
      2. Open /gaming.
      3. Choose zone, place, and valid time window.
      4. Submit booking.
      5. Verify success UI and persisted booking with place_id.
    Expected: Real gaming booking is created against a specific place.
    Evidence: .sisyphus/evidence/task-11-gaming-real-booking.png

  Scenario: Occupied gaming place/time is blocked
    Tool: Playwright
    Steps:
      1. Seed backend data with an occupied gaming place/time.
      2. Open /gaming and attempt to book that exact place/time.
    Expected: UI and backend reject the action consistently.
    Evidence: .sisyphus/evidence/task-11-gaming-conflict.png
  ```

  **Commit**: NO | Message: `feat(gaming): connect gaming flow to backend` | Files: `frontend/app/pages/gaming.vue`, related stores/API modules/tests, backend contract usage updates

- [ ] 12. Integrate home route and cross-route auth return/draft restoration behavior

  **What to do**:
  - Connect `/` to the real home catalog read model instead of purely static assumptions where catalog summaries are shown.
  - Ensure homepage CTAs route into live integrated surfaces without selector regressions.
  - Centralize post-auth return handling so drafts restore correctly to `/lounge`, `/event`, or `/gaming`.
  - Ensure the shell stays coherent when auth state changes across public routes.

  **Must NOT do**:
  - Do not redesign the homepage.
  - Do not add unrelated marketing sections or profile/dashboard flows.

  **Recommended Agent Profile**:
  - Category: `quick` - Reason: bounded route wiring and shell consistency once core flows are integrated
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: NO | Wave 3 | Blocks: 13, 14 | Blocked By: 7, 8, 9, 10, 11

  **References**:
  - Pattern: `frontend/app/pages/index.vue:74-100` - homepage CTA links and card surface.
  - Pattern: `frontend/app/components/AppHeader.vue:17-44` - auth-aware nav/action consistency.
  - Pattern: `frontend/app/app.vue:2-7` - shared shell where auth-aware global behavior lives.

  **Acceptance Criteria**:
  - [ ] Homepage still routes to `/gaming`, `/lounge`, and `/event` using stable selectors.
  - [ ] Return-to restoration works for all three booking routes after auth.
  - [ ] Header state stays correct across public routes after login/logout.

  **QA Scenarios**:
  ```text
  Scenario: Homepage CTA routing remains stable
    Tool: Playwright
    Steps:
      1. Open /.
      2. Click each CTA in turn: gaming, lounge, event.
      3. Assert route transitions succeed and page shells load.
    Expected: Homepage remains the stable entry point to all integrated routes.
    Evidence: .sisyphus/evidence/task-12-home-cta-routing.png

  Scenario: Draft restoration works across all booking routes
    Tool: Playwright
    Steps:
      1. Create partial drafts for lounge, event, and gaming in separate runs while signed out.
      2. Trigger protected submit and complete deterministic auth.
      3. Assert each route restores the correct draft after return.
    Expected: Return-to flow restores the right draft for each route without cross-route leakage.
    Evidence: .sisyphus/evidence/task-12-return-to-restoration.png
  ```

  **Commit**: NO | Message: `feat(home): connect home and auth return polish` | Files: `frontend/app/pages/index.vue`, shell/header/store glue, updated smoke tests

- [ ] 13. Write `docs/integration.md` and cross-link existing docs

  **What to do**:
  - Create a new canonical implementation write-up at `docs/integration.md`.
  - The document must include, at minimum:
    - purpose and scope of the integration pass
    - supersedes/complements note for `docs/frontend.md` and old mock-flow plans
    - architecture decisions (Google-only auth, same-origin cookie session, Pinia store boundaries)
    - backend endpoint inventory for public and protected flows
    - route → store → endpoint mapping table for `/`, `/login`, `/register`, `/lounge`, `/event`, `/gaming`
    - DB/model changes introduced for session/catalog/booking/gaming
    - env vars and local setup instructions
    - verification commands and expected outcomes
    - known limitations and explicit out-of-scope items
  - Add small cross-links from `docs/frontend.md` and relevant backend docs so the integration doc becomes discoverable.

  **Must NOT do**:
  - Do not dump raw changelog text without structure.
  - Do not leave doc ownership ambiguous; `docs/integration.md` is the canonical implementation write-up.

  **Recommended Agent Profile**:
  - Category: `writing` - Reason: technical documentation and architecture write-up
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 4 | Blocks: F1-F4 | Blocked By: 2, 12

  **References**:
  - Requirement: `docs/frontend.md:3-5` - current frontend doc still describes the pre-integration static phase and must be complemented.
  - Requirement: `docs/third_stage_requirements.md:7-14, 36-42` - booking, auth, and security requirements to reflect in the doc.
  - Requirement: `docs/sixth_stage_interfaces.md:7-13` - client interface scope to map against real integrated routes.
  - Plan context: `.sisyphus/plans/frontend-continuation.md` and `.sisyphus/plans/google-auth-booking-ui-refresh.md` - existing plan history to supersede/relate explicitly.

  **Acceptance Criteria**:
  - [ ] `docs/integration.md` exists and contains all required sections above.
  - [ ] The document contains a route → store → endpoint mapping table covering all in-scope client routes.
  - [ ] Existing docs reference the new integration document where appropriate.
  - [ ] Documentation is sufficient for a new engineer to run verification commands and understand the integration topology.

  **QA Scenarios**:
  ```text
  Scenario: Integration doc contains required sections
    Tool: Bash
    Steps:
      1. Verify `docs/integration.md` exists.
      2. Search for required headings/sections: scope, auth/session, pinia stores, route/store/endpoint matrix, env vars, verification, limitations.
    Expected: All required sections are present.
    Evidence: .sisyphus/evidence/task-13-doc-sections.txt

  Scenario: Cross-links from existing docs are present
    Tool: Bash
    Steps:
      1. Search `docs/frontend.md` and touched backend docs for links or references to `docs/integration.md`.
    Expected: Existing docs point readers to the new canonical integration write-up.
    Evidence: .sisyphus/evidence/task-13-doc-links.txt
  ```

  **Commit**: NO | Message: `docs(integration): add frontend-backend integration guide` | Files: `docs/integration.md`, targeted docs cross-links

- [ ] 14. Expand automated verification for backend contracts and real frontend flows

  **What to do**:
  - Add backend tests and/or contract checks covering the new session/catalog/booking logic:
    - auth session contract
    - zone contract repair
    - booking overlap/ownership/gaming conflict behavior
  - Expand Playwright smoke to verify real integrated behavior, not mock-only flows:
    - signed-out browse on `/`, `/gaming`, `/lounge`, `/event`
    - auth gating on protected submit
    - authenticated booking success on lounge/event/gaming
    - logout behavior
    - already-authenticated behavior on `/login` and `/register`
  - Add stable `curl` verification scripts or documented command blocks using cookie jars.
  - Save evidence under `.sisyphus/evidence/` and `.sisyphus/evidence/final-qa/`.

  **Must NOT do**:
  - Do not depend on manual Google login.
  - Do not leave endpoint names/status codes/response shapes undecided in the tests.

  **Recommended Agent Profile**:
  - Category: `unspecified-high` - Reason: multi-layer verification expansion across backend and frontend
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: NO | Wave 4 | Blocks: F1-F4 | Blocked By: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12

  **References**:
  - Pattern: `frontend/tests/smoke.spec.ts:3-70` - existing smoke file to expand from mock checks to real integrated checks.
  - Pattern: `frontend/playwright.config.ts:3-23` - current Playwright harness and web-server startup.
  - Pattern: `frontend/package.json:5-13` - current `smoke` script entry point.
  - Requirement: `docs/third_stage_requirements.md:10-12, 36-40` - availability prevention and auth/security must be verifiable.

  **Acceptance Criteria**:
  - [ ] `go test ./...` passes and covers new backend contract/invariant behavior.
  - [ ] `corepack pnpm run smoke` passes against real backend-integrated flows.
  - [ ] Cookie-jar `curl` checks exist and prove public/private endpoint behavior.
  - [ ] Evidence files exist for auth, catalog, booking, gaming, and docs verification.

  **QA Scenarios**:
  ```text
  Scenario: Protected and public API contracts are binary-verifiable
    Tool: Bash
    Steps:
      1. GET public catalog endpoints without auth.
      2. POST booking without auth.
      3. Seed authenticated session.
      4. POST booking with auth.
      5. POST logout and re-check session endpoint.
    Expected: Public reads succeed, unauthenticated write rejects, authenticated write succeeds, logout clears session.
    Evidence: .sisyphus/evidence/task-14-api-contracts.txt

  Scenario: Real integrated smoke suite passes
    Tool: Bash + Playwright
    Steps:
      1. Run `go test ./...` in backend.
      2. Run `corepack pnpm run lint`, `typecheck`, and `smoke` in frontend.
      3. Collect generated smoke evidence and summary logs.
    Expected: All verification commands pass against the integrated app.
    Evidence: .sisyphus/evidence/task-14-integrated-smoke.txt
  ```

  **Commit**: NO | Message: `test(integration): verify real frontend-backend flows` | Files: backend tests, frontend smoke tests, QA scripts/evidence files

## Final Verification Wave (MANDATORY — after ALL implementation tasks)
> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.
> **Do NOT auto-proceed after verification. Wait for user's explicit approval before marking work complete.**
> **Never mark F1-F4 as checked before getting user's okay.** Rejection or user feedback -> fix -> re-run -> present again -> wait for okay.
- [ ] F1. Plan Compliance Audit — `oracle`
- [ ] F2. Code Quality Review — `unspecified-high`
- [ ] F3. Real Manual QA — `unspecified-high` (+ Playwright)
- [ ] F4. Scope Fidelity Check — `deep`

## Commit Strategy
- No commits during execution unless the user explicitly requests them.
- If the user requests a commit after implementation passes verification, use one final commit: `feat(integration): connect frontend flows to backend with pinia`

## Success Criteria
### Verification Commands
```bash
cd backend && go test ./...
cd frontend && corepack pnpm run lint
cd frontend && corepack pnpm run typecheck
cd frontend && corepack pnpm run smoke
```

### Final Checklist
- [ ] All in-scope frontend routes use real backend data where planned
- [ ] Google cookie session is usable by browser flows and automated QA
- [ ] Booking ownership comes from session, not client payload
- [ ] Gaming availability is backed by real place/config data
- [ ] `docs/integration.md` is present and complete
- [ ] All automated verification passes
