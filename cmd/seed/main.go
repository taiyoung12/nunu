package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pgvector/pgvector-go"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Knowledge struct {
	ID        uint            `gorm:"primaryKey"`
	Category  string          `gorm:"index;not null"`
	Title     string          `gorm:"not null"`
	Content   string          `gorm:"type:text;not null"`
	Tags      string          `gorm:"type:text"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)"`
}

func (Knowledge) TableName() string { return "knowledges" }

func main() {
	dsn := "host=localhost user=nunu password=nunu dbname=nunu port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("db connect error:", err)
	}

	// 기존 지식 삭제
	db.Exec("DELETE FROM knowledges")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}
	client := openai.NewClient(strings.TrimSpace(apiKey))

	knowledges := []Knowledge{
		// ===== 스키마 =====
		{Category: "schema", Title: "membership_tiers 테이블",
			Content: `테이블명: membership_tiers
설명: 멤버십 등급 정의
컬럼: id(PK), name(등급명: BRONZE/SILVER/GOLD/DIAMOND), min_spending(최소 누적구매액), point_rate(적립률), benefits(혜택설명), created_at`,
			Tags: `["schema","membership_tiers","멤버십","등급"]`},

		{Category: "schema", Title: "customers 테이블",
			Content: `테이블명: customers
설명: 고객 정보 (10만건)
컬럼: id(PK), name(고객명), email(이메일,UNIQUE), phone(전화번호), gender(성별:M/F/NULL), birth_date(생년월일), grade(등급:DIAMOND/GOLD/SILVER/BRONZE), total_spending(총구매액), point_balance(포인트잔액), is_active(활성여부), created_at, updated_at`,
			Tags: `["schema","customers","고객","회원","등급","포인트"]`},

		{Category: "schema", Title: "customer_addresses 테이블",
			Content: `테이블명: customer_addresses
설명: 고객 배송지 (고객당 1~3개, 약 20만건)
컬럼: id(PK), customer_id(FK→customers), label(집/회사/부모님댁), recipient_name, phone, zip_code(우편번호), address(주소), address_detail(상세주소), is_default(기본배송지여부), created_at
주소 지역: 서울/경기/인천/부산/대구/대전/광주/울산/세종/강원/충북/충남/전북/전남/경북/경남/제주`,
			Tags: `["schema","customer_addresses","배송지","주소","지역"]`},

		{Category: "schema", Title: "brands 테이블",
			Content: `테이블명: brands
설명: 브랜드 정보 (200개)
컬럼: id(PK), name(브랜드명), logo_url, country(국가), description, is_active, created_at`,
			Tags: `["schema","brands","브랜드"]`},

		{Category: "schema", Title: "sellers 테이블",
			Content: `테이블명: sellers
설명: 판매자/입점업체 정보 (500개)
컬럼: id(PK), name(업체명), business_number(사업자번호), representative(대표자), email, phone, address, commission_rate(수수료율%), status(ACTIVE/SUSPENDED/INACTIVE), joined_at`,
			Tags: `["schema","sellers","판매자","입점업체","수수료"]`},

		{Category: "schema", Title: "categories 테이블",
			Content: `테이블명: categories
설명: 상품 카테고리 (대/중/소 3단계, 116개)
컬럼: id(PK), parent_id(FK→categories, 상위카테고리), name, depth(0=대,1=중,2=소), sort_order, is_active, created_at
대카테고리: 패션, 디지털/가전, 식품/음료, 뷰티/건강, 홈/리빙, 스포츠/레저, 도서/문구, 키즈/반려, 자동차/공구, 여행/티켓`,
			Tags: `["schema","categories","카테고리","분류","대분류","중분류","소분류"]`},

		{Category: "schema", Title: "products 테이블",
			Content: `테이블명: products
설명: 상품 정보 (10만건)
컬럼: id(PK), seller_id(FK→sellers), brand_id(FK→brands), category_id(FK→categories), name(상품명), description, price(판매가), cost_price(원가), stock(재고), status(ACTIVE/SOLDOUT/HIDDEN), is_featured(추천상품), weight_gram(무게), rating_avg(평균별점), review_count(리뷰수), sales_count(판매수), created_at, updated_at`,
			Tags: `["schema","products","상품","가격","재고","별점","판매"]`},

		{Category: "schema", Title: "product_options 테이블",
			Content: `테이블명: product_options
설명: 상품 옵션 (사이즈/색상/소재 등, 20만건)
컬럼: id(PK), product_id(FK→products), option_type(사이즈/색상/소재/기타), option_value(XS~XXL, 블랙~핑크 등), additional_price(추가금액), stock(옵션재고), is_active`,
			Tags: `["schema","product_options","옵션","사이즈","색상"]`},

		{Category: "schema", Title: "product_images, product_tags 테이블",
			Content: `테이블명: product_images
설명: 상품 이미지 (상품당 1~5개, 30만건)
컬럼: id(PK), product_id(FK→products), image_url, sort_order, is_main(대표이미지여부), created_at

테이블명: product_tags
설명: 상품 태그 (20만건)
컬럼: id(PK), product_id(FK→products), tag(인기/신상/할인/베스트/추천/한정/시즌/이벤트/무료배송/당일출고/친환경/프리미엄/가성비 등)`,
			Tags: `["schema","product_images","product_tags","이미지","태그"]`},

		{Category: "schema", Title: "warehouses, inventory 테이블",
			Content: `테이블명: warehouses
설명: 물류 창고 (10개: 서울/경기북부/경기남부/인천/부산/대구/대전/광주/제주/강원)
컬럼: id(PK), name, address, region, capacity, is_active, created_at

테이블명: inventory
설명: 창고별 재고 (30만건)
컬럼: id(PK), warehouse_id(FK→warehouses), product_id(FK→products), quantity(수량), reserved(예약수량), updated_at`,
			Tags: `["schema","warehouses","inventory","창고","재고","물류"]`},

		{Category: "schema", Title: "orders 테이블",
			Content: `테이블명: orders
설명: 주문 (10만건)
컬럼: id(PK), customer_id(FK→customers), order_number(주문번호,UNIQUE), status(PENDING/CONFIRMED/SHIPPING/DELIVERED/CANCELLED), total_amount(총금액), discount_amount(할인액), shipping_fee(배송비), final_amount(최종결제액=총금액-할인+배송비), shipping_address(배송지), shipping_memo(배송메모), ordered_at(주문일시), confirmed_at, shipped_at, delivered_at, cancelled_at`,
			Tags: `["schema","orders","주문","배송","결제","주문상태"]`},

		{Category: "schema", Title: "order_items 테이블",
			Content: `테이블명: order_items
설명: 주문 상세 (주문당 1~5개 상품, 30만건)
컬럼: id(PK), order_id(FK→orders), product_id(FK→products), product_option_id(FK→product_options), quantity(수량), unit_price(단가), subtotal(소계), status(ORDERED 등)`,
			Tags: `["schema","order_items","주문상세","주문내역"]`},

		{Category: "schema", Title: "payments 테이블",
			Content: `테이블명: payments
설명: 결제 정보 (10만건)
컬럼: id(PK), order_id(FK→orders), payment_method(CARD/BANK_TRANSFER/KAKAO_PAY/NAVER_PAY/TOSS_PAY/PHONE/POINT), amount(결제금액), status(PENDING/PAID/CANCELLED), pg_provider(PG사: KG이니시스/NHN KCP/토스페이먼츠/나이스페이/카카오페이), pg_transaction_id, paid_at, cancelled_at, created_at`,
			Tags: `["schema","payments","결제","카드","페이","PG"]`},

		{Category: "schema", Title: "shipments, shipment_tracking 테이블",
			Content: `테이블명: shipments
설명: 배송 정보 (10만건, CONFIRMED 이상 주문)
컬럼: id(PK), order_id(FK→orders), carrier(택배사: CJ대한통운/한진택배/롯데택배/우체국택배/로젠택배), tracking_number(운송장번호), status(PREPARING/IN_TRANSIT/DELIVERED), shipped_at, delivered_at, created_at

테이블명: shipment_tracking
설명: 배송 추적 이력 (40만건)
컬럼: id(PK), shipment_id(FK→shipments), status(집하/간선상차/배달출발/배달완료), location(위치), description, tracked_at`,
			Tags: `["schema","shipments","shipment_tracking","배송","택배","운송장","추적"]`},

		{Category: "schema", Title: "refunds 테이블",
			Content: `테이블명: refunds
설명: 환불 (1만건)
컬럼: id(PK), order_id(FK→orders), payment_id(FK→payments), reason(단순변심/상품불량/오배송/사이즈안맞음/상품파손/다른상품수령/색상상이/배송지연), amount(환불금액), status(REQUESTED/PROCESSING/COMPLETED), requested_at, completed_at`,
			Tags: `["schema","refunds","환불"]`},

		{Category: "schema", Title: "returns, exchanges 테이블",
			Content: `테이블명: returns
설명: 반품 (8,000건)
컬럼: id(PK), order_id(FK→orders), order_item_id(FK→order_items), customer_id(FK→customers), reason(DEFECT/WRONG_ITEM/CHANGE_MIND/SIZE_ISSUE/DAMAGED/COLOR_DIFF), reason_detail, status(REQUESTED/COLLECTING/INSPECTING/COMPLETED/REJECTED), requested_at, completed_at

테이블명: exchanges
설명: 교환 (5,000건)
컬럼: id(PK), order_id(FK→orders), order_item_id(FK→order_items), customer_id(FK→customers), reason(SIZE_CHANGE/COLOR_CHANGE/DEFECT/WRONG_ITEM), new_product_id(FK→products), status(REQUESTED/COLLECTING/SHIPPING_NEW/COMPLETED), requested_at, completed_at`,
			Tags: `["schema","returns","exchanges","반품","교환"]`},

		{Category: "schema", Title: "reviews 테이블",
			Content: `테이블명: reviews
설명: 상품 리뷰 (10만건)
컬럼: id(PK), product_id(FK→products), customer_id(FK→customers), order_item_id(FK→order_items), rating(1~5), content(리뷰내용), image_urls(이미지URL JSON), is_verified(인증구매여부), helpful_count(도움됨수), created_at, updated_at`,
			Tags: `["schema","reviews","리뷰","별점","평점"]`},

		{Category: "schema", Title: "wishlists, cart_items 테이블",
			Content: `테이블명: wishlists
설명: 찜/위시리스트 (10만건)
컬럼: id(PK), customer_id(FK→customers), product_id(FK→products), created_at

테이블명: cart_items
설명: 장바구니 (5만건)
컬럼: id(PK), customer_id(FK→customers), product_id(FK→products), product_option_id(FK→product_options), quantity, added_at`,
			Tags: `["schema","wishlists","cart_items","위시리스트","찜","장바구니"]`},

		{Category: "schema", Title: "coupons, coupon_usage 테이블",
			Content: `테이블명: coupons
설명: 쿠폰 (200개)
컬럼: id(PK), code(쿠폰코드,UNIQUE), name(쿠폰명), discount_type(PERCENT/FIXED), discount_value(할인값), min_order_amount(최소주문금액), max_discount(최대할인), total_quantity(총발행수), used_quantity(사용수), starts_at, expires_at, is_active, created_at

테이블명: coupon_usage
설명: 쿠폰 사용 내역 (10만건)
컬럼: id(PK), coupon_id(FK→coupons), customer_id(FK→customers), order_id, discount_amount, used_at`,
			Tags: `["schema","coupons","coupon_usage","쿠폰","할인"]`},

		{Category: "schema", Title: "promotions, promotion_products 테이블",
			Content: `테이블명: promotions
설명: 프로모션/이벤트 (100개)
컬럼: id(PK), name(프로모션명), description, promotion_type(SEASON_SALE/FLASH_SALE/BUNDLE_DEAL/CATEGORY_SALE/BRAND_SALE), discount_rate(할인율%), starts_at, ends_at, is_active, created_at

테이블명: promotion_products
설명: 프로모션 대상 상품 (5만건)
컬럼: id(PK), promotion_id(FK→promotions), product_id(FK→products), promotion_price(프로모션가격)`,
			Tags: `["schema","promotions","promotion_products","프로모션","이벤트","세일"]`},

		{Category: "schema", Title: "point_history 테이블",
			Content: `테이블명: point_history
설명: 포인트 적립/사용 내역 (20만건)
컬럼: id(PK), customer_id(FK→customers), order_id(FK→orders), type(EARN/USE/EXPIRE/ADMIN), amount(양수=적립,음수=사용), balance_after(변동후잔액), description(주문적립/이벤트적립/포인트만료/관리자지급/리뷰적립/생일적립), created_at`,
			Tags: `["schema","point_history","포인트","적립","사용"]`},

		{Category: "schema", Title: "gift_cards, gift_card_transactions 테이블",
			Content: `테이블명: gift_cards
설명: 기프트카드 (5,000개)
컬럼: id(PK), code(코드,UNIQUE), initial_balance(초기잔액), current_balance(현재잔액), purchaser_id(FK→customers), recipient_email, status(ACTIVE/USED/EXPIRED), expires_at, created_at

테이블명: gift_card_transactions
설명: 기프트카드 사용 내역 (1만건)
컬럼: id(PK), gift_card_id(FK→gift_cards), order_id(FK→orders), amount, type(USE/CHARGE/REFUND), created_at`,
			Tags: `["schema","gift_cards","gift_card_transactions","기프트카드","상품권"]`},

		{Category: "schema", Title: "seller_settlements 테이블",
			Content: `테이블명: seller_settlements
설명: 판매자 정산 (~1만건, 월별)
컬럼: id(PK), seller_id(FK→sellers), settlement_period(정산월:YYYY-MM), total_sales(총매출), commission(수수료), shipping_subsidy(배송지원금), net_amount(정산금액=매출-수수료+배송지원), status(PENDING/SETTLED), settled_at, created_at`,
			Tags: `["schema","seller_settlements","정산","수수료","매출"]`},

		{Category: "schema", Title: "inquiries, inquiry_replies 테이블",
			Content: `테이블명: inquiries
설명: 1:1 문의 (3만건)
컬럼: id(PK), customer_id(FK→customers), category(주문/결제/배송/교환/반품/상품문의/기타), title, content, order_id(FK→orders), product_id(FK→products), status(OPEN/IN_PROGRESS/RESOLVED/CLOSED), created_at

테이블명: inquiry_replies
설명: 문의 답변 (2.5만건)
컬럼: id(PK), inquiry_id(FK→inquiries), author_type(STAFF/CUSTOMER), content, created_at`,
			Tags: `["schema","inquiries","inquiry_replies","문의","CS","고객센터"]`},

		{Category: "schema", Title: "product_qna 테이블",
			Content: `테이블명: product_qna
설명: 상품 Q&A (5만건)
컬럼: id(PK), product_id(FK→products), customer_id(FK→customers), question(질문), answer(답변,NULL이면 미답변), is_secret(비밀글), answered_at, created_at`,
			Tags: `["schema","product_qna","QNA","질문","답변","상품문의"]`},

		{Category: "schema", Title: "search_logs 테이블",
			Content: `테이블명: search_logs
설명: 검색 로그 (10만건)
컬럼: id(PK), customer_id(FK→customers,NULL=비로그인), keyword(검색어), result_count(검색결과수), clicked_product_id(FK→products,클릭상품), searched_at`,
			Tags: `["schema","search_logs","검색","키워드","인기검색어"]`},

		{Category: "schema", Title: "page_views 테이블",
			Content: `테이블명: page_views
설명: 페이지뷰 로그 (10만건)
컬럼: id(PK), customer_id(FK→customers,NULL=비로그인), session_id, page_type(HOME/PRODUCT/CATEGORY/SEARCH/CART/ORDER/MYPAGE/EVENT/BRAND), reference_id, referrer(유입경로URL), device_type(MOBILE/PC/TABLET), viewed_at`,
			Tags: `["schema","page_views","페이지뷰","트래픽","방문"]`},

		{Category: "schema", Title: "login_history 테이블",
			Content: `테이블명: login_history
설명: 로그인 이력 (10만건)
컬럼: id(PK), customer_id(FK→customers), ip_address, device_type(MOBILE/PC/TABLET), os(iOS/Android/Windows/macOS/Linux), browser(Chrome/Safari/Samsung Internet 등), is_success(성공여부), logged_at`,
			Tags: `["schema","login_history","로그인","접속","디바이스"]`},

		{Category: "schema", Title: "customer_notifications, notification_settings 테이블",
			Content: `테이블명: customer_notifications
설명: 고객 알림 (10만건)
컬럼: id(PK), customer_id(FK→customers), type(ORDER/DELIVERY/PROMOTION/POINT/COUPON/REVIEW/SYSTEM), title, content, is_read(읽음여부), reference_type, reference_id, created_at

테이블명: notification_settings
설명: 알림 설정 (~48만건)
컬럼: id(PK), customer_id(FK→customers), channel(PUSH/SMS/EMAIL/KAKAO), type(ORDER/DELIVERY/PROMOTION/POINT), is_enabled(활성여부), updated_at`,
			Tags: `["schema","customer_notifications","notification_settings","알림","푸시","설정"]`},

		{Category: "schema", Title: "shipping_policies, delivery_areas 테이블",
			Content: `테이블명: shipping_policies
설명: 배송 정책 (4건: 기본배송/제주도서산간/새벽배송/당일배송)
컬럼: id(PK), name, base_fee(기본배송비), free_threshold(무료배송기준), additional_fee(추가비용), region, is_active

테이블명: delivery_areas
설명: 배송 가능 지역 (35건)
컬럼: id(PK), region(시도), sub_region(시군구), delivery_days(배송소요일), additional_fee(추가비용), is_available`,
			Tags: `["schema","shipping_policies","delivery_areas","배송비","배송정책","지역"]`},

		{Category: "schema", Title: "product_bundles, bundle_items 테이블",
			Content: `테이블명: product_bundles
설명: 묶음 상품 (1,000개)
컬럼: id(PK), name, description, bundle_price(묶음가), original_price(정상가합계), status(ACTIVE/INACTIVE), created_at

테이블명: bundle_items
설명: 묶음 구성 상품 (2,000건)
컬럼: id(PK), bundle_id(FK→product_bundles), product_id(FK→products), quantity`,
			Tags: `["schema","product_bundles","bundle_items","묶음","세트","패키지"]`},

		{Category: "schema", Title: "daily_sales, category_stats, product_rankings 테이블",
			Content: `테이블명: daily_sales
설명: 일별 매출 통계 (365일)
컬럼: id(PK), sale_date, total_orders(주문수), total_revenue(총매출), total_refunds(환불액), net_revenue(순매출), new_customers(신규고객수), returning_customers(재방문고객수), avg_order_value(평균주문금액), created_at

테이블명: category_stats
설명: 카테고리별 일간 통계 (3,480건)
컬럼: id(PK), category_id(FK→categories), stat_date, product_count, order_count, revenue, avg_rating, return_rate(반품률%)

테이블명: product_rankings
설명: 상품 랭킹 (10만건)
컬럼: id(PK), product_id(FK→products), ranking_type(BEST_SELLER/NEW_ARRIVAL/TOP_RATED/MOST_WISHED/HOT_DEAL), category_id(FK→categories), rank_position, score, ranking_date, created_at`,
			Tags: `["schema","daily_sales","category_stats","product_rankings","통계","매출","랭킹","순위"]`},

		{Category: "schema", Title: "faq, notices, banners 테이블",
			Content: `테이블명: faq
설명: FAQ (50건, 카테고리: 주문/결제, 배송, 교환/반품, 회원/포인트, 기타)
컬럼: id(PK), category, question, answer, sort_order, is_active, created_at

테이블명: notices
설명: 공지사항 (100건)
컬럼: id(PK), title, content, category(공지/이벤트/배송/정책), is_pinned(상단고정), view_count, created_at

테이블명: banners
설명: 배너 (30건)
컬럼: id(PK), title, image_url, link_url, position(MAIN_TOP/MAIN_MIDDLE/CATEGORY_TOP/SIDEBAR/POPUP), sort_order, starts_at, ends_at, is_active, created_at`,
			Tags: `["schema","faq","notices","banners","FAQ","공지사항","배너"]`},

		{Category: "schema", Title: "tax_rates 테이블",
			Content: `테이블명: tax_rates
설명: 세율 정보 (4건: 일반10%/면세식품0%/주류10%/담배10%)
컬럼: id(PK), category, rate(세율%), description, effective_from, is_active`,
			Tags: `["schema","tax_rates","세율","부가세","면세"]`},

		// ===== 테이블 관계도 =====
		{Category: "schema", Title: "테이블 관계도 (ERD 요약)",
			Content: `핵심 테이블 관계:
[고객 중심]
- customers 1:N customer_addresses (배송지)
- customers 1:N orders (주문)
- customers 1:N reviews (리뷰)
- customers 1:N wishlists (찜)
- customers 1:N cart_items (장바구니)
- customers 1:N point_history (포인트)
- customers 1:N inquiries (문의)
- customers 1:N login_history (로그인)
- customers 1:N coupon_usage (쿠폰사용)

[상품 중심]
- categories(대→중→소 self-join via parent_id) 1:N products
- brands 1:N products
- sellers 1:N products
- products 1:N product_options / product_images / product_tags
- products 1:N order_items / reviews / wishlists / product_qna

[주문 중심]
- orders 1:N order_items (주문상세)
- orders 1:1 payments (결제)
- orders 1:1 shipments → 1:N shipment_tracking (배송추적)
- orders 1:N refunds / returns / exchanges

[판매자]
- sellers 1:N products
- sellers 1:N seller_settlements (정산)

[프로모션]
- promotions 1:N promotion_products
- coupons 1:N coupon_usage`,
			Tags: `["schema","ERD","관계","조인","테이블관계"]`},

		// ===== 예시 SQL =====
		{Category: "example_sql", Title: "지역별 고객 수 조회",
			Content: `-- 특정 지역 고객 수 (customer_addresses 기준)
SELECT COUNT(DISTINCT ca.customer_id)
FROM customer_addresses ca
WHERE ca.address LIKE '%대전%';

-- 시도별 고객 수 집계
SELECT
  SPLIT_PART(ca.address, ' ', 1) AS city,
  COUNT(DISTINCT ca.customer_id) AS customer_count
FROM customer_addresses ca
GROUP BY city
ORDER BY customer_count DESC;`,
			Tags: `["example_sql","고객","지역","주소","지역별"]`},

		{Category: "example_sql", Title: "등급별/성별 고객 통계",
			Content: `-- 등급별 고객 수
SELECT grade, COUNT(*) AS cnt FROM customers GROUP BY grade ORDER BY cnt DESC;

-- 성별 고객 수
SELECT gender, COUNT(*) FROM customers GROUP BY gender;

-- 등급별 평균 구매액
SELECT grade, AVG(total_spending)::INT AS avg_spending FROM customers GROUP BY grade ORDER BY avg_spending DESC;`,
			Tags: `["example_sql","고객","등급","성별","통계"]`},

		{Category: "example_sql", Title: "카테고리별/브랜드별 매출",
			Content: `-- 대카테고리별 매출
SELECT c.name AS category, SUM(oi.subtotal) AS total_sales
FROM categories c
JOIN products p ON c.id = p.category_id
JOIN order_items oi ON p.id = oi.product_id
JOIN orders o ON oi.order_id = o.id
WHERE c.depth = 0 AND o.status != 'CANCELLED'
GROUP BY c.name ORDER BY total_sales DESC;

-- 브랜드별 매출 TOP 10
SELECT b.name, SUM(oi.subtotal) AS revenue
FROM brands b JOIN products p ON b.id = p.brand_id
JOIN order_items oi ON p.id = oi.product_id
JOIN orders o ON oi.order_id = o.id WHERE o.status != 'CANCELLED'
GROUP BY b.name ORDER BY revenue DESC LIMIT 10;`,
			Tags: `["example_sql","카테고리","브랜드","매출"]`},

		{Category: "example_sql", Title: "베스트셀러/인기상품",
			Content: `-- 판매량 TOP 10
SELECT p.name, p.sales_count, p.rating_avg, p.review_count
FROM products p WHERE p.status = 'ACTIVE'
ORDER BY p.sales_count DESC LIMIT 10;

-- 카테고리별 베스트셀러
SELECT DISTINCT ON (p.category_id) p.category_id, c.name AS category, p.name, p.sales_count
FROM products p JOIN categories c ON p.category_id = c.id
ORDER BY p.category_id, p.sales_count DESC;`,
			Tags: `["example_sql","베스트셀러","인기상품","판매량","TOP"]`},

		{Category: "example_sql", Title: "월별 매출 추이",
			Content: `-- 월별 매출 (orders 기준)
SELECT TO_CHAR(ordered_at, 'YYYY-MM') AS month,
  COUNT(*) AS order_count, SUM(final_amount) AS revenue
FROM orders WHERE status != 'CANCELLED'
GROUP BY month ORDER BY month;

-- daily_sales 테이블 활용
SELECT sale_date, total_orders, total_revenue, net_revenue, new_customers
FROM daily_sales ORDER BY sale_date DESC LIMIT 30;`,
			Tags: `["example_sql","월별","매출","추이","일별"]`},

		{Category: "example_sql", Title: "결제수단별/PG사별 통계",
			Content: `-- 결제수단별 건수/금액
SELECT payment_method, COUNT(*) AS cnt, SUM(amount) AS total
FROM payments WHERE status = 'PAID'
GROUP BY payment_method ORDER BY total DESC;

-- PG사별 통계
SELECT pg_provider, COUNT(*) AS cnt, SUM(amount) AS total
FROM payments WHERE status = 'PAID'
GROUP BY pg_provider ORDER BY total DESC;`,
			Tags: `["example_sql","결제","결제수단","PG","카드","페이"]`},

		{Category: "example_sql", Title: "환불/반품/교환 분석",
			Content: `-- 환불 사유별 통계
SELECT reason, COUNT(*) AS cnt, SUM(amount) AS total_refund
FROM refunds GROUP BY reason ORDER BY cnt DESC;

-- 반품률 계산
SELECT
  (SELECT COUNT(*) FROM returns)::FLOAT / NULLIF((SELECT COUNT(*) FROM orders WHERE status = 'DELIVERED'), 0) * 100 AS return_rate_pct;

-- 교환 사유별
SELECT reason, COUNT(*) FROM exchanges GROUP BY reason ORDER BY count DESC;`,
			Tags: `["example_sql","환불","반품","교환","반품률"]`},

		{Category: "example_sql", Title: "검색어 분석",
			Content: `-- 인기 검색어 TOP 20
SELECT keyword, COUNT(*) AS search_count
FROM search_logs
GROUP BY keyword ORDER BY search_count DESC LIMIT 20;

-- 검색 후 클릭률
SELECT keyword, COUNT(*) AS searches,
  COUNT(clicked_product_id) AS clicks,
  ROUND(COUNT(clicked_product_id)::NUMERIC / COUNT(*) * 100, 1) AS click_rate
FROM search_logs
GROUP BY keyword HAVING COUNT(*) >= 10
ORDER BY search_count DESC LIMIT 20;`,
			Tags: `["example_sql","검색","인기검색어","키워드","클릭률"]`},

		{Category: "example_sql", Title: "판매자 정산/수수료 분석",
			Content: `-- 판매자별 총 매출/수수료
SELECT s.name, SUM(ss.total_sales) AS total_sales, SUM(ss.commission) AS total_commission, SUM(ss.net_amount) AS net_amount
FROM sellers s JOIN seller_settlements ss ON s.id = ss.seller_id
GROUP BY s.name ORDER BY total_sales DESC LIMIT 10;

-- 월별 정산 추이
SELECT settlement_period, SUM(total_sales), SUM(commission), SUM(net_amount)
FROM seller_settlements WHERE status = 'SETTLED'
GROUP BY settlement_period ORDER BY settlement_period;`,
			Tags: `["example_sql","판매자","정산","수수료"]`},

		{Category: "example_sql", Title: "고객 행동 분석 (페이지뷰, 디바이스)",
			Content: `-- 디바이스별 방문 비율
SELECT device_type, COUNT(*) AS views,
  ROUND(COUNT(*)::NUMERIC / (SELECT COUNT(*) FROM page_views) * 100, 1) AS pct
FROM page_views GROUP BY device_type ORDER BY views DESC;

-- 페이지 유형별 조회수
SELECT page_type, COUNT(*) FROM page_views GROUP BY page_type ORDER BY count DESC;

-- 유입 경로 분석
SELECT referrer, COUNT(*) FROM page_views WHERE referrer IS NOT NULL
GROUP BY referrer ORDER BY count DESC;`,
			Tags: `["example_sql","페이지뷰","디바이스","유입","트래픽"]`},

		{Category: "example_sql", Title: "쿠폰 사용 분석",
			Content: `-- 쿠폰별 사용 현황
SELECT c.name, c.discount_type, c.total_quantity, c.used_quantity,
  ROUND(c.used_quantity::NUMERIC / NULLIF(c.total_quantity, 0) * 100, 1) AS usage_rate
FROM coupons c ORDER BY c.used_quantity DESC LIMIT 10;

-- 쿠폰 할인 총액
SELECT SUM(discount_amount) AS total_discount FROM coupon_usage;`,
			Tags: `["example_sql","쿠폰","할인","사용률"]`},

		{Category: "example_sql", Title: "택배사별 배송 분석",
			Content: `-- 택배사별 배송 건수
SELECT carrier, COUNT(*) AS cnt FROM shipments GROUP BY carrier ORDER BY cnt DESC;

-- 택배사별 평균 배송 소요일
SELECT carrier,
  AVG(EXTRACT(EPOCH FROM (delivered_at - shipped_at)) / 86400)::NUMERIC(3,1) AS avg_days
FROM shipments WHERE delivered_at IS NOT NULL AND shipped_at IS NOT NULL
GROUP BY carrier ORDER BY avg_days;`,
			Tags: `["example_sql","택배","배송","소요일"]`},

		// ===== 용어집 =====
		{Category: "glossary", Title: "고객 등급 체계",
			Content: `등급(grade): BRONZE(기본,적립1%) → SILVER(30만이상,적립2%) → GOLD(100만이상,적립3%) → DIAMOND(300만이상,적립5%)`,
			Tags: `["glossary","등급"]`},

		{Category: "glossary", Title: "주문 상태 정의",
			Content: `주문상태(status): PENDING(대기) → CONFIRMED(확인) → SHIPPING(배송중) → DELIVERED(배송완료) / CANCELLED(취소)`,
			Tags: `["glossary","주문상태"]`},

		{Category: "glossary", Title: "결제수단 종류",
			Content: `결제수단(payment_method): CARD(신용/체크카드), BANK_TRANSFER(무통장입금), KAKAO_PAY(카카오페이), NAVER_PAY(네이버페이), TOSS_PAY(토스페이), PHONE(휴대폰결제), POINT(포인트결제)`,
			Tags: `["glossary","결제수단"]`},

		{Category: "glossary", Title: "상품 상태 / 배송 상태",
			Content: `상품상태: ACTIVE(판매중), SOLDOUT(품절), HIDDEN(숨김)
배송상태(shipments.status): PREPARING(준비중) → IN_TRANSIT(배송중) → DELIVERED(배달완료)`,
			Tags: `["glossary","상품상태","배송상태"]`},

		{Category: "glossary", Title: "포인트 타입 / 프로모션 타입",
			Content: `포인트타입: EARN(적립), USE(사용), EXPIRE(만료), ADMIN(관리자지급)
프로모션타입: SEASON_SALE(시즌세일), FLASH_SALE(타임딜), BUNDLE_DEAL(묶음할인), CATEGORY_SALE(카테고리할인), BRAND_SALE(브랜드할인)`,
			Tags: `["glossary","포인트","프로모션"]`},
	}

	ctx := context.Background()

	for i, k := range knowledges {
		text := k.Title + " " + k.Content
		resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
			Input:      []string{text},
			Model:      openai.SmallEmbedding3,
			Dimensions: 1536,
		})
		if err != nil {
			log.Fatalf("embedding error for %q: %v", k.Title, err)
		}
		if len(resp.Data) == 0 {
			log.Fatalf("no embedding for %q", k.Title)
		}
		k.Embedding = pgvector.NewVector(resp.Data[0].Embedding)
		knowledges[i] = k

		if err := db.WithContext(ctx).Create(&k).Error; err != nil {
			log.Fatalf("insert error for %q: %v", k.Title, err)
		}
		fmt.Printf("✓ [%d/%d] %s - %s\n", i+1, len(knowledges), k.Category, k.Title)
	}

	fmt.Printf("\nDone! %d knowledge entries seeded.\n", len(knowledges))
}
