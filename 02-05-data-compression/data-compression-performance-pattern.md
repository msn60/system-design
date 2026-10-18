# Data Compression — Performance Pattern

Data Compression به معنای کاهش حجم داده پیش از انتقال یا ذخیره‌سازی است تا Byte کمتری از Network یا Storage عبور کند. در یک API، جریان متداول عبارت است از `Original Data → Compress → Smaller Payload → Transfer → Decompress → Original Data`. برای مثال، یک JSON Response با حجم `500 KB` ممکن است پس از gzip به `80 KB` برسد. با این حال، Compression رایگان نیست: **Network Bandwidth، Transfer Time و گاهی Storage/Disk I/O کاهش می‌یابند، اما CPU Cost، Compression/Decompression Latency و پیچیدگی عملیاتی افزایش پیدا می‌کنند.** Compression یک Trade-off میان CPU و هزینه جابه‌جایی داده است. کاهش Memory Usage نیز تنها در شرایطی رخ می‌دهد که داده به‌شکل فشرده نگهداری شود یا از پردازش Streaming استفاده کنیم؛ Compress کردن Response به‌تنهایی الزاماً حافظه مصرفی Application را کم نمی‌کند.

## Mental Model

صد پتو را تصور کنید که باید با کامیون حمل شوند. حمل مستقیم آن‌ها فضای زیادی اشغال می‌کند، اما Vacuum Packing حجم بار را کم می‌کند؛ در مقابل، فشرده‌سازی و بازکردن بسته‌ها زمان و انرژی نیاز دارد. در Backend نیز `Compression → Less Data Movement + More CPU Work` است. اگر Network کند، محدود یا گران باشد، هزینه CPU ممکن است ارزش داشته باشد؛ اگر Payload کوچک یا CPU از قبل Bottleneck باشد، Compression حتی می‌تواند Performance را بدتر کند.

## چرا Compression یک Performance Pattern است؟

Bottleneck سیستم همیشه CPU نیست و ممکن است Network، Cross-region Bandwidth یا Cloud Egress باشد. فرض کنید سرویس روزانه `100 Million` Response با متوسط اندازه `100 KB` ارسال می‌کند. بدون Compression تقریباً `10 TB` داده منتقل می‌شود؛ اگر متوسط Payload فشرده `20 KB` شود، حجم انتقال تقریباً `2 TB` و صرفه‌جویی حدود `8 TB` خواهد بود. این کاهش می‌تواند روی Response Time، Egress Cost، Mobile Data Usage، Service-to-Service Traffic، Storage و Disk I/O اثر بگذارد. البته این یک برآورد حجمی است و به‌تنهایی به معنی کاهش متناظر در End-to-End Latency نیست؛ Bottleneck و CPU Cost باید اندازه‌گیری شوند.

## Compression Ratio و نوع داده

دو بیان رایج را باید از هم تفکیک کرد: `Compressed / Original` نسبت اندازه خروجی به ورودی است و `1 − (Compressed / Original)` درصد صرفه‌جویی را نشان می‌دهد. برای مثال، تبدیل `1 MB → 200 KB` به معنی نسبت اندازه حدود `0.20` و **۸۰٪ کاهش حجم** است. گاهی اصطلاح *Compression Ratio* را برعکس، یعنی `Original / Compressed` بیان می‌کنند؛ در این مثال `5:1`. هنگام گزارش Metric باید مشخص کنیم کدام تعریف را استفاده کرده‌ایم.

داده‌های دارای الگو و تکرار، معمولاً به‌خوبی Compress می‌شوند: `JSON`, `XML`, `HTML`, `CSS`, `JavaScript`, `CSV`, `Logs`, `Plain Text`, `SVG`. برای مثال، در هزار Object مشابه JSON نام فیلدهایی مانند `user_id`، `username` و `country` بارها تکرار می‌شود و Compressor می‌تواند از این تکرار استفاده کند:

```json
{
  "user_id": 1001,
  "username": "mehdi",
  "country": "Iran"
}
```

