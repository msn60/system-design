#  جمع‌بندی مهندسی Backend Performance

در سیستم‌های Backend مدرن، Performance صرفاً به معنی “سریع بودن” نیست. یک سیستم ممکن است: در شرایط عادی سریع باشد، اما تحت فشار collapse کند، latency غیرقابل پیش‌بینی تولید کند، یا در بار واقعی production ناپایدار شود. به همین دلیل، تعریف حرفه‌ای Performance در سیستم‌های production-scale شامل چند ویژگی کلیدی است:  Performance خوب یعنی:

## 1. قابل پیش‌بینی بودن (Predictability)

رفتار سیستم باید تحت workload های مختلف قابل پیش‌بینی باشد. سیستمی که:  یک بار 50ms و بار دیگر 5s پاسخ می‌دهد، حتی اگر average خوبی داشته باشد، از دید production system سالم محسوب نمی‌شود. Predictability مهم است زیرا: 

- timeoutها قابل تنظیم نیستند
- retry storm ایجاد می‌شود
- queue buildup رخ می‌دهد
- autoscaling تصمیم اشتباه می‌گیرد
- tail latency انفجاری می‌شود

در سیستم‌های distributed، نوسان latency اغلب خطرناک‌تر از latency متوسط بالا است. 

## 2. قابل مقیاس بودن (Scalability)

سیستم باید بتواند با رشد: تعداد کاربران، حجم داده، concurrency،  request rate همچنان behavior قابل قبولی حفظ کند. Scalability فقط به معنی “سرور بیشتر” نیست. بسیاری از bottleneckها با scale افقی حل نمی‌شوند:

- database contention
- distributed lock contention
- hot partitions
- queue saturation
- cache stampede
- fan-out explosion

به همین دلیل scalability باید در سطح:

- architecture
- data model
- caching strategy
- partitioning
- dependency management

طراحی شود.

## 3. قابل مشاهده بودن (Observability)

سیستمی که قابل مشاهده نباشد، قابل بهینه‌سازی نیست. در production، تنها دانستن اینکه latency بالا رفته، کافی نیست. باید بتوان تشخیص داد:

- کدام dependency کند شده
- queue کجا buildup کرده
- کدام query باعث tail latency شده
- آیا bottleneck از CPU است یا IO یا lock contention

 سه ستون اصلی Observability عبارتند از: 

### Metrics

برای: 

- latency
- throughput
- saturation
- error rate

### Logs

برای:

- debugging
- correlation
- incident analysis

### Distributed Tracing

برای:

- request path analysis
- identifying slow spans
- dependency timing breakdown

## 4. مقاوم بودن تحت فشار (Resilience Under Load)

سیستم production باید تحت load بالا degrade شود، نه collapse. مثلا در یک فروشگاه اینترنتی: recommendation system می‌تواند موقتاً غیرفعال شود، اما checkout نباید از کار بیفتد. 

## 5. کنترل‌شده بودن رفتار سیستم (Controlled Behavior)

یکی از مهم‌ترین تفاوت‌های سیستم mature و immature: سیستم mature تحت فشار، رفتار کنترل‌شده دارد، نه behavior تصادفی. 

مثال رفتار کنترل‌نشده: retry بدون محدودیت، queue نامحدود، goroutine explosion،  uncontrolled memory growth این‌ها معمولاً باعث: cascading failure - node instability - GC collapse - OOM kill می‌شوند.
