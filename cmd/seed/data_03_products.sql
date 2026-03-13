-- 상품 100,000개
INSERT INTO products (seller_id, brand_id, category_id, name, description, price, cost_price, stock, status, is_featured, weight_gram, rating_avg, review_count, sales_count, created_at)
SELECT
    1 + (floor(random() * 500))::INT,
    1 + (floor(random() * 200))::INT,
    -- 카테고리 (id 1~116 범위)
    1 + (i % 116),
    (ARRAY['프리미엄','클래식','베이직','슬림핏','오버핏','빈티지','모던','럭셔리','에코','내추럴',
           '울트라','프로','미니','맥스','라이트','하이퍼','네오','퓨어','소프트','리얼'])[1 + (i % 20)]
    || ' ' ||
    (ARRAY['티셔츠','청바지','운동화','백팩','스마트폰','노트북','이어폰','텀블러','선크림','비타민',
           '소파','책상','프라이팬','요가매트','만년필','기타','캐리어','목걸이','손목시계','충전기',
           '패딩','셔츠','로퍼','크로스백','태블릿','TV','스피커','도마','토너','오메가3',
           '침대','의자','밀폐용기','덤벨','다이어리','피아노','방향제','반지','체중계','케이블',
           '가디건','코트','슬리퍼','에코백','모니터','냉장고','헤드셋','식기세트','립스틱','콜라겐',
           '바지','원피스','스니커즈','지갑','워치','에어팟','마우스','보틀','세럼','유산균',
           '매트리스','서랍장','칼세트','폼롤러','스케치북','드럼패드','블랙박스','귀걸이','안마기','허브'])[1 + (i % 70)]
    || ' ' || (ARRAY['A','B','C','S','X','Z','Pro','Max','Lite','SE','Plus','Ultra','Mini','Air','Neo'])[1 + (i % 15)]
    || '-' || (10000 + i)::TEXT,
    '상품 설명입니다. SKU: SKU-' || LPAD(i::TEXT, 7, '0'),
    CASE
        WHEN (i % 116) BETWEEN 6 AND 11 THEN (200000 + (random() * 1800000))::INT
        WHEN (i % 116) BETWEEN 22 AND 27 THEN (100000 + (random() * 900000))::INT
        WHEN (i % 116) BETWEEN 12 AND 16 THEN (1000 + (random() * 49000))::INT
        WHEN (i % 116) BETWEEN 34 AND 38 THEN (5000 + (random() * 25000))::INT
        ELSE (10000 + (random() * 290000))::INT
    END,
    NULL,
    (floor(random() * 1000))::INT,
    CASE WHEN random() < 0.85 THEN 'ACTIVE' WHEN random() < 0.95 THEN 'SOLDOUT' ELSE 'HIDDEN' END,
    CASE WHEN random() < 0.1 THEN TRUE ELSE FALSE END,
    (100 + (floor(random() * 9900)))::INT,
    (1 + (random() * 4))::NUMERIC(2,1),
    0, 0,
    NOW() - (random() * INTERVAL '730 days')
FROM generate_series(1, 100000) AS i;

-- 상품 옵션 (~300,000건, 상품당 1~5개)
INSERT INTO product_options (product_id, option_type, option_value, additional_price, stock, is_active)
SELECT
    p.id,
    CASE opt_num
        WHEN 1 THEN '사이즈'
        WHEN 2 THEN '색상'
        WHEN 3 THEN '소재'
        ELSE '기타'
    END,
    CASE opt_num
        WHEN 1 THEN (ARRAY['XS','S','M','L','XL','XXL','FREE'])[1 + ((p.id + opt_num) % 7)]
        WHEN 2 THEN (ARRAY['블랙','화이트','네이비','그레이','베이지','레드','블루','카키','브라운','핑크'])[1 + ((p.id + opt_num) % 10)]
        WHEN 3 THEN (ARRAY['면','폴리','린넨','울','캐시미어'])[1 + (p.id % 5)]
        ELSE (ARRAY['옵션A','옵션B','옵션C'])[1 + (p.id % 3)]
    END,
    CASE WHEN opt_num > 1 AND random() < 0.3 THEN (1000 * (1 + (floor(random() * 10))::INT)) ELSE 0 END,
    (floor(random() * 200))::INT,
    CASE WHEN random() < 0.9 THEN TRUE ELSE FALSE END
FROM products p
CROSS JOIN generate_series(1, 4) AS opt_num
WHERE opt_num <= 2 OR random() < 0.4;

-- 상품 이미지 (~300,000건, 상품당 1~5개)
INSERT INTO product_images (product_id, image_url, sort_order, is_main, created_at)
SELECT
    p.id,
    'https://cdn.example.com/products/' || p.id || '/img_' || img_num || '.jpg',
    img_num,
    CASE WHEN img_num = 1 THEN TRUE ELSE FALSE END,
    p.created_at
FROM products p
CROSS JOIN generate_series(1, 5) AS img_num
WHERE img_num <= 2 OR random() < 0.5;

-- 상품 태그 (~200,000건)
INSERT INTO product_tags (product_id, tag)
SELECT
    p.id,
    (ARRAY['인기','신상','할인','베스트','추천','한정','시즌','이벤트','무료배송','당일출고',
           '친환경','프리미엄','가성비','국내제작','수입','오가닉','핸드메이드','한정판','콜라보','리뉴얼'])[1 + ((p.id + tag_num * 3) % 20)]
FROM products p
CROSS JOIN generate_series(1, 3) AS tag_num
WHERE tag_num = 1 OR random() < 0.5;

-- 묶음 상품 1,000개
INSERT INTO product_bundles (name, description, bundle_price, original_price, status, created_at)
SELECT
    '알뜰 세트 #' || i,
    '인기 상품 묶음 할인',
    (30000 + (floor(random() * 200000)))::INT,
    (50000 + (floor(random() * 300000)))::INT,
    CASE WHEN random() < 0.8 THEN 'ACTIVE' ELSE 'INACTIVE' END,
    NOW() - (random() * INTERVAL '180 days')
FROM generate_series(1, 1000) AS i;

-- 묶음 구성 상품
INSERT INTO bundle_items (bundle_id, product_id, quantity)
SELECT
    b.id,
    1 + (floor(random() * 100000))::INT,
    1 + (floor(random() * 3))::INT
FROM product_bundles b
CROSS JOIN generate_series(1, 4) AS s
WHERE s <= 2 OR random() < 0.5;

-- 재고 (창고별, ~200,000건)
INSERT INTO inventory (warehouse_id, product_id, quantity, reserved, updated_at)
SELECT
    w_id,
    p.id,
    (floor(random() * 500))::INT,
    (floor(random() * 50))::INT,
    NOW() - (random() * INTERVAL '30 days')
FROM products p
CROSS JOIN generate_series(1, 3) AS w_idx
CROSS JOIN LATERAL (SELECT 1 + (floor(random() * 10))::INT AS w_id) wh
WHERE w_idx = 1 OR random() < 0.3;
