-- CREATE OR REPLACE FUNCTION update_movie_status_trigger() 
-- RETURNS TRIGGER 
-- AS $$
-- BEGIN
--   IF TG_OP = 'INSERT' THEN
--     SELECT update_movie_status(NEW);
--   ELSIF TG_OP = 'DELETE' THEN
--     SELECT update_movie_status(OLD);
--   END IF;
  
--   RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER tickets_update_movie_status 
-- AFTER INSERT OR DELETE ON tickets 
-- FOR EACH ROW 
-- EXECUTE FUNCTION update_movie_status_trigger();

CREATE OR REPLACE FUNCTION update_movie_status(ticket tickets) 
RETURNS VOID 
AS $$
BEGIN
  UPDATE movies SET status = 'closed' 
  WHERE id = ticket.movie_id AND (
    SELECT COUNT(*) FROM tickets WHERE movie_id = ticket.movie_id 
  ) >= 20;
END;
$$ LANGUAGE plpgsql;


-- DROP TRIGGER IF EXISTS tickets_update_movie_status ON tickets;
-- DROP FUNCTION IF EXISTS update_movie_status();
