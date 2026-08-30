ALTER TABLE expenses DROP CONSTRAINT IF EXISTS expenses_split_type_check;
ALTER TABLE expenses ADD CONSTRAINT expenses_split_type_check CHECK (split_type IN ('equal', 'manual_amount', 'percentage'));

ALTER TABLE expense_splits ADD COLUMN percentage_basis_points BIGINT NOT NULL DEFAULT 0 CHECK (percentage_basis_points >= 0 AND percentage_basis_points <= 10000);
