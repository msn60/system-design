# Pagination

Pagination یعنی تقسیم یک Result Set بزرگ به بخش‌های کوچک‌تر و قابل‌کنترل، به‌گونه‌ای که در هر Request فقط بخشی از داده بارگذاری و ارسال شود. هدف اصلی، کاهش Database Work، Memory Usage، Response Size، Serialization Cost، Network Transfer و Initial Latency است.

بدون Pagination:

```
Client Request
→ Load All Matching Rows
→ Serialize Entire Result
→ Send Large Response
```

با Pagination:

```
Client Request
→ Load Limited Number of Rows
→ Return One Page
→ Client Requests Next Page
```

برای مثال، اگر جدول Product ده میلیون Row داشته باشد، منطقی نیست API همه آن‌ها را در یک Response برگرداند. به‌جای:

```
GET /products
```

می‌توان از یکی از این مدل‌ها استفاده کرد:

```
GET /products?page=1&page_size=50
```

یا:

```
GET /products?limit=50&cursor=...
```

## Mental Model

یک کتابخانه با یک میلیون کتاب را در نظر بگیرید. اگر کاربر بخواهد فهرست کتاب‌ها را ببیند، منطقی نیست فهرست تمام یک میلیون کتاب یک‌جا تولید شود. روش مناسب این است که ابتدا مثلاً پنجاه کتاب نمایش داده شود و در صورت نیاز، بخش بعدی درخواست شود.

```
Large Dataset
→ Small Ordered Chunks
→ Controlled Resource Usage
```

Pagination فقط یک قابلیت UI نیست؛ یک Performance Pattern است که Work هر Request را محدود می‌کند.

## اثر Pagination بر Performance

فرض کنید Query زیر یک میلیون Row برمی‌گرداند:

```
SELECT product_id, name, price
FROM products
ORDER BY product_id;
```

بدون Pagination:

```
Read 1M Rows from DB
→ Hold Large Result in Memory
→ Serialize 1M Objects
→ Transfer Large Payload
→ Client Parses Entire Response
```

با Pagination فقط تعداد محدودی Row خوانده می‌شود:

```
SELECT product_id, name, price
FROM products
ORDER BY product_id
FETCH FIRST 50 ROWS ONLY;
```

Pagination معمولاً این موارد را بهبود می‌دهد:

```
Time to First Result ↓
Response Payload ↓
Application Memory ↓
Client Memory ↓
Network Usage ↓
Serialization Cost ↓
```

البته اگر Client در نهایت تمام Pageها را بخواند، مجموع داده منتقل‌شده تقریباً همان خواهد بود. مزیت اصلی این است که هزینه به‌صورت تدریجی پرداخت می‌شود و بسیاری از Clientها اصلاً همه داده‌ها را مصرف نمی‌کنند:

```
Without Pagination
→ Pay Entire Cost Upfront

With Pagination
→ Pay Cost Incrementally and Only as Needed
```

## Offset-based Pagination

در Offset Pagination، Client شماره Page یا تعداد Rowهایی را که باید Skip شوند مشخص می‌کند.

مثلاً:

```
GET /products?page=3&page_size=20

Offset:
offset = (page - 1) × pageSize
```

برای Page سوم:

```
offset = (3 - 1) × 20 = 40

Query:
SELECT product_id, name, price
FROM products
ORDER BY product_id
OFFSET 40 ROWS
FETCH NEXT 20 ROWS ONLY;

Flow:
Skip First 40 Rows
→ Return Next 20 Rows
```

### مزایای Offset Pagination

```
Simple API
Easy UI Navigation
Supports Arbitrary Page Jump
Easy Total Page Calculation
Familiar User Experience
```

برای مثال Client می‌تواند مستقیم به Page 100 برود:

```
GET /products?page=100&page_size=20
```

## مشکل Large Offset

مشکل اصلی Offset Pagination این است که Database معمولاً باید Rowهای قبلی را پیدا و از آن‌ها عبور کند.

مثلاً:

```
OFFSET 999980 ROWS
FETCH NEXT 20 ROWS ONLY;
```

Database ممکن است عملاً این کار را انجام دهد:

```
Find / Walk Through 999,980 Rows
→ Discard Them
→ Return 20 Rows
```

با Page Size برابر پنجاه:

