-- ============================================================
-- 커머스 플랫폼 스키마 (50 테이블)
-- ============================================================

-- 1. 멤버십 등급 정의
CREATE TABLE membership_tiers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    min_spending INT NOT NULL DEFAULT 0,
    point_rate NUMERIC(3,2) NOT NULL DEFAULT 0.01,
    benefits TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. 고객
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(200) UNIQUE NOT NULL,
    phone VARCHAR(20),
    gender VARCHAR(10),
    birth_date DATE,
    grade VARCHAR(20) DEFAULT 'BRONZE',
    total_spending BIGINT DEFAULT 0,
    point_balance INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. 고객 배송지
CREATE TABLE customer_addresses (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    label VARCHAR(50),
    recipient_name VARCHAR(100),
    phone VARCHAR(20),
    zip_code VARCHAR(10),
    address TEXT NOT NULL,
    address_detail TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. 브랜드
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    logo_url TEXT,
    country VARCHAR(50),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. 판매자
CREATE TABLE sellers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    business_number VARCHAR(20),
    representative VARCHAR(100),
    email VARCHAR(200),
    phone VARCHAR(20),
    address TEXT,
    commission_rate NUMERIC(4,2) DEFAULT 10.00,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    joined_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. 카테고리
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES categories(id),
    name VARCHAR(100) NOT NULL,
    depth INT DEFAULT 0,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. 상품
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    seller_id INT REFERENCES sellers(id),
    brand_id INT REFERENCES brands(id),
    category_id INT REFERENCES categories(id),
    name VARCHAR(300) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    cost_price INT,
    stock INT NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    is_featured BOOLEAN DEFAULT FALSE,
    weight_gram INT,
    rating_avg NUMERIC(2,1) DEFAULT 0,
    review_count INT DEFAULT 0,
    sales_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. 상품 옵션
CREATE TABLE product_options (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    option_type VARCHAR(50),
    option_value VARCHAR(100),
    additional_price INT DEFAULT 0,
    stock INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

-- 9. 상품 이미지
CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    image_url TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    is_main BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 10. 상품 태그
CREATE TABLE product_tags (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    tag VARCHAR(50) NOT NULL
);

-- 11. 창고
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    region VARCHAR(50),
    capacity INT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 12. 재고 (창고별)
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    warehouse_id INT NOT NULL REFERENCES warehouses(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 0,
    reserved INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 13. 쿠폰
CREATE TABLE coupons (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    discount_type VARCHAR(20) NOT NULL,
    discount_value INT NOT NULL,
    min_order_amount INT DEFAULT 0,
    max_discount INT,
    total_quantity INT,
    used_quantity INT DEFAULT 0,
    starts_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 14. 쿠폰 사용 내역
CREATE TABLE coupon_usage (
    id SERIAL PRIMARY KEY,
    coupon_id INT NOT NULL REFERENCES coupons(id),
    customer_id INT NOT NULL REFERENCES customers(id),
    order_id INT,
    discount_amount INT NOT NULL,
    used_at TIMESTAMPTZ DEFAULT NOW()
);

-- 15. 프로모션
CREATE TABLE promotions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    promotion_type VARCHAR(50),
    discount_rate NUMERIC(4,2),
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 16. 프로모션 대상 상품
CREATE TABLE promotion_products (
    id SERIAL PRIMARY KEY,
    promotion_id INT NOT NULL REFERENCES promotions(id),
    product_id INT NOT NULL REFERENCES products(id),
    promotion_price INT
);

-- 17. 배송 정책
CREATE TABLE shipping_policies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    base_fee INT NOT NULL DEFAULT 3000,
    free_threshold INT DEFAULT 50000,
    additional_fee INT DEFAULT 0,
    region VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE
);

-- 18. 배송 가능 지역
CREATE TABLE delivery_areas (
    id SERIAL PRIMARY KEY,
    region VARCHAR(50) NOT NULL,
    sub_region VARCHAR(50),
    delivery_days INT DEFAULT 1,
    additional_fee INT DEFAULT 0,
    is_available BOOLEAN DEFAULT TRUE
);

-- 19. 주문
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    order_number VARCHAR(30) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    total_amount INT NOT NULL DEFAULT 0,
    discount_amount INT DEFAULT 0,
    shipping_fee INT DEFAULT 0,
    final_amount INT NOT NULL DEFAULT 0,
    shipping_address TEXT,
    shipping_memo TEXT,
    ordered_at TIMESTAMPTZ DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ
);

-- 20. 주문 상세
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    product_id INT NOT NULL REFERENCES products(id),
    product_option_id INT REFERENCES product_options(id),
    quantity INT NOT NULL,
    unit_price INT NOT NULL,
    subtotal INT NOT NULL,
    status VARCHAR(20) DEFAULT 'ORDERED'
);

-- 21. 결제
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    payment_method VARCHAR(30) NOT NULL,
    amount INT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    pg_provider VARCHAR(50),
    pg_transaction_id VARCHAR(100),
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 22. 배송
CREATE TABLE shipments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    carrier VARCHAR(50),
    tracking_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'PREPARING',
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 23. 배송 추적
CREATE TABLE shipment_tracking (
    id SERIAL PRIMARY KEY,
    shipment_id INT NOT NULL REFERENCES shipments(id),
    status VARCHAR(50) NOT NULL,
    location VARCHAR(200),
    description TEXT,
    tracked_at TIMESTAMPTZ DEFAULT NOW()
);

-- 24. 환불
CREATE TABLE refunds (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    payment_id INT REFERENCES payments(id),
    reason TEXT,
    amount INT NOT NULL,
    status VARCHAR(20) DEFAULT 'REQUESTED',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- 25. 반품
CREATE TABLE returns (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    order_item_id INT REFERENCES order_items(id),
    customer_id INT NOT NULL REFERENCES customers(id),
    reason VARCHAR(50),
    reason_detail TEXT,
    status VARCHAR(20) DEFAULT 'REQUESTED',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- 26. 교환
CREATE TABLE exchanges (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id),
    order_item_id INT REFERENCES order_items(id),
    customer_id INT NOT NULL REFERENCES customers(id),
    reason VARCHAR(50),
    new_product_id INT REFERENCES products(id),
    status VARCHAR(20) DEFAULT 'REQUESTED',
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

-- 27. 리뷰
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    customer_id INT NOT NULL REFERENCES customers(id),
    order_item_id INT REFERENCES order_items(id),
    rating INT CHECK (rating BETWEEN 1 AND 5),
    content TEXT,
    image_urls TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    helpful_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 28. 위시리스트
CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    product_id INT NOT NULL REFERENCES products(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 29. 장바구니
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    product_id INT NOT NULL REFERENCES products(id),
    product_option_id INT REFERENCES product_options(id),
    quantity INT NOT NULL DEFAULT 1,
    added_at TIMESTAMPTZ DEFAULT NOW()
);

-- 30. 포인트 내역
CREATE TABLE point_history (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    order_id INT REFERENCES orders(id),
    type VARCHAR(20) NOT NULL,
    amount INT NOT NULL,
    balance_after INT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 31. 기프트카드
CREATE TABLE gift_cards (
    id SERIAL PRIMARY KEY,
    code VARCHAR(30) UNIQUE NOT NULL,
    initial_balance INT NOT NULL,
    current_balance INT NOT NULL,
    purchaser_id INT REFERENCES customers(id),
    recipient_email VARCHAR(200),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 32. 기프트카드 사용 내역
CREATE TABLE gift_card_transactions (
    id SERIAL PRIMARY KEY,
    gift_card_id INT NOT NULL REFERENCES gift_cards(id),
    order_id INT REFERENCES orders(id),
    amount INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 33. 판매자 정산
CREATE TABLE seller_settlements (
    id SERIAL PRIMARY KEY,
    seller_id INT NOT NULL REFERENCES sellers(id),
    settlement_period VARCHAR(20),
    total_sales BIGINT NOT NULL,
    commission BIGINT NOT NULL,
    shipping_subsidy BIGINT DEFAULT 0,
    net_amount BIGINT NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING',
    settled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 34. 1:1 문의
CREATE TABLE inquiries (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    category VARCHAR(50),
    title VARCHAR(300) NOT NULL,
    content TEXT NOT NULL,
    order_id INT REFERENCES orders(id),
    product_id INT REFERENCES products(id),
    status VARCHAR(20) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 35. 문의 답변
CREATE TABLE inquiry_replies (
    id SERIAL PRIMARY KEY,
    inquiry_id INT NOT NULL REFERENCES inquiries(id),
    author_type VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 36. 상품 Q&A
CREATE TABLE product_qna (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    customer_id INT NOT NULL REFERENCES customers(id),
    question TEXT NOT NULL,
    answer TEXT,
    is_secret BOOLEAN DEFAULT FALSE,
    answered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 37. FAQ
CREATE TABLE faq (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50),
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 38. 공지사항
CREATE TABLE notices (
    id SERIAL PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50),
    is_pinned BOOLEAN DEFAULT FALSE,
    view_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 39. 배너
CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    image_url TEXT NOT NULL,
    link_url TEXT,
    position VARCHAR(50),
    sort_order INT DEFAULT 0,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 40. 묶음 상품
CREATE TABLE product_bundles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    bundle_price INT NOT NULL,
    original_price INT NOT NULL,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 41. 묶음 구성 상품
CREATE TABLE bundle_items (
    id SERIAL PRIMARY KEY,
    bundle_id INT NOT NULL REFERENCES product_bundles(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT DEFAULT 1
);

-- 42. 세율
CREATE TABLE tax_rates (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50),
    rate NUMERIC(4,2) NOT NULL,
    description TEXT,
    effective_from DATE,
    is_active BOOLEAN DEFAULT TRUE
);

-- 43. 검색 로그
CREATE TABLE search_logs (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id),
    keyword VARCHAR(200) NOT NULL,
    result_count INT DEFAULT 0,
    clicked_product_id INT REFERENCES products(id),
    searched_at TIMESTAMPTZ DEFAULT NOW()
);

-- 44. 페이지뷰
CREATE TABLE page_views (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id),
    session_id VARCHAR(100),
    page_type VARCHAR(50),
    reference_id INT,
    referrer VARCHAR(500),
    device_type VARCHAR(20),
    viewed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 45. 로그인 이력
CREATE TABLE login_history (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    ip_address VARCHAR(50),
    device_type VARCHAR(20),
    os VARCHAR(50),
    browser VARCHAR(50),
    is_success BOOLEAN DEFAULT TRUE,
    logged_at TIMESTAMPTZ DEFAULT NOW()
);

-- 46. 알림
CREATE TABLE customer_notifications (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    reference_type VARCHAR(50),
    reference_id INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 47. 일별 매출 통계
CREATE TABLE daily_sales (
    id SERIAL PRIMARY KEY,
    sale_date DATE NOT NULL,
    total_orders INT DEFAULT 0,
    total_revenue BIGINT DEFAULT 0,
    total_refunds BIGINT DEFAULT 0,
    net_revenue BIGINT DEFAULT 0,
    new_customers INT DEFAULT 0,
    returning_customers INT DEFAULT 0,
    avg_order_value INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 48. 카테고리별 통계
CREATE TABLE category_stats (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id),
    stat_date DATE NOT NULL,
    product_count INT DEFAULT 0,
    order_count INT DEFAULT 0,
    revenue BIGINT DEFAULT 0,
    avg_rating NUMERIC(2,1) DEFAULT 0,
    return_rate NUMERIC(4,2) DEFAULT 0
);

-- 49. 상품 랭킹
CREATE TABLE product_rankings (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    ranking_type VARCHAR(50) NOT NULL,
    category_id INT REFERENCES categories(id),
    rank_position INT NOT NULL,
    score NUMERIC(10,2),
    ranking_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 50. 알림 설정
CREATE TABLE notification_settings (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id),
    channel VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- nunu 앱 내부 테이블
CREATE TABLE memories (
    id BIGSERIAL PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    user_question TEXT NOT NULL,
    summary TEXT,
    sql_used TEXT,
    tools_used TEXT,
    embedding vector(1536),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE knowledges (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    title VARCHAR(300) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT,
    embedding vector(1536),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    channel_id TEXT,
    user_id TEXT,
    thread_ts TEXT,
    question TEXT,
    answer TEXT,
    sql_used TEXT,
    tools_used TEXT,
    success BOOLEAN DEFAULT FALSE,
    feedback TEXT,
    duration BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
