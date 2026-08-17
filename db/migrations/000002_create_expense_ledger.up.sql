CREATE TABLE expenses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  payer_participant_id UUID NOT NULL REFERENCES participants(id),
  description TEXT NOT NULL,
  amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
  currency CHAR(3) NOT NULL DEFAULT 'THB',
  expense_date DATE NOT NULL,
  split_type TEXT NOT NULL CHECK (split_type IN ('equal')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX expenses_group_id_idx ON expenses(group_id);
CREATE INDEX expenses_payer_participant_id_idx ON expenses(payer_participant_id);
CREATE INDEX expenses_expense_date_idx ON expenses(expense_date);

CREATE TABLE expense_splits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
  participant_id UUID NOT NULL REFERENCES participants(id),
  amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (expense_id, participant_id)
);

CREATE INDEX expense_splits_expense_id_idx ON expense_splits(expense_id);
CREATE INDEX expense_splits_participant_id_idx ON expense_splits(participant_id);
