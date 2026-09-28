# Cache Strategy: Read Through

اگر خود Cache را «ابزار» در نظر بگیریم، Cache Strategy در واقع «روش استفاده از ابزار» است. بسیاری از مهندس‌ها تصور می‌کنند Read Through نسخه پیشرفته‌تر Cache Aside است، اما در عمل معمولاً چنین نیست. مهم‌ترین چیزی که باید یاد بگیریم این نیست که Read Through چیست، بلکه این است که دقیقاً چه تفاوتی با Cache Aside دارد و در چه شرایطی ارزش استفاده دارد.

## Problem Statement

در Cache Aside، Application خودش مسئول مدیریت Cache است:

```
Application → Check Cache → If Miss → Read DB → Update Cache
```

تمام منطق Cache داخل Application قرار دارد. سوال اینجاست که آیا Application واقعاً باید مسئول مدیریت Cache باشد و آیا نمی‌توان این مسئولیت را به خود Cache سپرد؟ پاسخ به این سؤال، **Read Through** است.

## Mental Model

اگر بخواهیم تفاوت Cache Aside و Read Through را در یک جمله حفظ کنیم:

```
Cache Aside: Application manages cache
Read Through: Cache manages cache
```

در Cache Aside، Application تصمیم می‌گیرد که Cache Hit رخ داده یا Cache Miss، چه زمانی از Database بخواند و چه زمانی Cache را به‌روزرسانی کند. در Read Through، Application فقط داده را درخواست می‌کند و Cache خودش تصمیم می‌گیرد که داده از Cache بازگردانده شود یا از Database بارگذاری و سپس Cache شود. به بیان ساده:

```
Cache Aside: Application → Cache + Database
Read Through: Application → Cache
```

## Architecture Flow

### Cache Aside

```
Client → Application → Redis
           ↓
          Miss
           ↓
       Database
           ↓
    Redis Update
          ↓
       Response
```

در این الگو، Application از وجود Database و Cache آگاه است و منطق مدیریت Cache را در اختیار دارد.

### Read Through

```
Client → Application → Cache
              ↓
            Miss
              ↓
           Database
              ↓
          Cache Update
              ↓
           Response
```

ظاهر Flow تقریباً مشابه Cache Aside است، اما تفاوت اصلی در Ownership قرار دارد. در Cache Aside، همیشه Application مالک منطق Cache است؛ در Read Through، عملا  Cache Layer مالک این منطق است.

**مثال واقعی**: فرض کنیم سرویس کاربران داریم.

```
Cache Aside
func GetUser(id string) User {
    user, found := cache.Get(id)
    if found {
        return user
    }
    user = db.Get(id)
    cache.Set(id, user)
    return user
}
```

در این حالت Application تمام جزئیات Cache را می‌داند.

```
Read Through
func GetUser(id string) User {
    return cache.Get(id)
}
```

اما پشت صحنه، Cache Layer مسئول اجرای فرآیند زیر است:

```
Cache Miss → Load DB → Store Cache → Return
```

بنابراین **پیچیدگی از Application حذف شده و به Cache Layer منتقل** می‌شود.

## Performance Impact

از دید Performance، تفاوت قابل توجهی بین Cache Aside و Read Through وجود ندارد. فرض کنیم:

```
DB Read = 100ms
Cache Read = 2ms
Hit Ratio = 95%
```

در هر دو Pattern:

```
95% Requests → 2ms
5% Requests → 100ms
```

بنابراین:

- Latency ≈ برابر
- Throughput ≈ برابر
- DB Load ≈ برابر

این نکته بسیار مهم است. Read Through عمدتاً یک Pattern مربوط به Code Organization، Abstraction و Ownership است، نه یک Performance Optimization جدید.

## مزایا

### سادگی Application

در Cache Aside معمولاً منطق Cache در نقاط مختلف Application تکرار می‌شود:

- Get Cache
- Check Hit
- Get DB
- Update Cache

در Read Through، این پیچیدگی حذف شده و معمولاً یک فراخوانی ساده کافی است:

```
user := cache.Get(id)
```

###  Centralized Cache Logic

تمام منطق Cache مانند TTL، Cache Miss Handling، Loading، Refresh و Metrics در یک محل متمرکز می‌شود. این موضوع باعث می‌شود رفتار Cache در کل سیستم یکسان و استاندارد باشد.

### کاهش تکرار

اگر سرویس‌های متعددی مانند User Service، Product Service، Wallet Service و Order Service داشته باشیم، همه می‌توانند از یک مکانیزم مشترک Cache استفاده کنند و نیازی به تکرار منطق Cache در هر سرویس وجود نخواهد داشت.

## معایب

### وابستگی بیشتر

Application به Cache Layer وابسته می‌شود. در نتیجه اگر Cache Layer دچار مشکل شود، تعداد زیادی از سرویس‌ها به‌طور همزمان تحت تأثیر قرار می‌گیرند.

### انعطاف کمتر

فرض کنیم سیاست Cache به شکل زیر باشد:

```
User TTL = 1 Hour
Product TTL = 5 Minutes
Wallet TTL = 30 Seconds
```

در Cache Aside، Application می‌تواند این سیاست‌ها را به‌سادگی کنترل کند. اما در Read Through، این منطق باید داخل Cache Layer مدیریت شود که معمولاً پیچیدگی بیشتری ایجاد می‌کند.

### Debug سخت‌تر

در Cache Aside مسیر کامل Cache Hit، Cache Miss، DB Read و Cache Update مستقیماً داخل Application قابل مشاهده است. اما در Read Through بخشی از این رفتار پشت Cache Layer پنهان می‌شود. در نتیجه هنگام بروز مشکلاتی مانند Hit Ratio Drop، Stale Data یا Latency Spike، پیدا کردن علت اصلی می‌تواند دشوارتر باشد.

## چرا Read Through کمتر از Cache Aside استفاده می‌شود؟

در نگاه اول ممکن است این سوال مطرح شود که اگر Read Through کد کمتری نیاز دارد و بخشی از پیچیدگی Cache را از Application پنهان می‌کند، چرا اکثر تیم‌های Backend همچنان Cache Aside را ترجیح می‌دهند؟

**دلیل اصلی، Control است**. تیم‌های Backend معمولاً ترجیح می‌دهند سیاست‌های Cache مانند TTL، Refresh Strategy، Invalidation Policy، Cache Policy و Fallback Logic را مستقیماً داخل Application کنترل کنند. در Cache Aside تمام این تصمیم‌ها در اختیار Application است و تیم می‌تواند برای هر نوع داده یا هر Use Case رفتار متفاوتی تعریف کند.

برای مثال ممکن است در یک سیستم تصمیم بگیریم:

```
User TTL = 1 Hour
Product TTL = 5 Minutes
Wallet TTL = 30 Seconds
```

یا برای برخی داده‌ها از Refresh Strategy خاصی استفاده کنیم و برای برخی دیگر هنگام Update، Cache را Invalid کنیم. در Cache Aside پیاده‌سازی چنین سیاست‌هایی معمولاً ساده‌تر و شفاف‌تر است.

در معماری Microservice، عموما Application Ownership معمولاً ارزش بیشتری از مخفی کردن Complexity دارد. بسیاری از تیم‌ها ترجیح می‌دهند منطق Cache را به‌صورت شفاف داخل Service نگه دارند تا بتوانند رفتار سیستم را راحت‌تر درک، Debug و تغییر دهند.

به همین دلیل در بسیاری از شرکت‌های بزرگ و سیستم‌های Production، زمانی که تیم Backend مالک کامل Application و Redis است، Cache Aside انتخاب رایج‌تری نسبت به Read Through محسوب می‌شود. در مقابل، Read Through بیشتر در محیط‌هایی دیده می‌شود که یک Shared Cache Layer، Data Access Layer یا Framework مشترک وجود دارد و هدف اصلی کاهش تکرار کد و یکپارچه‌سازی رفتار Cache در کل سیستم است.

## Performance-Oriented Trade-off