```
Page 1 → Skip 0
Page 100 → Skip 4,950
Page 10,000 → Skip 499,950
```

بنابراین:

```
Small Offset → Usually Acceptable
Large Offset → Increasing DB Work
```

حتی Index نیز این مشکل را کاملاً حذف نمی‌کند. Index می‌تواند Full Table Scan را کاهش دهد، اما Database ممکن است هنوز مجبور باشد تعداد زیادی Index Entry را طی کند:

```
Index Scan
→ Visit First 500K Entries
→ Discard Them
→ Return Next 50
```

### مثال عددی

فرض کنید:

```
Table Rows = 20 Million
Page Size = 50
Requested Page = 200,000
```

Offset تقریباً:

```
(200,000 - 1) × 50
≈ 10 Million Rows
```

Database برای برگرداندن فقط پنجاه Row باید از حدود ده میلیون Position عبور کند. به همین دلیل Offset برای Deep Pagination روی Datasetهای بزرگ مناسب نیست.

## Cursor-based / Keyset Pagination

در Cursor Pagination به‌جای گفتن «N Row را Skip کن»، موقعیت آخرین Row قبلی مشخص می‌شود.

اگر آخرین Product Page قبلی `product_id = 1250` باشد:

```
GET /products?limit=20&after_id=1250

Query:
SELECT product_id, name, price
FROM products
WHERE product_id > :after_id
ORDER BY product_id
FETCH FIRST 20 ROWS ONLY;

Flow:
Seek to product_id = 1250
→ Continue from Index
→ Read Next 20 Rows
```

### Mental Model

Offset Pagination شبیه این است که هر بار کتاب را از صفحه اول باز کنیم و هزاران خط جلو برویم:

```
Start from Beginning
→ Skip Many Rows
→ Reach Desired Position
```

Cursor Pagination شبیه Bookmark است:

```
Open at Bookmark
→ Continue Reading
```

### مزایای Cursor Pagination

```
Stable Performance for Deep Pagination
Efficient Index Seek
Suitable for Infinite Scroll
Lower Database Work
Better for Large Datasets
```

Database به‌جای پیمایش Rowهای قبلی، معمولاً Range Scan را از Cursor شروع می‌کند.

مدل هزینه مفهومی:

```
Cursor Pagination
≈ Index Seek + Page Size

Offset Pagination
≈ Offset + Page Size
```

### محدودیت Cursor Pagination

Cursor Pagination برای حرکت Sequential بسیار مناسب است، اما Jump مستقیم به Page دلخواه را به‌صورت طبیعی پشتیبانی نمی‌کند:

```
Page 1 → Cursor A
Page 2 → Cursor B
Page 3 → Cursor C

Trade-off:
Efficient Sequential Navigation
↔
Poor Arbitrary Page Jump
```

برای Feed، Timeline و Infinite Scroll این محدودیت معمولاً مشکلی نیست. برای Admin Tableهایی که کاربر باید مستقیم به Page مشخص برود، Offset ممکن است مناسب‌تر باشد.

## Stable Ordering

Pagination بدون Ordering پایدار قابل اعتماد نیست.

Query نامناسب:

```
SELECT product_id, name
FROM products
FETCH FIRST 20 ROWS ONLY;
```

بدون `ORDER BY`، Database تضمین نمی‌کند Rowها در Requestهای مختلف با ترتیب یکسان برگردند.

Pagination باید Ordering صریح داشته باشد:

```
ORDER BY product_id
```

حتی Ordering روی Column غیرUnique نیز کافی نیست.

مثلاً:

```
ORDER BY created_at
```

اگر چند Row `created_at` یکسان داشته باشند، ترتیب میان آن‌ها قطعی نیست.

روش بهتر:

```
ORDER BY created_at, product_id
```

در اینجا `product_id` نقش Tie-breaker را دارد.

اصل مهم:

```
Pagination Requires a Deterministic Total Order
```

## Tie-breaker

فرض کنید:

```
Product 10 → 10:00:00
Product 11 → 10:00:00
Product 12 → 10:00:00
```

اگر فقط `created_at` مبنای Sort باشد، ممکن است ترتیب این Rowها بین Requestها تغییر کند و Duplicate یا Missing Record ایجاد شود.

با:

```
ORDER BY created_at, product_id
```

ترتیب کامل و Deterministic می‌شود.

