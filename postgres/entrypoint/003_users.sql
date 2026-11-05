-- Table: yug_users

CREATE TABLE IF NOT EXISTS yus_users (
    ID bigserial primary key,
	yus_code character varying(32) unique NOT NULL,
    yus_name text default '',
    yus_passwd text default '',
    yus_memo text default '',
    yus_seq integer DEFAULT 256,
    yus_active boolean DEFAULT true,
    yus_system boolean DEFAULT false,
	CreatedAt timestamp without time zone DEFAULT now(),
	UpdatedAt timestamp without time zone DEFAULT now(),
	DeletedAt timestamp without time zone,
    UNIQUE (yus_code)
);

CREATE OR REPLACE FUNCTION prevent_delete_if_system()
RETURNS trigger AS $$
BEGIN
    IF OLD.yus_system = true THEN
        RAISE EXCEPTION 'Delete blocked: this row cannot be deleted.';
    END IF;

    RETURN OLD; -- required for BEFORE DELETE
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS yus_prevent_delete ON yus_users;
CREATE TRIGGER yus_prevent_delete
BEFORE DELETE ON yus_users
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_if_system();


INSERT INTO yus_users (yus_code, yus_name, yus_passwd, yus_memo, yus_active, yus_system) VALUES 
('ADMIN','Administrator','admin','',true,true) ON CONFLICT (yus_code) DO NOTHING;