| **موضوع** | **Cache Aside** | **Read Through** |
|---|---|---|
| **مدیریت Cache** | Application | Cache Layer |
| **مدیریت Miss** | Application | Cache Layer |
| **انعطاف** | زیاد | کمتر |
| **سادگی کد** | کمتر | بیشتر |
| **Debugging** | ساده‌تر | سخت‌تر |
| **محبوبیت** | بسیار زیاد | کمتر |
| **Performance** | تقریباً برابر | تقریباً برابر |
| **Control** | زیاد | کمتر |

## Production Insight

یکی از اشتباهات رایج این است که تصور کنیم Read Through نسخه پیشرفته‌تر Cache Aside است. در واقع Read Through بیشتر یک تصمیم معماری برای مدیریت Complexity است تا یک بهینه‌سازی Performance.

اگر تیم کنترل کامل Application، Redis و Cache Policy را در اختیار دارد، معمولاً Cache Aside انتخاب بهتری است. اما اگر سازمان دارای Shared Cache Layer، Common Data Access Layer یا Framework-Based Access باشد، Read Through می‌تواند ارزشمند شود.

## Mental Model نهایی

اگر فقط یک جمله از این فصل به خاطر بسپاری:

- **Cache Aside**: در این حالت Application می‌داند Cache وجود دارد و خودش آن را مدیریت می‌کند.
- **Read Through**: در این حالت Application فقط داده می‌خواهد و Cache مسئول پیدا کردن، بارگذاری و ذخیره کردن آن داده است.

## Review Questions & Answers

**سوال 1**: فرض کن در یک User Service با Go و Redis کار می‌کنی و تیم Backend مالک کامل Application و Redis است. آیا Read Through واقعاً مزیت مهمی نسبت به Cache Aside ایجاد می‌کند؟

پاسخ: معمولاً خیر. در چنین شرایطی تیم کنترل کامل Application و Redis را در اختیار دارد و نیازی به مخفی کردن منطق Cache پشت یک Cache Layer وجود ندارد. مزیت اصلی Read Through کاهش Boilerplate Code است، اما در مقابل بخشی از Control از Application گرفته می‌شود. به همین دلیل در اکثر Microservice های مدرن، Cache Aside انتخاب رایج‌تری است.

**سوال 2**: اگر فردا بخواهی TTLهای متفاوت برای User، Product و Wallet داشته باشی، کدام Pattern انعطاف بیشتری به تو می‌دهد؟

پاسخ: Cache Aside. زیرا Application مستقیماً سیاست‌های Cache را مدیریت می‌کند و می‌تواند برای هر نوع داده TTL، Refresh Strategy یا Invalidation Strategy متفاوتی تعریف کند. در Read Through این منطق باید به Cache Layer منتقل شود و معمولاً پیچیده‌تر می‌شود.

**سوال 3**: اگر Hit Ratio ناگهان افت کند، Debug کردن در کدام Pattern ساده‌تر است؟

پاسخ: معمولاً Cache Aside. زیرا تمام مسیر Cache Hit، Cache Miss، DB Read و Cache Update به‌صورت شفاف داخل Application قرار دارد و Trace کردن رفتار سیستم راحت‌تر است. در Read Through بخشی از این رفتار داخل Cache Layer پنهان شده و برای پیدا کردن علت افت Hit Ratio یا افزایش Latency باید لایه‌های بیشتری بررسی شوند.

## جمع‌بندی

از دید Performance، تفاوت محسوسی بین Cache Aside و Read Through وجود ندارد. هر دو Latency را کاهش می‌دهند، Throughput را افزایش می‌دهند و فشار روی Database را کم می‌کنند. **تفاوت اصلی در Ownership و محل قرارگیری منطق Cache است.** Cache Aside کنترل و انعطاف بیشتری به Application می‌دهد، در حالی که Read Through سادگی بیشتری برای Application فراهم می‌کند. به همین دلیل در اکثر Microservice های مدرن، Cache Aside انتخاب پیش‌فرض و رایج‌تر است.
