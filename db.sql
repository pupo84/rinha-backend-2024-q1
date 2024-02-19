DROP TABLE public.balances;
DROP TABLE public.transactions;
DROP TABLE public.users;

CREATE TABLE public.users (
	id int8 NOT NULL,
	"name" varchar NOT NULL,
	email varchar NOT NULL,
	"document" varchar NOT NULL,
	"limit" int8 DEFAULT 0 NOT NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT users_email_uk UNIQUE (email),
	CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE INDEX users_document_idx ON public.users USING btree (document);

CREATE TABLE public.balances (
	id int8 NOT NULL,
	user_id int8 NOT NULL,
	amount int8 DEFAULT 0 NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT balances_pk PRIMARY KEY (id),
	CONSTRAINT balances_user_id_uk UNIQUE (user_id),
	CONSTRAINT balances_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.transactions (
	id int8 NOT NULL,
	user_id int8 NOT NULL,
	"type" public."transaction_type" NOT NULL,
	amount int8 NOT NULL,
	description varchar NOT NULL,
	created_at timestamptz NOT NULL,
	CONSTRAINT transactions_pk PRIMARY KEY (id),
	CONSTRAINT transactions_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);

DROP SEQUENCE public.balances_seq;

CREATE SEQUENCE public.balances_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

DROP SEQUENCE public.transactions_seq;

CREATE SEQUENCE public.transactions_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;

DROP SEQUENCE public.users_seq;

CREATE SEQUENCE public.users_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
