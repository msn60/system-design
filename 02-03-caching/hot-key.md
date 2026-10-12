# Cache Failure Mode: Hot Key

پس از بررسی Cache Stampede و Cache Avalanche، به Failure Mode بعدی یعنی Hot Key می‌رسیم. Hot Key زمانی ایجاد می‌شود که سهم بسیار بزرگی از Traffic روی یک Cache Key یا تعداد بسیار کمی Key متمرکز شود.

## مفهوم اصلی

Flow ساده Hot Key به این شکل است:

```
Huge Request Volume → Same Cache Key → Same Redis Node or Shard → Resource Saturation
```

نکته مهم این است که در Hot Key ممکن است Cache کاملاً سالم باشد و تقریباً تمام Requestها Cache Hit شوند. برخلاف Cache Stampede یا Cache Avalanche، الزاماً Cache Miss یا فشار مستقیم روی Database نداریم؛ مشکل این است که فشار روی یک نقطه کوچک از Cache Layer متمرکز شده است.

مثال:

```
Total Traffic = 100K Requests/sec
Traffic for product:123 = 70K Requests/sec
All Other Keys = 30K Requests/sec
```

در این حالت `product:123` یک Hot Key محسوب می‌شود.

## Mental Model

یک فروشگاه بزرگ را تصور کنید که چندین صندوق دارد، اما همه مشتری‌ها فقط به یک صندوق مراجعه می‌کنند. فروشگاه در مجموع ظرفیت زیادی دارد، ولی چون Traffic توزیع نشده است، همان یک صندوق به Bottleneck تبدیل می‌شود.

```
Large Total Capacity + Uneven Traffic Distribution = Local Saturation
```

در Cache نیز ممکن است Redis Cluster چندین Node داشته باشد، اما یک Key فقط روی یک Shard قرار بگیرد. اگر بیشتر Traffic مربوط به همان Key باشد، Shardهای دیگر بیکار می‌مانند و یک Shard Saturate می‌شود:

```
Shard A = 90% Traffic
Shard B = 3%
Shard C = 4%
Shard D = 3%
```

بنابراین Hot Key بیشتر یک مشکل Traffic Skew است تا کمبود ظرفیت کلی.

## Problem Statement

فرض کنید یک News Platform خبر مهمی منتشر کرده است:

```
Key = article:breaking-news
Traffic = 200K Requests/sec
Redis Cluster = 8 Shards
Average Cache Read = 1ms
```

تمام Requestها برای همان Key به یک Shard می‌روند:

```
200K Requests/sec → article:breaking-news → Shard 3
```

حتی اگر هفت Shard دیگر تقریباً خالی باشند، Shard 3 ممکن است با CPU Saturation، Network Throughput Saturation، Event Loop Delay، Connection Queue Growth، افزایش P95/P99 Latency، Timeout، Replica Lag و در نهایت Node Failure مواجه شود.

اگر Shard از دسترس خارج شود، همان Hot Key ممکن است Cache Miss شود و Hot Key Failure به Cache Stampede تبدیل شود:

```
Hot Key Saturates Shard → Shard Failure → Cache Miss → Massive Source Load
```

بنابراین Hot Key می‌تواند آغازگر Failure Modeهای دیگر باشد.

## چه چیزی یک Key را Hot می‌کند؟

Hot Key معمولاً یکی یا چند مورد از این ویژگی‌ها را دارد: تعداد Read بسیار بالا، تعداد Write بسیار بالا، Payload بزرگ، عملیات CPU-intensive مانند Serialization یا Compression، دسترسی گسترده از تعداد زیادی Instance، تمرکز Traffic در یک بازه کوتاه یا تعلق Key به یک Shard شلوغ یا کم‌ظرفیت.

Hot بودن فقط به Request Count وابسته نیست. یک Key با Payload بسیار بزرگ نیز ممکن است با Traffic کمتر Hot شود:

```
Key A = 100K Reads/sec × 1KB
Key B = 10K Reads/sec × 1MB
```

در این مثال ممکن است Key B از نظر Network و Memory Bandwidth فشار بیشتری ایجاد کند.

## علت‌های رایج Hot Key

### محتوای بسیار محبوب

