# Cache Warming

پس از بررسی Failure Modeهای اصلی Cache شامل Cache Stampede، Cache Avalanche، Hot Key و Cache Penetration، به مبحث Cache Warming می‌رسیم. Cache Warming یعنی داده‌های مهم و پرتکرار را پیش از آنکه Request واقعی مجبور شود آن‌ها را از Source اصلی بارگذاری کند، داخل Cache قرار دهیم.

```
Service Startup / Deployment / Cache Recovery
→ Identify Important Data
→ Load Data from Source
→ Populate Cache
→ Accept or Increase Traffic
```

هدف این است که سیستم در لحظه Startup، Restart، Deployment یا Recovery با Cache کاملاً خالی وارد Traffic واقعی نشود. بدون Warming، ورود Requestها به Cache خالی می‌تواند تعداد زیادی Cache Miss ایجاد کند و فشار شدیدی به Database یا Source اصلی وارد سازد:

```
Cache Empty → Real Requests Arrive → Massive Cache Misses
→ Database Load Spike → Latency and Timeout Increase
```

در مقابل، Cache Warming ابتدا بخشی از داده‌های مهم را آماده می‌کند و سپس Traffic به‌صورت کنترل‌شده وارد سیستم می‌شود:

```
Cache Empty → Preload Critical Data → Traffic Starts Gradually
→ Important Requests Hit Cache
```

## Mental Model

یک رستوران را تصور کنید که ساعت ۸ شب وارد شلوغ‌ترین بازه کاری خود می‌شود. اگر آشپزخانه دقیقاً هنگام ورود مشتریان شروع به آماده‌سازی مواد اولیه کند، اولین سفارش‌ها با تأخیر زیادی مواجه می‌شوند و سفارش‌های جدید نیز روی آن‌ها انباشته خواهند شد. روش بهتر این است که مواد اولیه، سس‌ها و اقلام پرمصرف پیش از آغاز ساعت شلوغی آماده شوند.

```
No Warming → Customers Arrive Before Preparation
Cache Warming → Preparation Happens Before Peak Traffic
```

در Cache نیز لازم نیست همه داده‌های سیستم از قبل Load شوند؛ هدف اصلی این است که داده‌های حیاتی و Hot پیش از ورود Traffic آماده باشند.

## زمان‌های نیاز به Cache Warming

Cache Warming معمولاً در شرایط زیر اهمیت پیدا می‌کند:

```
Application Startup
Redis Restart
Cache Flush
Deployment
Failover
Cache Key Version Change
Large Migration
Traffic Spike Preparation
Scheduled Peak Hours
Disaster Recovery
New Region Launch
```

برای مثال، بعد از Restart شدن Redis ممکن است Cache کاملاً خالی باشد:

```
Redis Starts Empty → Application Instances Receive Traffic
→ Every Request Misses → Database Becomes Temporary Cache Builder
```

اگر Traffic بالا باشد، Database ممکن است قبل از آنکه Cache به‌اندازه کافی پر شود Saturate شود.

## Cold Cache و Cold Start

Cache خالی یا Cacheای که بخش بزرگی از داده‌های پرتکرار را ندارد، Cold Cache نامیده می‌شود:

```
Cold Cache → Low Hit Ratio → High Source Load
Warm Cache → High Hit Ratio → Lower Source Load
```

Cold Start فقط به Startup فیزیکی Application یا Redis محدود نیست. شرایط زیر نیز می‌توانند یک Cold Start منطقی ایجاد کنند:

```
New Cache Namespace
New Key Prefix
New Region
New Application Version
New Serialization Format
Shard Replacement
Mass Invalidation
```

ممکن است Redis میلیون‌ها Key داشته باشد، اما نسخه جدید Application به‌دلیل تغییر Prefix، Serialization یا Schema نتواند از آن‌ها استفاده کند. در این حالت Cache از نظر فیزیکی پر، اما از دید Application جدید Cold است.

## Problem Statement

فرض کنید یک فروشگاه اینترنتی مشخصات زیر را دارد:

```
Products = 5 Million
Hot Products = 20,000
Traffic = 40K Requests/sec
Cache Hit Ratio Before Restart = 95%
Database Safe Capacity = 5K Queries/sec
```

پیش از Restart، فقط پنج درصد Requestها به Database می‌رسند:

```
40K Requests/sec × 5% Miss = 2K DB Queries/sec
```

Database در محدوده امن خود قرار دارد. اما پس از Restart شدن Redis:

```
Cache Hit Ratio ≈ 0%
40K Requests/sec → Up to 40K DB Queries/sec
```

این Load حدود هشت برابر ظرفیت امن Database است و می‌تواند این Failure Chain را ایجاد کند:

```
Cold Cache → Mass Miss → DB Pool Saturation → Query Latency Increase
→ Request Timeout → Retry Amplification → Recovery Slows Down
```

Cache Warming باید حداقل Hot Productها و داده‌های حیاتی را پیش از ورود کامل Traffic بارگذاری کند.

## تفاوت Cache Warming با Preloading کامل

یک تصور اشتباه این است که برای Cache Warming باید کل Database وارد Cache شود. این کار معمولاً پرهزینه و غیرضروری است:

```
Millions of Rows → Huge Source Load → Large Cache Memory
→ Long Startup → Many Unused Entries
```

هدف درست این است:

```
Warm the Most Valuable Data, Not All Data
```

