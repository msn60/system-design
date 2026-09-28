# Latency (مهم‌ترین مفهوم بک‌اندی)

Latency یک عدد واحد نیست. اجزای Latency در Backend عموما به صورت زیر می‌باشد: 

```
Total Response Time = Network + Load Balancer + Auth + Application Logic + DB + Cache + External Services
```

 در Design/System Doc باید مشخص باشد:

- کجاها وقت تلف می شود؟
- کدام بخش بیشترین سهم را دارد؟ کدام بخش bottleneck است؟
- کدام قابل بهینه‌سازی است؟ در‌واقع کدام بخش دست ماست و می‌توانیم روی آن تأثیر بگذاریم
- کدام غیرقابل کنترل است؟ کدام دست ما نیست

## Response Time، Latency و TTFB چه فرقی دارند؟

این سه مفهوم شبیه هم‌اند ولی دقیقاً یکی نیستند.

### Latency

Latency معمولاً یعنی زمان رفت‌وبرگشت request تا response. در بک‌اند وقتی می‌گوییم latency، اغلب منظورمان این است:

- Client sends request  
- Server receives and processes
- Server sends response
- Client receives response

یعنی کل زمانی که request طول می‌کشد. Latency می‌تواند از زاویه‌های مختلف اندازه‌گیری شود، به عنوان مثال:

- Client-side latency: کل زمان از دید کلاینت
- Server-side latency: زمانی که request داخل سرور پردازش شده
- DB latency: زمان query دیتابیس
- Network latency: زمان شبکه

پس Latency یک واژه عمومی‌تر است. مثلاً می‌گویی: DB latency is high یا Network latency is high یا API latency P95 is 400ms. این  یعنی latency همیشه الزاماً فقط کل response نیست؛ می‌تواند latency یک بخش از سیستم هم باشد.

Latency یک مفهوم عمومی برای اندازه‌گیری تاخیر در انجام یک عملیات است. این عملیات می‌تواند مربوط به کل API، شبکه، دیتابیس، کش، صف، سرویس داخلی یا یک dependency خارجی باشد.  بنابراین Latency فقط به زمان شروع پاسخ یا شروع پردازش محدود نمی‌شود، بلکه می‌تواند در هر نقطه‌ای از مسیر درخواست اندازه‌گیری شود. نمونه‌هایی از Latency عبارت‌اند از:

- Network Latency
- Database Latency
- Redis Latency
- gRPC Latency
- Queue Latency
- API Latency

به بیان دقیق‌تر:  Latency نشان‌دهنده‌ی مدت‌زمان تاخیر در یک بخش مشخص از مسیر پردازش request است. این بخش می‌تواند شبکه، صف، اپلیکیشن، دیتابیس، کش، سرویس خارجی یا کل API باشد.

#### Network Latency

Network Latency مربوط به تاخیر های شبکه‌ای قبل از رسیدن درخواست به اپلیکیشن یا هنگام برگشت پاسخ به کلاینت است. مهم‌ترین پارامترهای آن عبارت‌اند از:

- RTT (Round Trip Time)
- One-way Latency
- DNS Resolution Time
- TCP Connection Time
- TLS Handshake Time

برای مثال:

- DNS Lookup:       20ms
- TCP Connect:      40ms
- TLS Handshake:    70ms
- Request Send:      5ms

این زمان‌ها هنوز وارد business logic نشده‌اند، اما مستقیماً روی تجربه‌ی کاربر اثر می‌گذارند.

نکته‌ی مهم این است که این نوع latency معمولاً از سمت client، gateway، proxy، CDN یا ابزارهای synthetic monitoring قابل مشاهده است و در backend monitoring خام همیشه به‌صورت مستقیم دیده نمی‌شود.

#### Application Latency

Application Latency مربوط به تاخیر هایی است که داخل اپلیکیشن رخ می‌دهد؛ یعنی بعد از اینکه request به سرویس رسیده، اما هنوز response آماده نشده است. نمونه‌هایی از Application Latency: 

- Queue Time
- Handler Dispatch Time
- Middleware Time
- Business Logic Time
- Serialization Time
- Time Before First Write

برای مثال، ممکن است request وارد سرویس شده باشد اما به دلیل فشار بالا، کمبود worker، پر بودن goroutine pool یا محدودیت connection pool، مدتی در صف منتظر بماند. نمونه:

- Total API Latency: 800ms
- DB Latency:         50ms
- Business Logic:    100ms
- Queue Time:        600ms

