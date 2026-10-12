# Cache Failure Mode: Cache Penetration

پس از بررسی Cache Stampede، Cache Avalanche و Hot Key، به Failure Mode بعدی یعنی Cache Penetration می‌رسیم. Cache Penetration زمانی رخ می‌دهد که Requestها برای داده‌هایی ارسال شوند که اساساً در Source اصلی، معمولاً Database، وجود ندارند. چون Cache معمولاً فقط داده‌های موجود را نگه می‌دارد، هر Request ابتدا Cache Miss می‌شود، سپس Database را Query می‌کند و در نهایت نتیجه `Not Found` دریافت می‌کند. اگر این نتیجه Cache نشود، Request بعدی دوباره همان مسیر را تکرار می‌کند:

```
Request for Nonexistent Key → Cache Miss → Database Miss → Return Not Found → Repeat
```

برای مثال، اگر Request زیر برای Productی ارسال شود که وجود ندارد:

```
GET /products/999999999
```

Cache مقداری برای آن ندارد و Database نیز Rowای پیدا نمی‌کند. اگر هزاران Request برای همین ID یا برای IDهای نامعتبر مختلف ارسال شوند، Cache عملاً نمی‌تواند Database Load را کاهش دهد.

## Mental Model

Cache را مانند مسئول پذیرشی در نظر بگیرید که فقط اطلاعات افراد ثبت‌شده را نگه می‌دارد. اگر مراجعه‌کننده نام فردی را بپرسد که وجود ندارد، مسئول پذیرش هر بار مجبور می‌شود به بایگانی اصلی مراجعه کند و دوباره بررسی کند.

مدل نامطلوب:

```
Unknown Name → Reception Has No Record → Check Archive → Not Found
Unknown Name Again → Check Archive Again → Not Found
```

مدل مطلوب:

```
Unknown Name → Check Archive Once → Remember "Not Found" Temporarily
```

بنابراین در Cache Penetration مشکل این نیست که داده موجود از Cache افتاده است؛ مشکل این است که سیستم برای داده‌ای که هرگز وجود نداشته، بارها به Source مراجعه می‌کند.

## Problem Statement

فرض کنید سرویس Product این مشخصات را دارد:

```
Valid Products = 10 Million
Traffic = 50K Requests/sec
Database Capacity = 8K Queries/sec
```

در شرایط عادی بیشتر Requestها برای Productهای موجود هستند و Cache بخش بزرگی از Traffic را پاسخ می‌دهد. اما یک Bot یا Client معیوب شروع به ارسال IDهای تصادفی می‌کند:

```
GET /products/927381928371
GET /products/712837182731
GET /products/111111111111
```

چون این Productها وجود ندارند، مسیر هر Request چنین خواهد بود:

```
Request → Cache Miss → DB Query → No Rows → 404
```

اگر نرخ Requestهای نامعتبر به 20 هزار Request در ثانیه برسد:

```
20K Invalid Requests/sec → 20K DB Misses/sec
```

در حالی که Database فقط 8 هزار Query در ثانیه ظرفیت امن دارد. در نتیجه با وجود سالم بودن Cache، Database Saturate می‌شود.

## چرا Cache در این وضعیت بی‌اثر می‌شود؟

Cache معمولاً فقط مقادیر موجود را نگه می‌دارد:

```
product:123 → Product Data
product:456 → Product Data
```

برای داده ناموجود هیچ Entry ثبت نمی‌شود:

```
product:999999 → Nothing in Cache
```

در نتیجه هر Request جدید همان Cache Miss را تکرار می‌کند. این وضعیت ممکن است به‌دلیل Client Bug، Bot Traffic، Enumeration Attack، Random ID Scanning، لینک‌های قدیمی یا شکسته، داده حذف‌شده، Race میان Create و Read، Requestهای malformed، ورودی بدون Validation یا حمله هدفمند برای فشار روی Database ایجاد شود.

## تفاوت Cache Penetration با سایر Failure Modeها

در Cache Stampede، داده معمولاً وجود دارد، اما Cache Entry آن منقضی یا حذف شده است:

```
Existing Data + Missing Cache → Many Requests Load Same Source Value
```

در Cache Penetration، داده اساساً وجود ندارد:

```
Nonexistent Data → Cache Miss + DB Miss Repeatedly
```

