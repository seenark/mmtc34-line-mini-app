# Anyone who signs in with LINE becomes a Member without approval

Milestone 1 made every newcomer an Applicant who saw nothing until an Admin approved them. The owner dropped the gate: approving classmates one by one was busywork, and the group only shares the LIFF link inside its own LINE group. Signing in with LINE and filling in a profile now makes a Member at once. A reader will find a money-collecting app that lets any LINE user in. The owner accepted that trade-off on purpose: a forwarded link lets a stranger see open Campaigns, their bank details and the "who ordered" roster (nicknames and sizes), and place Orders, but never other Members' Orders, real names or contact data.

## Consequences

- Applicants that exist at deploy time become Members.
- Linking can no longer copy a LINE account onto an Unlinked Member, as ADR 0002 described, because the person is already a Member when an Admin notices them. An Admin instead links the Unlinked Member's Orders, from any one of them, to that Member, and the emptied Unlinked Member is removed.
- ADR 0005 still holds for promoting Admins. The Super Admin no longer approves anyone.
