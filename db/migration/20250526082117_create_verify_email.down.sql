DROP TABLE IF EXISTS "verify_email" CASCADE;

ALTER TABLE "users" DROP COLUMN  "is_verified_email";