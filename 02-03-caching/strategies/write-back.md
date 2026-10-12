# Cache Strategy: Write Back

در مسیر بررسی Performance Patternهای مربوط به Caching، ابتدا Caching Fundamentals، سپس Cache Aside، Read Through و Write Through بررسی شدند. Cache Aside و Read Through بیشتر روی Read Path تمرکز داشتند، در حالی که Write Through وارد Write Path شد و Cache و Database را در مسیر هم‌زمان نوشتن قرار داد. Write Back نیز یک Write Strategy است، اما تفاوت مهم آن با Write Through این است که Database را از مسیر فوری Write خارج می‌کند.

## مفهوم اصلی

در Write Back، زمانی که داده‌ای نوشته یا تغییر داده می‌شود، Application ابتدا آن را در Cache یا یک Buffer سریع می‌نویسد و پاسخ موفقیت را بدون انتظار برای Database برمی‌گرداند. سپس Cache، Background Worker یا یک Pipeline پردازش غیرهم‌زمان، تغییرات را بعداً به Database منتقل می‌کند.

Flow اصلی این الگو به شکل زیر است:

```
Application → Write Cache → Return Success → Async Flush to Database
```

در این الگو، Database دیگر در Critical Path عملیات Write قرار ندارد؛ بنابراین caller برای کامل شدن DB Write منتظر نمی‌ماند. همین ویژگی، Write Back را از نظر کاهش Write Latency و افزایش Write Throughput بسیار جذاب می‌کند.

## Mental Model

در Write Through، هر Write باید همان لحظه هم در Cache و هم در Database اعمال شود:

```
Write Through: Write Cache + Write Database → Return Success
```

اما Write Back می‌گوید:

```
Write Back: Write Cache → Return Success → Write Database Later
```

به بیان ساده، ابتدا داده را در یک محل سریع ثبت می‌کنیم و سپس در فرصت مناسب آن را به Storage اصلی منتقل می‌کنیم.

برای درک بهتر، یک فروشگاه را در نظر بگیرید. فروشنده برای اینکه مشتری را منتظر سیستم حسابداری نگذارد، هر خرید را ابتدا در یک دفتر موقت کنار دست خود ثبت می‌کند و پاسخ فروش را فوراً می‌دهد. سپس هر چند دقیقه، اطلاعات دفتر موقت را به سیستم حسابداری اصلی منتقل می‌کند. در این مثال، Cache همان دفتر موقت سریع، Database سیستم حسابداری اصلی و Flush Worker مسئول انتقال اطلاعات از دفتر موقت به سیستم اصلی است.

## تفاوت با Write Through

در Write Through، Database همچنان بخشی از مسیر synchronous نوشتن است:

```
Application → Write Cache → Write Database → Return Success
```

یا بسته به طراحی:

```
Application → Write Database → Write Cache → Return Success
```

در Write Back، مسیر فوری Write فقط تا Cache یا Buffer ادامه پیدا می‌کند:

```
Application → Write Cache → Return Success
```

و عملیات پایدارسازی بعداً انجام می‌شود:

Cache یا Worker → Flush to Database

بنابراین تفاوت اصلی این است:

Write Through → Database در Critical Write Path قرار دارد

Write Back → Database به Async Write Path منتقل می‌شود

## دلیل تأثیر Write Back بر Performance

Database معمولاً یکی از کندترین و محدودترین Dependencyهای سیستم است. اگر هر Write مجبور باشد به Database برسد، Latency و Throughput نوشتن به ظرفیت Database محدود می‌شوند. Write Back با خارج کردن Database از مسیر فوری، پاسخ را سریع‌تر برمی‌گرداند و امکان تجمیع یا Batch کردن Writeهای بعدی را فراهم می‌کند.

فرض کنید:

```
Database Write = 50ms
Cache Write = 2ms
```

