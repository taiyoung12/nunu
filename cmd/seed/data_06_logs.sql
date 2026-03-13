-- 검색 로그 100,000건
INSERT INTO search_logs (customer_id, keyword, result_count, clicked_product_id, searched_at)
SELECT
    CASE WHEN random() < 0.7 THEN 1 + (floor(random() * 100000))::INT ELSE NULL END,
    (ARRAY['티셔츠','원피스','운동화','노트북','아이폰','에어팟','선크림','비타민','소파','냉장고',
           '청바지','패딩','가디건','백팩','스마트워치','이어폰','세럼','오메가3','침대','TV',
           '맥북','갤럭시','나이키','아디다스','유니클로','다이슨','크림','마스크팩','요가매트','캐리어',
           '키보드','마우스','텀블러','프라이팬','강아지사료','고양이간식','기저귀','분유','색연필','기타줄',
           '남자코트','여자가방','커플운동화','아기옷','캠핑텐트','자전거','골프채','수영복','등산화','블루투스스피커'])[1 + (floor(random() * 50))::INT],
    (floor(random() * 500))::INT,
    CASE WHEN random() < 0.4 THEN 100001 + (floor(random() * 100000))::INT ELSE NULL END,
    NOW() - (random() * INTERVAL '90 days')
FROM generate_series(1, 100000);

-- 페이지뷰 100,000건
INSERT INTO page_views (customer_id, session_id, page_type, reference_id, referrer, device_type, viewed_at)
SELECT
    CASE WHEN random() < 0.6 THEN 1 + (floor(random() * 100000))::INT ELSE NULL END,
    'sess_' || (floor(random() * 500000))::INT,
    (ARRAY['HOME','PRODUCT','CATEGORY','SEARCH','CART','ORDER','MYPAGE','EVENT','BRAND'])[1 + (floor(random() * 9))::INT],
    CASE WHEN random() < 0.5 THEN (floor(random() * 100000))::INT ELSE NULL END,
    (ARRAY['https://google.com','https://naver.com','direct','https://instagram.com','https://facebook.com','https://youtube.com',NULL])[1 + (floor(random() * 7))::INT],
    (ARRAY['MOBILE','MOBILE','MOBILE','PC','PC','TABLET'])[1 + (floor(random() * 6))::INT],
    NOW() - (random() * INTERVAL '90 days')
FROM generate_series(1, 100000);

-- 로그인 이력 100,000건
INSERT INTO login_history (customer_id, ip_address, device_type, os, browser, is_success, logged_at)
SELECT
    1 + (floor(random() * 100000))::INT,
    (1 + (floor(random() * 254)))::TEXT || '.' || (floor(random() * 255))::TEXT || '.' || (floor(random() * 255))::TEXT || '.' || (1 + (floor(random() * 254)))::TEXT,
    (ARRAY['MOBILE','MOBILE','MOBILE','PC','PC','TABLET'])[1 + (floor(random() * 6))::INT],
    (ARRAY['iOS','Android','Windows','macOS','Linux'])[1 + (floor(random() * 5))::INT],
    (ARRAY['Chrome','Safari','Samsung Internet','Edge','Firefox','KakaoTalk','Naver'])[1 + (floor(random() * 7))::INT],
    CASE WHEN random() < 0.95 THEN TRUE ELSE FALSE END,
    NOW() - (random() * INTERVAL '90 days')
FROM generate_series(1, 100000);

-- 알림 100,000건
INSERT INTO customer_notifications (customer_id, type, title, content, is_read, reference_type, reference_id, created_at)
SELECT
    1 + (floor(random() * 100000))::INT,
    (ARRAY['ORDER','DELIVERY','PROMOTION','POINT','COUPON','REVIEW','SYSTEM'])[1 + (floor(random() * 7))::INT],
    CASE (i % 7)
        WHEN 0 THEN '주문이 접수되었습니다'
        WHEN 1 THEN '상품이 배송 중입니다'
        WHEN 2 THEN '특별 할인 이벤트!'
        WHEN 3 THEN '포인트가 적립되었습니다'
        WHEN 4 THEN '새 쿠폰이 발급되었습니다'
        WHEN 5 THEN '리뷰 작성 부탁드립니다'
        WHEN 6 THEN '시스템 점검 안내'
    END,
    '알림 상세 내용입니다.',
    CASE WHEN random() < 0.6 THEN TRUE ELSE FALSE END,
    (ARRAY['ORDER','PRODUCT','COUPON','POINT',NULL])[1 + (i % 5)],
    CASE WHEN random() < 0.7 THEN (floor(random() * 100000))::INT ELSE NULL END,
    NOW() - (random() * INTERVAL '90 days')
FROM generate_series(1, 100000) AS i;

