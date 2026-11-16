<div dir="rtl" align="right">

<h1 dir="rtl" align="right">تاخیر - Latency (یکی از مهم‌ترین مفاهیم بک‌اندی)</h1>

<p dir="rtl" align="right">در نظر داشته باشید که Latency یک عدد واحد نیست. اجزای Latency در Backend عموما به صورت زیر می‌باشد:</p>

<div dir="ltr" align="left">

```
Total Response Time = Network + Load Balancer + Auth + Application Logic + DB + Cache + External Services
```

</div>

<p dir="rtl" align="right">در Design/System Doc باید مشخص باشد:</p>

<ul dir="rtl">
  <li>کجاها وقت تلف می شود؟</li>
  <li>کدام بخش بیشترین سهم را دارد؟ کدام بخش bottleneck است؟</li>
  <li>کدام قابل بهینه‌سازی است؟ در‌واقع کدام بخش دست ماست و می‌توانیم روی آن تأثیر بگذاریم</li>
  <li>کدام غیرقابل کنترل است؟ کدام دست ما نیست</li>
</ul>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right"><strong>1) این فرمول فقط برای اجزای پشت‌سرهم (sequential) یک جمع ساده است.</strong> اگر دو بخش موازی اجرا شوند (مثلاً دو call همزمان به DB و یک سرویس خارجی)، زمان آن‌ها جمع نمی‌شود و فقط کندترین‌شان در مسیر اصلی (critical path) اثر می‌گذارد:</p>

<ul dir="rtl">
  <li>پشت‌سرهم: 80ms (DB) + 300ms (Payment RPC) = 380ms</li>
  <li>موازی: max(80ms, 300ms) = 300ms</li>
</ul>

<p dir="rtl" align="right">بنابراین قبل از بهینه‌سازی باید بپرسیم کدام بخش روی critical path است. بهینه‌سازی بخشی که موازی و سریع‌تر از بقیه است، هیچ اثری روی Latency نهایی ندارد.</p>

<p dir="rtl" align="right"><strong>2) میانگین‌ها جمع می‌شوند، ولی Percentile ها نه.</strong> اگر میانگین DB برابر 80ms و میانگین سرویس خارجی 300ms باشد، میانگین کل (در حالت پشت‌سرهم) برابر 380ms است. اما P95 کل، جمع P95 اجزا نیست. چرا؟ چون بعید است هر دو جزء همزمان در بدترین حالت خودشان باشند. پس فرمول بالا را برای اجزا با میانگین یا برای یک request مشخص (در trace) به کار ببرید، نه برای جمع کردن P99 ها (<a href="percentiles.md">Percentile-based Thinking</a>).</p>

<p dir="rtl" align="right"><strong>3) جواب دو سؤال آخر («دست ماست» و «دست ما نیست»):</strong></p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">دست ماست (قابل کنترل)</th>
      <th align="right">دست ما نیست (یا فقط غیرمستقیم)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">منطق برنامه، کپی‌های غیرضروری، serialization</td>
      <td align="right">سرعت نور (فاصله‌ی فیزیکی کاربر تا سرور)</td>
    </tr>
    <tr>
      <td align="right">query و index دیتابیس، N+1</td>
      <td align="right">کیفیت شبکه‌ی کاربر (موبایل، Wi-Fi ضعیف)</td>
    </tr>
    <tr>
      <td align="right">اندازه‌ی pool و timeout و صف</td>
      <td align="right">Latency سرویس شخص ثالث (مثل درگاه پرداخت)</td>
    </tr>
    <tr>
      <td align="right">Cache، Batching، Pagination</td>
      <td align="right">رفتار DNS و ISP کاربر</td>
    </tr>
    <tr>
      <td align="right">محل deploy (نزدیک کردن سرور به کاربر با CDN یا multi-region)</td>
      <td align="right">زمان GC یا مشکلات سخت‌افزاری ابر</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">برای موارد ستون دوم راه‌حل «کاهش اثر» است، نه «حل»: timeout، retry سنجیده، cache، fallback و جدا کردن مسیر از مسیر بحرانی.</p>

<p dir="rtl" align="right"><strong>4) مقیاس عددی (فقط order of magnitude و وابسته به سخت‌افزار):</strong></p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">عملیات</th>
      <th align="right">زمان تقریبی</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">خواندن از حافظه‌ی اصلی (RAM)</td>
      <td align="right">حدود 100ns</td>
    </tr>
    <tr>
      <td align="right">خواندن تصادفی از SSD</td>
      <td align="right">حدود 0.1ms</td>
    </tr>
    <tr>
      <td align="right">رفت‌وبرگشت شبکه داخل یک data center</td>
      <td align="right">حدود 0.5ms</td>
    </tr>
    <tr>
      <td align="right">seek روی هارد دیسک (HDD)</td>
      <td align="right">حدود 10ms</td>
    </tr>
    <tr>
      <td align="right">رفت‌وبرگشت شبکه بین قاره‌ها</td>
      <td align="right">حدود 100 تا 150ms</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">این جدول یک نتیجه‌ی مهم دارد: یک call به Redis در همان data center (حدود 1ms) هزاران برابر کندتر از خواندن از حافظه‌ی process است، و یک call بین‌قاره‌ای صدها برابر کندتر از آن. پس تعداد call های شبکه‌ای در یک request معمولاً مهم‌تر از «بهینه‌تر نوشتن» یک حلقه است.</p>

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