در Write Through، هر Write ممکن است بیش از `50ms` طول بکشد، زیرا Database نیز باید قبل از بازگشت پاسخ آپدیت شود. در Write Back، عملیات اولیه می‌تواند در حدود `2ms` کامل شود و Database بعداً به‌روزرسانی شود.

اثرهای اصلی عبارت‌اند از:

Write Latency ↓ شدید

Write Throughput ↑ شدید

فشار لحظه‌ای روی Database ↓

امکان Batch و Aggregation ↑

## مثال عددی: شمارنده بازدید

فرض کنید یک سیستم شمارش View برای ویدئوها داریم. اگر برای هر بازدید مستقیماً یک Database Write انجام شود:

```
1M View/min → 1M Database Write/min
```

چنین حجمی می‌تواند Database را Saturate کند. در Write Back، هر بازدید ابتدا در Redis یا یک Buffer سریع ثبت می‌شود:

```
View Event → Increment Redis Counter → Return
```

سپس هر چند ثانیه یا چند دقیقه، شمارنده‌های تجمیع‌شده به Database منتقل می‌شوند:

```
Redis Counters → Aggregate → Periodic Batch Flush → Database
```

در نتیجه، به‌جای اجرای یک میلیون Write مستقل، ممکن است فقط در بازه‌های زمانی مشخص یک مقدار تجمیع‌شده برای هر ویدئو نوشته شود. این کار تعداد عملیات Database، Network Round Trip و Transaction Overhead را به‌شدت کاهش می‌دهد.

## مثال واقعی

در سیستم‌هایی مانند پلتفرم‌های ویدئویی، شبکه‌های اجتماعی یا سرویس‌های Analytics، شمارنده‌هایی مانند View Count، Like Count و Impression Count ممکن است نرخ Write بسیار بالایی داشته باشند. نوشتن synchronous تمام این تغییرات در Database اصلی معمولاً از نظر Performance منطقی نیست.

در چنین سناریوهایی، مقدار ابتدا در یک لایه سریع مانند Redis، In-Memory Store، Stream Processor یا State Store ثبت می‌شود و سپس با Batch، Aggregation یا Background Flush به Storage پایدار منتقل می‌شود.

بااین‌حال، همین الگو برای موجودی حساب بانکی، انتقال وجه یا اطلاعات حساس پرداخت خطرناک است. اگر Cache یا Buffer قبل از Flush به Database از بین برود، ممکن است داده‌ای که سیستم موفق اعلام کرده است هرگز در Storage اصلی ثبت نشود. بنابراین Write Back برای همه انواع داده مناسب نیست.

## مزایا

مهم‌ترین مزیت Write Back کاهش شدید Write Latency است. Application منتظر تکمیل Database Write نمی‌ماند و پس از ثبت داده در Cache یا Buffer سریع می‌تواند پاسخ بدهد.

مزیت دوم افزایش Write Throughput است. Cache یا Buffer معمولاً ظرفیت Write بسیار بیشتری نسبت به Database دارد و می‌تواند Burstهای ترافیکی را جذب کند.

مزیت سوم کاهش فشار لحظه‌ای روی Database است. به‌جای ارسال هر تغییر به Database، تغییرات جمع‌آوری و با تعداد عملیات کمتر منتقل می‌شوند.

مزیت چهارم امکان Coalescing یا ادغام Writeهاست. اگر یک Counter صد بار تغییر کند، لازم نیست صد Write مستقل در Database انجام شود؛ می‌توان فقط مقدار نهایی یا Delta تجمیع‌شده را ذخیره کرد:

```
100 Counter Updates → 1 Aggregated Database Write
```

این ویژگی Write Back را برای Counterها، Metrics و داده‌های قابل تجمیع بسیار مؤثر می‌کند.

## Eventual Consistency و Dirty Data

در Write Back بین زمان نوشتن در Cache و زمان Flush به Database فاصله وجود دارد. در این بازه:

Cache = مقدار جدید

Database = مقدار قدیمی

بنابراین Write Back ذاتاً نوعی Eventual Consistency ایجاد می‌کند؛ Database در نهایت با Cache هماهنگ می‌شود، اما این هماهنگی فوری نیست.

داده‌ای که در Cache تغییر کرده ولی هنوز به Database منتقل نشده است، Dirty Data نامیده می‌شود:

Clean Data → Cache و Database هماهنگ‌اند

Dirty Data → Cache جدیدتر از Database است

یک سیستم Write Back باید بداند کدام Keyها یا Updateها Dirty هستند، چه زمانی ایجاد شده‌اند و آیا با موفقیت Flush شده‌اند یا نه. بدون Dirty Tracking قابل‌اعتماد، Write Back نمی‌تواند به‌صورت Production-grade عمل کند.

## Failure Modes

### Data Loss

خطر اصلی Write Back از دست رفتن داده‌ای است که در Cache ثبت شده اما هنوز به Database منتقل نشده است:

```
Write Cache → Return Success → Cache Crash → Database Never Updated
```

از دید client، عملیات موفق بوده است، اما پس از Crash هیچ اثری از آن در Storage اصلی باقی نمی‌ماند. این خطر زمانی بیشتر است که Cache کاملاً Volatile باشد یا Persistence و Replication مناسبی نداشته باشد.

### Stale Reads از Database

تا قبل از Flush، Database مقدار قدیمی را نگه می‌دارد. اگر بخشی از سیستم از Cache و بخش دیگری مستقیماً از Database بخواند، ممکن است دو نسخه متفاوت مشاهده شود:

```
Service A → Read Cache → New Value
Service B → Read Database → Old Value
```

بنابراین معماری باید مشخص کند Readهای مرتبط با داده Write Back شده از کدام منبع انجام شوند و چه میزان Staleness قابل قبول است.

### Flush Failure

ممکن است Background Worker نتواند داده را به Database منتقل کند؛ برای مثال Database Down باشد، Timeout رخ دهد یا Constraint و Conflict ایجاد شود. سیستم باید برای این وضعیت Retry محدود، Exponential Backoff، Jitter، Dead Letter Queue یا Recovery Plan داشته باشد.

### Ordering Problem

اگر چند Update متوالی برای یک Key ثبت شوند، ترتیب Flush اهمیت دارد. فرض کنید:

```
Update 1 → balance = 100
Update 2 → balance = 150
```

اگر Update دوم ابتدا و Update اول بعداً در Database نوشته شود، مقدار نهایی به اشتباه `100` خواهد شد. برای جلوگیری از این مشکل می‌توان از Version، Sequence Number، Timestamp قابل اعتماد، Compare-and-Set یا Idempotent Upsert استفاده کرد.

### Cache Eviction

اگر Cache تحت Memory Pressure یک Dirty Entry را پیش از Flush شدن Evict کند، داده از دست می‌رود. بنابراین استفاده از یک Cache معمولی با Eviction آزاد برای Write Back خطرناک است، مگر اینکه Dirty Data در یک مکان پایدار دیگر نیز ثبت شده باشد یا Eviction برای داده‌های Dirty ممنوع شود.

### Worker Backlog

اگر نرخ ورود Writeها از ظرفیت Flush بیشتر شود، Queue یا Dirty Set رشد می‌کند:

```
Incoming Write Rate > Database Flush Rate → Queue Growth → Flush Lag ↑
```

در ادامه Memory Consumption، Staleness و خطر Data Loss افزایش پیدا می‌کنند. این همان نقطه‌ای است که Backpressure، Load Shedding یا محدودیت پذیرش Writeها ممکن است ضروری شود.

### Process Crash و Shutdown ناقص

