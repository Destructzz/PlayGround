# Google-Linked Booking & Desktop UI Refresh

## TL;DR
> **Summary**: Convert the current mock-driven lounge/event frontend into a Google-authenticated booking flow backed by the existing PlayGround backend, while refreshing desktop UI with Nuxt UI components and tighter radii.
> **Deliverables**:
> - Cookie-based Google auth session for frontend/backend integration
> - First-booking phone capture with DB prefill and inline editability
> - Immutable booking-time contact snapshots stored with each booking
> - Auth-aware lounge and event flows on the existing shared booking backbone
> - Desktop UI refresh built from Nuxt UI primitives with reduced rounding
> - Minimal smoke maintenance only; no new test framework or suite expansion
> **Effort**: Large
> **Parallel**: YES - 3 waves
> **Critical Path**: 1 → 2 → 3 → 4 → 6 → 7 → 8 → 9

## Context
### Original Request
Improve the frontend around lounge/event selection, make Google account the required identity for all bookings, collect only Google + phone at account/booking entry, defer all other profile completion to booking/purchase with DB-prefill and inline editing, and make desktop UI look more polished using existing Nuxt UI components with smaller radii.

### Interview Summary
- Google account is mandatory for every booking; no guest flow.
- Phone is **not** collected immediately after Google login; it is collected only at the first booking/purchase attempt if missing.
- Prefilled user data must remain editable directly inside the booking/purchase form.
- Scope includes both frontend and the required backend/API/DB work.
- Desktop surfaces should be rebuilt with framework-native Nuxt UI components and less oversized rounding.
- Automated coverage must stay bounded: no new framework or suite expansion beyond maintaining the existing smoke path.

### Metis Review (gaps addressed)
- Locked auth gate to the booking submit action, not page entry.
- Locked frontend↔backend auth contract to cookie-backed session, return-path preservation, and explicit auth-state endpoint.
- Locked field ownership rules:
  - Google-owned canonical fields: `google_id`, `email`, `avatar_url`
  - User-editable canonical fields: `full_name`, `phone`
  - Booking snapshot fields: `contact_name`, `contact_email`, `contact_phone`
- Locked scope to **new bookings only**; no legacy user/bookings migration or backfill.
- Locked desktop refresh to lounge/event/auth/header surfaces only; mobile must remain functional but is not a redesign target.

### Oracle Review (gaps addressed)
- Reuse the existing unified `users` / `bookings` / `zones` model; do not split lounge/event account systems.
- Add immutable booking-time contact snapshots so later profile edits do not rewrite history.
- Make phone collection a booking-time readiness step.
- Centralize UI refresh through shared Nuxt UI patterns and reduced radii, not ad hoc per-page styling.
- Add explicit duplicate-submit, session-expiry, and callback-failure handling.

## Work Objectives
### Core Objective
Ship one coherent booking identity system where browsing lounge/event pages remains public, but completing any booking requires a Google-authenticated session, captures phone on the first booking if missing, prefills later bookings from DB, and presents the entire desktop journey through cleaned-up Nuxt UI-based interfaces.

### Deliverables
- Backend auth/session endpoints and middleware for frontend-driven Google login
- Backend schema updates for booking contact snapshots and per-flow details
- Backend `me` endpoints for auth state and contact updates
- Frontend auth/session composables and booking-intent persistence
- Google-only `/login` and `/register` entry surfaces
- Auth-aware header/account actions
- Refreshed lounge and event booking UIs using Nuxt UI primitives
- Updated existing smoke test file to match the new auth-aware flow without expanding suite scope

### Definition of Done (verifiable conditions with commands)
- [ ] `go test ./...` passes in `/home/destruct/Projects/Web/PlayGround/backend`
- [ ] `corepack pnpm run lint` passes in `/home/destruct/Projects/Web/PlayGround/frontend`
- [ ] `corepack pnpm run typecheck` passes in `/home/destruct/Projects/Web/PlayGround/frontend`
- [ ] `corepack pnpm run smoke` passes in `/home/destruct/Projects/Web/PlayGround/frontend`
- [ ] `curl` verification proves unauthenticated booking create is rejected and authenticated create succeeds in non-production QA mode
- [ ] Lounge and event both preserve booking intent across auth return and show deterministic error recovery states

