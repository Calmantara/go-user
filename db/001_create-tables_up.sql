CREATE TABLE public.users (
	id bigserial NOT NULL,
	"name" text NULL,
	address text NULL,
	email text NULL,
	"password" text NULL,
	created_by text NULL,
	updated_by text NULL,
	record_flag text NULL DEFAULT 'ACTIVE'::text,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
	CONSTRAINT unique_email UNIQUE (email)
);
CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);
CREATE INDEX idx_users_record_flag ON public.users USING btree (record_flag);

CREATE TABLE public.photos (
	id bigserial NOT NULL,
	user_id int NULL,
	"name" text NULL,
	created_by text NULL,
	updated_by text NULL,
	record_flag text NULL DEFAULT 'ACTIVE'::text,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT photos_pkey PRIMARY KEY (id),
	CONSTRAINT fk_users_photos FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX idx_photos_deleted_at ON public.photos USING btree (deleted_at);
CREATE UNIQUE INDEX idx_photos_name_user_id ON public.photos USING btree (user_id, name);
CREATE INDEX idx_photos_record_flag ON public.photos USING btree (record_flag);

CREATE TABLE public.credit_card_tokens (
	id bigserial NOT NULL,
	user_id int NULL,
	"token" text NULL,
	created_by text NULL,
	updated_by text NULL,
	record_flag text NULL DEFAULT 'ACTIVE'::text,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	CONSTRAINT credit_card_tokens_pkey PRIMARY KEY (id),
	CONSTRAINT fk_users_credit_card_token FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX idx_credit_card_tokens_deleted_at ON public.credit_card_tokens USING btree (deleted_at);
CREATE INDEX idx_credit_card_tokens_record_flag ON public.credit_card_tokens USING btree (record_flag);