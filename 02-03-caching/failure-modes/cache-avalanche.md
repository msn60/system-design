#  Cache Failure Mode: Cache Avalanche

پس از بررسی Cache Stampede، به Failure Mode بعدی یعنی Cache Avalanche می‌رسیم. Cache Avalanche زمانی رخ می‌دهد که تعداد زیادی Cache Entry تقریباً هم‌زمان منقضی، Evict، حذف یا غیرقابل‌دسترسی شوند و در نتیجه حجم بزرگی از Requestها به‌صورت ناگهانی به Database یا Source اصلی منتقل شود.

## مفهوم اصلی

Flow خطرناک Cache Avalanche به این شکل است:

```
Many Cache Entries Lost Together → Massive Cache Misses → Broad Database Load Spike
```

برخلاف Cache Stampede که معمولاً تعداد زیادی Request روی یک Key یا تعداد محدودی Key بسیار محبوب متمرکز می‌شوند، در Cache Avalanche تعداد زیادی Key مختلف تقریباً هم‌زمان Miss می‌شوند:

```
Cache Stampede → Many Requests × One Hot Key
Cache Avalanche → Many Requests × Many Missing Keys
```

## Mental Model

یک فروشگاه بزرگ را تصور کنید که هزاران کالا در قفسه‌های نزدیک مشتری نگهداری می‌شوند و انبار اصلی در طبقه دیگری قرار دارد. اگر فقط یک کالای محبوب از قفسه تمام شود، تعداد زیادی مشتری برای همان کالا به انبار مراجعه می‌کنند؛ این وضعیت شبیه Cache Stampede است. اما اگر تمام قفسه‌ها تقریباً هم‌زمان خالی شوند، مشتریان برای هزاران کالای مختلف به انبار اصلی هجوم می‌برند؛ این وضعیت Cache Avalanche است.

```
Stampede → One Shelf Empty → Crowd Requests Same Item
Avalanche → All Shelves Empty → Entire Store Hits Warehouse
```

## Problem Statement

فرض کنید یک فروشگاه اینترنتی ده هزار Product را در Redis Cache کرده است:

```
Cached Products = 10,000
TTL = 30 Minutes
Database Capacity = 2,000 Queries/sec
Traffic = 20,000 Requests/sec
```

اگر تمام Productها هنگام Startup یا اجرای یک Batch Job تقریباً هم‌زمان Cache شده باشند، TTL آن‌ها نیز تقریباً هم‌زمان تمام می‌شود:

```
10,000 Keys Created at 10:00 → TTL = 30 Minutes → 10,000 Keys Expire around 10:30
```

در زمان Expiration، درخواست‌های مربوط به Productهای مختلف همگی Cache Miss دریافت می‌کنند و به Database منتقل می‌شوند:

```
Requests for Product A → Cache Miss → DB
Requests for Product B → Cache Miss → DB
Requests for Product C → Cache Miss → DB
Requests for Thousands of Products → DB
```

Database ناگهان به‌جای چندصد Query، با هزاران یا ده‌ها هزار Query مواجه می‌شود. نتیجه ممکن است این زنجیره باشد:

```
Mass Expiration → DB Load Spike → Pool Saturation → Query Latency ↑ → Request Timeout → Retry Traffic ↑ → Database Collapse
```

## علت‌های اصلی Cache Avalanche

### Expiration هم‌زمان تعداد زیادی Key

رایج‌ترین علت این است که Keyهای زیادی TTL یکسان داشته باشند و تقریباً هم‌زمان ساخته شده باشند:

```
100K Keys + Base TTL = 1 Hour + Created by Same Batch → Mass Expiration after 1 Hour
```

این مشکل ممکن است در Cache Warming، Batch Preloading، Deployment، Migration، Rebuild کامل Cache، Scheduled Job یا Import حجیم داده رخ دهد.

### Cache Restart یا Failover

اگر Redis Restart شود، Node از دسترس خارج شود یا Failover باعث از بین رفتن بخشی از Cache شود، تعداد زیادی Key ممکن است ناگهان غیرقابل‌دسترسی شوند:

```
Redis Restart → Cache Empty or Temporarily Unavailable → Reads Fall Back to DB
```

حتی اگر داده Redis کاملاً از بین نرفته باشد، زمان Reconnect یا Failover ممکن است باعث شود Application تعداد زیادی Cache Miss یا Cache Error مشاهده کند.

### Cache Flush

اجرای اشتباهی عملیات پاک‌سازی گسترده می‌تواند Avalanche ایجاد کند:

```
FLUSHDB / FLUSHALL / Bulk Delete → Entire Cache Lost → Full Traffic Moves to DB
```

Bulk Invalidation در Deployment یا Release نیز می‌تواند همین اثر را داشته باشد.

### Memory Pressure و Eviction گسترده

اگر Redis به سقف Memory برسد، بسته به Eviction Policy ممکن است تعداد زیادی Key حذف شوند:

```
Memory Full → Continuous Evictions → Growing Miss Rate → DB Load ↑
```

در این حالت Avalanche ممکن است تدریجی آغاز شود، اما در بازه کوتاهی حجم بزرگی از داده از Cache حذف شود.

### Cache Cluster Failure

اگر یک Shard یا Node در Redis Cluster از دسترس خارج شود، تمام Keyهای متعلق به آن Shard هم‌زمان غیرقابل‌دسترسی می‌شوند:

```
Shard 3 Down → 25% of Keys Unavailable → 25% of Traffic Hits DB
```

این نوع Avalanche محدود به بخشی از Keyspace است، اما همچنان می‌تواند برای Database بسیار سنگین باشد.

### Deployment یا تغییر Cache Key Format

اگر نسخه جدید Application فرمت Cache Key را تغییر دهد، Keyهای قدیمی دیگر قابل استفاده نیستند:

```
Old Key: product:123
New Key: v2:product:123
```

پس از Deployment، تمام Requestها برای Keyهای جدید Cache Miss می‌شوند، حتی اگر Redis پر از داده باشد. این وضعیت در عمل یک Logical Cache Flush است.

### تغییر Serialization یا Schema

اگر فرمت داده Cache تغییر کند و Application نتواند مقادیر قبلی را Decode کند، Keyها عملاً نامعتبر می‌شوند:

```
Old Cached Payload → Decode Error → Treat as Miss → Reload DB
```

اگر این تغییر بدون Compatibility، Versioning یا Migration تدریجی انجام شود، احتمال Avalanche بالا می‌رود.

## تفاوت Cache Avalanche با Cache Stampede

در Cache Stampede معمولاً یک Hot Key یا تعداد محدودی Key منقضی می‌شوند و Requestهای زیادی همان Query گران را تکرار می‌کنند:

```
One Hot Key Expires → Thousands of Requests → Same Expensive Query Repeated
```

در Cache Avalanche، هزاران Key مختلف تقریباً هم‌زمان منقضی یا غیرقابل‌دسترسی می‌شوند:

```
Thousands of Keys Expire → Requests for Many Different Objects → Thousands of Different Queries
```

در Stampede، Request Coalescing روی Key بسیار مؤثر است، چون همه Requestها نتیجه یک Load مشترک را می‌خواهند. اما در Avalanche، حتی اگر برای هر Key از `singleflight` استفاده شود، ممکن است هزاران Key متفاوت وجود داشته باشد:

```
10K Missing Keys → At Least 10K Distinct Source Loads
```

پس `singleflight` به‌تنهایی Avalanche را حل نمی‌کند. علاوه بر Per-Key Deduplication، به TTL Jitter، Cache Warming تدریجی، Global Source Limiting، Backpressure و Source Protection نیاز داریم.

## اثر روی Database

Avalanche معمولاً فشار گسترده‌تری نسبت به Stampede ایجاد می‌کند، چون Queryها ممکن است متفاوت باشند و Database نتواند از Shared Result یا الگوهای تکراری یکسان بهره زیادی ببرد. اثرهای رایج عبارت‌اند از:

```
DB Query Rate ↑
CPU Usage ↑
Disk I/O ↑
Buffer Cache Miss ↑
Connection Pool Wait ↑
Lock Contention ↑
Replication Lag ↑
P95/P99 Query Latency ↑
```

