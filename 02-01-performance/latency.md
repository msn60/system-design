<div dir="rtl" align="right">

<h1 dir="rtl" align="right">تاخیر - Latency (مهم‌ترین مفهوم بک‌اندی)</h1>

<p dir="rtl" align="right">در نظر داشته باشید که Latency یک عدد واحد نیست. اجزای Latency در Backend عموما به صورت زیر می‌باشد:</p>

</div>

```
Total Response Time = Network + Load Balancer + Auth + Application Logic + DB + Cache + External Services
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">در Design/System Doc باید مشخص باشد:</p>

<ul dir="rtl">
  <li>کجاها وقت تلف می شود؟</li>
  <li>کدام بخش بیشترین سهم را دارد؟ کدام بخش bottleneck است؟</li>
  <li>کدام قابل بهینه‌سازی است؟ در‌واقع کدام بخش دست ماست و می‌توانیم روی آن تأثیر بگذاریم</li>
  <li>کدام غیرقابل کنترل است؟ کدام دست ما نیست</li>
</ul>

<h2 dir="rtl" align="right">مقدارهای Response Time، Latency و TTFB چه فرقی دارند؟</h2>

<p dir="rtl" align="right">این سه مفهوم شبیه هم‌اند ولی دقیقاً یکی نیستند.</p>

<h3 dir="rtl" align="right">Latency</h3>

<p dir="rtl" align="right">عموما Latency یعنی زمان رفت‌وبرگشت request تا response. در بک‌اند وقتی می‌گوییم latency، اغلب منظورمان این است:</p>

<ul dir="rtl">
  <li>Client sends request</li>
  <li>Server receives and processes</li>
  <li>Server sends response</li>
  <li>Client receives response</li>
</ul>

<p dir="rtl" align="right">یعنی کل زمانی که request طول می‌کشد. Latency می‌تواند از زاویه‌های مختلف اندازه‌گیری شود، به عنوان مثال:</p>

<ul dir="rtl">
  <li>Client-side latency: کل زمان از دید کلاینت</li>
  <li>Server-side latency: زمانی که request داخل سرور پردازش شده</li>
  <li>DB latency: زمان query دیتابیس</li>
  <li>Network latency: زمان شبکه</li>
</ul>

<p dir="rtl" align="right">پس Latency یک واژه عمومی‌تر است. مثلاً می گوییم: DB latency is high یا Network latency is high یا API latency P95 is 400ms. این  یعنی latency همیشه الزاماً فقط کل response نیست؛ می‌تواند latency یک بخش از سیستم هم باشد.</p>

<p dir="rtl" align="right">عموما Latency یک مفهوم عمومی برای اندازه‌گیری تاخیر در انجام یک عملیات است. این عملیات می‌تواند مربوط به کل API، شبکه، دیتابیس، کش، صف، سرویس داخلی یا یک dependency خارجی باشد.  بنابراین Latency فقط به زمان شروع پاسخ یا شروع پردازش محدود نمی‌شود، بلکه می‌تواند در هر نقطه‌ای از مسیر درخواست اندازه‌گیری شود. نمونه‌هایی از Latency عبارت‌اند از:</p>

<ul dir="rtl">
  <li>Network Latency</li>
  <li>Database Latency</li>
  <li>Redis Latency</li>
  <li>gRPC Latency</li>
  <li>Queue Latency</li>
  <li>API Latency</li>
</ul>

<p dir="rtl" align="right">به بیان دقیق‌تر:  Latency نشان‌دهنده‌ی مدت‌زمان تاخیر در یک بخش مشخص از مسیر پردازش request است. این بخش می‌تواند شبکه، صف، اپلیکیشن، دیتابیس، کش، سرویس خارجی یا کل API باشد.</p>

<h4 dir="rtl" align="right">Network Latency</h4>

<p dir="rtl" align="right">همواره Network Latency مربوط به تاخیر های شبکه‌ای قبل از رسیدن درخواست به اپلیکیشن یا هنگام برگشت پاسخ به کلاینت است. مهم‌ترین پارامترهای آن عبارت‌اند از:</p>

<ul dir="rtl">
  <li>RTT (Round Trip Time)</li>
  <li>One-way Latency</li>
  <li>DNS Resolution Time</li>
  <li>TCP Connection Time</li>
  <li>TLS Handshake Time</li>
</ul>

<p dir="rtl" align="right">برای مثال:</p>

<ul dir="rtl">
  <li>DNS Lookup:       20ms</li>
  <li>TCP Connect:      40ms</li>
  <li>TLS Handshake:    70ms</li>
  <li>Request Send:      5ms</li>
</ul>

<p dir="rtl" align="right">این زمان‌ها هنوز وارد business logic نشده‌اند، اما مستقیماً روی تجربه‌ی کاربر اثر می‌گذارند.</p>

<p dir="rtl" align="right">نکته‌ی مهم این است که این نوع latency معمولاً از سمت client، gateway، proxy، CDN یا ابزارهای synthetic monitoring قابل مشاهده است و در backend monitoring خام همیشه به‌صورت مستقیم دیده نمی‌شود.</p>

<h4 dir="rtl" align="right">Application Latency</h4>

<p dir="rtl" align="right">تاخیر اپلیکیشن یا Application Latency مربوط به تاخیر هایی است که داخل اپلیکیشن رخ می‌دهد؛ یعنی بعد از اینکه request به سرویس رسیده، اما هنوز response آماده نشده است. نمونه‌هایی از Application Latency&rlm;:</p>

<ul dir="rtl">
  <li>Queue Time</li>
  <li>Handler Dispatch Time</li>
  <li>Middleware Time</li>
  <li>Business Logic Time</li>
  <li>Serialization Time</li>
  <li>Time Before First Write</li>
</ul>

<p dir="rtl" align="right">برای مثال، ممکن است request وارد سرویس شده باشد اما به دلیل فشار بالا، کمبود worker، پر بودن goroutine pool یا محدودیت connection pool، مدتی در صف منتظر بماند. نمونه:</p>

<ul dir="rtl">
  <li>Total API Latency: 800ms</li>
  <li>DB Latency:         50ms</li>
  <li>Business Logic:    100ms</li>
  <li>Queue Time:        600ms</li>
</ul>

<p dir="rtl" align="right">در این مثال، مشکل اصلی دیتابیس نیست؛ بلکه request قبل از شروع پردازش مدت زیادی در صف مانده است.</p>

<h4 dir="rtl" align="right">Internal / Dependency Latency</h4>

<p dir="rtl" align="right">بخش مهم دیگری از Latency مربوط به dependency های داخلی یا خارجی سرویس است. نمونه‌ها:</p>

<ul dir="rtl">
  <li>Database Call Latency</li>
  <li>Cache Latency</li>
  <li>RPC / gRPC Latency</li>
  <li>Message Broker Latency</li>
  <li>External Service Latency</li>
</ul>

<p dir="rtl" align="right">برای مثال:</p>

<ul dir="rtl">
  <li>PostgreSQL Query Latency: 80ms</li>
  <li>Redis Latency:             5ms</li>
  <li>Payment RPC Latency:     300ms</li>
  <li>Kafka Produce Latency:    20ms</li>
</ul>

<p dir="rtl" align="right">در تحلیل عملکرد سیستم، بررسی این بخش‌ها کمک می‌کند مشخص شود bottleneck اصلی کجاست: دیتابیس، کش، سرویس خارجی، صف پیام یا خود application logic&rlm;.</p>

<h3 dir="rtl" align="right">Response Time</h3>

<p dir="rtl" align="right">مدت زمان پاسخ یا Response Time معمولاً از دید کاربر یا کلاینت اندازه‌گیری می‌شود. یعنی: از لحظه ارسال request تا کامل دریافت شدن response (از وقتی کلاینت request را می‌فرستد تا وقتی کل response را دریافت می‌کند)</p>

<p dir="rtl" align="right">برای APIهای معمولی، latency و response time خیلی وقت‌ها تقریباً به‌جای هم استفاده می‌شوند. اما در سیستم‌های بزرگ‌تر ممکن است دقیق‌تر تفکیک شوند.</p>

<p dir="rtl" align="right">مثلاً:</p>

<ul dir="rtl">
  <li>Server-side latency: 120ms</li>
  <li>Network overhead:    40ms</li>
  <li>Client parsing:      20ms</li>
  <li>Total response time: 180ms</li>
</ul>

<p dir="rtl" align="right">بنابراین Response Time کل زمان request/response از دید client است؛ از ارسال request تا دریافت آخرین بایت پاسخ. برای APIهای معمولی، Response Time معمولاً مهم‌ترین معیار user-facing است. مانند:</p>

<ul dir="rtl">
  <li>POST /login</li>
</ul>

</div>

```
P95 Response Time < 300ms
```

<div dir="rtl" align="right">

<ul dir="rtl">
  <li>POST /payment/authorize</li>
</ul>

</div>

```
P95 Response Time < 700ms
```

<div dir="rtl" align="right">

<ul dir="rtl">
  <li>GET /feed</li>
</ul>

</div>

```
P95 Response Time < 800ms
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">در این نوع APIها، معمولاً response کوچک است و کاربر یا کلاینت منتظر دریافت کامل نتیجه است؛ بنابراین Response Time معیار اصلی تجربه‌ی کاربر محسوب می‌شود.</p>