<p dir="rtl" align="right">پس Latency یک واژه عمومی‌تر است. مثلاً می گوییم: DB latency is high یا Network latency is high یا API latency P95 is 400ms&rlm;. این یعنی latency همیشه الزاماً فقط کل response نیست؛ می‌تواند latency یک بخش از سیستم هم باشد.</p>

<p dir="rtl" align="right">عموما Latency یک مفهوم عمومی برای اندازه‌گیری تاخیر در انجام یک عملیات است. این عملیات می‌تواند مربوط به کل API، شبکه، دیتابیس، کش، صف، سرویس داخلی یا یک dependency خارجی باشد. بنابراین Latency فقط به زمان شروع پاسخ یا شروع پردازش محدود نمی‌شود، بلکه می‌تواند در هر نقطه‌ای از مسیر درخواست اندازه‌گیری شود. نمونه‌هایی از Latency عبارت‌اند از:</p>

<ul dir="rtl">
  <li>Network Latency</li>
  <li>Database Latency</li>
  <li>Redis Latency</li>
  <li>gRPC Latency</li>
  <li>Queue Latency</li>
  <li>API Latency</li>
</ul>

<p dir="rtl" align="right">به بیان دقیق‌تر: Latency نشان‌دهنده‌ی مدت‌زمان تاخیر در یک بخش مشخص از مسیر پردازش request است. این بخش می‌تواند شبکه، صف، اپلیکیشن، دیتابیس، کش، سرویس خارجی یا کل API باشد.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>این واژه‌ها در همه‌ی منابع یک تعریف ندارند.</strong> مثلاً در حوزه‌ی شبکه، Latency معمولاً یعنی تاخیر یک‌طرفه یا RTT. در کتاب Designing Data-Intensive Applications، Latency فقط مدتی است که request منتظر پردازش می‌ماند و Response Time کل چیزی است که کلاینت می‌بیند. تعریف این فایل (Latency به‌عنوان تاخیر یک بخش و Response Time به‌عنوان کل زمان از دید کلاینت) رایج‌ترین استفاده در مستندات مهندسی است، ولی در مصاحبه یا هنگام خواندن یک سند، بهتر است همیشه بپرسید «از دید کجا؟».</li>
  <li><strong>Latency با Throughput یکی نیست.</strong> Latency زمان یک request است و Throughput تعداد request در واحد زمان. رابطه‌ی آن‌ها از طریق Concurrency در <a href="performance.md">Performance</a> آمده است.</li>
</ul>

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
  <li>DNS Lookup: 20ms</li>
  <li>TCP Connect: 40ms</li>
  <li>TLS Handshake: 70ms</li>
  <li>Request Send: 5ms</li>
</ul>

<p dir="rtl" align="right">این زمان‌ها هنوز وارد business logic نشده‌اند، اما مستقیماً روی تجربه‌ی کاربر اثر می‌گذارند.</p>

<p dir="rtl" align="right">نکته‌ی مهم این است که این نوع latency معمولاً از سمت client، gateway، proxy، CDN یا ابزارهای synthetic monitoring قابل مشاهده است و در backend monitoring خام همیشه به‌صورت مستقیم دیده نمی‌شود.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>این هزینه‌ها فقط برای connection جدید است.</strong> اگر کلاینت همان connection را دوباره استفاده کند (keep-alive یا HTTP/2)، DNS و TCP و TLS دوباره پرداخت نمی‌شوند. به همین دلیل کلاینت‌های سرویس به سرویس باید connection ها را reuse کنند. در Go، <code>http.Client</code> به‌صورت پیش‌فرض این کار را می‌کند، ولی فقط وقتی که body پاسخ کامل خوانده و بسته شود. در غیر این صورت هر request یک connection تازه می‌سازد. جزئیات pool در <a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a> آمده است.</li>
  <li><strong>TCP حدود 1 RTT و TLS حدود 1 تا 2 RTT</strong> (بسته به نسخه‌ی TLS) هزینه دارد. در مثال بالا، اگر RTT برابر 40ms باشد، عددهای 40ms برای TCP و 70ms برای TLS منطقی هستند.</li>
  <li><strong>یک کف فیزیکی وجود دارد.</strong> نور در فیبر با سرعت حدود 200,000 km/s حرکت می‌کند. پس برای فاصله‌ی 10,000 km، تاخیر یک‌طرفه حداقل حدود 50ms و RTT حداقل حدود 100ms است. هیچ بهینه‌سازی نرم‌افزاری آن را کمتر نمی‌کند، و تنها راه، نزدیک کردن سرور (یا Cache) به کاربر است (CDN، multi-region).</li>