در مقابل، فرمت‌هایی مانند `JPEG`, `PNG`, `WebP`, `MP4`, `MP3`, `ZIP`, `gzip` و بسیاری از PDFها قبلاً فشرده شده‌اند. اجرای gzip روی JPEG معمولاً کاهش چشمگیری ایجاد نمی‌کند و گاهی Headerهای اضافی، خروجی را کمی بزرگ‌تر می‌کنند. اصل کاربردی: `Already Compressed Data → Usually Skip Recompression`. این یک قاعده تجربی است، نه ممنوعیت مطلق؛ ملاک نهایی اندازه واقعی خروجی و CPU Cost است.

## Small Payload Problem و Minimum Threshold

Compression برای Payloadهای بسیار کوچک معمولاً صرفه اقتصادی ندارد؛ برای مثال، یک Response صدبایتی ممکن است پس از اضافه‌شدن Header و Metadata الگوریتم، تقریباً همان اندازه یا حتی بزرگ‌تر شود و در عین حال CPU مصرف کند: `100 Bytes → Compression Work + Headers → Possibly 110 Bytes`. بنابراین بهتر است یک Minimum Size Threshold تعریف کنیم، مثلاً `1 KB` یا `2 KB`. این اعداد صرفاً نقطه شروع‌اند و مقدار مناسب باید با Benchmark برای Payload، الگوریتم، سطح Compression و Workload واقعی تعیین شود.

## Trade-off اصلی: CPU در برابر Bandwidth و Latency

فرض کنید `1 MB` داده به `100 KB` تبدیل شود. روی یک مسیر کند، سناریوی فرضی زیر به نفع Compression است:

```text
Without Compression: Transfer ≈ 100ms
With Compression:    Compress 5ms + Transfer 10ms + Decompress 2ms ≈ 17ms
```

اما روی یک شبکه داخلی سریع، با فرض `1ms` انتقال بدون Compression، همین هزینه پردازشی ممکن است وضعیت را بدتر کند:

```text
Without Compression: Transfer ≈ 1ms
With Compression:    Compress 5ms + Transfer 0.1ms + Decompress 2ms ≈ 7.1ms
```

این اعداد آموزشی‌اند و سایر هزینه‌ها را حذف کرده‌اند. نتیجه مهندسی: `Slow/Expensive Network → Compression Often Valuable`، اما `Fast Network + CPU-sensitive Workload → Compression May Hurt`. همچنین CPU زیاد می‌تواند Request Queue را طولانی‌تر کند و در نتیجه P95/P99 را افزایش دهد، حتی اگر حجم پاسخ کمتر شده باشد.

## Compression Level

الگوریتم‌هایی مانند gzip سطح Compression قابل‌تنظیم دارند. Level پایین معمولاً سریع‌تر و دارای خروجی بزرگ‌تر است؛ Level بالا CPU بیشتری مصرف می‌کند و غالباً خروجی کوچک‌تری تولید می‌کند، اما با *Diminishing Returns*. برای مثال، فرضی است که Level 1 خروجی `130 KB` و Level 9 خروجی `100 KB` بدهد، اما Level 9 چند برابر CPU مصرف کند. در Dynamic API Response معمولاً `DefaultCompression` یا سطح سریع‌تر مناسب‌تر از `BestCompression` است. افزایش Level همیشه به همان نسبت کاهش اندازه ایجاد نمی‌کند.

## gzip و HTTP Compression

`gzip` از متداول‌ترین گزینه‌ها برای HTTP است. Client توانایی خود را با Header زیر اعلام می‌کند:

```http
Accept-Encoding: gzip
```

Server در صورت انتخاب gzip، Response را فشرده و Header زیر را اضافه می‌کند:

```http
Content-Encoding: gzip
Vary: Accept-Encoding
```

Flow: `Client Supports gzip → Server Compresses → Network Transfer → Client Decompresses`. در Browserها Decompression معمولاً به‌صورت Transparent انجام می‌شود. Header `Vary: Accept-Encoding` برای Cache و CDN مهم است؛ زیرا نسخه فشرده و غیرفشرده یک Resource نباید با هم اشتباه شوند. در پیاده‌سازی واقعی، `Accept-Encoding` باید با Parsing درست Tokenها و کیفیت‌های `q` و نه صرفاً جست‌وجوی Substring بررسی شود.

### Response Compression