<h4 dir="rtl" align="right">Total Response Time</h4>

<p dir="rtl" align="right">کل زمان پاسخ دهی Total Response Time معیار پایه برای سنجش کامل شدن پاسخ است، اما باید مشخص شود از دید کدام نقطه اندازه‌گیری می‌شود:</p>

<ul dir="rtl">
  <li>Client-observed Response Time</li>
  <li>Gateway-observed Response Time</li>
  <li>Server-side Response Time</li>
</ul>

<p dir="rtl" align="right">این سه مقدار می‌توانند با هم متفاوت باشند. مثال:</p>

<ul dir="rtl">
  <li>Backend Server Time:       180ms</li>
  <li>Gateway Observed Time:     220ms</li>
  <li>Client Observed Time:      350ms</li>
</ul>

<p dir="rtl" align="right">تفاوت این اعداد می‌تواند ناشی از network latency، TLS overhead، download time، CDN، proxy یا موقعیت جغرافیایی client باشد.</p>

<h4 dir="rtl" align="right">Download Time</h4>

<p dir="rtl" align="right">زمان دانلود یا همان Download Time فاصله‌ی بین دریافت اولین بایت و دریافت آخرین بایت response است. رابطه‌ی ساده:</p>

</div>

```
Download Time ≈ TTLB - TTFB
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">نکته در مورد TTLB:  مخفف Time To Last Byte است و نشان می‌دهد آخرین بایت response چه زمانی به client رسیده است. از دید client، TTLB معمولاً معادل عملی Full Response Time&rlm;.</p>

<ul dir="rtl">
  <li>TTFB = زمان تا شروع response</li>
  <li>TTLB = زمان تا پایان response</li>
  <li>Response Time ≈ TTLB از دید client</li>
</ul>

<p dir="rtl" align="right">مقایسه‌ی TTFB و TTLB برای تحلیل عملکرد بسیار مفید است.  مثال: TTFB = 150ms و TTLB = 2150ms پس</p>

</div>

```
Download Time ≈ 2000ms
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">برداشت این است که سرور نسبتاً زود شروع به پاسخ‌دهی کرده، اما دریافت کامل payload زمان زیادی برده است. دلیل آن می‌تواند بزرگ بودن response، کندی شبکه یا محدودیت client باشد. Download Time در موارد زیر اهمیت بیشتری دارد:</p>

