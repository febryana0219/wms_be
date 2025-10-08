-- public.warehouses definition

-- Drop table

-- DROP TABLE public.warehouses;

CREATE TABLE public.warehouses (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	code varchar(50) NOT NULL,
	"name" varchar(100) NOT NULL,
	address text NULL,
	phone varchar(20) NULL,
	email varchar(100) NULL,
	manager varchar(100) NULL,
	capacity int4 DEFAULT 0 NULL,
	current_utilization int4 DEFAULT 0 NULL,
	is_active bool DEFAULT true NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT warehouses_code_key UNIQUE (code),
	CONSTRAINT warehouses_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_warehouses_code ON public.warehouses USING btree (code);
CREATE INDEX idx_warehouses_is_active ON public.warehouses USING btree (is_active);

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar(100) NOT NULL,
	email varchar(100) NOT NULL,
	password_hash varchar(255) NOT NULL,
	"role" public."user_role" DEFAULT 'staff'::user_role NULL,
	warehouse_id uuid NULL,
	is_active bool DEFAULT true NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_users_email ON public.users USING btree (email);
CREATE INDEX idx_users_role ON public.users USING btree (role);
CREATE INDEX idx_users_warehouse_id ON public.users USING btree (warehouse_id);

-- public.users foreign keys
ALTER TABLE public.users ADD CONSTRAINT users_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);

-- DROP TABLE public.refresh_tokens;

CREATE TABLE public.refresh_tokens (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	user_id uuid NULL,
	token_hash text NOT NULL,
	expires_at timestamptz NOT NULL,
	revoked_at timestamptz NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);
CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);

-- public.refresh_tokens foreign keys
ALTER TABLE public.refresh_tokens ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- DROP TABLE public.products;

CREATE TABLE public.products (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	sku varchar(100) NOT NULL,
	"name" varchar(100) NOT NULL,
	category varchar(100) NULL,
	description text NULL,
	price numeric(15, 2) DEFAULT 0.00 NOT NULL,
	min_stock int4 DEFAULT 0 NOT NULL,
	stock int4 DEFAULT 0 NOT NULL,
	reserved_stock int4 DEFAULT 0 NOT NULL,
	warehouse_id uuid NOT NULL,
	is_active bool DEFAULT true NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	available_stock int4 GENERATED ALWAYS AS (stock - reserved_stock) STORED NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_products_available_stock ON public.products USING btree (available_stock);
CREATE INDEX idx_products_category ON public.products USING btree (category);
CREATE INDEX idx_products_is_active ON public.products USING btree (is_active);
CREATE INDEX idx_products_search ON public.products USING gin (to_tsvector('english'::regconfig, (((name)::text || ' '::text) || COALESCE(description, ''::text))));
CREATE INDEX idx_products_sku ON public.products USING btree (sku);
CREATE INDEX idx_products_stock ON public.products USING btree (stock);
CREATE INDEX idx_products_warehouse_id ON public.products USING btree (warehouse_id);


-- DROP FUNCTION public.update_product_timestamp();

CREATE OR REPLACE FUNCTION public.update_product_timestamp()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$function$;

-- Table Triggers
create trigger trigger_update_product_timestamp before
update on public.products for each row execute function update_product_timestamp();


-- DROP FUNCTION public.update_warehouse_utilization();

CREATE OR REPLACE FUNCTION public.update_warehouse_utilization()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    total_stock int;
BEGIN
    -- =========================
    -- HANDLE DELETE
    -- =========================
    IF TG_OP = 'DELETE' THEN
        UPDATE warehouses
        SET current_utilization = current_utilization - OLD.stock
        WHERE id = OLD.warehouse_id;
        RETURN OLD;
    END IF;

    -- =========================
    -- HANDLE INSERT
    -- =========================
    IF TG_OP = 'INSERT' THEN
        -- Hitung total stok + new stock
        SELECT COALESCE(SUM(stock),0) + NEW.stock INTO total_stock
        FROM products
        WHERE warehouse_id = NEW.warehouse_id;

        -- Jika melebihi capacity, batalkan
        IF total_stock > (SELECT capacity FROM warehouses WHERE id = NEW.warehouse_id) THEN
            RAISE EXCEPTION 'Warehouse capacity exceeded!';
        END IF;

        -- Update current_utilization
        UPDATE warehouses
        SET current_utilization = total_stock
        WHERE id = NEW.warehouse_id;

        RETURN NEW;
    END IF;

    -- =========================
    -- HANDLE UPDATE
    -- =========================
    IF TG_OP = 'UPDATE' THEN
        -- Jika warehouse berubah, kurangi dari warehouse lama
        IF NEW.warehouse_id <> OLD.warehouse_id THEN
            UPDATE warehouses
            SET current_utilization = current_utilization - OLD.stock
            WHERE id = OLD.warehouse_id;
        END IF;

        -- Hitung total stok warehouse baru
        SELECT COALESCE(SUM(stock),0) - COALESCE(OLD.stock,0) + NEW.stock INTO total_stock
        FROM products
        WHERE warehouse_id = NEW.warehouse_id;

        -- Jika melebihi capacity, batalkan
        IF total_stock > (SELECT capacity FROM warehouses WHERE id = NEW.warehouse_id) THEN
            RAISE EXCEPTION 'Warehouse capacity exceeded!';
        END IF;

        -- Update current_utilization
        UPDATE warehouses
        SET current_utilization = total_stock
        WHERE id = NEW.warehouse_id;

        RETURN NEW;
    END IF;

    RETURN NULL; -- default
END;
$function$;

create trigger trg_update_warehouse_utilization before
insert or delete or update on public.products for each row execute function update_warehouse_utilization();

-- public.products foreign keys
ALTER TABLE public.products ADD CONSTRAINT products_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);

