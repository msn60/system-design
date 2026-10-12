# Performance Pattern 1: Connection Pooling

Connection Pooling یعنی به‌جای اینکه برای هر درخواست جدید یک connection تازه به دیتابیس، Redis، سرویس خارجی یا message broker ساخته شود، تعدادی connection از قبل ایجاد و مدیریت می‌شود و request ها از همان connection های موجود استفاده می‌کنند. هدف اصلی:

- کاهش هزینه ساخت connection
- کنترل تعداد connection های همزمان
- افزایش throughput
- کاهش latency
- جلوگیری از overload شدن dependency

## مسئله‌ای که Connection Pooling حل می‌کند

ساخت connection معمولاً ارزان نیست. برای مثال اتصال به PostgreSQL ممکن است شامل این موارد باشد: 

- TCP handshake
- TLS handshake
- Authentication
- Session initialization
- Connection state setup

اگر برای هر request یک connection جدید ساخته شود، سیستم هم کند می‌شود و هم dependency را تحت فشار قرار می‌دهد. 

مثلا یک الگوی اشتباه به این صورت است:

```
HTTP Request
    → Create DB Connection
    → Execute Query
    → Close Connection
```

در load بالا این رفتار می‌تواند باعث موارد زیر شود: 

- افزایش شدید latency
- افزایش مصرف CPU روی Database
- افزایش تعداد connection های همزمان
- timeout شدن request ها
- کاهش throughput
- ایجاد connection exhaustion

## ایده اصلی Connection Pool

در Connection Pooling، تعدادی connection از قبل ایجاد می‌شوند و در یک pool نگه‌داری می‌شوند. زمانی که request جدیدی وارد می‌شود، به‌جای ساخت connection جدید، یکی از connections های آزاد pool استفاده می‌شود. پس از پایان کار نیز connection بسته نمی‌شود، بلکه دوباره به pool بازمی‌گردد تا requestهای بعدی از آن استفاده کنند. رفتار کلی سیستم به این شکل است: 

```
Request
    → دریافت connection از pool
    → اجرای query یا request
    → بازگرداندن connection به pool
```

این کار باعث می‌شود: 

- هزینه connection setup حذف شود
- response سریع‌تر شود
- فشار روی dependency کاهش پیدا کند
- performance پایدارتر شود

## Pool Size و اهمیت آن

یکی از مهم‌ترین بخش‌های Connection Pooling، تنظیم اندازه pool است. 

اگر pool بیش از حد کوچک باشد: 

- request ها منتظر connection می‌مانند
- queue تشکیل می‌شود
- p95 و p99 latency افزایش پیدا می‌کند

اگر pool بیش از حد بزرگ باشد:

- Database overload می‌شود
- memory usage بالا می‌رود
- lock contention بیشتر می‌شود
- performance کلی سیستم افت می‌کند

بنابراین هدف این نیست که pool را تا جای ممکن بزرگ کنیم؛ هدف این است که اندازه‌ای انتخاب شود که با ظرفیت واقعی dependency هماهنگ باشد. 

## Monitoring و Observability

Connection Pool بدون observability تقریباً غیرقابل مدیریت است. مهم‌ترین metric هایی که باید monitor شوند:

- تعداد connection های فعال
- تعداد connection های idle
- wait count
- wait duration
- pool saturation
- DB connection errors

در Go، پکیج \`database/sql\` این اطلاعات را ارائه می‌دهد: stats := db.Stats() و از طریق آن می‌توان bottleneckهای pool را تشخیص داد. 

##  Timeout و Connection Pool

یکی از اشتباهات رایج این است که request ها بدون timeout منتظر connection بمانند. در سیستم‌های production، معمولاً request نباید زمان زیادی برای گرفتن connection صبر کند. در غیر این صورت: 

- queue buildup رخ می‌دهد
- go routine ها زیاد می‌شوند
- memory pressure ایجاد می‌شود
- retry storm ممکن است شروع شود

به همین دلیل معمولاً query ها همراه با context timeout اجرا می‌شوند. 

## اشتباهات رایج در Connection Pooling

### ساخت connection در هر request

این یکی از رایج‌ترین anti-pattern ها است و باعث می‌شود عملاً pooling بی‌اثر شود. 

### Pool بسیار بزرگ

بزرگ بودن بیش از حد pool معمولاً Database را به bottleneck تبدیل می‌کند.

### Transactions های طولانی

اگر transaction مدت زیادی باز بماند، connection برای مدت طولانی در اختیار همان request باقی می‌ماند و pool سریع‌تر اشباع می‌شود.

### نداشتن timeout

بدون timeout، request ها می‌توانند مدت زیادی منتظر connection یا query باقی بمانند و باعث overload شدن کل سیستم شوند.

## ارتباط Connection Pooling با Performance و Reliability

Connection Pooling فقط یک optimization ساده نیست، بلکه یکی از مکانیزم‌های اصلی کنترل فشار روی dependency ها است. اگر به‌درستی تنظیم شود:

- latency کاهش پیدا می‌کند
- throughput افزایش پیدا می‌کند
- tail latency پایدارتر می‌شود
- dependency ها محافظت می‌شوند

اما اگر اشتباه تنظیم شود، می‌تواند باعث:

- queue buildup
- timeout storm
- retry storm
- database saturation
- cascading failure

شود.

## جمع‌بندی

Connection Pooling یکی از بنیادی‌ترین الگوهای Performance در Backend است که با reuse کردن connectionها، هزینه ساخت connection جدید را کاهش می‌دهد و concurrency را کنترل می‌کند. 

در سیستم‌های production، تنظیم صحیح pool، timeout، transaction duration و observability اهمیت بسیار زیادی دارند، زیرا Connection Pool مستقیماً روی: 

- latency
- throughput
- p95/p99
- error rate
- system stability

اثر می‌گذارد. 