نمونه‌های رایج عبارت‌اند از Breaking News، Viral Post، Popular Product، Live Match Score، Home Page Configuration و Trending List. چنین داده‌هایی ممکن است میلیون‌ها بار خوانده شوند.

### Shared Global Configuration

گاهی تمام Requestهای سیستم یک Key مشترک را می‌خوانند:

```
config:global
feature-flags:all
exchange-rate:latest
```

اگر هر Request برای دریافت تنظیمات به Redis مراجعه کند، همان Key به Hot Key تبدیل می‌شود.

### Counterهای پرترافیک

Keyهایی مانند موارد زیر ممکن است Write Rate بسیار بالایی داشته باشند:

```
video:123:view-count
post:999:like-count
campaign:77:impression-count
```

### Poor Key Design

گاهی چند داده مستقل داخل یک Key بزرگ قرار می‌گیرند:

```
all-products
all-users
global-dashboard
```

در نتیجه تمام Readها و Writeها به همان Key وابسته می‌شوند.

### Cache Stampede Recovery

پس از Cache Miss و بازسازی یک Key بسیار محبوب، تمام Requestها دوباره روی همان Key متمرکز می‌شوند. ممکن است Stampede حل شده باشد، اما Hot Key همچنان باقی بماند.

### Temporal Hotspot

بعضی Keyها همیشه Hot نیستند و فقط در بازه مشخصی Hot می‌شوند:

```
Live Match During Game
Flash Sale Product
Election Result Page
New Release Announcement
```

تشخیص و ظرفیت‌سنجی این نوع Hot Key دشوارتر است.

## Read Hot Key و Write Hot Key

Hot Key می‌تواند Read-heavy یا Write-heavy باشد و راهکارهای این دو یکسان نیستند.

### Read Hot Key

```
Many Readers → Same Key
```

مثال:

```
GET live-score:match-123
```

مشکلات اصلی شامل Network Load، Redis CPU، Connection Load، Serialization Cost و Single-Shard Bottleneck هستند. راهکارهای متداول شامل L1 Local Cache، Replication، Client-Side Cache، Key Replication، CDN و Request Coalescing در مسیرهای گران است.

### Write Hot Key

```
Many Writers → Same Key
```

مثال:

```
INCR video:123:view-count
```

مشکلات اصلی شامل Write Contention، اجرای ترتیبی Commandها، Replication Traffic، Persistence Overhead و رشد Command Queue است. راهکارهای متداول شامل Counter Sharding، Batching، Aggregation، Write Back، Stream Buffer و Partitioned Counters هستند.

## اثر Hot Key روی Redis Cluster

Redis Cluster، Keyها را براساس Hash Slot توزیع می‌کند:

```
Key → Hash Slot → Redis Node
```

یک Key واحد نمی‌تواند به‌طور خودکار بین چند Shard تقسیم شود:

```
product:123 → One Hash Slot → One Primary Node
```

افزودن Node بیشتر الزاماً Hot Key را حل نمی‌کند، چون همان Key همچنان روی یک Node قرار دارد:

```
4 Nodes → Hot Key on Node 2
8 Nodes → Hot Key Still on One Node
```

Horizontal Scaling زمانی مؤثر است که Traffic میان Keyهای متعدد توزیع شده باشد.

## اثر روی Replication

اگر Hot Key Read-heavy باشد و Read از Replicaها انجام شود، بخشی از Traffic را می‌توان توزیع کرد:

```
Primary → Replica 1 / Replica 2 → Distributed Reads
```

اما این روش Trade-offهایی مانند Replica Lag، Stale Reads، Failover Complexity، Uneven Replica Load و Replication Network Overhead دارد.

برای Write Hot Key، Replica ظرفیت Write را افزایش نمی‌دهد، زیرا Write همچنان باید ابتدا روی Primary اعمال شود.

## اثر روی Network

گاهی CPU Redis هنوز آزاد است، اما Network Interface یا Bandwidth به Bottleneck تبدیل می‌شود. فرض کنید:

```
Payload = 500KB
Requests = 20K/sec
```

حجم تقریبی خروجی:

```
500KB × 20K/sec ≈ 10GB/sec
```

در این حالت Hot Key بیشتر یک Network Hotspot است تا CPU Hotspot. Application نیز هزینه Network Copy، Decode، Allocation و Garbage Collection را پرداخت می‌کند.

## اثر روی Application

