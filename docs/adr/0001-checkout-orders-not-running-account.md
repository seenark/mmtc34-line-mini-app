# Checkout Orders instead of a running account per Campaign

The original handoff modelled one Participation per Member per Campaign: a running account where Items are added at any time and paid off in parts with multiple Slips. We replaced it with ordinary shop-style Orders: a Member builds a cart, checks out, and pays the full Order with one Slip. A Member may place several Orders in one Campaign. Members already know this flow from online shopping, and one Order with one Payment reduces Slip review to "amount matches or not". Partial or mismatched amounts become Admin Adjustments.

## Consequences

- Address and Delivery Method are captured per Order, not per Member per Campaign.
- Members cannot change what an Order costs after checkout (Options, quantities, Delivery Method). While the Campaign is open they may correct an Item's details (size, print name, number) and the shipping name, address and phone, even once paid, so a typo never needs an Admin. Admins edit anything at any time.
- There is no payment deadline per Order. Closing the Campaign is the only cut-off: unpaid Orders are cancelled then, so nothing lingers after close.
- Exception: an Order an Admin places on a Member's behalf is never cancelled by the close. Admins take those requests from the group chat after the cut-off, and cancel them explicitly if needed.
