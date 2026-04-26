# Frontend Continuation Plan

## TL;DR

> **Quick Summary**: Continue the Nuxt frontend with two new public mock-driven routes — `/lounge` and `/event` — plus small homepage/navigation fixes so the new areas are reachable and consistent with the current UI.
>
> **Deliverables**:
> - Lounge page with catalog + booking flow on mock data
> - Event page with listing + registration flow on mock data
> - Shared UI/mocks extracted from current page-level code
> - Minimal smoke setup for route and flow verification
>
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 3 waves
> **Critical Path**: 1 → 4 → 7 → 8 → F1-F4

---

## Context

### Original Request
User asked to go into `/home/destruct/Projects/Web/PlayGround`, prepare a plan for continuing frontend development, and focus for now on interface-only work without backend.

### Interview Summary
**Key Discussions**:
- Priority is **main screens first**.
- Main targets are **lounge zone** and **event zone**.
- These must be **separate pages/routes**, not just blocks on the homepage.
- UI depth should be **near-finished mock UX** with validation and explicit success/error/loading/empty states.
- Work is **frontend-only**, no backend integration.
- Style should **continue the current UI language**.
- Responsive priority is **desktop-first**.
- Testing choice is **minimal smoke setup**, not a full browser-test platform rollout.

**Research Findings**:
- `frontend/app/pages/gaming.vue` is the richest current interaction pattern and the best source for local mock flow behavior.
- `frontend/app/pages/index.vue` already advertises Lounge/Event but the CTA actions are not wired into full route flows.
- `frontend/app/components/AppHeader.vue` contains `/about` and `/features` links that do not match the currently discovered pages.
- `docs/frontend.md` requires mock/static data, reusable components, Nuxt routing, and stylistic consistency.
- `docs/sixth_stage_interfaces.md` suggests broader client capabilities, so scope must be explicitly narrowed for this iteration.

### Metis Review
**Identified Gaps** (addressed in this plan):
- Locked route defaults to `/lounge` and `/event`.
- Locked flow mode to **public mock flows** with **reset-on-refresh** behavior.
- Locked scope to exclude profile, booking management, upcoming visits, and real auth/backend.
- Limited homepage changes to **CTA routing and consistency polish**, not redesign.
- Included explicit handling of dead header links as a bounded cleanup item.
- Converted “minimal smoke” into concrete smoke tasks and QA scenarios.

---

## Work Objectives

### Core Objective
Extend the current Nuxt UI with two production-shaped mock user journeys — lounge booking and event registration — while keeping the implementation static/local, visually consistent, and easy to integrate with backend later.

### Concrete Deliverables
- `frontend/app/pages/lounge.vue`
- `frontend/app/pages/event.vue`
- Shared mock/state/component files under `frontend/app/components` and/or `frontend/app/composables` as needed
- Homepage CTA wiring to the new routes
- Header/navigation cleanup required to avoid dead-end UX on new pages
- Minimal smoke verification setup and smoke specs/scripts

### Definition of Done
- [ ] `pnpm run lint` passes in `frontend/`
- [ ] `pnpm run typecheck` passes in `frontend/`
- [ ] Smoke verification for `/`, `/lounge`, and `/event` passes
- [ ] Lounge and event flows both demonstrate happy-path and failure-path states on mock data

### Must Have
- Public mock-accessible `/lounge` page with catalog and booking flow
- Public mock-accessible `/event` page with listing and registration flow
- Deterministic loading/empty/success/error states for both routes
- Shared-component extraction where repeated UI would otherwise be duplicated
- Minimal smoke harness for route and flow verification

### Must NOT Have (Guardrails)
- No backend/API integration
- No real authentication, OAuth, profile, booking management, or upcoming-visits screens
- No payment, notification, analytics, CMS/admin, or booking-engine generalization
- No full redesign of homepage/header/theme
- No persistence requirement beyond in-memory/local route session unless explicitly added later

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** — all verification must be agent-executable.

### Test Decision
- **Infrastructure exists**: Partial project quality checks exist (`lint`, `typecheck`), but no confirmed frontend browser smoke setup
- **Automated tests**: Tests-after (minimal smoke)
- **Framework**: Playwright smoke setup + existing `pnpm run lint` / `pnpm run typecheck`

### QA Policy
Every implementation task includes runnable QA scenarios with evidence under `.sisyphus/evidence/`.