در این مثال، مشکل اصلی دیتابیس نیست؛ بلکه request قبل از شروع پردازش مدت زیادی در صف مانده است.

#### Internal / Dependency Latency

بخش مهم دیگری از Latency مربوط به dependency های داخلی یا خارجی سرویس است. نمونه‌ها: 

- Database Call Latency
- Cache Latency
- RPC / gRPC Latency
- Message Broker Latency
- External Service Latency

برای مثال: 

- PostgreSQL Query Latency: 80ms
- Redis Latency:             5ms
- Payment RPC Latency:     300ms
- Kafka Produce Latency:    20ms

در تحلیل عملکرد سیستم، بررسی این بخش‌ها کمک می‌کند مشخص شود bottleneck اصلی کجاست: دیتابیس، کش، سرویس خارجی، صف پیام یا خود application logic.

### Response Time

Response Time معمولاً از دید کاربر یا کلاینت اندازه‌گیری می‌شود. یعنی: از لحظه ارسال request تا کامل دریافت شدن response (از وقتی کلاینت request را می‌فرستد تا وقتی کل response را دریافت می‌کند)

برای APIهای معمولی، latency و response time خیلی وقت‌ها تقریباً به‌جای هم استفاده می‌شوند. اما در سیستم‌های بزرگ‌تر ممکن است دقیق‌تر تفکیک شوند.

مثلاً:

- Server-side latency: 120ms
- Network overhead:    40ms
- Client parsing:      20ms
- Total response time: 180ms

Response Time کل زمان request/response از دید client است؛ از ارسال request تا دریافت آخرین بایت پاسخ. برای APIهای معمولی، Response Time معمولاً مهم‌ترین معیار user-facing است. مانند:

- POST /login

```
P95 Response Time < 300ms
```

- POST /payment/authorize

```
P95 Response Time < 700ms
```

- GET /feed

```
P95 Response Time < 800ms
```

در این نوع APIها، معمولاً response کوچک است و کاربر یا کلاینت منتظر دریافت کامل نتیجه است؛ بنابراین Response Time معیار اصلی تجربه‌ی کاربر محسوب می‌شود.

#### Total Response Time

Total Response Time معیار پایه برای سنجش کامل شدن پاسخ است، اما باید مشخص شود از دید کدام نقطه اندازه‌گیری می‌شود:

- Client-observed Response Time
- Gateway-observed Response Time
- Server-side Response Time

این سه مقدار می‌توانند با هم متفاوت باشند. مثال:

- Backend Server Time:       180ms
- Gateway Observed Time:     220ms
- Client Observed Time:      350ms

تفاوت این اعداد می‌تواند ناشی از network latency، TLS overhead، download time، CDN، proxy یا موقعیت جغرافیایی client باشد.

#### Download Time

Download Time فاصله‌ی بین دریافت اولین بایت و دریافت آخرین بایت response است. رابطه‌ی ساده: 

```
Download Time ≈ TTLB - TTFB
```

نکته در مورد TTLB:  مخفف Time To Last Byte است و نشان می‌دهد آخرین بایت response چه زمانی به client رسیده است. از دید client، TTLB معمولاً معادل عملی Full Response Time است.

- TTFB = زمان تا شروع response
- TTLB = زمان تا پایان response
- Response Time ≈ TTLB از دید client

مقایسه‌ی TTFB و TTLB برای تحلیل عملکرد بسیار مفید است.  مثال: TTFB = 150ms و TTLB = 2150ms پس 

```
Download Time ≈ 2000ms
```

برداشت این است که سرور نسبتاً زود شروع به پاسخ‌دهی کرده، اما دریافت کامل payload زمان زیادی برده است. دلیل آن می‌تواند بزرگ بودن response، کندی شبکه یا محدودیت client باشد. Download Time در موارد زیر اهمیت بیشتری دارد: 

- File Download
- Image / Video Delivery
- Large JSON Payload
- CSV Export
- Report Generation
- Streaming

####  Server Processing Time

Server Processing Time بخشی از Response Time است که داخل backend صرف پردازش request می‌شود. این زمان می‌تواند شامل موارد زیر باشد:

- CPU Processing Time
- Database Query Time
- External RPC Time
- Internal Function Execution Time
- Template Rendering Time
- Business Logic Time
- Serialization Time

نکته‌ی مهم این است که Server Processing Time فقط یکی از اجزای Response Time است. Response Time از دید client علاوه بر پردازش سرور، شامل network، gateway، download time و سایر overhead ها نیز می‌شود.