یعنی داده‌هایی ابتدا Warm شوند که بیشترین اثر را بر Traffic، Latency، ظرفیت Source یا Business دارند.

## انتخاب داده‌های مناسب برای Warming

Candidateهای رایج برای Warming عبارت‌اند از:

```
Hot Keys
Frequently Accessed Products
Global Configuration
Feature Flags
Home Page Data
Popular Articles
Live Events
Tenant Configuration
Reference Data
Exchange Rates
Routing Tables
Session Metadata for Active Users
```

داده مناسب برای Warming معمولاً یک یا چند ویژگی زیر را دارد:

```
High Read Frequency
High Source Load Cost
High Business Criticality
Predictable Access Pattern
Small or Moderate Size
Low Update Frequency
Acceptable Staleness
```

Warm کردن داده‌هایی که احتمال استفاده از آن‌ها پایین است، فقط Source Load و Memory Consumption را افزایش می‌دهد.

## Hot-Key Prioritization

همه Keyها ارزش یکسانی ندارند. فرض کنید:

```
Total Keys = 10 Million
Top 10K Keys = 80% of Traffic
```

Warm کردن همان ده هزار Key می‌تواند بیشتر Traffic را پوشش دهد:

```
Warm Top 10K Keys → Protect 80% of Traffic
→ Let Long-Tail Data Load Lazily
```

این روش بسیار بهتر از Warm کردن کل ده میلیون Key است.

فهرست Hot Keyها می‌تواند از منابع زیر استخراج شود:

```
Historical Access Logs
Top-K Metrics
Database Query Frequency
Business Priority List
Recent Popularity
Scheduled Events
Product Campaigns
Previous Cache Snapshot
```

صرفاً تعداد Request نیز نباید معیار باشد. ممکن است یک Key کم‌درخواست، Query بسیار سنگینی داشته باشد یا برای Business حیاتی باشد. بنابراین اولویت می‌تواند ترکیبی از Frequency، Source Cost و Business Criticality باشد.

## Eager Warming و Lazy Loading

دو رویکرد اصلی برای پر کردن Cache وجود دارد.

### Eager Warming

در Eager Warming، داده پیش از ورود Request واقعی بارگذاری می‌شود:

```
Startup → Load Important Keys → Mark Service Ready
```

مزایا:

```
Predictable Initial Latency
Higher Initial Hit Ratio
Lower Cold-Start Load
Better Source Protection
```

معایب:

```
Longer Startup
Source Load During Warming
Possibility of Loading Unused Data
Failure-Handling Complexity
```

### Lazy Loading

در Lazy Loading، داده هنگام اولین Request بارگذاری می‌شود:

```
First Request → Cache Miss → Load Source → Cache Result
```

مزایا:

```
Simple
Only Used Data Is Cached
Faster Startup
```

معایب:

```
Cold Requests Are Slow
Initial Traffic Hits Source
Stampede Risk
Unpredictable Recovery
```

در Production معمولاً طراحی ترکیبی مناسب‌تر است:

```
Critical and Hot Data → Eager Warming
Long-Tail Data → Lazy Loading
```

این مدل هم از Source محافظت می‌کند و هم از Load و Memory غیرضروری جلوگیری می‌کند.

## Blocking Warming و Background Warming

یکی از تصمیم‌های مهم این است که آیا Service تا پایان Warming از دریافت Traffic جلوگیری کند یا هم‌زمان با Warming وارد مدار شود.

### Blocking Warming

```
Process Starts → Warm Required Data → Mark Ready → Receive Traffic
```

این روش برای داده‌هایی مناسب است که بدون آن‌ها Service نمی‌تواند رفتار صحیح یا ایمنی داشته باشد.

مزایا:

```
Required Cache State Guaranteed
No Traffic Before Critical Data Exists
```

معایب:

```
Slow Startup
Deployment Delay
Risk of Readiness Never Becoming True
```

اگر Warming صدها هزار Key طول بکشد یا Source دچار اختلال باشد، Deployment ممکن است برای مدت طولانی متوقف بماند.

### Background Warming

```
Process Starts → Becomes Ready → Warms Cache in Background
```

مزایا:

```
Fast Startup
Warming Failure Does Not Fully Block Deployment
```

معایب:

```
Real Traffic Competes with Warming
Initial Hit Ratio Is Low
Source May Be Overloaded
```

طراحی متعادل‌تر معمولاً چندمرحله‌ای است:

```
Start → Warm Minimum Critical Set → Mark Ready
→ Continue Background Warming → Gradually Ramp Traffic
```

## Minimum Viable Warm Set

لازم نیست Service برای Warm شدن همه داده‌ها صبر کند. می‌توان مجموعه‌ای حداقلی از داده‌های ضروری تعریف کرد:

```
Global Configuration
Top 1,000 Products
Homepage Data
Critical Tenant Configuration
Current Exchange Rates
```

پس از آماده شدن این مجموعه، Service می‌تواند Ready شود و Warming باقی داده‌ها در Background ادامه پیدا کند:

```
Critical Warm Set Complete → Ready
Long-Tail Warming Continues → Background
```

Minimum Viable Warm Set باید کوچک، مشخص و قابل‌اندازه‌گیری باشد. اگر این مجموعه بیش از حد بزرگ انتخاب شود، عملاً دوباره به Blocking Warming کامل تبدیل می‌شود.

## Gradual Warming

Cache نباید با بیشترین سرعت ممکن Warm شود. طراحی زیر خطرناک است:

