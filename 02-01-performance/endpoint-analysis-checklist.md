<div dir="rtl" align="right">

<h1 dir="rtl" align="right">چک‌لیست تحلیل Performance برای هر Endpoint</h1>

<p dir="rtl" align="right">برای هر endpoint باید بتوانید مسیر کامل request را تحلیل کنید. هدف این نیست که فقط: response time = 700ms را ببینید. هدف این است که بدانید: <strong>این 700ms دقیقاً کجا مصرف شده است</strong>.</p>

<h2 dir="rtl" align="right">سوال‌های کلیدی برای هر Endpoint</h2>

<h3 dir="rtl" align="right">1. TTFB چقدر است؟</h3>

<p dir="rtl" align="right">TTFB (Time To First Byte) نشان می‌دهد: سیستم چقدر طول می‌کشد تا اولین بخش response را ارسال کند. TTFB بالا معمولاً نشان‌دهنده یکی از این موارد است:</p>

<ul dir="rtl">
  <li>queue delay</li>
  <li>DB wait</li>
  <li>dependency latency</li>
  <li>connection pool exhaustion</li>
  <li>slow serialization</li>
</ul>

<p dir="rtl" align="right">مثلا  اگر: TTFB = 900ms و Total Response Time = 950ms باشد، مشکل معمولاً قبل از شروع response است. اما اگر: TTFB = 50ms و Total Response Time = 3s احتمالاً:</p>

<ul dir="rtl">
  <li>payload بزرگ است</li>
  <li>streaming کند است</li>
  <li>client/network کند است</li>
</ul>

<h2 dir="rtl" align="right">2. Response Time p95/p99 چقدر است؟</h2>

<p dir="rtl" align="right">Average latency تقریباً همیشه ناکافی است. در distributed systems: همواره tail latency تعیین‌کننده تجربه واقعی کاربران است. زیرا اگر 99% requests خوب باشند، 1% بدترین request ها می‌توانند کل سیستم را degrade کنند.</p>

<h3 dir="rtl" align="right">دلایل رایج P99 بالا</h3>

<ul dir="rtl">
  <li>GC pauses</li>
  <li>lock contention</li>
  <li>slow queries</li>
  <li>dependency spikes</li>
  <li>queue buildup</li>
  <li>noisy neighbors</li>
  <li>network retransmission</li>
</ul>

<h2 dir="rtl" align="right">3. آیا Queue وجود دارد؟</h2>

<p dir="rtl" align="right">تقریباً هر queue می‌تواند: buffer مفید یا محل پنهان شدن latency باشد. Queue خوب باید:</p>

<ul dir="rtl">
  <li>bounded باشد</li>
  <li>observable باشد</li>
  <li>backpressure ایجاد کند</li>
</ul>

<p dir="rtl" align="right">Queue خطرناک:</p>

<ul dir="rtl">
  <li>Queue نامحدود:</li>
</ul>

</div>

```
traffic spike → queue growth → memory growth → latency explosion → timeout storm 
```

<div dir="rtl" align="right">

<h2 dir="rtl" align="right">4. آیا Database Bottleneck است؟</h2>

<p dir="rtl" align="right">Database معمولاً اولین bottleneck واقعی سیستم‌های backend است. علائم DB Bottleneck عبارتند از:</p>

<ul dir="rtl">
  <li>query latency بالا</li>
  <li>lock wait</li>
  <li>connection pool exhaustion</li>
  <li>replication lag</li>
  <li>sequential scans</li>
  <li>high IO wait</li>
</ul>

<p dir="rtl" align="right">سوال‌های مهم که می توانیم برای دیتابیس از خود بپرسیم:</p>

<ul dir="rtl">
  <li>آیا query ها index مناسب دارند؟</li>
  <li>آیا pagination وجود دارد؟</li>
  <li>آیا N+1 query رخ می‌دهد؟</li>
  <li>آیا fan-out query داریم؟</li>
  <li>آیا hot rows وجود دارد؟</li>
</ul>

<h2 dir="rtl" align="right">5. آیا External Dependency وجود دارد؟</h2>

<p dir="rtl" align="right">هر dependency خارجی:</p>

<ul dir="rtl">
  <li>latency</li>
  <li>failure</li>
  <li>retry</li>
  <li>unpredictability</li>
</ul>

<p dir="rtl" align="right">را وارد سیستم می‌کند. یک قانون مهم این است که هر dependency خارجی،SLO شما را ضعیف‌تر می‌کند. به عنوان مثال اگر endpoint شما: auth service و payment gateway و recommendation service را call کند، latency نهایی مجموع behavior همه آن‌هاست.</p>

<h2 dir="rtl" align="right">6. سیستم Under Load چگونه رفتار می‌کند؟</h2>

<p dir="rtl" align="right">بسیاری از سیستم‌ها در load کم سالم‌اند، اما تحت pressure collapse می‌کنند. سوال‌های مهم که در این مورد باید از خود بپرسیم:</p>

<ul dir="rtl">
  <li>آیا queue build up می‌شود؟</li>
  <li>آیا retry storm رخ می‌دهد؟</li>
  <li>آیا auto scaling به‌موقع عمل می‌کند؟</li>
  <li>آیا connection pool اشباع می‌شود؟</li>
  <li>آیا tail latency انفجاری می‌شود؟</li>
  <li>آیا backpressure وجود دارد؟</li>
</ul>

</div>