### Time To First Byte یا TTFB

TTFB مخفف Time To First Byte است و به زمانی گفته می‌شود که از لحظه‌ی ارسال request تا دریافت اولین بایت response توسط client طول می‌کشد. به بیان ساده: TTFB نشان می‌دهد سیستم چه زمانی شروع به پاسخ‌دهی کرده است.  نکته‌ی مهم این است که TTFB الزاماً فقط server processing time نیست. به عنوان نمونه این می تواند یک flow ساده ارسال درخواست باشد:

- Client sends request
- Network
- Backend starts processing
- Backend begins sending response
- Client receives first byte  ← TTFB
- Client receives full response ← Full Response Time

این به این معنی است که  سرور سریع شروع کرده جواب بدهد، ولی کامل شدن response زمان بیشتری گرفته است. مثلاً برای دانلود فایل، streaming API، یا صفحه وب، TTFB مهم است.  فرض کنید Instagram یک feed برمی‌گرداند. اگر backend بتواند زود اولین بخش response را بفرستد، کلاینت می‌تواند سریع‌تر شروع به render کردن کند، حتی اگر کل payload هنوز کامل نیامده باشد. مثال: TTFB: 150ms و Full response completed: 900ms یعنی کاربر یا کلاینت خیلی زود فهمیده که server زنده است و دارد جواب می‌دهد.

اگر TTFB از دید client یا browser اندازه‌گیری شود، می‌تواند شامل موارد زیر باشد:

- DNS Resolution
- TCP Connection
- TLS Handshake
- Request Network Time
- Load Balancer / Gateway Overhead
- Queue Time
- Server Processing Before First Write
- Serialization Before First Write
- Network Return Time

برای مثال:

- DNS:                       20ms
- TCP:                       40ms
- TLS:                       60ms
- Request Network Time:      20ms
- Gateway Overhead:          10ms
- Backend Queue Time:        50ms
- Backend Processing:       120ms
- First Byte Return Time:    20ms

```
TTFB = 340ms
```

بنابراین اگر TTFB بالا باشد، الزاماً نمی‌توان نتیجه گرفت که فقط backend کند است. مشکل ممکن است از DNS، TLS، شبکه، gateway، صف، دیتابیس، serialization یا application logic باشد.

#### کاربردهای اصلی TTFB

TTFB مخصوصاً در سناریوهایی اهمیت دارد که شروع پاسخ برای تجربه‌ی کاربر مهم‌تر از کامل شدن کل پاسخ است. نمونه‌ها:

- Web Pages
- Streaming APIs
- File Downloads
- Large Responses
- Server-Sent Events
- gRPC Streaming
- Video / Audio Streaming
- LLM-like Token Streaming

#### Upstream First Byte Time

در سیستم‌هایی که چند hop دارند، فقط TTFB نهایی کافی نیست. باید بررسی شود هر hop چه زمانی اولین بایت را از upstream خود دریافت کرده است. برای مثال: Client → CDN → API Gateway → Service A → Service B . در چنین ساختاری ممکن است چند نوع first byte time داشته باشیم:

- Service B → Service A First Byte Time
- Service A → Gateway First Byte Time
- Gateway → Client First Byte Time

این متریک‌ها مخصوصاً در API Gateway ها و reverse proxy هایی مانند Nginx، Envoy و Traefik بسیار مفید هستند. برای مثال: 

- Gateway TTFB to Client:       700ms
- Upstream First Byte Time:     650ms
- Gateway Overhead:              50ms

در این حالت، مشکل احتمالاً در upstream service است، نه خود gateway. 

#### Serialization Latency و اثر آن بر TTFB

Serialization Latency زمانی است که سیستم برای آماده‌سازی اولین خروجی قابل ارسال صرف می‌کند؛ مثلاً تبدیل struct به JSON.اگر backend مجبور باشد کل response را در memory بسازد و سپس آن را serialize کند، TTFB ممکن است بالا برود.

طراحی اول:

- Fetch all data
- Build full object
- Serialize full JSON
- Send response

در این حالت، اولین بایت زمانی ارسال می‌شود که کل response آماده شده باشد.

طراحی دوم:

- Start response early
- Send chunks gradually

در این حالت، TTFB کاهش پیدا می‌کند، هرچند ممکن است full response time همچنان بالا باشد.

### جمع‌بندی:

اگر بخواهیم تفاوت اصلی این سه را در یک جمله خلاصه کنیم: 

