# 01 - Add Participant Management

**What to build:** Owners can add one registered User to a Group by email, see active Participants in the Group, and receive clear errors for unknown or duplicate emails. Regular Participants can view the Participant list but cannot manage Participants.

**Blocked by:** None - can start immediately.

**Status:** resolved

- [x] An Owner can add one existing User to a Group by email.
- [x] Adding an unknown email returns a clear validation error.
- [x] Adding a User who is already a Participant returns a clear duplicate error.
- [x] Regular Participants cannot add Participants.
- [x] Users can see only Participants for Groups they participate in.
- [x] OpenAPI documents the Participant management routes.
- [x] Backend tests cover Owner authorization, Regular Participant rejection, unknown email, duplicate email, and successful add.
- [x] Frontend tests cover the add-by-email form success and error states.