```
Start 10,000 Goroutines → Load 1 Million Rows → DB Saturation
```

Warming خودش یک Load Generator است و می‌تواند Source را Down کند.

طراحی بهتر:

```
Keys → Batches → Bounded Workers
→ Rate-Limited Source Reads → Bounded Cache Writes
```

برای مثال:

```
Batch Size = 500
Concurrent Workers = 20
Maximum DB Reads = 1,000/sec
Delay Between Batches = 100ms
```

مقادیر مناسب باید با Load Test و مشاهده Metricهای Production تعیین شوند و عدد ثابت عمومی برای همه سیستم‌ها وجود ندارد.

## Capacity Budget

فرض کنید Database چنین وضعیتی دارد:

```
Safe Capacity = 10K Queries/sec
Normal Live Load = 6K Queries/sec
Safety Margin = 2K Queries/sec
```

ظرفیت مناسب Warming احتمالاً فقط دو هزار Query در ثانیه است:

```
Warming Budget = 2K Queries/sec
```

چهار هزار Query ظرفیت ظاهراً آزاد وجود دارد، اما دو هزار مورد آن باید برای Traffic Spike، Queryهای کند و تغییرات طبیعی Load حفظ شود.

مدل ظرفیت:

```
Source Capacity
= Live Traffic
+ Warming Traffic
+ Retry Traffic
+ Safety Headroom
```

اصل اصلی این است که Warming باید از Spare Capacity استفاده کند، نه اینکه با Traffic حیاتی رقابت کند.

## Concurrency Limiting

برای جلوگیری از فشار بیش از حد، تعداد Loadهای هم‌زمان باید محدود شود. یک مدل ساده در Go:

```
type WarmLimiter struct {
	sem chan struct{}
}

func NewWarmLimiter(limit int) *WarmLimiter {
	if limit <= 0 {
		panic("warm concurrency limit must be positive")
	}

	return &WarmLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (l *WarmLimiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *WarmLimiter) Release() {
	<-l.sem
}

Flow:
Warm Job → Acquire Slot → Load Source
→ Write Cache → Release Slot
```

در پیاده‌سازی واقعی، `Release` باید فقط پس از `Acquire` موفق و معمولاً با `defer` فراخوانی شود.

برای سناریوهای عمومی‌تر یا Weighted Concurrency می‌توان از `golang.org/x/sync/semaphore` نیز استفاده کرد.

## تفاوت Concurrency Limit و Rate Limit

این دو مفهوم یکسان نیستند:

Concurrency Limit → چند عملیات هم‌زمان اجرا شوند؟

Rate Limit → چند عملیات در واحد زمان آغاز شوند؟

برای مثال:

```
20 Concurrent Queries
Each Query = 2ms
→ Potentially Very High QPS
```

بنابراین محدود کردن Concurrency به 20 لزوماً QPS را به 20 محدود نمی‌کند. در Warming حساس بهتر است هر دو کنترل وجود داشته باشند:

```
Concurrency Limit + QPS Rate Limit
```

در Go می‌توان از `golang.org/x/time/rate` برای Rate Limiting استفاده کرد.

## Batch Loading

به‌جای اجرای Query جداگانه برای هر Key، در صورت پشتیبانی Source می‌توان داده‌ها را به‌صورت Batch دریافت کرد.

روش ضعیف:

```
1,000 IDs → 1,000 Separate Queries
```

روش بهتر:

```
1,000 IDs → 10 Queries × 100 IDs
```

نمونه Query:

```
SELECT product_id, name, price
FROM products
WHERE product_id IN (:ids);
```

مزایا:

```
Fewer Network Round Trips
Lower Query Parsing Overhead
Better Throughput
```

معایب:

```
Large Queries
Higher Temporary Memory Usage
Uneven Batch Latency
Partial Result Handling
Database Parameter Limits
```

Batch Size باید با Benchmark و Load Test تعیین شود. Batch بسیار کوچک مزایای Batching را کاهش می‌دهد و Batch بسیار بزرگ می‌تواند Latency، Memory و Lock Duration را افزایش دهد.

## Pipeline برای Cache Write

هنگام استفاده از Redis می‌توان Cache Writeها را Pipeline کرد:

```
Load Batch from DB → Build Cache Entries
→ Pipeline SET Commands → Flush
```

Pipeline تعداد Network Round Tripها را کاهش می‌دهد، اما Pipeline بسیار بزرگ می‌تواند مشکلات زیر را ایجاد کند:

```
Memory Usage ↑
Latency Spike ↑
Redis Event Loop Delay ↑
Failure Recovery Complexity ↑
```

بنابراین تعداد Commandها و حجم Byteهای هر Pipeline باید محدود باشد.

## TTL Jitter در Cache Warming

اگر همه Keyها هنگام Warming با TTL یکسان نوشته شوند، بعداً تقریباً هم‌زمان Expire خواهند شد و Cache Avalanche ایجاد می‌شود.

روش خطرناک:

```
100K Keys Warmed at 10:00
TTL = 1 Hour
→ 100K Keys Expire around 11:00
```

روش بهتر:

```
Base TTL = 1 Hour
Jitter = ±10%
```

نمونه:

```
Key A = 55m
Key B = 63m
Key C = 58m
Key D = 65m
```

رابطه مهم میان این دو مبحث:

```
Cache Warming Without TTL Jitter
→ Future Cache Avalanche
```

Jitter باید با Freshness Requirement سازگار باشد و نباید عمر داده را بیش از حد مجاز افزایش دهد.

