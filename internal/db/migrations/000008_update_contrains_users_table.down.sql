
ALTER TABLE z_users
ADD CONSTRAINT z_users_parent_id_fkey
FOREIGN KEY (parent_id) REFERENCES z_users(user_id);