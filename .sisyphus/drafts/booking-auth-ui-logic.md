# Draft: Booking auth UI logic

## Requirements (confirmed)
- User wants a UI plan for booking logic centered around Google sign-in.
- Google sign-in is used for user identification across devices.
- After sign-in, user should be able to book a resource: gaming PC, lounge, or event.
- For lounge and event, booking should capture only a numeric participant count; per-person names are not needed.
- For gaming PC, booking is always for a single person.
- Booking should also capture a separate master contact phone for the reservation.
- Google sign-in level for this iteration: mock Google UI.
- Master phone behavior: prefill from user profile after sign-in, but editable per booking.

## Technical Decisions
- Planning mode only: no implementation in this step.
- Frontend logic should be planned around current mock-driven UI unless repo evidence requires otherwise.
- Booking owner should be treated as an authenticated user, not an anonymous guest.
- Lounge/event booking model should use participant count only, plus one master booking owner/contact.

## Research Findings
- Current auth UI is in `frontend/app/pages/login.vue` and `frontend/app/pages/register.vue`; both are static forms with no OAuth flow yet.
- Current booking UI is split across `frontend/app/pages/gaming.vue`, `frontend/app/pages/lounge.vue`, and `frontend/app/pages/event.vue`.
- Current lounge/event forms store only one person name plus one contact/email; gaming has no person/contact model yet.
- Shared booking mocks/state live in `frontend/app/utils/experienceData.ts`, `frontend/app/composables/useExperienceMockState.ts`, and `frontend/app/composables/useMockSubmission.ts`.
- Docs require OAuth login (`docs/sixth_stage_interfaces.md`), user identity via `google_id` plus `phone` (`docs/fifth_stage_database.md`), and bookings owned by `user_id` with `participants` count.
- Docs do not define an anonymous "guest booking" concept; the booking model is centered on an authenticated client/user.

## Open Questions
- Does this iteration include only booking creation, or also a "my bookings" screen after sign-in?

## Scope Boundaries
- INCLUDE: Google sign-in entry, booking creation flow, person-name collection, master contact phone, resource-specific booking UX.
- EXCLUDE: backend implementation until explicitly confirmed by requirements.