اگر Pending Writeها فقط داخل Memory Process نگهداری شوند، Crash یا Deployment می‌تواند آن‌ها را از بین ببرد. همچنین در Graceful Shutdown باید مشخص باشد که Pending Writeها Flush می‌شوند، به Queue پایدار منتقل می‌شوند یا برای پردازش Instance بعدی باقی می‌مانند.

## تکنیک‌های Production برای کاهش ریسک

### Durability

اگر Redis یا Cache دیگری برای نگهداری Dirty Data استفاده شود، باید درباره Persistence، Replication و Recovery آن تصمیم گرفته شود. در Redis، مکانیزم‌هایی مانند RDB Snapshot و AOF می‌توانند ریسک Data Loss را کاهش دهند، اما هرکدام Trade-off خاص خود را در Performance، Recovery Time و Durability دارند.

بااین‌حال، فعال بودن Persistence به‌تنهایی Write Back را کاملاً امن نمی‌کند؛ زیرا هنوز ممکن است بین زمان Write، Replication و Flush Window فاصله وجود داشته باشد.

### Durable Queue یا Write-Ahead Log

برای کاهش خطر Data Loss می‌توان قبل از اعلام Success، تغییر را در یک Queue یا Log پایدار ثبت کرد:

```
Write Cache → Append to Durable Queue → Return Success → Worker Flushes to Database
```

در این مدل، اگر Cache یا Application Crash کند، تغییرات از Queue قابل بازیابی هستند. Kafka، RabbitMQ با تنظیمات Durability مناسب یا یک Write-Ahead Log اختصاصی می‌توانند چنین نقشی داشته باشند.

### Batch Flush

Worker می‌تواند چندین تغییر را در یک عملیات Database ارسال کند:

```
1000 Cache Updates → 1 Batch Database Write
```

Batching تعداد Round Tripها و Transaction Overhead را کاهش می‌دهد و Throughput را افزایش می‌دهد. بااین‌حال، Batch بزرگ‌تر ممکن است Flush Latency و مقدار داده در معرض خطر را افزایش دهد؛ بنابراین Batch Size و Flush Interval باید بر اساس SLO و ظرفیت Database تنظیم شوند.

### Coalescing و Aggregation

اگر چند Update برای یک Key وجود داشته باشد، می‌توان آن‌ها را قبل از Flush ادغام کرد. برای Counterها، Deltaها جمع می‌شوند و برای داده‌های State-based، ممکن است فقط آخرین Version نوشته شود. این تکنیک حجم Write را کاهش می‌دهد، اما باید Semantics عملیات حفظ شود؛ همه عملیات‌ها قابل ادغام نیستند.

### Retry کنترل‌شده

اگر Database موقتاً در دسترس نباشد، Worker باید با Exponential Backoff و Jitter تلاش مجدد انجام دهد. Retry نامحدود و سریع می‌تواند Retry Storm ایجاد کند و وضعیت Database را بدتر کند. همچنین باید حداکثر Retry، Deadline و مسیر انتقال به Dead Letter Queue مشخص شود.

### Idempotency

Flush ممکن است به‌دلیل Timeout یا خطای نامشخص دوباره اجرا شود. بنابراین Database Write باید تا حد امکان Idempotent باشد. استفاده از Operation ID، Event ID، Version و Unique Constraint می‌تواند مانع اعمال چندباره یک تغییر شود.

### Reconciliation

حتی با وجود Queue و Retry، بهتر است یک Reconciliation Job دوره‌ای اختلاف Cache و Database را بررسی و اصلاح کند. این Job برای کشف داده‌های جاافتاده، Flushهای ناقص و Divergenceهای طولانی‌مدت مفید است.

## ترکیب Write Back با Queue

در سیستم‌های واقعی، Write Back اغلب با Message Broker یا Durable Queue ترکیب می‌شود. اگر فقط Cache آپدیت شود و Database بعداً نوشته شود، خطر Data Loss بالا باقی می‌ماند. ثبت تغییر در Queue پایدار امکان Recovery و Replay را فراهم می‌کند.

