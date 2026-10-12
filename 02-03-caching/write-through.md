# Cache Strategy: Write Through

تا اینجا Cache Aside و Read Through بیشتر **روی Read Path تمرکز داشتند**؛ یعنی سؤال اصلی این بود که وقتی داده را می‌خوانیم، Cache چه نقشی دارد. اما در Write Through وارد Write Path می‌شویم و **سؤال اصلی تغییر می‌کند: وقتی داده تغییر می‌کند، Cache چه رفتاری باید داشته باشد؟** این همان نقطه‌ای است که بحث‌هایی مانند Stale Data، Cache Invalidation و Consistency وارد تصمیم‌گیری می‌شوند، البته در این فصل تمرکز همچنان روی Performance است.

## Problem Statement

در Cache Aside، هنگام Update معمولاً ابتدا Database به‌روزرسانی می‌شود و سپس Cache حذف می‌شود: Write DB → Delete Cache. برای مثال:

```
UpdateUser(user)
redis.Del(userID)
```

در این روش، درخواست بعدی با Cache Miss مواجه می‌شود و دوباره داده را از Database می‌خواند و Cache را می‌سازد: Cache Miss → Read DB → Rebuild Cache. این روش ساده و رایج است، اما یک نکته دارد: بعد از Update، تا زمانی که اولین Read بعدی اتفاق بیفتد، Cache برای آن داده خالی است. سوال Write Through از همین‌ جا شروع می‌شود: چرا بعد از Update، همان لحظه Cache را هم به‌روزرسانی نکنیم؟

## Mental Model

ساده‌ترین تعریف Write Through این است: هر Write همزمان در Cache و Database اعمال می‌شود. Flow ذهنی آن به شکل Application → Cache → Database است. در مقابل، در Cache Aside معمولاً Flow نوشتن به شکل Application → Database → Invalidate Cache بود. **قانون اصلی در Write Through این است که Cache و DB با هم تغییر می‌کنند و هدف این است که بعد از هر Write، باید Cache همچنان آماده و گرم باقی بماند**.

## Write Flow

فرض کنید کاربر ایمیل خود را تغییر می‌دهد. قبل از تغییر، هم Redis و هم Database مقدار قدیمی را دارند: Redis: user:123, email=old@example.com و DB: email=old@example.com. وقتی درخواست Update با مقدار email=new@example.com می‌رسد، در Write Through داده هم در Cache و هم در Database نوشته می‌شود: Application → Write Cache → Write DB → Success. بعد از اتمام عملیات، هر دو منبع مقدار جدید را دارند: Redis: email=new@example.com و DB: email=new@example.com. نتیجه این است که Cache بعد از Write همچنان گرم است و درخواست Read بعدی احتمالاً Cache Hit خواهد بود.

## Pseudo Flow

یک نمونه ساده از Write Through می‌تواند این شکل را داشته باشد:

```
func UpdateUser(user User) error {
    err := db.Update(user)
    if err != nil {
        return err
    }
    err = cache.Set(user.ID, user)
    if err != nil {
        return err
    }
    return nil
}
```

در بعضی پیاده‌سازی‌ها ممکن است ابتدا Cache و سپس Database نوشته شود:

```
func UpdateUser(user User) error {
    cache.Set(user.ID, user)
    db.Update(user)
    return nil
}
```

اما **ترتیب نوشتن بسیار مهم است**، چون اگر یکی از دو عملیات موفق شود و دیگری fail شود، احتمال ناسازگاری بین Cache و Database ایجاد می‌شود. در سیستم‌های Production معمولاً باید این رفتار به‌صورت دقیق طراحی شود: آیا Database منبع حقیقت است؟ اگر Cache Update fail شد، آیا کل Update باید fail شود یا فقط Cache باید بعداً اصلاح شود؟ این‌ها تصمیم‌های مهمی هستند که روی Consistency، Availability و Performance اثر می‌گذارند.

## تفاوت با Cache Aside

