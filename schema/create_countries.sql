CREATE TABLE countries (
    id integer NOT NULL,
    name character varying(255),
    slug character varying(255),
    caption text,
    body text,
    image character varying(255)
);

CREATE SEQUENCE countries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE countries_id_seq OWNED BY countries.id;
ALTER TABLE ONLY countries ALTER COLUMN id SET DEFAULT nextval('countries_id_seq'::regclass);
