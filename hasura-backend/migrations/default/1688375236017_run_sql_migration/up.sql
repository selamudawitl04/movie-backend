CREATE OR REPLACE FUNCTION calculate_movie_rating_v3(movie movies)
RETURNS NUMERIC
AS $$
BEGIN
  RETURN ROUND((SELECT AVG(rating) FROM ratings WHERE movie_id = movie.id), 2);
END;
$$ LANGUAGE plpgsql;
