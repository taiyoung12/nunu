-- 배송 추적 (~300,000건)
INSERT INTO shipment_tracking (shipment_id, status, location, description, tracked_at)
SELECT
    s.id,
    status_val,
    loc,
    desc_val,
    s.created_at + (step_num * INTERVAL '8 hours')
FROM shipments s
CROSS JOIN generate_series(1, 4) AS step_num,
LATERAL (SELECT
    CASE step_num
        WHEN 1 THEN '집하'
        WHEN 2 THEN '간선상차'
        WHEN 3 THEN '배달출발'
        WHEN 4 THEN '배달완료'
    END AS status_val,
    CASE step_num
        WHEN 1 THEN (ARRAY['서울 강남 물류센터','경기 용인 허브','부산 강서 물류센터','대전 유성 물류센터','인천 중구 허브'])[1 + (s.id % 5)]
        WHEN 2 THEN (ARRAY['수도권 메가허브','중부 간선터미널','영남 간선터미널','호남 간선터미널'])[1 + (s.id % 4)]
        WHEN 3 THEN (ARRAY['강남구 배달점','마포구 배달점','분당구 배달점','해운대구 배달점','수성구 배달점'])[1 + (s.id % 5)]
        WHEN 4 THEN '수령인 전달'
    END AS loc,
    CASE step_num
        WHEN 1 THEN '물품을 인수하였습니다'
        WHEN 2 THEN '간선 운송 중입니다'
        WHEN 3 THEN '배달을 시작합니다'
        WHEN 4 THEN '배달이 완료되었습니다'
    END AS desc_val
) info
WHERE (s.status = 'DELIVERED' AND step_num <= 4)
   OR (s.status = 'IN_TRANSIT' AND step_num <= 3)
   OR (s.status = 'PREPARING' AND step_num <= 1);

-- 리뷰 100,000건
INSERT INTO reviews (product_id, customer_id, order_item_id, rating, content, image_urls, is_verified, helpful_count, created_at)
SELECT
    100001 + (floor(random() * 100000))::INT,
    1 + (floor(random() * 100000))::INT,
    NULL,
    rating,
    CASE rating
        WHEN 1 THEN (ARRAY['최악입니다. 환불 원합니다.','기대 이하네요. 실망입니다.','사진과 너무 달라요.','불량품이 왔어요. 절대 비추.','배송도 늦고 품질도 별로입니다.','돈이 아깝습니다.','다시는 안 삽니다.','포장도 엉망이고 상품도 엉망.'])[ 1 + (floor(random()*8))::INT]
        WHEN 2 THEN (ARRAY['별로예요. 가격 대비 아쉽습니다.','그냥 그래요. 기대 이하.','디자인은 괜찮은데 품질이...','사이즈가 안 맞아요.','색상이 사진과 달라요.','재질이 좀 별로네요.','이 가격에 이 퀄리티는 아쉽습니다.'])[ 1 + (floor(random()*7))::INT]
        WHEN 3 THEN (ARRAY['무난합니다.','가격 대비 괜찮아요.','그럭저럭 쓸만해요.','나쁘진 않아요. 보통.','가격 생각하면 OK.','디자인은 좋은데 내구성이 걱정.','선물용으로 무난합니다.','평균은 합니다.'])[ 1 + (floor(random()*8))::INT]
        WHEN 4 THEN (ARRAY['좋아요! 만족합니다.','가성비 최고입니다.','품질이 좋네요. 재구매 의사 있어요.','배송도 빠르고 제품도 좋아요.','친구한테도 추천했어요.','기대 이상이었습니다.','포장도 꼼꼼하고 만족해요.','퀄리티 대비 가격이 착해요.','두 번째 구매인데 역시 좋습니다.'])[ 1 + (floor(random()*9))::INT]
        WHEN 5 THEN (ARRAY['완벽합니다! 강력 추천!','인생템 찾았습니다!','세 번째 재구매 중이에요.','선물로도 최고입니다.','품질 대박! 가격도 착해요.','이건 진짜 사야 합니다.','역대급 만족도!','와 진짜 대박이에요.','모든 면에서 완벽합니다.','다른 색상도 추가 구매할 예정!'])[ 1 + (floor(random()*10))::INT]
    END,
    CASE WHEN random() < 0.3 THEN '["https://cdn.example.com/reviews/' || i || '_1.jpg"]' ELSE NULL END,
    CASE WHEN random() < 0.7 THEN TRUE ELSE FALSE END,
    (floor(random() * 50))::INT,
    NOW() - (random() * INTERVAL '365 days')
FROM generate_series(1, 100000) AS i,
LATERAL (SELECT (ARRAY[2,3,3,4,4,4,4,5,5,5])[1 + (floor(random()*10))::INT] AS rating) r;

-- 위시리스트 100,000건
INSERT INTO wishlists (customer_id, product_id, created_at)
SELECT DISTINCT ON (cid, pid)
    cid, pid, NOW() - (random() * INTERVAL '365 days')