- **Frontend/UI**: Playwright
- **Static verification**: `pnpm run lint`, `pnpm run typecheck`
- **Evidence**: screenshots, Playwright output, and terminal logs

---

## Execution Strategy

### Parallel Execution Waves

```text
Wave 1 (Start Immediately - routing + shared foundations):
├── Task 1: Route contract + homepage CTA wiring [quick]
├── Task 2: Shared mock datasets + state helpers [quick]
└── Task 3: Shared UI primitives for cards/forms/status blocks [quick]

Wave 2 (After Wave 1 - core route builds):
├── Task 4: Lounge page shell + catalog experience (depends: 1,2,3) [visual-engineering]
├── Task 5: Event page shell + listing experience (depends: 1,2,3) [visual-engineering]
└── Task 6: Global shell/header cleanup for new routes (depends: 1,3) [quick]

Wave 3 (After Wave 2 - flow completion + smoke):
├── Task 7: Lounge booking and event registration states/validation (depends: 4,5,2) [deep]
└── Task 8: Minimal smoke setup and smoke specs (depends: 4,5,6,7) [unspecified-high]

Wave FINAL (After ALL tasks):
├── Task F1: Plan compliance audit (oracle)
├── Task F2: Code quality review (unspecified-high)
├── Task F3: Real manual QA execution (unspecified-high)
└── Task F4: Scope fidelity check (deep)
```

### Dependency Matrix

- **1**: Blocked By: none | Blocks: 4, 5, 6
- **2**: Blocked By: none | Blocks: 4, 5, 7
- **3**: Blocked By: none | Blocks: 4, 5, 6
- **4**: Blocked By: 1, 2, 3 | Blocks: 7, 8
- **5**: Blocked By: 1, 2, 3 | Blocks: 7, 8
- **6**: Blocked By: 1, 3 | Blocks: 8
- **7**: Blocked By: 2, 4, 5 | Blocks: 8
- **8**: Blocked By: 4, 5, 6, 7 | Blocks: F1-F4

### Agent Dispatch Summary

- **Wave 1**: T1 `quick`, T2 `quick`, T3 `quick`
- **Wave 2**: T4 `visual-engineering`, T5 `visual-engineering`, T6 `quick`
- **Wave 3**: T7 `deep`, T8 `unspecified-high`
- **FINAL**: F1 `oracle`, F2 `unspecified-high`, F3 `unspecified-high`, F4 `deep`

---

## TODOs

- [ ] 1. Wire route contract and homepage entry points

  **What to do**:
  - Create/confirm route contract for `/lounge` and `/event`.
  - Update homepage CTA actions in `frontend/app/pages/index.vue` so Lounge/Event are real route entries.
  - Keep homepage changes intentionally narrow: CTA wiring, labels, and any minimal state/affordance polish needed to avoid dead ends.

  **Must NOT do**:
  - Do not redesign the homepage layout.
  - Do not add unrelated content sections.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: focused route wiring and bounded page-entry adjustments.
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - `playwright`: not needed for implementation; reserved for QA.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3)
  - **Blocks**: 4, 5, 6
  - **Blocked By**: None

  **References**:
  - `frontend/app/pages/index.vue:30-68` - Existing homepage cards for Gaming/Lounge/Event; use this as the exact entry-point area to convert Lounge/Event from dead-end actions into route links.
  - `frontend/app/app.vue:1-12` - Global shell behavior; confirms new routes will render inside the shared app layout.
  - `docs/frontend.md:21-25` - Routing must reflect intended page structure using Nuxt file routing.

  **Acceptance Criteria**:
  - [ ] Homepage Lounge CTA routes to `/lounge`.
  - [ ] Homepage Event CTA routes to `/event`.
  - [ ] Direct browser navigation to `/lounge` and `/event` resolves to real pages once implemented.

  **QA Scenarios**:
  ```text
  Scenario: Homepage lounge CTA navigation
    Tool: Playwright
    Preconditions: Nuxt dev server running in frontend/
    Steps:
      1. Open http://127.0.0.1:3000/
      2. Click the Lounge CTA button/link in the second card
      3. Assert URL is /lounge
    Expected Result: Browser reaches the lounge route without console/runtime errors
    Failure Indicators: CTA is still a plain button, route stays on /, 404 page appears
    Evidence: .sisyphus/evidence/task-1-home-to-lounge.png

  Scenario: Homepage event CTA navigation
    Tool: Playwright
    Preconditions: Nuxt dev server running in frontend/
    Steps:
      1. Open http://127.0.0.1:3000/
      2. Click the Event CTA button/link in the third card
      3. Assert URL is /event
    Expected Result: Browser reaches the event route without console/runtime errors
    Evidence: .sisyphus/evidence/task-1-home-to-event.png
  ```

  **Commit**: NO