اگر Read Replica استفاده شود، Avalanche می‌تواند Replication Lag را افزایش دهد. در نتیجه Replica داده قدیمی برمی‌گرداند یا از مدار خارج می‌شود و فشار بیشتری به Primary منتقل می‌شود.

## اثر روی Connection Pool

فرض کنید:

```
Service Instances = 20
DB Pool per Instance = 100
Total Possible DB Connections = 2,000
```

در حالت عادی فقط بخش کوچکی از درخواست‌ها به DB می‌روند. اما هنگام Avalanche ممکن است تقریباً تمام Readها Database را صدا بزنند:

```
20K Requests/sec → Mostly Cache Miss → DB Pool Saturation
```

Connection Pool سریعاً پر می‌شود:

```
Active Connections = Max → Pool Waiters ↑ → Request Queue ↑ → Timeout ↑
```

افزایش اندازه Pool همیشه راه‌حل نیست. اگر Pool دو برابر شود، ممکن است فقط Database سریع‌تر Saturate شود. Pool باید نقش Bulkhead و محافظ Source را داشته باشد، نه اینکه تمام Load ورودی را بدون محدودیت به Database منتقل کند.

## اثر روی Tail Latency

قبل از Avalanche ممکن است شرایط چنین باشد:

```
Cache Read = 2ms
P99 = 8ms
```

هنگام Avalanche:

```
DB Query = 100ms–2s
Pool Wait = 500ms–3s
P99 = Several Seconds
```

حتی Requestهایی که برای داده ساده هستند نیز ممکن است به‌دلیل Queue، CPU Contention، I/O Contention یا Shared Pool کند شوند. نمودار Latency معمولاً چنین الگویی دارد:

```
Stable → Broad Latency Spike → Timeout Wave → Slow Recovery
```

Recovery نیز ممکن است کند باشد، زیرا حتی پس از بازگشت Cache، Requestهای در صف، Retryها و Queryهای قبلی همچنان در سیستم باقی مانده‌اند.

## Retry Amplification

اگر ده هزار Request Timeout شوند و هرکدام دو بار Retry داشته باشند:

```
10K Original Requests × 3 Attempts = 30K Total Attempts
```

در Avalanche، Retry روی Keyهای مختلف انجام می‌شود و Load گسترده را چند برابر می‌کند. Retry باید محدود، فقط برای خطاهای Retryable، همراه Exponential Backoff و Jitter و تحت Retry Budget باشد. Retry فوری و بدون Jitter می‌تواند موج دوم Avalanche ایجاد کند.

## راهکار اصلی: TTL Jitter

مهم‌ترین راهکار برای جلوگیری از Expiration هماهنگ، اضافه کردن Jitter به TTL است.

بدون Jitter:

```
All Keys TTL = 30m → All Keys Expire Together
```

با Jitter:

```
Base TTL = 30m
Jitter = ±5m
```

نمونه:

```
Key A = 26m42s
Key B = 31m18s
Key C = 28m55s
Key D = 34m02s
```

در نتیجه Expirationها در زمان پخش می‌شوند:

```
Mass Expiration → Distributed Expiration
```

در Go می‌توان TTL را به شکل زیر تولید کرد:

```
func jitteredTTL(base, jitter time.Duration) time.Duration {
    if base <= 0 {
        return 0
    }
    if jitter <= 0 {
        return base
    }

    maxDelta := int64(jitter)
    delta := time.Duration(rand.Int64N(2*maxDelta+1) - maxDelta)
    ttl := base + delta
    if ttl <= 0 {
        return base
    }

    return ttl
}
```

در نسخه‌های جدید Go می‌توان از تابع‌های سطح package در `math/rand/v2` استفاده کرد و نیازی به Seed کردن مکرر Generator در هر فراخوانی نیست. برای TTL Jitter نیز Random رمزنگاری‌شده لازم نیست؛ Random معمولی کافی است.

Jitter نباید آن‌قدر بزرگ باشد که Freshness Requirement را خراب کند. اگر داده حداکثر ده دقیقه می‌تواند قدیمی باشد، TTL نباید به‌صورت تصادفی تا سی دقیقه افزایش پیدا کند.

