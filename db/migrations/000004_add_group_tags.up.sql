CREATE TABLE tags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES groups(id),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT tags_group_name_unique UNIQUE (group_id, name)
);

CREATE TABLE tag_participants (
  tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  participant_id UUID NOT NULL REFERENCES participants(id),
  PRIMARY KEY (tag_id, participant_id)
);

ALTER TABLE expenses ADD COLUMN tag_id UUID REFERENCES tags(id);
ALTER TABLE expenses DROP CONSTRAINT expenses_split_type_check;
ALTER TABLE expenses ADD CONSTRAINT expenses_split_type_check CHECK (split_type IN ('equal', 'manual_amount', 'percentage', 'tag'));
