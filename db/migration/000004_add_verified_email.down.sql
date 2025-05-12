ALTER TABLE "users"
DROP COLUMN "is_email_verified";

DROP TABLE "verify_emails" CASCADE;