## ترتیب Warming

ترتیب Warming باید براساس ارزش داده و احتمال استفاده تعیین شود:

```
Priority 1 → Global Configuration
Priority 2 → Current Hot Products
Priority 3 → Homepage Data
Priority 4 → Recently Active Users
Priority 5 → Long-Tail Data
```

می‌توان Candidateها را در Priority Queue قرار داد:

```
High-Value Keys → Warm First
Low-Value Keys → Warm Later or Lazily
```

اولویت می‌تواند براساس ترکیبی از عوامل زیر محاسبه شود:

```
Traffic Frequency
Business Criticality
Source Cost
Recency
Expected Campaign
Data Size
Freshness Requirement
```

برای مثال، Key کوچک با Traffic بالا و Query گران باید پیش از Key بزرگی که احتمال استفاده کمی دارد Warm شود.

## Traffic Ramp-Up

حتی پس از پایان Warming اولیه، بهتر است Traffic به‌تدریج وارد شود:

```
10% Traffic → Observe Metrics → 25% → 50% → 100%
```

این روش اجازه می‌دهد پیش از ورود Traffic کامل، Hit Ratio، Database Load، Redis Latency و Error Rate بررسی شوند.

در محیط‌هایی مانند Kubernetes می‌توان Readiness، Canary Deployment و Progressive Delivery را با Warming هماهنگ کرد.

اصل مهم:

```
Process Ready ≠ Cache Ready ≠ System Ready
```

ممکن است Process سالم و Port آن باز باشد، اما Cache هنوز Cold باشد و Source ظرفیت Traffic کامل را نداشته باشد.

## Warming در چند Application Instance

اگر 100 Instance هم‌زمان Start شوند و هرکدام Warming یکسانی را اجرا کنند:

```
100 Instances × Same Warm Job
→ Duplicate Source Loads
→ Duplicate Cache Writes
```

این وضعیت Warming Storm نامیده می‌شود.

راهکارهای رایج:

```
Single Leader Warmer
Distributed Lock
Dedicated Warming Service
Partitioned Warm Work
Shared Job Queue
Idempotent Workers
```

## Leader-Based Warming

در این مدل یک Instance به‌عنوان Leader انتخاب می‌شود:

```
Leader → Executes Warming
Followers → Wait or Serve Controlled Traffic
```

مزیت اصلی:

```
No Duplicate Work
```

معایب:

```
Leader Failure
Election Complexity
Single Warming Throughput
```

در صورت Failure رهبر، Job باید قابل Resume باشد و Leader جدید نباید تمام Work موفق را از ابتدا تکرار کند.

## Partitioned Warming

در این مدل هر Worker یا Instance بخشی از Keyها را Warm می‌کند:

```
Instance 1 → Partition 0
Instance 2 → Partition 1
Instance 3 → Partition 2
```

تقسیم ساده:

```
partition = hash(key) % partitionCount
```

مزیت این مدل افزایش Throughput و توزیع Work است. معایب آن شامل Membership، Rebalancing، Duplicate Processing و مدیریت Failure یک Partition است.

بهتر است تعداد Partitionها از تعداد Instanceها مستقل و بیشتر باشد تا در تغییر تعداد Workerها، توزیع Work انعطاف‌پذیر باقی بماند.

## Dedicated Warming Worker

در سیستم‌های بزرگ می‌توان Warming را از Application اصلی جدا کرد:

```
Scheduler → Warming Queue → Dedicated Workers → Cache
```

مزایا:

```
Independent Scaling
Controlled Concurrency
Clear Observability
No Startup Coupling
```

معایب:

```
Additional Infrastructure
Queue Management
Consistency Complexity
Ownership Complexity
```

این مدل برای Warming گسترده، Scheduled Warming و چند Consumer معمولاً قابل‌کنترل‌تر است.

## Idempotency

Warming Job باید Idempotent باشد؛ یعنی اجرای دوباره آن نتیجه مخرب یا متفاوتی ایجاد نکند.

عملیات زیر معمولاً Idempotent است:

```
SET product:123 = current-value
```

اما عملیات زیر برای Warming مناسب نیست:

```
INCR product:123:view-count
```

اجرای دوباره `INCR` مقدار را به‌اشتباه افزایش می‌دهد. Warming باید State را Set یا Replace کند، نه اینکه Mutation غیرقابل‌تکرار ایجاد کند.

## Version Awareness

در Deployment جدید ممکن است Cache Key یا Payload Version تغییر کند:

```
v2:product:123
```

اگر نسخه قدیمی و جدید Application هم‌زمان فعال باشند، باید Migration Strategy مشخصی وجود داشته باشد:

```
Dual Read
Dual Write
Versioned Namespace
Backward-Compatible Payload
Gradual Migration
```

در غیر این صورت Versionها ممکن است Valueهای ناسازگار را روی Cache یکدیگر بنویسند یا Cache موجود را غیرقابل استفاده کنند.

## Source داده‌های Warming

داده‌های موردنیاز برای Warming می‌توانند از منابع مختلف دریافت شوند:

```
Primary Database
Read Replica
Analytics Store
Previous Cache Snapshot
Object Storage Snapshot
Event Log
Precomputed Export
```

Primary Database همیشه بهترین Source نیست. در صورت وجود Read Replica یا Snapshot معتبر، می‌توان فشار را از Primary برداشت. بااین‌حال، Staleness، Replication Lag و Consistency باید بررسی شوند.