- [ ] 2. Extract shared mock datasets and flow helpers

  **What to do**:
  - Move lounge/event mock content, capacities, and deterministic UI-state fixtures out of page-local clutter into reusable frontend files.
  - Include explicit fixtures/flags for loading, empty, success, and error states.
  - Keep persistence local/in-memory; reset on refresh is acceptable for this iteration.

  **Must NOT do**:
  - Do not create API clients or fake backend layers.
  - Do not over-abstract into a generic booking engine.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: bounded data modeling and helper extraction.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3)
  - **Blocks**: 4, 5, 7
  - **Blocked By**: None

  **References**:
  - `frontend/app/pages/gaming.vue:235-339` - Current pattern for local mock data, selections, capacity-like taken slots, and UI helper functions.
  - `docs/frontend.md:5-13` - This stage explicitly uses static data and minimal logic.
  - `docs/frontend.md:39-45` - Encourages component/data organization that reduces duplication and supports later backend integration.

  **Acceptance Criteria**:
  - [ ] Lounge and event mock data live outside page templates.
  - [ ] At least one deterministic fixture exists for each: loading, empty, success, error.
  - [ ] Remaining-capacity/sold-out logic is represented in mock form without backend dependencies.

  **QA Scenarios**:
  ```text
  Scenario: Mock empty state can be triggered deterministically
    Tool: Playwright
    Preconditions: App exposes a deterministic mock fixture or query-driven mode for empty state
    Steps:
      1. Open /lounge with the empty-state trigger enabled
      2. Assert catalog empty-state block is visible
      3. Open /event with the empty-state trigger enabled
      4. Assert event empty-state block is visible
    Expected Result: Both pages render explicit empty states without crashes
    Evidence: .sisyphus/evidence/task-2-empty-states.png

  Scenario: Capacity data marks unavailable options
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open /lounge
      2. Assert at least one place/slot has disabled styling/behavior
      3. Open /event
      4. Assert at least one sold-out/full event state is rendered if fixture provides it
    Expected Result: Unavailable options are visibly non-interactive
    Evidence: .sisyphus/evidence/task-2-capacity-states.png
  ```

  **Commit**: NO

- [ ] 3. Extract shared UI primitives for new routes

  **What to do**:
  - Extract repeated UI patterns needed across lounge/event flows: hero/header blocks, cards, status/empty panels, selection blocks, or action footers.
  - Keep styling aligned with existing cyan/dark premium aesthetic.
  - Prefer a small number of reusable pieces over page-specific duplication.

  **Must NOT do**:
  - Do not create a full design system refactor.
  - Do not migrate unrelated existing pages unless the shared primitive is actually reused.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: small shared UI extraction from known page patterns.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2)
  - **Blocks**: 4, 5, 6
  - **Blocked By**: None

  **References**:
  - `frontend/app/pages/gaming.vue:43-171` - Existing premium section framing, selection controls, and booking-panel layout patterns worth reusing.
  - `frontend/app/pages/index.vue:34-68` - Existing card style and CTA presentation to keep visual continuity.
  - `frontend/app/assets/css/main.css:1-18` - Theme/font baseline.
  - `frontend/app/app.config.ts:1-8` - App-level color direction.

  **Acceptance Criteria**:
  - [ ] Shared UI pieces used by both lounge and event pages are componentized.
  - [ ] New routes do not duplicate large structural blocks unnecessarily.
  - [ ] Shared primitives preserve the current visual language.

  **QA Scenarios**:
  ```text
  Scenario: Shared status component renders across both routes
    Tool: Playwright
    Preconditions: Default fixture plus state trigger available
    Steps:
      1. Open /lounge with loading fixture enabled
      2. Assert shared loading/status block is visible
      3. Open /event with error fixture enabled
      4. Assert shared error/status block is visible
    Expected Result: Common UI primitives render correctly on both routes
    Evidence: .sisyphus/evidence/task-3-shared-status-blocks.png

  Scenario: Shared action/footer blocks preserve styling
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open /lounge and capture CTA block
      2. Open /event and capture CTA block
      3. Compare for consistent typography, spacing, and theme treatment
    Expected Result: Both routes look like the same product family
    Evidence: .sisyphus/evidence/task-3-shared-visual-family.png
  ```

  **Commit**: NO

