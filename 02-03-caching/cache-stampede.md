# Cache Failure Mode: Cache Stampede

پس از بررسی Caching Fundamentals و Cache Strategyهای مختلف شامل Cache Aside، Read Through، Write Through، Write Back و Refresh Ahead، وارد بخش Cache Failure Modes می‌شویم. در این بخش دیگر تمرکز اصلی روی نحوه استفاده از Cache نیست، بلکه روی شرایطی است که Cache ظاهراً سالم است، اما رفتار هم‌زمان Requestها باعث فشار شدید روی Database یا Source اصلی می‌شود. اولین Failure Mode مهم، Cache Stampede است.

## مفهوم اصلی

Cache Stampede زمانی رخ می‌دهد که تعداد زیادی Request تقریباً هم‌زمان برای یک داده واحد با Cache Miss مواجه شوند و همه آن‌ها برای بارگذاری همان داده به Database یا Source اصلی مراجعه کنند.

```
Hot Key Expires → Thousands of Requests Miss → Thousands of Source Queries
```

در حالت مطلوب، فقط یک Request باید داده را از Source بارگذاری کند، Cache را دوباره بسازد و نتیجه را در اختیار بقیه Requestها قرار دهد. اما اگر هیچ Coordination یا Deduplication وجود نداشته باشد، تمام Requestها همان عملیات گران را مستقل از یکدیگر انجام می‌دهند.

Cache Stampede گاهی با اصطلاح‌های Thundering Herd، Dogpile Effect یا Cache Miss Storm نیز توصیف می‌شود. این اصطلاح‌ها همیشه دقیقاً مترادف نیستند، اما در زمینه Cache معمولاً به هجوم هم‌زمان Requestها به Source پس از Miss، Expiration یا Eviction اشاره دارند.

## Mental Model

فرض کنید هزار نفر پشت یک در بسته ایستاده‌اند. وقتی در باز می‌شود، به‌جای اینکه یک نفر وارد شود و نتیجه را برای بقیه بیاورد، هر هزار نفر هم‌زمان وارد می‌شوند.

```
One Missing Value + Many Concurrent Requests = Many Identical Source Loads
```

مدل نامطلوب:

```
Many Requests → Many DB Loads for Same Key → Database Saturation
```

مدل مطلوب:

```
Many Requests → One Controlled DB Load → Cache Rebuilt → Shared Result
```

پس مشکل اصلی Cache Stampede فقط یک Cache Miss نیست؛ مشکل این است که هزاران Request تلاش می‌کنند همان Miss را به‌صورت مستقل جبران کنند.

## Problem Statement

فرض کنید API صفحه اصلی یک فروشگاه این مشخصات را دارد:

```
Traffic = 20,000 Requests/sec
Cache TTL = 5 Minutes
Database Query = 800ms
Cache Read = 3ms
```

تا زمانی که Cache معتبر است، درخواست‌ها با Latency بسیار کم پاسخ داده می‌شوند:

```
Request → Cache Hit → 3ms Response
```

پس از پنج دقیقه، Key منقضی می‌شود. اگر در همان لحظه هزاران Request وارد شوند:

```
20K Requests → Cache Miss → 20K Database Queries
```

Database که شاید فقط ظرفیت چندصد Query سنگین در ثانیه را داشته باشد، ناگهان با هزاران Query یکسان مواجه می‌شود. نتیجه می‌تواند شامل افزایش CPU و I/O، پر شدن Connection Pool، افزایش Query Latency، Timeout، Retry و در نهایت ناپایداری گسترده سیستم باشد.

```
Cache Expiration → DB Load Spike → Pool Exhaustion → Timeout → Retry → More Load
```

این Feedback Loop می‌تواند حتی بعد از بازسازی Cache نیز سیستم را برای مدتی ناپایدار نگه دارد.

## شرایط ایجاد Cache Stampede

Stampede معمولاً از ترکیب چند عامل ایجاد می‌شود:

```
Hot Key + Expiration or Eviction + Concurrent Requests + Slow Source + No Deduplication
```

هرچه Key محبوب‌تر، Source Load کندتر و تعداد Requestهای هم‌زمان بیشتر باشد، شدت Stampede بالاتر خواهد بود.

### Expiration یک Hot Key