## Snapshot-Based Warming

برای Datasetهای نسبتاً ثابت می‌توان Snapshot ساخت:

```
Nightly Snapshot → Store in Object Storage
→ Load into Cache at Startup → Apply Recent Changes
```

مزایا:

```
Fast Bulk Load
Low Primary DB Pressure
Predictable Recovery
```

معایب:

```
Snapshot Staleness
Complex Incremental Update
Storage and Version Management
```

Snapshot باید Version، زمان تولید، Schema و Compatibility مشخص داشته باشد. همچنین پس از Load باید تغییرات بین زمان Snapshot و زمان فعلی اعمال شوند.

## Event-Based Warming

Business Eventها می‌توانند Cache را پیش‌دستانه Populate یا Refresh کنند:

```
ProductCreated
ProductUpdated
CampaignStarted
MatchStarted

Flow:
Business Event → Consumer → Populate or Refresh Cache
```

این مدل Cache را پیش از Request آماده می‌کند، اما Reliability در Event Delivery، Ordering، Duplicate Event، Retry و Idempotency ضروری است.

Cache نباید تنها Consumer حقیقت باشد؛ در صورت ازدست‌رفتن Event باید مسیر Recovery یا Reconciliation وجود داشته باشد.

## Scheduled Warming

برخی Traffic Patternها قابل پیش‌بینی‌اند:

```
Flash Sale at 18:00
Football Match at 20:00
Daily Report at 08:00
Salary Portal at Month End
```

می‌توان پیش از شروع Peak، داده مرتبط را Warm کرد:

```
Known Peak Time → Warm Relevant Data
→ Validate Hit Ratio → Increase Capacity → Open Traffic
```

این نوع Proactive Warming زمانی مؤثر است که برنامه رویداد و Working Set از قبل قابل پیش‌بینی باشد.

## User-Specific Warming

Warm کردن داده تمام Userها معمولاً پرهزینه است، اما برای Userهای فعال ممکن است مفید باشد. برای مثال پس از Login:

```
Login Success → Load Profile → Permissions
→ Preferences → Recent Activity → Populate Cache
```

این کار می‌تواند Latency Requestهای بعدی را کاهش دهد. بااین‌حال، Login نباید با تعداد زیادی Query اضافی Block شود. داده ضروری می‌تواند Synchronous و داده اختیاری به‌صورت Background Warm شود.

## Consistency در Cache Warming

ممکن است Warmer داده‌ای را بخواند که هم‌زمان Update می‌شود:

```
Warmer Reads Old Value → Update Happens
→ Warmer Writes Old Value to Cache
```

در این حالت Warming می‌تواند مقدار تازه Cache را با مقدار قدیمی جایگزین کند.

راهکارهای معمول:

```
Version Check
UpdatedAt Comparison
Write-Through on Update
Event Ordering
Conditional Cache Write
Short TTL
```

می‌توان Version را داخل Value ذخیره کرد:

```
{
  "version": 42,
  "data": {}
}
```

Warmer نباید Version قدیمی‌تر را روی Version جدیدتر بنویسد. در Redis می‌توان این رفتار را با Script یا Transaction کنترل کرد و در Application نیز قبل از Write، Versionها را مقایسه کرد.

## Failure Handling

Warming باید برای خطاهای مختلف رفتار مشخصی داشته باشد:

```
DB Timeout
Redis Unavailable
Partial Batch Failure
Serialization Error
Context Cancellation
Rate Limit Wait Timeout
Deployment Shutdown
```

اصول مناسب:

```
Retry Only Retryable Errors
Use Backoff and Jitter
Limit Retry Count
Record Failed Keys
Continue Independent Items
Do Not Block Forever
```

یک Key خراب نباید کل Job را متوقف کند، مگر آنکه جزو Critical Warm Set باشد.

## Partial Success

فرض کنید از ده هزار Key، تعداد ۹٬۹۵۰ مورد موفق و ۵۰ مورد ناموفق باشند.

رفتار نامناسب:

```
50 Failures → Entire Warm Job Failed → Discard Success
```

رفتار مناسب‌تر:

```
Keep 9,950 Successful Entries
→ Record 50 Failures
→ Classify Errors
→ Retry Failed Keys Separately
```

برای Critical Set ممکن است سیاست سخت‌گیرانه‌تری اعمال شود؛ برای مثال Failure یک Configuration حیاتی می‌تواند مانع Readiness شود، ولی Failure چند Product کم‌مصرف نباید Deployment را متوقف کند.

## Retry و Backoff

Retry فوری می‌تواند Source را بیشتر تحت فشار قرار دهد:

```
Failure → Immediate Retry → More Failure → Retry Storm
```

روش مناسب:

```
Retryable Failure → Exponential Backoff
→ Jitter → Maximum Attempts → Dead-Letter or Report
```

Retry باید Retry Budget داشته باشد و Errorهای Permanent مانند Validation یا Serialization ناسازگار نباید بی‌دلیل Retry شوند.

## Cancellation و Graceful Shutdown

Warming در Go باید از `context.Context` پشتیبانی کند:

```
Deployment Shutdown → Cancel Context → Stop New Work
→ Finish or Abort In-Flight Work → Exit Cleanly
```

پس از Cancellation نباید Workerهای جدید ایجاد شوند. عملیات Source و Cache نیز باید Context را دریافت کنند تا Shutdown برای مدت نامحدود منتظر Timeoutهای طولانی نماند.