</ul>

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
  <li>DB Latency: 50ms</li>
  <li>Business Logic: 100ms</li>
  <li>Queue Time: 600ms</li>
</ul>

<p dir="rtl" align="right">در این مثال، مشکل اصلی دیتابیس نیست؛ بلکه request قبل از شروع پردازش مدت زیادی در صف مانده است.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>در این مثال 50 + 100 + 600 = 750ms است، یعنی <strong>50ms از 800ms توضیح داده نشده</strong> (احتمالاً middleware، serialization یا GC). این خودش یک درس است: وقتی جمع اجزا با کل نمی‌خواند، یعنی بخشی از مسیر اندازه‌گیری نشده و باید instrument شود.</li>
  <li><strong>Queue Time اگر اندازه نگیرید دیده نمی‌شود.</strong> معمولاً ساعت متریک «زمان handler» وقتی شروع می‌شود که handler اجرا شده است، یعنی صف قبل از آن اصلاً در عدد حساب نمی‌شود. برای دیدنش باید یا زمان رسیدن request را ثبت کنیم (مثلاً timestamp در gateway) یا زمان انتظار برای منابعی مثل connection pool را جدا اندازه بگیریم. در Go مقدار <code>WaitDuration</code> در <code>db.Stats()</code> نشان می‌دهد request ها چقدر برای گرفتن یک connection دیتابیس منتظر مانده‌اند (<a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a>).</li>
</ul>

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
  <li>Redis Latency: 5ms</li>
  <li>Payment RPC Latency: 300ms</li>
  <li>Kafka Produce Latency: 20ms</li>
</ul>

<p dir="rtl" align="right">در تحلیل عملکرد سیستم، بررسی این بخش‌ها کمک می‌کند مشخص شود bottleneck اصلی کجاست: دیتابیس، کش، سرویس خارجی، صف پیام یا خود application logic&rlm;.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>هر dependency باید <strong>timeout جدا</strong> داشته باشد. اگر Payment RPC در حالت عادی 300ms است، timeout حدود 1 تا 2s منطقی است، نه 30s&rlm;. در غیر این صورت کندی همان dependency، منابع سرویس شما را نگه می‌دارد.</li>
  <li>وقتی یک request به چند dependency وابسته باشد، <strong>احتمال اینکه دست‌کم یکی کند باشد</strong> بالا می‌رود (Tail Latency Amplification). این موضوع در <a href="percentiles.md">Percentile-based Thinking</a> توضیح داده شده است.</li>
  <li>برای دیدن سهم هر dependency در یک request مشخص، Distributed Tracing بهترین ابزار است (<a href="measurement-tools.md">ابزارهای اندازه‌گیری</a>).</li>
</ul>

<h3 dir="rtl" align="right">Response Time</h3>

<p dir="rtl" align="right">مدت زمان پاسخ یا Response Time معمولاً از دید کاربر یا کلاینت اندازه‌گیری می‌شود. یعنی: از لحظه ارسال request تا کامل دریافت شدن response (از وقتی کلاینت request را می‌فرستد تا وقتی کل response را دریافت می‌کند)</p>

<p dir="rtl" align="right">برای APIهای معمولی، latency و response time خیلی وقت‌ها تقریباً به‌جای هم استفاده می‌شوند. اما در سیستم‌های بزرگ‌تر ممکن است دقیق‌تر تفکیک شوند.</p>

<p dir="rtl" align="right">مثلاً:</p>

<ul dir="rtl">
  <li>Server-side latency: 120ms</li>
  <li>Network overhead: 40ms</li>
  <li>Total response time: 160ms</li>
</ul>

<p dir="rtl" align="right">اگر کلاینت بعد از دریافت پاسخ، آن را parse هم بکند (مثلاً 20ms)، کل زمان تا «قابل استفاده شدن نتیجه» برابر 180ms می‌شود. <strong>این 20ms طبق تعریف بالا جزو Response Time نیست</strong>، چون بعد از دریافت آخرین بایت اتفاق می‌افتد.</p>

<p dir="rtl" align="right">بنابراین Response Time کل زمان request/response از دید client است؛ از ارسال request تا دریافت آخرین بایت پاسخ. برای APIهای معمولی، Response Time معمولاً مهم‌ترین معیار user-facing است. مانند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">API</th>
      <th align="right">هدف نمونه (فقط مثال)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">POST /login</td>
      <td align="right">P95 Response Time &lt; 300ms</td>
    </tr>
    <tr>
      <td align="right">POST /payment/authorize</td>
      <td align="right">P95 Response Time &lt; 700ms</td>
    </tr>
    <tr>
      <td align="right">GET /feed</td>
      <td align="right">P95 Response Time &lt; 800ms</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">در این نوع APIها، معمولاً response کوچک است و کاربر یا کلاینت منتظر دریافت کامل نتیجه است؛ بنابراین Response Time معیار اصلی تجربه‌ی کاربر محسوب می‌شود.</p>

