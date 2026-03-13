-- ============================================================
-- 기준 데이터 (membership_tiers, brands, sellers, categories, warehouses, shipping_policies, delivery_areas, tax_rates, faq, notices, banners, coupons, promotions)
-- ============================================================

-- 멤버십 등급
INSERT INTO membership_tiers (name, min_spending, point_rate, benefits) VALUES
('BRONZE', 0, 0.01, '기본 적립 1%'),
('SILVER', 300000, 0.02, '적립 2%, 무료배송 월 2회'),
('GOLD', 1000000, 0.03, '적립 3%, 무료배송 무제한, 생일쿠폰'),
('DIAMOND', 3000000, 0.05, '적립 5%, 무료배송 무제한, 전용 CS, VIP 라운지');

-- 브랜드 200개
INSERT INTO brands (name, country, description, is_active)
SELECT
    (ARRAY['나이키','아디다스','뉴발란스','컨버스','반스','퓨마','리복','아식스','언더아머','노스페이스',
           '파타고니아','유니클로','자라','H&M','무인양품','이케아','다이슨','삼성','LG','애플',
           '소니','보스','JBL','샤오미','화웨이','구찌','프라다','루이비통','샤넬','에르메스',
           '설화수','이니스프리','라네즈','에뛰드','미샤','닥터자르트','코스알엑스','마녀공장','롬앤','클리오',
           'CJ','풀무원','오뚜기','농심','삼양','동원','매일유업','남양유업','빙그레','롯데제과'])[1 + (i % 50)]
    || CASE WHEN i > 50 THEN ' ' || (ARRAY['코리아','글로벌','프리미엄','스탠다드','홈','키즈','스포츠','뷰티','테크','라이프'])[1 + (i % 10)] ELSE '' END,
    (ARRAY['한국','미국','일본','독일','프랑스','이탈리아','스웨덴','중국','영국','스페인'])[1 + (i % 10)],
    '우수 브랜드',
    CASE WHEN random() < 0.95 THEN TRUE ELSE FALSE END
FROM generate_series(1, 200) AS i;

-- 판매자 500개
INSERT INTO sellers (name, business_number, representative, email, phone, address, commission_rate, status, joined_at)
SELECT
    (ARRAY['스타','블루','그린','레드','골드','실버','퍼스트','탑','베스트','프라임',
           '스마트','해피','럭키','드림','판다','이글','라이온','폭스','베어','호크'])[1 + (i % 20)]
    || (ARRAY['마켓','스토어','몰','샵','트레이딩','커머스','리테일','딜','플러스','허브'])[1 + ((i/20) % 10)]
    || CASE WHEN i > 200 THEN '-' || i::TEXT ELSE '' END,
    LPAD((100 + (i % 900))::TEXT, 3, '0') || '-' || LPAD((10 + (i % 90))::TEXT, 2, '0') || '-' || LPAD((10000 + i)::TEXT, 5, '0'),
    '대표' || i,
    'seller' || i || '@shop.com',
    '02-' || LPAD((1000 + (i % 9000))::TEXT, 4, '0') || '-' || LPAD((1000 + ((i*3) % 9000))::TEXT, 4, '0'),
    (ARRAY['서울시 강남구','서울시 마포구','경기도 성남시','부산시 해운대구','대구시 수성구'])[1 + (i % 5)],
    (8 + (random() * 12))::NUMERIC(4,2),
    CASE WHEN random() < 0.9 THEN 'ACTIVE' WHEN random() < 0.5 THEN 'SUSPENDED' ELSE 'INACTIVE' END,
    NOW() - (random() * INTERVAL '1095 days')
FROM generate_series(1, 500) AS i;

-- 카테고리 (대/중/소 3단계, ~150개)
INSERT INTO categories (parent_id, name, depth, sort_order, is_active) VALUES
(NULL,'패션',0,1,TRUE),(NULL,'디지털/가전',0,2,TRUE),(NULL,'식품/음료',0,3,TRUE),
(NULL,'뷰티/건강',0,4,TRUE),(NULL,'홈/리빙',0,5,TRUE),(NULL,'스포츠/레저',0,6,TRUE),
(NULL,'도서/문구',0,7,TRUE),(NULL,'키즈/반려',0,8,TRUE),(NULL,'자동차/공구',0,9,TRUE),(NULL,'여행/티켓',0,10,TRUE);

