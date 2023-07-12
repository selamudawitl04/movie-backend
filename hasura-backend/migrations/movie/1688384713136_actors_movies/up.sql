CREATE OR REPLACE FUNCTION count_movies_for_actor(actor actors) 
RETURNS INT
STABLE AS $$
BEGIN
 RETURN(SELECT COUNT( DISTINCT movie_id) FROM movies_actors WHERE actor_id = actor.id);
--  RETURN movie_count;
END;
$$ LANGUAGE plpgsql;