<ul dir="rtl">
  <li>File Download</li>
  <li>Image / Video Delivery</li>
  <li>Large JSON Payload</li>
  <li>CSV Export</li>
  <li>Report Generation</li>
  <li>Streaming</li>
</ul>

<h4 dir="rtl" align="right">Server Processing Time</h4>

<p dir="rtl" align="right">زمان پردازش سرور یا Server Processing Time بخشی از Response Time است که داخل backend صرف پردازش request می‌شود. این زمان می‌تواند شامل موارد زیر باشد:</p>

<ul dir="rtl">
  <li>CPU Processing Time</li>
  <li>Database Query Time</li>
  <li>External RPC Time</li>
  <li>Internal Function Execution Time</li>
  <li>Template Rendering Time</li>
  <li>Business Logic Time</li>
  <li>Serialization Time</li>
</ul>

<p dir="rtl" align="right">نکته‌ی مهم این است که Server Processing Time فقط یکی از اجزای Response Time است. Response Time از دید client علاوه بر پردازش سرور، شامل network، gateway، download time و سایر overhead ها نیز می‌شود.</p>

<h3 dir="rtl" align="right">Time To First Byte یا TTFB</h3>

<p dir="rtl" align="right">عنوان TTFB مخفف Time To First Byte است و به زمانی گفته می‌شود که از لحظه‌ی ارسال request تا دریافت اولین بایت response توسط client طول می‌کشد. به بیان ساده: TTFB نشان می‌دهد سیستم چه زمانی شروع به پاسخ‌دهی کرده است.  نکته‌ی مهم این است که TTFB الزاماً فقط server processing time نیست. به عنوان نمونه این می تواند یک flow ساده ارسال درخواست باشد:</p>

