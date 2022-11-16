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
  "duration" varchar,
  "skills" text[] NOT NULL,
  "fulltime" boolean DEFAULT true,
  "percent_required" numeric DEFAULT 1,
  "project_id" varchar,
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "resource_assignments" (
  "resource_request_id" int,
  "resource_id" int,
  "percentage" numeric DEFAULT 1,
  "created_at" datetime DEFAULT (now()),
  "updated_at" datetime,
  "version" int DEFAULT 1,
  PRIMARY KEY(resource_request_id, resource_id)
);

CREATE INDEX "idx_resource_clearance" ON "resources" ("clearance_id");

CREATE INDEX "idx_resource_position" ON "resources" ("position_id");

ALTER TABLE "resources" ADD FOREIGN KEY ("position_id") REFERENCES "positions" ("id");

ALTER TABLE "resources" ADD FOREIGN KEY ("clearance_id") REFERENCES "clearances" ("id");

ALTER TABLE "resource_assignments" ADD FOREIGN KEY ("resource_request_id") REFERENCES "resource_requests" ("id");

ALTER TABLE "resource_assignments" ADD FOREIGN KEY ("resource_id") REFERENCES "resources" ("id");