-- DROP TABLE public.transactions;

CREATE TABLE public.transactions (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"type" public."transaction_type" NOT NULL,
	product_id uuid NOT NULL,
	quantity int4 NOT NULL,
	warehouse_id uuid NOT NULL,
	to_warehouse_id uuid NULL,
	reference_number varchar(100) NULL,
	notes text NULL,
	created_by uuid NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT transactions_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_transactions_created_at ON public.transactions USING btree (created_at DESC);
CREATE INDEX idx_transactions_created_by ON public.transactions USING btree (created_by);
CREATE INDEX idx_transactions_product_id ON public.transactions USING btree (product_id);
CREATE INDEX idx_transactions_reference ON public.transactions USING btree (reference_number);
CREATE INDEX idx_transactions_to_warehouse_id ON public.transactions USING btree (to_warehouse_id);
CREATE INDEX idx_transactions_type ON public.transactions USING btree (type);
CREATE INDEX idx_transactions_warehouse_id ON public.transactions USING btree (warehouse_id);

-- public.transactions foreign keys
ALTER TABLE public.transactions ADD CONSTRAINT transactions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE public.transactions ADD CONSTRAINT transactions_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);
ALTER TABLE public.transactions ADD CONSTRAINT transactions_to_warehouse_id_fkey FOREIGN KEY (to_warehouse_id) REFERENCES public.warehouses(id);
ALTER TABLE public.transactions ADD CONSTRAINT transactions_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);


-- DROP FUNCTION public.transfer_stock(uuid, uuid, uuid, int4, varchar, text, uuid);

CREATE OR REPLACE FUNCTION public.transfer_stock(p_product_id uuid, p_from_warehouse_id uuid, p_to_warehouse_id uuid, p_quantity integer, p_reference_number character varying, p_notes text, p_created_by uuid)
 RETURNS boolean
 LANGUAGE plpgsql
AS $function$
DECLARE
  src_id            uuid;
  src_sku           text;
  src_stock         int;
  src_reserved      int;
  target_id         uuid;
  from_util         int;
  from_capacity     int;
  to_util           int;
  to_capacity       int;
