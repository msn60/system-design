# Performance

How quickly and efficiently a system responds to user requests.

**Performance: product owner's view**

Slow performance can lead to user frustration, abandonment, and lost revenue.

- **How quickly should the system respond** to a user action to feel instantaneous?
- What is the **expected number of users** the system should handle without slowing down?
- Are there specific critical user journeys that must have a high level of performance?

کارایی یا Performance یعنی: سیستم با چه سرعت، پایداری و بهره‌وری به درخواست‌های کاربر پاسخ می‌دهد، در شرایط مختلف و تحت بارهای متفاوت.

<div dir="rtl" align="right">

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">سه سؤال بالا فقط وقتی مفیدند که جواب‌شان <strong>عدد</strong> باشد. چند عدد مرجع:</p>

<ul dir="rtl">
  <li><strong>«فوری» یعنی چقدر؟</strong> بر اساس تحقیقات معروف Jakob Nielsen درباره‌ی زمان پاسخ: حدود 0.1 ثانیه یعنی کاربر آن را «فوری» حس می‌کند، حدود 1 ثانیه یعنی جریان فکر کاربر قطع نمی‌شود، و حدود 10 ثانیه حد تحمل توجه کاربر است. بعد از آن باید حتماً نشانه‌ی پیشرفت نشان دهیم.</li>
  <li>برای صفحه‌های وب، Google برای Core Web Vitals این مقادیر را «خوب» می‌داند: LCP کمتر از 2.5 ثانیه، INP کمتر از 200ms و CLS کمتر از 0.1. نحوه‌ی اندازه‌گیری‌شان در <a href="measurement-tools.md">ابزارهای اندازه‌گیری</a> آمده است.</li>
  <li><strong>«چند کاربر؟»</strong> باید به RPS در زمان peak تبدیل شود، چون سیستم با «کاربر» کار نمی‌کند، با request کار می‌کند (<a href="throughput-and-capacity.md">Throughput &amp; Capacity</a>).</li>
  <li><strong>«کدام مسیرها مهم‌اند؟»</strong> جواب آن در بخش «نقش Performance در User Journey» همین فایل آمده است.</li>
</ul>

<p dir="rtl" align="right">این عددها راهنما هستند، نه قانون. عدد نهایی هر سیستم باید از نیاز بیزینس بیاید و به شکل NFR قابل‌اندازه‌گیری نوشته شود (<a href="../01-non-functional-requirements/non-functional-requirements.md">Non-Functional Requirements</a>).</p>

<h2 dir="rtl" align="right">آیا performance همان سرعت است؟</h2>

<p dir="rtl" align="right">اما در عمل Performance فقط «سرعت» نیست. چهار لایه مهم دارد:</p>

<h3 dir="rtl" align="right">عملکرد ادراک شده - Perceived Performance</h3>

<p dir="rtl" align="right">چیزی که کاربر حس می‌کند، نه لزوماً آنچه واقعاً اتفاق افتاده</p>

<p dir="rtl" align="right">آیا سیستم «فوری» به نظر می‌رسد؟ آیا کاربر حس می‌کند سیستم زنده است؟ آیا در حین پردازش بازخورد می‌گیرد؟</p>

<p dir="rtl" align="right">گاهی سیستم کند است، ولی چون Feedback مناسب دارد، کاربر آن را سریع تلقی می‌کند.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>نمونه‌های رایج: Skeleton به‌جای صفحه‌ی خالی، Progress bar، و <strong>Optimistic UI</strong> (مثلاً دکمه‌ی لایک همان لحظه تغییر می‌کند و تأیید سرور بعداً می‌رسد).</li>
  <li>محدودیت: Perceived Performance جایگزین Actual Performance نمی‌شود. Skeleton جلوی timeout را نمی‌گیرد. در Optimistic UI هم اگر سرور خطا بدهد، باید تغییر را به کاربر برگردانیم (rollback). پس این لایه برای پنهان کردن یک سیستم خراب نیست، برای بهتر کردن تجربه‌ی یک سیستم سالم است.</li>
</ul>

<h3 dir="rtl" align="right">عملکرد واقعی - Actual Performance</h3>

<p dir="rtl" align="right">اعداد و واقعیت فنی:</p>