<p dir="rtl" align="right"><strong>نکته:</strong> عددهای جدول بالا نمونه هستند و از نیاز بیزینس می‌آیند، نه از یک قانون جهانی. برای نوشتن آن‌ها به شکل NFR به <a href="../01-non-functional-requirements/non-functional-requirements.md">Non-Functional Requirements</a> و برای نوشتن به شکل SLO به <a href="sli-slo-sla.md">SLI / SLO / SLA</a> نگاه کنید.</p>

<h4 dir="rtl" align="right">Total Response Time</h4>

<p dir="rtl" align="right">کل زمان پاسخ دهی Total Response Time معیار پایه برای سنجش کامل شدن پاسخ است، اما باید مشخص شود از دید کدام نقطه اندازه‌گیری می‌شود:</p>

<ul dir="rtl">
  <li>Client-observed Response Time</li>
  <li>Gateway-observed Response Time</li>
  <li>Server-side Response Time</li>
</ul>

<p dir="rtl" align="right">این سه مقدار می‌توانند با هم متفاوت باشند. مثال:</p>

<ul dir="rtl">
  <li>Backend Server Time: 180ms</li>
  <li>Gateway Observed Time: 220ms</li>
  <li>Client Observed Time: 350ms</li>
</ul>

<p dir="rtl" align="right">تفاوت این اعداد می‌تواند ناشی از network latency، TLS overhead، download time، CDN، proxy یا موقعیت جغرافیایی client باشد.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>دام رایج:</strong> اگر داشبورد شما فقط Server-side Response Time را نشان دهد (مثل اکثر histogram های داخل اپلیکیشن)، همیشه خوش‌بین‌تر از تجربه‌ی واقعی کاربر است. در مثال بالا، سرور 180ms می‌گوید ولی کاربر 350ms می‌بیند.</li>
  <li>پس هنگام نوشتن SLO باید صریح بنویسید <strong>کجا اندازه می‌گیریم</strong>. یک سازش رایج، اندازه‌گیری در gateway یا load balancer است، چون شبکه‌ی بین کاربر و gateway را نمی‌بیند ولی کل زمان داخلی را می‌بیند. برای تجربه‌ی واقعی کاربر به Real User Monitoring (RUM) یا synthetic monitoring نیاز داریم (<a href="measurement-tools.md">ابزارهای اندازه‌گیری</a>).</li>
</ul>

<h4 dir="rtl" align="right">Download Time</h4>

<p dir="rtl" align="right">زمان دانلود یا همان Download Time فاصله‌ی بین دریافت اولین بایت و دریافت آخرین بایت response است. رابطه‌ی ساده:</p>

<div dir="ltr" align="left">

```
Download Time ≈ TTLB - TTFB
```

</div>

<p dir="rtl" align="right">نکته در مورد TTLB: مخفف Time To Last Byte است و نشان می‌دهد آخرین بایت response چه زمانی به client رسیده است. از دید client، TTLB معمولاً معادل عملی Full Response Time&rlm;.</p>

<ul dir="rtl">
  <li>TTFB = زمان تا شروع response</li>
  <li>TTLB = زمان تا پایان response</li>
  <li>Response Time ≈ TTLB از دید client</li>
</ul>

<p dir="rtl" align="right">مقایسه‌ی TTFB و TTLB برای تحلیل عملکرد بسیار مفید است. مثال: TTFB = 150ms و TTLB = 2150ms پس Download Time ≈ 2000ms&rlm;.</p>

<p dir="rtl" align="right">برداشت این است که سرور نسبتاً زود شروع به پاسخ‌دهی کرده، اما دریافت کامل payload زمان زیادی برده است. دلیل آن می‌تواند بزرگ بودن response، کندی شبکه یا محدودیت client باشد. Download Time در موارد زیر اهمیت بیشتری دارد:</p>

<ul dir="rtl">
  <li>File Download</li>
  <li>Image / Video Delivery</li>
  <li>Large JSON Payload</li>
  <li>CSV Export</li>
  <li>Report Generation</li>
  <li>Streaming</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>Download Time را می‌شود تخمین زد:</strong> زمان تقریبی = حجم ÷ پهنای باند. یک response با حجم 5MB (یعنی حدود 40 مگابیت) روی یک اتصال 20Mbps حدود 2 ثانیه طول می‌کشد. این دقیقاً عدد مثال بالاست و نشان می‌دهد کندی آن مشکل سرور نیست.</li>
  <li>راه‌های کم کردنش، کوچک کردن payload است: فشرده‌سازی (gzip یا brotli)، فقط فیلدهای لازم را فرستادن، و Pagination به‌جای فرستادن همه‌ی داده (<a href="../02-04-batching-and-pagination/pagination.md">Pagination</a>).</li>
  <li><strong>اگر TTFB کم و Download Time زیاد باشد، افزودن CPU یا سرور کمکی نمی‌کند.</strong> مشکل در حجم داده یا شبکه است. تشخیص این حالت‌ها در <a href="endpoint-analysis-checklist.md">چک‌لیست تحلیل Endpoint</a> آمده است.</li>
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

