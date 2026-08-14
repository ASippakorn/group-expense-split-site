INSERT INTO users (id, email, password_hash)
VALUES
  ('10000000-0000-0000-0000-000000000001', 'ada@example.com', 'replace-with-local-hash'),
  ('10000000-0000-0000-0000-000000000002', 'ben@example.com', 'replace-with-local-hash'),
  ('10000000-0000-0000-0000-000000000003', 'chai@example.com', 'replace-with-local-hash')
ON CONFLICT (email) DO NOTHING;

INSERT INTO groups (id, name, default_currency, description, owner_id)
VALUES (
  '20000000-0000-0000-0000-000000000001',
  'Bangkok Food Crawl',
  'THB',
  'Demo group for local review.',
  '10000000-0000-0000-0000-000000000001'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO participants (group_id, user_id, role, active)
VALUES
  ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'owner', true),
  ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002', 'participant', true),
  ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000003', 'participant', true)
ON CONFLICT (group_id, user_id) DO NOTHING;
