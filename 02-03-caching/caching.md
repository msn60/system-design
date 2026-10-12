# Performance Pattern: Caching

در این بخش، Caching را به‌عنوان یک Performance Pattern بررسی می‌کنیم، نه به‌عنوان موضوعی در Distributed Systems یا Reliability. تمرکز اصلی روی این است که Cache چگونه Latency را کاهش می‌دهد، Throughput را افزایش می‌دهد، مصرف منابع را کمتر می‌کند و Tail Latency (P95/P99) را بهبود می‌بخشد. ممکن است در طول بحث به Consistency یا Availability اشاره شود، اما فقط در حد Trade-off.

## Problem Statement

Caching یکی از مهم‌ترین Pattern ها Performance است. ایده اصلی بسیار ساده است: به‌جای اینکه هر بار یک عملیات گران را دوباره انجام دهیم، نتیجه آن را ذخیره می‌کنیم و در درخواست‌های بعدی از همان نتیجه استفاده می‌کنیم.

عملیات گران می‌تواند شامل موارد زیر باشد:

- Database Query
- API Call
- File Read
- Computation
- Aggregation
- Feed Generation

مثال: فرض کنید endpoint زیر را داریم:  GET /users/123 Flow فعلی: Client → API → PostgreSQL → Response

و Query زیر برای هر درخواست اجرا می‌شود: SELECT \* FROM users WHERE id = 123;

اگر اجرای Query حدود 50ms طول بکشد و سیستم 1000 request/sec دریافت کند، عملاً 1000 DB query/sec روی دیتابیس ایجاد می‌شود. اگر Cache اضافه کنیم:

```
Client → API → Redis → PostgreSQL
```

در این حالت:

```
Cache Hit  → 1ms
Cache Miss → 50ms
```

نتیجه:

- Latency کاهش پیدا می‌کند.
- DB Load کاهش پیدا می‌کند.
- Throughput افزایش پیدا می‌کند.

### Mental Model

برای درک Cache، فرض کن هر بار برای پاسخ دادن به یک سوال باید به آرشیو شرکت مراجعه کنی. آرشیو:

- بزرگ
- دور
- کند

است. اما اگر سوالات پرتکرار را روی میز خودت نگه داری:

- سریع
- ارزان
- در دسترس

خواهند بود. در این تشبیه:

- Cache = میز کنار دست
- Database = آرشیو اصلی

### چرا Cache این‌قدر موثر است؟

دلیل اصلی اختلاف بسیار زیاد سرعت بین لایه‌های مختلف سیستم است. 

| **Operation** | **Latency تقریبی** |
|---|---|
| CPU Cache | چند ns |
| RAM | ده‌ها ns |
| Redis | حدود 0.1 تا 1 ms |
| Local SSD | صدها μs |
| Database Query | چند ms تا صدها ms |
| Remote Service | ده‌ها تا صدها ms |

به همین دلیل گاهی فقط اضافه کردن Redis می‌تواند: 

```
P95 = 300ms
↓
P95 = 20ms
```

را رقم بزند.

### انواع Cache از دید Performance

### L1 Cache (In-Memory Cache)

داخل خود Process قرار دارد. مثال: map\[string]User 

مزایا:

- فوق‌العاده سریع
- بدون Network Call
- Latency بسیار کم

معایب:

- Memory محدود
- بین Replica ها Share نمی‌شود

مثال:

```
API Instance A → Cache A
API Instance B → Cache B
```

### L2 Cache (Distributed Cache)

معمولاً Redis یا Memcached:  

```
Flow: API → Redis → DB
```

مزایا:

- Shared
- Scalable
- ظرفیت بیشتر

معایب:

- Network Call دارد
- پیچیدگی بیشتری دارد

### Cache Hit و Cache Miss

Cache Hit یعنی داده داخل Cache پیدا شده است. 

```
Flow: Request → Cache → Found → Return
```

مثلاً: 1ms

Cache Miss یعنی داده داخل Cache وجود ندارد. 

```
Flow: Request → Cache → Miss → DB → Cache Update → Return
```

مثلاً: 60ms

### Cache Hit Ratio

مهم‌ترین Metric در سیستم‌های Cache محور است. فرمول:

```
Hit Ratio = Hits / Total Requests
```

مثال:

```
1000 Request
900 Hit
100 Miss
```

نتیجه: Hit Ratio = 90%

هرچه Hit Ratio بیشتر باشد:

- Latency کمتر می‌شود.
- DB Load کمتر می‌شود.
- Throughput بیشتر می‌شود.

### تاثیر روی SLI ها

فرض کنیم:

- DB Read = 100ms
- Redis Read = 2ms

بدون Cache عموما  P95 = 100ms اما با Hit Ratio برابر 90٪:

- 90% Requests → 2ms
- 10% Requests → 100ms

نتیجه:

- P50 به‌شدت بهتر می‌شود.
- P95 به‌شدت بهتر می‌شود.
- مصرف CPU دیتابیس کاهش پیدا می‌کند.

### Cache از چه چیزی محافظت می‌کند؟

بسیاری از افراد فکر می‌کنند Cache فقط برای کاهش Latency است. در عمل، نقش مهم‌تر Cache محافظت از Dependency، به‌خصوص Database است. مثال: 10K RPS

- بدون Cache:

10K Query/sec روی DB

- با Hit Ratio برابر 95٪:

500 Query/sec روی DB

یعنی Cache عملاً 95٪ از بار Database را حذف کرده است.

### Trade-offs

هیچ Performance Pattern رایگان نیست و Cache هم هزینه دارد.

| **مزایا** | **هزینه‌ها** |
|---|---|
| Latency ↓ | Data Freshness ↓ |
| Throughput ↑ | Invalidation Complexity ↑ |
| DB Load ↓ | Memory Cost ↑ |
| Scalability ↑ | Operational Complexity ↑ |

### مهم‌ترین نکته در Caching

بزرگ‌ترین مشکل Cache معمولاً Storage نیست. **بزرگ‌ترین مشکل: Cache Invalidation است**. یعنی: چه زمانی Cache را پاک کنیم؟ مثال:  User Profile داخل Redis ذخیره شده است. کاربر ایمیل خود را تغییر می‌دهد. حالا سؤال این است: Redis معتبر است؟ یا DB؟ این همان نقطه‌ای است که بحث **Cache Consistency** مطرح می‌شود.

#### Production Insight

در Production معمولاً Cache برای سه نوع داده استفاده می‌شود: 

- Hot Data

مانند: User Profile و Product و Feature Flag 

- Expensive Data

مانند: Feed Generation و  Analytics و Aggregations

- Frequently Accessed Data

مانند: Configuration و Country List و Currency List

تا اینجا فقط خود Pattern اصلی Caching را بررسی کردیم. در ادامه باید Cache Strategy ها را بررسی کنیم:

- Cache Aside
- Read Through
- Write Through
- Write Back
- Refresh Ahead

و سپس Failure Mode های مهم Cache:

- Cache Stampede
- Cache Avalanche
- Hot Key
- Cache Penetration

این موارد دقیقاً همان چیزهایی هستند که در سیستم‌های واقعی می‌توانند Performance را به‌شدت تحت تاثیر قرار دهند.

به همین دلیل، بعد از آشنایی با Fundamentals، اولین Strategy مهمی که باید بررسی شود Cache Aside است؛ زیرا بخش بزرگی از سیستم‌های Backend دنیا عملاً از همین الگو استفاده می‌کنند.