در Stampede، `singleflight` بسیار مؤثر است، چون Requestها معمولاً نتیجه یک Load مشترک را می‌خواهند. در Penetration، اگر Requestها برای میلیون‌ها ID متفاوت باشند، `singleflight` کمک محدودی می‌کند.

در Cache Avalanche، تعداد زیادی Key موجود تقریباً هم‌زمان از Cache خارج می‌شوند:

```
Many Existing Keys Lost → Broad DB Load
```

اما در Penetration، Keyها معمولاً از ابتدا وجود نداشته‌اند:

```
Many Invalid Keys Requested → Repeated DB Misses
```

راهکارهای Avalanche بیشتر شامل TTL Jitter، Cache Warming، HA و Source Limiting هستند، در حالی که Penetration بیشتر با Negative Caching، Bloom Filter، Validation و Rate Limiting کنترل می‌شود.

در Hot Key، یک Key موجود Traffic بسیار زیادی دارد:

```
Existing Hot Key → Huge Cache Traffic
```

اما در Penetration، Requestها برای Keyهای ناموجود هستند:

```
Nonexistent Keys → Cache Miss → DB Miss
```

ممکن است یک Key ناموجود خودش Hot شود؛ برای مثال هزاران Request برای `product:999999` ارسال شوند. در این حالت Negative Caching بسیار مؤثر است.

## اثر روی Database

Query ناموفق نیز برای Database هزینه دارد. حتی اگر هیچ Row پیدا نشود، Database ممکن است Query را Parse کند، Connection بگیرد، Index را بررسی کند، در B-Tree حرکت کند، Buffer Page بخواند، Visibility را کنترل کند و نتیجه خالی برگرداند.

```
Parse Query → Acquire Connection → Traverse Index → Check Rows → Return Empty Result
```

اگر Index مناسب وجود نداشته باشد، هزینه بیشتر می‌شود و حتی ممکن است Table Scan رخ دهد. اثرهای رایج عبارت‌اند از:

```
DB Query Rate ↑
Connection Pool Wait ↑
CPU Usage ↑
Index I/O ↑
P95/P99 Latency ↑
Timeout ↑
Retry ↑
Legitimate Traffic Slows Down
```

بنابراین Query ناموفق الزاماً ارزان نیست.

## اثر روی Connection Pool

فرض کنید:

```
DB Pool Size = 100
Invalid Requests = 10K/sec
DB Query Duration = 20ms
```

حتی Queryهای کوتاه نیز می‌توانند Pool را پر کنند:

```
Invalid Requests → Pool Saturation → Legitimate Requests Wait
```

در نتیجه کاربران واقعی که Product موجود درخواست می‌کنند نیز با Latency بیشتر مواجه می‌شوند.

## Negative Caching

اصلی‌ترین راهکار Cache Penetration این است که نتیجه `Not Found` نیز برای مدت کوتاهی Cache شود:

```
product:999999 → NOT_FOUND → TTL 30 Seconds

Flow:
First Request → Cache Miss → DB Miss → Cache NOT_FOUND
Next Requests → Cache Hit NOT_FOUND → Return 404 Without DB
```

این روش Negative Caching نام دارد.

یک نمونه ساده در Go:

```
var ErrProductNotFound = errors.New("product not found")

type CachedProduct struct {
    Product  Product
    NotFound bool
}

func GetProduct(ctx context.Context, id string) (Product, error) {
    key := "product:" + id

    cached, err := cache.Get(ctx, key)
    if err == nil {
        if cached.NotFound {
            return Product{}, ErrProductNotFound
        }

        return cached.Product, nil
    }

    product, err := db.GetProduct(ctx, id)
    if errors.Is(err, sql.ErrNoRows) {
        negative := CachedProduct{
            NotFound: true,
        }

        if cacheErr := cache.Set(ctx, key, negative, 30*time.Second); cacheErr != nil {
            // Record the cache failure without hiding the Not Found result.
        }

        return Product{}, ErrProductNotFound
    }
    if err != nil {
        return Product{}, err
    }

    positive := CachedProduct{
        Product: product,
    }

    if cacheErr := cache.Set(ctx, key, positive, 5*time.Minute); cacheErr != nil {
        // Record the cache failure without hiding the DB result.
    }

    return product, nil
}
```

