-- Jadwalkan fungsi expire_pending_orders untuk dijalankan setiap 5 menit
SELECT cron.schedule('expire_pending_orders_job', '*/5 * * * *', $$SELECT expire_pending_orders();$$);
