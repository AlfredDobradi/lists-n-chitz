CREATE FUNCTION set_all_links_deleted() RETURNS TRIGGER AS $$
BEGIN
	UPDATE users_x_links SET status = 0 WHERE hashlink IN (SELECT hash FROM links WHERE links.status = 0) AND users_x_links.status <> 0;
	RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER cascade_url_deletion AFTER UPDATE OF status ON links
FOR EACH STATEMENT
EXECUTE PROCEDURE public.set_all_links_deleted();