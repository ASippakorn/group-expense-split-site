ALTER TABLE expenses DROP CONSTRAINT IF EXISTS expenses_split_type_check;
ALTER TABLE expenses ADD CONSTRAINT expenses_split_type_check CHECK (split_type IN ('equal', 'manual_amount', 'percentage'));
ALTER TABLE expenses DROP COLUMN IF EXISTS tag_id;
DROP TABLE IF EXISTS tag_participants;
DROP TABLE IF EXISTS tags;