- [ ] 4. Build the lounge route shell and catalog experience

  **What to do**:
  - Implement `frontend/app/pages/lounge.vue`.
  - Render a lounge catalog with cards/sections for available areas/tables/zones.
  - Support a primary booking entry path from the page shell into the later booking flow.
  - Include deterministic loading and empty catalog states.

  **Must NOT do**:
  - Do not implement payment or real reservation persistence.
  - Do not add profile/my-bookings sections.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: this is a route-level UX build requiring layout polish and continuity with current style.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 5, 6)
  - **Blocks**: 7, 8
  - **Blocked By**: 1, 2, 3

  **References**:
  - `frontend/app/pages/gaming.vue:1-171` - Main example for immersive route shell, section sequencing, and selectable booking-related UI.
  - `frontend/app/pages/index.vue:46-56` - Existing Lounge positioning/copy direction from homepage.
  - `docs/frontend.md:8-13` - Single visual style, simplicity, reusable components, responsive layout.
  - `docs/sixth_stage_interfaces.md:8-13` - Client-facing capability includes viewing zones/services and creating bookings.

  **Acceptance Criteria**:
  - [ ] `/lounge` renders a real page with route-specific title and shell content.
  - [ ] Lounge catalog renders from extracted mock data.
  - [ ] Empty and loading states are triggerable and visually complete.
  - [ ] The page contains a clear CTA into lounge booking.

  **QA Scenarios**:
  ```text
  Scenario: Lounge page happy-path shell and catalog
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open http://127.0.0.1:3000/lounge
      2. Assert page heading includes Lounge-related copy
      3. Assert at least one catalog card/zone block is visible
      4. Assert booking CTA is visible and enabled
    Expected Result: Lounge route loads as a complete public page
    Evidence: .sisyphus/evidence/task-4-lounge-shell.png

  Scenario: Lounge empty state
    Tool: Playwright
    Preconditions: Empty fixture enabled
    Steps:
      1. Open /lounge with empty fixture
      2. Assert empty-state title/message and fallback CTA are visible
    Expected Result: No broken layout or blank page when catalog is empty
    Evidence: .sisyphus/evidence/task-4-lounge-empty.png
  ```

  **Commit**: NO

- [ ] 5. Build the event route shell and listing experience

  **What to do**:
  - Implement `frontend/app/pages/event.vue`.
  - Render event listing/showcase from mock data with date/slot/capacity metadata.
  - Support a primary registration entry path from listing to later registration flow.
  - Include deterministic loading and empty states.

  **Must NOT do**:
  - Do not add calendar sync, reminders, or admin tooling.
  - Do not broaden into a generic event platform.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: route-level UI and listing hierarchy with current brand styling.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4, 6)
  - **Blocks**: 7, 8
  - **Blocked By**: 1, 2, 3

  **References**:
  - `frontend/app/pages/index.vue:58-67` - Existing Event CTA/copy source on homepage.
  - `frontend/app/pages/gaming.vue:96-171` - Section composition and progressive reveal pattern for detailed selection areas.
  - `docs/frontend.md:21-25` - Page structure should align with Nuxt routing and planned interfaces.
  - `docs/sixth_stage_interfaces.md:8-13` - Supports viewing services/zones and creating bookings from client UI.

  **Acceptance Criteria**:
  - [ ] `/event` renders a real page with route-specific title and shell content.
  - [ ] Event list/showcase renders from extracted mock data.
  - [ ] Empty and loading states are triggerable and visually complete.
  - [ ] The page contains a clear CTA into event registration.

  **QA Scenarios**:
  ```text
  Scenario: Event page happy-path shell and listing
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open http://127.0.0.1:3000/event
      2. Assert page heading includes event-related copy
      3. Assert at least one event card/row is visible
      4. Assert registration CTA is visible and enabled for an available event
    Expected Result: Event route loads as a complete public page
    Evidence: .sisyphus/evidence/task-5-event-shell.png

  Scenario: Event empty state
    Tool: Playwright
    Preconditions: Empty fixture enabled
    Steps:
      1. Open /event with empty fixture
      2. Assert empty-state title/message and fallback CTA are visible
    Expected Result: No broken layout or blank page when there are no events
    Evidence: .sisyphus/evidence/task-5-event-empty.png
  ```

  **Commit**: NO

