CREATE TYPE gender AS ENUM ('Unknown', 'Male', 'Female', 'Not Specified');

ALTER TABLE resources ADD COLUMN sex gender;
UPDATE resources SET sex = 'Male';
ALTER TABLE resources ALTER COLUMN sex SET NOT NULL;
ALTER TABLE resources ALTER COLUMN sex SET DEFAULT 'Unknown';