<p dir="rtl" align="right">عنوان TTFB مخفف Time To First Byte است و به زمانی گفته می‌شود که از لحظه‌ی ارسال request تا دریافت اولین بایت response توسط client طول می‌کشد. به بیان ساده: TTFB نشان می‌دهد سیستم چه زمانی شروع به پاسخ‌دهی کرده است. نکته‌ی مهم این است که TTFB الزاماً فقط server processing time نیست. به عنوان نمونه این می تواند یک flow ساده ارسال درخواست باشد:</p>

<ul dir="rtl">
  <li>Client sends request</li>
  <li>Network</li>
  <li>Backend starts processing</li>
  <li>Backend begins sending response</li>
  <li>Client receives first byte ← TTFB</li>
  <li>Client receives full response ← Full Response Time</li>
</ul>

<p dir="rtl" align="right">این به این معنی است که سرور سریع شروع کرده جواب بدهد، ولی کامل شدن response زمان بیشتری گرفته است. مثلاً برای دانلود فایل، streaming API، یا صفحه وب، TTFB مهم است. فرض کنید Instagram یک feed برمی‌گرداند. اگر backend بتواند زود اولین بخش response را بفرستد، کلاینت می‌تواند سریع‌تر شروع به render کردن کند، حتی اگر کل payload هنوز کامل نیامده باشد. مثال: TTFB: 150ms و Full response completed: 900ms یعنی کاربر یا کلاینت خیلی زود فهمیده که server زنده است و دارد جواب می‌دهد.</p>

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
  <li>DNS: 20ms</li>
  <li>TCP: 40ms</li>
  <li>TLS: 60ms</li>
  <li>Request Network Time: 20ms</li>
  <li>Gateway Overhead: 10ms</li>
  <li>Backend Queue Time: 50ms</li>
  <li>Backend Processing: 120ms</li>
  <li>First Byte Return Time: 20ms</li>
</ul>

<p dir="rtl" align="right">جمع این اجزا می‌شود <strong>TTFB = 340ms</strong>&rlm;.</p>

<p dir="rtl" align="right">بنابراین اگر TTFB بالا باشد، الزاماً نمی‌توان نتیجه گرفت که فقط backend کند است. مشکل ممکن است از DNS، TLS، شبکه، gateway، صف، دیتابیس، serialization یا application logic باشد.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>TTFB دو تعریف رایج دارد.</strong> بعضی ابزارها آن را از لحظه‌ی شروع navigation اندازه می‌گیرند (شامل DNS و TCP و TLS، مثل مثال بالا). بعضی دیگر فقط از لحظه‌ی ارسال request تا اولین بایت پاسخ (بدون DNS و TCP و TLS). پس هنگام مقایسه‌ی دو عدد TTFB، مطمئن شوید با یک تعریف اندازه‌گیری شده‌اند.</li>
  <li><strong>اندازه‌گیری سریع با curl.</strong> هر کدام از این مقادیر، زمان <strong>تجمعی</strong> از شروع است، نه زمان همان مرحله. پس برای مدت هر مرحله باید از مقدار قبلی کم کرد:</li>
</ul>

<div dir="ltr" align="left">

```bash
curl -o /dev/null -s -w "dns=%{time_namelookup} tcp=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer} total=%{time_total}\n" https://example.com/feed
```

</div>

<ul dir="rtl">
  <li>در مثال 340ms، اگر connection را reuse کنیم، DNS و TCP و TLS (جمعاً 120ms) حذف می‌شوند و TTFB حدود 220ms می‌شود.</li>
  <li>برای LLM و پاسخ‌های token به token، اصطلاح رایج‌تر <strong>TTFT</strong> (Time To First Token) است و معنی‌اش همان TTFB است.</li>
  <li><strong>تله‌ی رایج در streaming:</strong> اگر بین سرور و کلاینت یک proxy یا middleware باشد که پاسخ را buffer می‌کند (مثلاً <code>proxy_buffering</code> در Nginx یا فشرده‌سازی که تا پر شدن بافر صبر می‌کند)، کلاینت اولین بایت را دیرتر می‌بیند حتی اگر کد شما زود Flush کرده باشد. اگر streaming می‌خواهید، این لایه‌ها را هم بررسی کنید.</li>
  <li>مقایسه‌ی TTFB و Total برای تشخیص نوع مشکل (قبل از شروع پاسخ، یا بعد از آن) در <a href="endpoint-analysis-checklist.md">چک‌لیست تحلیل Endpoint</a> آمده است.</li>
</ul>

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

<p dir="rtl" align="right">در سیستم‌هایی که چند hop دارند، فقط TTFB نهایی کافی نیست. باید بررسی شود هر hop چه زمانی اولین بایت را از upstream خود دریافت کرده است. برای مثال: Client → CDN → API Gateway → Service A → Service B&rlm;. در چنین ساختاری ممکن است چند نوع first byte time داشته باشیم:</p>