- [ ] 6. Clean up global shell/header behavior exposed by the new routes

  **What to do**:
  - Decide and implement the minimal header cleanup needed once `/lounge` and `/event` become visible public routes.
  - Fix, hide, or otherwise neutralize dead `/about` and `/features` links in a bounded way.
  - Ensure new pages fit cleanly into the existing shell.

  **Must NOT do**:
  - Do not redesign the global header.
  - Do not introduce a larger navigation IA project.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: small but necessary consistency cleanup.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4, 5)
  - **Blocks**: 8
  - **Blocked By**: 1, 3

  **References**:
  - `frontend/app/components/AppHeader.vue:13-24` - Current nav contains links that do not match discovered pages.
  - `frontend/app/app.vue:1-12` - New pages will inherit this shell and header behavior.
  - `docs/frontend.md:8-13` - Consistency and simplicity are explicit requirements.

  **Acceptance Criteria**:
  - [ ] New public routes do not expose obviously dead navigation destinations.
  - [ ] Header behavior remains consistent across `/`, `/lounge`, `/event`, and `/gaming`.
  - [ ] Login route exception remains intact if still required.

  **QA Scenarios**:
  ```text
  Scenario: Header remains consistent on new public routes
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open /
      2. Capture header
      3. Open /lounge and /event
      4. Assert header is present and visually consistent
    Expected Result: Shared shell feels coherent across public routes
    Evidence: .sisyphus/evidence/task-6-header-consistency.png

  Scenario: Dead navigation does not remain exposed
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open /
      2. Inspect header nav items
      3. Verify clicking any visible nav item does not lead to a known missing page
    Expected Result: No user-visible dead-end nav links remain
    Evidence: .sisyphus/evidence/task-6-no-dead-nav.png
  ```

  **Commit**: NO

- [ ] 7. Complete booking/registration flows with validation and state handling

  **What to do**:
  - Finish lounge booking flow: selection, required fields, capacity-aware disable/error, confirmation/success.
  - Finish event registration flow: selection, required fields, sold-out handling, confirmation/success.
  - Provide deterministic loading/error/success/reset behavior and prevent obvious double-submit bugs.

  **Must NOT do**:
  - Do not build real auth-gating.
  - Do not persist bookings across refresh unless trivial and explicitly local.

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: combined flow-state logic, validation, and edge-case handling across two routes.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: 8
  - **Blocked By**: 2, 4, 5

  **References**:
  - `frontend/app/pages/gaming.vue:116-339` - Existing booking-panel interaction model for date/zone/place/time selection and unavailable-slot handling.
  - `frontend/app/pages/login.vue:33-63` - Existing form styling direction for required input interaction.
  - `docs/frontend.md:5-13` - Minimal logic on static data, but enough UI behavior for clear UX.

  **Acceptance Criteria**:
  - [ ] Lounge flow blocks booking when required inputs are missing.
  - [ ] Lounge flow blocks over-capacity/unavailable selection with explicit feedback.
  - [ ] Event flow blocks registration when event is sold out or fields are invalid.
  - [ ] Both flows surface loading, success, and error states deterministically.
  - [ ] Double-submit is prevented while mock submission is in progress.

  **QA Scenarios**:
  ```text
  Scenario: Lounge booking happy path
    Tool: Playwright
    Preconditions: Default fixture enabled with available capacity
    Steps:
      1. Open /lounge
      2. Select an available zone/table/slot
      3. Fill required fields with concrete data (e.g. name 'Ivan', party size '2')
      4. Submit booking
      5. Assert success state/message is shown
    Expected Result: Successful mock booking confirmation appears
    Evidence: .sisyphus/evidence/task-7-lounge-success.png

  Scenario: Lounge booking failure path
    Tool: Playwright
    Preconditions: Default fixture enabled
    Steps:
      1. Open /lounge
      2. Try to submit without one required field or choose an unavailable slot
      3. Assert validation/error message is shown and submission does not complete
    Expected Result: Booking is blocked gracefully with explicit feedback
    Evidence: .sisyphus/evidence/task-7-lounge-error.png

  Scenario: Event registration happy path
    Tool: Playwright
    Preconditions: Default fixture enabled with at least one available event
    Steps:
      1. Open /event
      2. Choose an available event
      3. Fill required fields with concrete data (e.g. name 'Anna', email 'anna@example.com')
      4. Submit registration
      5. Assert success state/message is shown
    Expected Result: Successful mock registration confirmation appears
    Evidence: .sisyphus/evidence/task-7-event-success.png

  Scenario: Event registration failure path
    Tool: Playwright
    Preconditions: Default fixture enabled with sold-out or invalid case available
    Steps:
      1. Open /event
      2. Select a sold-out event or submit invalid/missing fields
      3. Assert error/validation message is shown and CTA stays blocked or fails gracefully
    Expected Result: Registration cannot proceed when constraints are violated
    Evidence: .sisyphus/evidence/task-7-event-error.png
  ```

  **Commit**: NO

