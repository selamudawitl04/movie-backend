-- Could not auto-generate a down migration.
-- Please write an appropriate down migration for the SQL below:
-- CREATE OR REPLACE FUNCTION calculate_movie_rating(id movies)
-- RETURNS NUMERIC
-- AS $$
-- BEGIN
--   RETURN ROUND((SELECT AVG(ratings) FROM movies WHERE movie_id = id), 2);
-- END;
-- $$ LANGUAGE plpgsql;