INSERT INTO categories (parent_id, name, depth, sort_order, is_active) VALUES
(1,'남성의류',1,1,TRUE),(1,'여성의류',1,2,TRUE),(1,'신발',1,3,TRUE),(1,'가방/지갑',1,4,TRUE),(1,'쥬얼리/시계',1,5,TRUE),(1,'언더웨어',1,6,TRUE),
(2,'스마트폰/태블릿',1,1,TRUE),(2,'노트북/PC',1,2,TRUE),(2,'TV/영상',1,3,TRUE),(2,'음향기기',1,4,TRUE),(2,'생활가전',1,5,TRUE),(2,'주방가전',1,6,TRUE),
(3,'신선식품',1,1,TRUE),(3,'가공식품',1,2,TRUE),(3,'음료/커피',1,3,TRUE),(3,'건강식품',1,4,TRUE),(3,'간식/제과',1,5,TRUE),
(4,'스킨케어',1,1,TRUE),(4,'메이크업',1,2,TRUE),(4,'향수',1,3,TRUE),(4,'헤어케어',1,4,TRUE),(4,'건강용품',1,5,TRUE),
(5,'가구',1,1,TRUE),(5,'침구',1,2,TRUE),(5,'수납/정리',1,3,TRUE),(5,'주방용품',1,4,TRUE),(5,'인테리어소품',1,5,TRUE),(5,'생활용품',1,6,TRUE),
(6,'피트니스',1,1,TRUE),(6,'아웃도어',1,2,TRUE),(6,'구기종목',1,3,TRUE),(6,'수상스포츠',1,4,TRUE),(6,'자전거',1,5,TRUE),
(7,'소설/문학',1,1,TRUE),(7,'경제경영',1,2,TRUE),(7,'자기계발',1,3,TRUE),(7,'학습/참고서',1,4,TRUE),(7,'문구/사무',1,5,TRUE),
(8,'유아동복',1,1,TRUE),(8,'장난감',1,2,TRUE),(8,'유아용품',1,3,TRUE),(8,'반려동물용품',1,4,TRUE),
(9,'자동차용품',1,1,TRUE),(9,'공구',1,2,TRUE),(9,'안전용품',1,3,TRUE),
(10,'국내여행',1,1,TRUE),(10,'해외여행',1,2,TRUE),(10,'티켓/입장권',1,3,TRUE);

-- 소카테고리 (~100개)
INSERT INTO categories (parent_id, name, depth, sort_order, is_active)
SELECT
    parent_id,
    sub_name,
    2,
    row_number() OVER (PARTITION BY parent_id ORDER BY sub_name),
    TRUE
FROM (VALUES
(11,'티셔츠'),(11,'셔츠'),(11,'바지'),(11,'아우터'),(11,'정장'),
(12,'원피스'),(12,'블라우스'),(12,'스커트'),(12,'니트'),(12,'코트'),
(13,'운동화'),(13,'구두'),(13,'샌들'),(13,'부츠'),(13,'슬리퍼'),
(14,'백팩'),(14,'크로스백'),(14,'토트백'),(14,'지갑'),(14,'클러치'),
(17,'스마트폰'),(17,'태블릿'),(17,'스마트워치'),(17,'케이스/필름'),
(18,'노트북'),(18,'데스크탑'),(18,'모니터'),(18,'키보드/마우스'),
(21,'에어컨'),(21,'공기청정기'),(21,'청소기'),(21,'세탁기'),(21,'냉장고'),
(23,'과일/채소'),(23,'정육/수산'),(23,'유제품'),(23,'계란'),
(28,'스킨/토너'),(28,'에센스/세럼'),(28,'크림/로션'),(28,'마스크팩'),(28,'선케어'),
(33,'소파'),(33,'침대'),(33,'책상'),(33,'의자'),(33,'옷장'),
(39,'러닝머신'),(39,'덤벨'),(39,'요가매트'),(39,'풀업바'),
(40,'텐트'),(40,'침낭'),(40,'등산화'),(40,'배낭'),
(48,'캐리어'),(48,'여행파우치'),(48,'여권케이스')
) AS t(parent_id, sub_name);

