# چک‌لیست تحلیل Performance برای هر Endpoint

برای هر endpoint باید بتوانید مسیر کامل request را تحلیل کنید. هدف این نیست که فقط: response time = 700ms را ببینید. هدف این است که بدانید: **این 700ms دقیقاً کجا مصرف شده است**.

## سوال‌های کلیدی برای هر Endpoint

### 1. TTFB چقدر است؟

TTFB (Time To First Byte) نشان می‌دهد: سیستم چقدر طول می‌کشد تا اولین بخش response را ارسال کند. TTFB بالا معمولاً نشان‌دهنده یکی از این موارد است:

- queue delay
- DB wait
- dependency latency
- connection pool exhaustion
- slow serialization

مثلا  اگر: TTFB = 900ms و Total Response Time = 950ms باشد، مشکل معمولاً قبل از شروع response است. اما اگر: TTFB = 50ms و Total Response Time = 3s احتمالاً:

- payload بزرگ است
- streaming کند است
- client/network کند است

## 2. Response Time p95/p99 چقدر است؟

Average latency تقریباً همیشه ناکافی است. در distributed systems: همواره tail latency تعیین‌کننده تجربه واقعی کاربران است. زیرا اگر 99% requests خوب باشند، 1% بدترین request ها می‌توانند کل سیستم را degrade کنند.

### دلایل رایج P99 بالا

- GC pauses
- lock contention
- slow queries
- dependency spikes
- queue buildup
- noisy neighbors
- network retransmission

## 3. آیا Queue وجود دارد؟

تقریباً هر queue می‌تواند: buffer مفید یا محل پنهان شدن latency باشد. Queue خوب باید: 

- bounded باشد
- observable باشد
- backpressure ایجاد کند

Queue خطرناک:

- Queue نامحدود:  

```
traffic spike → queue growth → memory growth → latency explosion → timeout storm 
```

## 4. آیا Database Bottleneck است؟

Database معمولاً اولین bottleneck واقعی سیستم‌های backend است. علائم DB Bottleneck عبارتند از: 

- query latency بالا
- lock wait
- connection pool exhaustion
- replication lag
- sequential scans
- high IO wait

سوال‌های مهم که می توانیم برای دیتابیس از خود بپرسیم: 

- آیا query ها index مناسب دارند؟
- آیا pagination وجود دارد؟
- آیا N+1 query رخ می‌دهد؟
- آیا fan-out query داریم؟
- آیا hot rows وجود دارد؟

## 5. آیا External Dependency وجود دارد؟

هر dependency خارجی: 

- latency
- failure
- retry
- unpredictability

را وارد سیستم می‌کند. یک قانون مهم این است که هر dependency خارجی،SLO شما را ضعیف‌تر می‌کند. به عنوان مثال اگر endpoint شما: auth service و payment gateway و recommendation service را call کند، latency نهایی مجموع behavior همه آن‌هاست.

## 6. سیستم Under Load چگونه رفتار می‌کند؟

بسیاری از سیستم‌ها در load کم سالم‌اند، اما تحت pressure collapse می‌کنند. سوال‌های مهم که در این مورد باید از خود بپرسیم: 

- آیا queue build up می‌شود؟
- آیا retry storm رخ می‌دهد؟
- آیا auto scaling به‌موقع عمل می‌کند؟
- آیا connection pool اشباع می‌شود؟
- آیا tail latency انفجاری می‌شود؟
- آیا backpressure وجود دارد؟
