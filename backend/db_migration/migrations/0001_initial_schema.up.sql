CREATE TABLE if not exists achievement
(
                        id serial PRIMARY KEY,
                        achievement_name character varying(255) NULL,
                        display_name character varying(255) NULL,
                        emoji character varying(255) NULL,
                        points character varying(255) NULL,
                        created_at timestamp without time zone NULL,
                        updated_at timestamp without time zone NULL,
                        is_active boolean NULL
);

CREATE TABLE if not exists kleos
(
                        id serial PRIMARY KEY,
                        sender_id character varying(255) NULL,
                        message text NULL,
                        achievement character varying(255) NULL,
                        receiver_id character varying(255) NULL,
                        year character varying NULL,
                        month character varying NULL,
                        week character varying NULL,
                        day character varying NULL,
                        created_at timestamp without time zone NULL,
                        updated_at timestamp without time zone NULL,
                        deleted_at timestamp without time zone NULL
);

CREATE TABLE if not exists managers
(
                        id serial PRIMARY KEY,
                        email character varying(255) NULL,
                        slack_user_id character varying(255) NULL,
                        user_name character varying(255) NULL,
                        real_name character varying(255) NULL,
                        created_at timestamp without time zone NULL,
                        updated_at timestamp without time zone NULL,
                        deleted_at timestamp without time zone NULL
);

CREATE TABLE if not exists user_count
(
                        id serial PRIMARY KEY,
                        created_at timestamp without time zone NOT NULL DEFAULT now(),
                        updated_at timestamp without time zone NULL DEFAULT now(),
                        deleted_at timestamp without time zone NULL,
                        user_id character varying(255) NULL,
                        month character varying(255) NULL,
                        week character varying(255) NULL,
                        given_count bigint NULL,
                        received_count bigint NULL
);

CREATE TABLE if not exists users (
                        id serial PRIMARY KEY,
                        created_at timestamp without time zone NULL,
                        updated_at timestamp without time zone NULL,
                        deleted_at timestamp without time zone NULL,
                        slack_user_id character varying(255) NULL,
                        user_name character varying(255) NULL,
                        email character varying(255) NULL,
                        slack_image_url character varying(255) NULL,
                        real_name character varying(255) NULL,
                        manager_id integer NULL,
                        given_count integer NULL,
                        received_count integer NULL
);

CREATE TABLE if not exists shedlock
(
                        id SERIAL PRIMARY KEY,
                        name character varying(255) NOT NULL,
                        lock_until timestamp without time zone NULL,
                        locked_at timestamp without time zone NULL,
                        locked_by character varying(255) NULL,
                        locked_value boolean NULL,
                        unlocked_by character varying(255) NULL
);