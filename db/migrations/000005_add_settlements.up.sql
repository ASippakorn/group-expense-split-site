CREATE TABLE settlements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  payer_participant_id UUID NOT NULL REFERENCES participants(id),
  receiver_participant_id UUID NOT NULL REFERENCES participants(id),
  amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
  currency CHAR(3) NOT NULL DEFAULT 'THB',
  settlement_date DATE NOT NULL,
  note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (payer_participant_id <> receiver_participant_id)
);
CREATE INDEX settlements_group_id_idx ON settlements(group_id);
