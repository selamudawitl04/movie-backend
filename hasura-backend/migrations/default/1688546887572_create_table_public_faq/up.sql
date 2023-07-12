CREATE TABLE "public"."faqs" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "question" text NOT NULL, "answer" text NOT NULL, "date" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("id") );
CREATE EXTENSION IF NOT EXISTS pgcrypto;
