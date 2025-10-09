
CREATE VIEW public.dashboard_summary AS
SELECT 
  (SELECT count(*) AS count FROM public.products) AS total_products,
  (SELECT count(*) AS count FROM public.warehouses) AS total_warehouses,
  (SELECT count(*) AS count FROM public.orders) AS total_orders,
  (SELECT count(*) AS count FROM public.transactions) AS total_transactions,
  (SELECT count(*) AS count FROM public.warehouses WHERE (warehouses.is_active = true)) AS active_warehouses,
  (SELECT count(*) AS count FROM public.orders WHERE (orders.status = 'pending_payment'::public.order_status)) AS pending_orders;

CREATE VIEW public.low_stock_products AS
SELECT
  p.id,
  p.sku,
  p.name,
  p.category,
  p.description,
  p.price,
  p.min_stock,
  p.stock,
  p.reserved_stock,
  p.warehouse_id,
  p.is_active,
  p.created_at,
  p.updated_at,
  p.available_stock,
  w.name AS warehouse_name,
  (p.min_stock - p.available_stock) AS shortage_quantity
FROM (public.products p
JOIN public.warehouses w ON ((p.warehouse_id = w.id)))
WHERE (
  (p.available_stock <= p.min_stock) 
  AND (p.is_active = true)
  AND (w.is_active = true)
)
ORDER BY p.available_stock;

CREATE VIEW public.transaction_histories AS
SELECT
  t.id,
  t.type,
  t.product_id,
  p.name AS product_name,
  p.sku,
  t.quantity,
  t.warehouse_id,
  w1.name AS warehouse_name,
  t.to_warehouse_id,
  w2.name AS to_warehouse_name,
  t.reference_number,
  t.notes,
  t.created_by,
  u.name AS created_by_name,
  t.created_at
FROM ((((public.transactions t
JOIN public.products p ON ((t.product_id = p.id)))
JOIN public.warehouses w1 ON ((t.warehouse_id = w1.id)))
LEFT JOIN public.warehouses w2 ON ((t.to_warehouse_id = w2.id)))
JOIN public.users u ON ((t.created_by = u.id)));