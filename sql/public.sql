INSERT INTO developers (name, department, geolocation, last_known_ip)
SELECT 
    first_names.name || ' ' || last_names.name AS name,
    departments.department,
    POINT(
        random() * 180 - 90,  -- широта (-90 до 90)
        random() * 360 - 180   -- долгота (-180 до 180)
    ) AS geolocation,
    ( 
        (random() * 255)::int || '.' ||
        (random() * 255)::int || '.' ||
        (random() * 255)::int || '.' ||
        (random() * 255)::int
    )::inet AS last_known_ip
FROM 
    (VALUES ('James'), ('Mary'), ('John'), ('Patricia'), ('Robert')) AS first_names(name),
    (VALUES ('Smith'), ('Johnson'), ('Williams'), ('Brown'), ('Jones')) AS last_names(name),
    (VALUES ('backend'), ('frontend'), ('ios'), ('android')) AS departments(department)
ORDER BY random()
LIMIT 20;

SELECT * FROM developers;

SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE table_name = 'developers';

EXPLAIN ANALYZE SELECT * FROM developers WHERE name LIKE 'James%';

EXPLAIN ANALYZE SELECT * FROM developers WHERE department = 'backend';

EXPLAIN ANALYZE SELECT * FROM developers WHERE last_known_ip = '192.168.1.10'::inet;

EXPLAIN ANALYZE SELECT * FROM developers WHERE is_available = TRUE;

CREATE INDEX idx_developers_name ON developers (name varchar_pattern_ops);

CREATE INDEX idx_developers_department ON developers (department);

CREATE INDEX idx_developers_geolocation ON developers (
    (geolocation[0]), (geolocation[1])
);

CREATE INDEX idx_developers_ip ON developers (last_known_ip);
DROP INDEX IF EXISTS idx_developers_department;


INSERT INTO developers (name, department, geolocation, last_known_ip)
SELECT 
    (ARRAY['James', 'Mary', 'John', 'Patricia', 'Robert'])[floor(random()*5 + 1)::int] || ' ' ||
    (ARRAY['Smith', 'Johnson', 'Williams', 'Brown', 'Jones'])[floor(random()*5 + 1)::int] AS name,

    (ARRAY['backend', 'frontend', 'ios', 'android'])[floor(random()*5 + 1)::int] AS department,

    POINT( 
        (random()*180 - 90)::numeric(8,5), 
        (random()*360 - 180)::numeric(8,5)
    ) AS geolocation,

    ( (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text
    )::inet AS last_known_ip

FROM generate_series(1,5000);