یکی از رایج‌ترین کاربردها، فشرده‌سازی پاسخ API است: `Database → 500 KB JSON → gzip → 70 KB → Client`. نمونه‌ها شامل Search Results، Product Catalog، Analytics Response، Reports و Configurationهای بزرگ‌اند. Response Compression برای Read-heavy API با Payloadهای متنی بزرگ معمولاً Candidate مناسبی است.

### Request Compression

در Bulk API، Telemetry، Log Upload یا Analytics می‌توان Request را هم فشرده کرد: `10 MB JSON → gzip → 1 MB → Server`. در این حالت Client از `Content-Encoding: gzip` استفاده می‌کند و Server مسیر `Receive → Decompress → Parse JSON` را اجرا می‌کند. **محدودیت Request باید هم روی حجم Compressed و هم روی حجم Decompressed اعمال شود**؛ صرف `Content-Length` یا محدودکردن Body قبل از Decompression کافی نیست.

## Brotli

`Brotli` معمولاً برای Text Payloadها Compression Ratio بسیار مناسبی دارد و در Static Web Assets مانند `HTML`, `CSS` و `JavaScript` پرکاربرد است. تفاوت کاربردی با gzip:

| ویژگی | gzip | Brotli |
|---|---|---|
| سازگاری و کاربرد | بسیار گسترده، رایج در API | رایج در وب و Static Assets |
| سرعت Compression | معمولاً مناسب Dynamic Response | در سطح‌های بالا ممکن است CPU-intensive باشد |
| نسبت فشرده‌سازی متن | خوب | اغلب بهتر |
| کاربرد پیشنهادی | Dynamic APIs و سازگاری عمومی | Precompressed Static Assets و موارد Benchmark‌شده |

Brotli هم می‌تواند برای Dynamic Response استفاده شود، اما انتخاب باید براساس Benchmark باشد. مزیت بزرگ آن در Static Assets این است که هزینه Compression یک‌بار در Build پرداخت می‌شود:

```text
Build Time → app.js → app.js.gz / app.js.br
Request → Choose Variant Based on Accept-Encoding → Serve Existing File
```

## Static Compression در برابر Dynamic Compression

**Static Compression:** فایل یک‌بار در Build Pipeline فشرده و هزاران بار استفاده می‌شود: `Original Asset → Build-time Compress → Store .gz/.br → Serve`. در این مدل CPU Cost آنلاین فشرده‌سازی تقریباً حذف می‌شود و برای فایل‌های ثابت انتخاب مناسبی است.

**Dynamic Compression:** پاسخ در زمان اجرای Request تولید می‌شود: `DB Query → Build JSON → Compress JSON → Send`. چون Response ممکن است برای هر Request متفاوت باشد، هزینه CPU، Allocation و احتمالاً Buffering در Runtime پرداخت می‌شود. برای نمونه `GET /users/123/orders` ممکن است به Dynamic Compression نیاز داشته باشد.

## HTTP Content Negotiation و Cache Variants

Client ممکن است چند الگوریتم را اعلام کند:

```http
Accept-Encoding: gzip, br
```

Server یکی را انتخاب می‌کند و با `Content-Encoding` برمی‌گرداند. `Vary: Accept-Encoding` باید در طراحی CDN/Proxy/Cache در نظر گرفته شود تا Cache نسخه مناسب را ارائه کند. هنگام تولید Variantهای مختلف، `Content-Length` باید با Representation فشرده‌شده سازگار باشد یا در حالت Streaming حذف شود. `ETag` و سایر Metadataهای مرتبط با Representation نیز باید با Encoding و رفتار Cache هماهنگ باشند. همچنین Responseهایی که از قبل `Content-Encoding` دارند نباید ناخواسته دوباره فشرده شوند.

## Streaming Compression و Buffering Trade-off

روش پرمصرف از نظر Memory این است: `Generate 500 MB → Buffer Entire Result → Compress → Send`. روش مناسب‌تر برای خروجی‌های بزرگ: `Generate Chunk → Compress Chunk → Send → Repeat`. Streaming Compression می‌تواند Memory Usage و Time to First Byte را کاهش دهد، اما Compressor ممکن است برای Compression Ratio بهتر مقداری داده Buffer کند. اگر Eventها با تأخیر Flush شوند، Client داده را دیرتر دریافت می‌کند: `Event Generated → Compression Buffer → Delayed Delivery`. در SSE، Live Updates و Streamingهای Latency-sensitive باید Flush Policy، Compression Overhead و در صورت نیاز غیرفعال‌کردن Compression بررسی شوند. Streaming به معنی حذف کامل Buffering نیست.