### Must Have
- Google-only booking identity with no guest booking path
- Cookie-backed session endpoint usable by Nuxt frontend
- Booking submit gate that preserves the user’s selected zone/event context across auth redirects
- First-booking phone capture inline in the booking flow if `users.phone` is missing
- DB-prefill for later bookings with inline editability before submit
- Immutable booking contact snapshots stored per booking
- Shared account/session backbone for both lounge and event
- Desktop UI refresh using existing Nuxt UI components and smaller radii

### Must NOT Have (guardrails, AI slop patterns, scope boundaries)
- No email/password auth, guest flow, SMS OTP, or multi-provider OAuth expansion
- No separate account or booking systems for lounge vs event
- No legacy-user migration/backfill project
- No broad mobile redesign, global rebrand, or full design-system rewrite
- No new test framework, no second Playwright suite, and no CI rollout in this iteration
- No manual-only acceptance criteria such as “looks better” without selectors/commands/evidence

## Verification Strategy
> ZERO HUMAN INTERVENTION - all verification is agent-executed.
- Test decision: **tests-after on existing infrastructure only**
- Frameworks/tools: existing Playwright smoke + backend `go test` + frontend lint/typecheck + `curl`
- Auth QA rule: because real Google credentials are not agent-safe, add a **non-production-only dev session seeding endpoint** that produces the same session cookie shape as the real callback; use it for automated post-auth verification, while still verifying real OAuth redirect initiation.
- QA policy: every task must include happy-path and failure-path scenarios with exact selectors/commands.
- Evidence: `.sisyphus/evidence/task-{N}-{slug}.{ext}`

## Execution Strategy
### Parallel Execution Waves
> Target: 5-8 tasks per wave. Shared auth/backend foundations are extracted into Wave 1 for maximum parallelism.

Wave 1: backend/session contract + data model + frontend session plumbing
- Task 1: Backend Google session contract and auth middleware
- Task 2: Backend user/profile and booking schema hardening
- Task 3: Frontend auth session and booking-intent state

Wave 2: auth surfaces + shared desktop UI + lounge conversion
- Task 4: Google-only auth entry surfaces and header/account state
- Task 5: Shared desktop UI foundation with Nuxt UI and reduced radii
- Task 6: Lounge booking flow integration

Wave 3: event conversion + recovery hardening + smoke maintenance
- Task 7: Event booking flow integration
- Task 8: Cross-flow recovery/error-state hardening
- Task 9: Existing smoke maintenance and final project gates

### Dependency Matrix (full, all tasks)
- **1**: Blocked By: none | Blocks: 3, 4, 6, 7, 8, 9
- **2**: Blocked By: none | Blocks: 6, 7, 8, 9
- **3**: Blocked By: 1 | Blocks: 4, 6, 7, 8, 9
- **4**: Blocked By: 1, 3 | Blocks: 6, 7
- **5**: Blocked By: none | Blocks: 6, 7, 8
- **6**: Blocked By: 1, 2, 3, 4, 5 | Blocks: 8, 9
- **7**: Blocked By: 1, 2, 3, 4, 5 | Blocks: 8, 9
- **8**: Blocked By: 1, 2, 3, 5, 6, 7 | Blocks: 9
- **9**: Blocked By: 1, 2, 3, 6, 7, 8 | Blocks: F1-F4

### Agent Dispatch Summary (wave → task count → categories)
- **Wave 1**: 3 tasks → `deep`, `unspecified-high`, `quick`
- **Wave 2**: 3 tasks → `quick`, `visual-engineering`, `deep`
- **Wave 3**: 3 tasks → `deep`, `unspecified-high`, `quick`

## TODOs
> Implementation + Test = ONE task. Never separate.
> Every task below includes exact QA scenarios and no executor judgment calls.