## Composite Cursor

اگر Ordering روی چند Column باشد، Cursor نیز باید همه Columnهای لازم را نگهداری کند.

مثلاً:

```
ORDER BY created_at DESC, product_id DESC
```

Cursor باید شامل:

```
created_at
product_id
```

باشد.

Query Page بعدی:

```
SELECT product_id, name, created_at
FROM products
WHERE
    created_at < :cursor_created_at
    OR (
        created_at = :cursor_created_at
        AND product_id < :cursor_product_id
    )
ORDER BY created_at DESC, product_id DESC
FETCH FIRST 20 ROWS ONLY;

Flow:
Earlier created_at
OR
Same created_at + Smaller product_id
```

برخی Databaseها از Row-value Comparison پشتیبانی می‌کنند:

```
WHERE (created_at, product_id) < (:created_at, :product_id)
```

اما Syntax و Support به Database بستگی دارد.

## Concurrent Writes و Offset Pagination

فرض کنید Page اول:

```
100
99
98
97
96
```

سپس Record جدید `101` Insert شود. ترتیب جدید:

```
101
100
99
98
97
96
95
...
```

اگر Page دوم با `OFFSET 5` گرفته شود، از `96` شروع می‌شود؛ بنابراین `96` دوباره دیده می‌شود:

```
Concurrent Insert
→ Offset Positions Shift
→ Duplicate Record
```

Delete نیز می‌تواند Missing Record ایجاد کند.

مثلاً اگر Page اول:

```
100
99
98
97
96
```

باشد و `99` حذف شود:

```
100
98
97
96
95
94
```

Page دوم با `OFFSET 5` از `94` شروع می‌شود و `95` از دست می‌رود:

```
Concurrent Delete
→ Position Shift
→ Missing Record
```

## Concurrent Writes و Cursor Pagination

اگر Page اول:

```
100
99
98
97
96
```

باشد، Cursor مقدار `96` را نگه می‌دارد.

Page بعدی:

```
WHERE product_id < 96
ORDER BY product_id DESC
```

اگر `101` Insert شود، Cursor جابه‌جا نمی‌شود:

```
Newer Record Inserted
→ Cursor Position Remains Stable
→ Continue Below 96
```

Cursor Pagination در برابر Shift شدن Positionها مقاوم‌تر است، اما Consistency کامل تضمین نمی‌کند. اگر Sort Key یک Row تغییر کند یا داده‌ای در میانه Range Insert شود، ممکن است رفتار متفاوتی مشاهده شود.

## Snapshot Consistency

اگر لازم باشد همه Pageها دقیقاً یک Snapshot ثابت را نمایش دهند، Pagination عادی کافی نیست.

مثلاً:

```
Page 1 at 10:00
Page 2 at 10:05
Page 3 at 10:10
```

در این فاصله Dataset تغییر کرده است. برای گزارش مالی، Audit یا Export ممکن است لازم باشد تمام Pageها یک Snapshot ثابت را ببینند.

راهکارهای ممکن:

```
Database Snapshot / Repeatable Read
Snapshot ID
As-of Timestamp
Versioned Dataset
Materialized Result
Export Job
```

مثلاً:

```
GET /transactions?limit=100&snapshot_id=abc123

Trade-off:
Strong Cross-Page Consistency
↔
More Storage and Operational Complexity
```

برای Feed اجتماعی معمولاً این سطح از Consistency لازم نیست، ولی برای Billing و Audit می‌تواند ضروری باشد.

## طراحی Cursor در API

بهتر است Cursor برای Client Opaque باشد:

```
GET /products?limit=20&cursor=eyJjcmVhdGVkX2F0IjoiLi4uIn0

Response:
{
  "items": [],
  "next_cursor": "eyJjcmVhdGVkX2F0IjoiLi4uIn0",
  "has_more": true
}
```

Opaque Cursor این مزایا را دارد:

```
Implementation Details Hidden
Composite Cursor Supported
Schema Can Evolve
Client Is Not Coupled to DB Columns
```

## Base64 به‌تنهایی کافی نیست

Base64 فقط Encoding است، نه Encryption و نه Integrity Protection. Client می‌تواند Cursor را Decode و تغییر دهد.

برای Cursor حساس می‌توان از این موارد استفاده کرد:

```
HMAC-signed Cursor
Encrypted Token
Server-Side Cursor ID
```