## gRPC Compression و تفاوت Serialization با Compression

جریان مفهومی در gRPC: `Protobuf Message → Optional Message Compression → HTTP/2 → Decompression`. **Protobuf فشرده‌سازی نیست**؛ یک Binary Serialization Format است. معمولاً از JSON کوچک‌تر است، اما Messageهای Protobuf با Patternهای تکراری هنوز می‌توانند Compress شوند. بنابراین `Serialization Format ≠ Compression` و حتی در بعضی داده‌ها اندازه `JSON + gzip` ممکن است به Raw Protobuf نزدیک شود یا از آن کوچک‌تر باشد؛ CPU و Latency هر دو باید اندازه‌گیری شوند.

Compression برای تمام RPCها مناسب نیست. اگر Request حدود `200 Bytes`، Response حدود `500 Bytes` و نرخ `100K RPC/s` روی شبکه داخلی سریع باشد، CPU Compression ممکن است از صرفه‌جویی Network بیشتر شود. اما برای Responseهای `100 KB`، `1 MB` یا `10 MB` ممکن است تصمیم متفاوتی بگیریم. در gRPC باید پشتیبانی Client/Server، الگوریتم، Message Size و هزینه هر فراخوانی سنجیده شود؛ Policy بر اساس Payload Size مفید است.

## Compression در Message Broker و ترکیب با Batching

Producer می‌تواند چند Message را Batch و سپس فشرده کند: `Messages → Batch → Compress → Broker`. فشرده‌سازی یک Batch در بسیاری از Workloadها بهتر از Compression جداگانه هر Message است، زیرا Patternهای تکراری بین Messageها نیز قابل استفاده می‌شوند:

```json
{"user_id":1,"event":"click"}
{"user_id":2,"event":"click"}
{"user_id":3,"event":"click"}
```

در این مثال، Field Nameها و مقدار `click` تکرار شده‌اند. بنابراین `Batching + Compression` می‌تواند Network و Storage Efficiency را بهبود دهد. در مقابل، زمان تشکیل Batch، CPU و Latency Messageها باید اندازه‌گیری شوند.

## Compression، Encryption و TLS

در حالت معمول ترتیب فنی این است: `Original Data → Compress → Encrypt → Transfer`؛ زیرا Ciphertext معمولاً مانند داده تصادفی به نظر می‌رسد و دیگر به‌خوبی Compress نمی‌شود. در HTTPS مسیر معمول `Application Response → HTTP Compression → TLS Encryption → Network` است. بااین‌حال، **Compression پیش از Encryption همیشه بدون ریسک امنیتی نیست**: اگر Secret و داده قابل‌کنترل توسط مهاجم در یک Response فشرده شوند، تغییر اندازه خروجی می‌تواند اطلاعاتی درباره Secret آشکار کند (مانند BREACH). بنابراین در Responseهای حساس باید Security Review انجام شود و در صورت نیاز Compression محدود یا غیرفعال شود.

## Compression Bomb و محدودیت Decompression

Payload فشرده ممکن است کوچک، اما پس از Decompression بسیار بزرگ باشد؛ برای مثال `1 MB Compressed → 10 GB Decompressed`. نتیجه می‌تواند CPU/Memory/Disk Exhaustion باشد. در Production باید محدودیت‌های جداگانه تعریف شوند:

```text
Maximum Compressed Request = 10 MB
Maximum Decompressed Request = 100 MB
Decompression Deadline / CPU Budget = Bounded
```

این اعداد نمونه‌اند. اجرای Decompression باید در حالت Streaming و با محدودیت حجم خروجی انجام شود؛ `Content-Length` فقط اندازه Representation روی Wire را گزارش می‌کند. Error نامعتبر بودن Encoding، Payload بیش از حد بزرگ و Timeout نیز باید واضح و کنترل‌شده مدیریت شوند.

## Compression و Cache