BEGIN
  IF p_quantity <= 0 THEN
    RAISE EXCEPTION 'quantity must be > 0';
  END IF;

  -- Lock source product row (mengunci agar tidak terjadi race)
  SELECT id, sku, stock, reserved_stock
  INTO src_id, src_sku, src_stock, src_reserved
  FROM products
  WHERE id = p_product_id
    AND warehouse_id = p_from_warehouse_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'source product not found in source warehouse';
  END IF;

  -- Cek available stock di source (stock - reserved)
  IF (src_stock - COALESCE(src_reserved,0)) < p_quantity THEN
    RAISE EXCEPTION 'not enough available stock in source product';
  END IF;

  -- Lock warehouse rows (FOR UPDATE) untuk mencegah race
  SELECT current_utilization, capacity
  INTO from_util, from_capacity
  FROM warehouses
  WHERE id = p_from_warehouse_id
  FOR UPDATE;

  SELECT current_utilization, capacity
  INTO to_util, to_capacity
  FROM warehouses
  WHERE id = p_to_warehouse_id
  FOR UPDATE;

  IF from_util IS NULL THEN
    RAISE EXCEPTION 'source warehouse not found';
  END IF;
  IF to_util IS NULL THEN
    RAISE EXCEPTION 'destination warehouse not found';
  END IF;

  -- cek apakah pengurangan dari sumber tidak membuat negative utilization
  IF (from_util - p_quantity) < 0 THEN
    RAISE EXCEPTION 'transfer would make source warehouse utilization negative';
  END IF;

  -- cek apakah tujuan punya kapasitas untuk menampung tambahan
  IF (to_util + p_quantity) > to_capacity THEN
    RAISE EXCEPTION 'destination warehouse does not have enough capacity';
  END IF;

  -- Cari apakah product dengan sku yang sama sudah ada di tujuan (lock jika ada)
  SELECT id INTO target_id
  FROM products
  WHERE sku = src_sku
    AND warehouse_id = p_to_warehouse_id
  FOR UPDATE;

  IF target_id IS NULL THEN
    -- jika belum ada, copy product (masukkan initial stock = p_quantity)
    INSERT INTO products (
      name, sku, description, price, stock, warehouse_id, category, min_stock, created_at, updated_at, is_active
    )
    SELECT name, sku, description, price, p_quantity, p_to_warehouse_id, category, min_stock, now(), now(), is_active
    FROM products
    WHERE id = p_product_id
    RETURNING id INTO target_id;
  ELSE
    -- jika ada, tambahkan stock
    UPDATE products
    SET stock = stock + p_quantity,
        updated_at = now()
    WHERE id = target_id;
  END IF;

  -- Kurangi stock di source
  UPDATE products
  SET stock = stock - p_quantity,
      updated_at = now()
  WHERE id = src_id;

  RETURN TRUE;
EXCEPTION
  WHEN OTHERS THEN
    -- rethrow supaya transaksi rollback dan caller dapat error
    RAISE;
END;
$function$;


-- DROP TABLE public.inbounds;

CREATE TABLE public.inbounds (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	product_id uuid NOT NULL,
	warehouse_id uuid NOT NULL,
	quantity int4 NOT NULL,
	supplier_name varchar(100) NOT NULL,
	supplier_contact varchar(100) NULL,
	reference_number varchar(100) NULL,
	unit_cost numeric(15, 2) NULL,
	total_cost numeric(15, 2) NULL,
	notes text NULL,
	received_date timestamptz NOT NULL,
	created_by uuid NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT inbound_records_pkey PRIMARY KEY (id),
	CONSTRAINT inbound_records_quantity_check CHECK ((quantity > 0))
);
CREATE INDEX idx_inbound_records_created_by ON public.inbounds USING btree (created_by);
CREATE INDEX idx_inbound_records_product_id ON public.inbounds USING btree (product_id);
CREATE INDEX idx_inbound_records_received_date ON public.inbounds USING btree (received_date DESC);
CREATE INDEX idx_inbound_records_reference_number ON public.inbounds USING btree (reference_number);
CREATE INDEX idx_inbound_records_supplier_name ON public.inbounds USING btree (supplier_name);
CREATE INDEX idx_inbound_records_warehouse_id ON public.inbounds USING btree (warehouse_id);

-- DROP FUNCTION public.fn_update_stock_inbound();

