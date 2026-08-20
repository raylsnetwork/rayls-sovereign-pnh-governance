-- Drop all tables in reverse order
DROP TABLE IF EXISTS public.token_freeze_audits CASCADE;
DROP TABLE IF EXISTS public.token_freeze_states CASCADE;
DROP TABLE IF EXISTS public.private_networks CASCADE;
DROP TABLE IF EXISTS public.header_flag_events CASCADE;
DROP TABLE IF EXISTS public.header_proof_events CASCADE;
DROP TABLE IF EXISTS public.flagged_transactions CASCADE;
DROP TABLE IF EXISTS public.revert_data_transactions CASCADE;
DROP TABLE IF EXISTS public.enygma_transactions CASCADE;
DROP TABLE IF EXISTS public.transactions CASCADE;
DROP TABLE IF EXISTS public.last_processed_transaction CASCADE;
DROP TABLE IF EXISTS public.last_processed_block CASCADE;
DROP TABLE IF EXISTS public.tokens CASCADE;
DROP TABLE IF EXISTS public.participants CASCADE;
DROP TABLE IF EXISTS public.balances CASCADE;

DROP SCHEMA IF EXISTS public CASCADE;