CREATE OR REPLACE FUNCTION count_movies_for_directors(director directors) 
RETURNS INT
STABLE AS $$
BEGIN
 RETURN(SELECT COUNT( DISTINCT director_id) FROM movies WHERE director_id = director.id);
--  RETURN movie_count;
END;
$$ LANGUAGE plpgsql;