- [ ] 1. Implement backend Google session contract and auth middleware

  **What to do**:
  - Keep existing Google provider wiring as the auth origin.
  - Change auth flow from “callback returns JSON only” to a browser-friendly session contract:
    - `GET /api/v1/auth/google?return_to=<frontend-path>` starts OAuth and stores `return_to`
    - `GET /api/v1/auth/google/callback` upserts the user, sets HTTP-only cookie `playground_session`, and redirects to the saved `return_to`
    - `GET /api/v1/auth/session` returns `{ authenticated, user, bookingReady }`
    - `POST /api/v1/auth/logout` clears `playground_session`
  - Add non-production-only `POST /api/v1/auth/dev-login` guarded by `X-Dev-Auth-Key` so automated QA can seed the same cookie shape without real Google credentials.
  - Add backend auth middleware that resolves the session user and injects authenticated user identity into request context.

  **Must NOT do**:
  - Do not introduce JWT, local password auth, or a second session mechanism.
  - Do not require login before browsing `/lounge` or `/event`.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: cross-cutting auth/session work across handlers, middleware, and redirect flow
  - Skills: `[]` - no extra skills needed
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 3, 4, 6, 7, 8, 9 | Blocked By: none

  **References**:
  - Pattern: `backend/internal/http/handlers/auth.go:24-99` - existing Google OAuth begin/callback flow to extend instead of replacing
  - Pattern: `backend/internal/app/app.go:51-61` - current Goth Google provider registration
  - Pattern: `backend/internal/http/server/server.go:48-77` - route registration area for auth and booking endpoints
  - Pattern: `/home/destruct/Projects/Web/arh/frontend/app/stores/auth.ts:63-80` - example of frontend-side authenticated user bootstrap against cookie-backed API

  **Acceptance Criteria**:
  - [ ] `GET /api/v1/auth/session` returns `401` or `{ authenticated: false }` when no session cookie exists.
  - [ ] Successful callback or `dev-login` sets `playground_session` and redirects/returns a session that `GET /api/v1/auth/session` recognizes.
  - [ ] `POST /api/v1/auth/logout` clears the cookie and subsequent `GET /api/v1/auth/session` is unauthenticated.
  - [ ] `return_to` survives the auth round-trip.

  **QA Scenarios**:
  ```text
  Scenario: Real OAuth initiation preserves return path
    Tool: Playwright
    Steps:
      1. Open http://127.0.0.1:3000/lounge
      2. Trigger the primary booking action while signed out
      3. Assert browser navigates to /api/v1/auth/google with a return_to value for /lounge
    Expected: Auth starts from the correct backend route and includes the intended frontend return path
    Evidence: .sisyphus/evidence/task-1-auth-begin.png

  Scenario: Dev login seeds a valid session cookie
    Tool: Bash
    Steps:
      1. POST /api/v1/auth/dev-login with X-Dev-Auth-Key in non-production
      2. Reuse the returned cookie jar to call GET /api/v1/auth/session
    Expected: Session endpoint returns authenticated user payload and bookingReady boolean
    Evidence: .sisyphus/evidence/task-1-dev-login.txt
  ```

  **Commit**: NO | Message: `feat(auth): add browser session contract for google login` | Files: `backend/internal/http/handlers/auth.go`, `backend/internal/http/server/server.go`, `backend/internal/http/middleware/*`, related session files

- [ ] 2. Harden backend user/profile and booking data model

  **What to do**:
  - Add a migration that enforces unique Google identity with a unique index on `users.google_id` where non-null.
  - Keep `users` as the canonical latest-profile record.
  - Extend `bookings` with immutable snapshot columns:
    - `contact_name TEXT NOT NULL`
    - `contact_email VARCHAR NOT NULL`
    - `contact_phone VARCHAR NOT NULL`
    - `details_json JSONB NOT NULL DEFAULT '{}'::jsonb`
  - Add authenticated endpoints:
    - `GET /api/v1/me` → current user profile + `bookingReady`
    - `PATCH /api/v1/me/contact` → update `full_name` and/or `phone`
  - Normalize phone on the backend by stripping spaces, dashes, and parentheses; accept only 10-15 digits after normalization, then store with a leading `+`.
  - Update booking create contract to stop accepting arbitrary client `user_id`; derive the user from session context instead.
  - On booking create, atomically update missing/edited canonical `full_name`/`phone` and insert the booking snapshot.
  - Add duplicate-submit protection: reject a second active booking for the same `user_id + service_id + start_time` with `409`.

  **Must NOT do**:
  - Do not create separate booking tables for lounge and event.
  - Do not mutate old bookings when a user later changes profile data.

  **Recommended Agent Profile**:
  - Category: `unspecified-high` - Reason: schema, repository, validation, and handler contract updates
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 6, 7, 8, 9 | Blocked By: none

  **References**:
  - API/Type: `backend/sql/schema.sql:12-70` - current `users` and `bookings` tables to extend, not replace
  - API/Type: `backend/internal/domain/user.go:3-8` - current upsert request omits phone and must be expanded by explicit profile/update contract
  - API/Type: `backend/internal/domain/booking.go:7-25` - current booking request trusts `user_id`; must be session-derived
  - Pattern: `backend/internal/http/server/server.go:48-77` - placement for `me` endpoints alongside existing auth/booking routes

  **Acceptance Criteria**:
  - [ ] DB migration adds unique Google identity protection and booking snapshot columns.
  - [ ] `PATCH /api/v1/me/contact` updates canonical `users.full_name` / `users.phone` and returns normalized values.
  - [ ] Booking create rejects missing/invalid phone when the canonical user record lacks a valid phone.
  - [ ] Booking create stores immutable snapshot contact fields and optional per-flow metadata in `details_json`.
  - [ ] Duplicate active submits return `409`.

  **QA Scenarios**:
  ```text
  Scenario: First booking updates canonical phone and snapshot fields atomically
    Tool: Bash
    Steps:
      1. Seed authenticated session for a user without phone
      2. Call PATCH /api/v1/me/contact with full_name and phone
      3. Call authenticated booking create for a lounge payload
      4. Inspect API response and DB-visible fields through existing read endpoint or debug response
    Expected: users.phone is normalized, booking includes contact_name/contact_email/contact_phone snapshot, and no partial success occurs
    Evidence: .sisyphus/evidence/task-2-contact-and-booking.txt

  Scenario: Duplicate booking is rejected
    Tool: Bash
    Steps:
      1. Submit the same authenticated booking twice with identical service_id and start_time
      2. Inspect second response
    Expected: First request succeeds; second returns 409 duplicate/active-booking error
    Evidence: .sisyphus/evidence/task-2-duplicate-booking.txt
  ```

  **Commit**: NO | Message: `feat(backend): persist booking snapshots and profile readiness` | Files: `backend/sql/schema.sql`, migrations, domain/service/repo/handler files

