# TTL Strategy

TTL یا `Time To Live` مشخص می‌کند یک Cache Entry چه مدت معتبر باقی بماند. پس از پایان این زمان، Entry منقضی می‌شود و Cache دیگر نباید آن را به‌عنوان داده معتبر عادی برگرداند.

```
Write Cache Entry → TTL Countdown Starts → TTL Reaches Zero
→ Entry Expires → Next Request Becomes Cache Miss
```

برای مثال، اگر `product:123` ساعت 10:00 با TTL پنج‌دقیقه‌ای در Cache نوشته شود، تقریباً تا ساعت 10:05 قابل استفاده است. پس از آن، Request بعدی باید داده را دوباره از Source اصلی Load کند یا از یک Refresh Strategy استفاده شود.

انتخاب TTL فقط پاسخ به این سؤال نیست که «داده چند دقیقه داخل Cache بماند». TTL مستقیماً بر Data Freshness، Cache Hit Ratio، Database Load، Tail Latency، Memory Consumption، Invalidation Complexity، Availability و Failure Recovery اثر می‌گذارد. بنابراین TTL یک عدد تصادفی یا Configuration ساده نیست، بلکه یک تصمیم معماری و Business-oriented است.

## Mental Model

Cache را مانند نسخه چاپ‌شده برنامه حرکت قطارها در نظر بگیرید. اگر این نسخه هر یک دقیقه تعویض شود، اطلاعات بسیار تازه است، اما هزینه تولید و توزیع آن بالا می‌رود. اگر فقط هفته‌ای یک بار تعویض شود، هزینه کاهش می‌یابد، اما احتمال نمایش اطلاعات قدیمی افزایش پیدا می‌کند.

```
Short TTL → Frequent Refresh → Fresher Data → Higher Source Load
Long TTL → Less Refresh → Higher Hit Ratio → More Stale Data
```

هیچ TTLای ذاتاً درست یا غلط نیست. TTL مناسب به نرخ تغییر داده، حداکثر Staleness قابل قبول، هزینه Load مجدد، اهمیت Business، ظرفیت Source و کیفیت Invalidation Strategy بستگی دارد.

## تعادل میان Freshness و Performance

فرض کنید Product Price در Cache قرار دارد. اگر TTL پنج ثانیه باشد، قیمت معمولاً تازه است، اما Expiration و Source Load بیشتر می‌شوند:

```
TTL = 5 Seconds
→ Better Freshness
→ Frequent Expiration
→ More Cache Misses
→ More Database Queries
```

اگر TTL یک ساعت باشد، Hit Ratio افزایش و Database Load کاهش پیدا می‌کند، اما ممکن است کاربر قیمت قدیمی ببیند:

```
TTL = 1 Hour
→ Higher Hit Ratio
→ Lower Source Load
→ Higher Staleness Risk
```

رابطه کلی:

```
TTL Shorter
→ Freshness ↑
→ Cache Hit Ratio ↓
→ Source Load ↑

TTL Longer
→ Freshness ↓
→ Cache Hit Ratio ↑
→ Source Load ↓
```

این رابطه مطلق نیست. با Invalidation دقیق، Refresh Ahead، Event-driven Update یا Soft/Hard TTL می‌توان TTL بلندتری داشت و همچنان Freshness مناسبی حفظ کرد.

## TTL یکسان برای همه داده‌ها

استفاده از یک TTL عمومی برای تمام Cacheها یک Anti-pattern رایج است:

```
Default TTL = 5 Minutes
Applied to Everything
```

انواع داده رفتارهای متفاوتی دارند. Country List ممکن است ماه‌ها تغییر نکند، Product Description کم‌تغییر باشد، Feature Flag گاهی تغییر کند، Inventory هر ثانیه تغییر کند و Live Match Score چند بار در دقیقه Update شود. بنابراین TTL باید براساس ماهیت داده تعیین شود، نه صرفاً براساس راحتی Configuration.

### داده تقریباً ثابت

نمونه‌ها:

```
Country Codes
Currency Metadata
Static Reference Tables
Protocol Definitions
```

TTL این داده‌ها ممکن است چند ساعت یا چند روز باشد:

```
TTL = 6 Hours to 24 Hours
```

در صورت وجود Event-based Invalidation، TTL طولانی‌تر نیز ممکن است منطقی باشد.

### داده کم‌تغییر

نمونه‌ها:

```
Product Description
User Preferences
Tenant Configuration
Feature Configuration
```

TTL احتمالی:

```
TTL = 5 Minutes to 1 Hour
```

مدت دقیق به Freshness Requirement و کیفیت Invalidation بستگی دارد.

### داده پرتغییر

نمونه‌ها:

```
Inventory
Exchange Rate
Live Score
Dynamic Pricing
Account Balance
```

TTL ممکن است بین یک ثانیه تا یک دقیقه باشد. بااین‌حال، داده‌های حساس نباید فقط به Cache متکی باشند. برای مثال Inventory در صفحه Product می‌تواند Cache شود، اما هنگام ثبت سفارش باید موجودی نهایی از Source معتبر و با عملیات Atomic بررسی شود.

## Business Staleness Budget

به‌جای پرسیدن «TTL مناسب چقدر است؟» بهتر است بپرسیم:

حداکثر قدیمی بودن قابل قبول این داده چقدر است؟

این مقدار را می‌توان `Staleness Budget` نامید:

```
Live Match Score → Up to 1 Second Stale
Product Description → Up to 30 Minutes Stale
Exchange Rate Display → Up to 1 Minute Stale
Payment Authorization Data → Possibly No Staleness Allowed
```

