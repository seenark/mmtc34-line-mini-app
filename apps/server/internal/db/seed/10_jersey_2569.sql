-- Idempotent: each INSERT is guarded by the Campaign title / Option name.
INSERT INTO campaigns (type, title, description, bank_name, bank_account_no, bank_account_name, status, close_at, created_at)
SELECT 'merch', 'เสื้อรุ่น 2569',
       'สั่งเสื้อรุ่น 2569 เลือกไซซ์ ใส่ชื่อและเบอร์ที่จะพิมพ์บนเสื้อ แล้วโอนเงินตามยอดที่ระบบแจ้ง',
       'ธนาคารกสิกรไทย', '000-0-00000-0', 'บัญชีตัวอย่าง', 'open', NULL, '2026-10-24T00:00:00Z'
WHERE NOT EXISTS (SELECT 1 FROM campaigns WHERE title = 'เสื้อรุ่น 2569');

INSERT INTO options (campaign_id, name, description, price, sizes, needs_print_name, needs_number, sort_order)
SELECT c.id, 'เสื้อรุ่น', '', 260,
       '["K2XS","K3XS","KXS","2XS","S","M","L","XL","2XL","3XL","4XL","6XL"]', 1, 1, 0
FROM campaigns c
WHERE c.title = 'เสื้อรุ่น 2569'
  AND NOT EXISTS (SELECT 1 FROM options o WHERE o.campaign_id = c.id AND o.name = 'เสื้อรุ่น');