Hot Key فقط Redis را تحت فشار قرار نمی‌دهد. Application برای هر Request ممکن است این Flow را اجرا کند:

```
Redis Read → Network Copy → Deserialize → Allocate Objects → Build Response
```

اگر داده ثابت باشد ولی در هر Request دوباره Decode شود، CPU و GC Application افزایش پیدا می‌کند. در چنین شرایطی L1 Cache یا Precomputed Serialized Response می‌تواند مؤثر باشد.

## اثر روی Tail Latency

در حالت عادی:

```
Redis Read = 1ms
```

هنگام Saturation:

```
Command Queue Wait = 5ms
Network Queue = 10ms
Application Decode = 4ms
P99 = 20ms+
```

ممکن است Average Latency همچنان قابل قبول باشد، اما P99 افزایش پیدا کند. اگر همان Redis Node Keyهای دیگری نیز داشته باشد، Hot Key می‌تواند Latency Keyهای همسایه را نیز افزایش دهد:

```
One Hot Key → Node Saturation → Neighboring Keys Slow Down
```

این وضعیت نمونه‌ای از Noisy Neighbor در Cache Layer است.

## تشخیص Hot Key

Hot Key باید از روی Distribution Traffic تشخیص داده شود، نه فقط Total Request Rate.

نشانه‌های رایج:

```
One Redis Node CPU Much Higher Than Others
One Shard Network Much Higher
Few Keys Dominate Command Count
P99 High Only on One Shard
Eviction or Latency Localized to One Node
```

Metricهای مهم شامل Requests per Key، Bytes Read per Key، Bytes Written per Key، Commands per Shard، CPU per Redis Node، Network In/Out per Node، Command Latency، Key Size، Connection Count و Replication Lag هستند.

در Redis می‌توان از Hot Key Sampling، Command Statistics، Slow Log و ابزارهای Observability استفاده کرد. اجرای Scan یا تحلیل سنگین Keyspace در Production باید با احتیاط انجام شود، چون خود این عملیات می‌تواند Load ایجاد کند.

## راهکار اول: L1 Local Cache

برای Read Hot Key، یکی از مؤثرترین راهکارها نگهداری نسخه‌ای Local در هر Application Instance است:

```
Request → L1 In-Memory Cache → Redis → Database
```

فرض کنید Traffic برابر 100K Request/sec و تعداد Instanceها برابر 100 باشد. اگر هر Instance مقدار را محلی نگه دارد، Redis Load می‌تواند از 100K Read/sec به تعداد کمی Refresh در ثانیه کاهش پیدا کند:

```
Before L1 = 100K Redis Reads/sec
After L1 = About 100 Redis Refreshes/sec
```

مزایا شامل کاهش Redis Load، کاهش Network Traffic، Latency کمتر و Resilience بیشتر Application است. معایب شامل Per-Instance Staleness، Memory Duplication، Invalidation Complexity و امکان ارائه نسخه‌های متفاوت توسط Instanceهاست.

L1 Cache باید TTL کوتاه، Size Limit و Eviction Policy مشخص داشته باشد.

## L1 Cache با TTL کوتاه

برای داده‌های بسیار Hot می‌توان TTL محلی بسیار کوتاه تعریف کرد:

```
Redis TTL = 5 Minutes
L1 TTL = 1 Second
```

حتی TTL یک‌ثانیه‌ای می‌تواند Load را شدیداً کاهش دهد. فرض کنید:

```
Traffic = 50K Requests/sec
Instances = 50
L1 TTL = 1s
```

هر Instance تقریباً یک بار در ثانیه Redis را می‌خواند:

```
50K Redis Reads/sec → About 50 Redis Reads/sec
```

در مقابل، داده ممکن است حداکثر حدود یک ثانیه Stale باشد.

## Client-Side Caching

برخی Cache Systemها می‌توانند Client را از تغییر Key مطلع کنند. Client نسخه محلی را نگه می‌دارد و در زمان Invalidation آن را حذف می‌کند:

```
Redis Value → Client Local Cache
Key Changed → Invalidation Notification → Client Removes Local Copy
```

این مدل Load را کاهش می‌دهد، اما Connection State، Invalidation Delivery و Recovery را پیچیده‌تر می‌کند. چون ممکن است Notification از دست برود، TTL همچنان باید به‌عنوان Safety Net باقی بماند.

