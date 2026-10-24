# Checkout Orders instead of a running account per Campaign

The original handoff modelled one Participation per Member per Campaign: a running account where Items are added at any time and paid off in parts with multiple Slips. We replaced it with ordinary shop-style Orders: a Member builds a cart, checks out, and pays the full Order with one Slip. A Member may place several Orders in one Campaign. Members already know this flow from online shopping, and one Order with one Payment reduces Slip review to "amount matches or not". Partial or mismatched amounts become Admin Adjustments.

## Consequences

- Address and Delivery Method are captured per Order, not per Member per Campaign.
- Members cannot edit an Order after checkout; Admins edit anything at any time.
- There is no payment deadline per Order. Closing the Campaign is the only cut-off: unpaid Orders are cancelled then, so nothing lingers after close.