در Cache Aside هنگام Write معمولاً Database به‌روزرسانی و Cache حذف می‌شود: Write DB → Delete Cache. سپس در Read بعدی، Cache Miss رخ می‌دهد و داده جدید از Database خوانده و Cache دوباره ساخته می‌شود: Read → Cache Miss → DB → Cache Rebuild. در Write Through هنگام Write، عملا Cache نیز همزمان با Database به‌روزرسانی می‌شود: Write Cache → Write DB. بنابراین در Read بعدی معمولاً Cache Hit رخ می‌دهد: Read → Cache Hit.

از نظر ذهنی، Cache Aside یعنی Lazy Update؛ پس Cache فقط وقتی دوباره ساخته می‌شود که کسی داده را بخواهد. اما Write Through یعنی Immediate Update؛ به محض تغییر داده، Cache نیز به‌روزرسانی می‌شود.

## Performance Impact

از دید Read Performance، عملا Write Through بسیار جذاب است، چون بعد از هر Update، Cache آماده است و احتمال Cache Miss کمتر می‌شود. اگر User Profile زیاد خوانده شود، بعد از هر تغییر، نسخه جدید آن در Cache وجود دارد و درخواست‌های بعدی مستقیماً از Cache پاسخ می‌گیرند. در نتیجه Hit Ratio بالا می‌رود، Read Latency کاهش پیدا می‌کند و Cold Read کمتر اتفاق می‌افتد.

اما این مزیت در سمت Write هزینه دارد. در Write Through هر Write باید حداقل دو عملیات انجام دهد: نوشتن در Cache و نوشتن در Database. بنابراین Write Latency معمولاً بیشتر از حالتی می‌شود که فقط Database را Update کنیم و Cache را Delete کنیم. پس Write Through به‌طور مستقیم Read Path را بهتر می‌کند، اما می‌تواند Write Path را کندتر و پیچیده‌تر کند.

## مزایا

اولین مزیت Write Through این است که Cache همیشه گرم‌تر می‌ماند. بعد از هر Update، داده جدید داخل Cache قرار دارد و درخواست بعدی احتمالاً Cache Hit خواهد بود. مزیت دوم این است که احتمال Stale Data کمتر می‌شود، چون Cache و Database به‌صورت هماهنگ تغییر می‌کنند. مزیت سوم، بهبود Read Latency برای داده‌های Hot است؛ یعنی داده‌هایی که زیاد خوانده می‌شوند و کمتر تغییر می‌کنند، مثل User Profile، Product Catalog، Configuration و Feature Flags.

## معایب

مهم‌ترین عیب Write Through افزایش Write Latency است، چون هر Write باید هم Cache و هم Database را درگیر کند. در Cache Aside معمولاً Flow به شکل Write DB → Delete Cache است، اما در Write Through به شکل Write Cache → Write DB یا Write DB → Write Cache خواهد بود. بنابراین هر Write حداقل دو عملیات دارد و این می‌تواند Latency مسیر Write را افزایش دهد.

عیب دوم کاهش انعطاف در Availability است. فرض کنید Database سالم است اما Redis Down شده است. حالا سؤال مهم این است: آیا باید Write را fail کنیم یا اجازه دهیم Database Update شود و Cache بعداً اصلاح شود؟ اگر بگوییم هر دو باید حتماً موفق شوند، Availability کاهش پیدا می‌کند، چون خرابی Cache می‌تواند مسیر Write را هم خراب کند. اگر اجازه دهیم Write در Database موفق شود ولی Cache Update fail شود، **وارد مسئله ناسازگاری موقت یا Stale Cache** می‌شویم.

عیب سوم Write Amplification است. فرض کنید سیستم 1M Update/day دارد اما فقط 10K Read/day. در این حالت ممکن است Cache را بارها و بارها برای داده‌هایی Update کنیم که اصلاً خوانده نمی‌شوند. یعنی هزینه نوشتن در Cache پرداخت می‌شود، اما منفعت Read Performance آن به‌دست نمی‌آید. در چنین سناریویی Write Through ممکن است انتخاب اشتباهی باشد.

## چه زمانی مناسب است؟