## راهکار دوم: Key Replication

برای Read Hot Key می‌توان چند Copy از همان Value با Keyهای مختلف ایجاد کرد:

```
product:123:replica:0
product:123:replica:1
product:123:replica:2
product:123:replica:3
```

Client یکی از Replica Keyها را براساس Hash یا Random انتخاب می‌کند:

```
Request → Choose Replica Key → Redis Cluster
```

اگر Replica Keyها روی Hash Slotهای متفاوت قرار بگیرند، Traffic بین چند Node توزیع می‌شود:

```
One Logical Value → Multiple Physical Keys → Multiple Shards
```

مزایا شامل Read Load Distribution، Network Distribution و کاهش فشار یک Node است. معایب شامل Write Amplification، Consistency Complexity، Invalidation Complexity، Memory بیشتر و نیاز به Replica Selection Logic است.

هر Update باید همه Replicaها را به‌روزرسانی یا Invalidate کند.

## نکته Hash Tag در Redis Cluster

Keyهایی که Hash Tag یکسان دارند روی یک Slot قرار می‌گیرند:

```
product:{123}:replica:0
product:{123}:replica:1
```

به‌دلیل `{123}` هر دو Key روی همان Shard قرار می‌گیرند و Load Distribution رخ نمی‌دهد.

برای پخش روی Slotهای مختلف نباید Hash Tag یکسان استفاده شود:

```
product:123:replica:0
product:123:replica:1
product:123:replica:2
```

در مقابل، Multi-Key Atomic Operation روی Slotهای مختلف سخت‌تر می‌شود.

## راهکار سوم: Counter Sharding

برای Write Hot Key مانند View Counter، به‌جای یک Counter می‌توان چند Counter مستقل داشت:

```
video:123:views:0
video:123:views:1
video:123:views:2
...
video:123:views:31
```

هر Writer یک Shard انتخاب می‌کند:

```
Writer → Hash or Random Shard → INCR
```

برای خواندن مقدار کل، Counterهای Shardها با هم جمع می‌شوند:

```
Total Views = Sum(All Counter Shards)

Flow:
Many Writes → Distributed Counter Shards → Periodic Aggregation
```

مثلاً:

```
1M Writes/sec
32 Counter Shards
≈ 31,250 Writes/sec per Shard
```

مزایا شامل توزیع Write، کاهش Single-Key Contention و استفاده بهتر از Cluster است. معایب شامل نیاز به Aggregation در Read، Temporary Inconsistency، تعداد Key بیشتر و پیچیدگی Reset و Migration است.

برای Readهای پرتکرار می‌توان مقدار Aggregate شده را جداگانه Cache کرد و هر چند ثانیه Update نمود.

## راهکار چهارم: Batching و Aggregation

به‌جای Update کردن Cache برای هر Event، تغییرات در Application، Queue یا Stream جمع می‌شوند:

```
1000 View Events → Aggregate Delta = +1000 → One Redis or DB Update

Flow:
Events → Local Buffer or Queue → Batch Aggregation → Cache Update
```

این روش Write Throughput را افزایش می‌دهد، اما تأخیر نمایش مقدار دقیق، Recovery Complexity و احتمال از دست رفتن داده‌های Buffer را افزایش می‌دهد.

## راهکار پنجم: CDN یا Edge Cache

اگر Hot Key محتوای عمومی و قابل ارائه در HTTP باشد، بهتر است Request حتی به Application یا Redis نرسد:

```
Client → CDN or Edge Cache → Application → Redis
```

نمونه‌های مناسب شامل Public Product Detail، News Article، Static Home Page Fragment، Popular API Response و Image Metadata هستند. CDN می‌تواند میلیون‌ها Request را در Edge پاسخ دهد.

این راهکار برای داده‌های Personalized، Private یا بسیار Dynamic ممکن است مناسب نباشد.

## راهکار ششم: Precomputed Response

گاهی Hot Key داده خام نیست، بلکه هزینه ساخت Response مشکل ایجاد می‌کند. به‌جای Cache کردن Object و Serialize کردن آن در هر Request، می‌توان Response آماده را Cache کرد:

```
Cached Domain Object → Deserialize → Transform → Serialize JSON
```

در مقابل:

```
Cached Serialized Response → Return Bytes
```

روش دوم CPU و Allocation Application را کاهش می‌دهد، اما Versioning، Content Negotiation و انعطاف Response را پیچیده‌تر می‌کند.

## راهکار هفتم: Request Coalescing

اگر Requestهای هم‌زمان علاوه بر Redis Read، Computation یا Upstream Call گرانی داشته باشند، می‌توان Work تکراری را Coalesce کرد:

```
Many Requests → One Shared Computation → Shared Result
```

اما اگر Redis Hit بسیار سریع باشد، استفاده از `singleflight` روی تمام Readها ممکن است خودش Lock Contention ایجاد کند. Request Coalescing بیشتر برای Cache Miss، Refresh یا Computation گران مناسب است.

## راهکار هشتم: Read Replica

برای Read-heavy Hot Key می‌توان Readها را میان Replicaهای Redis توزیع کرد:

```
Primary + Replica 1 + Replica 2 → Distributed Reads
```

این روش نیازمند پذیرش Staleness و Replica Lag است و Read Routing و Failover را پیچیده‌تر می‌کند. Read Replica ظرفیت Write Hot Key را افزایش نمی‌دهد.

## راهکار نهم: Dedicated Cache Instance

اگر یک Hot Key یا مجموعه کوچکی از Keyها برای Business بسیار مهم باشند، می‌توان آن‌ها را به Cache Cluster یا Instance جدا منتقل کرد:

```
General Cache → Cluster A
Hot Global Data → Cluster B
Counters → Cluster C
```

این کار Bulkhead ایجاد می‌کند و مانع می‌شود Hot Key سایر Cache Data را تحت تأثیر قرار دهد. هزینه آن افزایش Infrastructure Cost، Operational Complexity و Routing Logic است.

## راهکار دهم: Split کردن داده بزرگ

اگر Hot Key یک Value بزرگ دارد، می‌توان آن را به بخش‌های کوچک‌تر تقسیم کرد. به‌جای:

```
homepage:all-data
```

می‌توان از این Keyها استفاده کرد:

```
homepage:products
homepage:banner
homepage:trending
homepage:recommendations
```

مزایا شامل Payload کوچک‌تر، Selective Read، TTL مستقل، Refresh مستقل و Network Transfer کمتر است. در مقابل تعداد Round Tripها ممکن است افزایش پیدا کند. Pipeline یا `MGET` می‌تواند مفید باشد، اما در Redis Cluster، Multi-Key Operation روی Slotهای مختلف محدودیت دارد.

## Hot Key و Serialization

گاهی سهم اصلی Latency متعلق به Redis نیست:

```
Redis Read = 1ms
JSON Decode = 5ms
Transformation = 3ms
```

در این حالت Scale کردن Redis به‌تنهایی مشکل را حل نمی‌کند. راهکارهای مناسب شامل L1 Cache برای Object Decode شده، Cache کردن Serialized Response، استفاده از Encoding کارآمدتر، حذف Transformation تکراری و کاهش Payload است.

## Hot Key و Compression

Compression می‌تواند Network را کاهش دهد، اما CPU را افزایش می‌دهد:

```
No Compression → Higher Network, Lower CPU
Compression → Lower Network, Higher CPU
```

تصمیم باید براساس Bottleneck واقعی باشد. اگر Network Saturate شده باشد، Compression می‌تواند مفید باشد؛ اگر CPU مشکل اصلی باشد، ممکن است شرایط را بدتر کند.

## Hot Key و Expiration

اگر Hot Key Expire شود، خطر آن بسیار بیشتر از یک Key معمولی است:

```
Hot Key Hit → Redis Pressure
Hot Key Expiry → Cache Stampede
```

برای Hot Keyها معمولاً ترکیبی از Refresh Ahead، Soft TTL، Hard TTL، Serve Stale، `singleflight`، L1 Cache و TTL Jitter مناسب است. Hot Key نباید بدون برنامه Expire شود.

## Hot Key و Eviction

اگر Redis تحت Memory Pressure باشد و Hot Key Evict شود، ممکن است Stampede ایجاد شود. Metricهای مهم عبارت‌اند از:

```
used_memory
maxmemory
evicted_keys
keyspace_hits
keyspace_misses
```

