CREATE TABLE IF NOT EXISTS iit_items (
    ID bigserial primary key,
	iit_code character varying(32) unique NOT NULL,
    iit_ict_code varchar(16) DEFAULT '',
    iit_pvd_code varchar(16) DEFAULT '',
    iit_ium_code varchar(16) DEFAULT '',
    iit_ity_code varchar(16) DEFAULT '',
    iit_ytx_code varchar(16) DEFAULT '',
    iit_oct_code varchar(16) DEFAULT '',
    iit_iwh_code varchar(16) DEFAULT '',
    iit_gac_inven varchar(16) DEFAULT '',
    iit_ean13 varchar(32) DEFAULT '',
    iit_name text default '',
    iit_namealt text default '',
    iit_memo text default '',
    iit_retroday integer DEFAULT 1 NOT NULL,
    iit_price_base numeric(16,2) DEFAULT 0,
    iit_price_part numeric(16,2) DEFAULT 0,
    iit_bbu_code varchar(16) DEFAULT '',
    iit_ist_code varchar(16) DEFAULT '',
    iit_dclimit numeric(6,2) DEFAULT 0,
    iit_weight numeric(12,3) DEFAULT 0,
    iit_volume numeric(12,3) DEFAULT 0,
    iit_taxexem numeric(12,2) DEFAULT 0,
    iit_vendor_price numeric(16,4) DEFAULT 0,
    iit_qty_base integer DEFAULT 1,
    iit_seq integer DEFAULT 256,
    iit_allowdc boolean DEFAULT false,
    iit_qdiscqty integer DEFAULT 0,
    iit_qdiscrate numeric(6,2) DEFAULT 0,
    iit_buy boolean DEFAULT true,
    iit_sell boolean DEFAULT true,
    iit_inventory boolean DEFAULT true,
    iit_taxable boolean DEFAULT true,
    iit_active boolean DEFAULT true,
    iit_system boolean DEFAULT false,
	CreatedAt timestamp without time zone DEFAULT now(),
	UpdatedAt timestamp without time zone DEFAULT now(),
	DeletedAt timestamp without time zone,
    UNIQUE (iit_code)
);

CREATE TABLE IF NOT EXISTS iis_itemstatus (
    ID bigserial primary key,
	iis_iit_code character varying(32) unique NOT NULL,
    iis_lastbuy date DEFAULT current_date,
    iis_cost_avg numeric(16,4) DEFAULT 0,
    iis_cost_std numeric(16,4) DEFAULT 0,
    iis_cost_rtrn numeric(16,4) DEFAULT 0,
    iis_cost_lst numeric(16,4) DEFAULT 0,
    iis_qty_oh integer DEFAULT 0,
    iis_pqty_oh integer DEFAULT 0,
    iis_qty_bk integer DEFAULT 0,
    iis_pqty_bk integer DEFAULT 0,
    iis_qty_po integer DEFAULT 0,
    iis_pqty_po integer DEFAULT 0,
    iis_qty_it integer DEFAULT 0,
    iis_pqty_it integer DEFAULT 0,
    iis_qty_al integer DEFAULT 0,
    iis_pqty_al integer DEFAULT 0,
	CreatedAt timestamp without time zone DEFAULT now(),
	UpdatedAt timestamp without time zone DEFAULT now(),
	DeletedAt timestamp without time zone,
    UNIQUE (iis_iit_code)
);


CREATE OR REPLACE FUNCTION prevent_delete_if_system()
RETURNS trigger AS $$
BEGIN
    IF OLD.iit_system = true THEN
        RAISE EXCEPTION 'Delete blocked: this row cannot be deleted.';
    END IF;

    RETURN OLD; -- required for BEFORE DELETE
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS add_status_if_insert ON iit_items;
CREATE TRIGGER iit_prevent_delete
BEFORE DELETE ON iit_items
FOR EACH ROW
EXECUTE FUNCTION prevent_delete_if_system();

CREATE OR REPLACE FUNCTION prevent_delete_if_system()
RETURNS trigger AS $$
BEGIN
    INSERT INTO iis_itemstatus (iis_iit_code) VALUES (NEW.iit_code ON CONFLICT (iis_iit_code) DO NOTHING;
    RETURN NEW; 
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS iit_add_status ON iit_items;
CREATE TRIGGER iit_add_status
AFTER INSERT ON iit_items
FOR EACH ROW
EXECUTE FUNCTION add_status_if_insert();


