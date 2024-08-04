CREATE TABLE if not exists hrbp
(
    id serial PRIMARY KEY,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL DEFAULT now(),
    email character varying(255) NULL,
    name character varying(255) NULL,
    deleted_at timestamp(6) without time zone NULL
);

CREATE TABLE if not exists departments
(
    id serial PRIMARY KEY,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NULL DEFAULT now(),
    deleted_at timestamp without time zone NULL,
    name character varying(255) NULL
);

ALTER TABLE users
    ADD COLUMN if not exists department_id integer NULL,
    ADD COLUMN if not exists hrbp_id integer NULL,
    ADD COLUMN if not exists employee_id character varying(255) NULL;