رایج‌ترین حالت، منقضی شدن یک Key بسیار پرترافیک است:

```
homepage:data Expires → All Home Page Requests Miss
```

### Cache Eviction

ممکن است Key به‌دلیل Memory Pressure از Cache حذف شود، حتی اگر TTL آن تمام نشده باشد:

```
Redis Memory Pressure → Hot Key Evicted → Sudden Miss Storm
```

### Cache Restart یا Flush

بعد از Restart شدن Redis، پاک شدن Cache، Failover یا Deployment ممکن است داده‌های Cache از بین بروند. اگر تمرکز درخواست‌ها روی یک Key باشد، Stampede ایجاد می‌شود. اگر تعداد زیادی Key هم‌زمان از بین بروند، مشکل به Cache Avalanche نزدیک می‌شود.

### Refresh Failure

ممکن است Refresh Ahead یا Background Refresh نتواند داده را به‌موقع تازه کند و Key در نهایت Expire شود. سپس Requestهای هم‌زمان همگی Source را صدا می‌زنند.

### Cold Start

پس از راه‌اندازی یک Service یا Cache جدید، داده هنوز Cache نشده است. اگر Traffic بلافاصله وارد شود، چندین Request یکسان ممکن است هم‌زمان Source را فراخوانی کنند.

## نمونه ساده Cache Aside در Go

یک پیاده‌سازی معمول Cache Aside ممکن است چنین باشد:

```
func GetProduct(ctx context.Context, id string) (Product, error) {
    key := "product:" + id

    product, err := cache.GetProduct(ctx, key)
    if err == nil {
        return product, nil
    }

    product, err = db.GetProduct(ctx, id)
    if err != nil {
        return Product{}, err
    }

    if err := cache.SetProduct(ctx, key, product, 5*time.Minute); err != nil {
        // Observe the cache write failure.
    }

    return product, nil
}
```

این کد در حالت عادی درست به نظر می‌رسد، اما هیچ محافظتی در برابر Concurrency ندارد. اگر هزار Request هم‌زمان Cache Miss دریافت کنند، همه آن‌ها `db.GetProduct` را اجرا خواهند کرد:

```
1000 Cache Misses → 1000 Identical DB Queries
```

مشکل از خود Cache Aside نیست؛ مشکل این است که Missها با یکدیگر Coordination ندارند.

## اثر روی Connection Pool

یکی از اولین بخش‌هایی که هنگام Stampede آسیب می‌بیند، Connection Pool است. فرض کنید:

```
DB Pool Size = 100
Concurrent Misses = 5,000
```

فقط 100 Request می‌توانند Connection دریافت کنند و 4,900 Request دیگر منتظر می‌مانند:

```
5000 Requests → 100 Active DB Connections → 4900 Pool Waiters
```

اگر Query اصلی 800 میلی‌ثانیه طول بکشد، Pool Wait بالا می‌رود. در ادامه Requestها Timeout می‌شوند و Client یا Service بالادستی ممکن است Retry کند:

```
Cache Miss → Pool Wait → Timeout → Retry → Additional Load
```

در این وضعیت، کوچک بودن Pool لزوماً مشکل اصلی نیست. افزایش بی‌رویه Pool نیز ممکن است فقط فشار بیشتری را به Database منتقل کند. Pool باید یک محدودکننده محافظتی باشد، نه راه‌حل Stampede.

## اثر روی Tail Latency

قبل از Stampede:

```
Cache Hit Latency = 3ms
```

هنگام Stampede ممکن است شرایط چنین شود:

```
DB Query = 800ms
Pool Wait = 2s
Request Timeout = 3s
```

در نتیجه P95، P99 و P99.9 ناگهان افزایش پیدا می‌کنند. Average Latency ممکن است همچنان برای مدت کوتاهی قابل‌قبول به نظر برسد، اما Tail Latency سریع‌تر مشکل را آشکار می‌کند:

```
Stable Low Latency → Sudden Sharp Spike → Gradual Recovery
```

## چرا Retry وضعیت را بدتر می‌کند؟

اگر Requestها Timeout شوند و Retry Policy تهاجمی باشد، هر Request اولیه می‌تواند چند Request جدید تولید کند:

```
10K Original Requests × 3 Attempts → Up to 30K Total Requests
```

