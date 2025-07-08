BEGIN;
ALTER TABLE z_entity
ADD CONSTRAINT user_id_unique UNIQUE (user_id);
COMMIT;