یک Flow رایج:

```
Application → Write Cache → Publish Durable Event → Return Success → Worker Consumes Event → Update Database
```

یک Flow دیگر:

```
Application → Publish Durable Event → Update Cache → Return Success → Worker Updates Database
```

ترتیب دقیق به این بستگی دارد که Queue، Cache یا Database کدام‌یک منبع اصلی پذیرش Write در نظر گرفته شود. اگر سیستم فقط پس از ثبت موفق Event در Durable Queue پاسخ Success بدهد، Durability بهتر می‌شود؛ اما Publish Latency به مسیر Write اضافه خواهد شد.

در این نقطه طراحی کمی وارد Event-Driven Architecture می‌شود، اما از دید Performance اهمیت دارد، زیرا Database را از Critical Path خارج و Queue را به‌عنوان Buffer پایدار وارد سیستم می‌کند.

## چه زمانی Write Back مناسب است؟

Write Back برای workloadهایی مناسب است که حجم Write بالا، Latency نوشتن مهم و Eventual Consistency قابل قبول باشد. همچنین زمانی بیشترین ارزش را دارد که Writeها قابل Batch، Merge یا Aggregate باشند.

نمونه‌های مناسب عبارت‌اند از View Count، Like Count، Impression Count، Analytics Counters، Metrics Aggregation، Telemetry، برخی Session Stateها و Gaming Leaderboard در سناریوهایی که تأخیر کوتاه‌مدت قابل قبول است.

Write Back معمولاً برای Bank Balance، Payment Transaction، موجودی حساس انبار، Order Final State، تغییر Credential و داده‌های Security-sensitive مناسب نیست؛ زیرا تأخیر یا از دست رفتن Persistence می‌تواند اثر جدی و غیرقابل‌قبول داشته باشد.

## مقایسه Write Through و Write Back

| **موضوع** | **Write Through** | **Write Back** |
|---|---|---|
| مسیر Write | Cache و Database در مسیر synchronous | ابتدا Cache یا Buffer، سپس Database به‌صورت async |
| Write Latency | بیشتر | کمتر |
| Write Throughput | محدودتر به Database | بیشتر |
| Consistency | قوی‌تر و سریع‌تر | Eventual |
| خطر Data Loss | کمتر | بیشتر |
| امکان Batching | محدودتر | زیاد |
| پیچیدگی عملیاتی | متوسط | زیاد |
| مناسب برای | داده‌های Hot و Read-heavy | داده‌های Write-heavy و قابل تجمیع |
| مثال | Product Catalog و User Profile | View Count، Metrics و Analytics |

## نکات پیاده‌سازی در Go

در Go نباید Write Back را صرفاً با ساختن یک Goroutine بدون مدیریت پیاده‌سازی کرد:

```
go func() {
    _ = db.Update(ctx, data)
}()
```

این Anti-pattern است؛ زیرا با Crash شدن Process، عملیات از بین می‌رود، خطا ممکن است مشاهده نشود، Context اصلی ممکن است قبل از پایان کار Cancel شود و هیچ Retry یا Backpressure مشخصی وجود ندارد.

Flow مناسب‌تر برای سیستم Production:

```
Write Cache → Enqueue Durable Job → Return → Worker Pool → Batch Database Write
```

اگر Queue داخل Process باشد، باید Bounded باشد تا با افزایش ترافیک Memory بدون محدودیت رشد نکند. بااین‌حال، Bounded Channel درون Process همچنان Durable نیست و با Crash از بین می‌رود؛ بنابراین برای داده مهم باید از Queue پایدار استفاده شود.

در Go، پیاده‌سازی باید شامل Context Timeout برای Database، Worker Pool محدود، Retry با Exponential Backoff و Jitter، Idempotent Database Write، Graceful Shutdown، ثبت Error و Metric، محدودیت Queue و مدیریت Backpressure باشد.

