-- 고객 100,000명
INSERT INTO customers (name, email, phone, gender, birth_date, grade, total_spending, point_balance, is_active, created_at)
SELECT
    (ARRAY['김','이','박','최','정','강','조','윤','장','임','한','오','서','신','권','황','안','송','류','홍',
           '전','고','문','양','손','배','백','허','유','남','심','노','하','곽','성','차','주','우','구','민'])[1 + (i % 40)]
    || (ARRAY['민준','서연','도윤','서윤','시우','지우','예준','하윤','주원','하은',
              '지호','서현','준서','민서','현우','지민','준영','수빈','건우','예은',
              '성민','다은','우진','채원','태영','소율','민재','유진','승현','나은',
              '재현','은서','동현','소희','정훈','미래','상우','보라','진우','혜원',
              '승준','가영','태민','예지','시현','다현','윤호','소연','재민','은비'])[1 + (i % 50)],
    'u' || i || '@' || (ARRAY['gmail.com','naver.com','kakao.com','daum.net','hanmail.net','outlook.com','nate.com','icloud.com'])[1 + (i % 8)],
    '010-' || LPAD((1000 + (i % 9000))::TEXT, 4, '0') || '-' || LPAD((1000 + ((i * 7 + 13) % 9000))::TEXT, 4, '0'),
    CASE WHEN random() < 0.48 THEN 'M' WHEN random() < 0.96 THEN 'F' ELSE NULL END,
    DATE '1960-01-01' + (floor(random() * 22000))::INT,
    CASE
        WHEN random() < 0.05 THEN 'DIAMOND'
        WHEN random() < 0.15 THEN 'GOLD'
        WHEN random() < 0.40 THEN 'SILVER'
        ELSE 'BRONZE'
    END,
    (floor(random() * 5000000))::BIGINT,
    (floor(random() * 50000))::INT,
    CASE WHEN random() < 0.92 THEN TRUE ELSE FALSE END,
    NOW() - (random() * INTERVAL '1095 days')
FROM generate_series(1, 100000) AS i;

-- 고객 배송지 (~200,000건, 고객당 1~3개)
INSERT INTO customer_addresses (customer_id, label, recipient_name, phone, zip_code, address, address_detail, is_default, created_at)
SELECT
    c.id,
    (ARRAY['집','회사','부모님댁'])[addr_num],
    c.name,
    c.phone,
    LPAD((10000 + (c.id * addr_num) % 90000)::TEXT, 5, '0'),
    (ARRAY[
        '서울시 강남구 역삼대로','서울시 서초구 서초대로','서울시 마포구 월드컵북로','서울시 송파구 올림픽로','서울시 영등포구 여의대방로',
        '서울시 강서구 화곡로','서울시 노원구 동일로','서울시 성북구 보문로','서울시 용산구 이태원로','서울시 종로구 종로',
        '서울시 관악구 관악로','서울시 동작구 상도로','서울시 광진구 능동로','서울시 중랑구 면목로','서울시 은평구 진관로',
        '경기도 성남시 분당구 불정로','경기도 수원시 영통구 광교로','경기도 고양시 일산동구 중앙로','경기도 용인시 수지구 성복로','경기도 안양시 동안구 시민대로',
        '경기도 부천시 원미구 길주로','경기도 화성시 동탄대로','경기도 파주시 금바위로','경기도 김포시 풍무로','경기도 광명시 오리로',
        '인천시 남동구 구월로','인천시 부평구 부평대로','인천시 연수구 센트럴로',
        '부산시 해운대구 해운대로','부산시 수영구 광안해변로','부산시 부산진구 중앙대로','부산시 사하구 낙동대로',
        '대구시 수성구 달구벌대로','대구시 달서구 월배로','대구시 중구 동성로',
        '대전시 유성구 대학로','대전시 서구 둔산로','대전시 중구 대종로',
        '광주시 서구 상무대로','광주시 북구 용봉로',
        '울산시 남구 삼산로','울산시 중구 성남로',
        '세종시 조치원읍 세종로',
        '강원도 춘천시 중앙로','강원도 원주시 원일로','강원도 강릉시 경포로',
        '충북 청주시 상당구 상당로','충남 천안시 동남구 만남로','충남 아산시 배방읍 희망로',
        '제주시 연동 노형로'
    ])[1 + ((c.id + addr_num * 7) % 50)] || ' ' || (1 + (c.id % 300))::TEXT,
    (100 + (c.id % 2000))::TEXT || '동 ' || (100 + ((c.id * 3) % 1500))::TEXT || '호',
    CASE WHEN addr_num = 1 THEN TRUE ELSE FALSE END,
    c.created_at + (random() * INTERVAL '30 days')
FROM customers c
CROSS JOIN generate_series(1, 3) AS addr_num
WHERE addr_num = 1 OR random() < 0.5;