- [ ] 3. Build frontend auth session and booking-intent state

  **What to do**:
  - Add frontend composables instead of introducing Pinia into PlayGround:
    - `useAuthSession()` for session fetch, logout, auth state, and booking readiness
    - `useBookingIntent()` for serializing pending lounge/event form state into `sessionStorage`
  - Fetch `GET /api/v1/auth/session` once on app bootstrap and whenever login/logout transitions occur.
  - Persist booking intent immediately before redirecting to Google auth.
  - After auth return, restore the pending intent into the original route and reopen the same flow step.
  - Store only same-tab temporary intent; clear it after successful booking or explicit reset.

  **Must NOT do**:
  - Do not add a new global state library.
  - Do not store sensitive auth tokens in localStorage.

  **Recommended Agent Profile**:
  - Category: `quick` - Reason: bounded composable/state plumbing in the Nuxt frontend
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 1 | Blocks: 4, 6, 7, 8, 9 | Blocked By: 1

  **References**:
  - Pattern: `frontend/app/pages/lounge.vue:67-258` - current booking form state that must become serializable/restorable
  - Pattern: `frontend/app/pages/event.vue:68-244` - current event state and validation structure to preserve
  - Pattern: `/home/destruct/Projects/Web/arh/frontend/app/stores/auth.ts:10-89` - example of Nuxt-side auth bootstrap against cookie-based endpoints
  - Test: `frontend/tests/smoke.spec.ts:16-47` - selectors that should stay stable where feasible

  **Acceptance Criteria**:
  - [ ] Signed-in state is available to header and booking pages via a shared composable.
  - [ ] Pending lounge/event intent is restored after auth return without losing selected zone/event and entered fields.
  - [ ] Successful booking or manual reset clears stored intent.

  **QA Scenarios**:
  ```text
  Scenario: Lounge intent survives auth return
    Tool: Playwright
    Steps:
      1. Open /lounge signed out
      2. Select a zone, party size, slot, and guest name
      3. Trigger auth-required submit, complete dev-login helper flow, return to /lounge
      4. Assert previously selected values are restored
    Expected: The same booking context is visible after auth return
    Evidence: .sisyphus/evidence/task-3-lounge-intent.png

  Scenario: Event intent clears after success
    Tool: Playwright
    Steps:
      1. Complete an authenticated event booking
      2. Reload /event
      3. Assert no stale booking intent is auto-restored
    Expected: Temporary intent storage is removed after success
    Evidence: .sisyphus/evidence/task-3-intent-clear.png
  ```

  **Commit**: NO | Message: `feat(frontend): add auth session and booking intent composables` | Files: `frontend/app/composables/*`, app bootstrap files