در Graceful Shutdown باید پذیرش Job جدید متوقف شود، Jobهای Pending تا Deadline مشخص Flush شوند یا اطمینان حاصل شود که در Queue پایدار قرار گرفته‌اند. استفاده صرف از `time.Sleep` یا انتظار نامحدود برای Workerها مناسب نیست؛ باید از `context.Context`، `sync.WaitGroup` یا ساختارهای Lifecycle کنترل‌شده استفاده شود.

اگر عملیات‌ها Batch می‌شوند، Batch Size و Flush Interval باید Configurable باشند. Batch بسیار کوچک منفعت Performance را کاهش می‌دهد و Batch بسیار بزرگ Staleness، Memory Consumption و Data-at-Risk را افزایش می‌دهد.

## Observability

Write Back بدون Observability بسیار خطرناک است، زیرا ممکن است Application سریع پاسخ دهد اما Database به‌تدریج عقب بماند. متریک‌های مهم عبارت‌اند از:

- `dirty_keys_count`: تعداد Keyهایی که هنوز Flush نشده‌اند.
- `flush_queue_depth`: تعداد Jobهای در انتظار.
- `flush_lag_seconds`: فاصله زمانی Cache و Database.
- `oldest_unflushed_update_age`: سن قدیمی‌ترین Update پردازش‌نشده.
- `flush_success_total` و `flush_failure_total`: تعداد Flushهای موفق و ناموفق.
- `cache_write_duration`: Latency نوشتن در Cache.
- `db_flush_duration`: Latency عملیات Database.
- `retry_count`: تعداد Retryهای Flush.
- `dead_letter_count`: تعداد Jobهایی که پس از Retryهای مجاز منتقل شده‌اند.
- `cache_db_divergence_count`: تعداد اختلاف‌های کشف‌شده بین Cache و Database.

افزایش `flush_lag_seconds` یا `oldest_unflushed_update_age` نشان می‌دهد Database از Cache عقب افتاده است. رشد مداوم `flush_queue_depth` نشان می‌دهد ظرفیت مصرف کمتر از نرخ تولید است. افزایش `flush_failure_total` و `retry_count` نیز می‌تواند نشانه خرابی یا Saturation در Database باشد.

Alerting نباید فقط بر اساس Error انجام شود. ممکن است هیچ Flushی fail نشده باشد، اما Lag به‌قدری زیاد شده باشد که SLO مربوط به Freshness نقض شود. بنابراین Lag و Age نیز باید SLI یا حداقل Operational Indicator باشند.

## Production Insight

Write Back یکی از مؤثرترین و در عین حال خطرناک‌ترین Performance Patternهاست. این الگو با خارج کردن Database از مسیر فوری Write، Latency را کاهش و Throughput را افزایش می‌دهد، اما در مقابل Durability، Ordering، Recovery و Consistency را پیچیده‌تر می‌کند.

جمله کلیدی این الگو چنین است:

Write Back سرعت را با تأخیر در Persistence می‌خرد.

برای Counterها، Analytics و Metrics این Trade-off اغلب قابل قبول است. اما برای پول، وضعیت نهایی سفارش، Credential و داده‌های حساس، این تأخیر یا احتمال Data Loss معمولاً قابل قبول نیست.

## مثال Review: Video View Counter

فرض کنید یک سرویس شمارش بازدید ویدئو داریم:

```
100K View/sec
Database Capacity = 5K Write/sec
Accuracy Requirement = Eventual
Maximum Acceptable Delay = 1 Minute
```

### آیا Write Back مناسب است؟

بله. نرخ Write ورودی بیست برابر ظرفیت Database است و Database نمی‌تواند هر View را جداگانه ذخیره کند. از طرف دیگر، دقت کاملاً لحظه‌ای نیاز نیست و تأخیر تا یک دقیقه قابل قبول است. بنابراین می‌توان Viewها را در Redis، State Store یا Stream Processor تجمیع و سپس به‌صورت دوره‌ای به Database منتقل کرد.