Cursor مناسب می‌تواند این اطلاعات را داشته باشد:

```
Version
Sort Position
Filter Context
Sort Context
Signature
```

## Cursor و Filter

Cursor یک Query نباید با Filter متفاوت استفاده شود.

مثلاً Cursor این Query:

```
GET /products?category=book&sort=created_at_desc
```

نباید برای این Query معتبر باشد:

```
GET /products?category=phone&cursor=...
```

در غیر این صورت:

```
Cursor from Query A
+ Filter from Query B
→ Inconsistent Pagination
```

Server باید Filter Context را داخل Cursor قرار دهد یا هنگام Decode آن را Validate کند.

## Cursor Versioning

ساختار Cursor ممکن است بعداً تغییر کند. بهتر است Version داشته باشد:

```
{
  "v": 1,
  "created_at": "...",
  "id": 123
}
```

Server می‌تواند Versionهای قدیمی را برای مدتی پشتیبانی کند یا Error مشخصی برای Version نامعتبر برگرداند.

## Page Size

Page Size باید Limit داشته باشد.

Request خطرناک:

```
GET /products?page_size=1000000
```

Policy نمونه:

```
Default Page Size = 20
Maximum Page Size = 100

Flow:
No Limit Provided → Use Default
Valid Limit → Accept
Limit Too Large → Clamp or Reject
```

### Trade-off Page Size

Page کوچک:

```
Lower Response Latency
Lower Memory
More Requests
More Round Trips
```

Page بزرگ:

```
Fewer Requests
Better Transfer Efficiency
Higher Response Time
Higher Memory
Higher Failure Cost
```

بنابراین:

```
Small Page
→ Better Responsiveness
→ More Request Overhead

Large Page
→ Better Throughput
→ Higher Per-Request Cost
```

Page Size باید بر اساس UI، Payload Size، Database Cost و Latency SLO تنظیم شود.

## تعیین `has_more`

روش رایج این است که یک Row بیشتر از Limit درخواست شود.

مثلاً:

```
Requested Limit = 20
Database Limit = 21
```

اگر 21 Row برگردد:

```
Return First 20
has_more = true
```

اگر حداکثر 20 Row برگردد:

```
has_more = false
```

مزیت:

```
No Separate COUNT Query Required
```

## Total Count

برخی APIها این اطلاعات را برمی‌گردانند:

```
{
  "items": [],
  "page": 3,
  "page_size": 20,
  "total_items": 1278342,
  "total_pages": 63918
}
```

اما `COUNT(*)` روی Dataset بزرگ یا Filter پیچیده ممکن است گران باشد:

```
Fetch One Page = Cheap
Exact COUNT = Potentially Expensive
```

گزینه‌های ممکن:

```
Exact Count
Approximate Count
Cached Count
Precomputed Count
No Count
```

برای Feed و Infinite Scroll معمولاً Count نیاز نیست:

```
{
  "items": [],
  "next_cursor": "...",
  "has_more": true
}
```

اصل مهم:

```
Do Not Compute Exact Count Unless the Product Actually Needs It
```

## Index Design

Pagination مؤثر به Index سازگار با Filter و Ordering نیاز دارد.

مثلاً:

```
SELECT product_id, created_at, name
FROM products
WHERE category_id = :category_id
  AND (
      created_at < :cursor_created_at
      OR (
          created_at = :cursor_created_at
          AND product_id < :cursor_product_id
      )
  )
ORDER BY created_at DESC, product_id DESC
FETCH FIRST 20 ROWS ONLY;
```

Index مناسب مفهومی:

```
(category_id, created_at DESC, product_id DESC)
```

قاعده عمومی:

```
Filter Prefix
→ Sort Columns
→ Tie-breaker
```

اگر Index با Filter و Sort هماهنگ نباشد، Database ممکن است تعداد زیادی Row یا Index Entry را بررسی کند.

## Covering Index

اگر Query فقط چند Column نیاز داشته باشد، Covering Index می‌تواند Table Lookup را کاهش دهد.

مزیت:

```
Faster Paginated Reads
```

هزینه:

```
More Index Storage
Higher Write Overhead
```

بنابراین:

```
Read Performance
↔
Storage and Write Cost
```

## Pagination روی Join