-- 1:1 문의 30,000건
INSERT INTO inquiries (customer_id, category, title, content, order_id, product_id, status, created_at)
SELECT
    1 + (floor(random() * 100000))::INT,
    (ARRAY['주문/결제','배송','교환/반품','상품문의','기타'])[1 + (floor(random() * 5))::INT],
    (ARRAY['배송이 안 옵니다','환불 요청합니다','상품이 불량입니다','사이즈 교환 문의','결제 오류','포인트 문의','쿠폰 적용 안됨','주문 취소 요청','배송지 변경 요청','영수증 발급 요청'])[1 + (floor(random() * 10))::INT] || ' #' || i,
    '문의 내용입니다. 빠른 답변 부탁드립니다.',
    CASE WHEN random() < 0.5 THEN (SELECT id FROM orders ORDER BY random() LIMIT 1) ELSE NULL END,
    CASE WHEN random() < 0.3 THEN 100001 + (floor(random() * 100000))::INT ELSE NULL END,
    (ARRAY['OPEN','OPEN','IN_PROGRESS','RESOLVED','CLOSED'])[1 + (floor(random() * 5))::INT],
    NOW() - (random() * INTERVAL '180 days')
FROM generate_series(1, 30000) AS i;

-- 문의 답변 25,000건
INSERT INTO inquiry_replies (inquiry_id, author_type, content, created_at)
SELECT
    inq.id,
    CASE WHEN rep_num = 1 THEN 'STAFF' ELSE (ARRAY['STAFF','CUSTOMER'])[1 + (floor(random() * 2))::INT] END,
    CASE WHEN rep_num = 1 THEN '안녕하세요, 고객님. 문의 주셔서 감사합니다. 확인 후 안내드리겠습니다.'
         ELSE '추가 문의드립니다.' END,
    inq.created_at + (rep_num * INTERVAL '4 hours')
FROM inquiries inq
CROSS JOIN generate_series(1, 2) AS rep_num
WHERE inq.status != 'OPEN' OR rep_num = 1
ORDER BY random()
LIMIT 25000;

-- 상품 Q&A 50,000건
INSERT INTO product_qna (product_id, customer_id, question, answer, is_secret, answered_at, created_at)
SELECT
    100001 + (floor(random() * 100000))::INT,
    1 + (floor(random() * 100000))::INT,
    (ARRAY['실제 색상이 사진과 같나요?','배송은 얼마나 걸리나요?','사이즈가 정사이즈인가요?','세탁 방법이 궁금합니다','재입고 예정이 있나요?',
           '선물 포장 가능한가요?','A/S는 어떻게 되나요?','반품 가능한가요?','다른 색상도 있나요?','무게가 어느정도 되나요?'])[1 + (floor(random() * 10))::INT],
    CASE WHEN random() < 0.7 THEN '안녕하세요! 문의 감사합니다. ' ||
        (ARRAY['네, 실제 색상과 동일합니다.','보통 1~2일 내 출고됩니다.','정사이즈로 나옵니다.','세탁 라벨을 확인해주세요.','재입고 시 알림 신청해주세요.',
               '요청사항에 기재해주시면 됩니다.','구매일로부터 1년간 가능합니다.','수령 후 7일 이내 가능합니다.','상품 페이지에서 옵션 확인 부탁드립니다.','상세 페이지 하단을 참고해주세요.'])[1 + (floor(random() * 10))::INT]
    ELSE NULL END,
    CASE WHEN random() < 0.1 THEN TRUE ELSE FALSE END,
    CASE WHEN random() < 0.7 THEN NOW() - (random() * INTERVAL '180 days') ELSE NULL END,
    NOW() - (random() * INTERVAL '180 days')
FROM generate_series(1, 50000);

-- 기프트카드 5,000개
INSERT INTO gift_cards (code, initial_balance, current_balance, purchaser_id, recipient_email, status, expires_at, created_at)
SELECT
    'GC-' || LPAD(i::TEXT, 8, '0'),
    (ARRAY[10000,30000,50000,100000,200000])[1 + (floor(random() * 5))::INT],
    bal,
    1 + (floor(random() * 100000))::INT,
    'gift' || i || '@example.com',
    CASE WHEN bal > 0 THEN 'ACTIVE' WHEN random() < 0.5 THEN 'USED' ELSE 'EXPIRED' END,
    NOW() + (random() * INTERVAL '365 days'),
    NOW() - (random() * INTERVAL '365 days')
FROM generate_series(1, 5000) AS i,
LATERAL (SELECT (floor(random() * (ARRAY[10000,30000,50000,100000,200000])[1 + (floor(random() * 5))::INT]))::INT AS bal) b;

-- 기프트카드 사용 내역 10,000건
INSERT INTO gift_card_transactions (gift_card_id, order_id, amount, type, created_at)
SELECT
    gc.id,
    (SELECT id FROM orders ORDER BY random() LIMIT 1),
    CASE WHEN random() < 0.7 THEN (1000 + (floor(random() * 30000)))::INT ELSE gc.initial_balance END,
    (ARRAY['USE','USE','USE','CHARGE','REFUND'])[1 + (floor(random() * 5))::INT],
    gc.created_at + (random() * INTERVAL '180 days')