باید بین `Cache Miss` و `Cached Not Found` تفاوت صریح وجود داشته باشد. اگر هر دو با `nil` یا Zero Value نمایش داده شوند، رفتار مبهم می‌شود.

## TTL مناسب برای Negative Cache

TTL نتیجه منفی معمولاً باید کوتاه‌تر از TTL داده موجود باشد:

```
Positive Cache TTL = 10 Minutes
Negative Cache TTL = 30 Seconds
```

دلیل این تفاوت آن است که ممکن است داده کمی بعد ایجاد شود:

```
10:00:00 → Request product:123 → Not Found → Negative Cache
10:00:05 → Product 123 Created
10:00:10 → Request product:123
```

اگر Negative TTL بسیار طولانی باشد، Request همچنان `404` می‌گیرد، با اینکه داده ایجاد شده است. بنابراین:

TTL خیلی کوتاه → DB Protection کم

TTL خیلی بلند → False Not Found Staleness زیاد

اگر سیستم می‌داند داده چه زمانی ایجاد می‌شود، بهتر است هنگام Create، Negative Cache حذف یا با مقدار واقعی جایگزین شود:

```
Create Product 123 → Delete Negative Cache
```

یا:

```
Create Product 123 → Write Positive Cache
```

این روش Consistency را بهتر می‌کند و اجازه می‌دهد TTL منفی کمی طولانی‌تر باشد.

## خطر Cache Pollution

اگر مهاجم میلیون‌ها ID تصادفی ارسال کند و همه آن‌ها Negative Cache شوند، Cache با Entryهای کم‌ارزش پر می‌شود:

```
1M Invalid IDs → 1M Negative Cache Entries
```

نتیجه می‌تواند افزایش Memory Pressure، Eviction داده‌های مفید و حتی افزایش خطر Cache Avalanche باشد:

```
Memory Pressure ↑ → Useful Keys Evicted → Miss Rate ↑
```

بنابراین Negative Caching به‌تنهایی کافی نیست و باید TTL کوتاه، Memory Limit، Admission Policy یا Rate Limiting نیز وجود داشته باشد.

Negative Caching زمانی بیشترین اثر را دارد که Requestهای تکراری برای همان Key ناموجود ارسال شوند:

```
100K Requests for product:999999 → One DB Query + Negative Cache Hits
```

اما اگر هر Request یک ID تصادفی و جدید داشته باشد، تقریباً هر Request یک Cache Miss جدید ایجاد می‌کند. در این شرایط Bloom Filter و Validation مؤثرتر هستند.

## Input Validation

قبل از مراجعه به Cache یا Database باید ورودی بررسی شود. اگر Product ID باید عددی مثبت با حداکثر طول مشخص باشد، Requestهای زیر باید در همان Application Layer رد شوند:

```
/products/-1
/products/abc
/products/999999999999999999999
```

Flow مناسب:

```
Invalid Input → Return 400 → No Cache Lookup → No DB Query
```

نمونه Go:

```
func parseProductID(raw string) (int64, error) {
    id, err := strconv.ParseInt(raw, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid product id: %w", err)
    }

    if id <= 0 {
        return 0, errors.New("product id must be positive")
    }

    return id, nil
}
```

Validation ساده می‌تواند حجم زیادی از Traffic بی‌ارزش را پیش از ورود به Cache Layer حذف کند.

فقط Syntax کافی نیست. ممکن است ID از نظر فرمت معتبر باشد، اما خارج از Range منطقی سیستم قرار بگیرد. اگر بزرگ‌ترین Product ID حدود 50 میلیون است، Requestی مانند زیر مشکوک است:

```
GET /products/999999999999
```

می‌توان Range منطقی را بررسی کرد، اما این Optimization باید با احتیاط استفاده شود؛ زیرا Range ممکن است تغییر کند یا IDها Sequence ساده نباشند.

## Bloom Filter

Bloom Filter یک ساختار احتمالی کم‌حافظه است که بررسی می‌کند آیا یک عنصر احتمالاً در مجموعه وجود دارد یا قطعاً وجود ندارد.

```
Flow:
Request Key → Bloom Filter
Definitely Not Present → Return Not Found Without DB
Maybe Present → Check Cache → Database if Needed
```

ویژگی اصلی Bloom Filter:

False Positive ممکن است

False Negative در طراحی صحیح نباید وجود داشته باشد

