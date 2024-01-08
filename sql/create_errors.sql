-- public.errors definition

-- Drop table

-- DROP TABLE public.errors;

CREATE TABLE public.errors (
	id serial4 NOT NULL,
	cnpj varchar(30) NOT NULL,
	field varchar(255) NULL,
	value varchar(5000) NULL,
	tag varchar(50) NULL,
	file varchar(50) NULL,
	create_at timestamp NULL,
	CONSTRAINT errors_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_errors_cnpj ON public.errors USING hash (cnpj);