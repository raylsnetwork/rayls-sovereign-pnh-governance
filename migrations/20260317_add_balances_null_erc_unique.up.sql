-- The existing idx_balances_unique index (chain_id, resource_id, erc_id) uses a B-tree
-- where NULLs are treated as distinct, allowing duplicate (chain_id, resource_id, NULL) rows.
-- This partial index closes the gap for balances without an erc_id.
CREATE UNIQUE INDEX IF NOT EXISTS idx_balances_unique_null_erc
    ON public.balances (chain_id, resource_id)
    WHERE erc_id IS NULL;
