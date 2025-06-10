ALTER TABLE users ADD COLUMN role varchar(10) NOT NULL DEFAULT 'user';
CREATE INDEX idx_users_role ON users(role);