ممکن است Bloom Filter بگوید Key احتمالاً وجود دارد، در حالی که وجود ندارد؛ در این صورت فقط یک DB Query اضافی انجام می‌شود. اما اگر بگوید Key قطعاً وجود ندارد، می‌توان Request را بدون مراجعه به Database رد کرد.

Mental Model آن مانند نگهبانی است که نمی‌تواند همیشه بگوید چه کسی داخل ساختمان است، اما می‌تواند با اطمینان بگوید بعضی افراد قطعاً داخل نیستند:

```
Definitely Not Present → Reject
Maybe Present → Continue Checking
```

Flow کامل:

```
Request → Validate Input → Bloom Filter → Cache Lookup → Database if Needed
```

اگر Bloom Filter نتیجه `Definitely Absent` بدهد، Service مستقیماً `404` برمی‌گرداند. اگر نتیجه `Maybe Present` باشد، Cache و سپس Database بررسی می‌شوند.

برای Datasetهای بزرگ، Bloom Filter نسبت به نگهداری تمام IDها در Set معمولی Memory بسیار کمتری نیاز دارد، اما Trade-offهایی دارد:

Memory کمتر

Lookup سریع

False Positive محدود

Deletion دشوار

نیاز به Rebuild یا Counting Bloom Filter

## به‌روزرسانی Bloom Filter

وقتی داده جدید ایجاد می‌شود، Bloom Filter نیز باید Update شود:

```
Create Product → Commit DB → Add Product ID to Bloom Filter
```

اگر Update Bloom Filter fail شود و سیستم با قطعیت نتیجه `Not Present` بدهد، ممکن است داده موجود اشتباهاً ناموجود تشخیص داده شود. به همین دلیل در حالت عدم اطمینان بهتر است Request به Database عبور کند.

Bloom Filter استاندارد حذف را به‌خوبی پشتیبانی نمی‌کند، چون Bitها میان عناصر مختلف مشترک هستند. راهکارهای رایج شامل Periodic Rebuild، Counting Bloom Filter یا پذیرش False Positive پس از Delete است. False Positive پس از حذف معمولاً خطرناک نیست؛ فقط یک DB Query اضافی ایجاد می‌کند.

Bloom Filter می‌تواند Local در هر Application Instance، Shared در یک Service، داخل Redis Module یا به‌صورت Snapshot فایل نگهداری شود. Bloom Filter محلی Latency کمی دارد، اما Update و Versioning میان Instanceها سخت‌تر است. Bloom Filter مشترک Consistency بهتری دارد، اما Network Dependency جدید ایجاد می‌کند.

## Rate Limiting

اگر Traffic نامعتبر از Client، IP، Token یا Tenant مشخصی می‌آید، Rate Limiting باید قبل از Cache و Database اعمال شود:

```
Client → Rate Limit → Validation → Cache → DB
```

برای مثال:

```
Maximum 100 Product Lookups/sec per Client
```

Rate Limiting در برابر Botها، Enumeration و Abuse مفید است.

گاهی می‌توان Clientهایی را که درصد بسیار بالایی از Requestهایشان `404` است، محدود کرد:

```
Client A: 1000 Requests → 950 Not Found → Suspicious
```

این رفتار به Adaptive Rate Limiting یا Abuse Detection نزدیک است، اما باید مراقب False Positive بود؛ بعضی Clientهای واقعی ممکن است داده‌های حذف‌شده یا لینک‌های قدیمی را Query کنند.

## Bot و Attack Protection

Cache Penetration می‌تواند حمله‌ای عمدی باشد:

```
Random IDs → Force Cache Miss → Exhaust Database
```

لایه‌های دفاعی می‌توانند شامل WAF، API Gateway Rate Limit، IP Reputation، Authentication، Request Quotas، Captcha در Flowهای عمومی و Anomaly Detection باشند. Cache Layer نباید تنها دفاع سیستم در برابر Abuse باشد.

## Request Coalescing

اگر Requestهای زیادی برای همان Key ناموجود برسند، `singleflight` می‌تواند DB Queryهای تکراری را یکی کند:

```
1000 Requests for missing product:999 → One DB Query → Shared Not Found
```

سپس نتیجه منفی Cache می‌شود. اما اگر Keyها متفاوت باشند، `singleflight` تأثیر محدودی دارد.