Flow ساده:

```
View Event → Increment Redis Counter → Return → Periodic Aggregate Flush → Database
```

Flow مقاوم‌تر:

View Event → Durable Queue → Aggregator → Redis یا State Store → Periodic Database Flush

در Flow دوم، Event ابتدا در یک Queue پایدار ثبت می‌شود؛ بنابراین Crash شدن Redis یا Aggregator الزاماً باعث از دست رفتن Viewها نمی‌شود.

### چه ریسک‌هایی وجود دارد؟

اگر Redis بدون Persistence یا Durable Queue استفاده شود، Viewهایی که هنوز Flush نشده‌اند ممکن است از دست بروند. اگر Flush Worker از نرخ ورودی عقب بماند، Database مقدار قدیمی نشان خواهد داد و Queue Depth افزایش پیدا می‌کند. اگر ترتیب یا منطق Aggregation اشتباه باشد، Counter نهایی نادرست می‌شود. همچنین ممکن است بعضی ویدئوهای بسیار محبوب Hot Key ایجاد کنند و یک Shard یا Redis Node را تحت فشار قرار دهند.

### چگونه ریسک‌ها کاهش پیدا می‌کنند؟

برای کاهش Data Loss می‌توان از Durable Queue، Redis Persistence و Replication استفاده کرد. برای افزایش Throughput باید Batching و Aggregation انجام شود. Database Writeها باید Idempotent باشند تا Retry باعث Double Count نشود. Counterهای Hot می‌توانند بین چند Key یا Shard تقسیم و سپس Merge شوند. همچنین باید Flush Lag، Queue Depth، Oldest Pending Event، Retry Rate و اختلاف Cache و Database مانیتور شوند.

یک Reconciliation Job دوره‌ای نیز می‌تواند شمارنده‌های Database را با Log یا State پایدار مقایسه و اختلاف‌ها را اصلاح کند.

## جمع‌بندی و مدل نهایی

Write Back یعنی برای سریع‌تر شدن مسیر Write، داده ابتدا در Cache یا Buffer سریع ثبت و Database بعداً به‌صورت Async به‌روزرسانی می‌شود:

```
Write Back: Write Fast Now → Persist Later
```

مقایسه نهایی:

Write Through: Cache و Database اکنون نوشته می‌شوند → امن‌تر ولی کندتر

Write Back: Cache اکنون و Database بعداً نوشته می‌شود → سریع‌تر ولی پرریسک‌تر

Write Back زمانی مناسب است که حجم Write بالا، Write Latency مهم، داده قابل Batch یا Aggregate و Eventual Consistency قابل قبول باشد. اگر Durability فوری، Ordering دقیق یا Strong Consistency حیاتی باشد، Write Back معمولاً انتخاب مناسبی نیست.

مدل ذهنی نهایی این است: **Write Back، Database را از Critical Path نوشتن خارج می‌کند و Performance را افزایش می‌دهد، اما مسئولیت حفظ داده را از Database به معماری Queue، Worker، Persistence، Retry، Idempotency و Observability منتقل می‌کند.**

Write Back برای workloadهایی مناسب است که **حجم Write بالا، حساسیت به Write Latency زیاد، امکان Batch/Aggregation وجود داشته باشد و Eventual Consistency قابل قبول باشد**. این الگو برای داده‌هایی نامناسب است که **نباید گم شوند، باید فوراً پایدار شوند یا Strong Consistency و Ordering دقیق نیاز دارند**.

## مناسب برای