## مدل پیاده‌سازی در Go

Flow کلی:

```
Load Candidate Keys
→ Sort by Priority
→ Split into Batches
→ Apply Rate Limit
→ Dispatch to Bounded Workers
→ Load from Source
→ Validate and Serialize
→ Write Cache with Jittered TTL
→ Record Metrics
→ Retry Failed Items
```

برای Production بهتر است از Worker Pool ثابت و Queue محدود استفاده شود، نه ساخت یک Goroutine برای هر Key.

یک ساختار ساده:

```
type WarmItem struct {
	Key      string
	Priority int
}

type Source interface {
	Load(ctx context.Context, key string) ([]byte, error)
}

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type Warmer struct {
	source  Source
	cache   Cache
	workers int
	queue   chan WarmItem
}
```

مدل Worker:

```
func (w *Warmer) worker(ctx context.Context, results chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return

		case item, ok := <-w.queue:
			if !ok {
				return
			}

			value, err := w.source.Load(ctx, item.Key)
			if err != nil {
				results <- fmt.Errorf("load warm item %q: %w", item.Key, err)
				continue
			}

			ttl := jitteredTTL(10*time.Minute, time.Minute)
			if err := w.cache.Set(ctx, item.Key, value, ttl); err != nil {
				results <- fmt.Errorf("cache warm item %q: %w", item.Key, err)
			}
		}
	}
}
```

مدل Dispatch:

```
func (w *Warmer) Warm(ctx context.Context, items []WarmItem) error {
	if w.workers <= 0 {
		return errors.New("worker count must be positive")
	}

	results := make(chan error, len(items))

	var wg sync.WaitGroup
	for range w.workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.worker(ctx, results)
		}()
	}

enqueue:
	for _, item := range items {
		select {
		case <-ctx.Done():
			break enqueue
		case w.queue <- item:
		}
	}

	close(w.queue)
	wg.Wait()
	close(results)

	var errs []error
	for err := range results {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if ctx.Err() != nil {
		errs = append(errs, ctx.Err())
	}

	return errors.Join(errs...)
}
```

این نمونه فقط چارچوب آموزشی است. در Production باید Queue Lifecycle، استفاده مجدد از Warmer، Batch Loading، Rate Limiting، Retry Policy، Error Classification، Metrics و Maximum Error Retention با دقت طراحی شوند. همچنین نگهداری تمام Errorها برای میلیون‌ها Key می‌تواند Memory زیادی مصرف کند؛ در این حالت باید Errorها Aggregate یا Sample شوند و Keyهای شکست‌خورده در Storage یا Queue جدا ثبت شوند.

## جلوگیری از Unbounded Goroutines

```
Anti-pattern:
for _, key := range keys {
	go warm(key)
}
```

اگر یک میلیون Key وجود داشته باشد، ممکن است یک میلیون Goroutine ساخته شود و Memory، Scheduler و Connection Pool تحت فشار قرار گیرند.

روش مناسب:

```
Fixed Worker Count
Bounded Queue
Backpressure
```

برای مثال:

```
Workers = 32
Queue Capacity = 1,000
```

وقتی Queue پر شود، Producer منتظر می‌ماند و Backpressure ایجاد می‌شود.

## Double-Check پیش از Load یا Write

ممکن است در حین Warming، یک Request واقعی Cache را برای همان Key پر کرده باشد:

```
Warm Candidate → Cache Already Filled? → Skip Source Load
```

Double-Check می‌تواند Duplicate Work را کاهش دهد، اما خودش یک Cache Read اضافه ایجاد می‌کند. استفاده از آن باید براساس هزینه Source و احتمال هم‌پوشانی تصمیم‌گیری شود.

گاهی بررسی پیش از Load کافی نیست، زیرا بین Check و Source Load ممکن است Request دیگری Cache را پر کند. در این حالت `singleflight` یا Version-aware Write مؤثرتر است.

## Cache Warming و singleflight

اگر Warming و Request واقعی هم‌زمان یک Key را Load کنند، `singleflight` می‌تواند Work را یکی کند:

```
Warmer + Real Request → Same Key → One Source Load
```

بهتر است Request Path و Warmer از یک Loader مشترک استفاده کنند تا Cache Recheck، Error Classification، TTL و Deduplication یکسان باشند.

بااین‌حال، `singleflight` جای Global Concurrency Limit را نمی‌گیرد؛ زیرا Warming ممکن است هزاران Key متفاوت داشته باشد و برای هر Key یک Load مستقل همچنان اجرا شود.

## Readiness Strategy

Readiness می‌تواند چند سطح داشته باشد:

```
Process Ready
Dependency Ready
Critical Cache Ready
Full Warm Ready
```

برای مثال Service زمانی Ready شود که:

```
Redis Connected
Database Healthy
Global Config Loaded
Top 1,000 Keys Warmed
```

اما برای Warm شدن باقی صد هزار Key منتظر نماند.

Readiness Condition باید محدود و قابل‌دستیابی باشد. اگر معیار به درصدی وابسته باشد که به‌دلیل چند Key خراب هیچ‌گاه کامل نمی‌شود، Deployment گیر خواهد کرد.

## Observability

Metricهای مهم Warming عبارت‌اند از:

```
cache_warm_started_total
cache_warm_completed_total
cache_warm_failed_total
cache_warm_duration_seconds
cache_warm_items_total
cache_warm_items_succeeded_total
cache_warm_items_failed_total
cache_warm_inflight
cache_warm_queue_depth
cache_warm_source_qps
cache_warm_cache_write_qps
cache_warm_retry_total
cache_warm_skipped_total
cache_hit_ratio
```