- Latency: تاخیر در یک مسیر یا بخش از سیستم، مفهوم عمومی‌تر (برای تحلیل تاخیر در هر بخش از سیستم استفاده می‌شود)
- Response Time: کل زمان کامل شدن پاسخ از دید کلاینت(برای فهمیدن اینکه کل پاسخ چه زمانی کامل می‌شود استفاده می‌شود)
- TTFB: زمان رسیدن اولین byte پاسخ به کلاینت (برای فهمیدن اینکه پاسخ‌دهی چه زمانی شروع می‌شود استفاده می‌شود)

فرض کن این API را داریم:

```
GET /feed
```

زمان‌ها این شکلی‌اند:

```
Client sends request:              0ms
Network to backend:               20ms
Load balancer:                     5ms
Auth:                             15ms
App logic:                        60ms
DB query:                        100ms
Backend starts response:         200ms
First byte reaches client:        220ms
Response body download:           300ms
Client receives full response:    520ms
```

حالا:

```
TTFB = 220ms
Response Time = 520ms
```

و latency بسته به context می‌تواند این‌ها باشد:

```
Network latency = 20ms رفت + برگشت/بخش شبکه
DB latency = 100ms
Server-side latency = حدود 180ms تا 200ms
Client-observed latency = 520ms
```

پس وقتی کسی فقط می‌گوید latency، باید بپرسیم: Latency از دید کجا؟ کلاینت؟ سرور؟ دیتابیس؟ شبکه؟

### کی از کدام استفاده کنیم؟

وقتی درباره تجربه کاربر یا API حرف می‌زنیم: Response Time  مثلاً در Design Doc بنویسی: 

P95 response time for GET /feed should be under 500ms.

این برای user-facing APIها خیلی مناسب است. مثال:

- Login API response time P95 < 300ms
- Checkout API response time P95 < 700ms
- Search API response time P95 < 500ms

اینجا مهم است کاربر چقدر منتظر می‌ماند.

وقتی می خواهیم  bottleneck پیدا می‌کنیم: Latency 

Latency برای تحلیل داخلی سیستم بهتر است. در واقع برای debugging و bottleneck analysis استفاده می‌شود. مثلاً:

- DB latency increased from 20ms to 200ms.
- Redis latency is stable.
- External payment provider latency is high.
- Network latency between services increased.

اینجا نمی‌گویید response time، چون دارید درباره یک بخش خاص حرف می‌زنیم.

مثلاً در debugging:

- Total response time is 800ms.
- DB latency is 500ms.
- So the bottleneck is probably database access.

وقتی response بزرگ، دانلود فایل، streaming یا web page داری: TTFB

TTFB مخصوصاً وقتی مهم است که:

- response بزرگ است
- فایل دانلود می‌شود
- HTML page render می‌شود
- API streaming داری
- Server-Sent Events داری
- gRPC streaming داری

کاربر باید زود بفهمد سیستم شروع به پاسخ کرده است. مثلاً: 

TTFB should be under 200ms for the homepage.

برای سایت‌ها، TTFB خیلی مهم است چون browser تا اولین byte را نگیرد، نمی‌تواند شروع خوبی برای render داشته باشد. برای streaming هم مهم است. مثلاً در ChatGPT-like response:

- User sends prompt
- TTFB = time until first token appears
- Full response time = time until whole answer is generated

اگر اولین token بعد از ۵۰۰ms بیاید، کاربر حس می‌کند سیستم responsive است، حتی اگر کل جواب ۱۰ ثانیه طول بکشد.

نکته‌ی کلیدی این است که TTFB بخشی از Response Time است؛ زیرا Response Time از شروع request تا دریافت آخرین بایت ادامه دارد، در حالی که TTFB فقط تا دریافت اولین بایت را اندازه‌گیری می‌کند. 

همچنین Latency مفهومی عمومی‌تر از هر دو است و می‌تواند برای کل request یا هر بخش داخلی سیستم، مانند شبکه، دیتابیس، کش، RPC، صف و business logic استفاده شود.

بنابراین در یک Design Doc یا System Design Discussion، بهتر است معیارها به‌صورت دقیق و قابل اندازه‌گیری بیان شوند. برای مثال:

- For read APIs, P95 response time should be under 300ms.
- For payment authorization, P99 response time should be under 1s.
- The main latency contributors should be tracked separately: DB latency, cache latency, RPC latency, and queue time.
- For streaming or large responses, TTFB should be under 200ms and TTLB should be monitored separ