FROM gift_cards gc
CROSS JOIN generate_series(1, 3) AS s
WHERE s = 1 OR random() < 0.3
LIMIT 10000;

-- 판매자 정산 12,000건 (월별)
INSERT INTO seller_settlements (seller_id, settlement_period, total_sales, commission, shipping_subsidy, net_amount, status, settled_at, created_at)
SELECT
    s.id,
    TO_CHAR(NOW() - ((m-1) || ' months')::INTERVAL, 'YYYY-MM'),
    total_s,
    comm,
    ship_sub,
    total_s - comm + ship_sub,
    CASE WHEN m > 1 THEN 'SETTLED' ELSE 'PENDING' END,
    CASE WHEN m > 1 THEN (NOW() - ((m-1) || ' months')::INTERVAL + INTERVAL '15 days') ELSE NULL END,
    NOW() - ((m-1) || ' months')::INTERVAL
FROM sellers s
CROSS JOIN generate_series(1, 24) AS m,
LATERAL (SELECT (100000 + (floor(random() * 50000000)))::BIGINT AS total_s) ts,
LATERAL (SELECT (total_s * s.commission_rate / 100)::BIGINT AS comm) c,
LATERAL (SELECT (floor(random() * 100000))::BIGINT AS ship_sub) ss
WHERE random() < 0.8
LIMIT 12000;

-- 알림 설정 (~200,000건)
INSERT INTO notification_settings (customer_id, channel, type, is_enabled, updated_at)
SELECT
    c.id,
    channel,
    ntype,
    CASE WHEN random() < 0.7 THEN TRUE ELSE FALSE END,
    c.created_at + (random() * INTERVAL '30 days')
FROM (SELECT id, created_at FROM customers ORDER BY random() LIMIT 50000) c
CROSS JOIN (VALUES ('PUSH'),('SMS'),('EMAIL'),('KAKAO')) AS ch(channel)
CROSS JOIN (VALUES ('ORDER'),('DELIVERY'),('PROMOTION'),('POINT')) AS nt(ntype)
WHERE random() < 0.6;

-- 프로모션 대상 상품 50,000건
INSERT INTO promotion_products (promotion_id, product_id, promotion_price)
SELECT
    pr.id,
    100001 + (floor(random() * 100000))::INT,
    (floor(random() * 200000))::INT
FROM promotions pr
CROSS JOIN generate_series(1, 800) AS s
WHERE random() < 0.6
LIMIT 50000;

-- 상품 랭킹 100,000건
INSERT INTO product_rankings (product_id, ranking_type, category_id, rank_position, score, ranking_date, created_at)
SELECT
    100001 + (floor(random() * 100000))::INT,
    (ARRAY['BEST_SELLER','NEW_ARRIVAL','TOP_RATED','MOST_WISHED','HOT_DEAL'])[1 + (floor(random() * 5))::INT],
    1 + (floor(random() * 116))::INT,
    row_number() OVER (PARTITION BY ranking_date ORDER BY random()),
    (random() * 100)::NUMERIC(10,2),
    (NOW() - ((d-1) || ' days')::INTERVAL)::DATE,
    NOW() - ((d-1) || ' days')::INTERVAL
FROM generate_series(1, 30) AS d
CROSS JOIN generate_series(1, 3400) AS r
LIMIT 100000;

-- 일별 매출 통계 (365일)
INSERT INTO daily_sales (sale_date, total_orders, total_revenue, total_refunds, net_revenue, new_customers, returning_customers, avg_order_value)
SELECT
    (NOW() - ((d-1) || ' days')::INTERVAL)::DATE,
    orders_cnt,
    rev,
    ref,
    rev - ref,
    new_c,
    ret_c,
    CASE WHEN orders_cnt > 0 THEN rev / orders_cnt ELSE 0 END
FROM generate_series(1, 365) AS d,
LATERAL (SELECT (150 + (floor(random() * 300)))::INT AS orders_cnt) oc,
LATERAL (SELECT (orders_cnt * (50000 + (floor(random() * 100000))))::BIGINT AS rev) r,
LATERAL (SELECT (rev * (random() * 0.08))::BIGINT AS ref) rf,
LATERAL (SELECT (20 + (floor(random() * 100)))::INT AS new_c) nc,
LATERAL (SELECT (orders_cnt - new_c) AS ret_c) rc;

-- 카테고리별 통계 (30일 * 116카테고리)
INSERT INTO category_stats (category_id, stat_date, product_count, order_count, revenue, avg_rating, return_rate)
SELECT
    cat.id,
    (NOW() - ((d-1) || ' days')::INTERVAL)::DATE,
    (10 + (floor(random() * 500)))::INT,
    (floor(random() * 200))::INT,
    (floor(random() * 50000000))::BIGINT,
    (2.5 + (random() * 2.5))::NUMERIC(2,1),
    (random() * 5)::NUMERIC(4,2)
FROM categories cat
CROSS JOIN generate_series(1, 30) AS d;