## Source-Side Protection

Database باید مستقل از Cache محافظت شود:

```
Connection Pool Limit
Query Timeout
Concurrency Limit
Rate Limit
Circuit Breaker
Load Shedding
Read Replica
```

اگر Bloom Filter یا Negative Cache fail کند، Database نباید کل Traffic نامعتبر را بدون محدودیت بپذیرد.

اصل مهم:

```
Invalid Traffic Must Not Consume Unlimited Source Capacity
```

## Query Optimization و Indexing

چون Cache Penetration باعث تعداد زیادی Query با نتیجه `Not Found` می‌شود، Index مناسب اهمیت زیادی دارد:

```
SELECT *
FROM products
WHERE product_id = :id;
```

باید روی `product_id` Index مناسب وجود داشته باشد. بدون Index، هر Request نامعتبر ممکن است Table Scan ایجاد کند.

بااین‌حال، Index فقط هزینه هر Query را کم می‌کند و مشکل تعداد Queryها را حل نمی‌کند:

```
Indexing = Damage Reduction
Negative Cache or Bloom Filter = Query Prevention
```

## Cache کردن Range یا Metadata

گاهی می‌توان Metadata ساده‌ای Cache کرد که محدوده معتبر را نشان دهد:

```
min_product_id
max_product_id
active_tenant_ids
valid_category_ids
```

Requestهای خارج از محدوده می‌توانند بدون DB رد شوند. این روش برای IDهای Sequential یا Domainهای محدود مناسب است، اما برای UUID یا شناسه‌های پراکنده کاربرد کمتری دارد.

## UUID و Cache Penetration

در سیستم‌هایی که UUID استفاده می‌شود، Range Validation ممکن نیست. مهاجم می‌تواند UUIDهای معتبر از نظر Format ولی ناموجود تولید کند:

```
GET /users/550e8400-e29b-41d4-a716-446655440000
```

در این حالت Validation فقط Format را بررسی می‌کند. برای کاهش Database Load، Negative Cache، Bloom Filter، Rate Limiting و Authentication اهمیت بیشتری دارند.

## Tenant-Aware Validation

در سیستم Multi-Tenant باید قبل از DB Query بررسی شود که Resource ID در Scope Tenant فعلی قرار دارد.

Flow نامطلوب:

```
Tenant A Requests Arbitrary IDs → Global DB Lookup
```

Flow مناسب:

```
Authenticate Tenant → Validate Scope → Query by Tenant + ID
```

این موضوع علاوه بر Performance، برای Security نیز ضروری است.

## Race میان Create و Negative Cache

یکی از مشکلات مهم Negative Caching:

```
Request A → DB Miss → Negative Cache
Request B → Create Data
Request C → Negative Cache Hit → Incorrect 404
```

راهکارهای اصلی:

```
Invalidate Negative Cache on Create
Short Negative TTL
Versioned Events
Write Positive Cache after Create
```

در سیستم Eventual Consistent ممکن است مقدار کمی False Not Found قابل قبول باشد، اما در سیستم‌های حساس باید Invalidation دقیق‌تری وجود داشته باشد.

جهت معکوس نیز ممکن است رخ دهد:

```
Data Deleted → Positive Cache Still Exists → Stale Success
```

این مسئله مستقیماً Cache Penetration نیست، اما نشان می‌دهد Positive و Negative Cache هر دو به Invalidation Strategy نیاز دارند.

## Temporary Error را Negative Cache نکنیم

نباید هر Error از Database را به `Not Found` تبدیل و Cache کرد. خطاهایی مانند DB Timeout، Connection Error، Permission Error یا Serialization Error به معنای عدم وجود داده نیستند.

```
Anti-pattern:
Any DB Error → Cache NOT_FOUND
```

این رفتار می‌تواند Outage موقت Database را به هزاران `404` اشتباه تبدیل کند.

فقط خطای قطعی عدم وجود داده باید Negative Cache شود:

```
if errors.Is(err, sql.ErrNoRows) {
    // Safe candidate for negative caching.
}
```

## Status Code مناسب

ورودی نامعتبر، داده ناموجود، Rate Limit و Source Failure باید از هم تفکیک شوند:

```
Malformed ID → 400 Bad Request
Valid ID but Missing Data → 404 Not Found
Rate Limit Exceeded → 429 Too Many Requests
Source Unavailable → 503 Service Unavailable
```