- **View Count و Impression Count:** مانند شمارش بازدید ویدئو یا نمایش تبلیغ؛ می‌توان رویدادها را در Redis یا یک Buffer سریع جمع کرد و هر چند ثانیه به DB منتقل کرد.
- **Like Count و Reaction Count:** اگر تأخیر چندثانیه‌ای در نمایش عدد نهایی قابل قبول باشد.
- **Metrics و Telemetry:** مانند CPU metrics، application metrics و event counters که معمولاً قابل Batch و Aggregate هستند.
- **Analytics Counters:** مثل تعداد کلیک، بازدید صفحه، conversion count و engagement metrics.
- **Gaming Leaderboard در برخی سناریوها:** وقتی تغییر رتبه با کمی تأخیر قابل قبول باشد و State در یک Store سریع نگهداری شود.
- **Session State با امکان Recovery:** در صورتی که از دست رفتن محدود داده قابل قبول باشد یا Queue پایدار برای بازیابی وجود داشته باشد.
- **High-write workloads قابل تجمیع:** مانند Counterهایی که به‌جای هزاران Write مستقل، می‌توان Delta یا مقدار نهایی آن‌ها را ذخیره کرد.
- **سیستم‌های Eventual Consistent:** جایی که DB می‌تواند چند ثانیه یا چند دقیقه از Cache عقب‌تر باشد.

مثال مناسب:

```
100K View/sec
DB Capacity = 5K Write/sec
```

Delay تا 1 دقیقه قابل قبول

در این حالت Write Back منطقی است:

```
View Event → Increment Redis Counter → Return → Periodic Batch Flush → DB
```

## نامناسب برای

- **Bank Balance و Ledger مالی:** موجودی و تراکنش‌های مالی نباید بعد از اعلام موفقیت گم شوند یا با تأخیر نامطمئن ثبت شوند.
- **Payment Transaction:** تأیید پرداخت بدون Persistence فوری می‌تواند به ناسازگاری یا از دست رفتن تراکنش منجر شود.
- **Order Final State:** وضعیت نهایی سفارش، پرداخت یا تحویل معمولاً باید Durable و قابل Audit باشد.
- **Inventory حساس به Oversell:** اگر به‌روزرسانی موجودی دیر به DB برسد، ممکن است یک کالا چند بار فروخته شود.
- **Credential و Security Data:** تغییر Password، Token، Permission یا Access Policy نباید فقط در Cache باقی بماند.
- **داده‌های نیازمند Strong Consistency:** جایی که همه Readerها باید فوراً یک مقدار یکسان ببینند.
- **عملیات وابسته به Ordering دقیق:** اگر جابه‌جایی ترتیب Updateها نتیجه نهایی را خراب کند.
- **داده‌های غیرقابل تکرار یا غیرقابل جبران:** جایی که در صورت Crash امکان Replay یا Reconciliation وجود ندارد.

مثال نامناسب:

```
Transfer Money → Update Cache → Return Success → Cache Crash
```

در این حالت ممکن است سیستم انتقال وجه را موفق اعلام کرده باشد، اما تغییر هرگز در Database یا Ledger پایدار ثبت نشود.

## معیار تصمیم‌گیری سریع

Write Back مناسب است اگر بیشتر پاسخ‌ها «بله» باشند:

- Write Rate بالاست؟
- Write Latency اهمیت زیادی دارد؟
- Writeها قابل Batch، Merge یا Aggregate هستند؟
- تأخیر در Persistence قابل قبول است؟
- Eventual Consistency قابل قبول است؟
- Durable Queue، Retry و Recovery داریم؟
- از دست رفتن محدود داده قابل جبران است؟

Write Back مناسب نیست اگر یکی از این موارد حیاتی باشد:

- Durability فوری
- Strong Consistency
- Auditability
- Ordering دقیق
- عدم تحمل Data Loss
- صحت مالی یا امنیتی

**مدل نهایی:** Write Back برای داده‌های پرتعداد، قابل‌تجمیع و کم‌ریسک مناسب است؛ نه برای داده‌هایی که اعلام موفقیت آن‌ها باید به معنی ثبت قطعی و فوری در Storage اصلی باشد.