Pagination مستقیم روی Join ممکن است نتیجه Logical اشتباه بدهد.

مثلاً:

```
SELECT p.product_id, p.name, t.tag_name
FROM products p
JOIN product_tags pt ON ...
JOIN tags t ON ...
ORDER BY p.product_id;
```

هر Product ممکن است چند Tag داشته باشد:

```
20 SQL Rows
≠
20 Products
```

اگر `LIMIT 20` مستقیم روی Joined Result اعمال شود، ممکن است فقط چند Product واقعی برگردند.

راهکار مناسب:

```
First Page Primary Entity IDs
→ Load Related Data for Those IDs
```

مثلاً ابتدا:

```
SELECT product_id
FROM products
WHERE ...
ORDER BY product_id
FETCH FIRST 20 ROWS ONLY;
```

سپس Relationها به‌صورت Batch بارگذاری شوند.

## Pagination و N+1

Pagination تعداد Entityهای Page را محدود می‌کند، اما N+1 Problem را حل نمی‌کند.

```
Anti-pattern:
Load 20 Products
→ 20 Queries for Category
→ 20 Queries for Seller
```

روش مناسب:

```
Page Primary Entities
→ Batch Load Related Data
→ Assemble Response
```

بنابراین Pagination و Batching اغلب مکمل یکدیگرند.

## Pagination و Cache

Pageها را می‌توان Cache کرد:

```
products:category:10:page:1:size:20
```

اما Offset Pageها در Dataset متغیر سریع Invalid می‌شوند:

```
Insert at Beginning
→ Page 1 Changes
→ Page 2 Changes
→ Page 3 Changes
→ ...
```

به همین دلیل Cache کردن Pageهای ابتدایی، Filterهای محبوب و Datasetهای نسبتاً ثابت معمولاً ارزش بیشتری دارد.

Cursor Pageها نیز قابل Cache هستند، اما Cardinality Cursor بالا و Reuse پایین ممکن است Cache Efficiency را محدود کند.

## Reverse Pagination

برای Cursor Pagination ممکن است Client به `next` و `previous` نیاز داشته باشد:

```
{
  "next_cursor": "...",
  "previous_cursor": "..."
}

Forward:
WHERE id < :cursor
ORDER BY id DESC

Backward:
WHERE id > :cursor
ORDER BY id ASC
```

سپس Result ممکن است برای نمایش با ترتیب اصلی Reverse شود. Reverse Pagination یکی از نقاط افزایش Complexity در Cursor Pagination است.

## Pagination در Search Engine

Deep Pagination در Search Engine نیز پرهزینه است.

مدل Offset:

```
from = 1,000,000
size = 20
```

Search Engine ممکن است مجبور شود تعداد بسیار زیادی Result را Rank و نگهداری کند تا فقط تعداد کمی Result برگرداند.

راهکارهای رایج مفهومی:

```
Search After
Cursor
Point-in-Time Snapshot
Scroll-like API for Bulk Processing
```

برای User-facing Search، Cursor یا `search_after` معمولاً بهتر از Deep Offset است. برای Export کامل نیز Streaming یا Scroll-like Processing مناسب‌تر است.

## Pagination در Batch Job

Pagination فقط برای API نیست. برای پردازش کامل Dataset نیز کاربرد دارد.

مثلاً برای پردازش صد میلیون User:

روش نامناسب:

```
OFFSET 50000000 ROWS
FETCH NEXT 1000 ROWS ONLY;
```

روش مناسب:

```
WHERE user_id > :last_processed_id
ORDER BY user_id
FETCH FIRST 1000 ROWS ONLY;

Flow:
Read Page
→ Process Page
→ Persist Last Processed ID
→ Continue
```

مزیت مهم:

```
Checkpoint = Last Processed ID
```

اگر Job Crash کند:

```
Resume from Last Processed ID
```

برای Batch Jobهای بزرگ، Keyset Pagination معمولاً هم Performance بهتر و هم Resume Semantics ساده‌تری دارد.

## پیاده‌سازی Cursor Pagination در Go

```
Request:
type ListProductsRequest struct {
	Limit  int
	Cursor string
}

Response:
type ListProductsResponse struct {
	Items      []Product `json:"items"`
	NextCursor string    `json:"next_cursor,omitempty"`
	HasMore    bool      `json:"has_more"`
}
```

ساختار Cursor:

```
type productCursor struct {
	Version   int       `json:"v"`
	CreatedAt time.Time `json:"created_at"`
	ProductID int64     `json:"product_id"`
}

Encode:
func encodeCursor(cursor productCursor) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("marshal product cursor: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

Decode:
func decodeCursor(raw string) (productCursor, error) {
	var cursor productCursor

	data, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return cursor, fmt.Errorf("decode product cursor: %w", err)
	}

	if err := json.Unmarshal(data, &cursor); err != nil {
		return cursor, fmt.Errorf("unmarshal product cursor: %w", err)
	}

	if cursor.Version != 1 {
		return cursor, errors.New("unsupported cursor version")
	}

	return cursor, nil
}
```

این نمونه فقط Encoding انجام می‌دهد و در Production برای جلوگیری از Tampering باید Integrity Protection مانند HMAC اضافه شود.

## Service Flow در Go

```
Validate Limit
→ Decode Cursor
→ Query limit + 1
→ Determine has_more
→ Trim Extra Row
→ Build Next Cursor
→ Return
```

نمونه:

```
func (s *Service) ListProducts(
	ctx context.Context,
	req ListProductsRequest,
) (ListProductsResponse, error) {
	limit := normalizeLimit(req.Limit)

	var cursor *productCursor
	if req.Cursor != "" {
		decoded, err := decodeCursor(req.Cursor)
		if err != nil {
			return ListProductsResponse{}, ErrInvalidCursor
		}

		cursor = &decoded
	}

	products, err := s.repo.ListProducts(ctx, cursor, limit+1)
	if err != nil {
		return ListProductsResponse{}, fmt.Errorf("list products: %w", err)
	}

	hasMore := len(products) > limit
	if hasMore {
		products = products[:limit]
	}

	response := ListProductsResponse{
		Items:   products,
		HasMore: hasMore,
	}

	if hasMore && len(products) > 0 {
		last := products[len(products)-1]

		next, err := encodeCursor(productCursor{
			Version:   1,
			CreatedAt: last.CreatedAt,
			ProductID: last.ID,
		})
		if err != nil {
			return ListProductsResponse{}, fmt.Errorf("encode next cursor: %w", err)
		}

		response.NextCursor = next
	}

	return response, nil
}

Page Size Validation:
const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func normalizeLimit(limit int) int {
	switch {
	case limit <= 0:
		return defaultPageSize
	case limit > maxPageSize:
		return maxPageSize
	default:
		return limit
	}
}
```

می‌توان به‌جای Clamp کردن Limit بزرگ، Request را Reject کرد؛ مهم این است که Contract صریح باشد.

## Error Semantics

Errorهای رایج:

```
Invalid Cursor → 400 Bad Request
Unsupported Cursor Version → 400 Bad Request
Invalid Limit → 400 Bad Request
Source Timeout → 503 or 504
Permission Failure → 403
```

جزئیات داخلی Cursor نباید به Client Leak شوند:

```
Internal:
cursor signature mismatch

Client:
invalid cursor
```

## Failure Modeهای اصلی

### Ordering ناپایدار

```
No ORDER BY
or
Non-unique ORDER BY
→ Duplicate and Missing Items
```

### Large Offset

```
Deep Page
→ High DB Work
→ P95/P99 Increase
```

### Page Size نامحدود

```
Huge Limit
→ Memory Spike
→ Large Response
→ Timeout
```

### Cursor قابل دست‌کاری

```
Unsigned Cursor
→ Client Changes Position or Filter
```

### Cursor ناسازگار با Filter

```
Cursor from One Query
→ Used with Another Query
→ Incorrect Pagination
```

### Exact COUNT گران

```
Fast Page Query
+ Expensive COUNT
→ Slow Endpoint
```

### Pagination اشتباه روی Join

```
Limit Applied to Joined Rows
→ Fewer Logical Entities
→ Incomplete Result
```

### Mutable Sort Key

اگر Cursor براساس `updated_at` باشد و Row مرتب Update شود:

```
Row Moves Between Pages
→ Duplicate or Missing Result
```

برای Pagination پایدار بهتر است Sort Key تا حد امکان Immutable باشد یا Snapshot Semantics تعریف شود.

### Deleted Cursor Row

Keyset Pagination معمولاً به وجود Row مربوط به Cursor وابسته نیست:

