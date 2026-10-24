# Account deletion keeps Slips and Order lines as financial evidence

When a Member deletes their account, we erase their contact data (real name, Cohort, address, phone, LINE link, and the shipping fields on their Orders) but keep their nickname, Orders, Items, Payments, Adjustments and Slip images. The owner chose this during design: the class needs a complete record of who paid what for each Campaign, and a Slip is the only evidence of a bank transfer. A reviewer will see that a deleted Member's Slip still opens for Admins, which looks like a PDPA gap; it is a deliberate retention of accounting evidence, not an oversight.

## Consequences

- Slip images can show the payer's full name and bank account. Only Admins can open them after deletion.
- Free text the Member typed into an Order, such as a print name, stays with the Order.
- Erased values can remain in SQLite free pages until the database is vacuumed.
