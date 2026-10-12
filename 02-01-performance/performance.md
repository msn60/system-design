# Performance

how quickly and efficiently a system responds to user requests.

**Performance: product owner’s view**

Slow performance can lead to user frustration, abandonment, and lost revenue.

- **How quickly should the system respond to** a user action to feel instantaneous?
- What is the **expected number of users** the system should handle without slowing down?
- Are there specific critical user journeys that must have a high level of performance?

<div dir="rtl" align="right">

<p dir="rtl" align="right">Performance یعنی: سیستم با چه سرعت، پایداری و بهره‌وری به درخواست‌های کاربر پاسخ می‌دهد، در شرایط مختلف و تحت بارهای متفاوت.</p>

<h2 dir="rtl" align="right">آیا performance همان سرعت است؟</h2>

<p dir="rtl" align="right">اما در عمل Performance فقط «سرعت» نیست. چهار لایه مهم دارد:</p>

<h3 dir="rtl" align="right">1.1 عملکرد ادراک شده - Perceived Performance</h3>

<p dir="rtl" align="right">چیزی که کاربر حس می‌کند، نه لزوماً آنچه واقعاً اتفاق افتاده</p>

<p dir="rtl" align="right">آیا سیستم «فوری» به نظر می‌رسد؟ آیا کاربر حس می‌کند سیستم زنده است؟ آیا در حین پردازش بازخورد می‌گیرد؟</p>

<p dir="rtl" align="right">گاهی سیستم کند است، ولی چون Feedback مناسب دارد، کاربر آن را سریع تلقی می‌کند.</p>

<h3 dir="rtl" align="right">1.2 عملکرد واقعی - Actual Performance</h3>

<p dir="rtl" align="right">اعداد و واقعیت فنی:</p>

<ul dir="rtl">
  <li>زمان پاسخ</li>
  <li>توان پردازش</li>
  <li>مصرف منابع</li>
  <li>تحمل بار همزمان</li>
</ul>

<h3 dir="rtl" align="right">1.3 پایداری عملکرد - Consistency</h3>

<ul dir="rtl">
  <li>آیا همیشه سریع است؟</li>
  <li>یا گاهی سریع و گاهی خیلی کند؟</li>
</ul>

<p dir="rtl" align="right"><strong>ناپایداری از کندی خطرناک‌تر است.</strong></p>

<h3 dir="rtl" align="right">1.4 اثر مقیاس پذیری - Scalability Impact</h3>

<ul dir="rtl">
  <li>وقتی کاربران زیاد می‌شوند چه اتفاقی می‌افتد؟</li>
  <li>آیا تجربه افت می‌کند یا ثابت می‌ماند؟</li>
</ul>

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

<p dir="rtl" align="right">عموما این مسیرها باید سریع‌تر از میانگین باشند و Animation + Skeleton + Feedback داشته باشند. ضمن اینکه باید Error Rate نزدیک صفر داشته باشند</p>

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

<h2 dir="rtl" align="right">نقش Performance در System Design</h2>

<p dir="rtl" align="right">در سیستم دیزاین، Performance فقط مسئله Backend نیست.</p>

<p dir="rtl" align="right"><strong>تصمیم های طراحی یا Design Decisions</strong> مؤثر بر <strong>Performance&rlm;:</strong></p>

<ul dir="rtl">
  <li>حجم تصاویر</li>
  <li>Animation Duration</li>
  <li>Use of Skeletons</li>
  <li>Pagination vs Infinite Scroll</li>
  <li>Client-side vs Server-side rendering</li>
  <li>Preload / Prefetch</li>
  <li>Component Complexity</li>
</ul>

<p dir="rtl" align="right">از دید بک‌اند، Performance یعنی: توانایی سیستم در پاسخ‌دهی سریع، پایدار و قابل پیش‌بینی به درخواست‌ها، تحت بارهای مختلف، با حداقل مصرف منابع و بدون افت کیفیت. Performance یعنی سیستم بتواند: درخواست‌ها را <strong>سریع،</strong> <strong>پایدار،</strong> <strong>قابل پیش‌بینی</strong> و <strong>با مصرف منطقی منابع</strong> پاسخ دهد؛ حتی وقتی بار سیستم تغییر می‌کند.</p>

<p dir="rtl" align="right">چهار ستون اصلی دارد:</p>

<ul dir="rtl">
  <li>Latency – هر درخواست چقدر طول می‌کشد</li>
  <li>Throughput – چند درخواست در واحد زمان پردازش می‌شود</li>
  <li>Concurrency – چند درخواست همزمان بدون افت کیفیت</li>
  <li>Stability under Load – رفتار سیستم وقتی تحت فشار است</li>
</ul>

<p dir="rtl" align="right"><strong>در نظر داشته باشید که Performance ≠ Speed</strong></p>

<p dir="rtl" align="right">یکی از خطاهای رایج بک‌اند: «اگه سریع باشه، Performance خوبه»</p>

<p dir="rtl" align="right">مثلاً ممکن است یک endpoint در تست local فقط در 20ms جواب بدهد، اما وقتی ۵۰۰ کاربر همزمان می‌آیند:</p>

<ul dir="rtl">
  <li>در آن زمان latency برود روی ۳ ثانیه</li>
  <li>به سرعت connection pool دیتابیس پر شود</li>
  <li>تعداد go routine ها زیاد شوند</li>
  <li>مصرف CPU بالا برود</li>
  <li>تعداد timeout زیاد شود</li>
  <li>بر اثر فشار وارد شده retry storm ایجاد شود</li>
</ul>

<p dir="rtl" align="right">اینجا endpoint در حالت عادی سریع بوده، ولی performance واقعی ندارد.</p>

<p dir="rtl" align="right">در واقع Performance ترکیب این‌هاست:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">عامل</th>
      <th align="right">توضیح</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">Latency</td>
      <td align="right">زمان پاسخ</td>
    </tr>
    <tr>
      <td align="right">Variance</td>
      <td align="right">نوسان زمان پاسخ</td>
    </tr>
    <tr>
      <td align="right">Saturation</td>
      <td align="right">پر شدن منابع</td>
    </tr>
    <tr>
      <td align="right">Backpressure</td>
      <td align="right">واکنش به فشار</td>
    </tr>
    <tr>
      <td align="right">Failure Behavior</td>
      <td align="right">رفتار در خرابی</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">سیستمی که گاهی 50ms و گاهی 5s جواب می‌دهد، بدتر از سیستمی است که همیشه 300ms جواب می‌دهد.</p>

</div>