در این حالت Retry به‌جای Recovery، Load Amplification ایجاد می‌کند. Retry باید محدود، فقط برای خطاهای Retryable و همراه با Exponential Backoff و Jitter باشد. همچنین Retry Budget باید مشخص کند چه میزان Retry برای سیستم قابل قبول است.

## راهکار اصلی: Request Coalescing

مهم‌ترین راهکار مقابله با Cache Stampede این است که Requestهای هم‌زمان مربوط به یک Key را به یک Source Load تبدیل کنیم.

```
1000 Requests for product:123 → One DB Query → One Cache Set → Shared Result
```

این روش Request Coalescing، Miss Deduplication یا Single-Flight Loading نامیده می‌شود.

در Go می‌توان از `golang.org/x/sync/singleflight` استفاده کرد:

```
var group singleflight.Group

func GetProduct(ctx context.Context, id string) (Product, error) {
    key := "product:" + id

    product, err := cache.GetProduct(ctx, key)
    if err == nil {
        return product, nil
    }

    value, err, _ := group.Do(key, func() (any, error) {
        // Recheck the cache because another request may have populated it
        // before this loader began.
        product, err := cache.GetProduct(ctx, key)
        if err == nil {
            return product, nil
        }

        product, err = db.GetProduct(ctx, id)
        if err != nil {
            return Product{}, err
        }

        if err := cache.SetProduct(ctx, key, product, 5*time.Minute); err != nil {
            // Record the cache write failure without hiding the DB result.
        }

        return product, nil
    })
    if err != nil {
        return Product{}, err
    }

    product, ok := value.(Product)
    if !ok {
        return Product{}, fmt.Errorf("unexpected product result type %T", value)
    }

    return product, nil
}
```

Flow این مدل:

```
Many Goroutines → singleflight by Key → One Loader Execution → Shared Result
```

Double Check داخل `singleflight` مهم است، زیرا ممکن است یک Request دیگر بین Cache Miss اولیه و شروع Loader، Cache را ساخته باشد.

`singleflight` فقط داخل همان Process عمل می‌کند. اگر Service بیست Instance داشته باشد، هر Instance ممکن است یک Query اجرا کند:

```
20 Service Instances → Up to 20 Identical DB Queries
```

این وضعیت بسیار بهتر از هزاران Query است، اما Deduplication کامل در سطح Cluster نیست.

## Distributed Lock

برای Deduplication میان چند Instance می‌توان از Distributed Lock استفاده کرد:

```
Cache Miss → Acquire Distributed Lock → Recheck Cache → Load DB → Update Cache → Release Lock
```

نمونه مفهومی با Redis:

```
SET lock:product:123 unique-token NX PX 5000
```

اما Distributed Lock پیچیدگی‌های مهمی دارد:

```
Lock Expiration
Owner Crash
Slow Source Load
Lock Renewal
Token Ownership
Safe Unlock
Network Partition
Clock and Timeout Assumptions
```

هنگام آزاد کردن Lock باید فقط Owner همان Lock مجاز به حذف آن باشد. حذف ساده Lock بدون بررسی Token می‌تواند Lock متعلق به Request دیگری را حذف کند.

Distributed Lock نباید اولین انتخاب برای هر Cache Miss باشد. ترتیب عملی مناسب‌تر معمولاً چنین است:

```
In-Process singleflight → Serve Stale → Refresh Ahead → TTL Jitter → Distributed Lock if Required
```

اگر تعداد محدود Query در سطح Instanceها برای Database قابل تحمل باشد، Distributed Lock ممکن است فقط Complexity و Failure Modeهای جدید ایجاد کند.

## Double-Check Locking

بعد از گرفتن Lock باید Cache دوباره بررسی شود:

```
Cache Miss → Acquire Lock → Check Cache Again → Load DB Only If Still Missing
```

بدون Double Check ممکن است Requestها با فاصله زمانی Lock را دریافت کنند و هرکدام دوباره Source را بخوانند، حتی اگر Cache قبلاً توسط Request قبلی ساخته شده باشد.

## Serve Stale

یکی از مؤثرترین روش‌ها برای جلوگیری از Stampede، ارائه موقت مقدار قدیمی در زمان Refresh است:

```
Soft TTL Passed → Serve Stale Value → One Background Refresh
```

این رفتار به Stale-While-Revalidate نزدیک است:

```
Fresh Window → Serve Fresh
Stale Window → Serve Stale + Refresh
Hard Expiry → Reload, Fallback or Fail
```

مثال:

```
Soft TTL = 5 Minutes
Hard TTL = 30 Minutes
```

اگر Refresh در دقیقه پنجم آغاز شود ولی Database موقتاً Down باشد، مقدار قبلی تا سقف Hard TTL قابل ارائه است. این روش فقط برای داده‌هایی مناسب است که مقدار محدودی Staleness در آن‌ها قابل قبول باشد.

## Refresh Ahead

Refresh Ahead احتمال وقوع Stampede را کاهش می‌دهد، زیرا Cache را پیش از Expiration تازه می‌کند:

```
Near Expiry → Deduplicated Background Refresh → TTL Reset
```

بااین‌حال، Refresh Ahead به‌تنهایی کافی نیست. اگر Refresh fail شود یا چند Request هم‌زمان آن را Trigger کنند، همچنان به Deduplication، Worker Limit، Cooldown و Backoff نیاز داریم.

## TTL Jitter

اگر Keyهای زیادی TTL دقیقاً یکسان داشته باشند، ممکن است هم‌زمان منقضی شوند. برای پخش زمان Expiration و Refresh، باید TTL به‌اندازه محدودی تصادفی شود:

```
Base TTL = 10 Minutes
Jitter = ±60 Seconds
```

مثال:

```
Key A TTL = 9m22s
Key B TTL = 10m41s
Key C TTL = 9m57s
```

TTL Jitter بیشتر برای جلوگیری از Expiration هماهنگ چند Key کاربرد دارد و در Cache Avalanche اهمیت ویژه‌ای خواهد داشت، اما برای Hot Keyها نیز به پخش Refresh Load کمک می‌کند.

## Loader Timeout و Waiter Timeout

اگر Request مالک Load شود ولی Source Query گیر کند، بقیه Requestها نباید بدون محدودیت منتظر بمانند. بهتر است دو Timeout مستقل تعریف شود:

```
Loader Timeout
Waiter Timeout
```

مثال:

```
Source Load Timeout = 2s
Waiter Maximum Wait = 500ms
```

پس از Waiter Timeout می‌توان یکی از این رفتارها را انتخاب کرد:

```
Serve Stale
Return Controlled Error
Use Fallback
Retry Later with Backoff
```

این تصمیم به Business Requirement، SLO و میزان Staleness قابل قبول بستگی دارد.

## Negative Caching

Stampede فقط برای داده‌های موجود رخ نمی‌دهد. اگر تعداد زیادی Request برای یک ID نامعتبر ارسال شود و Cache فقط مقادیر موجود را ذخیره کند، تمام Requestها ممکن است Database را صدا بزنند:

```
GET product:999999 → Not Found
Repeated Invalid Requests → Repeated DB Misses
```

برای کاهش این فشار می‌توان نتیجه منفی را با TTL کوتاه Cache کرد:

```
product:999999 → NOT_FOUND → TTL 30s
```

این روش Negative Caching نام دارد و در Cache Penetration نقش اصلی‌تری دارد. TTL نتیجه منفی باید کوتاه‌تر باشد تا اگر داده کمی بعد ایجاد شد، نتیجه Not Found برای مدت طولانی باقی نماند.

## Probabilistic Early Refresh

به‌جای آغاز Refresh دقیقاً در یک نقطه ثابت، می‌توان با نزدیک شدن به Expiration احتمال Refresh را افزایش داد:

```
Far from Expiry → Low Refresh Probability
Near Expiry → High Refresh Probability
```

این روش که به Probabilistic Early Expiration یا Early Recomputation نزدیک است، Refreshها را در زمان پخش می‌کند و احتمال هجوم هم‌زمان را کاهش می‌دهد.

## Queue و Bounded Worker Pool

برای Refreshهای Background بهتر است از Queue کنترل‌شده و Worker Pool محدود استفاده شود:

```
Miss or Near Expiry → Enqueue Refresh Job → Deduplicate → Bounded Worker Pool → Source
```

Queue باید Bounded باشد؛ Queue نامحدود فقط فشار را از Database به Memory منتقل می‌کند. اگر Queue پر شود، رفتار باید از قبل مشخص باشد:

```
Drop Duplicate Job
Skip Refresh
Serve Stale
Apply Backpressure
Return Controlled Error
```

برای Cache Refresh معمولاً حذف Job تکراری یا Serve Stale منطقی‌تر از Block کردن کل Request Path است.

## Cache Stampede در سیستم توزیع‌شده

در یک Service چند-Instance، حل Stampede دشوارتر است:

```
Instance A → Cache Miss
Instance B → Cache Miss
Instance C → Cache Miss
```

هر Instance ممکن است داخل خودش `singleflight` داشته باشد، اما از Loadهای Instanceهای دیگر خبر ندارد.

راهکارهای ممکن:

```
Distributed Lock
Shared Refresh Queue
Leader-Based Refresh
Event-Based Refresh
Serve Stale
Source-Side Concurrency Limit
```

هدف همیشه نباید کاهش تعداد Query به دقیقاً یک Query در کل Cluster باشد. اگر بیست Instance داشته باشیم و هرکدام فقط یک Query اجرا کنند، ممکن است بیست Query برای Database کاملاً قابل تحمل باشد. راهکار باید متناسب با ظرفیت واقعی Source و هزینه Complexity انتخاب شود.

## Source-Side Protection

Cache نباید تنها سد دفاعی سیستم باشد. Database یا Dependency اصلی نیز باید بتواند از خود محافظت کند:

```
Connection Pool Limit
Query Timeout
Concurrency Limit
Rate Limit
Circuit Breaker
Load Shedding
Admission Control
```

اگر Cache از دسترس خارج شود، Source نباید تمام Traffic را بدون محدودیت بپذیرد. محدودیت Concurrency و Pool باعث می‌شود Source به‌جای Collapse کامل، بخشی از Load را به‌صورت کنترل‌شده رد کند.

## تفاوت Cache Stampede و Cache Avalanche

این دو Failure Mode نزدیک‌اند، اما یکسان نیستند.

Cache Stampede معمولاً درباره تعداد زیادی Request برای یک Key یا تعداد محدودی Key بسیار محبوب است:

```
One Hot Key Expires → Many Requests Hit Source
```

Cache Avalanche معمولاً زمانی رخ می‌دهد که تعداد زیادی Key تقریباً هم‌زمان منقضی، Evict یا حذف شوند:

```
Many Keys Expire Together → Broad Source Load Spike
```

Stampede بیشتر روی Concurrency برای یک داده تمرکز دارد؛ Avalanche روی Expiration یا Loss گسترده مجموعه‌ای از داده‌ها متمرکز است.

## تفاوت Cache Stampede و Hot Key

Hot Key یعنی یک Key سهم بسیار زیادی از Traffic Cache را دریافت می‌کند:

```
One Key → Huge Cache Traffic
```

ممکن است Hot Key همچنان Cache Hit باشد و فقط Redis Node را تحت فشار قرار دهد. اما اگر همان Hot Key Expire شود، می‌تواند Stampede بسیار شدیدی ایجاد کند.

```
Hot Key = Traffic Concentration
Cache Stampede = Concurrent Source Load After Miss
```

## شرایط افزایش خطر Stampede

خطر Cache Stampede در این شرایط بیشتر است:

- Key بسیار Hot باشد.
- TTL کوتاه باشد.
- Source Load گران یا کند باشد.
- Request Rate بالا باشد.
- تعداد Service Instanceها زیاد باشد.
- Cache Entry ناگهان Evict شود.
- Miss Deduplication وجود نداشته باشد.
- Retry Policy تهاجمی باشد.
- Connection Pool محدود و Queryها طولانی باشند.
- Cache Restart یا Deployment رخ دهد.
- مقدار Stale قابل ارائه نباشد.
- Refresh Ahead fail شود.
- Timeout و Backpressure مشخص نباشد.

## چه زمانی راهکار ساده کافی است؟

اگر Traffic محدود، تعداد Instanceها کم و Source Query ارزان باشد، `singleflight` داخل Process ممکن است کافی باشد:

```
Traffic = 50 Requests/sec
Instances = 2
DB Query = 10ms
```

حتی اگر هر Instance یک Query بزند، Source تحت فشار جدی قرار نمی‌گیرد.

اما در شرایط زیر:

```
Traffic = 20K Requests/sec
Instances = 30
DB Query = 800ms
```

