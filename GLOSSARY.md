# Class Campaign

A LINE MINI App where a school class runs recurring money-collecting campaigns (merch orders, donations, parties), replacing copy-pasted lists and slip albums in the LINE group.

## People

**Applicant**:
A person who has opened the app and filled in their profile but is not yet approved by an Admin. Sees nothing until approved.
_Avoid_: Pending user, guest

**Member**:
A person approved by an Admin, or created by an Admin on someone's behalf. Usually has a LINE account linked; may not.
_Avoid_: User, friend, customer

**Unlinked Member**:
A Member created by an Admin with no LINE account yet, for people who order through an Admin instead of using the app. An Admin links it when that person joins.
_Avoid_: Guest, placeholder, offline user

**Admin**:
A Member who may create Campaigns and edit any data at any time. There can be several.
_Avoid_: Owner, organiser, treasurer

**Cohort**:
The class year (รุ่น) a Member says they belong to. Shown only; it grants nothing.
_Avoid_: Group, class, batch

## Campaigns

**Group**:
The community sharing this app; its name may change. Exactly one exists today.
_Avoid_: Tenant, organisation, cohort

**Campaign**:
One collection effort with a purpose, a bank account and QR code to pay into, and a status of draft, open, or closed. Closing takes no new Orders and cancels Orders still awaiting payment.
_Avoid_: Event (an event is one kind of Campaign), project, post

**Option**:
A purchasable thing offered in a Campaign, with a price and the details each unit needs (for example a jersey design needing size, print name, number). A donation Option lets the Member choose the price.
_Avoid_: Product, SKU, variant

**Sponsor Position**:
A spot on a Campaign's merch where a sponsor's name or logo goes, assigned by an Admin.
_Avoid_: Slot, logo spot

## Ordering

**Order**:
One checkout by a Member in a Campaign: the Items they chose, how they receive them, and its Payment. A Member may place many Orders in one Campaign.
_Avoid_: Participation, entry, registration, cart

**Item**:
A line in an Order: an Option, its details, and a quantity. All units in one Item are identical.
_Avoid_: Order line, product

**Delivery Method**:
How an Order reaches the Member: shipped to the address captured on that Order, or picked up.
_Avoid_: Shipping option

**Balance**:
What an Order still owes: Item totals minus its confirmed Payments, plus Adjustments. A cancelled Order owes nothing for its Items.
_Avoid_: Debt, outstanding, total

## Money

**Payment**:
A Slip submitted for the full amount of an Order, which an Admin confirms or rejects (with a reason). Only confirmed Payments count. May be recorded by an Admin without a Slip, for cash.
_Avoid_: Transaction, transfer

**Slip**:
The image of a bank transfer receipt attached to a Payment as evidence.
_Avoid_: Receipt, proof

**Adjustment**:
An Admin-only correction to what an Order owes, positive or negative, always with a note: refunds, discounts, over- or under-transfers.
_Avoid_: Refund, negative payment, credit