## Percentage-Based Jitter

به‌جای مقدار ثابت، Jitter می‌تواند درصدی از TTL باشد:

```
Base TTL = 1 Hour
Jitter = ±10%
```

Actual TTL = 54m تا 66m

این روش برای Keyهایی با TTLهای متفاوت مناسب‌تر است:

```
func ttlWithPercentJitter(base time.Duration, percent float64) time.Duration {
    if base <= 0 {
        return 0
    }
    if percent <= 0 {
        return base
    }

    maxJitter := time.Duration(float64(base) * percent)
    if maxJitter <= 0 {
        return base
    }

    maxDelta := int64(maxJitter)
    delta := time.Duration(rand.Int64N(2*maxDelta+1) - maxDelta)
    ttl := base + delta
    if ttl <= 0 {
        return base
    }

    return ttl
}
```

## Cache Warming تدریجی

پس از Restart یا Deployment نباید تمام Cache را با یک Batch بزرگ و هم‌زمان Warm کرد، چون همین کار TTLها را Sync می‌کند.

طراحی ضعیف:

```
Startup → Load 100K Keys at Once → All TTLs Start Together
```

طراحی بهتر:

```
Startup → Prioritize Hot Keys → Warm in Batches → Add Jitter → Rate Limit Source Reads
```

نمونه تنظیمات:

```
Batch Size = 500
Delay Between Batches = 200ms
Concurrent Loaders = 20
```

Cache Warming باید تدریجی، اولویت‌دار و متناسب با ظرفیت Database انجام شود.

## Refresh Ahead

Refresh Ahead می‌تواند پیش از Expiration، Keyها را تازه کند:

```
Near Expiry → Background Refresh → TTL Reset
```

اما اگر هزاران Key هم‌زمان وارد Refresh Window شوند، خود Refresh Ahead می‌تواند Avalanche ایجاد کند. بنابراین باید همراه با موارد زیر باشد:

```
TTL Jitter
Bounded Worker Pool
Refresh Rate Limit
Deduplication
Retry Backoff
```

Refresh Ahead بدون Concurrency Control فقط Avalanche را از زمان Expiration به زمان Refresh منتقل می‌کند.

## Serve Stale

در صورت Expiration، Cache Error یا خرابی Source می‌توان برای مدت محدودی مقدار قبلی را ارائه کرد:

```
Soft TTL Passed → Serve Stale → Refresh in Background
Hard TTL Passed → Reload, Fallback or Fail
```

مثال:

```
Soft TTL = 10m
Hard TTL = 1h
```

اگر Redis یا Database موقتاً مشکل داشته باشد، Application می‌تواند تا سقف Hard TTL مقدار قبلی را ارائه کند؛ البته فقط برای داده‌هایی که Staleness آن‌ها قابل قبول است.

## Multi-Level Cache

استفاده از Local In-Memory Cache در کنار Redis می‌تواند اثر خرابی Redis را کاهش دهد:

```
Request → L1 Local Cache → L2 Redis → Database
```

اگر Redis Down شود، بخشی از Hot Data همچنان در L1 Cache هر Instance باقی می‌ماند.

مزایا:

```
Redis Load ↓
Network Dependency ↓
Partial Protection During Redis Failure
```

معایب:

```
Memory Duplication
Per-Instance Staleness
Invalidation Complexity
Lower Cross-Instance Consistency
```

L1 Cache باید TTL و Memory Limit داشته باشد و باید در نظر داشت که هر Instance نسخه مستقل خود را نگه می‌دارد.

## Rate Limiting و Concurrency Limiting

وقتی Cache Miss گسترده رخ می‌دهد، نباید تمام Requestها اجازه داشته باشند هم‌زمان Database را صدا بزنند:

```
Cache Misses → Source Concurrency Limiter → Database
```

در Go می‌توان از Semaphore استفاده کرد. یک مدل ساده مبتنی بر Channel:

```
type LoaderLimiter struct {
    sem chan struct{}
}

func NewLoaderLimiter(limit int) *LoaderLimiter {
    if limit <= 0 {
        panic("loader concurrency limit must be positive")
    }

    return &LoaderLimiter{
        sem: make(chan struct{}, limit),
    }
}

func (l *LoaderLimiter) Acquire(ctx context.Context) error {
    select {
    case l.sem <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (l *LoaderLimiter) Release() {
    <-l.sem
}
```

برای پیاده‌سازی‌های عمومی‌تر در Go می‌توان از `golang.org/x/sync/semaphore` نیز استفاده کرد، به‌خصوص زمانی که Weighted Semaphore نیاز است.

این محدودیت Database را از Collapse محافظت می‌کند، اما Queue و Wait Time باید محدود باشند. اگر هزاران Request پشت Semaphore منتظر بمانند، Memory Consumption و Tail Latency افزایش پیدا می‌کنند.

## Load Shedding

در شرایط بحرانی، بهتر است بخشی از Requestها سریع و کنترل‌شده رد شوند تا کل سیستم Collapse نکند:

```
Source Capacity Reached → Reject Low-Priority Requests → Protect Core Traffic
```

رفتارهای ممکن:

```
Return 503
Serve Stale
Use Default Response
Skip Optional Data
Degrade Feature
```

برای مثال صفحه اصلی فروشگاه می‌تواند بدون Recommendation یا Trending پاسخ داده شود، به‌جای اینکه کل صفحه Timeout شود.

## Graceful Degradation

Cache Avalanche نباید الزاماً باعث Failure کامل Endpoint شود. بخش‌های غیرضروری می‌توانند حذف یا با داده ساده جایگزین شوند:

```
Home Page:
Product List = Required
Recommendations = Optional
Personalized Banner = Optional
Trending = Optional
```

در زمان Avalanche:

```
Return Product List → Skip Recommendations → Use Static Banner → Serve Previous Trending Data
```

این رفتار Availability را بالا می‌برد و فشار روی Source را کاهش می‌دهد.

## Negative Caching

اگر بخشی از Missها مربوط به داده‌های ناموجود باشد، Negative Caching جلوی Queryهای تکراری را می‌گیرد:

```
product:999999 → NOT_FOUND → TTL 30s
```

در Avalanche ناشی از Traffic نامعتبر یا Botها، Negative Caching اهمیت بیشتری پیدا می‌کند.

## Bloom Filter

برای Datasetهای بزرگ و Requestهای نامعتبر، Bloom Filter می‌تواند پیش از Database بررسی کند که آیا Key احتمالاً وجود دارد یا نه:

```
Request Key → Bloom Filter
Definitely Not Present → Reject without DB
Maybe Present → Continue Cache/DB Lookup
```

Bloom Filter ممکن است False Positive داشته باشد، اما در طراحی صحیح False Negative ندارد؛ یعنی ممکن است بگوید Key احتمالاً وجود دارد در حالی که وجود ندارد، اما نباید Key موجود را قطعاً ناموجود اعلام کند. این راهکار بیشتر در Cache Penetration اهمیت دارد، اما در Avalanche ناشی از حمله یا Traffic نامعتبر نیز مفید است.

## Source Protection

حتی با بهترین Cache Design، Database باید مستقل از Cache محافظت شود:

```
Connection Pool Limit
Query Timeout
Statement Timeout
Concurrency Limit
Rate Limit
Circuit Breaker
Load Shedding
Read Replica
Admission Control
```

اگر Cache کاملاً Down شود، Source نباید کل Traffic را بدون محدودیت بپذیرد.

اصل مهم:

```
Cache Failure Must Not Automatically Become Database Failure
```

## Circuit Breaker

اگر Database یا Redis به‌وضوح Saturate شده باشد، ادامه ارسال Request فقط وضعیت را بدتر می‌کند:

```
Failure Rate High → Open Circuit → Fast Fail or Serve Fallback
```

Circuit Breaker باید با دقت تنظیم شود. اگر همه Instanceها هم‌زمان وارد Half-Open شوند و Probe بفرستند، ممکن است Recovery Storm ایجاد شود. Probeها باید محدود و همراه Jitter باشند.