TTL باید با این Budget سازگار باشد، اما TTL لزوماً برابر Maximum Guaranteed Staleness نیست:

```
TTL ≠ Maximum Guaranteed Staleness
```

Source ممکن است بلافاصله پس از Cache Write تغییر کند. همچنین Replica Lag، Refresh Delay، Snapshot Age، Clock Skew و تأخیر Eventها روی Staleness واقعی اثر می‌گذارند.

## Fixed TTL

در ساده‌ترین Strategy، برای هر نوع داده TTL ثابتی تعریف می‌شود:

```
Product Cache TTL = 10 Minutes
User Profile TTL = 5 Minutes
Country List TTL = 24 Hours
```

مزایا:

```
Simple
Predictable
Easy to Configure
Easy to Observe
```

معایب:

```
Ignores Data Hotness
Ignores Per-Entry Update Frequency
May Synchronize Expiration
May Be Too Long for Some Entries
May Be Too Short for Others
```

Fixed TTL نقطه شروع مناسبی است، اما در سیستم‌های بزرگ معمولاً به‌تنهایی کافی نیست.

## TTL Jitter

اگر تعداد زیادی Key با TTL یکسان و در بازه زمانی نزدیک نوشته شوند، ممکن است هم‌زمان Expire شوند و Cache Avalanche ایجاد کنند.

بدون Jitter:

```
100K Keys Written at 10:00
TTL = 30 Minutes
→ Most Keys Expire around 10:30
```

با Jitter:

```
Base TTL = 30 Minutes
Jitter = ±5 Minutes
```

نمونه:

```
Key A = 26m20s
Key B = 31m10s
Key C = 34m05s
Key D = 28m45s
```

نتیجه:

```
Synchronized Expiration
→ Distributed Expiration
→ Smoother Source Load
```

### پیاده‌سازی TTL Jitter در Go

در نسخه‌های جدید Go می‌توان از `math/rand/v2` استفاده کرد:

```
package cachettl

import (
	"math/rand/v2"
	"time"
)

func WithJitter(base, jitter time.Duration) time.Duration {
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

نمونه استفاده:

```
ttl := WithJitter(30*time.Minute, 5*time.Minute)
```

برای TTL Jitter به Random رمزنگاری‌شده نیاز نیست؛ هدف فقط پخش‌کردن زمان Expiration است.

### Percentage-Based Jitter

برای TTLهای مختلف، Jitter درصدی معمولاً مناسب‌تر است:

```
Base TTL = 1 Hour
Jitter = ±10%
Actual TTL = 54 to 66 Minutes