دو رویکرد اصلی برای Cache Representation وجود دارد: `Store Uncompressed → Compress on Every Response` یا `Store Compressed → Send Directly`. روش دوم CPU per Request را کاهش می‌دهد، اما ممکن است به Variantهای جداگانه مانند `resource`, `resource:gzip`, `resource:br` نیاز باشد. CDNها معمولاً انتخاب Variant و `Accept-Encoding` را مدیریت می‌کنند. در Cache Key و Invalidation باید مشخص باشد چه Representationی ذخیره شده است.

## Double Compression و محل مسئولیت

Anti-pattern مهم، فشرده‌سازی تکراری در چند Layer است:

```text
Application Compress → Reverse Proxy Compress Again
Already Compressed Response → CDN Recompress
```

این کار CPU مصرف می‌کند، معمولاً Size را چندان کاهش نمی‌دهد و ممکن است Compatibility/Encoding را خراب کند. در معماری `Go Service → Nginx/API Gateway → CDN → Client` باید مالک Compression مشخص باشد. گاهی بهتر است Application خروجی غیرفشرده تولید کند و Proxy/CDN مسئول Public HTTP Compression باشد؛ در ارتباط داخلی Service-to-Service ممکن است Application یا gRPC مستقیماً مسئول شود. **به‌جای Compression در تمام Layerها، یک Policy آگاه از Topology تعیین می‌کنیم.**

## Go: پیاده‌سازی آموزشی gzip Response Middleware

نمونه زیر منطق پایه Wrap کردن Response Writer را نشان می‌دهد؛ **Production-ready نیست** و عمداً مشکلات آن را پس از کد بررسی می‌کنیم.

```go
package main

import (
    "compress/gzip"
    "io"
    "net/http"
    "strings"
)

type gzipResponseWriter struct {
    http.ResponseWriter
    io.Writer
}

func (w *gzipResponseWriter) Write(p []byte) (int, error) {
    return w.Writer.Write(p)
}

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Add("Vary", "Accept-Encoding")

        gz := gzip.NewWriter(w)
        defer gz.Close()

        next.ServeHTTP(&gzipResponseWriter{
            ResponseWriter: w,
            Writer:         gz,
        }, r)
    })
}
```

این نسخه برای نمایش مفهوم `Wrap Writer → Compress Response` است؛ در Production مشکلات مهم دارد: (۱) تشخیص `Accept-Encoding` با `strings.Contains` کافی نیست؛ (۲) Threshold ندارد و JSON کوچک، JPEG، ZIP یا Streaming Endpoint را هم Compress می‌کند؛ (۳) ممکن است Handler بعداً `Content-Length` یا `Content-Encoding` ناسازگار تنظیم کند؛ (۴) باید وضعیت‌های بدون Body، `HEAD`، Range Response و Content-Type را مدیریت کند؛ (۵) ممکن است قابلیت‌های اختیاری `http.ResponseWriter` مانند `http.Flusher` را پنهان کند؛ (۶) Error نهایی `Close/Flush` را مدیریت نمی‌کند؛ (۷) اگر Header زود نوشته شود، تصمیم بر اساس اندازه و Content-Type واقعی دشوار می‌شود. برای همین در سیستم واقعی، استفاده از Middleware آزموده‌شده در Proxy/Framework یا پیاده‌سازی دقیق با Testهای مناسب معمولاً منطقی‌تر است.

### Content-Type و Compression Policy

برای داده‌های متنی نظیر `application/json`, `application/xml`, `text/html`, `text/plain`, `text/css`, `application/javascript`, `image/svg+xml` معمولاً Compression مفید است. برای `image/jpeg`, `image/png`, `image/webp`, `video/*`, `audio/*`, `application/zip`, `application/gzip` معمولاً باید Skip شود. علاوه بر Content-Type، `Minimum Payload Size`، `Already Compressed`، `Streaming`، `Status Code` و `Client Capability` نیز باید بررسی شوند.

### Compression Level در Go

```go
func newGzipWriter(w io.Writer, fast bool) (*gzip.Writer, error) {
    level := gzip.DefaultCompression
    if fast {
        level = gzip.BestSpeed
    }
    return gzip.NewWriterLevel(w, level)
}
```

در این نمونه، برای هر Response تنها یک سطح Compression انتخاب می‌شود. `gzip.BestSpeed` معمولاً سریع‌تر و دارای خروجی بزرگ‌تر است، در حالی که `gzip.DefaultCompression` تعادل عمومی ارائه می‌دهد. `gzip.BestCompression` نیز در Go وجود دارد، اما برای Dynamic API معمولاً فقط در صورت اثبات سود با Benchmark منطقی است.