-- 창고 10개
INSERT INTO warehouses (name, address, region, capacity, is_active) VALUES
('서울 물류센터','서울시 강서구 마곡동','서울',500000,TRUE),
('경기 북부센터','경기도 파주시 월롱면','경기북부',800000,TRUE),
('경기 남부센터','경기도 용인시 처인구','경기남부',800000,TRUE),
('인천 허브','인천시 중구 운서동','인천',600000,TRUE),
('부산 물류센터','부산시 강서구 대저동','부산',400000,TRUE),
('대구 물류센터','대구시 달성군 현풍읍','대구',300000,TRUE),
('대전 물류센터','대전시 유성구 관평동','대전',300000,TRUE),
('광주 물류센터','광주시 광산구 하남동','광주',200000,TRUE),
('제주 물류센터','제주시 조천읍','제주',100000,TRUE),
('강원 물류센터','강원도 원주시 문막읍','강원',150000,TRUE);

-- 배송 정책
INSERT INTO shipping_policies (name, base_fee, free_threshold, additional_fee, region, is_active) VALUES
('기본배송',3000,50000,0,'전국',TRUE),
('제주/도서산간',3000,50000,3000,'제주/도서산간',TRUE),
('새벽배송',5000,80000,0,'수도권',TRUE),
('당일배송',7000,100000,0,'서울',TRUE);

-- 배송 지역
INSERT INTO delivery_areas (region, sub_region, delivery_days, additional_fee, is_available) VALUES
('서울','강남구',1,0,TRUE),('서울','서초구',1,0,TRUE),('서울','마포구',1,0,TRUE),('서울','송파구',1,0,TRUE),('서울','영등포구',1,0,TRUE),
('서울','강서구',1,0,TRUE),('서울','노원구',1,0,TRUE),('서울','성북구',1,0,TRUE),('서울','용산구',1,0,TRUE),('서울','종로구',1,0,TRUE),
('경기','성남시',1,0,TRUE),('경기','수원시',1,0,TRUE),('경기','고양시',1,0,TRUE),('경기','용인시',1,0,TRUE),('경기','안양시',1,0,TRUE),
('경기','부천시',1,0,TRUE),('경기','화성시',2,0,TRUE),('경기','파주시',2,0,TRUE),
('인천','남동구',1,0,TRUE),('인천','부평구',1,0,TRUE),('인천','연수구',1,0,TRUE),
('부산','해운대구',2,0,TRUE),('부산','수영구',2,0,TRUE),('부산','부산진구',2,0,TRUE),
('대구','수성구',2,0,TRUE),('대구','달서구',2,0,TRUE),
('대전','유성구',2,0,TRUE),('대전','서구',2,0,TRUE),
('광주','서구',2,0,TRUE),('울산','남구',2,0,TRUE),
('제주','제주시',3,3000,TRUE),('제주','서귀포시',3,3000,TRUE),
('강원','춘천시',2,0,TRUE),('강원','원주시',2,0,TRUE),('강원','강릉시',3,0,TRUE);

-- 세율
INSERT INTO tax_rates (category, rate, description, effective_from, is_active) VALUES
('일반',10.00,'부가가치세','2020-01-01',TRUE),
('면세식품',0.00,'면세 농수산물','2020-01-01',TRUE),
('주류',10.00,'주류 부가세','2020-01-01',TRUE),
('담배',10.00,'담배 부가세','2020-01-01',TRUE);

-- 쿠폰 200개
INSERT INTO coupons (code, name, discount_type, discount_value, min_order_amount, max_discount, total_quantity, used_quantity, starts_at, expires_at, is_active)
SELECT
    'CPN' || LPAD(i::TEXT, 5, '0'),
    CASE (i % 5)
        WHEN 0 THEN '신규가입 ' || (i % 20 + 5) || '% 할인'
        WHEN 1 THEN (i % 50 + 1) * 1000 || '원 즉시할인'
        WHEN 2 THEN '카테고리 특별 ' || (i % 15 + 5) || '% 쿠폰'
        WHEN 3 THEN '생일축하 ' || (i % 30 + 10) || '% 할인'
        WHEN 4 THEN '재구매 감사 ' || (i % 10 + 3) || '천원 쿠폰'
    END,
    CASE WHEN i % 2 = 0 THEN 'PERCENT' ELSE 'FIXED' END,
    CASE WHEN i % 2 = 0 THEN (i % 20 + 5) ELSE (i % 50 + 1) * 1000 END,
    CASE WHEN i % 2 = 0 THEN (i % 5 + 1) * 10000 ELSE (i % 3 + 1) * 10000 END,
    CASE WHEN i % 2 = 0 THEN (i % 5 + 1) * 5000 ELSE NULL END,
    (100 + (i * 50)),
    (floor(random() * (100 + (i * 50))))::INT,
    NOW() - INTERVAL '30 days' - (random() * INTERVAL '180 days'),
    NOW() + (random() * INTERVAL '180 days'),
    CASE WHEN random() < 0.7 THEN TRUE ELSE FALSE END