<ul dir="rtl">
  <li>زمان پاسخ</li>
  <li>توان پردازش</li>
  <li>مصرف منابع</li>
  <li>تحمل بار همزمان</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> هر کدام از این چهار مورد یک متریک و یک فایل مرتبط دارد:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">مورد</th>
      <th align="right">متریک</th>
      <th align="right">بخش مرتبط</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">زمان پاسخ</td>
      <td align="right">Latency (P50 و P95 و P99)</td>
      <td align="right"><a href="latency.md">Latency</a>، <a href="percentiles.md">Percentile-based Thinking</a></td>
    </tr>
    <tr>
      <td align="right">توان پردازش</td>
      <td align="right">Throughput (RPS)</td>
      <td align="right"><a href="throughput-and-capacity.md">Throughput &amp; Capacity</a></td>
    </tr>
    <tr>
      <td align="right">مصرف منابع</td>
      <td align="right">CPU، حافظه، Disk و Network I/O</td>
      <td align="right"><a href="advanced-metrics.md">Advanced Metrics</a></td>
    </tr>
    <tr>
      <td align="right">تحمل بار همزمان</td>
      <td align="right">Concurrency (تعداد request در حال پردازش)</td>
      <td align="right">بخش «چهار ستون» همین فایل</td>
    </tr>
  </tbody>
</table>

<h3 dir="rtl" align="right">پایداری عملکرد - Consistency</h3>

<ul dir="rtl">
  <li>آیا همیشه سریع است؟</li>
  <li>یا گاهی سریع و گاهی خیلی کند؟</li>
</ul>

<p dir="rtl" align="right"><strong>ناپایداری از کندی خطرناک‌تر است.</strong></p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>«Consistency» در این بخش یعنی <strong>پایدار بودن زمان پاسخ</strong>. این با Consistency داده (حرف C در CAP) فرق دارد. برای جلوگیری از اشتباه، می‌توانید در ذهن‌تان آن را «Predictability» بنامید.</li>
  <li>جمله‌ی «ناپایداری از کندی خطرناک‌تر است» در حالت کلی درست نیست و فقط وقتی درست است که نوسان باعث timeout، retry و صف شود. توضیح این زنجیره در <a href="backend-performance-summary.md">جمع‌بندی Backend Performance</a> (بخش Predictability) آمده است، و در انتهای همین فایل یک مثال عددی دارد.</li>
</ul>

<h3 dir="rtl" align="right">اثر مقیاس پذیری - Scalability Impact</h3>

<ul dir="rtl">
  <li>وقتی کاربران زیاد می‌شوند چه اتفاقی می‌افتد؟</li>
  <li>آیا تجربه افت می‌کند یا ثابت می‌ماند؟</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> تفاوت Performance و Scalability: Performance یعنی سیستم <strong>با بار فعلی</strong> چقدر خوب جواب می‌دهد. Scalability یعنی وقتی بار <strong>رشد کند</strong>، این کیفیت چقدر حفظ می‌شود. سیستمی که برای 100 کاربر سریع است و برای 10,000 کاربر نه، Performance خوب ولی Scalability ضعیف دارد. جزئیات Scale در <a href="../00-challenges-and-goals/challenges.md">چالش 1</a> آمده است.</p>

<h2 dir="rtl" align="right">نقش Performance در User Journey</h2>

<p dir="rtl" align="right">همه بخش‌ها نیاز به Performance یکسان ندارند.</p>

<h3 dir="rtl" align="right">Critical User Journeys</h3>

<p dir="rtl" align="right">مسیرهایی که اگر کند باشند، کاربر می‌رود:</p>

<ul dir="rtl">
  <li>Login / Sign Up</li>
  <li>Search</li>
  <li>Checkout / Payment</li>
  <li>Navigation اولیه</li>
  <li>First Content Load</li>
</ul>

<p dir="rtl" align="right">عموماً این مسیرها باید سریع‌تر از میانگین باشند و Animation + Skeleton + Feedback داشته باشند. ضمن اینکه باید Error Rate نزدیک صفر داشته باشند.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>«سریع‌تر از میانگین» در عمل یعنی هدف سخت‌گیرانه‌تر، مثلاً P99 کمتر و Availability بالاتر برای این مسیرها (نمونه در <a href="sli-slo-sla.md">SLI / SLO / SLA</a>).</li>
  <li>فهرست مسیرهای بحرانی به محصول بستگی دارد. Search برای یک فروشگاه اینترنتی بحرانی است، ولی برای یک اپلیکیشن بانکی شاید «انتقال وجه» و «ورود» بحرانی‌تر باشند. یک روش ساده برای پیدا کردنشان: «اگر این مسیر کند شود، کدام عدد بیزینسی (درآمد، ثبت‌نام، نگه‌داشت کاربر) مستقیماً افت می‌کند؟»</li>