اگر همه خطاها `404` شوند، Observability و Abuse Detection ضعیف می‌شود.

## TTL Jitter برای Negative Cache

اگر تعداد زیادی Negative Entry هم‌زمان ساخته شوند و TTL یکسان داشته باشند، ممکن است هم‌زمان Expire شوند و موج جدیدی از DB Missها ایجاد کنند:

```
100K Invalid IDs Cached at 10:00
Negative TTL = 1 Minute
→ 100K Entries Expire around 10:01
```

برای پخش زمان Expiration می‌توان Jitter اضافه کرد:

```
Negative TTL = 30s ± 10s
```

این موضوع زمانی مفید است که Bot همان IDها را تکرار کند.

## Memory Management

Negative Cache باید Memory Budget مشخص داشته باشد. گزینه‌های معمول عبارت‌اند از:

```
Short TTL
Maximum Entry Count
Separate Cache Namespace
Dedicated Memory Quota
Admission Policy
Approximate LFU or LRU
```

گاهی بهتر است Negative Cache از Positive Cache جدا شود تا Entryهای نامعتبر موجب Eviction داده‌های مفید نشوند:

```
Positive Cache → Shared Redis
Negative Cache → Small Local Cache or Separate Namespace
```

مزیت این مدل آن است که Invalid Traffic نمی‌تواند به‌سادگی داده‌های ارزشمند را از Cache خارج کند؛ عیب آن نیز افزایش Complexity زیرساخت و منطق Application است.

Negative Cache محلی با TTL کوتاه نیز می‌تواند برای Requestهای تکراری روی همان Instance مفید باشد:

```
Request → L1 Negative Cache → Shared Cache → DB
```

در چند Instance ممکن است هر Instance یک DB Query انجام دهد، اما همچنان بسیار بهتر از هزاران Query است.

## Go Implementation Model

یک Flow مناسب در Go:

```
Request → Parse and Validate Input → Apply Rate Limit → Check Bloom Filter → Check Positive/Negative Cache → Enter singleflight by Key → Recheck Cache → Query DB with Timeout → Cache Positive or Negative Result → Return
```

اجزای مهم:

```
context.Context
Typed Error Classification
Separate Positive and Negative TTL
singleflight
TTL Jitter
Rate Limit
Source Timeout
Metrics
```

برای صریح کردن وضعیت Cache بهتر است از مدل Type-safe استفاده شود:

```
type CacheState uint8

const (
    CacheStateValue CacheState = iota + 1
    CacheStateNotFound
)

type CachedResult[V any] struct {
    State CacheState
    Value V
}
```

این مدل از ابهام میان Zero Value و `Not Found` جلوگیری می‌کند.

Errorها نیز باید دسته‌بندی شوند:

```
Not Found → Negative Cache Candidate
Invalid Input → Reject Immediately
Timeout → Do Not Negative Cache
Temporary DB Error → Retry or Fallback Policy
Permission Error → Return Authorization Error
Decode Error → Observe and Repair Cache
```

این تفکیک برای جلوگیری از Cache کردن نتایج اشتباه حیاتی است.

## Observability

Metricهای مهم عبارت‌اند از:

```
cache_negative_hit_total
cache_negative_miss_total
negative_cache_write_total
negative_cache_entry_count
invalid_request_total
not_found_total
bloom_filter_reject_total
bloom_filter_maybe_total
db_not_found_query_total
rate_limited_total
distinct_missing_keys
missing_key_repeat_rate
```

همچنین باید Database Query Rate، Connection Pool Wait، نرخ `404`، نرخ `400`، نرخ `429`، Top Missing Key Patternها و Clientهای مشکوک مانیتور شوند.

نباید تمام IDهای ناموجود به‌عنوان Label در Metric ثبت شوند، چون Cardinality بسیار بالا می‌رود:

```
missing_key="999999999"
```

روش‌های مناسب‌تر شامل Top-K Sampling، Hashed Buckets، Endpoint Label، Tenant Label، Key Pattern و Approximate Distinct Count هستند. برای Log نیز Sampling لازم است؛ ثبت Log برای هر `404` در Traffic بالا می‌تواند خود به یک مشکل Performance تبدیل شود.

نشانه‌های Cache Penetration:

```
Cache Miss Rate ↑
DB Not Found Rate ↑
404 Rate ↑
Distinct Requested Keys ↑
```

Negative Cache Hit Rate پایین همراه با Distinct Key بالا

```
Traffic from Few Clients or IPs ↑
DB CPU and Pool Wait ↑
```

اگر Cache Miss زیاد باشد و Database عمدتاً `Not Found` برگرداند، احتمال Cache Penetration بالاست.

## Production Example: Product Enumeration Attack

فرض کنید:

```
Normal Traffic = 20K Requests/sec
Attack Traffic = 50K Random Product IDs/sec
DB Capacity = 10K Queries/sec
```

بدون Protection:

```
50K Random IDs → 50K Cache Misses → 50K DB Misses
```

Database سریعاً Saturate می‌شود.

طراحی دفاعی:

```
API Gateway Rate Limit → Validate ID Format → Bloom Filter → Negative Cache → DB Concurrency Limit
```

اگر Bloom Filter بتواند 99 درصد IDهای تصادفی را قطعاً ناموجود تشخیص دهد:

```
50K Requests/sec → 500 Requests/sec Reach DB
```

Database از Collapse محافظت می‌شود.

## Production Example: Deleted User Profile

فرض کنید یک لینک قدیمی برای User حذف‌شده در موتور جست‌وجو یا Clientها باقی مانده است:

```
GET /users/123
```

این URL روزانه میلیون‌ها بار Request می‌شود. چون User حذف شده، بدون Negative Cache هر Request Database را Query می‌کند.

راهکار مناسب:

```
users:123 → NOT_FOUND → Negative TTL 5 Minutes
```

اگر احتمال ایجاد دوباره همان User ID بسیار کم باشد، TTL می‌تواند نسبتاً طولانی‌تر باشد. اگر IDها قابل استفاده مجدد باشند، TTL باید کوتاه‌تر یا Invalidation روی Create وجود داشته باشد.

## Production Example: Username Availability

در سرویس ثبت‌نام، کاربران مرتب بررسی می‌کنند:

```
Is username "alex" available?
```

اگر Username موجود نباشد، این نتیجه ممکن است سریع تغییر کند، چون کاربر دیگری می‌تواند همان لحظه آن را ثبت کند. Negative Cache بلندمدت خطرناک است:

```
alex = AVAILABLE Cached for 10m
Another User Registers alex
Cache Still Says AVAILABLE
```

در این سناریو TTL باید بسیار کوتاه باشد و تصمیم نهایی هنگام ثبت‌نام حتماً با Unique Constraint در Database انجام شود. Cache فقط برای Performance است و Source of Truth برای uniqueness باید Database باشد.

## Review Scenario: Order API

فرض کنید API دریافت Order این مشخصات را دارد:

```
Traffic = 30K Requests/sec
Valid Order IDs = UUID
35% Requests Return 404
DB Capacity = 12K Queries/sec
Requests Come from Many Anonymous Clients
```

### مشکل چیست؟

حدود 10,500 Request در ثانیه به `Not Found` می‌رسند:

```
30K × 35% = 10.5K Not Found/sec
```

اگر بیشتر این Requestها Cache Miss باشند، تقریباً تمام ظرفیت Database صرف Queryهای ناموفق می‌شود.

### آیا Input Validation کافی است؟

خیر. UUIDها ممکن است از نظر Format کاملاً معتبر باشند، اما در Database وجود نداشته باشند. Validation فقط UUIDهای malformed را حذف می‌کند.

### آیا Negative Caching کافی است؟

فقط زمانی که همان UUIDهای ناموجود تکرار شوند. اگر هر Request UUID جدیدی داشته باشد، Negative Cache ممکن است Cache Pollution ایجاد کند و اثر محدودی داشته باشد.

### بهترین ترکیب راهکار چیست؟

```
UUID Format Validation
+ Rate Limiting per Client or IP
+ Bloom Filter for Existing Orders
+ Short Negative Cache
+ DB Concurrency Limit
+ Abuse Detection
```

### Bloom Filter چه کمکی می‌کند؟

Bloom Filter می‌تواند بسیاری از UUIDهای تصادفی را بدون مراجعه به Database رد کند. اگر نتیجه `Definitely Not Present` باشد، Service مستقیماً `404` برمی‌گرداند. اگر نتیجه `Maybe Present` باشد، Cache و Database بررسی می‌شوند.