## Bulkhead

می‌توان Pool یا Concurrency منابع مختلف را جدا کرد:

```
Critical Product Reads → Pool A
Recommendation Reads → Pool B
Analytics Reads → Pool C
```

اگر Recommendation دچار Avalanche شود، نباید Connectionهای مسیرهای حیاتی مانند Product اصلی یا Payment را مصرف کند. این جداسازی Bulkhead نام دارد و Failure Domainها را محدود می‌کند.

## Cache Cluster High Availability

برای کاهش Avalanche ناشی از Redis Failure باید Cache Layer نیز High Availability داشته باشد:

```
Replication
Sentinel or Cluster Failover
Multiple Availability Zones
Health Checks
Connection Retry with Backoff
Capacity Headroom
```

اما High Availability به معنی Zero Failure نیست. Failover ممکن است چند ثانیه طول بکشد و در این فاصله Application باید رفتار مشخصی داشته باشد.

## رفتار Application هنگام Redis Failure

یک Anti-pattern رایج:

```
Redis Error → Treat Everything as Cache Miss → Send All Traffic to DB
```

این رفتار Redis Failure را فوراً به Database Failure تبدیل می‌کند.

رفتار بهتر:

```
Redis Failure Detected → Limit DB Fallback → Serve Local or Stale Data → Degrade Optional Features → Apply Load Shedding
```

برای هر Endpoint باید از قبل مشخص شود که Cache Failure چه اثری دارد و چه Fallbackی مجاز است.

## Cache Key Versioning

در زمان تغییر فرمت Key نباید تمام Traffic ناگهان به Namespace جدید منتقل شود.

روش خطرناک:

```
Deploy v2 → Read Only v2 Keys → All v1 Cache Becomes Useless
```

روش امن‌تر:

```
Read v2 → If Miss Read v1 → If v1 Hit Convert and Write v2 → Else Load DB
```

یا در دوره Migration می‌توان از Dual Write استفاده کرد:

```
Write v1 + v2 During Transition
```

این روش Migration Load را در زمان پخش می‌کند.

## Serialization Compatibility

هنگام تغییر Cache Payload بهتر است Version داخل Value ذخیره شود:

```
{
  "version": 2,
  "data": {
    "id": "123",
    "name": "Product"
  }
}
```

در صورت مواجهه با فرمت قدیمی می‌توان آن را Convert کرد، از Reader قبلی استفاده کرد، به‌تدریج Refresh کرد یا فقط همان Key را Invalidate کرد. یک Decode Error نباید باعث پاک‌سازی گسترده Cache شود.

## راهکار Recovery پس از Redis Restart

اگر Cache پس از Restart خالی باشد، باید Recovery Plan مشخصی وجود داشته باشد:

```
Redis Starts → Do Not Send Full Traffic Immediately → Warm Critical Keys → Gradually Increase Traffic → Keep DB Fallback Limited
```

راهکارهای ممکن:

```
Startup Warming
Traffic Ramp-Up
Readiness Gate
Rate-Limited Rebuild
Restore from Persistence
L1 Cache
Serve Stale Snapshot
```

Readiness صرفاً نباید بررسی کند Redis Process بالا آمده است؛ باید مشخص شود Cache و Dependencyها ظرفیت پذیرش Traffic را دارند.

## Go Implementation Model

یک مدل Production-oriented در Go می‌تواند چنین باشد:

```
Request → L1 Cache → L2 Cache → On Miss Check Source Limiter → Enter Per-Key singleflight → Recheck Cache → Load Source with Timeout → Write Cache with Jittered TTL → Return
```

اجزای مهم:

```
context.Context
singleflight
Bounded Source Concurrency
Jittered TTL
Timeout
Stale Fallback
Structured Metrics
Graceful Shutdown
```

در Avalanche، `singleflight` همچنان مفید است، زیرا برای هر Key از Queryهای تکراری جلوگیری می‌کند؛ اما چون Keyها متعدد هستند، باید همراه Global Concurrency Limit استفاده شود:

```
Per-Key Deduplication + Global Source Limit
```