Capacity Headroom اهمیت زیادی دارد و Redis نباید دائماً نزدیک سقف Memory کار کند.

## Rate Limiting

اگر Hot Key مربوط به API عمومی باشد، Rate Limiting می‌تواند Load را کاهش دهد:

```
Client, IP or User → Rate Limit → Cache
```

Rate Limit بهتر است براساس Consumer، Tenant یا Endpoint طراحی شود تا کاربران واقعی به‌صورت غیرضروری محدود نشوند.

## Load Shedding

در فشار شدید می‌توان بخشی از Requestها را رد یا پاسخ ساده برگرداند:

```
Hot Key Saturated → Serve Stale / Default / 429 / 503
```

هدف این است که Redis Node یا Application به Failure کامل نرسد.

## Graceful Degradation

اگر Hot Key مربوط به داده Optional مانند Trending Content باشد، در شرایط فشار می‌توان یکی از این رفتارها را داشت:

```
Serve Previous Trending List
Serve Static Default
Disable Personalization
Return Smaller Result Set
```

داده Optional نباید مسیرهای حیاتی سیستم را Down کند.

## Source of Truth و Freshness

در طراحی Hot Key معمولاً باید بین Freshness و Load تعادل برقرار کرد:

```
Freshness ↑ → Refresh Frequency ↑ → Load ↑
Staleness Allowed ↑ → Load ↓
```

هرچه TTL محلی کوتاه‌تر باشد، داده تازه‌تر ولی Redis Load بیشتر است. این تصمیم باید بر اساس Business Requirement انجام شود.

## Production Example: Live Match Score

فرض کنید:

```
Traffic = 100K Requests/sec
Key = match:123:score
Updates = 2/sec
Redis Cluster = 6 Nodes
Staleness Allowed = 1 Second
```

اگر همه Requestها مستقیم Redis را بخوانند:

```
100K Reads/sec → One Redis Shard
```

راهکار مناسب این است که Redis آخرین Score را نگه دارد و هر Application Instance یک L1 Cache با TTL حدود 500 میلی‌ثانیه تا یک ثانیه داشته باشد:

```
Redis stores latest score
Each Instance refreshes L1 every 500ms–1s
Requests are served locally
```

با 100 Instance:

```
100K Application Requests/sec → About 100–200 Redis Reads/sec
```

Staleness حداکثر حدود یک ثانیه است که طبق Requirement قابل قبول است.

## Production Example: Viral Video Counter

فرض کنید:

```
Views = 500K Events/sec
Key = video:999:views
```

یک `INCR` روی یک Key می‌تواند همان Redis Node را Saturate کند. راهکار مناسب:

```
32 Counter Shards
Each Event Chooses One Shard
Periodic Aggregator Sums Shards
Public Read Uses Aggregated Cached Value

Flow:
View Event → Sharded Counter → Periodic Sum → Public View Count
```

در این مدل نمایش Counter ممکن است چند ثانیه تأخیر داشته باشد، اما Write Throughput بسیار بیشتر می‌شود.

## نکات پیاده‌سازی L1 Cache در Go

یک مدل ساده می‌تواند چنین باشد:

```
type LocalEntry[V any] struct {
    Value     V
    ExpiresAt time.Time
}

type LocalCache[K comparable, V any] struct {
    mu      sync.RWMutex
    entries map[K]LocalEntry[V]
}
```

اما یک `map` ساده برای Production کافی نیست. باید Memory Limit، Eviction Policy، TTL، Concurrency Safety، Metrics، Background Cleanup و جلوگیری از Unbounded Growth در نظر گرفته شوند.

در بسیاری از موارد استفاده از Library آزمایش‌شده بهتر از ساخت Cache عمومی از صفر است. انتخاب Library باید براساس نیاز TTL، Eviction، Metrics، Admission Policy و کنترل Memory انجام شود.

## Counter Sharding در Go

برای توزیع Eventها نباید فقط خود Key ثابت Hash شود، چون همه Eventها همان Shard را انتخاب می‌کنند. باید یک عامل متغیر مانند Event ID، Worker ID یا Random Shard استفاده شود:

```
func counterShard(eventID string, shardCount uint64) uint64 {
    if shardCount == 0 {
        panic("shard count must be positive")
    }

    h := fnv.New64a()
    _, _ = h.Write([]byte(eventID))
    return h.Sum64() % shardCount
}
```