FROM generate_series(1, 200) AS i;

-- 프로모션 100개
INSERT INTO promotions (name, description, promotion_type, discount_rate, starts_at, ends_at, is_active)
SELECT
    (ARRAY['봄맞이','여름 특가','가을 세일','겨울 할인','블랙프라이데이','연말 결산','신년 특가','추석 선물','설날 특가','어린이날'])[1 + (i % 10)]
    || ' ' || (ARRAY['대전','페스티벌','빅세일','특별전','기획전','타임딜','시즌오프','클리어런스','한정판매','감사제'])[1 + ((i/10) % 10)],
    '특별 할인 이벤트',
    (ARRAY['SEASON_SALE','FLASH_SALE','BUNDLE_DEAL','CATEGORY_SALE','BRAND_SALE'])[1 + (i % 5)],
    (5 + (i % 40))::NUMERIC(4,2),
    NOW() - (random() * INTERVAL '365 days'),
    NOW() + (random() * INTERVAL '90 days'),
    CASE WHEN random() < 0.6 THEN TRUE ELSE FALSE END
FROM generate_series(1, 100) AS i;

-- FAQ 50개
INSERT INTO faq (category, question, answer, sort_order, is_active)
SELECT
    (ARRAY['주문/결제','배송','교환/반품','회원/포인트','기타'])[1 + (i % 5)],
    CASE (i % 10)
        WHEN 0 THEN '주문 후 결제 수단을 변경할 수 있나요?'
        WHEN 1 THEN '배송은 보통 얼마나 걸리나요?'
        WHEN 2 THEN '반품/교환은 어떻게 하나요?'
        WHEN 3 THEN '포인트는 어떻게 적립되나요?'
        WHEN 4 THEN '회원 탈퇴는 어떻게 하나요?'
        WHEN 5 THEN '해외 배송이 가능한가요?'
        WHEN 6 THEN '영수증 발급은 어떻게 하나요?'
        WHEN 7 THEN '무통장 입금 기한은 얼마인가요?'
        WHEN 8 THEN '쿠폰은 어디서 확인하나요?'
        WHEN 9 THEN '선물 포장 가능한가요?'
    END || ' (FAQ-' || i || ')',
    '답변 내용입니다. 자세한 사항은 고객센터로 문의해주세요.',
    i,
    TRUE
FROM generate_series(1, 50) AS i;

-- 공지사항 100개
INSERT INTO notices (title, content, category, is_pinned, view_count, created_at)
SELECT
    (ARRAY['시스템 점검','이벤트','배송','정책 변경','서비스'])[1 + (i % 5)] || ' 안내 #' || i,
    '공지 내용입니다.',
    (ARRAY['공지','이벤트','배송','정책'])[1 + (i % 4)],
    CASE WHEN i <= 5 THEN TRUE ELSE FALSE END,
    (floor(random() * 10000))::INT,
    NOW() - (random() * INTERVAL '365 days')
FROM generate_series(1, 100) AS i;

-- 배너 30개
INSERT INTO banners (title, image_url, link_url, position, sort_order, starts_at, ends_at, is_active)
SELECT
    '배너 ' || i,
    'https://cdn.example.com/banners/banner_' || i || '.jpg',
    '/promotions/' || i,
    (ARRAY['MAIN_TOP','MAIN_MIDDLE','CATEGORY_TOP','SIDEBAR','POPUP'])[1 + (i % 5)],
    i,
    NOW() - INTERVAL '7 days',
    NOW() + INTERVAL '30 days',
    CASE WHEN random() < 0.8 THEN TRUE ELSE FALSE END
FROM generate_series(1, 30) AS i;