- [ ] 4. Rewrite login/register and header into Google-only auth surfaces

  **What to do**:
  - Keep `/login` and `/register` as public routes, but convert both to a single Google-only CTA surface.
  - `/login` copy = “Continue with Google”; `/register` copy = “Create account with Google”. Both call the same backend auth begin endpoint.
  - Remove username/password/nickname/email form fields from both pages.
  - Update header actions:
    - signed out → single primary Google CTA
    - signed in → avatar/name summary + logout action
  - Keep the existing route names so navigation doesn’t break.

  **Must NOT do**:
  - Do not leave dead email/password forms in place.
  - Do not add a separate profile page in this iteration.

  **Recommended Agent Profile**:
  - Category: `quick` - Reason: focused surface rewrite on bounded files
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 6, 7 | Blocked By: 1, 3

  **References**:
  - Pattern: `frontend/app/components/AppHeader.vue:15-45` - existing nav/action area to make auth-aware
  - Pattern: `frontend/app/pages/login.vue:1-89` - current placeholder login page to replace with Google-only entry
  - Pattern: `frontend/app/pages/register.vue:1-92` - current placeholder register page to replace with Google-only entry
  - Pattern: `/home/destruct/Projects/Web/arh/frontend/app/layouts/default.vue:31-43` - example of auth-aware top-level navigation behavior

  **Acceptance Criteria**:
  - [ ] `/login` and `/register` contain one Google CTA each and no local credential fields.
  - [ ] Header reflects signed-out vs signed-in state without page reload hacks.
  - [ ] Logout clears session and returns header to signed-out state.

  **QA Scenarios**:
  ```text
  Scenario: Login page starts Google auth
    Tool: Playwright
    Steps:
      1. Open /login
      2. Click the primary Google CTA
      3. Assert navigation reaches /api/v1/auth/google
    Expected: Login page uses Google-only auth start
    Evidence: .sisyphus/evidence/task-4-login-google.png

  Scenario: Header updates after logout
    Tool: Playwright
    Steps:
      1. Seed an authenticated session
      2. Open /
      3. Click logout from the header
      4. Assert header shows the signed-out Google CTA again
    Expected: Header state tracks session correctly
    Evidence: .sisyphus/evidence/task-4-header-logout.png
  ```

  **Commit**: NO | Message: `feat(frontend): replace placeholder auth screens with google entry` | Files: `frontend/app/pages/login.vue`, `frontend/app/pages/register.vue`, `frontend/app/components/AppHeader.vue`

- [ ] 5. Establish shared desktop UI foundation with Nuxt UI and reduced radii

  **What to do**:
  - Refactor shared lounge/event/auth shell pieces to use Nuxt UI primitives where appropriate:
    - buttons → `UButton`
    - cards/panels → `UCard`
    - alerts/status blocks → `UAlert`
    - form rows → `UForm` + `UFormField` + `UInput`
    - badges/status labels → `UBadge`
  - Keep the current dark/cyan premium visual direction, but standardize radii:
    - interactive controls max radius = `rounded-xl`
    - cards/panels max radius = `rounded-2xl`
    - remove oversized circular/`rounded-full` usage except for true pills/badges
  - Centralize these style constraints in shared components and theme/global CSS, not per-page one-offs.

  **Must NOT do**:
  - Do not redesign unrelated pages.
  - Do not preserve the existing large rounded login/register shells.

  **Recommended Agent Profile**:
  - Category: `visual-engineering` - Reason: UI-system cleanup and component-driven desktop polish
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 6, 7, 8 | Blocked By: none

  **References**:
  - Pattern: `frontend/nuxt.config.ts:2-27` - existing Nuxt UI module setup
  - Pattern: `frontend/app/app.config.ts:1-8` - existing UI token entry point
  - Pattern: `frontend/app/components/ExperienceCard.vue:1-91` - current shared card primitive with custom styling to simplify
  - Pattern: `frontend/app/components/ExperienceFlowPanel.vue:1-27` - current shared flow wrapper to translate into cleaner Nuxt UI usage

  **Acceptance Criteria**:
  - [ ] Lounge/event/auth surfaces rely on shared Nuxt UI primitives rather than duplicated raw inputs/buttons/cards.
  - [ ] Large decorative rounding is removed from desktop forms/cards except where explicitly kept as pills/badges.
  - [ ] Shared components encapsulate the new desktop styling rules.

  **QA Scenarios**:
  ```text
  Scenario: Desktop lounge surface uses shared UI primitives consistently
    Tool: Playwright
    Steps:
      1. Open /lounge at desktop width
      2. Capture screenshot of catalog cards and booking panel
      3. Assert primary CTA, fields, badges, and alerts are visually aligned and rendered without overflow
    Expected: The page renders a consistent component system with tighter radii and no broken layout
    Evidence: .sisyphus/evidence/task-5-lounge-desktop.png

  Scenario: Auth surface no longer uses oversized decorative shell
    Tool: Playwright
    Steps:
      1. Open /login at desktop width
      2. Capture screenshot
      3. Compare visually against acceptance rule: no giant rounded neon wrapper, clean card-based layout only
    Expected: Login page matches the bounded desktop refresh rules
    Evidence: .sisyphus/evidence/task-5-login-desktop.png
  ```

  **Commit**: NO | Message: `refactor(ui): rebuild desktop booking surfaces with nuxt ui` | Files: shared frontend components, `app.config.ts`, `main.css`, related page templates

