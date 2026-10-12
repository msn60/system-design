#  جمع‌بندی مهندسی Backend Performance

در سیستم‌های Backend مدرن، Performance صرفاً به معنی “سریع بودن” نیست. یک سیستم ممکن است: در شرایط عادی سریع باشد، اما تحت فشار collapse کند، latency غیرقابل پیش‌بینی تولید کند، یا در بار واقعی production ناپایدار شود. به همین دلیل، تعریف حرفه‌ای Performance در سیستم‌های production-scale شامل چند ویژگی کلیدی است:  Performance خوب یعنی:

## 1. قابل پیش‌بینی بودن (Predictability)

رفتار سیستم باید تحت workload های مختلف قابل پیش‌بینی باشد. سیستمی که:  یک بار 50ms و بار دیگر 5s پاسخ می‌دهد، حتی اگر average خوبی داشته باشد، از دید production system سالم محسوب نمی‌شود. Predictability مهم است زیرا: 


<ul dir="rtl">
  <li>timeoutها قابل تنظیم نیستند</li>
  <li>retry storm ایجاد می‌شود</li>
  <li>queue buildup رخ می‌دهد</li>
  <li>autoscaling تصمیم اشتباه می‌گیرد</li>
  <li>tail latency انفجاری می‌شود</li>
</ul>


در سیستم‌های distributed، نوسان latency اغلب خطرناک‌تر از latency متوسط بالا است. 

## 2. قابل مقیاس بودن (Scalability)

سیستم باید بتواند با رشد: تعداد کاربران، حجم داده، concurrency،  request rate همچنان behavior قابل قبولی حفظ کند. Scalability فقط به معنی “سرور بیشتر” نیست. بسیاری از bottleneckها با scale افقی حل نمی‌شوند:

<ul dir="rtl">
  <li>database contention</li>
  <li>distributed lock contention</li>
  <li>hot partitions</li>
  <li>queue saturation</li>
  <li>cache stampede</li>
  <li>fan-out explosion</li>
</ul>

به همین دلیل scalability باید در سطح:

<ul dir="rtl">
  <li>architecture</li>
  <li>data model</li>
  <li>caching strategy</li>
  <li>partitioning</li>
  <li>dependency management</li>
</ul>

طراحی شود.

## 3. قابل مشاهده بودن (Observability)

سیستمی که قابل مشاهده نباشد، قابل بهینه‌سازی نیست. در production، تنها دانستن اینکه latency بالا رفته، کافی نیست. باید بتوان تشخیص داد:

<ul dir="rtl">
  <li>کدام dependency کند شده</li>
  <li>queue کجا buildup کرده</li>
  <li>کدام query باعث tail latency شده</li>
  <li>آیا bottleneck از CPU است یا IO یا lock contention</li>
</ul>

 سه ستون اصلی Observability عبارتند از: 

### Metrics

برای: 

<ul dir="rtl">
  <li>latency</li>
  <li>throughput</li>
  <li>saturation</li>
  <li>error rate</li>
</ul>

### Logs

برای:

<ul dir="rtl">
  <li>debugging</li>
  <li>correlation</li>
  <li>incident analysis</li>
</ul>

### Distributed Tracing

برای:

<ul dir="rtl">
  <li>request path analysis</li>
  <li>identifying slow spans</li>
  <li>dependency timing breakdown</li>
</ul>

## 4. مقاوم بودن تحت فشار (Resilience Under Load)

سیستم production باید تحت load بالا degrade شود، نه collapse. مثلا در یک فروشگاه اینترنتی: recommendation system می‌تواند موقتاً غیرفعال شود، اما checkout نباید از کار بیفتد. 

## 5. کنترل‌شده بودن رفتار سیستم (Controlled Behavior)

یکی از مهم‌ترین تفاوت‌های سیستم mature و immature: سیستم mature تحت فشار، رفتار کنترل‌شده دارد، نه behavior تصادفی. 

مثال رفتار کنترل‌نشده: retry بدون محدودیت، queue نامحدود، goroutine explosion،  uncontrolled memory growth این‌ها معمولاً باعث: cascading failure - node instability - GC collapse - OOM kill می‌شوند.
