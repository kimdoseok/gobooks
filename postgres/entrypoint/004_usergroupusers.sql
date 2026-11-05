-- Table: ygu_usergroupusers

CREATE TABLE IF NOT EXISTS ygu_usergroupusers (
    ID bigserial primary key,
	ygu_yug_code character varying(32) unique NOT NULL,
	ygu_yus_code character varying(32) unique NOT NULL,
    ygu_active boolean DEFAULT true,
	CreatedAt timestamp without time zone DEFAULT now(),
	UpdatedAt timestamp without time zone DEFAULT now(),
	DeletedAt timestamp without time zone,
    UNIQUE (ygu_yug_code, ygu_yus_code)
);

INSERT INTO ygu_usergroupusers (ygu_yus_code, ygu_yug_code, ygu_active) VALUES 
('ADMIN','ADMINS',true) ON CONFLICT (ygu_yus_code, ygu_yug_code) DO NOTHING;
