-- 주문 100,000건
INSERT INTO orders (customer_id, order_number, status, total_amount, discount_amount, shipping_fee, final_amount, shipping_address, shipping_memo, ordered_at, confirmed_at, shipped_at, delivered_at, cancelled_at)
SELECT
    1 + (floor(random() * 100000))::INT,
    'ORD-' || TO_CHAR(NOW() - (random() * INTERVAL '365 days'), 'YYYYMMDD') || '-' || LPAD(i::TEXT, 7, '0'),
    status,
    total,
    disc,
    ship_fee,
    total - disc + ship_fee,
    (ARRAY[
        '서울시 강남구 역삼대로','서울시 서초구 서초대로','서울시 마포구 월드컵북로','서울시 송파구 올림픽로',
        '경기도 성남시 분당구 불정로','경기도 수원시 영통구 광교로','경기도 고양시 일산동구 중앙로',
        '인천시 남동구 구월로','부산시 해운대구 해운대로','대구시 수성구 달구벌대로',
        '대전시 유성구 대학로','광주시 서구 상무대로','울산시 남구 삼산로',
        '제주시 연동 노형로','강원도 춘천시 중앙로'
    ])[1 + (i % 15)] || ' ' || (1 + (i % 500))::TEXT || ' ' || (100 + (i % 2000))::TEXT || '호',
    (ARRAY['부재 시 문 앞에 놓아주세요','경비실에 맡겨주세요','배송 전 연락 부탁드립니다','안전하게 배송 부탁합니다',NULL])[1 + (i % 5)],
    ord_date,
    CASE WHEN status IN ('CONFIRMED','SHIPPING','DELIVERED') THEN ord_date + INTERVAL '1 hour' ELSE NULL END,
    CASE WHEN status IN ('SHIPPING','DELIVERED') THEN ord_date + INTERVAL '1 day' ELSE NULL END,
    CASE WHEN status = 'DELIVERED' THEN ord_date + INTERVAL '2 days' + (random() * INTERVAL '3 days') ELSE NULL END,
    CASE WHEN status = 'CANCELLED' THEN ord_date + (random() * INTERVAL '1 day') ELSE NULL END
FROM generate_series(1, 100000) AS i,
LATERAL (SELECT NOW() - (random() * INTERVAL '365 days') AS ord_date) d,
LATERAL (SELECT (ARRAY['PENDING','CONFIRMED','SHIPPING','DELIVERED','DELIVERED','DELIVERED','DELIVERED','DELIVERED','CANCELLED','CANCELLED'])[1 + (floor(random() * 10))::INT] AS status) s,
LATERAL (SELECT (20000 + (floor(random() * 480000)))::INT AS total) t,
LATERAL (SELECT (floor(random() * total * 0.15))::INT AS disc) dc,
LATERAL (SELECT CASE WHEN total - disc >= 50000 THEN 0 ELSE 3000 END AS ship_fee) sf;

-- 주문 상세 (~300,000건, 주문당 1~5개)
INSERT INTO order_items (order_id, product_id, product_option_id, quantity, unit_price, subtotal, status)
SELECT
    o.id,
    100001 + (floor(random() * 100000))::INT,
    NULL,
    qty,
    price,
    qty * price,
    o.status
FROM orders o
CROSS JOIN generate_series(1, 5) AS item_num,
LATERAL (SELECT (1 + floor(random() * 4))::INT AS qty) q,
LATERAL (SELECT (5000 + floor(random() * 200000))::INT AS price) pr
WHERE item_num = 1 OR (item_num = 2 AND random() < 0.7) OR (item_num = 3 AND random() < 0.4) OR (item_num = 4 AND random() < 0.15) OR (item_num = 5 AND random() < 0.05);

-- 결제 100,000건 (주문당 1건)
INSERT INTO payments (order_id, payment_method, amount, status, pg_provider, pg_transaction_id, paid_at, cancelled_at, created_at)
SELECT
    o.id,
    (ARRAY['CARD','CARD','CARD','CARD','BANK_TRANSFER','KAKAO_PAY','NAVER_PAY','TOSS_PAY','PHONE','POINT'])[1 + (floor(random() * 10))::INT],
    o.final_amount,
    CASE o.status WHEN 'CANCELLED' THEN 'CANCELLED' WHEN 'PENDING' THEN 'PENDING' ELSE 'PAID' END,
    (ARRAY['KG이니시스','NHN KCP','토스페이먼츠','나이스페이','카카오페이'])[1 + (floor(random() * 5))::INT],
    'PG-' || o.id || '-' || (floor(random() * 1000000))::INT,
    CASE WHEN o.status != 'PENDING' THEN o.ordered_at + INTERVAL '5 minutes' ELSE NULL END,
    CASE WHEN o.status = 'CANCELLED' THEN o.cancelled_at ELSE NULL END,
    o.ordered_at
FROM orders o;

-- 배송 (CONFIRMED 이상인 주문)
INSERT INTO shipments (order_id, carrier, tracking_number, status, shipped_at, delivered_at, created_at)
SELECT
    o.id,
    (ARRAY['CJ대한통운','한진택배','롯데택배','우체국택배','로젠택배'])[1 + (floor(random() * 5))::INT],
    LPAD((floor(random() * 9999999999999))::BIGINT::TEXT, 13, '0'),
    CASE o.status
        WHEN 'CONFIRMED' THEN 'PREPARING'
        WHEN 'SHIPPING' THEN 'IN_TRANSIT'
        WHEN 'DELIVERED' THEN 'DELIVERED'
        ELSE 'PREPARING'
    END,
    o.shipped_at,
    o.delivered_at,
    COALESCE(o.confirmed_at, o.ordered_at)
FROM orders o
WHERE o.status IN ('CONFIRMED','SHIPPING','DELIVERED');
