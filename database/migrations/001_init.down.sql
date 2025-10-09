-- Hapus trigger dan function dalam urutan aman
DROP TRIGGER IF EXISTS trg_update_reserved_stock_on_insert ON public.order_items;
DROP TRIGGER IF EXISTS trigger_set_order_item_prices ON public.order_items;
DROP TRIGGER IF EXISTS trigger_calculate_order_total_insert ON public.order_items;
DROP TRIGGER IF EXISTS trigger_calculate_order_total_update ON public.order_items;
DROP TRIGGER IF EXISTS trigger_calculate_order_total_delete ON public.order_items;
DROP TRIGGER IF EXISTS trg_update_products_on_order_status ON public.orders;
DROP TRIGGER IF EXISTS trg_update_stock_outbound ON public.outbounds;
DROP TRIGGER IF EXISTS trg_update_stock_inbound ON public.inbounds;
DROP TRIGGER IF EXISTS trg_update_warehouse_utilization ON public.products;
DROP TRIGGER IF EXISTS trigger_update_product_timestamp ON public.products;

-- Drop functions
DROP FUNCTION IF EXISTS public.update_reserved_stock_on_insert;
DROP FUNCTION IF EXISTS public.set_order_item_prices;
DROP FUNCTION IF EXISTS public.calculate_order_total;
DROP FUNCTION IF EXISTS public.update_products_on_order_status_change;
DROP FUNCTION IF EXISTS public.fn_update_stock_outbound;
DROP FUNCTION IF EXISTS public.fn_update_stock_inbound;
DROP FUNCTION IF EXISTS public.transfer_stock;
DROP FUNCTION IF EXISTS public.update_warehouse_utilization;
DROP FUNCTION IF EXISTS public.update_product_timestamp;

-- Drop foreign keys (optional, hanya jika ingin bersih total)
ALTER TABLE IF EXISTS public.order_items DROP CONSTRAINT IF EXISTS order_items_order_id_fkey;
ALTER TABLE IF EXISTS public.order_items DROP CONSTRAINT IF EXISTS order_items_product_id_fkey;

ALTER TABLE IF EXISTS public.orders DROP CONSTRAINT IF EXISTS orders_warehouse_id_fkey;

ALTER TABLE IF EXISTS public.outbounds DROP CONSTRAINT IF EXISTS outbound_records_created_by_fkey;
ALTER TABLE IF EXISTS public.outbounds DROP CONSTRAINT IF EXISTS outbound_records_product_id_fkey;
ALTER TABLE IF EXISTS public.outbounds DROP CONSTRAINT IF EXISTS outbound_records_warehouse_id_fkey;

ALTER TABLE IF EXISTS public.inbounds DROP CONSTRAINT IF EXISTS inbound_records_created_by_fkey;
ALTER TABLE IF EXISTS public.inbounds DROP CONSTRAINT IF EXISTS inbound_records_product_id_fkey;
ALTER TABLE IF EXISTS public.inbounds DROP CONSTRAINT IF EXISTS inbound_records_warehouse_id_fkey;

ALTER TABLE IF EXISTS public.transactions DROP CONSTRAINT IF EXISTS transactions_created_by_fkey;
ALTER TABLE IF EXISTS public.transactions DROP CONSTRAINT IF EXISTS transactions_product_id_fkey;
ALTER TABLE IF EXISTS public.transactions DROP CONSTRAINT IF EXISTS transactions_to_warehouse_id_fkey;
ALTER TABLE IF EXISTS public.transactions DROP CONSTRAINT IF EXISTS transactions_warehouse_id_fkey;

ALTER TABLE IF EXISTS public.products DROP CONSTRAINT IF EXISTS products_warehouse_id_fkey;
ALTER TABLE IF EXISTS public.refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey;
ALTER TABLE IF EXISTS public.users DROP CONSTRAINT IF EXISTS users_warehouse_id_fkey;

-- Drop tables
DROP TABLE IF EXISTS public.order_items;
DROP TABLE IF EXISTS public.orders;
DROP TABLE IF EXISTS public.outbounds;
DROP TABLE IF EXISTS public.inbounds;
DROP TABLE IF EXISTS public.transactions;
DROP TABLE IF EXISTS public.products;
DROP TABLE IF EXISTS public.refresh_tokens;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.warehouses;

-- Drop enum types
DROP TYPE IF EXISTS public.order_status;
DROP TYPE IF EXISTS public.outbound_destination_type;
DROP TYPE IF EXISTS public.transaction_type;
DROP TYPE IF EXISTS public.user_role;
