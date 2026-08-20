CREATE SCHEMA IF NOT EXISTS public;


CREATE TABLE public.balances (
	id uuid PRIMARY KEY,
	resource_id text NOT NULL,
	chain_id text,
	amount numeric,
	erc_id numeric,
	created_at timestamptz,
	updated_at timestamptz
);
CREATE INDEX idx_balances_resource_id ON public.balances(resource_id);
CREATE INDEX idx_balances_chain_id ON public.balances(chain_id);
CREATE UNIQUE INDEX idx_balances_unique ON public.balances (chain_id, resource_id, erc_id);


CREATE TABLE public.participants (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	chain_id numeric NOT NULL,
	"name" text,
	owner_id text,
	status smallint,
	"role" smallint,
	allowed_to_broadcast boolean,
	is_flagged boolean DEFAULT false,
	flag_reason text,
	flagged_at timestamptz
);
CREATE INDEX idx_participants_chain_id ON public.participants(chain_id);


CREATE TABLE public.tokens (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	"name" text,
	symbol text,
	resource_id text UNIQUE NOT NULL,
	metadata_url text,
	erc_standard smallint,
	decimals smallint,
	issuer_id text NOT NULL,
	status smallint
);
CREATE INDEX idx_tokens_issuer_id_name ON public.tokens(issuer_id, name);


CREATE TABLE public.last_processed_block (
	id bigserial PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	"number" numeric NOT NULL
);


CREATE TABLE public.transactions (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	resource_id text NOT NULL,
	message_id text,
	from_chain_id numeric,
	to_chain_id numeric,
	amount numeric,
	"from" text,
	"to" text,
	tx_type smallint,
	msg_type smallint,
	teleport_status smallint DEFAULT NULL,
	protocol smallint DEFAULT NULL,
	is_flagged boolean,
	is_processed boolean DEFAULT false,
	hub_tx_hash text,
	hub_timestamp timestamptz,
	block_number numeric,
	log_index bigint,
	tx_hash_destination text,
	tx_hash_source text,
	source_timestamp timestamptz,
	destination_timestamp timestamptz,
	shared_id text,
	batch_id text,
	erc_id numeric,
	payload text,
	-- Aggregation fields for scalable pagination
	aggregation_type text DEFAULT 'transaction',
	aggregation_key text,
	CONSTRAINT transactions_unique UNIQUE (message_id, block_number, log_index)
);
CREATE INDEX idx_transactions_message_id ON public.transactions(message_id);
CREATE INDEX idx_transactions_resource_id ON public.transactions(resource_id);
CREATE INDEX idx_transactions_from_chain_id ON public.transactions(from_chain_id);
CREATE INDEX idx_transactions_to_chain_id ON public.transactions(to_chain_id);
CREATE INDEX idx_transactions_created_at ON public.transactions(created_at DESC);
CREATE INDEX idx_transactions_batch_id ON public.transactions(batch_id) WHERE batch_id IS NOT NULL AND batch_id != '';
CREATE INDEX idx_transactions_unprocessed ON public.transactions(created_at ASC) WHERE is_processed = false AND (teleport_status = 1 OR teleport_status = 0 OR teleport_status IS NULL);
CREATE INDEX idx_transactions_aggregation_key ON public.transactions(aggregation_key);
CREATE INDEX idx_transactions_list_pagination ON public.transactions(aggregation_key, created_at DESC) WHERE aggregation_key IS NOT NULL AND aggregation_key != '' AND tx_type NOT IN (0, 1);
CREATE INDEX idx_transactions_shared_id ON public.transactions(shared_id) WHERE shared_id IS NOT NULL AND shared_id != '';
CREATE INDEX idx_transactions_from_lower ON public.transactions (LOWER("from"));
CREATE INDEX idx_transactions_to_lower ON public.transactions (LOWER("to"));


CREATE TABLE public.enygma_transactions (
	transaction_id uuid PRIMARY KEY REFERENCES public.transactions(id) ON DELETE CASCADE,
	to_r_value_to_add numeric,
	reference_id text NOT NULL,
	updated_at timestamptz
);


CREATE TABLE public.revert_data_transactions (
	transaction_id uuid PRIMARY KEY REFERENCES public.transactions(id) ON DELETE CASCADE,
	tx_hash_destination_revert text,
	tx_hash_destination_revert_status smallint,
	tx_hash_source_revert text,
	tx_hash_source_revert_status smallint
);


CREATE TABLE public.flagged_transactions (
	id uuid PRIMARY KEY,
	created_at timestamptz,
	updated_at timestamptz,
	transaction_id uuid UNIQUE NOT NULL REFERENCES public.transactions(id) ON DELETE CASCADE
);


CREATE TABLE public.header_proof_events (
    id bigserial PRIMARY KEY,
    chain_id numeric NOT NULL,
    block_number numeric NOT NULL,
    block_hash varchar(66) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_header_proof_events_chain_block_desc ON public.header_proof_events(chain_id, block_number DESC);
CREATE UNIQUE INDEX idx_header_proof_events_chain_block ON public.header_proof_events(chain_id, block_number);


CREATE TABLE public.header_flag_events (
    id bigserial PRIMARY KEY,
    chain_id numeric NOT NULL,
    block_number numeric NOT NULL,
    reason smallint NOT NULL,
    initiator smallint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE(chain_id, block_number)
);
CREATE INDEX idx_header_flag_events_chain_id ON public.header_flag_events(chain_id);


CREATE TABLE public.private_networks (
    id bigserial PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    username varchar(255) UNIQUE NOT NULL,
    password varchar(255) NOT NULL
);


CREATE TABLE public.token_freeze_states (
	resource_id text NOT NULL,
	chain_id text NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	is_frozen boolean NOT NULL DEFAULT false,
	PRIMARY KEY (resource_id, chain_id),
	FOREIGN KEY (resource_id) REFERENCES public.tokens(resource_id) ON DELETE CASCADE
);
CREATE INDEX idx_token_freeze_states_resource_frozen
	ON public.token_freeze_states(resource_id) WHERE is_frozen = true;


CREATE TABLE public.token_freeze_audits (
	id uuid PRIMARY KEY,
	resource_id text NOT NULL,
	chain_id text NOT NULL,
	action smallint NOT NULL,
	block_number numeric NOT NULL,
	transaction_hash text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	UNIQUE(resource_id, chain_id, transaction_hash, action)
);