### Buffer Pool در Go

در API با QPS بالا، Allocationهای `gzip.Writer` و Bufferهای همراه آن ممکن است روی GC Pressure اثر بگذارند. پس از Profile کردن می‌توان از `sync.Pool` استفاده کرد: `Request → Get gzip.Writer → Reset(output) → Compress → Close/Flush → Reset/Clear References → Put Back`. اما `sync.Pool` تضمین نمی‌کند Object همیشه باقی بماند و Reuse نادرست Writer می‌تواند Data Corruption یا نگهداری Reference ناخواسته ایجاد کند. ترتیب صحیح تصمیم: `Measure Allocations → Identify Hot Path → Consider Pooling`، نه Pool کردن پیش‌فرض همه چیز.

## Precompressed Static Assets

برای فایل ثابت، یک Build Pipeline می‌تواند هم‌زمان `app.js`, `app.js.gz`, `app.js.br` را تولید کند. Request دارای `Accept-Encoding: br, gzip` به Variant سازگار هدایت می‌شود؛ در غیر این صورت فایل عادی ارسال می‌شود. مزیت اصلی: `CPU Cost Paid Once → Reused for Many Requests`. Metadata کش، Content-Encoding، Content-Type و Content-Length برای هر Variant باید صحیح تنظیم شود.

## Compression Threshold و Adaptive Compression

یک Policy ابتدایی: `Payload < 1 KB → Skip` و `Payload ≥ 1 KB → Consider Compression`. Threshold جهانی وجود ندارد و به Payload، نوع داده، CPU، Network، الگوریتم و Level بستگی دارد. Policy پیشرفته‌تر می‌تواند چنین باشد:

```text
Small JSON → No Compression
Medium JSON → gzip BestSpeed
Large JSON → gzip DefaultCompression
Static Assets → Precompressed Brotli
Already Compressed Media → Skip
```

در سیستم‌های بسیار بزرگ، Load-aware Compression هم ممکن است مفید باشد: `CPU Saturated → Lower Level/Skip` یا `Bandwidth Saturated → Consider Stronger Compression`. این قابلیت پیچیدگی و نوسان رفتاری ایجاد می‌کند و فقط با Metrics و Guardrail توجیه دارد.

## Latency SLO و Throughput

برای Endpoint با هدف فرضی `P99 < 50ms`:

```text
Without Compression: Generate 10ms + Transfer 30ms = 40ms
With Compression:    Generate 10ms + Compress 8ms + Transfer 5ms + Decompress 2ms = 25ms
```

اما روی مسیر سریع:

```text
Without Compression: Generate 10ms + Transfer 1ms = 11ms
With Compression:    Generate 10ms + Compress 8ms + Transfer 1ms ≈ 19ms
```

این اعداد صرفاً برای فهم Trade-off هستند. Compression در Network-bound System با کاهش Bytes per Request می‌تواند Throughput را افزایش دهد؛ در CPU-bound System ممکن است CPU Saturation ایجاد کند و Requests/sec را کاهش دهد. Compression و Connection Reuse مکمل‌اند: `Keep-Alive/HTTP2 → Lower Connection Overhead` در برابر `Compression → Lower Payload Bytes`.

## Production Example: Product Search API

فرض کنید `20K req/s` با متوسط JSON Response برابر `200 KB` داریم و Network Bottleneck اصلی است. حجم ترافیک تقریبی بدون Compression: `20K × 200 KB ≈ 4 GB/s`؛ با فشرده‌شدن Response به `40 KB`، حجم انتقال حدود `800 MB/s` می‌شود. Policy نمونه: `application/json + Payload > 2 KB + Client Supports gzip → Compress`. باید اثر آن را بر CPU، P95/P99، Network Bytes و Throughput اندازه‌گیری کرد؛ اگر CPU اشباع شود، ممکن است Bottleneck فقط از Network به CPU منتقل شده باشد.

## Production Example: Internal gRPC Service