- [ ] 6. Convert the lounge page into a real authenticated booking flow

  **What to do**:
  - Keep lounge browsing public.
  - On primary submit while signed out, persist booking intent and start Google auth.
  - On return with authenticated user:
    - show inline phone capture step only if `bookingReady === false`
    - prefill `full_name`, `email` (read-only), and `phone`
    - allow editing `full_name` and `phone` before final submit
  - Replace the generic `contactHandle` field with phone-driven booking data.
  - Map lounge payload as:
    - `zone_id`, `service_id`, `start_time`, `end_time`, `participants`
    - `details_json.type = "lounge"`
    - `details_json.slot_id = <selected slot>`
  - On success, show the server-confirmed booking state instead of the current mock-only success banner.

  **Must NOT do**:
  - Do not require phone before the user attempts to book.
  - Do not leave the old Telegram/free-form contact field in place.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: combines UI, validation, auth gating, and backend payload mapping
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 2 | Blocks: 8, 9 | Blocked By: 1, 2, 3, 4, 5

  **References**:
  - Pattern: `frontend/app/pages/lounge.vue:67-258` - current lounge catalog/form structure and stable selectors to evolve
  - Test: `frontend/tests/smoke.spec.ts:16-31` - existing lounge validation/happy-path coverage to preserve minimally
  - API/Type: `backend/internal/domain/booking.go:7-25` - booking create contract baseline that must be session-driven
  - API/Type: `backend/sql/schema.sql:58-70` - booking record baseline to extend with snapshot data

  **Acceptance Criteria**:
  - [ ] Signed-out users can browse lounge cards but cannot complete booking without Google auth.
  - [ ] First-time authenticated user without phone is blocked at inline phone step until valid phone is provided.
  - [ ] Returning user sees DB-prefilled name/email/phone, with email read-only and name/phone editable.
  - [ ] Successful submit creates an authenticated booking with `details_json.type = "lounge"` and stored contact snapshot.

  **QA Scenarios**:
  ```text
  Scenario: First lounge booking requires auth and phone
    Tool: Playwright
    Steps:
      1. Open /lounge signed out
      2. Select data-testid values for zone, party, and slot
      3. Click data-testid="lounge-submit"
      4. Complete dev-login helper return
      5. Fill required phone field and submit again
    Expected: Booking succeeds only after auth + valid phone; success state reflects the chosen lounge context
    Evidence: .sisyphus/evidence/task-6-lounge-first-booking.png

  Scenario: Invalid phone blocks lounge booking without losing selections
    Tool: Playwright
    Steps:
      1. Return from auth with missing phone
      2. Enter an invalid short phone value
      3. Submit the lounge form
    Expected: Explicit validation error is shown, selected zone/party/slot remain intact, and no booking is created
    Evidence: .sisyphus/evidence/task-6-lounge-invalid-phone.png
  ```

  **Commit**: NO | Message: `feat(lounge): wire google-authenticated booking flow` | Files: `frontend/app/pages/lounge.vue`, related composables/API helpers, backend booking/profile endpoints as needed