<ul dir="rtl">
  <li>Service B → Service A First Byte Time</li>
  <li>Service A → Gateway First Byte Time</li>
  <li>Gateway → Client First Byte Time</li>
</ul>

<p dir="rtl" align="right">این متریک‌ها مخصوصاً در API Gateway ها و reverse proxy هایی مانند Nginx، Envoy و Traefik بسیار مفید هستند. برای مثال:</p>

<ul dir="rtl">
  <li>Gateway TTFB to Client: 700ms</li>
  <li>Upstream First Byte Time: 650ms</li>
  <li>Gateway Overhead: 50ms</li>
</ul>

<p dir="rtl" align="right">در این حالت، مشکل احتمالاً در upstream service است، نه خود gateway&rlm;.</p>

<h4 dir="rtl" align="right">نقش Serialization Latency و اثر آن بر TTFB</h4>

<p dir="rtl" align="right">عموما Serialization Latency زمانی است که سیستم برای آماده‌سازی اولین خروجی قابل ارسال صرف می‌کند؛ مثلاً تبدیل struct به JSON&rlm;. اگر backend مجبور باشد کل response را در memory بسازد و سپس آن را serialize کند، TTFB ممکن است بالا برود.</p>

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

<p dir="rtl" align="right"><strong>نکات تکمیلی: مثال Go و اندازه‌گیری واقعی</strong></p>

<p dir="rtl" align="right">دو handler زیر یک endpoint را به دو شکل پیاده‌سازی می‌کنند. داده‌ی هر آیتم 100ms طول می‌کشد (شبیه‌سازی یک منبع کند) و 5 آیتم داریم:</p>

<div dir="ltr" align="left">

```go
// loadItem simulates a slow data source (100ms per item).
func loadItem(id int) Item {
	time.Sleep(100 * time.Millisecond)
	return Item{ID: id}
}

// Design 1: load everything, serialize everything, then send.
func bulkHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]Item, 0, 5)
	for i := 1; i <= 5; i++ {
		items = append(items, loadItem(i))
	}
	body, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body) // the first byte leaves only now
}

// Design 2: send each item (NDJSON) as soon as it is ready.
func streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	enc := json.NewEncoder(w)
	for i := 1; i <= 5; i++ {
		if err := enc.Encode(loadItem(i)); err != nil {
			return
		}
		flusher.Flush() // push this item to the client now
	}
}
```

</div>

<p dir="rtl" align="right">نتیجه‌ی اجرای واقعی این برنامه (با یک client که با <code>httptrace</code> اولین بایت را اندازه می‌گیرد):</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">طراحی</th>
      <th align="right">TTFB</th>
      <th align="right">TTLB (Response Time)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">bulk (طراحی اول)</td>
      <td align="right">حدود 503ms</td>
      <td align="right">حدود 503ms</td>
    </tr>
    <tr>
      <td align="right">stream (طراحی دوم)</td>
      <td align="right">حدود 100ms</td>
      <td align="right">حدود 503ms</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">نکته‌ی مهم: <strong>Response Time در هر دو یکی است، ولی TTFB پنج برابر کمتر شده است.</strong> کاربر در حالت دوم بعد از 100ms اولین آیتم را می‌بیند. شرایط این بهبود:</p>

<ul dir="rtl">
  <li>کلاینت باید بتواند داده را تکه‌تکه پردازش کند (مثلاً NDJSON یا SSE). برای یک JSON آرایه‌ای معمولی که کلاینت باید کامل بگیرد، این روش فایده‌ای ندارد.</li>
  <li>خطا بعد از ارسال اولین بایت سخت‌تر مدیریت می‌شود، چون status code قبلاً ارسال شده است (معمولاً 200) و دیگر نمی‌شود آن را عوض کرد.</li>
  <li>فایل کامل و قابل اجرای این مثال (شامل middleware ای که زمان اولین Write را در سمت سرور ثبت می‌کند و client اندازه‌گیر) در <a href="examples/latency-ttfb/main.go">examples/latency-ttfb/main.go</a> است. توضیح اجرا و آزمایش‌های پیشنهادی در <a href="examples/latency-ttfb/README.md">README همان پوشه</a> آمده است.</li>
</ul>

<p dir="rtl" align="right">برای اجرا، از ریشه‌ی ریپو این دستور را بزنید:</p>

<div dir="ltr" align="left">

```bash
go run ./02-01-performance/examples/latency-ttfb
```

</div>

<h3 dir="rtl" align="right">جمع‌بندی:</h3>

<p dir="rtl" align="right">اگر بخواهیم تفاوت اصلی این سه را در یک جمله خلاصه کنیم:</p>

<ul dir="rtl">
  <li>تاخیر - Latency: تاخیر در یک مسیر یا بخش از سیستم، مفهوم عمومی‌تر (برای تحلیل تاخیر در هر بخش از سیستم استفاده می‌شود)</li>
  <li>زمان پاسخ - Response Time: کل زمان کامل شدن پاسخ از دید کلاینت (برای فهمیدن اینکه کل پاسخ چه زمانی کامل می‌شود استفاده می‌شود)</li>
  <li>زمان TTFB: زمان رسیدن اولین byte پاسخ به کلاینت (برای فهمیدن اینکه پاسخ‌دهی چه زمانی شروع می‌شود استفاده می‌شود)</li>
