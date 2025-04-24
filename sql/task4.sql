INSERT INTO developers (name, department, geolocation, last_known_ip)
SELECT 
    (ARRAY['James', 'Mary', 'John', 'Patricia', 'Robert'])[floor(random()*5 + 1)::int] || ' ' ||
    (ARRAY['Smith', 'Johnson', 'Williams', 'Brown', 'Jones'])[floor(random()*5 + 1)::int] AS name,

    (ARRAY['backend', 'frontend', 'ios', 'android'])[floor(random()*4 + 1)::int] AS department,

    ST_SetSRID(
        ST_MakePoint( 
            (random()*180 - 90)::numeric(8,5), 
            (random()*360 - 180)::numeric(8,5)
        ), 4326
    )::geography AS geolocation,

    ( (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text
    )::inet AS last_known_ip

FROM generate_series(1,5000);

SELECT 
    name,
    department,
    geolocation,
    last_known_ip
FROM 
    developers
WHERE 
    ST_DWithin(
        geolocation, 
        ST_SetSRID(ST_MakePoint(20.5122, 54.7104), 4326)::geography, 
        10000  
    )
LIMIT 10000;

INSERT INTO developers (name, department, geolocation, last_known_ip)
SELECT 
    (ARRAY['James', 'Mary', 'John', 'Patricia', 'Robert'])[floor(random()*5 + 1)::int] || ' ' ||
    (ARRAY['Smith', 'Johnson', 'Williams', 'Brown', 'Jones'])[floor(random()*5 + 1)::int] AS name,

    (ARRAY['backend', 'frontend', 'ios', 'android'])[floor(random()*4 + 1)::int] AS department,

    ST_SetSRID(
        ST_MakePoint( 
            (20.5122 + (random() * 0.1 - 0.05))::numeric(8,5),  
            (54.7104 + (random() * 0.1 - 0.05))::numeric(8,5)   
        ), 4326
    )::geography AS geolocation,

    ( (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text || '.' ||
      (floor(random()*256)::int)::text
    )::inet AS last_known_ip

FROM generate_series(1,10);  