```
WHERE id < :cursor_id
```

حتی اگر خود Row حذف شده باشد، Range همچنان معتبر است.

## Observability

Metricهای مهم:

```
pagination_request_total
pagination_page_size
pagination_result_size
pagination_query_duration_seconds
pagination_offset
pagination_cursor_decode_failure_total
pagination_invalid_limit_total
pagination_has_more_total
pagination_count_query_duration_seconds
```

برای Offset Pagination:

```
Offset Distribution
Deep Pagination Request Count
Latency by Offset Bucket
```

برای Cursor Pagination:

```
Cursor Validation Failure
Cursor Version Usage
Page Result Size
Empty Page Rate
```

Metricهای Database و Response:

```
DB Rows Examined
DB Rows Returned
Query P95/P99
Response Payload Size
Serialization Duration
```

یک نسبت بسیار مهم:

```
Rows Examined / Rows Returned
```

اگر برای برگرداندن پنجاه Row صدها هزار Row بررسی شوند، Query Plan، Index یا Pagination Strategy مناسب نیست.

برای جلوگیری از Metric Cardinality بالا، Cursor یا مقدار دقیق Offset نباید مستقیماً Label شوند. بهتر است از Bucket استفاده شود:

```
pagination_type="offset"
pagination_type="cursor"
offset_bucket="100k+"
page_size_bucket="51-100"
```

## Production Example: Social Feed

Feed کاربر براساس:

```
ORDER BY created_at DESC, post_id DESC
```

مرتب می‌شود.

Page اول:

```
GET /feed?limit=20
```

Page بعد:

```
GET /feed?limit=20&cursor=...

Query:
WHERE
    created_at < :created_at
    OR (
        created_at = :created_at
        AND post_id < :post_id
    )
ORDER BY created_at DESC, post_id DESC
FETCH FIRST 21 ROWS ONLY;
```

Cursor Pagination مناسب است، زیرا:

```
Dataset Changes Frequently
Users Navigate Sequentially
Deep Offset Would Be Expensive
Exact Total Count Is Unnecessary
```

Postهای جدید در ابتدای Feed اضافه می‌شوند، اما Position Cursor قبلی را جابه‌جا نمی‌کنند.

## Production Example: Admin Product Table

فرض کنید Admin نیاز دارد:

```
Go to Page 1
Jump to Page 50
See Total Pages
Sort by Name or Price
```

Offset Pagination ممکن است مناسب‌تر باشد، زیرا Random Page Navigation مهم است.

بااین‌حال باید محدودیت وجود داشته باشد:

```
Prevent Very Deep Pages
Require Narrow Filters
Use Proper Index
Possibly Switch to Cursor for Large Results
```

یک Hybrid Design نیز ممکن است:

```
Early Pages → Offset
Very Deep Navigation → Cursor or Refine Search
```

## Production Example: Background Migration

برای پردازش صد میلیون User:

روش نامناسب:

```
OFFSET 50000000 ROWS
FETCH NEXT 1000 ROWS ONLY;
```

روش مناسب:

```
WHERE user_id > :last_processed_id
ORDER BY user_id
FETCH FIRST 1000 ROWS ONLY;

Flow:
Load 1000 Users
→ Process
→ Persist Checkpoint
→ Continue
```

در صورت Crash:

```
Resume from Last Processed ID
```

Keyset Pagination در این سناریو هم Work Database را کنترل می‌کند و هم Resume را ساده می‌سازد.

## Review Scenario: Transaction API

فرض کنید:

```
Transactions = 200 Million
Traffic = 10K Requests/sec
Default Page Size = 50
Current API Uses page and page_size
Users Mostly Read First 5 Pages
Some Clients Request Page 100,000
Ordering = created_at DESC
Many Transactions Share the Same Timestamp
New Transactions Are Continuously Inserted
API Runs COUNT(*) for Every Request
```

### مشکل Deep Pagination

Page 100,000 با Page Size پنجاه:

```
Offset ≈ 5 Million Rows
```

Database باید میلیون‌ها Position را طی کند تا پنجاه Row برگرداند.

### مشکل Ordering

`created_at` به‌تنهایی Unique نیست. Ordering باید Tie-breaker داشته باشد:

```
ORDER BY created_at DESC, transaction_id DESC
```

### مشکل Concurrent Insert