func WithPercentJitter(base time.Duration, percent float64) time.Duration {
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

Jitter نباید Freshness Requirement را نقض کند. اگر حداکثر Staleness مجاز ده دقیقه باشد، `TTL = 10m ± 5m` ممکن است بعضی Entryها را تا پانزده دقیقه نگه دارد. در چنین شرایطی می‌توان Jitter را فقط در جهت کاهش اعمال کرد:

```
Maximum TTL = 10 Minutes
Jitter = 2 Minutes
Actual TTL = 8 to 10 Minutes
```

در این مدل سقف Staleness افزایش پیدا نمی‌کند.

## Absolute Expiration

در Absolute Expiration، Entry در زمان مشخصی منقضی می‌شود و Read شدن آن زمان Expiration را تغییر نمی‌دهد:

```
Written at 10:00
TTL = 10 Minutes
→ Expires at 10:10
```

این مدل برای داده‌هایی مناسب است که باید دوره‌ای دوباره تأیید شوند:

```
Product Price
Exchange Rates
Configuration Snapshot
News Feed
```

مزایا:

```
Predictable Refresh
Staleness Bounded by Time
Simple Semantics
```

معایب:

```
Hot Key May Expire Under Heavy Traffic
Can Trigger Stampede
Popular Data Is Not Kept Automatically
```

برای Hot Keyها بهتر است Absolute Expiration با Refresh Ahead، Serve Stale، Soft/Hard TTL یا `singleflight` ترکیب شود.

## Sliding Expiration

در Sliding Expiration، هر Access موفق زمان Expiration را تمدید می‌کند:

```
TTL = 30 Minutes
Every Read Resets TTL to 30 Minutes
```

مثال:

```
Read at 10:10 → New Expiry = 10:40
Read at 10:35 → New Expiry = 11:05
```

این Strategy برای داده‌هایی مناسب است که تا زمان فعال بودن باید باقی بمانند:

```
User Sessions
Temporary Workflow State
Active Shopping Carts
Recently Used Objects
```

مزایا:

```
Frequently Used Data Stays Cached
Inactive Data Expires Naturally
```

معایب:

```
Hot Data May Never Expire
Stale Data Can Survive Indefinitely
Every Read May Require a Write or Touch
Replication and Persistence Load Increase
```

Sliding TTL برای داده Mutable مانند Product Price خطرناک است؛ زیرا Product محبوب ممکن است برای همیشه مقدار قدیمی را نگه دارد.

## ترکیب Sliding و Absolute Expiration

برای جلوگیری از زنده ماندن نامحدود Entry می‌توان Idle Timeout و Maximum Lifetime را ترکیب کرد:

```
Idle TTL = 30 Minutes
Maximum Lifetime = 24 Hours
```

Entry در یکی از این دو حالت منقضی می‌شود:

```
No Access for 30 Minutes
OR
24 Hours Since Creation
```

این مدل برای Session و Security-sensitive State مناسب است.

## Soft TTL و Hard TTL

در مدل ساده TTL، Entry فقط Valid یا Expired است:

```
Before Expiry → Use Value
After Expiry → Cache Miss
```

در مدل Soft/Hard TTL دو مرز داریم:

```
Soft TTL → Data Becomes Stale but May Still Be Used
Hard TTL → Data Must No Longer Be Served Normally
```

مثال:

```
Soft TTL = 5 Minutes
Hard TTL = 30 Minutes
```

رفتار:

```
Age < 5m
→ Serve Fresh

5m ≤ Age < 30m
→ Serve Stale
→ Refresh in Background

Age ≥ 30m
→ Reload Source or Apply Failure Policy

Flow:
Fresh Window
→ Stale-While-Revalidate Window
→ Hard Expiration
```

این Strategy Latency را کاهش می‌دهد و خطر Stampede را کمتر می‌کند.

### مدل Go برای Soft و Hard TTL

```
type CacheFreshness uint8

const (
	CacheFresh CacheFreshness = iota + 1
	CacheStale
	CacheExpired
)

type CachedValue[V any] struct {
	Value     V
	CreatedAt time.Time
	SoftTTL   time.Duration
	HardTTL   time.Duration
}

func (c CachedValue[V]) State(now time.Time) CacheFreshness {
	age := now.Sub(c.CreatedAt)

	switch {
	case age < c.SoftTTL:
		return CacheFresh
	case age < c.HardTTL:
		return CacheStale
	default:
		return CacheExpired
	}
}
```

در حالت `CacheStale` می‌توان Value را فوراً برگرداند و Refresh را در Background انجام داد. Refresh باید Per-Key Deduplicated، Rate-Limited و دارای Concurrency Limit باشد.

## Serve Stale

Serve Stale یعنی در شرایط مشخص، به‌جای ایجاد Cache Miss، مقدار قدیمی ارائه شود:

```
Database Temporarily Slow
Redis Recovery
Source Timeout
Traffic Spike
Cache Refresh Failure
```

مثال:

```
Product Description Is 12 Minutes Old
Hard TTL = 1 Hour
Database Is Down
→ Serve Stale Description
```

این روش برای همه داده‌ها مناسب نیست:

```
Account Balance
Authorization State
Inventory Reservation
One-Time Token
Payment Status
```

برای این داده‌ها، پاسخ قدیمی ممکن است از Failure صریح خطرناک‌تر باشد. Serve Stale باید براساس Data Classification و Business Risk تعریف شود.

## Stale-If-Error

در این Strategy، مقدار Stale فقط در صورت خطای Source ارائه می‌شود:

```
Cache Entry Soft-Expired
→ Try Refresh
→ Refresh Successful: Return Fresh
→ Refresh Failed: Return Stale
```

Error Classification اهمیت زیادی دارد:

```
Timeout → Maybe Serve Stale
Connection Failure → Maybe Serve Stale
Not Found → Do Not Blindly Serve Old Value
Permission Error → Do Not Serve Stale
Validation Error → Do Not Serve Stale
```

برای مثال، اگر Source بگوید Resource حذف شده است، ارائه نسخه قدیمی آن ممکن است اشتباه یا حتی ناامن باشد.

## Dynamic TTL

در Dynamic TTL، TTL براساس ویژگی‌های هر Entry تعیین می‌شود:

```
Update Frequency
Traffic Frequency
Business Criticality
Data Size
Source Load Cost
Staleness Tolerance
Current System Load
Confidence in Data
```

مثال:

```
Popular Product → TTL = 20 Minutes
Cold Product → TTL = 5 Minutes
Static Product Metadata → TTL = 1 Hour
Inventory Value → TTL = 5 Seconds
```

Dynamic TTL می‌تواند Cache Efficiency را افزایش دهد، اما Complexity، Debugging و Incident Analysis را سخت‌تر می‌کند.

## TTL براساس Update Frequency

داده‌ای که مرتب تغییر می‌کند، معمولاً به TTL کوتاه‌تری نیاز دارد:

```
Product A Changes Once per Month → TTL = 6 Hours
Product B Changes Every Minute → TTL = 20 Seconds
```

این اطلاعات می‌تواند از `updated_at`، Event Frequency یا Historical Metrics استخراج شود:

```
High Update Frequency → Short TTL
Low Update Frequency → Long TTL
```

TTL نباید فقط براساس فاصله Updateها تعیین شود؛ Business Staleness، Source Cost و Invalidation نیز باید در نظر گرفته شوند.

## TTL براساس Hotness

TTL کوتاه روی Hot Key می‌تواند خطر Stampede را افزایش دهد:

```
Hot Key + Frequent Expiration → High Source Load
```

TTL بلند نیز Staleness را افزایش می‌دهد. راهکار مناسب معمولاً ترکیبی است:

```
Longer Base TTL
+ Refresh Ahead
+ Soft/Hard TTL
+ L1 Cache
+ singleflight
+ TTL Jitter
```

برای Live Match Score، TTL بلند مناسب نیست و بهتر است Push-based Update یا L1 Cache کوتاه‌مدت استفاده شود. در مقابل، News Article محبوب و Immutable می‌تواند TTL بسیار بلندتری داشته باشد.

## TTL براساس Source Cost

اگر بازسازی Entry گران باشد، TTL بلندتر ممکن است منطقی باشد:

```
Simple Row Lookup = 2ms
Complex Aggregation = 3 Seconds
```

برای Aggregation گران، ترکیب زیر مناسب‌تر است:

```
Longer TTL
Refresh Ahead
Serve Stale
Background Recalculation
```

اما هزینه بالا نباید بهانه‌ای برای ارائه نامحدود داده قدیمی باشد؛ Freshness Limit همچنان باید تعریف شود.

## TTL براساس Data Size

Entry بزرگ‌تر Memory بیشتری مصرف می‌کند:

```
Key A = 2KB, TTL = 1 Hour
Key B = 20MB, TTL = 1 Hour
```

اثر این دو Entry بر Cache Capacity یکسان نیست. برای Valueهای بزرگ ممکن است این روش‌ها مناسب باشند:

```
Shorter TTL
Compression
Split Payload
Store Only Frequently Used Fields
Admission Control
Dedicated Cache
```

TTL Strategy باید Memory Economics را نیز در نظر بگیرد.

## Probabilistic Early Refresh

در این روش، هرچه Entry به Expiration نزدیک‌تر شود، احتمال Refresh توسط Request بیشتر می‌شود:

```
Far from Expiry → Almost No Refresh
Near Expiry → Higher Refresh Probability
```

هدف این است که Refreshها در زمان پخش شوند و پیش از Expiration کامل، یکی از Requestها داده را تازه کند.

مزایا:

```
Lower Stampede Risk
No Fixed Background Scheduler Required
Distributed Refresh Work
```

معایب:

```
More Complex
May Cause Extra Refreshes
Harder to Debug
Requires Deduplication
```

این روش باید با `singleflight` یا Distributed Lock ترکیب شود تا چند Request هم‌زمان Refresh نکنند.

## Refresh Ahead و TTL Strategy

در Refresh Ahead، Refresh Window پیش از پایان TTL آغاز می‌شود:

```
TTL = 10 Minutes
Refresh Window Starts at 8 Minutes

Flow:
Age < 8m → Serve Normally
Age ≥ 8m → Serve Current Value + Trigger Background Refresh
```

این مدل برای Hot Keyها مؤثر است، اما اگر همه Entryها هم‌زمان وارد Refresh Window شوند، خود Refresh Ahead می‌تواند Avalanche ایجاد کند:

```
100K Keys Enter Refresh Window Together
→ 100K Background Refreshes
```

دفاع مناسب:

```
TTL Jitter
Refresh Jitter
Bounded Worker Pool
Rate Limit
Per-Key Deduplication
```

## TTL در Negative Caching

Negative TTL معمولاً باید کوتاه‌تر از Positive TTL باشد:

```
Positive Product TTL = 10 Minutes
Negative Product TTL = 30 Seconds
```

TTL خیلی کوتاه باعث می‌شود Requestهای تکراری زود دوباره به Database برسند. TTL خیلی بلند ممکن است پس از ایجاد داده، False `404` ایجاد کند.

نمونه‌ها:

```
Deleted Legacy Article
→ Negative TTL Could Be Several Minutes

Username Availability
→ Negative TTL Must Be Very Short

Random Invalid Product ID
→ Short TTL + Rate Limit + Bloom Filter
```

Negative TTL نیز باید Jitter داشته باشد تا Entryهای منفی هم‌زمان Expire نشوند.

## TTL برای Session

Session TTL علاوه بر Performance، بخشی از Security Policy است. مدل رایج:

```
Idle Timeout = 30 Minutes
Absolute Lifetime = 12 Hours
```

اگر User سی دقیقه فعالیت نکند، Session منقضی می‌شود. حتی اگر پیوسته فعال باشد، پس از دوازده ساعت باید دوباره Authentication انجام شود.

## TTL برای Authentication و Authorization

Cache کردن Permission، Role یا Token State حساس است. اگر Permission لغو شود، اما TTL یک ساعت باشد:

```
Permission Revoked
→ Cache Still Allows Access for Up to 1 Hour
```

راهکار مناسب:

```
Short TTL
Event-Driven Invalidation
Token Version
Revocation List
Security-Sensitive Recheck
```

TTL بلند بدون Invalidation دقیق ممکن است Security Risk ایجاد کند.

## TTL برای Configuration

Configuration کم‌تغییر می‌تواند TTL بلند داشته باشد، اما تغییر آن باید سریع منتشر شود:

```
Long TTL
+ Event or Pub/Sub Invalidation
+ Periodic Safety Refresh
```

مثال:

```
TTL = 1 Hour
Config Change Event → Immediate Invalidation
```

در این مدل، TTL نقش Safety Net را دارد. اگر Event از دست برود، حداکثر پس از یک ساعت داده Refresh می‌شود.

## TTL برای داده مالی و حساس

برای Account Balance، Settlement Status یا Payment Authorization، Cache ممکن است برای نمایش یا Optimization مناسب باشد، اما تصمیم نهایی باید از Source معتبر گرفته شود:

```
Displayed Balance → May Use Short Cache
Withdrawal Decision → Must Check Authoritative Source
```

اصل مهم:

```
Cache Freshness Must Match Decision Criticality
```

هرچه تصمیم حساس‌تر و غیرقابل‌جبران‌تر باشد، اتکا به Cache قدیمی باید کمتر شود.

## No Expiration یا TTL نامحدود

ذخیره Entry بدون TTL ممکن است برای داده Immutable مناسب باشد، اما ریسک‌های مهمی دارد:

```
Stale Data Lives Forever
Memory Is Not Reclaimed Automatically
Invalidation Bug Becomes Permanent
Schema Changes Leave Orphaned Keys
Old Versions Accumulate
```

No-Expiration فقط زمانی منطقی است که:

```
Data Is Truly Immutable
Key Is Content-Addressed or Versioned
Memory Lifecycle Is Controlled
Deletion Strategy Exists
```

مثال:

```
blob:sha256:<content-hash>
```

چون Content Hash تغییر نمی‌کند، Value Immutable است. بااین‌حال، Garbage Collection همچنان لازم است. برای داده Mutable، No-Expiration معمولاً Anti-pattern است.

## Explicit Invalidation و TTL

دو روش اصلی کنترل Freshness:

```
Expiration-Based → Wait Until TTL Ends
Invalidation-Based → Source Changes and Cache Is Deleted or Updated
```

Invalidation دقیق اجازه TTL بلندتر را می‌دهد، اما ممکن است Event از دست برود، Consumer Down باشد یا عملیات Cache Delete شکست بخورد. بنابراین طراحی مناسب معمولاً چنین است:

```
Explicit Invalidation
+ Finite TTL as Safety Net
```

## TTL و Cache Key Versioning

هنگام تغییر Schema یا Serialization می‌توان Key را Version کرد:

```
v1:product:123
v2:product:123
```

نسخه قدیمی با TTL خود به‌تدریج حذف می‌شود:

```
Write New Version
→ Stop Refreshing Old Version
→ Let Old TTL Expire
```

اگر TTL نامحدود یا بسیار بلند باشد، Keyهای قدیمی مدت زیادی Memory مصرف می‌کنند و ممکن است Cleanup Job لازم شود.

## TTL و Memory Pressure

TTL فقط Freshness را کنترل نمی‌کند؛ روی Memory Consumption نیز اثر دارد. تقریب ساده:

```
Cache Memory
≈ Entry Arrival Rate × Average Lifetime × Average Entry Size
```

مثال:

```
100K New Entries per Minute
Average Size = 2KB
TTL = 60 Minutes
```

افزایش TTL تعداد Entryهای Resident را بیشتر می‌کند. اگر Redis دائماً به `maxmemory` برسد:

```
Long TTL
→ More Resident Keys
→ More Evictions
→ Lower Effective Hit Ratio
```

بنابراین TTL بلند همیشه Hit Ratio را بهتر نمی‌کند؛ ممکن است Memory Pressure باعث Eviction داده‌های مهم شود.

## Configured TTL و Effective TTL

TTL تنظیم‌شده لزوماً مدت واقعی ماندن Entry نیست:

```
Configured TTL = 1 Hour
Entry Evicted after 3 Minutes Due to Memory Pressure
```

پس:

```
Configured TTL ≠ Effective Cache Lifetime
```

Metricهای مهم:

```
Expired Keys
Evicted Keys
Average Entry Age at Eviction
Memory Usage
Hit Ratio
```

اگر Eviction زیاد باشد، افزایش TTL احتمالاً وضعیت را بدتر می‌کند.

## TTL در Multi-Level Cache

در Cache چندلایه، هر Layer می‌تواند TTL متفاوتی داشته باشد:

```
L1 Local Cache TTL = 1 Second
L2 Redis TTL = 10 Minutes
Database = Source of Truth

Flow:
Request → L1 → L2 → Database
```

نمونه Live Score:

```
L1 TTL = 500ms
L2 TTL = 5s
```

نمونه Product Description:

```
L1 TTL = 30s
L2 TTL = 30m
```

L1 TTL معمولاً کوتاه‌تر است تا Staleness میان Instanceها محدود شود. باید مراقب Staleness تجمعی بود؛ ممکن است L2 قدیمی باشد و L1 همان مقدار قدیمی را دوباره Cache کند. نگهداری Metadataهایی مانند `generated_at`، `updated_at` یا `version` کمک می‌کند Age واقعی مشخص باشد.

## TTL از زمان تولید داده یا Cache Write

در حالت معمول TTL از زمان Cache Write محاسبه می‌شود، اما اگر داده از Snapshot قدیمی آمده باشد:

```
Snapshot Generated 30 Minutes Ago
Cache Write Happens Now
TTL = 1 Hour
```

داده هنگام Expiration ممکن است عملاً 90 دقیقه قدیمی باشد. بهتر است این Timestampها همراه Value نگهداری شوند:

```
SourceUpdatedAt
GeneratedAt
CachedAt
```

سپس Freshness براساس Age واقعی داده ارزیابی شود.

## Clock و TTL

Redis Expiration را براساس Clock Cache Server مدیریت می‌کند، اما Application-level Soft TTL از Clock Application استفاده می‌کند. مشکلات احتمالی:

```
Clock Skew
Clock Adjustment
Incorrect Timestamps
Cross-Region Differences
```

در Go، `time.Time` درون همان Process می‌تواند Monotonic Component داشته باشد، اما پس از Serialize و Deserialize شدن معمولاً این بخش از بین می‌رود. در سیستم توزیع‌شده باید Clock Synchronization مناسب، Margin کافی و Semantics مشخص برای `updated_at` و `generated_at` وجود داشته باشد.

## TTL Configuration و Guardrail

TTLها بهتر است Configurable باشند، اما Configuration بدون Validation خطرناک است:

```
Product TTL: 10 Minutes → 0
```

در برخی Libraryها صفر به معنای No Expiration و در برخی دیگر به معنای Expire Immediately است. Semantics باید صریح باشد:

```
TTL > 0 → Expiring Entry
TTL = 0 → Explicitly Defined Behavior
TTL < 0 → Invalid Configuration
```

برای No-Expiration بهتر است Option یا Constant مشخص وجود داشته باشد، نه عدد جادویی صفر.

### Validation در Go

```
type TTLPolicy struct {
	PositiveTTL time.Duration
	NegativeTTL time.Duration
	Jitter      time.Duration
}

func (p TTLPolicy) Validate() error {
	if p.PositiveTTL <= 0 {
		return errors.New("positive TTL must be greater than zero")
	}
	if p.NegativeTTL <= 0 {
		return errors.New("negative TTL must be greater than zero")
	}
	if p.Jitter < 0 {
		return errors.New("TTL jitter must not be negative")
	}
	if p.Jitter >= p.PositiveTTL {
		return errors.New("TTL jitter must be smaller than positive TTL")
	}

	return nil
}
```

Validation باید هنگام Startup انجام شود تا Configuration اشتباه سریع و واضح Fail شود.

## Centralized TTL Policy

پراکنده‌کردن Durationها در Code یک Anti-pattern است:

```
cache.Set(ctx, key, value, 5*time.Minute)
```

اگر این مقدار در چندین فایل تکرار شود، تغییر، Audit و تست آن سخت خواهد شد.

مدل ساختاریافته:

```
type CacheTTLs struct {
	Product       time.Duration
	UserProfile   time.Duration
	Configuration time.Duration
	Negative      time.Duration
}
```

یا Policy پویا:

```
type CacheTTLPolicy interface {
	ProductTTL(product Product) time.Duration
	UserTTL(user User) time.Duration
	NegativeProductTTL() time.Duration
}
```

Policy باید نام‌های معنایی داشته باشد و فقط به `DefaultTTL` متکی نباشد.

## Adaptive TTL

در Strategy پیشرفته، TTL براساس Runtime Metrics تغییر می‌کند:

```
Database Under High Load
→ Increase TTL for Stale-Tolerant Data

Data Update Rate Increases
→ Decrease TTL

Data Remains Stable
→ Gradually Increase TTL
```

مزایا:

```
Better Resource Utilization
Adaptation to Real Workload
```

خطرها:

```
Feedback Loop Instability
Harder Debugging
Unexpected Staleness
Complex Incident Analysis
```

Adaptive TTL باید Bound داشته باشد:

```
Min TTL ≤ Calculated TTL ≤ Max TTL
```

تغییرات باید تدریجی، Observable و قابل Rollback باشند.

## Load-Aware TTL

در Degraded Mode می‌توان TTL داده‌های Stale-tolerant را افزایش داد:

```
Normal TTL = 5 Minutes
Degraded Mode TTL = 20 Minutes
```

این کار برای Recommendation یا Trending List ممکن است مناسب باشد، اما برای Payment Status یا Authorization مناسب نیست. این رفتار باید Explicit Degradation Policy باشد، نه تغییر مخفی و غیرقابل‌مشاهده.

## رابطه TTL با Cache Stampede

TTL Strategy نامناسب یکی از علت‌های Stampede است:

```
Hot Key
+ Hard Expiration
+ No singleflight
= Cache Stampede
```

راهکارها:

```
TTL Jitter
Soft/Hard TTL
Refresh Ahead
singleflight
Serve Stale
```

Hot Key نباید با یک Expiration ناگهانی هزاران Request را به Source منتقل کند.

## رابطه TTL با Cache Avalanche

```
Mass Write
+ Same TTL
= Synchronized Expiration
```

راهکارها:

```
Per-Key Jitter
Staggered Warming
Rate-Limited Refresh
Expiration Distribution Monitoring
```

TTL Strategy باید Distribution زمان Expiration را در نظر بگیرد، نه فقط Average TTL را.

## رابطه TTL با Hot Key

TTL کوتاه برای Hot Key Source Load و Stampede Risk را افزایش می‌دهد. TTL بلند نیز Staleness را بیشتر می‌کند. ترکیب مناسب:

```
L1 Short TTL
L2 Longer TTL
Refresh Ahead
Push-Based Invalidation
Serve Stale
```

هدف این است که Hot Key بدون Expiration Storm تازه باقی بماند.

## رابطه TTL با Cache Penetration

در Negative Cache، TTL باید از Source محافظت کند، اما False Not Found طولانی ایجاد نکند:

```
Repeated Missing Key → Negative TTL Is Effective
Random Unique Missing Keys → Negative TTL Alone Is Not Enough
```

برای Traffic تصادفی همچنان Validation، Bloom Filter، Rate Limiting و Source Protection لازم است.

## رابطه TTL با Cache Warming

اگر Warmer تمام Keyها را با TTL یکسان بنویسد:

```
Warming Now → Avalanche Later
```

Warmer باید از Central TTL Policy، TTL Jitter و Refresh Rate Limit استفاده کند.

## Observability

Metricهای مهم:

```
cache_entry_ttl_seconds
cache_entry_age_seconds
cache_expiration_total
cache_eviction_total
cache_refresh_total
cache_refresh_failure_total
cache_stale_served_total
cache_hard_expired_total
cache_negative_expiration_total
cache_hit_ratio
cache_miss_ratio
source_load_after_expiration
```

Distributionها نیز باید مشاهده شوند:

```
TTL Distribution
Entry Age Distribution
Expirations per Second
Refreshes per Second
Distinct Keys Expiring per Window
```

Average TTL به‌تنهایی کافی نیست. ممکن است Average مناسب باشد، اما تعداد زیادی Key در یک بازه کوتاه Expire شوند.

### Metric Cardinality

ثبت Key کامل به‌عنوان Label اشتباه است:

```
key="product:123"
```

روش‌های مناسب‌تر:

```
Key Type
Cache Namespace
TTL Bucket
Data Category
Sampled Top-K
```

مثال:

```
cache_entry_ttl_seconds{
    namespace="product",
    ttl_bucket="5m-15m"
}
```

## Alertهای مهم

```
Expiration Rate Spikes
Cache Miss Rate Rises After Expiration
Source QPS Spikes Periodically
Stale Served Rate Remains High
Refresh Failure Rate Increases
Eviction Rate High Despite Long TTL
Negative Cache Expiration Wave
Many Keys Share the Same Expiry Window
```

Pattern زیر معمولاً نشان‌دهنده TTL هماهنگ و نبود Jitter است:

```
Every 30 Minutes
→ Cache Miss Spike
→ DB QPS Spike
→ P99 Spike
```

## روش عملی انتخاب TTL

برای هر نوع داده باید این پرسش‌ها پاسخ داده شوند:

### داده هر چند وقت تغییر می‌کند؟

```
Seconds
Minutes
Hours
Days
Rarely
```

### حداکثر Staleness قابل قبول چقدر است؟

```
No Staleness
1 Second
30 Seconds
10 Minutes
1 Hour
```

### Load مجدد چقدر گران است؟

```
Simple Memory Lookup
Indexed DB Query
Remote API
Complex Aggregation
Machine Learning Inference
```

### آیا Invalidation قابل اعتماد وجود دارد؟

```
No Invalidation
Best-Effort Event
Reliable Event
Synchronous Update
```

### Traffic داده چقدر است؟

```
Cold
Normal
Hot
Extremely Hot
```

### هنگام Source Failure می‌توان Stale ارائه کرد؟

```
Always
Only for Some Errors
Only for Limited Time
Never
```

### Memory Cost چقدر است؟

```
Small
Medium
Large
Very Large
```

سپس Base TTL انتخاب، Jitter و Refresh Policy تعریف و نتیجه با Metricهای واقعی اصلاح می‌شود.

## ماتریس نمونه TTL

اعداد زیر فقط نمونه‌اند و قانون عمومی نیستند:

```
Country Reference Data
→ TTL: 12–24 Hours
→ Invalidate on Update

Product Description
→ TTL: 15–60 Minutes
→ Jitter + Invalidation

Product Price
→ TTL: 10–60 Seconds
→ Event Update + Short Safety TTL

Inventory Display
→ TTL: 1–5 Seconds
→ Final Check at Purchase

Live Match Score
→ L1: 250ms–1s
→ L2: 2–10s
→ Push or Refresh

Feature Flags
→ TTL: 30s–5m
→ Push Invalidation

User Profile
→ TTL: 5–30 Minutes
→ Invalidate on Update

Authorization Data
→ Very Short TTL
→ Explicit Invalidation Required

Negative Product Lookup
→ TTL: 15–60 Seconds
→ Jitter

Static Immutable Content
→ Long TTL or Versioned No-Expiry
```

## Production Example: Product Catalog

فرض کنید Product Description به‌ندرت تغییر می‌کند، Price چند بار در ساعت Update می‌شود، Inventory هر چند ثانیه تغییر می‌کند و Traffic برابر 50K Request/sec است.

استفاده از یک Key واحد:

```
product:123
```

باعث می‌شود مجبور شویم یک TTL برای کل Object انتخاب کنیم. TTL کوتاه باعث Reload بی‌دلیل Description می‌شود و TTL بلند Inventory را قدیمی می‌کند.

طراحی بهتر:

```
product:123:metadata
product:123:price
product:123:inventory
```

Policy نمونه:

```
Metadata TTL = 30 Minutes
Price TTL = 30 Seconds
Inventory TTL = 2 Seconds
```

این مثال نشان می‌دهد مشکل TTL همیشه با تغییر عدد حل نمی‌شود؛ گاهی Cache Key و Data Model باید براساس Change Profile تفکیک شوند.

## Production Example: Exchange Rate API

فرض کنید:

```
Rates Update Every 30 Seconds
Traffic = 100K Requests/sec
Staleness Up to 2 Seconds Is Acceptable
Application Instances = 100
```

طراحی مناسب:

```
L2 Redis Updated Every 30 Seconds by Event or Job
L1 Local TTL = 1 Second
L2 Hard TTL = 2 Minutes
Serve Last Known Rate During Short Source Failure
```

L1 فشار Redis را کاهش می‌دهد، Update دوره‌ای L2 را تازه نگه می‌دارد و Hard TTL مانع قدیمی‌ماندن نامحدود داده می‌شود. برای Transaction مالی، نرخ نهایی همچنان باید با Business Rule و Source authoritative تأیید شود.

## Production Example: Feature Flags

Feature Flagها به‌ندرت تغییر می‌کنند، اما تغییر باید سریع منتشر شود:

```
Local Cache TTL = 30 Seconds
Push Invalidation on Change
Periodic Full Refresh = 5 Minutes
```

نقش هر بخش:

```
Push → Fast Propagation
Short Local TTL → Safety Net
Periodic Refresh → Recovery from Missed Events
```

تکیه صرف بر TTL ممکن است اعمال تغییر را به تأخیر بیندازد و تکیه صرف بر Push در صورت ازدست‌رفتن Event خطرناک است.

## Production Example: User Session

فرض کنید:

```
Idle Timeout = 30 Minutes
Absolute Lifetime = 12 Hours
```

هر Activity، Idle Expiration را تمدید می‌کند، اما Absolute Expiration ثابت باقی می‌ماند:

```
User Active Every 10 Minutes
→ Idle Timeout Keeps Extending
→ Session Still Ends after 12 Hours
```

این Strategy میان User Experience و Security Boundary تعادل ایجاد می‌کند.

## Review Scenario: E-commerce Platform

فرض کنید:

```
Traffic = 80K Requests/sec
Products = 4 Million
Top 20K Products = 75% Traffic
Product Cache TTL = Fixed 10 Minutes
All Products Warmed in One Batch
Price Can Change Every 2 Minutes
Description Changes Weekly
Inventory Changes Every Second
Every 10 Minutes DB QPS and P99 Spike
```

### مشکل اول

همه Productها در یک Batch و با TTL ده‌دقیقه‌ای نوشته شده‌اند:

```
Batch Warming + Same TTL
→ Synchronized Expiration
→ Periodic Cache Avalanche
```

افزایش DB QPS و P99 در بازه‌های ده‌دقیقه‌ای نشانه واضح همین مشکل است.

### آیا فقط TTL Jitter کافی است؟

خیر. Jitter موج Expiration را پخش می‌کند، اما Description، Price و Inventory همچنان Change Profile متفاوتی دارند و با TTL یکسان Cache شده‌اند.

راه‌حل کامل:

```
Split Cache by Change Profile
+ Per-Type TTL
+ TTL Jitter
+ Refresh Strategy
```

### Key Design مناسب

```
product:{id}:metadata
product:{id}:price
product:{id}:inventory
```

Policy نمونه:

```
Metadata TTL = 30–60 Minutes + Jitter
Price TTL = 30–60 Seconds
Inventory TTL = 1–3 Seconds
```

تصمیم نهایی Purchase نباید فقط بر Inventory Cache متکی باشد.

### راهکار برای Top 20K Product

برای داده‌های Hot:

```
Refresh Ahead
Soft/Hard TTL
Per-Key singleflight
L1 Cache Where Appropriate
```

داده Long-Tail می‌تواند با Cache Aside معمولی Load شود.

### رفتار هنگام Database Failure

Policy باید براساس ریسک هر داده باشد:

```
Metadata → Serve Stale up to 6 Hours
Price → Serve Stale Briefly or Mark Unavailable
Inventory → Do Not Promise Availability from Stale Data
```

### Metricهای حیاتی

```
Expirations per Second
Distinct Keys Expiring per Minute
Cache Miss Rate by Data Type
Source QPS by Data Type
Stale Served Count
Refresh Success and Failure
Entry Age Distribution
DB P99
Eviction Rate
```

### نتیجه طراحی

```
Data Decomposition
→ Per-Type TTL
→ Jittered Expiration
→ Hot-Key Refresh Ahead
→ Stale Policy by Business Risk
→ Source Protection
```

## Anti-patternهای رایج

### یک TTL برای همه‌چیز

```
Everything = 5 Minutes
```

این مدل ماهیت متفاوت داده‌ها را نادیده می‌گیرد.

### TTL بدون Jitter برای Batchها

```
Mass Write + Same TTL
```

باعث Expiration هماهنگ و Avalanche می‌شود.

### TTL بسیار بلند برای جبران Database ضعیف

```
Database Is Slow → Set TTL to 24 Hours
```

این کار مشکل Performance را به مشکل Staleness تبدیل می‌کند.

### TTL بسیار کوتاه برای تازگی

```
TTL = 1 Second for Everything
```

ممکن است Cache تقریباً بی‌اثر و Source Load شدید شود.

### Cache کردن خطای موقت به‌عنوان Not Found

```
DB Timeout → Cache NOT_FOUND
```

Failure موقت به نتیجه اشتباه طولانی تبدیل می‌شود.

### Sliding TTL برای داده Mutable

Hot Key ممکن است هرگز Expire نشود و مقدار قدیمی دائماً باقی بماند.

### No Expiration بدون Versioning و Cleanup

Keyهای قدیمی برای همیشه Memory مصرف می‌کنند.

### تغییر TTL بدون مشاهده Metric

تغییر TTL می‌تواند Source QPS، Memory، Hit Ratio، Staleness و Tail Latency را تغییر دهد و باید مانند تغییر مهم Capacity مدیریت شود.

## Production Insight

TTL Strategy پاسخ ساده‌ای به سؤال «پنج دقیقه یا ده دقیقه؟» نیست. TTL باید همراه با Data Classification، Invalidation، Refresh، Source Protection، Memory Budget و Business Staleness طراحی شود.

مدل ضعیف:

```
Pick an Arbitrary TTL
→ Apply Everywhere
→ Hope for a Good Hit Ratio
```

مدل Production:

```
Classify Data
→ Define Staleness Budget
→ Measure Source Cost
→ Select Base TTL
→ Add Jitter
→ Define Refresh and Invalidation
→ Define Stale Policy
→ Observe and Tune
```

اصل مهم:

```
TTL Is Not Just an Expiration Timer;
It Is a Contract Between Freshness, Load and Availability
```

## جمع‌بندی و Mental Model نهایی

TTL مشخص می‌کند یک Cache Entry چه مدت بدون تأیید مجدد از Source قابل استفاده باشد. TTL کوتاه معمولاً Freshness را افزایش می‌دهد، اما Cache Efficiency را کاهش و Source Load را افزایش می‌دهد. TTL بلند Hit Ratio را بیشتر می‌کند، اما Staleness، Memory Usage و Invalidation Risk را بالا می‌برد.

اجزای اصلی یک TTL Strategy مناسب:

```
Per-Data-Type TTL
Business Staleness Budget
Fixed or Dynamic TTL
TTL Jitter
Absolute Expiration
Sliding Expiration with Maximum Lifetime
Soft TTL and Hard TTL
Serve Stale and Stale-If-Error
Refresh Ahead
Negative Cache TTL
Explicit Invalidation
Versioned Keys
Memory-Aware TTL
Multi-Level TTL
Centralized Policy
Configuration Validation
Observability
```

مدل ذهنی نهایی این است:

**TTL مدت زمانی نیست که صرفاً دوست داریم داده در Cache باقی بماند؛ TTL حداکثر زمانی است که با توجه به Freshness موردنیاز، هزینه Source، ظرفیت Cache، ریسک Business و رفتار Failure حاضر هستیم بدون بازبینی مجدد به آن نسخه از داده اعتماد کنیم.**
