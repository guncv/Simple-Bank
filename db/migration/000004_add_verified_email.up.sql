ALTER TABLE "users"
ADD COLUMN "is_email_verified" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE "verify_emails" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR NOT NULL REFERENCES "users"("username") ON DELETE CASCADE,
  "email" VARCHAR NOT NULL,
  "secret_code" VARCHAR NOT NULL,
  "is_used" BOOLEAN NOT NULL DEFAULT FALSE,
  "expired_at" TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