`singleflight` Queryهای تکراری همان Key را حذف می‌کند و Semaphore یا Worker Pool تعداد کل Queryهای متفاوت را محدود می‌سازد.

## Observability

برای تشخیص Cache Avalanche باید هم Cache و هم Database مشاهده شوند. Metricهای مهم عبارت‌اند از:

```
cache_hit_rate
cache_miss_rate
cache_key_expiration_rate
cache_eviction_rate
cache_error_rate
redis_connected_clients
redis_memory_usage
redis_evicted_keys
source_query_rate
source_query_inflight
db_pool_waiters
db_pool_wait_duration
request_timeout_rate
retry_rate
fallback_rate
stale_served_rate
load_shed_total
```

نشانه‌های Avalanche:

Cache Hit Rate ↓ گسترده

Cache Miss Rate ↑ برای Keyهای متعدد

DB Query Rate ↑ ناگهانی

```
Eviction Rate ↑
Redis Errors ↑
Pool Waiters ↑
P99 Latency ↑
Timeout and Retry ↑
```

تفاوت Observability در Stampede و Avalanche:

```
Stampede → One or Few Keys Dominate Misses
Avalanche → Misses Distributed Across Many Keys
```

بنابراین باید Distribution و تعداد Distinct Keyهای Miss شده نیز بررسی شود، نه فقط Miss Count کل.

## Alertهای مهم

Alertهای مناسب شامل این موارد هستند:

```
Cache Hit Ratio Drops Below Threshold
Eviction Rate Spikes
Redis Unavailable
DB Query Rate Exceeds Safe Capacity
DB Pool Wait Exceeds SLO
Fallback Rate Grows Rapidly
High Number of Distinct Missed Keys
Retry Rate and Timeout Rate Rise Together
```

Alert صرفاً روی Redis Down کافی نیست. ممکن است Redis سالم باشد، اما Expiration هماهنگ، Eviction گسترده یا تغییر Key Format باعث Avalanche شود.

## Production Example: Product Catalog

فرض کنید:

```
Products = 500K
Traffic = 30K Requests/sec
Cache TTL = 1 Hour
DB Capacity = 4K Queries/sec
Deployment Rebuilds Cache Keys
```

نسخه جدید Prefix Key را از `product:` به `v2:product:` تغییر می‌دهد. پس از Deployment:

```
All v2 Keys Missing → Large DB Fallback → DB Saturation
```

حتی با وجود Redis سالم، یک Logical Avalanche رخ داده است.

طراحی بهتر:

```
Read v2 → If Miss Read v1 → Promote to v2 → Use DB Only If Both Miss
```

همچنین:

```
TTL = 1h ± 10%
Warming Limited to Hot Products
Source Concurrency Capped
Optional Product Details Degraded Under Pressure
```

## Review Scenario: News Platform

فرض کنید یک News Platform داریم:

```
Cached Articles = 200K
Traffic = 40K Requests/sec
Cache TTL = 20 Minutes
Redis Cluster = 4 Shards
One Shard Goes Down
Each Shard Owns about 25% of Keys
DB Safe Capacity = 3K Queries/sec
```

### مشکل چیست؟

با از دست رفتن یک Shard، تقریباً 25 درصد Keyها دیگر در دسترس نیستند. اگر 25 درصد Traffic به Database منتقل شود:

```
40K Requests/sec × 25% = 10K DB Requests/sec
```

Database فقط 3K Query/sec ظرفیت امن دارد؛ بنابراین حدود سه برابر ظرفیت خود Load دریافت می‌کند و احتمال Saturation بسیار بالاست.

### آیا singleflight کافی است؟

خیر. `singleflight` از Queryهای تکراری برای یک Article جلوگیری می‌کند، اما هزاران Article مختلف روی Shard ازدست‌رفته وجود دارند. حداقل یک Query برای هر Article همچنان لازم خواهد بود.

### ترکیب راهکار مناسب چیست؟

```
Redis HA/Failover
L1 Local Cache
Serve Stale
Global DB Concurrency Limit
Per-Key singleflight
Load Shedding
Priority-Based Degradation
Gradual Cache Rebuild
```

