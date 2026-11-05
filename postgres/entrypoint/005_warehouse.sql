-- Table: iwh_warehouses

CREATE TABLE IF NOT EXISTS iwh_warehouses (
    ID bigserial primary key,
	iwh_code character varying(32) unique NOT NULL,
    iwh_name text default '',
    iwh_namealt text DEFAULT '',
    iwh_addr1 text DEFAULT '',
    iwh_addr2 text DEFAULT '',
    iwh_city text DEFAULT '',
    iwh_state text DEFAULT '',
    iwh_zip text DEFAULT '',
    iwh_country text DEFAULT '',
    iwh_tel text DEFAULT '',
    iwh_fax text DEFAULT '',
    iwh_contact text DEFAULT '',
    iwh_memo text default '',
    iwh_seq integer DEFAULT 256,
    iwh_active boolean DEFAULT true,
    iwh_system boolean DEFAULT false,
	CreatedAt timestamp without time zone DEFAULT now(),
	UpdatedAt timestamp without time zone DEFAULT now(),
	DeletedAt timestamp without time zone,
    UNIQUE (iwh_code)

);


CREATE OR REPLACE FUNCTION prevent_delete_if_system()
RETURNS trigger AS $$
BEGIN
    IF OLD.iwh_system = true THEN
        RAISE EXCEPTION 'Delete blocked: this row cannot be deleted.';
    END IF;

    RETURN OLD; -- required for BEFORE DELETE
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS iwh_prevent_delete ON iwh_warehouses;
CREATE TRIGGER iwh_prevent_delete
BEFORE DELETE ON iwh_warehouses
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_if_system();


INSERT INTO iwh_warehouses (iwh_code, iwh_name, iwh_namealt, iwh_active, iwh_system) VALUES 
('ADMINS','Administrators','Administrators',true,true) ON CONFLICT (iwh_code) DO NOTHING;

