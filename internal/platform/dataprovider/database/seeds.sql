-- Default pack sizes
INSERT INTO pack_sizes(size, active) VALUES (250, true)
ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active;

INSERT INTO pack_sizes(size, active) VALUES (500, true)
ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active;

INSERT INTO pack_sizes(size, active) VALUES (1000, true)
ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active;

INSERT INTO pack_sizes(size, active) VALUES (2000, true)
ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active;

INSERT INTO pack_sizes(size, active) VALUES (5000, true)
ON CONFLICT (size) DO UPDATE SET active = EXCLUDED.active;