علاوه بر Metricهای خود Warmer، باید Dependencyها نیز مشاهده شوند:

```
DB CPU
DB Query Rate
DB Pool Wait
DB Query P95/P99
Redis CPU
Redis Network
Redis Memory
Redis Command Latency
Application Readiness Time
Request P95/P99 During Startup
```

Metricها باید میان Live Traffic و Warming Traffic قابل تفکیک باشند تا مشخص شود Load ایجادشده مربوط به کاربران است یا Warmer.

## Alertهای مهم

```
Warming Duration Exceeds Threshold
Warming Failure Rate Is Too High
Source QPS Exceeds Warming Budget
DB Pool Wait Increases During Warming
Cache Write Errors Increase
Warm Queue Stops Progressing
Hit Ratio Remains Low After Warming
Multiple Instances Start Duplicate Warming
Critical Warm Set Is Incomplete
Retry Rate Continues Increasing
```

فقط Alert روی Failure نهایی کافی نیست؛ توقف Progress نیز باید تشخیص داده شود. ممکن است Job Fail نشده باشد، اما Queue به‌دلیل Dependency یا Deadlock جلو نرود.

## Success Criteria

صرفاً پایان یافتن Job به معنای موفقیت Warming نیست. معیارهای مناسب می‌توانند شامل موارد زیر باشند:

```
Critical Keys Warmed ≥ 99%
Cache Hit Ratio Reaches Target
DB Load Remains Below Warming Budget
No Significant P99 Regression
No Excessive Redis Eviction
Traffic Ramp-Up Completes Safely
Duplicate Work Remains Below Threshold
```

Success Criteria باید قبل از Deployment مشخص شود، نه اینکه پس از Incident به‌صورت مبهم تعریف شود.

## Production Example: E-commerce Deployment

فرض کنید:

```
Products = 2 Million
Hot Products = 15K
Traffic = 60K Requests/sec
DB Capacity = 8K Queries/sec
New Deployment Uses v2 Cache Keys
```

اگر نسخه جدید مستقیماً Traffic کامل دریافت کند:

```
v2 Cache Empty → 60K Requests/sec Miss → DB Saturation
```

طراحی مناسب:

```
Build Top 15K Product List
→ Warm v2 Keys in Batches
→ Apply TTL Jitter
→ Verify Critical-Key Coverage
→ Route 10% Traffic
→ Observe Metrics
→ Increase Traffic Gradually
→ Load Long Tail Lazily
```

بهتر است Warming از Read Replica یا Snapshot انجام شود تا Primary Database تحت فشار کمتری قرار بگیرد.

## Production Example: Live Match

فرض کنید:

```
Match Starts at 20:00
Expected Traffic = 200K Requests/sec
Relevant Keys = Score, Lineup, Timeline, Stats
```

چند دقیقه پیش از مسابقه:

```
Warm Match Metadata
→ Warm Team Data
→ Warm Initial Score
→ Warm Lineup
→ Start Background Refresh
→ Enable CDN or L1 Cache
```

در ساعت شروع، Traffic به Cache گرم وارد می‌شود. برای داده‌هایی مانند Score که مرتب تغییر می‌کنند، Warming اولیه باید با Refresh Ahead یا Push-Based Update ادامه پیدا کند.

## Production Example: Multi-Tenant Configuration

فرض کنید یک SaaS ده هزار Tenant دارد، اما فقط پانصد Tenant در هر ساعت فعال‌اند.

روش ضعیف:

```
Warm Configuration for All 10K Tenants
```

روش بهتر:

```
Warm Recently Active 500 Tenants
→ Warm High-Value Tenants
→ Lazy Load Remaining Tenants
```

Candidateها می‌توانند از Login Activity، Recent Requests، Scheduled Jobs و Business Priority استخراج شوند.

## Review Scenario: Catalog Service

فرض کنید:

```
Catalog Items = 8 Million
Top 50K Items = 85% Traffic
Traffic = 100K Requests/sec
DB Safe Capacity = 12K Queries/sec
Redis Was Fully Restarted
20 Application Instances Start Together
```

### مشکل اصلی

Cache کاملاً Cold است. اگر Traffic کامل وارد شود، بخش بزرگی از 100K Request/sec به Database می‌رسد، در حالی که Database فقط 12K Query/sec ظرفیت امن دارد.

اگر هر 20 Instance مستقلاً پنجاه هزار Item را Warm کنند:

```
20 × 50K = 1 Million Duplicate Warm Loads
```

در نتیجه Warming Storm ایجاد می‌شود.

### آیا باید هر هشت میلیون Item Warm شود؟

خیر. چون پنجاه هزار Item حدود 85 درصد Traffic را پوشش می‌دهند:

```
Warm 50K Hot Items → Protect 85% of Traffic
→ Load 7.95M Long-Tail Items Lazily
```

### طراحی مناسب

```
Dedicated or Leader Warmer
→ Load Top 50K from Read Replica
→ Batch and Rate Limit
→ Use Bounded Workers
→ Pipeline Redis Writes
→ Apply TTL Jitter
→ Warm Critical Set First
→ Verify Metrics
→ Ramp Traffic Gradually
```

### تعداد Worker مناسب

