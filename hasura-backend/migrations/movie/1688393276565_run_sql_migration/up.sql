CREATE OR REPLACE FUNCTION update_movie_status_trigger() 
RETURNS TRIGGER 
AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    PERFORM update_movie_status(NEW);
  ELSIF TG_OP = 'DELETE' THEN
    PERFORM update_movie_status(OLD);
  END IF;
  
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tickets_update_movie_status 
AFTER INSERT OR DELETE ON tickets 
FOR EACH ROW 
EXECUTE FUNCTION update_movie_status_trigger();