- [ ] 8. Add minimal smoke setup and smoke specs

  **What to do**:
  - Add the smallest practical browser-smoke setup for this frontend iteration.
  - Cover route reachability and one happy/failure flow for lounge and event.
  - Wire commands/scripts so smoke execution is deterministic for later agents/CI expansion.

  **Must NOT do**:
  - Do not turn this into a full E2E platform migration.
  - Do not add broad cross-browser/device matrices.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: tool setup plus end-to-end verification coverage.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3
  - **Blocks**: F1-F4
  - **Blocked By**: 4, 5, 6, 7

  **References**:
  - `frontend/package.json:5-12` - Existing scripts for build/dev/lint/typecheck; extend without breaking current package workflow.
  - `frontend/.github/workflows/ci.yml:27-34` - Existing CI currently runs lint and typecheck only; smoke setup should be compatible with future CI extension but not require full CI redesign now.
  - `docs/frontend.md:52-53` - Stage result should be a fully assembled interface ready for later integration.

  **Acceptance Criteria**:
  - [ ] A concrete smoke command exists for frontend route/flow checks.
  - [ ] Smoke coverage includes `/`, `/lounge`, and `/event`.
  - [ ] Smoke coverage includes at least one happy path and one failure path.
  - [ ] Existing lint/typecheck workflow remains usable.

  **QA Scenarios**:
  ```text
  Scenario: Smoke suite passes on main routes
    Tool: Bash + Playwright
    Preconditions: Dependencies installed in frontend/
    Steps:
      1. Start dev server in frontend/
      2. Run lint and typecheck
      3. Run the smoke command
      4. Assert smoke suite passes for /, /lounge, /event
    Expected Result: Route smoke checks pass without regressions
    Evidence: .sisyphus/evidence/task-8-smoke-pass.txt

  Scenario: Failure-path smoke is enforced
    Tool: Playwright
    Preconditions: Default or error fixture enabled
    Steps:
      1. Run smoke spec covering invalid lounge/event submissions
      2. Assert the suite checks for visible validation/error messages instead of allowing submission
    Expected Result: Failure-path assertions are part of automated smoke coverage
    Evidence: .sisyphus/evidence/task-8-smoke-failure-path.txt
  ```

  **Commit**: YES
  - Message: `feat(frontend): add lounge and event mock flows`
  - Files: `frontend/app/**`, `frontend/package.json`, smoke setup files
  - Pre-commit: `pnpm run lint && pnpm run typecheck && <smoke-command>`

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Verify `/lounge`, `/event`, homepage CTA wiring, mock states, and smoke setup all match this plan. Reject any backend/auth/profile creep.

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `pnpm run lint` and `pnpm run typecheck`, review changed files for duplication, dead code, and over-abstraction.

- [ ] F3. **Real Manual QA** — `unspecified-high` (+ Playwright)
  Execute all QA scenarios above, capture evidence under `.sisyphus/evidence/final-qa/`, and verify both happy and failure paths.

- [ ] F4. **Scope Fidelity Check** — `deep`
  Compare final diff against the plan and ensure no implementation drift into backend, auth, booking-management, or redesign work.

---

## Commit Strategy

- **Single final commit**: `feat(frontend): add lounge and event mock flows`

---

## Success Criteria

### Verification Commands
```bash
pnpm run lint
pnpm run typecheck
<smoke-command>
```

### Final Checklist
- [ ] All "Must Have" items present
- [ ] All "Must NOT Have" items absent
- [ ] `/lounge` and `/event` are reachable from homepage
- [ ] Mock states are deterministic and testable
- [ ] Smoke verification passes
