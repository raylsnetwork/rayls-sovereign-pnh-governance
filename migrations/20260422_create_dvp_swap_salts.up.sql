CREATE TABLE public.dvp_swap_salts (
    shared_id            TEXT PRIMARY KEY,
    initiator_self_salt  BYTEA NOT NULL,
    initiator_ctxt_salt  BYTEA NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