معمولاً به ترکیبی از راهکارها نیاز داریم:

```
Refresh Ahead + Serve Stale + singleflight + TTL Jitter + Source Protection
```

## الگوی عملی پیشنهادی

برای یک Hot Key مهم، Flow مناسب می‌تواند چنین باشد:

```
Request → Cache Lookup
Fresh Hit → Return
Near Expiry Hit → Return Current Value + Trigger Deduplicated Refresh
Stale but Within Hard TTL → Return Stale + Trigger Refresh
Miss → Enter singleflight → Recheck Cache → Load Source with Timeout → Cache Result → Return
Refresh Failure → Backoff + Continue Serving Stale Within Limit
Hard TTL Reached → Fallback, Controlled Load or Error
```

این Flow هم از Miss Storm جلوگیری می‌کند و هم Request Path را تا حد امکان سریع نگه می‌دارد.

## نکات پیاده‌سازی در Go

در Go، اجرای امن این الگو معمولاً به این اجزا نیاز دارد:

```
singleflight.Group
context.Context
Source Timeout
Bounded Worker Pool
Per-Key Refresh State
TTL Jitter
Structured Logging
Metrics
Graceful Shutdown
```

نباید `context.Context` مربوط به Request را بدون تصمیم آگاهانه برای Background Refresh استفاده کرد. پس از بازگشت Response، Context درخواست معمولاً Cancel می‌شود و Refresh ممکن است نیمه‌کاره بماند. Background Work باید Context مستقل ولی محدود، Timeout مشخص و Lifecycle کنترل‌شده داشته باشد.

همچنین نباید برای هر Cache Miss یا Near-Expiry یک Goroutine آزاد ساخته شود:

```
go refresh(key)
```

این Anti-pattern تحت بار زیاد می‌تواند تعداد Goroutineها، Memory Consumption و Source Queryها را کنترل‌ناپذیر کند. استفاده از Bounded Queue و Worker Pool محدود انتخاب امن‌تری است.

`singleflight` نیز باید با دقت استفاده شود. Requestهای منتظر ممکن است Contextهای متفاوتی داشته باشند؛ بنابراین باید رفتار Cancellation، Timeout و اشتراک‌گذاری Error مشخص باشد. همچنین Error یک Loader نباید باعث Retry هم‌زمان و فوری تمام Waiterها شود؛ Cooldown یا Backoff لازم است.

## Observability

برای تشخیص Cache Stampede فقط Hit Ratio کافی نیست. Metricهای مهم عبارت‌اند از:

```
cache_hit_total
cache_miss_total
cache_miss_rate
source_load_total
source_load_inflight
source_load_duration
singleflight_shared_total
refresh_inflight
refresh_failure_total
db_pool_wait_duration
db_pool_waiters
db_query_rate
request_timeout_total
request_retry_total
stale_served_total
```

نشانه‌های محتمل Stampede:

Cache Miss Rate ↑ ناگهانی

Source Query Rate ↑ هم‌زمان

```
DB Pool Wait ↑
P99 Latency ↑
Timeout Rate ↑
Retry Rate ↑
Database CPU or I/O ↑
```

Metric مربوط به Shared Load اهمیت زیادی دارد. اگر `singleflight_shared_total` بالا باشد، نشان می‌دهد تعداد زیادی Request نتیجه یک Load مشترک را دریافت کرده‌اند و از اجرای Queryهای تکراری جلوگیری شده است.

Alerting باید ترکیبی باشد. افزایش Cache Miss به‌تنهایی ممکن است خطرناک نباشد، اما افزایش هم‌زمان Miss Rate، Source Load، Pool Wait و P99 نشانه جدی Stampede است.

## مثال Production: صفحه مسابقه فوتبال

فرض کنید API صفحه یک مسابقه فوتبال این مشخصات را دارد:

```
Traffic = 50K Requests/sec
Key = match:123:summary
Cache TTL = 30 Seconds
Database Query = 400ms
Service Instances = 20
```

اگر Key منقضی شود و هیچ Protection وجود نداشته باشد:

```
50K Requests → Cache Miss → Potentially Thousands of DB Queries
```

با `singleflight` داخل هر Instance:

```
50K Requests → At Most 20 Concurrent DB Queries
```

