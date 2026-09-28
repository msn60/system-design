# SLI / SLO / SLA 

در طراحی سیستم‌های قابل‌اعتماد (Reliability Engineering / SRE)، سه مفهوم کلیدی وجود دارد:

- SLI = Service Level Indicator
- SLO = Service Level Objective
- SLA = Service Level Agreement

این سه مفهوم پایه اصلی:

- Reliability Engineering
- Monitoring
- Alerting
- Capacity Planning
- Incident Management
- Production Operations

هستند.

## 1) SLI چیست؟

### تعریف دقیق

SLI یک معیار کمی و قابل اندازه‌گیری است که کیفیت واقعی تجربه کاربران از یک سرویس را نشان می‌دهد.

نکته مهم:هر `Metric` ای `SLI` نیست، اما هر `SLI` یک `Metric` است`.`

### تفاوت Metric و SLI

سیستم‌ها هزاران metric دارند:

- CPU usage
- memory usage
- goroutine count
- disk IO
- TCP connections
- queue size

اما همه این‌ها SLI نیستند. چرا؟ **چون SLI باید مستقیماً کیفیت قابل مشاهده برای کاربر را نشان دهد**. مثال:

**کاربر این موارد را حس می‌کند**:

- آیا درخواست سریع پاسخ داده شد؟
- آیا درخواست fail شد؟
- آیا سرویس available بود؟
- آیا داده تازه بود؟

پس این‌ها SLI های خوبی هستند:

- **request latency**
- **availability**
- **success rate**
- **freshness**

اما:

- CPU usage
- heap size
- disk usage

معمولاً operational metrics هستند، نه SLI. 

### Mental Model مهم

به این شکل نگاه کن:

- Metrics:  تمام داده‌های سیستم
- SLIs: مهم‌ترین metricهایی که کیفیت واقعی سرویس را نشان می‌دهند
- SLOs: هدف مورد انتظار برای آن کیفیت
- SLAs: تعهد تجاری/حقوقی بر اساس آن کیفیت

## 2) چرا SLI مهم است؟

بدون SLI:

- نمی‌توان فهمید سیستم واقعاً سالم است یا نه
- نمی‌توان کیفیت سرویس را اندازه‌گیری کرد
- نمی‌توان reliability را مدیریت کرد
- alerting معنی‌دار ممکن نیست
- capacity planning دشوار می‌شود
- SLO و SLA قابل تعریف نیستند

در واقع: You cannot improve what you cannot measure. SLI .در واقع SLI اطلاعات فنی، واقعی و measurable از رفتار سیستم ارائه می‌دهد.

## 3) انواع رایج SLI ها

### 3.1) Latency SLIs

این‌ها رایج‌ترین SLI ها هستند. 

- **P50 latency**
- **P95 latency**
- **P99 latency**
- **TTFB**
- **request duration**

چرا Percentile مهم است؟ Average latency معمولاً گمراه‌کننده است. مثال:

- 99 request = 10ms
- 1 request = 10s

میانگین هنوز خوب به نظر می‌رسد، اما یک کاربر تجربه فاجعه‌باری داشته است.به همین دلیل در production: 

- P95
- P99
- Tail Latency

بسیار مهم‌تر از average هستند. مثال SLI :  می خواهیم 95% از request ها باید زیر 200ms پاسخ داده شوند.

### 3.2) Availability SLIs

Availability فقط به معنی up بودن سرور نیست. ممکن است:  سرویس پاسخ دهد  ولی بسیار کند باشد. از دید کاربر، این هنوز failure است. معمولاً availability بر اساس:Good Requests / Total Valid Requests محاسبه می‌شود.

مثال Good Request

- request is good if:
- \- status < 500
- AND
- \- latency < 300ms

**مثال SLI – در این سرویس 99.9% از request ها باید successful و زیر 300ms باشند.**

### 3.3) Error Rate SLIs

این SLI نشان می‌دهد چه درصدی از requestها fail شده‌اند. مثال:

- HTTP 5xx ratio
- gRPC error ratio
- timeout rate

مثال: Error Rate < 0.1%

### 3.4) Throughput SLIs

Throughput همیشه SLI نیست. گاهی فقط capacity metric است. اما اگر مستقیماً تجربه کاربر را تحت تاثیر قرار دهد، می‌تواند SLI باشد. مثال:

- message processing rate
- streaming throughput
- event ingestion rate

### 3.5) Freshness SLIs

در سیستم‌های event-driven و streaming بسیار مهم است. مثال:

```
     replication lag
     Kafka consumer lag
     data staleness
```

مثال: 95% از eventها باید زیر 5 ثانیه پردازش شوند.

## 4) SLO چیست؟

