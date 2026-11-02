-- Idempotent: each INSERT is guarded by the Campaign title / Option name.
-- Image names are the embedded seed/images/*.jpg files; Seed copies them into DATA_DIR/files/images.
INSERT INTO campaigns (type, title, description, bank_code, bank_account_no, bank_account_name, qr_image, status, close_at, created_at)
SELECT 'merch', 'เสื้อฟุตบอล MMTC 2026',
       'สั่งเสื้อฟุตบอล MMTC 2026 เลือกไซซ์ ใส่ชื่อและเบอร์ที่จะพิมพ์บนเสื้อ แล้วโอนเงินเข้าบัญชีของรุ่นตามยอดที่ระบบแจ้ง',
       '004', '101-8-58883-4', 'ชัฏณิจสิทธิ์ ขวัญคีรี', '', 'open',
       '2026-11-06T16:59:59Z', '2026-10-26T00:00:00Z'
WHERE NOT EXISTS (SELECT 1 FROM campaigns WHERE title = 'เสื้อฟุตบอล MMTC 2026');

INSERT INTO campaign_photos (campaign_id, name, sort_order)
SELECT c.id, '0cc8f1050e894c7bfb982b1fe8562ce1.jpg', 0
FROM campaigns c
WHERE c.title = 'เสื้อฟุตบอล MMTC 2026'
  AND NOT EXISTS (SELECT 1 FROM campaign_photos p WHERE p.campaign_id = c.id);

INSERT INTO options (campaign_id, name, description, image, size_chart, price, sizes, needs_print_name, needs_number, sort_order)
SELECT c.id, 'เสื้อ', 'ราคารวมค่าจัดส่งแล้ว',
       '0cc8f1050e894c7bfb982b1fe8562ce1.jpg', 'ea315e4dafc245d0f0689635c0d39223.jpg', 260,
       '["S","M","L","XL","2XL","3XL","4XL","5XL","6XL","7XL","8XL","K4XS","K3XS","K2XS","KXS","KS","2XS","XS"]',
       1, 1, 0
FROM campaigns c
WHERE c.title = 'เสื้อฟุตบอล MMTC 2026'
  AND NOT EXISTS (SELECT 1 FROM options o WHERE o.campaign_id = c.id AND o.name = 'เสื้อ');

INSERT INTO options (campaign_id, name, description, price, sizes, needs_print_name, needs_number, sort_order)
SELECT c.id, 'สมทบทุน', 'สมทบทุนเพิ่มเติมตามกำลังศรัทธา ไม่รับเสื้อ', NULL, '[]', 0, 0, 1
FROM campaigns c
WHERE c.title = 'เสื้อฟุตบอล MMTC 2026'
  AND NOT EXISTS (SELECT 1 FROM options o WHERE o.campaign_id = c.id AND o.name = 'สมทบทุน');