<ul dir="rtl">
  <li>Client sends request</li>
  <li>Network</li>
  <li>Backend starts processing</li>
  <li>Backend begins sending response</li>
  <li>Client receives first byte  ← TTFB</li>
  <li>Client receives full response ← Full Response Time</li>
</ul>

<p dir="rtl" align="right">این به این معنی است که  سرور سریع شروع کرده جواب بدهد، ولی کامل شدن response زمان بیشتری گرفته است. مثلاً برای دانلود فایل، streaming API، یا صفحه وب، TTFB مهم است.  فرض کنید Instagram یک feed برمی‌گرداند. اگر backend بتواند زود اولین بخش response را بفرستد، کلاینت می‌تواند سریع‌تر شروع به render کردن کند، حتی اگر کل payload هنوز کامل نیامده باشد. مثال: TTFB: 150ms و Full response completed: 900ms یعنی کاربر یا کلاینت خیلی زود فهمیده که server زنده است و دارد جواب می‌دهد.</p>

<p dir="rtl" align="right">اگر TTFB از دید client یا browser اندازه‌گیری شود، می‌تواند شامل موارد زیر باشد:</p>

<ul dir="rtl">
  <li>DNS Resolution</li>
  <li>TCP Connection</li>
  <li>TLS Handshake</li>
  <li>Request Network Time</li>
  <li>Load Balancer / Gateway Overhead</li>
  <li>Queue Time</li>
  <li>Server Processing Before First Write</li>
  <li>Serialization Before First Write</li>
  <li>Network Return Time</li>
</ul>

<p dir="rtl" align="right">برای مثال:</p>

<ul dir="rtl">
  <li>DNS:                       20ms</li>
  <li>TCP:                       40ms</li>
  <li>TLS:                       60ms</li>
  <li>Request Network Time:      20ms</li>
  <li>Gateway Overhead:          10ms</li>
  <li>Backend Queue Time:        50ms</li>
  <li>Backend Processing:       120ms</li>
  <li>First Byte Return Time:    20ms</li>
</ul>

</div>