### هنگام Failure چه داده‌ای باید اولویت داشته باشد؟

مقالات Hot و صفحات اصلی باید اولویت بالاتری داشته باشند. داده‌های قدیمی، Cold یا Optional نباید ظرفیت محدود Database را مصرف کنند:

```
Hot Content → Allow Controlled DB Load
Cold Content → Serve Stale or Controlled Error
Recommendations → Disable Temporarily
```

### اگر Redis Shard برگردد چه خطری وجود دارد؟

ممکن است همه Instanceها هم‌زمان تلاش کنند Cache را دوباره Warm کنند و Recovery Avalanche ایجاد شود. Recovery باید همراه Jitter، Rate Limit، Batching و Gradual Ramp-Up باشد.

### Metricهای حیاتی چیست؟

```
Per-Shard Cache Error Rate
Distinct Cache Miss Keys
DB Query Rate
DB Pool Wait
Fallback Rate
Stale Served Rate
Load Shed Count
Cache Rebuild Rate
P99 Latency
```

## تفاوت Cache Avalanche و Hot Key

Cache Avalanche به از دست رفتن یا Expiration گسترده Keyهای متعدد مربوط است:

```
Many Keys Lost → Broad Source Load
```

Hot Key به تمرکز شدید Traffic روی یک Key مربوط است:

```
One Key → Huge Cache Traffic
```

یک Redis Node ممکن است به‌دلیل Hot Key Saturate شود و Failure همان Node سپس برای Keyهای دیگر Avalanche ایجاد کند. بنابراین این Failure Modeها می‌توانند به یکدیگر متصل شوند.

## تفاوت Cache Avalanche و Cache Penetration

در Cache Avalanche، داده معمولاً وجود دارد ولی Cache آن را از دست داده است:

```
Existing Data + Missing Cache → DB Load
```

در Cache Penetration، Requestها برای داده‌هایی هستند که معمولاً اصلاً وجود ندارند:

```
Nonexistent Data → Cache Miss + DB Miss Repeatedly
```

راهکارهای اصلی نیز متفاوت‌اند:

```
Avalanche → TTL Jitter, HA, Warming, Source Limiting
Penetration → Negative Cache, Bloom Filter, Validation
```

## Production Insight

Cache Avalanche معمولاً نتیجه یک Failure منفرد نیست؛ بلکه حاصل ترکیب چند تصمیم ظاهراً ساده است:

```
Same TTL + Mass Cache Population + No Jitter + Unlimited DB Fallback + Aggressive Retry = Avalanche
```

همچنین Redis Down نباید به‌طور خودکار کل سیستم را Down کند. اگر معماری طوری طراحی شده باشد که Redis Failure تمام Traffic را بدون محدودیت به Database منتقل کند، Cache به‌جای Performance Layer به Single Point of Cascading Failure تبدیل شده است.

اصل مهم:

```
A Cache Outage Must Degrade the System, Not Automatically Collapse It
```

## جمع‌بندی و Mental Model نهایی

Cache Avalanche زمانی رخ می‌دهد که بخش بزرگی از Cache در یک بازه کوتاه از بین برود، منقضی یا نامعتبر شود و Traffic گسترده به Source اصلی منتقل شود:

```
Many Keys Lost Together → Many Broad Cache Misses → Source Overload
```

تفاوت اصلی با Stampede:

```
Stampede → Many Requests Compete for the Same Missing Value
Avalanche → Many Different Values Disappear at the Same Time
```

راهکار اصلی یک ابزار واحد نیست، بلکه ترکیبی از چند لایه دفاعی است:

```
TTL Jitter
Gradual Cache Warming
Refresh Rate Limiting
Serve Stale
Per-Key singleflight
Global Source Concurrency Limit
Load Shedding
Graceful Degradation
Cache HA
Source-Side Protection
```

مدل ذهنی نهایی این است: **در Cache Avalanche، مشکل این نیست که یک Key Miss شده است؛ مشکل این است که بخش بزرگی از Cache هم‌زمان ناپدید می‌شود و تمام Source Load ذخیره‌شده پشت Cache، ناگهان آزاد می‌شود.**