یا در صورت مناسب بودن Semantics می‌توان Shard را به‌صورت Random انتخاب کرد:

```
shard := rand.Uint64N(shardCount)
```

در نسخه‌های جدید Go، `math/rand/v2` APIهای جدیدتری مانند `Uint64N` ارائه می‌دهد و برای این نوع Load Distribution مناسب است. Random رمزنگاری‌شده برای Counter Sharding لازم نیست.

## Consistency در Key Replication

اگر یک Value در چند Physical Key تکرار شود، Update Strategy باید مشخص باشد:

```
Synchronous Update All Replicas
Asynchronous Update
Versioned Values
Invalidate All Replicas
Short TTL
```

Synchronous Update Consistency بهتری ایجاد می‌کند، اما Write Latency و Availability را کاهش می‌دهد. Asynchronous Update سریع‌تر است، ولی Replicaها ممکن است موقتاً متفاوت باشند.

استفاده از Version داخل Value می‌تواند به Reader در تشخیص نسخه قدیمی کمک کند:

```
{
  "version": 42,
  "data": {}
}
```

## تشخیص Hot Key قبل از Incident

Hotness باید به‌صورت پیوسته تحلیل شود:

```
Top Keys by QPS
Top Keys by Bytes
Top Keys by Writes
Top Keys by Error Rate
Traffic Share of Top 1 and Top 10 Keys
```

مثال:

```
Top 1 Key = 45% of Traffic
Top 10 Keys = 82% of Traffic
```

چنین Distributionی حتی در صورت سالم بودن فعلی Redis یک هشدار جدی محسوب می‌شود.

## Observability

Metricهای مهم عبارت‌اند از:

```
cache_requests_by_key
cache_bytes_read_by_key
cache_bytes_written_by_key
cache_key_size
cache_node_cpu
cache_node_network_in
cache_node_network_out
cache_command_latency
cache_connections
cache_replication_lag
l1_cache_hit_ratio
hot_key_replica_skew
counter_shard_skew
```

ثبت Label مستقیم برای تمام Keyها می‌تواند Metric Cardinality بسیار بالایی ایجاد کند. بهتر است از Sampling، Top-K، Key Pattern، Bucket یا Profiling دوره‌ای استفاده شود.

مثال Key Pattern:

```
product:{id} → key_pattern=product
article:{id} → key_pattern=article
```

برای شناسایی Key دقیق، Top-K Sampling یا ابزار Profiling لازم است.

## Alertهای مهم

```
One Redis Node CPU Much Higher Than Cluster Average
One Node Network Near Capacity
Top Key Traffic Share Above Threshold
Replica Lag Increasing
L1 Hit Ratio Dropping
Hot Key P99 Increasing
Counter Shard Imbalance
Cache Node Queue or Connection Growth
```

Alert باید Skew را نیز بررسی کند، نه فقط Total Cluster Capacity.

## تفاوت Hot Key با Cache Stampede

```
Hot Key → Key Exists and Receives Huge Traffic
Cache Stampede → Key Is Missing and Many Requests Load Source
```

Hot Key می‌تواند بدون Miss رخ دهد، اما اگر همان Key Expire یا Evict شود، ممکن است به Stampede تبدیل شود.

## تفاوت Hot Key با Cache Avalanche

```
Hot Key → Few Keys, Concentrated Traffic
Cache Avalanche → Many Keys, Broad Miss Traffic
```

در Hot Key مشکل توزیع نامتوازن Traffic است؛ در Avalanche مشکل از دست رفتن گسترده Cache است.

## تفاوت Hot Key با Cache Penetration

```
Hot Key → Existing Key with Excessive Traffic
Cache Penetration → Repeated Requests for Nonexistent Keys
```

راهکارهای اصلی:

```
Hot Key → L1 Cache, Replication, Sharding, Aggregation
Cache Penetration → Negative Caching, Bloom Filter, Input Validation
```

## Review Scenario: Exchange Rate API

فرض کنید API نرخ ارز این مشخصات را دارد:

```
Traffic = 80K Requests/sec
Key = exchange-rates:latest
Payload = 20KB
Updates = Once Every 30 Seconds
Redis Cluster = 4 Nodes
Application Instances = 80
Staleness up to 2 Seconds Is Acceptable
```