</ul>

<h3 dir="rtl" align="right">Non-Critical Flows</h3>

<p dir="rtl" align="right">مثل:</p>

<ul dir="rtl">
  <li>Analytics</li>
  <li>History</li>
  <li>Settings</li>
  <li>Export ها</li>
</ul>

<p dir="rtl" align="right">می‌توانند:</p>

<ul dir="rtl">
  <li>به صورت Lazy Load شوند</li>
  <li>به صورت Async باشند</li>
  <li>دارای Feedback طولانی‌تری باشند</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li>Lazy Load بیشتر اصطلاح سمت کلاینت است. معادل‌های سمت بک‌اند: پردازش async با صف و worker، پاسخ <code>202 Accepted</code> با بررسی وضعیت بعداً، و اولویت پایین‌تر.</li>
  <li><strong>نکته‌ی مهم طراحی:</strong> مسیر غیربحرانی نباید منابع مسیر بحرانی را بگیرد. مثلاً اگر Export یک میلیون ردیف و Checkout از یک connection pool و یک worker pool مشترک استفاده کنند، یک Export سنگین می‌تواند Checkout را کند کند. راه‌حل: pool یا صف جدا (bulkhead) و در زمان فشار، اول از همه همین درخواست‌ها را رد یا عقب بیندازیم (load shedding). این مفاهیم در <a href="../00-challenges-and-goals/challenges.md">چالش 6</a> معرفی شده‌اند و تکنیک‌شان در <a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a> می‌آید.</li>
</ul>

<h2 dir="rtl" align="right">نقش Performance در System Design</h2>

<p dir="rtl" align="right">در سیستم دیزاین، Performance فقط مسئله Backend نیست.</p>

<p dir="rtl" align="right"><strong>تصمیم‌های طراحی یا Design Decisions</strong> مؤثر بر <strong>Performance&rlm;:</strong></p>

<ul dir="rtl">
  <li>حجم تصاویر</li>
  <li>Animation Duration</li>
  <li>Use of Skeletons</li>
  <li>Pagination vs Infinite Scroll</li>
  <li>Client-side vs Server-side rendering</li>
  <li>Preload / Prefetch</li>
  <li>Component Complexity</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> این فهرست بیشتر تصمیم‌های سمت کلاینت و UX است. تصمیم Pagination در برابر Infinite Scroll روی بک‌اند هم اثر مستقیم دارد (مثلاً نوع pagination و هزینه‌ی Count) و در <a href="../02-04-batching-and-pagination/pagination.md">Pagination</a> بررسی شده است.</p>

<p dir="rtl" align="right">از دید بک‌اند، Performance یعنی سیستم بتواند درخواست‌ها را <strong>سریع،</strong> <strong>پایدار،</strong> <strong>قابل پیش‌بینی</strong> و <strong>با مصرف منطقی منابع</strong> پاسخ دهد؛ حتی وقتی بار سیستم تغییر می‌کند، و بدون افت کیفیت.</p>

<p dir="rtl" align="right">چهار ستون اصلی دارد:</p>

<ul dir="rtl">
  <li>Latency – هر درخواست چقدر طول می‌کشد</li>
  <li>Throughput – چند درخواست در واحد زمان پردازش می‌شود</li>
  <li>Concurrency – چند درخواست همزمان بدون افت کیفیت</li>
  <li>Stability under Load – رفتار سیستم وقتی تحت فشار است</li>
</ul>

<p dir="rtl" align="right"><strong>نکات تکمیلی</strong></p>