```
TTFB = 340ms
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">بنابراین اگر TTFB بالا باشد، الزاماً نمی‌توان نتیجه گرفت که فقط backend کند است. مشکل ممکن است از DNS، TLS، شبکه، gateway، صف، دیتابیس، serialization یا application logic باشد.</p>

<h4 dir="rtl" align="right">کاربردهای اصلی TTFB</h4>

<p dir="rtl" align="right">عموما TTFB مخصوصاً در سناریوهایی اهمیت دارد که شروع پاسخ برای تجربه‌ی کاربر مهم‌تر از کامل شدن کل پاسخ است. نمونه‌ها:</p>

<ul dir="rtl">
  <li>Web Pages</li>
  <li>Streaming APIs</li>
  <li>File Downloads</li>
  <li>Large Responses</li>
  <li>Server-Sent Events</li>
  <li>gRPC Streaming</li>
  <li>Video / Audio Streaming</li>
  <li>LLM-like Token Streaming</li>
</ul>

<h4 dir="rtl" align="right">Upstream First Byte Time</h4>

<p dir="rtl" align="right">در سیستم‌هایی که چند hop دارند، فقط TTFB نهایی کافی نیست. باید بررسی شود هر hop چه زمانی اولین بایت را از upstream خود دریافت کرده است. برای مثال: Client → CDN → API Gateway → Service A → Service B&rlm; . در چنین ساختاری ممکن است چند نوع first byte time داشته باشیم:</p>

<ul dir="rtl">
  <li>Service B → Service A First Byte Time</li>
  <li>Service A → Gateway First Byte Time</li>
  <li>Gateway → Client First Byte Time</li>
</ul>

<p dir="rtl" align="right">این متریک‌ها مخصوصاً در API Gateway ها و reverse proxy هایی مانند Nginx، Envoy و Traefik بسیار مفید هستند. برای مثال:</p>

<ul dir="rtl">
  <li>Gateway TTFB to Client:       700ms</li>
  <li>Upstream First Byte Time:     650ms</li>
  <li>Gateway Overhead:              50ms</li>
</ul>

<p dir="rtl" align="right">در این حالت، مشکل احتمالاً در upstream service است، نه خود gateway&rlm;.</p>

<h4 dir="rtl" align="right">نقش Serialization Latency و اثر آن بر TTFB</h4>

<p dir="rtl" align="right">عموما Serialization Latency زمانی است که سیستم برای آماده‌سازی اولین خروجی قابل ارسال صرف می‌کند؛ مثلاً تبدیل struct به JSON&rlm;.اگر backend مجبور باشد کل response را در memory بسازد و سپس آن را serialize کند، TTFB ممکن است بالا برود.</p>

<p dir="rtl" align="right">طراحی اول:</p>

<ul dir="rtl">
  <li>Fetch all data</li>
  <li>Build full object</li>
  <li>Serialize full JSON</li>
  <li>Send response</li>
</ul>

<p dir="rtl" align="right">در این حالت، اولین بایت زمانی ارسال می‌شود که کل response آماده شده باشد.</p>

<p dir="rtl" align="right">طراحی دوم:</p>

<ul dir="rtl">
  <li>Start response early</li>
  <li>Send chunks gradually</li>
</ul>

<p dir="rtl" align="right">در این حالت، TTFB کاهش پیدا می‌کند، هرچند ممکن است full response time همچنان بالا باشد.</p>

<h3 dir="rtl" align="right">جمع‌بندی:</h3>

<p dir="rtl" align="right">اگر بخواهیم تفاوت اصلی این سه را در یک جمله خلاصه کنیم:</p>

<ul dir="rtl">
  <li>تاخیر - Latency: تاخیر در یک مسیر یا بخش از سیستم، مفهوم عمومی‌تر (برای تحلیل تاخیر در هر بخش از سیستم استفاده می‌شود)</li>
  <li>زمان پاسخ - Response Time: کل زمان کامل شدن پاسخ از دید کلاینت(برای فهمیدن اینکه کل پاسخ چه زمانی کامل می‌شود استفاده می‌شود)</li>
  <li>زمان TTFB: زمان رسیدن اولین byte پاسخ به کلاینت (برای فهمیدن اینکه پاسخ‌دهی چه زمانی شروع می‌شود استفاده می‌شود)</li>
</ul>

<p dir="rtl" align="right">فرض کن این API را داریم:</p>

</div>

```
GET /feed
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">زمان‌ها این شکلی‌اند:</p>

</div>

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

<div dir="rtl" align="right">

<p dir="rtl" align="right">حالا:</p>

</div>

```
TTFB = 220ms
Response Time = 520ms
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">و latency بسته به context می‌تواند این‌ها باشد:</p>

</div>

```
Network latency = 20ms رفت + برگشت/بخش شبکه
DB latency = 100ms
Server-side latency = حدود 180ms تا 200ms
Client-observed latency = 520ms
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">پس وقتی کسی فقط می‌گوید latency، باید بپرسیم: Latency از دید کجا؟ کلاینت؟ سرور؟ دیتابیس؟ شبکه؟</p>

<h3 dir="rtl" align="right">کی از کدام استفاده کنیم؟</h3>

<p dir="rtl" align="right">وقتی درباره تجربه کاربر یا API حرف می‌زنیم: Response Time  مثلاً در Design Doc بنویسی:</p>

<p dir="rtl" align="right">P95 response time for GET /feed should be under 500ms&rlm;.</p>