SLO هدفی است که برای SLI تعیین می‌کنیم. مثال SLI:  P95 latency = 180 ms و SLO: P95 latency باید کمتر از 200ms باشد. 

**نکته مهم**: SLI = چیزی که اندازه می‌گیریم و SLO = هدفی که می‌خواهیم به آن برسیم

### Error Budget

 از مهم‌ترین مفاهیم SRE است. اگر SLO برابر باشد با: 99.9% availability یعنی: 0.1% failure قابل قبول است این مقدار Error Budget نام دارد.به این دلیل  Error Budget مهم است چون بین: feature velocity و reliability تعادل ایجاد می‌کند.

مثال واقعی - اگر error budget تمام شود:  deploy متوقف می‌شود. release freeze انجام می‌شود. تمرکز روی reliability می‌رود. این در Google SRE بسیار مهم است.

## 5) SLA چیست؟

**SLA یک تعهد رسمی و تجاری بر اساس SLO است**. معمولاً شامل:

- target ها
- جریمه مالی
- نحوه اندازه‌گیری
- maintenance window
- exclusions

مثال:  99.95% monthly availability , otherwise the customer receives service credits.

تفاوت اصلی

- SLI → Measurement
- SLO → Engineering Target
- SLA → Business Commitment

## 6) مثال واقعی در یک User Management Service

فرض کن سرویسی داری: GET /api/v1/users/{id}

**SLI های مناسب**

- User-facing SLIs
- Latency
- P95 request latency
- P99 request latency
- Availability
- successful requests ratio
- Error Rate
- HTTP 5xx rate
- timeout rate
- Freshness
- profile update propagation delay

```
Internal Operational Metrics
```

این‌ها مهم‌اند ولی الزاماً SLI نیستند:

- PostgreSQL query latency
- Redis hit ratio
- goroutine count
- GC pause duration
- queue depth

این distinction مهم است.

## 7) تاثیر حجم دیتا روی SLI

وقتی می‌گویند: حجم دیتا روی SLI تاثیر دارد، منظور این نیست که: حجم دیتا خودش SLI است بلکه: workload characteristics روی behavior سیستم اثر می‌گذارند و در نهایت SLI ها degrade می‌شوند. سه نوع مهم حجم دیتا به شرح زیر می‌باشند: 

### 7.1) Payload Size

مثال:

- request body بزرگ
- response body بزرگ

اثر:

- network latency
- serialization cost
- GC pressure
- memory allocation

### 7.2) Dataset Size

مثال:

- رشد جدول users
- رشد Kafka topic
- رشد Elasticsearch index

اثر:

- cache miss
- query slowdown
- replication lag
- disk IO increase

### 7.3) Per-Request Workload Size

مثال:

- تعداد rows scanned
- fan-out calls
- aggregation size

اثر:

- CPU usage
- DB load
- queue delay

## 8) Tail Latency و Data Size

یکی از مهم‌ترین نکات production این است: بزرگ شدن دیتا معمولاً اول P99 را خراب می‌کند. نه average را. چون:

- slow requests
- GC pauses
- disk seeks
- queue contention

در percentile های بالا ظاهر می‌شوند.

### 9) ارتباط با Go

در Go افزایش حجم دیتا مستقیماً روی این‌ها اثر می‌گذارد:

- allocations
- heap growth
- GC overhead
- serialization cost
- copy amplification

مثال خطرناک: json.Marshal(hugeStruct)  روی payload های بزرگ:

- allocation spike
- CPU spike
- P99 degradation

ایجاد می‌کند.

راهکارهای رایج Production

- streaming
- chunked response
- protobuf
- compression
- object storage
- pagination
- cache
- buffer pooling
- zero-copy optimization

## 10) اشتباه رایج تیم‌ها

### اشتباه اول

هر metric را SLI در نظر می‌گیرند. در نتیجه:

- dashboard های noisy
- alert fatigue
- unclear priorities

### اشتباه دوم

فقط average latency را نگاه می‌کنند. در حالی که: Tail Latency is what kills distributed systems.

### اشتباه سوم

SLI ها را بدون context workload بررسی می‌کنند. مثلاً: P95 = 80ms بدون دانستن:

- payload size
- concurrency
- RPS
- dataset size

تقریباً بی‌معناست.

## 11) Production-Oriented Insight

بسیاری از سیستم‌ها در staging خوب کار می‌کنند، اما در production با dataset واقعی collapse می‌کنند. مثال: staging DB = 100K  rows و  production DB = 2B rows. نتیجه:

- query plan change
- cache inefficiency
- replication lag
- P99 explosion

## 12) Interview-Oriented Insight

اگر در مصاحبه پرسیدند: What makes a good SLI? پاسخ mature این است:

```
A good SLI should:
- reflect user experience
- be measurable
- correlate with reliability
- be actionable
- avoid excessive noise
```