Write Through **برای داده‌هایی مناسب است که نسبت Read به Write در آن‌ها بسیار بالاست**؛ یعنی Read >> Write. این داده‌ها زیاد خوانده می‌شوند، کمتر تغییر می‌کنند و بهتر است بعد از هر تغییر، نسخه جدیدشان آماده خواندن باشد. مثال‌های مناسب **شامل User Profile، Product Catalog، Configuration و Feature Flags** هستند.

در مقابل، Write Through معمولاً برای workload هایی که Write زیاد دارند مناسب نیست؛ مانند Event Streams، Logs، Metrics یا High Write Workloads. در این سیستم‌ها، هزینه Update دائمی Cache ممکن است بیشتر از منفعت آن باشد.

**مثال واقعی**

فرض کنید یک فروشگاه اینترنتی داریم و اطلاعات Product روزانه فقط 100 بار Update می‌شود، اما 1M بار Read می‌شود. در این حالت Write Through منطقی است، چون هزینه 100 بار Update کردن Cache بسیار کم‌تر از منفعتی است که از 1 میلیون Read سریع‌تر به‌دست می‌آید.

اما فرض کنید سیستم Order Events داریم که روزانه 1M Write و فقط 50K Read دارد. در این حالت Write Through احتمالاً انتخاب مناسبی نیست، چون بیشتر انرژی سیستم صرف Update کردن Cache برای داده‌هایی می‌شود که نسبتاً کمتر خوانده می‌شوند. برای چنین workload هایی معمولاً الگوهای دیگری مثل event-driven processing، append-only storage یا حتی عدم استفاده از Cache در مسیر Write مناسب‌تر است.

## Production Insight

یکی از اشتباهات رایج این است که تصور کنیم Write Through همیشه بهتر از Cache Aside است، چون Cache را همیشه به‌روز نگه می‌دارد. **اما سؤال اصلی این است: Read بیشتر است یا Write؟** اگر Read >> Write باشد، Write Through می‌تواند ارزشمند باشد، چون هزینه اضافه روی Write در برابر کاهش Latency و افزایش Hit Ratio در Read قابل قبول است. اما اگر Write زیاد باشد، هزینه Update دائمی Cache می‌تواند بیشتر از منفعت آن شود.

به همین دلیل در بسیاری از سیستم‌های بزرگ، Cache Aside هنوز رایج‌تر از Write Through است. Cache Aside ساده‌تر است، کنترل بیشتری به Application می‌دهد و Cache را فقط زمانی می‌سازد که داده واقعاً خوانده شود. Write Through زمانی ارزشمند می‌شود که مطمئن باشیم داده پس از Update به‌دفعات زیاد خوانده خواهد شد.

## نکته پیاده‌سازی در Go

در Go، هنگام پیاده‌سازی Write Through باید برای هر دو عملیات Database و Cache از context.Context با timeout مشخص استفاده شود تا مسیر Write به‌صورت نامحدود منتظر Redis یا Database نماند. همچنین باید خط‌مشی خطا مشخص باشد: اگر db.Update موفق شد اما cache.Set fail شد، آیا باید خطا به کاربر برگردد یا فقط Cache برای اصلاح بعدی علامت‌گذاری شود؟ در بسیاری از سیستم‌های Production، همواره Database منبع حقیقت در نظر گرفته می‌شود و Cache یک لایه مشتق‌شده است؛ بنابراین اگر DB Update موفق باشد ولی Cache Update fail شود، ممکن است request موفق تلقی شود و Cache با TTL کوتاه، retry محدود، background repair یا invalidate جبران شود. اما برای داده‌هایی که Freshness بسیار مهم است، ممکن است سیاست سخت‌گیرانه‌تری انتخاب شود.

یک پیاده‌سازی Production-grade باید حداقل این موارد را مشخص کند: ترتیب نوشتن، timeout ها، retry policy، idempotency، رفتار هنگام Redis failure، نوع metric برای cache update failure و لاگ ساخت‌یافته برای تشخیص divergence بین Cache و Database.