- [ ] 7. Convert the event page into the same authenticated account backbone

  **What to do**:
  - Preserve public event browsing and sold-out display.
  - Reuse the same auth gate and `bookingReady` rules as lounge.
  - Replace free email entry with Google-derived read-only email.
  - Keep attendee name editable and sync edits back to canonical `users.full_name` on successful submit.
  - Use the same phone readiness rule as lounge.
  - Map event payload as:
    - shared booking fields
    - `details_json.type = "event"`
    - `details_json.attendance_mode = <selected mode>`
    - `details_json.event_id = <selected event id>`
  - Keep sold-out items unbookable even when authenticated.

  **Must NOT do**:
  - Do not create a second auth/account flow just for events.
  - Do not leave manual free-form attendee email editable.

  **Recommended Agent Profile**:
  - Category: `deep` - Reason: shared auth contract plus event-specific payload/validation mapping
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: 8, 9 | Blocked By: 1, 2, 3, 4, 5

  **References**:
  - Pattern: `frontend/app/pages/event.vue:68-244` - current event listing/registration structure and selectors
  - Test: `frontend/tests/smoke.spec.ts:33-47` - current event happy-path/validation coverage to preserve minimally
  - API/Type: `backend/sql/schema.sql:33-70` - existing `zones` and `bookings` backbone shared with lounge
  - Pattern: `backend/internal/http/handlers/auth.go:72-99` - current Google-derived identity fields available after callback

  **Acceptance Criteria**:
  - [ ] Signed-out users can browse events but must authenticate before completing registration.
  - [ ] Event booking reuses the same authenticated account and phone readiness model as lounge.
  - [ ] Google email is shown but not editable; attendee name is editable and persisted canonically on success.
  - [ ] Sold-out events stay non-interactive.

  **QA Scenarios**:
  ```text
  Scenario: Returning user skips phone step on event booking
    Tool: Playwright
    Steps:
      1. Seed authenticated session for a user with stored phone
      2. Open /event
      3. Select an available event and attendance mode
      4. Submit booking
    Expected: No phone gate appears; event booking succeeds using the same account session
    Evidence: .sisyphus/evidence/task-7-event-returning-user.png

  Scenario: Sold-out event remains blocked
    Tool: Playwright
    Steps:
      1. Open /event
      2. Locate a sold-out event card
      3. Attempt to interact with its CTA
    Expected: CTA remains disabled and no registration flow starts for that item
    Evidence: .sisyphus/evidence/task-7-event-sold-out.png
  ```

  **Commit**: NO | Message: `feat(event): reuse unified authenticated booking flow` | Files: `frontend/app/pages/event.vue`, related composables/API helpers, backend booking/profile endpoints as needed

- [ ] 8. Harden cross-flow recovery and explicit error states

  **What to do**:
  - Add one shared error/status strategy across lounge and event for:
    - OAuth callback failure
    - expired/missing session on submit
    - invalid phone
    - duplicate booking (`409`)
    - stale/unavailable zone or event
  - Add stable `data-testid` hooks for these banners and recovery actions.
  - On session expiry during booking, preserve intent, restart auth, and return to the same flow.
  - On callback failure, show a dismissible error banner with a retry CTA that restarts Google auth.
  - Keep all failure states non-destructive to current selections unless the selected entity is no longer valid.

  **Must NOT do**:
  - Do not silently drop user selections on recoverable errors.
  - Do not show generic unknown-error text when a precise state is available.

  **Recommended Agent Profile**:
  - Category: `unspecified-high` - Reason: shared resilience layer across routes and API error contracts
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: YES | Wave 3 | Blocks: 9 | Blocked By: 1, 2, 3, 5, 6, 7

  **References**:
  - Pattern: `frontend/app/pages/lounge.vue:226-258` - current validation/error/success output area
  - Pattern: `frontend/app/pages/event.vue:202-243` - current validation/error/success output area
  - Pattern: `backend/internal/http/handlers/auth.go:58-69` - callback failure path that must map to explicit UI recovery
  - Test: `frontend/tests/smoke.spec.ts:16-47` - stable selector pattern to preserve while adding new error states

  **Acceptance Criteria**:
  - [ ] Callback failure, session expiry, invalid phone, duplicate booking, and stale availability each render distinct UI states.
  - [ ] Recoverable failures preserve existing selections and input.
  - [ ] Retry actions point to concrete next steps (restart auth, correct phone, reselect item).

  **QA Scenarios**:
  ```text
  Scenario: Callback failure shows retryable error state
    Tool: Playwright
    Steps:
      1. Force callback error mode in non-production QA environment
      2. Return to /lounge or /event
      3. Observe the auth-failure banner and click retry
    Expected: A specific auth failure message is shown and retry restarts the Google flow
    Evidence: .sisyphus/evidence/task-8-auth-failure.png

  Scenario: Duplicate submit surfaces a precise conflict error
    Tool: Playwright
    Steps:
      1. Complete one booking successfully
      2. Re-submit the same payload immediately
      3. Observe UI response
    Expected: UI surfaces duplicate-booking conflict, preserves context, and avoids false success
    Evidence: .sisyphus/evidence/task-8-duplicate-submit.png
  ```

  **Commit**: NO | Message: `fix(booking): add deterministic recovery states across auth and booking flows` | Files: lounge/event/shared UI/composables and related API error mapping