Insertهای جدید Positionهای Offset را جابه‌جا می‌کنند و می‌توانند Duplicate یا Missing Transaction ایجاد کنند.

### مشکل COUNT

اجرای `COUNT(*)` در هر Request ممکن است حتی از Query Page گران‌تر باشد.

### طراحی مناسب

برای Clientهایی که Sequential Navigation دارند:

```
Cursor Pagination
+ Composite Cursor(created_at, transaction_id)
+ limit + 1
+ No Exact Count by Default

Query:
WHERE
    created_at < :cursor_created_at
    OR (
        created_at = :cursor_created_at
        AND transaction_id < :cursor_transaction_id
    )
ORDER BY created_at DESC, transaction_id DESC
FETCH FIRST 51 ROWS ONLY;

Index:
(created_at DESC, transaction_id DESC)
```

اگر Query به Account محدود باشد:

```
(account_id, created_at DESC, transaction_id DESC)
```

### Total Count در صورت نیاز

گزینه‌ها:

```
Run Count Only on Explicit Request
Cache Count
Use Approximate Count
Precompute Count
Restrict Count to Narrow Filters
```

### Metricهای حیاتی

```
Query P95/P99
Rows Examined / Rows Returned
Offset Distribution
Deep Page Request Rate
Count Query Duration
Response Size
Cursor Error Rate
Duplicate/Missing Incident Reports
```

## نکات Production-oriented

در طراحی واقعی باید پاسخ این پرسش‌ها روشن باشد:

```
What Is the Stable Sort Order?
Is the Sort Key Unique?
Is the Sort Key Mutable?
Offset or Cursor?
What Is the Maximum Page Size?
Is Deep Pagination Allowed?
Is Exact Total Count Required?
How Are Filters Bound to the Cursor?
Is the Cursor Signed?
What Index Supports the Query?
What Happens under Concurrent Inserts and Deletes?
Is Snapshot Consistency Required?
```

Pagination فقط افزودن `LIMIT` یا `OFFSET` به Query نیست؛ API Contract، Ordering، Index Design، Consistency و Cursor Security همگی بخشی از طراحی آن هستند.

## نکات Interview-oriented

مهم‌ترین مقایسه:

```
Offset Pagination
+ Simple
+ Supports Arbitrary Page Jump
+ Easy Total Pages
- Slow for Deep Pages
- Unstable under Concurrent Writes

Cursor Pagination
+ Efficient for Large Datasets
+ Stable Sequential Navigation
+ Suitable for Feeds and Infinite Scroll
- More Complex
- Hard to Jump to Arbitrary Page
- Requires Stable Indexed Ordering
```

در پاسخ مصاحبه‌ای باید به این موارد نیز اشاره شود:

```
Unique Tie-breaker
Composite Cursor
Page Size Limit
limit + 1 for has_more
Avoid Unnecessary COUNT
Index Matching Filter and Sort
```

## جمع‌بندی و Mental Model نهایی

Pagination یک Result Set بزرگ را به بخش‌های کوچک، Ordered و قابل‌کنترل تقسیم می‌کند:

```
Large Result Set
→ Stable Ordered Pages
→ Lower Per-Request Memory and Network Cost
→ Faster Initial Response

Offset Pagination:
Start from Beginning
→ Skip N Rows
→ Return Page

Cursor Pagination:
Start from Last Known Position
→ Continue from Index
→ Return Next Page
```

مدل تصمیم‌گیری:

```
Small and Relatively Stable Dataset
+ Need Arbitrary Page Jump
→ Offset Pagination

Large or Frequently Changing Dataset
+ Sequential Navigation
→ Cursor / Keyset Pagination
```

مدل Production مناسب:

```
Define Stable Total Order
→ Add Unique Tie-breaker
→ Select Offset or Cursor
→ Enforce Page Size
→ Build Matching Index
→ Handle Concurrent Changes
→ Avoid Expensive Count Unless Required
→ Observe Rows Examined, Payload and P99
```

مدل ذهنی نهایی این است: **Pagination فقط تقسیم خروجی به چند Page نیست؛ بلکه روشی برای محدود کردن Work هر Request است و زمانی درست عمل می‌کند که Ordering پایدار، Index مناسب و رفتار مشخصی در برابر تغییر هم‌زمان داده وجود داشته باشد.**
