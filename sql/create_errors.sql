-- public.errors definition

-- Drop table

-- DROP TABLE public.errors;

CREATE TABLE public.errors (
	id serial4 NOT NULL,
	cnpj varchar(30) NULL,
	error varchar(255) NULL,
	create_at timestamp NULL,
	CONSTRAINT errors_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_errors_cnpj ON public.errors USING hash (cnpj);