فرض کنید `Request = 300 Bytes`، `Response = 500 Bytes`، `100K RPC/s` و شبکه داخلی `10Gbps` است. Compression احتمالاً سود اندکی دارد و CPU آن می‌تواند هزینه بیشتری تحمیل کند. اگر Response به `500 KB` افزایش یابد، Compression باید دوباره Benchmark شود؛ اندازه Payload و Bottleneck واقعی مبنای تصمیم‌اند، نه صرف استفاده از gRPC یا Protobuf.

## Production Example: Telemetry Upload

Deviceها هر سی ثانیه یک Batch به حجم `2 MB JSON` ارسال می‌کنند و gzip آن را به `200 KB` می‌رساند. مزایا: `Mobile Bandwidth ↓ + Upload Time ↓ + Server Ingress ↓`. اما Server باید `Compressed Size Limit`, `Decompressed Size Limit`, `CPU Budget`, `Timeout` و Error Handling مناسب داشته باشد.

## Failure Modes و Anti-patternها

- **CPU Saturation:** فعال‌کردن Compression روی همه Responseها → CPU Usage، Queueing Delay و P99 بالا می‌رود. راهکار: Threshold، Level پایین‌تر، Skip Payload کوچک و در صورت تناسب Offload به Proxy/CDN.
- **Already Compressed Data:** اجرای gzip روی JPEG/MP4 یا ZIP اغلب CPU Waste است.
- **Small Payload Compression:** Header/CPU Cost ممکن است از Network Saving بیشتر شود.
- **Double Compression:** Application، Proxy و CDN نباید بدون قرارداد روشن روی یک Representation دوباره Compression انجام دهند.
- **Compression Bomb:** Payload فشرده کوچک می‌تواند پس از Decompression بسیار بزرگ شود. محدودیت روی هر دو حجم و Timeout لازم است.
- **Missing `Vary`:** Cache/CDN ممکن است Encoding نامناسب را به Client بدهد.
- **Incorrect `Content-Length`:** اندازه Uncompressed نباید بعد از Compression همان مقدار باقی بماند؛ مقدار صحیح Compressed تعیین یا Header حذف شود.
- **Buffering/Delayed Flush:** Responseهای Real-time ممکن است به‌دلیل Buffering Compressor دیر برسند.
- **Unnecessary Highest Level:** انتخاب بالاترین Compression Level برای همه Dynamic APIها ممکن است هزینه CPU زیادی و سود حجمی کمی داشته باشد.
- **Compression Without Benchmark:** کاهش Byte الزاماً بهبود Latency یا Throughput را تضمین نمی‌کند.

## Observability

Metricهای پیشنهادی:

```text
compression_requests_total
compression_skipped_total
compression_input_bytes_total
compression_output_bytes_total
compression_duration_seconds
decompression_duration_seconds
compressed_response_total
decompression_failure_total
compression_ratio
```

همراه با Metricهای سیستمی: `CPU Usage`, `CPU Time per Request`, `P95/P99 Latency`, `Network Bytes Sent/Received`, `Response Size`, `Throughput`. برای نسبت اندازه `Compressed / Original`، مقدار `0.20` یعنی خروجی ۲۰٪ اندازه اصلی و صرفه‌جویی ۸۰٪ است؛ مقدار نزدیک `0.95` نشان‌دهنده صرفه‌جویی حجمی کم است. بهتر است Compression Ratio، Latency و CPU به تفکیک Algorithm، Content Type و Size Bucket مقایسه شوند.

**Cardinality:** به‌جای URL کامل یا User ID در Metric Label از `content_type="json"`, `algorithm="gzip"`, `size_bucket="100kb-1mb"` استفاده شود. Labelهایی مثل `url="/users/83910471/orders"` ممکن است Cardinality را بسیار بالا ببرند.

**Alertهای مهم:** افزایش Compression CPU Time، کاهش سود Compression Ratio، بدترشدن P99 پس از فعال‌سازی، کاهش نیافتن Network Bytes، افزایش Decompression Failure و CPU Saturation. مثلاً اگر CPU حدود ۳۰٪ بیشتر شود، اما Network Bytes تنها ۲٪ کاهش یابد، Compression Policy باید بازنگری شود.

## روش عملی تصمیم‌گیری

ابتدا مشخص می‌کنیم: آیا داده Compressible است؟ Payload چقدر بزرگ است؟ Network Bottleneck است یا CPU؟ آیا داده از قبل فشرده شده؟ Latency مهم‌تر است یا Bandwidth؟ Response Static است یا Dynamic؟ تصمیم نمونه:

```text
Small Payload? → Skip
Already Compressed? → Usually Skip
Text / JSON / XML? → Good Candidate
Static Asset? → Precompress
Large Dynamic Payload? → Benchmark gzip/Brotli
CPU-bound? → Fast Level or Skip
```

## Review Scenario: Product API

فرض کنید:

```text
Traffic = 30K req/s
Average JSON Response = 150 KB
P99 = 120ms
CPU Usage = 35%
Network Link Utilization = 80%
```

**تحلیل:** Network نزدیک اشباع است و CPU هنوز Headroom دارد؛ بنابراین Compression Candidate مناسبی است. فرض کنیم gzip مقدار `150 KB → 30 KB` بدهد و هزینه Compression حدود `2ms` باشد. حجم Network حدود ۸۰٪ کم می‌شود و ممکن است Network Queueing، Transfer Time و P99 کاهش یابند. اگر CPU مثلاً به `55%` برسد و P99 بهتر شود، تصمیم احتمالاً مطلوب است؛ اما اگر CPU به `95%` برسد، Bottleneck از Network به CPU منتقل شده است. اقدامات اصلاحی: `Lower Compression Level → Increase Threshold → Offload Compression → Consider CPU Scaling`. اگر Response ترکیبی از JSON Metadata و JPEG Thumbnail باشد، JSON را فشرده می‌کنیم، اما JPEG را دوباره Compress نمی‌کنیم. در حضور CDN نیز بهتر است محل مسئولیت Compression دوباره بررسی شود.

نکته اصلی این سناریو آن است که **کاهش حجم خروجی کافی نیست؛ بهبود واقعی باید در CPU، Network، Throughput و P99 تأیید شود.**

## نکات Production-oriented

در Production باید پاسخ این پرسش‌ها روشن باشد: چه Content Typeهایی قابل Compression هستند؟ Minimum Size Threshold چقدر است؟ Algorithm و Level چیست؟ کدام Layer مالک Compression است؟ فایل‌های Static چگونه Precompress می‌شوند؟ Request Decompression چه محدودیت حجمی و زمانی دارد؟ Double Compression چگونه جلوگیری می‌شود؟ CPU و P99 چگونه اندازه‌گیری می‌شوند؟ Responseهای Streaming، Range و حساس امنیتی چه Policy دارند؟ به‌جای یک Boolean عمومی `compression=true` بهتر است Policy شامل `Algorithm`, `Threshold`, `Level`, `Content Types`, `Max Decompressed Size` و Exceptionها وجود داشته باشد.

## نکات Interview-oriented

در مصاحبه، پس از پیشنهاد Compression باید بلافاصله Trade-off مطرح شود:

```text
Benefits: Lower Bandwidth, Lower Transfer Time, Lower Egress Cost
Costs:    Higher CPU, Possible Added Latency, Operational Complexity
```

همچنین باید به عدم Compression مجدد Media فشرده، Threshold برای Payload کوچک، Precompression فایل ثابت، Benchmark برای Dynamic Response، ترتیب معمول Compress-before-Encrypt، محدودیت Decompression و نقش Cache/Proxy/CDN اشاره شود. بیان اینکه «Compression همیشه Performance را بهتر می‌کند» نادرست است.

## جمع‌بندی و Mental Model نهایی

Data Compression یک Performance Pattern برای کاهش حجم داده در انتقال یا ذخیره‌سازی است: `Large Payload → Compress → Transfer Fewer Bytes → Decompress`. مدل تصمیم‌گیری: `Large + Compressible + Network-bound → Compression Usually Valuable` و `Small + Already Compressed + CPU-bound → Compression Often Not Valuable`. مدل Production: `Classify Payload → Apply Threshold → Select Algorithm/Level → Assign Ownership → Avoid Double Compression → Bound Decompression → Measure CPU, Network and P99 → Tune`.

**Mental Model نهایی:** Compression یعنی **CPU مصرف کنیم تا Byte کمتری جابه‌جا شود**؛ این تصمیم زمانی مفید است که صرفه‌جویی واقعی در Network، Storage یا I/O از هزینه Compress/Decompress و اثر آن بر Latency و Throughput بیشتر باشد.
