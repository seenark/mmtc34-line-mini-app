# A password-based Super Admin bootstraps Admins

Identity is the LINE userId, and the first plan was to list Admin userIds in an env var. The owner cannot know their LINE userId before their first login, so that bootstrap could never be configured. Instead, `SUPER_ADMIN_USERNAME`/`SUPER_ADMIN_PASSWORD` define one Super Admin that signs in with a password, approves Applicants and promotes Members to Admin. It is stored as an ordinary Member row with the reserved `line_user_id` `superadmin`, so sessions, Admin gates and "reviewed by" work without a second identity model. A reader will find a non-LINE account in a LINE-only app; that is deliberate.

## Consequences

- The app now has one password login on the internet. It is rate-limited (5 failures per minute) and the password must be at least 12 characters.
- Nobody can demote or delete the Super Admin. Rotating its password does not end sessions; rotating `SESSION_SECRET` does.
- The Super Admin appears in the Admin member list with a badge but is hidden from pickers meant for classmates.
