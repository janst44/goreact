-- +goose Up
-- +goose StatementBegin
-- First add the column as nullable
ALTER TABLE todos
ADD COLUMN user_id UUID NULL REFERENCES users(id);

-- Create a system user for existing todos
INSERT INTO users (id, email, password, name)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'system@example.com',
    'not-a-real-password',
    'System User'
) ON CONFLICT (id) DO NOTHING;

-- Update existing todos with the system user
UPDATE todos
SET user_id = '00000000-0000-0000-0000-000000000000'
WHERE user_id IS NULL;

-- Now make the column NOT NULL
ALTER TABLE todos
ALTER COLUMN user_id SET NOT NULL;

-- Add the index
CREATE INDEX idx_todos_user_id ON todos(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_todos_user_id;
ALTER TABLE todos DROP COLUMN user_id;
-- +goose StatementEnd