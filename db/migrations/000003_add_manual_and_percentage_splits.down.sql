ALTER TABLE expense_splits DROP COLUMN IF EXISTS percentage_basis_points;

ALTER TABLE expenses DROP CONSTRAINT IF EXISTS expenses_split_type_check;
ALTER TABLE expenses ADD CONSTRAINT expenses_split_type_check CHECK (split_type IN ('equal'));