</ul>

<p dir="rtl" align="right">فرض کن این API را داریم:</p>

<div dir="ltr" align="left">

```
GET /feed
```

</div>

<p dir="rtl" align="right">زمان‌ها این شکلی‌اند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">مرحله</th>
      <th align="right">مدت همین مرحله</th>
      <th align="right">زمان تجمعی از شروع</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">Client sends request</td>
      <td align="right">-</td>
      <td align="right">0ms</td>
    </tr>
    <tr>
      <td align="right">Network to backend</td>
      <td align="right">20ms</td>
      <td align="right">20ms</td>
    </tr>
    <tr>
      <td align="right">Load balancer</td>
      <td align="right">5ms</td>
      <td align="right">25ms</td>
    </tr>
    <tr>
      <td align="right">Auth</td>
      <td align="right">15ms</td>
      <td align="right">40ms</td>
    </tr>
    <tr>
      <td align="right">App logic</td>
      <td align="right">60ms</td>
      <td align="right">100ms</td>
    </tr>
    <tr>
      <td align="right">DB query</td>
      <td align="right">100ms</td>
      <td align="right">200ms</td>
    </tr>
    <tr>
      <td align="right">Backend starts response</td>
      <td align="right">-</td>
      <td align="right">200ms</td>
    </tr>
    <tr>
      <td align="right">First byte reaches client</td>
      <td align="right">20ms</td>
      <td align="right">220ms</td>
    </tr>
    <tr>
      <td align="right">Response body download</td>
      <td align="right">300ms</td>
      <td align="right">520ms</td>
    </tr>
    <tr>
      <td align="right">Client receives full response</td>
      <td align="right">-</td>
      <td align="right">520ms</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">حالا:</p>

<ul dir="rtl">
  <li><strong>TTFB = 220ms</strong></li>
  <li><strong>Response Time = 520ms</strong></li>
  <li>Download Time = 520 - 220 = 300ms</li>
</ul>

<p dir="rtl" align="right">و latency بسته به context می‌تواند این‌ها باشد:</p>

<ul dir="rtl">
  <li>Network latency: 20ms رفت و 20ms برگشت (برای اولین بایت)، یعنی 40ms</li>
  <li>DB latency: 100ms</li>
  <li>Server-side latency (از رسیدن request به backend تا شروع پاسخ): 5 + 15 + 60 + 100 = 180ms</li>
  <li>Client-observed latency: 520ms</li>
</ul>

<p dir="rtl" align="right">پس وقتی کسی فقط می‌گوید latency، باید بپرسیم: Latency از دید کجا؟ کلاینت؟ سرور؟ دیتابیس؟ شبکه؟</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> نقشه‌ی همه‌ی اصطلاح‌ها روی مسیر یک request&rlm;:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">بازه</th>
      <th align="right">نام رایج</th>
      <th align="right">معمولاً چه کسی/چه ابزاری آن را می‌بیند</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">DNS و TCP و TLS</td>
      <td align="right">Network Latency</td>
      <td align="right">client، CDN، synthetic monitoring</td>
    </tr>
    <tr>
      <td align="right">انتظار در صف</td>
      <td align="right">Queue Time (Application Latency)</td>
      <td align="right">فقط با instrument صریح در سرور یا gateway</td>
    </tr>
    <tr>
      <td align="right">منطق برنامه، DB، Cache، RPC</td>
      <td align="right">Application و Internal Latency</td>
      <td align="right">metrics و tracing داخل سرویس</td>
    </tr>
    <tr>
      <td align="right">آماده‌سازی تا اولین Write</td>
      <td align="right">Serialization و Time Before First Write</td>
      <td align="right">metrics داخل سرویس</td>
    </tr>
    <tr>
      <td align="right">تا اولین بایت</td>
      <td align="right">TTFB</td>
      <td align="right">client و gateway</td>
    </tr>
    <tr>
      <td align="right">بعد از اولین بایت تا آخرین بایت</td>
      <td align="right">Download Time</td>
      <td align="right">client</td>
    </tr>
    <tr>
      <td align="right">کل مسیر</td>
      <td align="right">Response Time</td>
      <td align="right">client یا gateway یا server (سه عدد متفاوت)</td>
    </tr>
  </tbody>
</table>

<h3 dir="rtl" align="right">کی از کدام استفاده کنیم؟</h3>

<p dir="rtl" align="right">وقتی درباره تجربه کاربر یا API حرف می‌زنیم: Response Time&rlm;. مثلاً در Design Doc بنویسی:</p>

<p dir="rtl" align="right">P95 response time for GET /feed should be under 800ms&rlm;.</p>

<p dir="rtl" align="right">این برای user-facing APIها خیلی مناسب است. مثال:</p>