CREATE OR REPLACE FUNCTION public.fn_update_stock_inbound()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE products
    SET stock = stock + NEW.quantity
    WHERE id = NEW.product_id;

    RETURN NEW;
END;
$function$;

-- Table Triggers
create trigger trg_update_stock_inbound after
insert on public.inbounds for each row execute function fn_update_stock_inbound();
-- public.inbounds foreign keys
ALTER TABLE public.inbounds ADD CONSTRAINT inbound_records_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE public.inbounds ADD CONSTRAINT inbound_records_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);
ALTER TABLE public.inbounds ADD CONSTRAINT inbound_records_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);

-- DROP TABLE public.outbounds;

CREATE TABLE public.outbounds (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	product_id uuid NOT NULL,
	warehouse_id uuid NOT NULL,
	quantity int4 NOT NULL,
	destination_type public."outbound_destination_type" DEFAULT 'customer'::outbound_destination_type NULL,
	destination_name varchar(100) NOT NULL,
	destination_contact varchar(100) NULL,
	reference_number varchar(100) NULL,
	unit_price numeric(15, 2) NULL,
	total_price numeric(15, 2) NULL,
	notes text NULL,
	shipped_date timestamptz NOT NULL,
	created_by uuid NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT outbound_records_pkey PRIMARY KEY (id),
	CONSTRAINT outbound_records_quantity_check CHECK ((quantity > 0))
);
CREATE INDEX idx_outbound_records_created_by ON public.outbounds USING btree (created_by);
CREATE INDEX idx_outbound_records_destination_name ON public.outbounds USING btree (destination_name);
CREATE INDEX idx_outbound_records_destination_type ON public.outbounds USING btree (destination_type);
CREATE INDEX idx_outbound_records_product_id ON public.outbounds USING btree (product_id);
CREATE INDEX idx_outbound_records_reference_number ON public.outbounds USING btree (reference_number);
CREATE INDEX idx_outbound_records_shipped_date ON public.outbounds USING btree (shipped_date DESC);
CREATE INDEX idx_outbound_records_warehouse_id ON public.outbounds USING btree (warehouse_id);

-- DROP FUNCTION public.fn_update_stock_outbound();

CREATE OR REPLACE FUNCTION public.fn_update_stock_outbound()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE products
    SET stock = stock - NEW.quantity
    WHERE id = NEW.product_id;

    RETURN NEW;
END;
$function$;

-- Table Triggers
create trigger trg_update_stock_outbound after
insert on public.outbounds for each row execute function fn_update_stock_outbound();

-- public.outbounds foreign keys
ALTER TABLE public.outbounds ADD CONSTRAINT outbound_records_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);
ALTER TABLE public.outbounds ADD CONSTRAINT outbound_records_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);
ALTER TABLE public.outbounds ADD CONSTRAINT outbound_records_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);


-- DROP TABLE public.orders;

CREATE TABLE public.orders (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	order_number varchar(50) NOT NULL,
	customer_id varchar(100) NOT NULL,
	customer_name varchar(100) NOT NULL,
	status public."order_status" DEFAULT 'pending_payment'::order_status NULL,
	total_amount numeric(15, 2) DEFAULT 0.00 NOT NULL,
	warehouse_id uuid NOT NULL,
	notes text NULL,
	expires_at timestamptz NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT orders_order_number_key UNIQUE (order_number),
	CONSTRAINT orders_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_orders_created_at ON public.orders USING btree (created_at DESC);
CREATE INDEX idx_orders_customer_id ON public.orders USING btree (customer_id);
CREATE INDEX idx_orders_expires_at ON public.orders USING btree (expires_at);
CREATE INDEX idx_orders_order_number ON public.orders USING btree (order_number);
CREATE INDEX idx_orders_pending_expiration ON public.orders USING btree (expires_at) WHERE (status = 'pending_payment'::order_status);
CREATE INDEX idx_orders_status ON public.orders USING btree (status);
CREATE INDEX idx_orders_warehouse_id ON public.orders USING btree (warehouse_id);

-- DROP FUNCTION public.update_products_on_order_status_change();

CREATE OR REPLACE FUNCTION public.update_products_on_order_status_change()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    -- Pending_payment -> Expired
     IF NEW.status = 'expired' AND OLD.status <> 'expired' THEN
        UPDATE products p
        SET reserved_stock = reserved_stock - oi.quantity
        FROM order_items oi
        WHERE oi.order_id = NEW.id AND p.id = oi.product_id;

    -- Pending_payment -> Cancelled
    ELSIF NEW.status = 'cancelled' AND OLD.status <> 'cancelled' THEN
        UPDATE products p
        SET reserved_stock = reserved_stock - oi.quantity
        FROM order_items oi
        WHERE oi.order_id = NEW.id AND p.id = oi.product_id;

    -- processing -> Shipped
    ELSIF NEW.status = 'shipped' AND OLD.status <> 'shipped' THEN
        UPDATE products p
        SET reserved_stock = reserved_stock - oi.quantity,
            stock = stock - oi.quantity
        FROM order_items oi
        WHERE oi.order_id = NEW.id AND p.id = oi.product_id;
    END IF;

    RETURN NEW;
END;
$function$;

-- Table Triggers
create trigger trg_update_products_on_order_status after
update on public.orders for each row execute function update_products_on_order_status_change();

-- public.orders foreign keys
ALTER TABLE public.orders ADD CONSTRAINT orders_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id);