### مشکل چیست؟

تمام 80 هزار Request در ثانیه برای یک Key به یک Redis Shard می‌روند. علاوه بر QPS، حدود 1.6GB/sec داده خام از Redis خوانده می‌شود:

```
80K × 20KB ≈ 1.6GB/sec
```

این مقدار بدون در نظر گرفتن Protocol Overhead، Replication و Network Copy است؛ بنابراین Network و Redis Node می‌توانند Bottleneck شوند.

### بهترین راهکار چیست؟

استفاده از L1 Cache داخل هر Application Instance با TTL یک تا دو ثانیه مناسب است:

```
Request → L1 Cache → Redis Only on Local Expiry
```

با 80 Instance و TTL یک‌ثانیه‌ای، Redis تقریباً 80 Read/sec دریافت می‌کند، نه 80K Read/sec:

```
80K Redis Reads/sec → About 80 Redis Reads/sec
```

از آنجا که Staleness تا دو ثانیه قابل قبول است، این Trade-off منطقی است.

### آیا Key Replication لازم است؟

احتمالاً خیر، چون L1 Cache تقریباً تمام Load را حذف می‌کند. اگر L1 ممکن نباشد یا تعداد Instanceها بسیار کم باشد، Key Replication یا Read Replica قابل بررسی است.

### L1ها چگونه از Update مطلع شوند؟

گزینه‌های اصلی:

```
Short Local TTL
Pub/Sub Invalidation
Version Check
Push-Based Update
```

ساده‌ترین انتخاب معمولاً TTL کوتاه است. Pub/Sub می‌تواند Freshness را بهتر کند، اما TTL باید به‌عنوان Safety Net باقی بماند.

### چه Metricهایی مهم‌اند؟

```
L1 Hit Ratio
Redis QPS for the Key
Bytes/sec from Redis
Per-Node Network
P99 Cache Latency
Local Entry Age
Invalidation Failure
```

## Production Insight

Hot Key مشکل کمبود ظرفیت کل نیست؛ مشکل این است که ظرفیت موجود قابل استفاده نیست، چون Traffic روی یک نقطه متمرکز شده است:

```
Cluster Has Capacity ≠ Hot Key Has Capacity
```

افزودن Redis Node همیشه راه‌حل نیست. ابتدا باید مشخص شود Bottleneck در Read، Write، Network، CPU، Payload Size، Serialization یا Shard Distribution قرار دارد.

راهکار مناسب براساس نوع Hot Key انتخاب می‌شود:

```
Read Hot Key → L1 Cache, CDN, Replica, Key Replication
Write Hot Key → Counter Sharding, Batching, Aggregation, Write Back
Large Payload Hot Key → Split Data, Pre-serialize, Compression Decision
Global Shared Key → Local Cache, Push Update, Dedicated Cache
```

## جمع‌بندی و Mental Model نهایی

Hot Key زمانی رخ می‌دهد که سهم نامتناسبی از Traffic روی یک Key یا تعداد بسیار کمی Key متمرکز شود:

```
One or Few Keys → Disproportionate Traffic → One Shard or Node Saturates
```

این مشکل ممکن است در حالی رخ دهد که Cache Hit Ratio نزدیک به صد درصد باشد. بنابراین Hit Ratio بالا الزاماً به معنای سالم بودن Cache نیست.

مهم‌ترین راهکارها عبارت‌اند از:

```
L1 Local Cache
Short Local TTL
Client-Side Cache
Key Replication
Read Replica
Counter Sharding
Batching and Aggregation
CDN or Edge Cache
Payload Splitting
Precomputed Response
Dedicated Cache Cluster
Rate Limiting
Load Shedding
```

مدل ذهنی نهایی این است: **در Hot Key مشکل این نیست که داده در Cache وجود ندارد؛ مشکل این است که تمام سیستم برای دسترسی به همان داده موجود، از یک مسیر مشترک و محدود عبور می‌کند.**

نسخه زیر برای قرار دادن مستقیم در داکیومنت آماده شده است. محتوای مهم، مثال‌ها، Mental Modelها، راهکارهای Production، نکات Go، Observability و سناریوی Review حفظ شده‌اند؛ فقط فاصله‌های عمودی اضافی کم و Flowها تا حد امکان افقی شده‌اند.