<ul dir="rtl">
  <li>Login API response time P95 کمتر از 300ms</li>
  <li>Checkout API response time P95 کمتر از 700ms</li>
  <li>Search API response time P95 کمتر از 500ms</li>
</ul>

<p dir="rtl" align="right">اینجا مهم است کاربر چقدر منتظر می‌ماند.</p>

<p dir="rtl" align="right">وقتی می‌خواهیم bottleneck پیدا کنیم: <strong>Latency</strong></p>

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

<p dir="rtl" align="right">اگر اولین token بعد از 500ms بیاید، کاربر حس می‌کند سیستم responsive است، حتی اگر کل جواب 10 ثانیه طول بکشد.</p>

<p dir="rtl" align="right">نکته‌ی کلیدی این است که TTFB بخشی از Response Time است؛ زیرا Response Time از شروع request تا دریافت آخرین بایت ادامه دارد، در حالی که TTFB فقط تا دریافت اولین بایت را اندازه‌گیری می‌کند.</p>

<p dir="rtl" align="right">همچنین Latency مفهومی عمومی‌تر از هر دو است و می‌تواند برای کل request یا هر بخش داخلی سیستم، مانند شبکه، دیتابیس، کش، RPC، صف و business logic استفاده شود.</p>

<p dir="rtl" align="right">بنابراین در یک Design Doc یا System Design Discussion، بهتر است معیارها به‌صورت دقیق و قابل اندازه‌گیری بیان شوند. برای مثال:</p>

<ul dir="rtl">
  <li>For read APIs, P95 response time should be under 300ms&rlm;.</li>
  <li>For payment authorization, P99 response time should be under 1s&rlm;.</li>
  <li>The main latency contributors should be tracked separately: DB latency, cache latency, RPC latency, and queue time&rlm;.</li>
  <li>For streaming or large responses, TTFB should be under 200ms and TTLB should be monitored separately&rlm;.</li>
</ul>

<h2 dir="rtl" align="right">چرا Latency زیر بار ناگهان بالا می‌رود؟</h2>

<p dir="rtl" align="right">تا اینجا Latency را «اجزای مسیر یک request» دیدیم. ولی یک عامل دیگر هم هست: <strong>صف</strong>. وقتی منابع (CPU، worker، connection) به ظرفیت‌شان نزدیک می‌شوند، request ها در صف می‌مانند و Latency غیرخطی بالا می‌رود.</p>

<p dir="rtl" align="right">یک مدل ساده (صف تک‌سرور با ورود تصادفی، که به آن M/M/1 می‌گویند) این رابطه را می‌دهد:</p>

<div dir="ltr" align="left">

```
Response Time = Service Time / (1 - Utilization)
```

</div>

<p dir="rtl" align="right">فرض کنید Service Time (زمان خود پردازش) برابر 20ms باشد:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">میزان استفاده از منبع</th>
      <th align="right">Response Time</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">50%</td>
      <td align="right">40ms</td>
    </tr>
    <tr>
      <td align="right">80%</td>
      <td align="right">100ms</td>
    </tr>
    <tr>
      <td align="right">90%</td>
      <td align="right">200ms</td>
    </tr>
    <tr>
      <td align="right">95%</td>
      <td align="right">400ms</td>
    </tr>
    <tr>
      <td align="right">99%</td>
      <td align="right">2,000ms</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">یعنی از 50% به 90% استفاده، Latency پنج برابر شده است و از 90% به 99%، ده برابر دیگر. این یک مدل ساده‌شده است (سیستم واقعی چند worker دارد و رفتارش دقیقاً این نیست)، ولی الگو واقعی است:</p>

<ul dir="rtl">
  <li><strong>سیستم قبل از 100% استفاده، به‌شدت کند می‌شود.</strong> برای همین ظرفیت را برای 100% طراحی نمی‌کنیم و <strong>headroom</strong> نگه می‌داریم (معمولاً زیر 70 تا 80% در حالت عادی).</li>
  <li>همین مکانیزم توضیح می‌دهد چرا مثال «20ms در تست local ولی چند ثانیه زیر بار» در <a href="performance.md">Performance</a> اتفاق می‌افتد.</li>
  <li>ظرفیت و headroom در <a href="throughput-and-capacity.md">Throughput &amp; Capacity</a> بیشتر بررسی می‌شود.</li>
</ul>

<h2 dir="rtl" align="right">ادامه‌ی مسیر</h2>

<ul dir="rtl">
  <li><a href="percentiles.md">Percentile-based Thinking</a>: Latency را با P95 و P99 چطور بخوانیم</li>
  <li><a href="measurement-tools.md">ابزارهای اندازه‌گیری</a>: هر کدام از این زمان‌ها را با چه ابزاری اندازه بگیریم</li>
  <li><a href="endpoint-analysis-checklist.md">چک‌لیست تحلیل Endpoint</a>: از روی این اعداد چطور مشکل را پیدا کنیم</li>
  <li><a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a>: کنترل بخشی از Queue Time</li>
</ul>

</div>