-- DROP TABLE public.order_items;

CREATE TABLE public.order_items (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	order_id uuid NOT NULL,
	product_id uuid NOT NULL,
	quantity int4 NOT NULL,
	unit_price numeric(15, 2) NOT NULL,
	total_price numeric(15, 2) NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT order_items_pkey PRIMARY KEY (id),
	CONSTRAINT order_items_quantity_check CHECK ((quantity > 0))
);
CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);
CREATE INDEX idx_order_items_product_id ON public.order_items USING btree (product_id);

-- DROP FUNCTION public.calculate_order_total();

CREATE OR REPLACE FUNCTION public.calculate_order_total()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    -- Update order total when order items change
    UPDATE orders 
    SET total_amount = (
        SELECT COALESCE(SUM(total_price), 0)
        FROM order_items 
        WHERE order_id = COALESCE(NEW.order_id, OLD.order_id)
    )
    WHERE id = COALESCE(NEW.order_id, OLD.order_id);
    
    RETURN COALESCE(NEW, OLD);
END;
$function$;

-- Table Triggers
create trigger trigger_calculate_order_total_insert after
insert on public.order_items for each row execute function calculate_order_total();

create trigger trigger_calculate_order_total_update after
update on public.order_items for each row execute function calculate_order_total();

create trigger trigger_calculate_order_total_delete after
delete on public.order_items for each row execute function calculate_order_total();

-- DROP FUNCTION public.set_order_item_prices();

CREATE OR REPLACE FUNCTION public.set_order_item_prices()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    product_price numeric(15,2);
BEGIN
    -- Ambil harga produk dari tabel products
    SELECT price INTO product_price
    FROM public.products
    WHERE id = NEW.product_id;

    IF product_price IS NULL THEN
        RAISE EXCEPTION 'Product with id % not found', NEW.product_id;
    END IF;

    -- Set unit_price & total_price
    NEW.unit_price := product_price;
    NEW.total_price := product_price * NEW.quantity;

    RETURN NEW;
END;
$function$;

create trigger trigger_set_order_item_prices before
insert on public.order_items for each row execute function set_order_item_prices();

-- DROP FUNCTION public.update_reserved_stock_on_insert();

CREATE OR REPLACE FUNCTION public.update_reserved_stock_on_insert()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF EXISTS (SELECT 1 FROM orders WHERE id = NEW.order_id AND status = 'pending_payment') THEN
        UPDATE products
        SET reserved_stock = reserved_stock + NEW.quantity
        WHERE id = NEW.product_id;
    END IF;

    RETURN NEW;
END;
$function$;

create trigger trg_update_reserved_stock_on_insert after
insert on public.order_items for each row execute function update_reserved_stock_on_insert();

-- public.order_items foreign keys
ALTER TABLE public.order_items ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;
ALTER TABLE public.order_items ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);