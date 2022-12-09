ALTER TABLE resources ADD COLUMN active BOOLEAN;
UPDATE resources SET active = 't';
ALTER TABLE resources ALTER COLUMN active SET NOT NULL;
ALTER TABLE resources ALTER COLUMN active SET DEFAULT TRUE;

ALTER TABLE resource_requests ADD COLUMN closed BOOLEAN;
UPDATE resource_requests SET closed = 'f';
ALTER TABLE resource_requests ALTER COLUMN closed SET NOT NULL;
ALTER TABLE resource_requests ALTER COLUMN closed SET DEFAULT FALSE;

ALTER TABLE resource_assignments ADD COLUMN completed BOOLEAN;
UPDATE resource_assignments SET completed = 'f';
ALTER TABLE resource_assignments ALTER COLUMN completed SET NOT NULL;
ALTER TABLE resource_assignments ALTER COLUMN completed SET DEFAULT FALSE;