## Review Example

فرض کنید در یک User Profile Service این الگوی ترافیک را داریم:

```
1000 Read/sec
20 Write/sec
```

**آیا Write Through انتخاب خوبی است؟**

بله، احتمالاً انتخاب مناسبی است، چون نسبت Read به Write بسیار بالاست. سیستم در هر ثانیه 1000 خواندن و فقط 20 نوشتن دارد، بنابراین هزینه اضافه‌ای که برای Update کردن Cache در هر Write پرداخت می‌شود، در برابر منفعتی که از Cache Hit های زیاد در Read ها به‌دست می‌آید قابل قبول است.

ن**سبت به Cache Aside چه مزیتی ایجاد می‌کند؟**

مزیت اصلی این است که بعد از هر Update، عملا Cache آماده و به‌روز است. در Cache Aside معمولاً بعد از Update، همیشه Cache حذف می‌شود و اولین Read بعدی باید Cache Miss بخورد، از Database بخواند و Cache را دوباره بسازد. اما در Write Through، Read بعدی مستقیماً Cache Hit می‌شود. بنابراین Hit Ratio بهتر می‌ماند، Cold Read کمتر می‌شود و Read Latency مخصوصاً برای داده‌های Hot کاهش پیدا می‌کند.

**اگر Redis Down شود ولی DB سالم باشد، رفتار سیستم باید چه باشد؟**

این یک تصمیم طراحی است و به حساسیت داده بستگی دارد. اگر Database منبع حقیقت باشد و User Profile بتواند برای مدت کوتاهی بدون Cache کار کند، معمولاً بهتر است Update در Database انجام شود و Cache Failure باعث fail شدن کل Write نشود. در این حالت باید خطای Cache ثبت شود، metric مناسب افزایش پیدا کند و با TTL، invalidate، retry محدود یا background repair وضعیت Cache اصلاح شود. اما اگر Freshness در Cache برای سیستم حیاتی باشد، ممکن است تصمیم بگیریم Write بدون موفقیت Cache هم fail شود. این تصمیم Availability را کاهش می‌دهد ولی Consistency بیشتری ایجاد می‌کند.

**آیا Availability را فدای Consistency می‌کنیم یا برعکس؟**

برای User Profile معمولاً Database منبع حقیقت است و Cache فقط برای Performance استفاده می‌شود؛ بنابراین معمولاً Availability را بیشتر حفظ می‌کنیم و اجازه می‌دهیم DB Update موفق شود، حتی اگر Cache موقتاً fail شود. سپس Cache با مکانیزم‌های جبرانی اصلاح می‌شود. اما برای داده‌های حساس‌تر، مثل موجودی مالی یا وضعیت‌هایی که Stale Read می‌تواند آسیب جدی ایجاد کند، ممکن است سیاست سخت‌گیرانه‌تری لازم باشد. در هر صورت، تصمیم درست باید بر اساس اثر Stale Data روی Business گرفته شود.

## جمع‌بندی و مدل نهایی

مدل نهایی Write Through این است: در Cache Aside هنگام Write معمولاً Cache را حذف می‌کنیم، اما در Write Through هنگام Write، عملا Cache را همزمان با Database به‌روزرسانی می‌کنیم. این الگو Read Performance را بهتر می‌کند، Hit Ratio را بالا نگه می‌دارد و احتمال Cold Read را کاهش می‌دهد، اما در سمت Write هزینه دارد: Write Latency بیشتر می‌شود، Failure Handling سخت‌تر می‌شود و اگر Write زیاد باشد، Write Amplification ایجاد می‌کند.

Write Through زمانی مناسب است که داده زیاد خوانده و کم نوشته شود، مثل Product Catalog، User Profile، Configuration و Feature Flags. اگر workload نوشتنی سنگین باشد، مثل Logs، Metrics، Events یا Order Streams، معمولاً Write Through انتخاب خوبی نیست. **جمله کلیدی این است: Write Through برای بهبود Read های آینده، هزینه بیشتری روی Write فعلی پرداخت می‌کند**.
