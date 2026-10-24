-- Idempotent: each INSERT is guarded by the Campaign title and position name.
INSERT INTO sponsor_positions (campaign_id, name, sponsor_name, sort_order)
SELECT c.id, v.name, v.sponsor_name, v.sort_order
FROM campaigns c
JOIN (
    SELECT 'กลางอก' AS name, 'TK.LINE' AS sponsor_name, 10 AS sort_order
    UNION ALL SELECT 'ชายเสื้อด้านหน้าขวา', '', 20
    UNION ALL SELECT 'แขนขวา', '', 30
    UNION ALL SELECT 'แขนซ้าย', 'CodeSook', 40
    UNION ALL SELECT 'ชายเสื้อด้านหลังกลาง', '', 50
) v
WHERE c.title = 'เสื้อรุ่น 2569'
  AND NOT EXISTS (SELECT 1 FROM sponsor_positions p WHERE p.campaign_id = c.id AND p.name = v.name);