<p dir="rtl" align="right">این برای user-facing APIها خیلی مناسب است. مثال:</p>

<ul dir="rtl">
  <li>Login API response time P95 &lt; 300ms</li>
  <li>Checkout API response time P95 &lt; 700ms</li>
  <li>Search API response time P95 &lt; 500ms</li>
</ul>

<p dir="rtl" align="right">اینجا مهم است کاربر چقدر منتظر می‌ماند.</p>

<p dir="rtl" align="right">وقتی می خواهیم  bottleneck پیدا می‌کنیم: <strong>Latency</strong></p>

<p dir="rtl" align="right">عموما Latency برای تحلیل داخلی سیستم بهتر است. در واقع برای debugging و bottleneck analysis استفاده می‌شود. مثلاً:</p>

<ul dir="rtl">
  <li>DB latency increased from 20ms to 200ms&rlm;.</li>
  <li>Redis latency is stable&rlm;.</li>
  <li>External payment provider latency is high&rlm;.</li>
  <li>Network latency between services increased&rlm;.</li>
</ul>

<p dir="rtl" align="right">اینجا نمی‌گویید response time، چون دارید درباره یک بخش خاص حرف می‌زنیم.</p>

<p dir="rtl" align="right">مثلاً در debugging&rlm;:</p>

<ul dir="rtl">
  <li>Total response time is 800ms&rlm;.</li>
  <li>DB latency is 500ms&rlm;.</li>
  <li>So the bottleneck is probably database access&rlm;.</li>
</ul>

<p dir="rtl" align="right">وقتی response بزرگ، دانلود فایل، streaming یا web page داری: TTFB</p>

<p dir="rtl" align="right">زمان TTFB مخصوصاً وقتی مهم است که:</p>

<ul dir="rtl">
  <li>response بزرگ است</li>
  <li>فایل دانلود می‌شود</li>
  <li>HTML page render می‌شود</li>
  <li>API streaming داری</li>
  <li>Server-Sent Events داری</li>
  <li>gRPC streaming داری</li>
</ul>

<p dir="rtl" align="right">کاربر باید زود بفهمد سیستم شروع به پاسخ کرده است. مثلاً:</p>

<p dir="rtl" align="right">TTFB should be under 200ms for the homepage&rlm;.</p>

<p dir="rtl" align="right">برای سایت‌ها، TTFB خیلی مهم است چون browser تا اولین byte را نگیرد، نمی‌تواند شروع خوبی برای render داشته باشد. برای streaming هم مهم است. مثلاً در ChatGPT-like response&rlm;:</p>

<ul dir="rtl">
  <li>User sends prompt</li>
  <li>TTFB = time until first token appears</li>
  <li>Full response time = time until whole answer is generated</li>
</ul>

<p dir="rtl" align="right">اگر اولین token بعد از ۵۰۰ms بیاید، کاربر حس می‌کند سیستم responsive است، حتی اگر کل جواب ۱۰ ثانیه طول بکشد.</p>

<p dir="rtl" align="right">نکته‌ی کلیدی این است که TTFB بخشی از Response Time است؛ زیرا Response Time از شروع request تا دریافت آخرین بایت ادامه دارد، در حالی که TTFB فقط تا دریافت اولین بایت را اندازه‌گیری می‌کند.</p>

<p dir="rtl" align="right">همچنین Latency مفهومی عمومی‌تر از هر دو است و می‌تواند برای کل request یا هر بخش داخلی سیستم، مانند شبکه، دیتابیس، کش، RPC، صف و business logic استفاده شود.</p>

<p dir="rtl" align="right">بنابراین در یک Design Doc یا System Design Discussion، بهتر است معیارها به‌صورت دقیق و قابل اندازه‌گیری بیان شوند. برای مثال:</p>

<ul dir="rtl">
  <li>For read APIs, P95 response time should be under 300ms&rlm;.</li>
  <li>For payment authorization, P99 response time should be under 1s&rlm;.</li>
  <li>The main latency contributors should be tracked separately: DB latency, cache latency, RPC latency, and queue time&rlm;.</li>
  <li>For streaming or large responses, TTFB should be under 200ms and TTLB should be monitored separ</li>
</ul>

</div>
