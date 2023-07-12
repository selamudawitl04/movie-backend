CREATE OR REPLACE FUNCTION calculate_movie_rating_v2(movie_id INTEGER)
RETURNS NUMERIC
AS $$
BEGIN
  RETURN ROUND((SELECT AVG(ratings) FROM movies WHERE id = movie_id), 2);
END;
$$ LANGUAGE plpgsql;