با Distributed Coordination ممکن است تعداد Queryها به یک Query کاهش یابد، اما ممکن است این Complexity لازم نباشد. اگر Database ظرفیت 20 Query هم‌زمان را دارد، `singleflight` محلی همراه Serve Stale و Refresh Ahead کافی است.

یک طراحی مناسب:

```
Soft TTL = 25s
Hard TTL = 2m
Near Soft TTL → One Refresh per Instance
Serve Existing Summary During Refresh
Refresh Timeout = 1s
Retry with Backoff
TTL Jitter around Base TTL
```

اگر Database Down شود، Summary قبلی برای مدت محدودی ارائه می‌شود. میزان Staleness مجاز باید براساس Business Requirement مشخص شود.

## Review Scenario: Product Detail

فرض کنید سرویس Product Detail این مشخصات را دارد:

```
Traffic for product:123 = 10K Requests/sec
Cache TTL = 5 Minutes
DB Query = 600ms
DB Pool Size = 100
Service Instances = 10
Staleness up to 30 Seconds is Acceptable
```

### مشکل اصلی چیست؟

اگر `product:123` منقضی شود، هزاران Request هم‌زمان Cache Miss دریافت می‌کنند. Connection Pool فقط 100 Connection دارد و بقیه Requestها منتظر می‌مانند. Query Latency و Pool Wait افزایش پیدا می‌کند و Timeoutها می‌توانند Retry بیشتری تولید کنند.

### بهترین ترکیب راهکار چیست؟

ترکیب مناسب:

```
Refresh Ahead + Soft/Hard TTL + Serve Stale + singleflight per Instance + TTL Jitter
```

Refresh Ahead تلاش می‌کند داده را پیش از Expiration تازه کند. اگر Refresh به‌موقع کامل نشود، مقدار قبلی تا 30 ثانیه ارائه می‌شود. `singleflight` باعث می‌شود هر Instance فقط یک Source Query اجرا کند. TTL Jitter نیز زمان Refresh و Expiration Keyهای دیگر را پخش می‌کند.

### آیا Distributed Lock ضروری است؟

الزاماً خیر. با 10 Instance، حداکثر 10 Query هم‌زمان به Database می‌رسد. اگر Database این Load را تحمل کند، Distributed Lock فقط Complexity اضافه ایجاد می‌کند. اگر حتی 10 Query سنگین نیز خطرناک باشد، Shared Lock، Leader-Based Refresh یا Shared Refresh Queue قابل بررسی است.

### اگر Refresh fail شود چه کنیم؟

```
Refresh Failure → Continue Serving Stale up to 30s → Retry with Backoff
```

پس از Hard TTL، سیستم باید براساس نیاز Business از Fallback استفاده کند، Load مستقیم ولی محدود انجام دهد یا Controlled Error برگرداند.

### چه Metricهایی باید Alert شوند؟

```
Cache Miss Rate
Source Load Inflight
DB Pool Wait
Refresh Failure Rate
Stale Served Rate
P99 Latency
Request Timeout Rate
Retry Rate
```

## جمع‌بندی و Mental Model نهایی

Cache Stampede زمانی رخ می‌دهد که تعداد زیادی Request هم‌زمان برای داده‌ای Cache نشده، منقضی‌شده یا Evict‌شده، Source اصلی را فراخوانی کنند:

```
One Cache Miss Event → Many Concurrent Source Loads
```

مشکل اصلی فقط افزایش Cache Miss نیست؛ بلکه چندبرابر شدن Load روی Database، Connection Pool، CPU، I/O و Network است. Retryهای نامناسب نیز این فشار را تشدید می‌کنند.

مهم‌ترین راهکارها عبارت‌اند از:

```
Request Coalescing
singleflight
Double-Check after Lock
Serve Stale
Refresh Ahead
Soft TTL and Hard TTL
TTL Jitter
Bounded Worker Pool
Negative Caching
Timeout and Backoff
Source-Side Protection
```

مدل ذهنی نهایی:

**در Cache Stampede، مشکل این نیست که Cache یک بار Miss شده است؛ مشکل این است که هزاران Request تصمیم می‌گیرند همان Miss را مستقل از یکدیگر جبران کنند.**

هدف طراحی صحیح:

```
Many Requests → One Controlled Load → Shared Result
```