### اگر Order تازه ایجاد شود چه می‌شود؟

پس از Commit موفق Order باید ID آن به Bloom Filter اضافه و Negative Cache احتمالی حذف شود:

```
Create Order → Commit DB → Update Bloom Filter → Delete Negative Cache
```

اگر Update Bloom Filter fail شود، طراحی نباید Order موجود را به‌اشتباه قطعاً ناموجود اعلام کند. در حالت عدم اطمینان باید Database بررسی شود.

### Metricهای مهم چیست؟

```
Not Found Rate
Distinct Missing UUIDs
Bloom Filter Reject Rate
Negative Cache Hit Ratio
DB Not Found Query Rate
Rate-Limited Requests
DB Pool Wait
Top Abusive Clients
```

## چه زمانی Negative Caching مناسب است؟

Negative Caching زمانی مناسب است که Missing Keyها تکرار شوند، Creation Frequency پایین باشد، Staleness کوتاه قابل قبول باشد و نتیجه `Not Found` قطعی و مستقل از User باشد.

نمونه‌ها:

```
Deleted Article
Unknown Product ID
Missing Configuration Key
Repeated Broken URL
Invalid Legacy Resource ID
```

## چه زمانی باید با احتیاط استفاده شود؟

زمانی که داده ممکن است به‌زودی ایجاد شود، Availability آن سریع تغییر کند، نتیجه به User یا Tenant وابسته باشد، Authorization روی Visibility اثر بگذارد یا `Not Found` در واقع ناشی از Temporary Failure باشد.

نمونه‌ها:

```
Username Availability
Inventory Reservation
Newly Created Order
Permission-Based Resource Lookup
Eventually Consistent Replica Reads
```

در بعضی سیستم‌ها `404` برای پنهان کردن عدم مجوز نیز استفاده می‌شود. چنین نتیجه‌ای نباید بدون در نظر گرفتن User یا Tenant به‌صورت Global Cache شود.

## Security و Cache Key Scope

Negative Cache Key باید Scope لازم را داشته باشد.

روش خطرناک:

```
resource:123 → NOT_FOUND
```

ممکن است Resource برای User A قابل مشاهده نباشد، اما برای User B وجود داشته باشد. در این حالت Global Negative Cache اشتباه است.

روش مناسب‌تر:

```
tenant:42:resource:123
user-scope:abc:resource:123
```

البته افزودن User ID به Key می‌تواند Cardinality را افزایش دهد. Design باید براساس Authorization Model انجام شود.

## Production Insight

Cache Penetration معمولاً از ترکیب این عوامل ایجاد می‌شود:

```
Cache Only Stores Existing Values
+
Not Found Results Are Never Cached
+
Input Is Weakly Validated
+
DB Fallback Is Unlimited
=
Repeated DB Misses
```

راه‌حل اصلی یک ابزار واحد نیست:

```
Validation + Negative Caching + Bloom Filter + Rate Limiting + Source Protection
```

Request نامعتبر باید تا حد امکان در ارزان‌ترین لایه متوقف شود:

```
Gateway → Validation → Bloom Filter → Cache → Database
```

هرچه Request نامعتبر دیرتر متوقف شود، هزینه بیشتری به سیستم تحمیل می‌کند.

## جمع‌بندی و Mental Model نهایی

Cache Penetration زمانی رخ می‌دهد که Requestها برای داده‌هایی ارسال شوند که وجود ندارند و چون نتیجه `Not Found` Cache نشده است، هر Request دوباره به Database می‌رسد:

```
Nonexistent Key → Cache Miss → DB Miss → Repeat
```

راهکارهای اصلی عبارت‌اند از:

```
Input Validation
Negative Caching
Short and Jittered Negative TTL
Bloom Filter
Rate Limiting
Bot and Abuse Protection
singleflight for Repeated Missing Keys
DB Concurrency Limit
Correct Error Classification
```

مدل ذهنی نهایی این است: **در Cache Penetration، مشکل این نیست که Cache داده‌ای را از دست داده است؛ مشکل این است که سیستم بارها به دنبال داده‌ای می‌گردد که از ابتدا وجود نداشته است.**

هدف طراحی صحیح:

```
Invalid or Nonexistent Request → Reject Early → Protect Cache and Database
```