عدد ثابتی برای همه سیستم‌ها وجود ندارد. می‌توان با مقدار محافظه‌کارانه‌ای مانند 20 Worker شروع کرد و Metricهای زیر را مشاهده کرد:

```
DB QPS
DB Pool Wait
Query P99
Redis Write Latency
Warming Throughput
Error Rate
```

سپس Concurrency به‌صورت کنترل‌شده تنظیم شود. Adaptive Concurrency نیز ممکن است براساس DB Pool Wait یا Latency، سرعت Warming را کاهش یا افزایش دهد.

### رفتار در برابر 500 Key ناموفق

نباید موفقیت 49,500 Key کنار گذاشته شود. Keyهای ناموفق باید ثبت، دسته‌بندی و در صورت Retryable بودن با Backoff دوباره پردازش شوند.

اگر Keyهای ناموفق جزو Critical Set باشند، Traffic Ramp-Up یا Readiness ممکن است متوقف شود. اگر Long-Tail باشند، Service می‌تواند ادامه دهد و آن‌ها را Lazily Load کند.

### جلوگیری از Avalanche آینده

همه Keyها نباید TTL یکسان داشته باشند:

```
Base TTL = 30m
Jitter = ±15%
```

Refresh Ahead نیز باید Bounded و Rate-Limited باشد تا Expiration هماهنگ فقط به Refresh هماهنگ تبدیل نشود.

### Metricهای حیاتی

```
Warm Completion Percentage
Critical-Key Success Rate
Warm Source QPS
DB Pool Wait
Redis Write Error Rate
Cache Hit Ratio After Warm
Traffic Ramp Percentage
Request P99 Latency
Duplicate Warm Load Count
Retry Count
```

## تفاوت Cache Warming با Refresh Ahead

```
Cache Warming
→ Populate Empty or Cold Cache
→ Usually During Startup, Deployment or Recovery

Refresh Ahead
→ Refresh Existing Cache Before Expiry
→ Usually Continuous Runtime Behavior
```

این دو مکمل یکدیگرند:

```
Warming Creates Initial Cache State
Refresh Ahead Keeps It Warm
```

اگر Warming انجام شود ولی Refresh Strategy وجود نداشته باشد، Cache پس از مدتی دوباره Cold یا Expired خواهد شد.

## تفاوت Cache Warming با Cache Aside

در Cache Aside:

```
Request Causes Miss → Application Loads Source
→ Application Writes Cache
```

در Cache Warming:

```
System Predicts Future Need
→ Loads Data Before the Request
```

Cache Warming می‌تواند از همان Loader مربوط به Cache Aside استفاده کند، با این تفاوت که Loader به‌صورت پیش‌دستانه فراخوانی می‌شود. استفاده از Loader مشترک باعث می‌شود Error Handling، `singleflight`، TTL و Serialization در دو مسیر متفاوت نشوند.

## چه زمانی Cache Warming مناسب نیست؟

Warming همیشه مفید نیست. در شرایط زیر ممکن است هزینه آن بیشتر از منفعتش باشد:

```
Access Pattern Is Highly Unpredictable
Dataset Is Extremely Large
Most Warmed Data Is Never Used
Source Load Is Very Expensive
Cache Memory Is Limited
Data Changes Very Frequently
Startup Must Be Extremely Fast
```

اگر Hit Ratio داده‌های Warm شده پایین باشد، Warming فقط Source Load، Cache Write و Memory اضافی ایجاد می‌کند.

برای تصمیم‌گیری باید Metricهایی مانند درصد استفاده از داده‌های Warm‌شده، Time-to-First-Use و Memory Waste بررسی شوند.

## Production Insight

Cache Warming یک حلقه ساده روی تمام Keyها نیست. Warmer خودش یک Client بزرگ برای Database و Redis است و اگر کنترل نشود، می‌تواند همان Failureای را ایجاد کند که قرار بود از آن جلوگیری کند.

```
Anti-pattern:
Cache Is Empty → Warm Everything as Fast as Possible
→ Database and Redis Saturate
```

طراحی صحیح:

```
Identify Valuable Data
→ Prioritize
→ Batch
→ Limit Concurrency and Rate
→ Apply TTL Jitter
→ Observe Dependencies
→ Ramp Traffic Gradually
```

اصل مهم:

```
Warming Must Consume Spare Capacity,
Not Compete With Critical Traffic
```

## جمع‌بندی و Mental Model نهایی

Cache Warming یعنی آماده کردن داده‌های مهم پیش از رسیدن Traffic واقعی:

```
Cold Cache → Preload Critical and Hot Data
→ Controlled Traffic Entry → Higher Initial Hit Ratio
→ Lower Source Pressure
```

راهکارهای اصلی:

```
Hot-Key Prioritization
Minimum Viable Warm Set
Hybrid Eager and Lazy Loading
Bounded Worker Pool
Concurrency and Rate Limiting
Batch Source Reads
Pipelined Cache Writes
TTL Jitter
Leader or Dedicated Warmer
Gradual Traffic Ramp-Up
Retry with Backoff
Partial Failure Handling
Version-Aware Writes
Graceful Cancellation
Observability
```

مدل ذهنی نهایی این است: **Cache Warming به معنای پر کردن کامل Cache نیست؛ بلکه به معنای استفاده کنترل‌شده از ظرفیت آزاد سیستم برای آماده کردن ارزشمندترین داده‌ها پیش از آن است که Traffic واقعی مجبور شود هزینه Cold Cache را بپردازد.**