- [ ] 9. Maintain the existing smoke path and run final project gates

  **What to do**:
  - Update `frontend/tests/smoke.spec.ts` only as needed to reflect the new auth-aware flow, selectors, and success messages.
  - Do **not** add another Playwright spec file or second test suite.
  - Use the non-production dev-login helper for post-auth state setup inside smoke tests.
  - Keep route coverage for `/`, `/lounge`, and `/event`.
  - Ensure backend compile/test and frontend lint/typecheck/smoke all pass after the new contracts land.

  **Must NOT do**:
  - Do not introduce Cypress, Vitest, or a new Playwright project.
  - Do not expand smoke into a broad regression suite.

  **Recommended Agent Profile**:
  - Category: `quick` - Reason: bounded maintenance of an existing spec and verification commands
  - Skills: `[]`
  - Omitted: `[]`

  **Parallelization**: Can Parallel: NO | Wave 3 | Blocks: F1-F4 | Blocked By: 1, 2, 3, 6, 7, 8

  **References**:
  - Test: `frontend/tests/smoke.spec.ts:1-47` - only existing smoke file; update in place rather than adding a suite
  - Test: `frontend/playwright.config.ts:3-23` - current smoke runner config and dev server usage
  - Pattern: `frontend/package.json:5-13` - current lint/typecheck/smoke commands to preserve

  **Acceptance Criteria**:
  - [ ] Existing smoke file is updated in place and still covers homepage route entry plus lounge/event happy paths.
  - [ ] `corepack pnpm run lint`, `corepack pnpm run typecheck`, and `corepack pnpm run smoke` pass.
  - [ ] `go test ./...` passes.

  **QA Scenarios**:
  ```text
  Scenario: Full frontend smoke passes on existing Playwright entrypoint
    Tool: Bash
    Steps:
      1. Run corepack pnpm run smoke in /home/destruct/Projects/Web/PlayGround/frontend
    Expected: Existing smoke suite passes without adding new test files
    Evidence: .sisyphus/evidence/task-9-smoke.txt

  Scenario: Backend compile/test passes after auth and booking changes
    Tool: Bash
    Steps:
      1. Run go test ./... in /home/destruct/Projects/Web/PlayGround/backend
    Expected: Backend packages compile and tests pass
    Evidence: .sisyphus/evidence/task-9-go-test.txt
  ```

  **Commit**: NO | Message: `test(smoke): align existing smoke flow with google-auth bookings` | Files: `frontend/tests/smoke.spec.ts`, possibly tiny test helpers only if kept in the same existing test area

## Final Verification Wave (MANDATORY — after ALL implementation tasks)
> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.
> Do NOT auto-proceed after verification. Wait for user's explicit approval before marking work complete.
> Never mark F1-F4 as checked before getting user's okay.

- [ ] F1. Plan Compliance Audit — oracle
- [ ] F2. Code Quality Review — unspecified-high
- [ ] F3. Real Manual QA — unspecified-high (+ playwright for UI)
- [ ] F4. Scope Fidelity Check — deep

## Commit Strategy
- Do not create any git commits automatically during execution unless the user explicitly requests commits.
- If the user later requests commits, prefer one commit after Wave 1, one after Wave 2, and one after Wave 3 using the commit messages drafted per task grouping.

## Success Criteria
- Google account is the single required booking identity for both lounge and event flows.
- Phone is requested only at first booking if missing, not at initial sign-in.
- Later bookings are prefilled from DB and remain editable inline before submit.
- Existing bookings retain their original snapshot contact data even after later profile edits.
- Desktop UI uses shared Nuxt UI components with visibly tighter radii and cleaner layout structure.
- Existing smoke infrastructure remains the only automated frontend suite used in this iteration.