FROM (
    SELECT
        1 + (floor(random() * 100000))::INT AS cid,
        100001 + (floor(random() * 100000))::INT AS pid
    FROM generate_series(1, 110000)
) sub
LIMIT 100000;

-- 장바구니 50,000건
INSERT INTO cart_items (customer_id, product_id, product_option_id, quantity, added_at)
SELECT
    1 + (floor(random() * 100000))::INT,
    100001 + (floor(random() * 100000))::INT,
    NULL,
    1 + (floor(random() * 4))::INT,
    NOW() - (random() * INTERVAL '30 days')
FROM generate_series(1, 50000);

-- 환불 10,000건
INSERT INTO refunds (order_id, payment_id, reason, amount, status, requested_at, completed_at)
SELECT
    o.id,
    p.id,
    (ARRAY['단순변심','상품불량','오배송','사이즈 안맞음','상품 파손','다른 상품 수령','색상 상이','배송 지연'])[1 + (floor(random() * 8))::INT],
    (o.final_amount * (0.3 + random() * 0.7))::INT,
    CASE WHEN random() < 0.7 THEN 'COMPLETED' WHEN random() < 0.5 THEN 'PROCESSING' ELSE 'REQUESTED' END,
    o.ordered_at + (random() * INTERVAL '7 days'),
    CASE WHEN random() < 0.7 THEN o.ordered_at + (random() * INTERVAL '14 days') ELSE NULL END
FROM orders o
JOIN payments p ON p.order_id = o.id
WHERE o.status IN ('DELIVERED','CANCELLED')
ORDER BY random()
LIMIT 10000;

-- 반품 8,000건
INSERT INTO returns (order_id, order_item_id, customer_id, reason, reason_detail, status, requested_at, completed_at)
SELECT
    o.id,
    NULL,
    o.customer_id,
    (ARRAY['DEFECT','WRONG_ITEM','CHANGE_MIND','SIZE_ISSUE','DAMAGED','COLOR_DIFF'])[1 + (floor(random() * 6))::INT],
    '반품 사유 상세 내용입니다.',
    (ARRAY['REQUESTED','COLLECTING','INSPECTING','COMPLETED','REJECTED'])[1 + (floor(random() * 5))::INT],
    o.delivered_at + (random() * INTERVAL '7 days'),
    CASE WHEN random() < 0.6 THEN o.delivered_at + (random() * INTERVAL '14 days') ELSE NULL END
FROM orders o
WHERE o.status = 'DELIVERED' AND o.delivered_at IS NOT NULL
ORDER BY random()
LIMIT 8000;

-- 교환 5,000건
INSERT INTO exchanges (order_id, order_item_id, customer_id, reason, new_product_id, status, requested_at, completed_at)
SELECT
    o.id,
    NULL,
    o.customer_id,
    (ARRAY['SIZE_CHANGE','COLOR_CHANGE','DEFECT','WRONG_ITEM'])[1 + (floor(random() * 4))::INT],
    100001 + (floor(random() * 100000))::INT,
    (ARRAY['REQUESTED','COLLECTING','SHIPPING_NEW','COMPLETED'])[1 + (floor(random() * 4))::INT],
    o.delivered_at + (random() * INTERVAL '7 days'),
    CASE WHEN random() < 0.5 THEN o.delivered_at + (random() * INTERVAL '14 days') ELSE NULL END
FROM orders o
WHERE o.status = 'DELIVERED' AND o.delivered_at IS NOT NULL
ORDER BY random()
LIMIT 5000;

-- 쿠폰 사용 100,000건
INSERT INTO coupon_usage (coupon_id, customer_id, order_id, discount_amount, used_at)
SELECT
    1 + (floor(random() * 200))::INT,
    o.customer_id,
    o.id,
    o.discount_amount,
    o.ordered_at
FROM orders o
WHERE o.discount_amount > 0
ORDER BY random()
LIMIT 100000;

-- 포인트 내역 200,000건
INSERT INTO point_history (customer_id, order_id, type, amount, balance_after, description, created_at)
SELECT
    o.customer_id,
    o.id,
    'EARN',
    (o.final_amount * 0.01)::INT,
    (floor(random() * 50000))::INT,
    '주문 적립',
    o.ordered_at + INTERVAL '3 days'
FROM orders o
WHERE o.status = 'DELIVERED'
UNION ALL
SELECT
    1 + (floor(random() * 100000))::INT,
    NULL,
    (ARRAY['USE','EARN','EXPIRE','ADMIN'])[1 + (floor(random() * 4))::INT],
    CASE WHEN random() < 0.5 THEN (100 + (floor(random() * 5000)))::INT ELSE -(100 + (floor(random() * 3000)))::INT END,
    (floor(random() * 30000))::INT,
    (ARRAY['포인트 사용','이벤트 적립','포인트 만료','관리자 지급','리뷰 적립','생일 적립'])[1 + (floor(random() * 6))::INT],
    NOW() - (random() * INTERVAL '365 days')
FROM generate_series(1, 100000);
