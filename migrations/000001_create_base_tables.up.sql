CREATE TABLE "clearances" (
  "id" bigserial PRIMARY KEY,
  "description" varchar UNIQUE NOT NULL
);

CREATE TABLE "positions" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE NOT NULL
);

CREATE TABLE "resources" (
  "id" int PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "position_id" int NOT NULL,
  "clearance_id" int NOT NULL,
  "specialties" text[],
  "certifications" text[]
);

CREATE TABLE "resource_requests" (
  "id" bigserial PRIMARY KEY,
  "customer" varchar NOT NULL,
  "start_date" date,
  "end_date" date,
  "hours_per_week" int DEFAULT 40,
  "skills" text[] NOT NULL,
  "opportunity_id" varchar,
  "engagement_id" varchar,
  "created_at" timestamp(0) DEFAULT (now()),
  "updated_at" timestamp(0) DEFAULT (now()),
  "version" int DEFAULT 1
);

CREATE TABLE "resource_assignments" (
  "resource_request_id" int,
  "resource_id" int,
  "hours_per_week" int DEFAULT 40,
  "created_at" timestamp(0) DEFAULT (now()),
  "updated_at" timestamp(0) DEFAULT (now()),
  "version" int DEFAULT 1,
  PRIMARY KEY(resource_request_id, resource_id)
);

CREATE INDEX "idx_resource_clearance" ON "resources" ("clearance_id");

CREATE INDEX "idx_resource_position" ON "resources" ("position_id");

ALTER TABLE "resources" ADD FOREIGN KEY ("position_id") REFERENCES "positions" ("id");

ALTER TABLE "resources" ADD FOREIGN KEY ("clearance_id") REFERENCES "clearances" ("id");

ALTER TABLE "resource_assignments" ADD FOREIGN KEY ("resource_request_id") REFERENCES "resource_requests" ("id");

ALTER TABLE "resource_assignments" ADD FOREIGN KEY ("resource_id") REFERENCES "resources" ("id");
