#  Cache strategy: Cache Aside (Lazy Loading)

اگر خود **Cache** را **«ابزار»** در نظر بگیریم، **Cache Strategy** در واقع «**روش استفاده از ابزار**» است. در مصاحبه‌های System Design و حتی در محیط Production، وقتی کسی می‌گوید «ما Redis داریم»، تقریباً اطلاعات مفیدی نداده است. سؤال مهم‌تر این است که «چطور از Redis استفاده می‌کنید؟». تفاوت بین یک Cache موفق و یک Cache فاجعه‌بار معمولاً در Strategy مشخص می‌شود.

Cache Aside یا Lazy Loading مهم‌ترین و رایج‌ترین Cache Pattern در دنیاست و در اکثر سرویس‌های Go، Java، .NET و Node.js استفاده می‌شود.

## Problem Statement

فرض کنیم یک User Service داریم که برای دریافت اطلاعات کاربر مستقیماً به PostgreSQL متصل می‌شود:

User Service → PostgreSQL . برای هر درخواست، Query زیر اجرا می‌شود:

```
SELECT * FROM users WHERE id = ?;
```

در این حالت: DB Load ↑ و Latency ↑ و Throughput ↓

**هدف Cache Aside این است که ابتدا Cache بررسی شود و فقط در صورت نبودن داده به Database مراجعه شو**د.

```
Check Cache → If Miss → DB
```

## Mental Model

**ساده‌ترین روش برای به خاطر سپردن Cache Aside: همیشه Application مسئول Cache است**. یعنی Redis خودش تصمیم نمی‌گیرد چه زمانی داده خوانده، نوشته، حذف یا به‌روزرسانی شود. تمام تصمیم‌ها در Application گرفته می‌شود:

- Cache Hit؟
- Cache Miss؟
- Cache Update؟
- Cache Delete؟

بنابراین Application کنترل کامل رفتار Cache را در اختیار دارد.

## Request Flow

### Cache Hit

در حالت Cache Hit، داده در Redis وجود دارد:

```
Client → Application → Redis → Found → Return
```

مثلاً: Redis Read = 1ms هیچ تماس با Database انجام نمی‌شود.

### Cache Miss

در حالت Cache Miss، داده در Cache وجود ندارد:

```
Client → Application → Redis → Miss → Database → Cache Update → Return
```

مثلاً: Redis Check = 1ms و DB Query = 50ms و Redis Write = 1ms و کل زمان پاسخ: 52ms اما درخواست بعدی: 1ms

خواهد بود، چون داده داخل Cache قرار گرفته است.

### Pseudo Flow

نمونه ساده در Go:

```
func GetUser(id string) User {
    user, found := redis.Get(id)
    if found {
        return user
    }
    user = db.GetUser(id)
    redis.Set(id, user)
    return user
}
```

این دقیقاً پیاده‌سازی Cache Aside است.

### چرا به آن Lazy Loading می‌گویند؟

زیرا Cache فقط زمانی پر می‌شود که کسی واقعاً آن داده را درخواست کند. فرض کنید: 1M User در سیستم وجود دارد. اما فقط: 100K User فعال هستند. در این حالت Cache فقط برای همان 100K کاربر فعال ساخته می‌شود و نه برای کل یک میلیون رکورد. این موضوع باعث استفاده بهینه از Memory می‌شود.

## Performance Impact

فرض کنیم:

```
DB Read = 100ms
Redis Read = 2ms
Hit Ratio = 95%
```

در این حالت:

```
95% Requests → 2ms
5% Requests → 100ms
```

نتیجه:

- Latency ↓ شدید
- Throughput ↑ شدید
- DB CPU ↓ شدید
- Connection Usage ↓

هرچه Hit Ratio بیشتر باشد، فشار روی Database کمتر و ظرفیت سیستم بیشتر خواهد شد.

### چرا تقریباً همه جا استفاده می‌شود؟

دلیل محبوبیت Cache Aside سه ویژگی اصلی آن است:

- Simple
- Cheap
- Flexible

مزایای اصلی:

- پیاده‌سازی ساده
- مستقل از تکنولوژی Cache
- کنترل کامل توسط Application
- Debug آسان
- انعطاف بالا برای سیاست‌های مختلف Cache

به همین دلیل در اکثر Microservice های مدرن، Cache Aside رایج‌ترین Strategy است.

## بزرگ‌ترین مشکل Cache Aside

بزرگ‌ترین چالش Cache Aside: مشکل **Cache Invalidation** است. مثال: User Profile داخل Redis ذخیره شده است. کاربر ایمیل خود را تغییر می‌دهد. Database مقدار جدید را دارد اما Cache هنوز مقدار قدیمی را نگه داشته است. نتیجه: Stale Data پس کاربر اطلاعات قدیمی دریافت می‌کند.

### راه‌حل رایج: Cache Invalidate On Write

رایج‌ترین راه‌حل این است که بعد از Update موفق در Database، آنگاه Cache حذف شود.

```
Flow: Update DB → Delete Cache
```

مثال:

```
UpdateUser()
redis.Del(userID)
```

در درخواست بعدی:

```
Cache Miss → DB → New Cache
```

داده جدید مجدداً داخل Cache قرار می‌گیرد. این روش با نام: Cache Invalidate On Write شناخته می‌شود.

##  Failure Mode شماره 1: Cache Stampede (Thundering Herd)

فرض کنید کلید زیر در Cache وجود ندارد: User:123 همزمان: 500 Request برای آن می‌رسد. همه درخواست‌ها: Cache Miss می‌خورند. همه همزمان: DB Query اجرا می‌کنند. نتیجه: DB Explosion یا همان:  Cache Stampede یا Thundering Herd که این یکی از معروف‌ترین مشکلات سیستم‌های Cache محور است.

## Failure Mode شماره 2: Cold Start

فرض کنید Redis Restart شود. در این حالت: تمام Cache پاک می‌شود. و ناگهان تمام درخواست‌ها مستقیماً به Database می‌رسند: Redis Empty → DB نتیجه: DB Saturation و در ادامه:

- Latency ↑
- Connection Usage ↑
- Timeout ↑

## Production Techniques

### TTL

برای هر Cache Entry زمان اعتبار تعریف می‌شود: مثلا 5 min یا 10 min یا 1 hour

### Random TTL

به جای: 600 sec از: 550-650 sec استفاده می‌شود. هدف جلوگیری از Expire شدن همزمان تعداد زیادی Cache Entry است.

### Singleflight

در Go، کتابخانه زیر بسیار پرکاربرد است: 

```
golang.org/x/sync/singleflight
```

فرض کنید 100 درخواست همزمان برای یک Cache Miss برسد.

- بدون Singleflight: 

```
100 DB Query	
```

- با Singleflight:

```
1 DB Query
```

و بقیه درخواست‌ها منتظر همان نتیجه می‌مانند.

##  Production Insight

اگر فقط یک Cache Strategy در کل دوران کاری خود استفاده کنید، به احتمال زیاد همان Cache Aside خواهد بود. دلایل:

- Simple
- Predictable
- Easy to Debug

اما دقیقاً به خاطر همین محبوبیت، اکثر مشکلات Cache در Production نیز از همین الگو ناشی می‌شوند:

- Cache Stampede
- Stale Data
- Hot Keys
- Cold Start

بنابراین یک Senior Engineer باید Failure Mode های Cache Aside را به همان اندازه خود Pattern بشناسد.

## Production Debugging Example: افت Cache Hit Ratio

فرض کنید وضعیت فعلی سیستم:

```
DB Read = 50ms
Redis Read = 1ms
Hit Ratio = 90%
```

باشد. اگر Hit Ratio ناگهان به: 40% سقوط کند، چه اتفاقی می‌افتد؟ 

### اثر روی Latency

قبلاً:

```
90% → 1ms
10% → 50ms
```

اکثر درخواست‌ها از Redis پاسخ می‌گرفتند. اکنون:

```
40% → 1ms
60% → 50ms
```

بیشتر درخواست‌ها مجبورند به Database مراجعه کنند. نتیجه:

- P50 ↑
- P95 ↑
- P99 ↑

Mental Model: قبلاً: Redis = مسیر اصلی و DB = مسیر استثنا، اما اکنون: DB = مسیر اصلی و Redis = مسیر استثنا

### اثر روی Throughput

فرض کنیم: 10,000 RPS داریم. قبلاً:

```
9,000 Request → Redis
1,000 Request → DB
```

اکنون:

```
4,000 Request → Redis
6,000 Request → DB
```

یعنی: DB Load ≈ 6 برابر می‌شود. در نتیجه: Throughput Capacity ↓ زیرا Database معمولاً Bottleneck اصلی سیستم است.

### اثر روی Connection Pool

قبلاً فقط 10٪ درخواست‌ها از Connection Pool استفاده می‌کردند. اکنون: 60٪ درخواست‌ها به Database می‌رسند. نتیجه: Connection Pool Saturation سپس:

- Wait Count ↑
- Wait Duration ↑
- Queue Time ↑
- P95 ↑
- P99 ↑
- Timeout ↑

و اگر وضعیت ادامه پیدا کند: Retry Storm یا حتی: DB Collapse ممکن است رخ دهد.

### اولین Metric هایی که در Grafana بررسی می‌کنیم

#### Cache Hit Ratio

  - redis_hit_ratio
  - cache_hit_ratio

برای اطمینان از اینکه واقعاً مشکل از Cache است.

#### DB Query Rate

queries/sec برای بررسی افزایش ناگهانی Query ها.

#### DB Connection Pool

- active_connections
- wait_count
- wait_duration

#### Request Latency

  - P95
  - P99
  - request_duration

#### Redis Errors

  - Redis Restart
  - Redis Timeout
  - Redis Eviction

#### Redis Memory Usage

برای بررسی Eviction یا کمبود حافظه. 

این دقیقاً همان مدل فکری مورد نیاز برای Production Debugging در سیستم‌های Cache محور است.
