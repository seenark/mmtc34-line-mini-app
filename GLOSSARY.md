# MMTC Events

A LINE MINI App where a school class runs recurring money-collecting campaigns (merch orders, donations, parties), replacing copy-pasted lists and slip albums in the LINE group.

## People

**Member**:
A person who has opened the app and filled in their profile, or whom an Admin created on someone's behalf. Nobody approves Members. Usually has a LINE account linked; may not.
_Avoid_: User, friend, customer, applicant

**Unlinked Member**:
A Member created by an Admin with no LINE account, for people who order through an Admin instead of using the app. When that person joins, an Admin links their Orders to the Member who has their LINE account; an Unlinked Member left with no Orders disappears.
_Avoid_: Guest, placeholder, offline user

**Admin**:
A Member who may create Campaigns and edit any data at any time. There can be several.
_Avoid_: Owner, organiser, treasurer

**Admin Mode**:
The state in which an Admin sees and uses admin tools. Outside it an Admin uses the app as an ordinary Member, and the Orders they place are their own.
_Avoid_: Admin login, admin account

**Super Admin**:
The one Admin account defined in the server configuration and signed into with a password, not LINE. It exists to promote the first Admins; it cannot be demoted or deleted.
_Avoid_: Root, owner, bootstrap admin

**Cohort**:
The class year (รุ่น) a Member says they belong to. Shown only; it grants nothing.
_Avoid_: Group, class, batch

## Campaigns

**Group**:
The community sharing this app; its name may change. Exactly one exists today.
_Avoid_: Tenant, organisation, cohort

**Campaign**:
One collection effort with a purpose, a bank account and QR code to pay into, and a status of draft, open, or closed. Closing takes no new Orders and cancels Orders still awaiting payment.
_Avoid_: Event (an event is one kind of Campaign), project, post. The app shows a Campaign to users as `อีเวนต์` (event); it is still distinct from the `event` type, which the UI labels `กิจกรรม`.

**Option**:
A purchasable thing offered in a Campaign, with a price and the details each unit needs (for example a jersey design needing size, print name, number, with a size chart to choose from). A donation Option lets the Member choose the price.
_Avoid_: Product, SKU, variant

**Sponsor Position**:
A spot on a Campaign's merch where a sponsor's name or logo goes, assigned by an Admin.
_Avoid_: Slot, logo spot

## Ordering

**Order**:
One checkout in a Campaign, by a Member or by an Admin on a Member's behalf: the Items chosen, how they are received, and its Payment. A Member may place many Orders in one Campaign, each shipped to a different address if they like.
_Avoid_: Participation, entry, registration, cart

**Item**:
A line in an Order: an Option, its details, and a quantity. All units in one Item are identical. Its details may be corrected by the Member while the Campaign is open; its Option and quantity are fixed once the Order is placed.
_Avoid_: Order line, product

**Delivery Method**:
How an Order reaches the Member: shipped to the address chosen for that Order, or picked up. An Order holding only donations has none.
_Avoid_: Shipping option

**Balance**:
What an Order still owes: Item totals minus its confirmed Payments, plus Adjustments. A cancelled Order owes nothing for its Items.
_Avoid_: Debt, outstanding, total

## Money

**Payment**:
A Slip submitted for the full amount of an Order, which an Admin confirms or rejects (with a reason). Only confirmed Payments count. May be recorded by an Admin without a Slip, for cash.
_Avoid_: Transaction, transfer

**Slip**:
The images (one to six per upload) of the bank transfer receipts behind a Payment, kept as evidence. An Admin may add or remove images later; a Slip always keeps at least one.
_Avoid_: Receipt, proof

**Adjustment**:
An Admin-only correction to what an Order owes, positive or negative, always with a note: refunds, discounts, over- or under-transfers.
_Avoid_: Refund, negative payment, credit