<ul dir="rtl">
  <li><strong>Concurrency</strong> در این مجموعه تا الان تعریف نشده بود: تعداد request هایی که <strong>همین لحظه</strong> در حال پردازش‌اند. رابطه‌ی ساده‌ی آن با دو ستون دیگر (Little's Law):</li>
</ul>

<div dir="ltr" align="left">

```text
Concurrency = Throughput × Latency
```

</div>

<p dir="rtl" align="right">مثال: اگر 500 RPS داشته باشیم و Latency هر request برابر 20ms باشد، همزمان حدود 500 × 0.02 = 10 request در سیستم هستند. اگر Latency به 3s برسد و ورودی ثابت بماند، این عدد می‌شود 500 × 3 = 1,500. یعنی کندی، خودش فشار را 150 برابر می‌کند. این همان مکانیزم مثال بعدی است.</p>

<ul dir="rtl">
  <li><strong>Stability under Load</strong> هر دو سمت بحث را پوشش می‌دهد: Saturation و Backpressure و رفتار در خرابی. توضیح کامل‌شان در <a href="key-principles.md">اصول کلیدی</a> (اصل 3 و 4 و 5) و <a href="backend-performance-summary.md">جمع‌بندی Backend Performance</a> آمده است. در این فایل فقط به‌صورت خلاصه می‌آیند تا تکرار نشوند.</li>
</ul>

<p dir="rtl" align="right"><strong>در نظر داشته باشید که Performance ≠ Speed</strong></p>

<p dir="rtl" align="right">یکی از خطاهای رایج بک‌اند: «اگه سریع باشه، Performance خوبه»</p>

<p dir="rtl" align="right">مثلاً ممکن است یک endpoint در تست local فقط در 20ms جواب بدهد، اما وقتی 500 کاربر همزمان می‌آیند:</p>

<ul dir="rtl">
  <li>در آن زمان latency برود روی 3 ثانیه</li>
  <li>به سرعت connection pool دیتابیس پر شود</li>
  <li>تعداد goroutine ها زیاد شوند</li>
  <li>مصرف CPU بالا برود</li>
  <li>تعداد timeout زیاد شود</li>
  <li>بر اثر فشار وارد شده retry storm ایجاد شود</li>
</ul>

<p dir="rtl" align="right">اینجا endpoint در حالت عادی سریع بوده، ولی زیر بار واقعی قابل اتکا نیست.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> این زنجیره را با عدد دنبال کنیم (همه‌ی عددها فرضی‌اند و ورودی را ثابت در نظر گرفته‌ایم):</p>

<ol dir="rtl">
  <li><strong>حالت عادی:</strong> 500 RPS و هر request حدود 20ms است. حدود 10 request همزمان در سیستم هست و یک pool با 20 connection کافی است.</li>
  <li><strong>دیتابیس شلوغ می‌شود:</strong> latency به 3s می‌رسد. حالا 500 × 3 = 1,500 request همزمان در سیستم هست.</li>
  <li><strong>pool محدود است:</strong> فقط 20 request همزمان به DB می‌رسند و بقیه در صف انتظار pool می‌مانند.</li>
  <li><strong>retry:</strong> کلاینت بعد از 2s timeout می‌شود و دوباره می‌فرستد. ورودی تا حدود 1,000 RPS بالا می‌رود.</li>
  <li><strong>حلقه‌ی معیوب:</strong> فشار بیشتر، latency بیشتر، retry بیشتر و در نهایت collapse&rlm;.</li>
</ol>

<ul dir="rtl">
  <li>در Go هر request معمولاً یک goroutine دارد. goroutine کم‌هزینه است (در حد چند کیلوبایت)، پس 1,500 تای آن به‌تنهایی مشکل نیست. مشکل واقعی <strong>چیزهایی است که این goroutine ها نگه داشته‌اند</strong>: connection، حافظه‌ی بافرها و lock&rlm;.</li>
  <li>نکته‌ی اندازه‌گیری: همین سناریو را با تست load واقعی باید دید، نه فقط با تست local (<a href="load-testing.md">Load Testing</a>).</li>
</ul>

<p dir="rtl" align="right">در واقع Performance ترکیب این‌هاست:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">عامل</th>
      <th align="right">توضیح</th>
      <th align="right">چطور اندازه می‌گیریم؟</th>
      <th align="right">بخش مرتبط</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">Latency</td>
      <td align="right">زمان پاسخ</td>
      <td align="right">P50 و P95 و P99</td>
      <td align="right"><a href="latency.md">Latency</a>، <a href="percentiles.md">Percentiles</a></td>
    </tr>
    <tr>
      <td align="right">Variance</td>
      <td align="right">نوسان زمان پاسخ</td>
      <td align="right">نسبت P99 به P50، Jitter</td>
      <td align="right"><a href="advanced-metrics.md">Advanced Metrics</a></td>
    </tr>
    <tr>
      <td align="right">Saturation</td>
      <td align="right">پر شدن منابع</td>
      <td align="right">میزان استفاده از CPU و pool، طول صف</td>
      <td align="right"><a href="key-principles.md">اصول کلیدی</a> (اصل 3)</td>
    </tr>
    <tr>
      <td align="right">Backpressure</td>
      <td align="right">واکنش به فشار</td>
      <td align="right">تعداد درخواست رد شده (429 و 503)، طول صف محدودشده</td>
      <td align="right"><a href="key-principles.md">اصول کلیدی</a> (اصل 4)</td>
    </tr>
    <tr>
      <td align="right">Failure Behavior</td>
      <td align="right">رفتار در خرابی</td>
      <td align="right">Error rate در زمان فشار، degrade یا collapse</td>
      <td align="right"><a href="backend-performance-summary.md">جمع‌بندی Backend Performance</a> (بخش 4)</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">سیستمی که گاهی 50ms و گاهی 5s جواب می‌دهد، بدتر از سیستمی است که همیشه 300ms جواب می‌دهد.</p>

<p dir="rtl" align="right"><strong>نکات تکمیلی:</strong> چرا؟ سه دلیل عملی:</p>

<ul dir="rtl">
  <li><strong>Timeout قابل تنظیم نیست.</strong> اگر timeout را 1s بگذاریم، درخواست‌های 5s fail می‌شوند. اگر 6s بگذاریم، منابع را مدت زیادی نگه می‌دارند. برای سیستم 300ms ثابت، timeout حدود 1s همه‌چیز را حل می‌کند.</li>
  <li><strong>Capacity Planning غیرممکن می‌شود.</strong> طبق Little's Law بالا، تعداد request همزمان به Latency وابسته است. وقتی Latency نوسان دارد، تعداد connection و worker لازم هم نوسان دارد.</li>
  <li><strong>Retry و صف.</strong> درخواست‌های کند باعث timeout و retry می‌شوند و فشار را بیشتر می‌کنند.</li>
</ul>

<p dir="rtl" align="right">یک قید: این جمله به <strong>سهم</strong> درخواست‌های کند بستگی دارد. اگر فقط یکی از هر 10,000 درخواست 5s باشد، ممکن است در SLO شما بگنجد. معیار تصمیم، Percentile مورد نظر SLO است (مثلاً P99)، نه مقایسه‌ی دو عدد.</p>

<h2 dir="rtl" align="right">یک مثال: بودجه‌ی Performance (Performance Budget)</h2>

<p dir="rtl" align="right">وقتی یک هدف مثل «P95 جستجو کمتر از 300ms» داریم، مفید است آن را بین اجزا تقسیم کنیم، مثل یک بودجه (اعداد فرضی):</p>

<div dir="ltr" align="left">

```text
Network + Load Balancer   30ms
Application logic         40ms
Database / Cache         130ms
External service          70ms
Serialization             30ms
-----------------------------
Total                    300ms
```

</div>

<p dir="rtl" align="right">اگر اندازه‌گیری نشان دهد query دیتابیس به‌تنهایی 200ms طول می‌کشد، بودجه از همان‌جا شکسته شده است و می‌دانیم کجا باید بهینه کنیم. این همان اصل «Latency باید Decompose شود» در <a href="key-principles.md">اصول کلیدی</a> است. چک‌لیست بررسی هر endpoint در <a href="endpoint-analysis-checklist.md">Performance Analysis Checklist</a> آمده است.</p>

---

<h2 dir="rtl" align="right">نقشه‌ی مفاهیم این فایل</h2>

<p dir="rtl" align="right">این فایل سه فهرست دارد که با هم قاطی می‌شوند. هر کدام از زاویه‌ی متفاوتی نگاه می‌کنند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">فهرست</th>
      <th align="right">از چه زاویه‌ای؟</th>
      <th align="right">موارد</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">چهار لایه</td>
      <td align="right">کاربر و محصول</td>
      <td align="right">Perceived، Actual، Consistency، Scalability Impact</td>
    </tr>
    <tr>
      <td align="right">چهار ستون</td>
      <td align="right">مهندسی (چه چیزی را اندازه می‌گیریم)</td>
      <td align="right">Latency، Throughput، Concurrency، Stability under Load</td>
    </tr>
    <tr>
      <td align="right">پنج عامل</td>
      <td align="right">مهندسی (چه چیزی را خراب می‌کند)</td>
      <td align="right">Latency، Variance، Saturation، Backpressure، Failure Behavior</td>
    </tr>
  </tbody>
</table>

<h2 dir="rtl" align="right">ادامه‌ی مسیر</h2>

<ul dir="rtl">
  <li><a href="latency.md">Latency</a>: تعریف دقیق Latency، Response Time و TTFB</li>
  <li><a href="percentiles.md">Percentile-based Thinking</a>: چرا P95 و P99 مهم‌تر از Average هستند</li>
  <li><a href="throughput-and-capacity.md">Throughput &amp; Capacity</a>: ظرفیت سیستم را چطور حساب کنیم</li>
  <li><a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a>: اولین الگوی عملی بهبود Performance</li>
</ul>

</div>
