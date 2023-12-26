-- public.rfb_data definition

-- Drop table

-- DROP TABLE public.rfb_data;

CREATE TABLE public.rfb_data (
	id serial4 NOT NULL,
	cnpj varchar(16) NULL,
	cnpj_basico varchar(8) NULL,
	cnpj_ordem varchar(4) NULL,
	cnpj_dv varchar(2) NULL,
	identificador varchar(2) NULL,
	nome_fantasia varchar(250) NULL,
	situacao_cadastral varchar(2) NULL,
	data_situacao_cadastral date NULL,
	motivo_situacao_cadastral varchar(2) NULL,
	nome_cidade_exterior varchar(250) NULL,
	pais varchar(100) NULL,
	data_inicio date NULL,
	cnae_principal varchar(7) NULL,
	cnae_secundario varchar(3000) NULL,
	tipo_logradouro varchar(25) NULL,
	logradouro varchar(250) NULL,
	numero varchar(8) NULL,
	complemento varchar(100) NULL,
	bairro varchar(100) NULL,
	cep varchar(8) NULL,
	uf varchar(2) NULL,
	municipio varchar(100) NULL,
	ddd_1 varchar(3) NULL,
	telefone_1 varchar(9) NULL,
	ddd_2 varchar(3) NULL,
	telefone_2 varchar(9) NULL,
	ddd_fax varchar(3) NULL,
	fax varchar(9) NULL,
	email varchar(100) NULL,
	situacao_especial varchar(250) NULL,
	data_situacao_especial date NULL,
	create_at timestamp NULL,
	update_at timestamp NULL,
	hash varchar(64) NOT NULL,
	CONSTRAINT data_hash UNIQUE (hash),
	CONSTRAINT data_cnpj_key UNIQUE (cnpj),
	CONSTRAINT data_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_rfb_data_id ON public.rfb_data USING hash (id);
CREATE INDEX idx_rfb_data_cnpj ON public.rfb_data USING hash (cnpj);
CREATE INDEX idx_rfb_data_hash ON public.rfb_data USING hash (hash);