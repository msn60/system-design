<div dir="rtl" align="right">

<h1 dir="rtl" align="right">چالش‌های اصلی System Design در مقیاس بزرگ</h1>

<p dir="rtl" align="right">بیشتر مصاحبه‌های سیستم‌دیزاین یا طراحی سرویس‌های واقعی حول همین چند محور می‌چرخند و عموما این سؤالات مطرح می شوند:</p>

<blockquote dir="rtl">
  <p dir="rtl" align="right"><strong>نکته درباره‌ی این لیست:</strong> این هشت مورد هم‌سطح نیستند. Availability و Consistency «ویژگی» سیستم هستند، Caching و Queue «تکنیک» هستند، و Real-time «نوع workload» است. در هر مورد، متن اصلی اول آمده و بعد از آن بخش «نکات تکمیلی» مسئله را کامل‌تر می‌کند. در انتهای فایل هم چند چالش مهم که در این لیست نبودند اضافه شده‌اند.</p>
</blockquote>

<h2 dir="rtl" align="right">1) چطور scale کنیم؟</h2>

<p dir="rtl" align="right">وقتی تعداد کاربران زیاد می‌شود، یک سرور دیگر پاسخگو نیست. باید بتوانیم سیستم را طوری طراحی کنیم که با افزایش کاربران، بتوانیم مدیریت scale کردن را به راحتی انجام دهیم. عموما از دو روش vertical scale (افزایش قدرت سرور) و یا horizontal scale (اضافه کردن سرورهای بیشتر) استفاده می‌شود.</p>

<p dir="rtl" align="right">مثلا در روش vertical با مشکلاتی مانند سقف داشتن، گران بودن و Single point of failure مواجه هستیم.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<ul dir="rtl">
  <li>Horizontal scale هم مشکل خودش را دارد. این روش فقط وقتی کار می‌کند که سرویس <strong>stateless</strong> باشد. اگر session، cache محلی یا فایل‌های آپلودشده روی خود سرور نگه‌داری شوند، با چند instance خراب می‌شوند. این state باید به بیرون (مثلا Redis یا Object Storage) منتقل شود.</li>
  <li>سخت‌ترین بخش scale معمولاً سرویس‌ها نیستند، <strong>دیتابیس</strong> است. اضافه کردن instance به سرویس ساده است، ولی scale کردن دیتابیس به read replica و در نهایت sharding نیاز دارد.</li>
  <li>مثال عددی: فرض کنید یک instance حدود 1,000 RPS تحمل می‌کند و ما به 10,000 RPS نیاز داریم. پس حدود 10 instance پشت Load Balancer لازم است، و در عمل کمی بیشتر برای ظرفیت اضافه (headroom&rlm;).</li>
</ul>

<h2 dir="rtl" align="right">2) چطور highly available باشیم؟</h2>

<p dir="rtl" align="right">این مورد با system uptime بیان می‌شود. به عنوان مثال گفته می شود که availability سیستم 99 درصد یا 99.9 یا 99.99 می باشد (که نشان می‌دهد سیستم چند درصد از زمان در دسترس بوده است؛ هر عدد معادل یک مقدار downtime مجاز در طول سال است). برای این مورد عموما از redundancy (مثل لود بالانس) و multi-region استفاده می‌کنیم، و برای تشخیص instance خراب، health check هم به آن‌ها اضافه می‌شود.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">هر «نُه» اضافه، downtime مجاز را حدود ده برابر کمتر می‌کند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">Availability</th>
      <th align="right">Downtime مجاز در سال</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">99%</td>
      <td align="right">حدود 3.65 روز</td>
    </tr>
    <tr>
      <td align="right">99.9%</td>
      <td align="right">حدود 8.76 ساعت</td>
    </tr>
    <tr>
      <td align="right">99.99%</td>
      <td align="right">حدود 52 دقیقه</td>
    </tr>
    <tr>
      <td align="right">99.999%</td>
      <td align="right">حدود 5 دقیقه</td>
    </tr>
  </tbody>
</table>

<ul dir="rtl">
  <li><strong>Availability در زنجیره ضرب می‌شود.</strong> اگر یک request از سه سرویس پشت‌سرهم عبور کند و هر سرویس 99.9% باشد، availability کل می‌شود 0.999 × 0.999 × 0.999 ≈ 99.7%. یعنی حدود 26 ساعت downtime در سال، به‌جای 8.76 ساعت. این موضوع مستقیماً به مورد 6 (decouple) ربط دارد.</li>
  <li>Redundancy فقط وقتی کمک می‌کند که <strong>هیچ نقطه‌ای</strong> Single Point of Failure نباشد. خود Load Balancer و دیتابیس هم باید redundant باشند.</li>
  <li>از دید کاربر، «سرور up است» کافی نیست. سرویسی که پاسخ می‌دهد ولی بسیار کند است هم عملاً در دسترس نیست. این بحث در <a href="../02-01-performance/sli-slo-sla.md">SLI / SLO / SLA</a> توضیح داده شده است.</li>
</ul>

<h2 dir="rtl" align="right">3) چطور request از دست ندهیم؟</h2>

<p dir="rtl" align="right">وقتی ترافیک خیلی زیاد می‌شود، ممکن است سرور نتواند همه درخواست‌ها را هندل کند. یکی از راه حل های این مشکل، استفاده از صف یا message broker هایی مثل Kafka یا RabbitMQ می‌باشد.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<ul dir="rtl">
  <li>صف فقط <strong>burst کوتاه‌مدت</strong> را جذب می‌کند. اگر نرخ ورودی به‌طور میانگین از نرخ پردازش بیشتر باشد، صف بی‌نهایت بزرگ می‌شود. مثال:</li>
</ul>

<div dir="ltr" align="left">

```text
ورودی:   1,000 پیام در ثانیه
پردازش:    800 پیام در ثانیه
هر ثانیه 200 پیام به صف اضافه می‌شود
بعد از 10 دقیقه: 200 × 600 = 120,000 پیام در صف مانده است
```

</div>

<ul dir="rtl">
  <li>پس در کنار صف به راه‌حل‌های دیگر هم نیاز داریم: <strong>backpressure</strong> (به تولیدکننده بگوییم آرام‌تر بفرستد)، <strong>load shedding</strong> (بعضی درخواست‌ها را عمداً رد کنیم) و <strong>rate limiting</strong>&rlm;.</li>
  <li>ریسک‌های صف: پیام ممکن است <strong>دو بار</strong> پردازش شود (پس پردازش باید idempotent باشد) و ترتیب پیام‌ها ممکن است تغییر کند.</li>
  <li>درخواستی که کاربر منتظر جواب فوری آن است را نمی‌شود فقط در صف گذاشت. در این حالت معمولاً پاسخ <code>202 Accepted</code> برمی‌گردد و کاربر وضعیت را بعداً چک می‌کند.</li>
  <li>Kafka و RabbitMQ از یک جنس نیستند: Kafka یک log توزیع‌شده است که پیام‌ها در آن می‌مانند، و RabbitMQ یک broker سنتی است که پیام بعد از تأیید از صف حذف می‌شود.</li>
</ul>

<h2 dir="rtl" align="right">4) چطور bottleneck I/O را حل کنیم؟</h2>

<p dir="rtl" align="right">برخی سیستم‌ها CPU-bound نیستند بلکه I/O-bound هستند مثل:</p>

<ul dir="rtl">
  <li>file storage</li>
  <li>chat systems</li>
  <li>streaming</li>
  <li>database-heavy systems</li>
</ul>

<p dir="rtl" align="right">به عنوان مثال فید اینستاگرام، یکی از این موارد است که مشکل اصلی در آن DB reads می باشد نه CPU (یک workload از نوع read-heavy&rlm;). برای موارد مشابه این از تکنیک‌هایی مثل caching، read replica، index، batching و pagination استفاده می‌کنیم.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<ul dir="rtl">
  <li>چطور بفهمیم سیستم I/O-bound است؟ وقتی CPU پایین است ولی latency بالاست و بیشتر زمان request صرف انتظار برای دیتابیس یا شبکه می‌شود.</li>
  <li>«Read-heavy» نام یک <strong>نوع workload</strong> است، نه یک راه‌حل. راه‌حل‌ها این بخش‌ها هستند:
    <ul dir="rtl">
      <li><a href="../02-03-caching/caching.md">Caching</a></li>
      <li><a href="../02-04-batching-and-pagination/batching.md">Batching</a></li>
      <li><a href="../02-04-batching-and-pagination/pagination.md">Pagination</a></li>
    </ul>
  </li>
  <li>Chat و streaming بیشتر مسئله‌ی <strong>تعداد connection</strong> هستند تا I/O دیسک یا دیتابیس. این موضوع در مورد 7 آمده است.</li>
  <li>در Go، goroutine برای I/O همزمان کم‌هزینه است، ولی یک دیتابیس کند را سریع نمی‌کند. محدودیت واقعی معمولاً <a href="../02-02-connection-pooling/connection-pooling.md">Connection Pool</a> و ظرفیت خود دیتابیس است.</li>
</ul>

<h2 dir="rtl" align="right">5) چطور سیستم را monitor کنیم؟</h2>

<p dir="rtl" align="right">مساله اصلی در این مورد Observability می‌باشد، این یعنی فهمیدن  اینکه داخل سیستم چه اتفاقی می‌افتد. عموما برای این مورد، سه ستون اصلی داریم که شامل log و metrics و tracing می‌شوند.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">هر ستون به یک سؤال جواب می‌دهد:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">ستون</th>
      <th align="right">سؤالی که جواب می‌دهد</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">Metrics</td>
      <td align="right">چقدر؟ آیا مشکلی هست؟ (مثلاً P99 و error rate)</td>
    </tr>
    <tr>
      <td align="right">Logs</td>
      <td align="right">دقیقاً چه اتفاقی افتاد؟</td>
    </tr>
    <tr>
      <td align="right">Traces</td>
      <td align="right">در کدام قسمت مسیر request کند شد؟</td>
    </tr>
  </tbody>
</table>

<ul dir="rtl">
  <li><strong>Monitoring</strong> سؤال‌هایی را جواب می‌دهد که از قبل می‌دانیم باید بپرسیم (dashboard و alert&rlm;). <strong>Observability</strong> کمک می‌کند مشکلی را بررسی کنیم که از قبل پیش‌بینی نکرده بودیم.</li>
  <li>یک نقطه‌ی شروع خوب، چهار <strong>Golden Signals</strong> است: Latency، Traffic، Errors و Saturation&rlm;.</li>
  <li>Alerting بخش جدایی‌ناپذیر این بحث است. بدون alert، metric فقط وقتی کمک می‌کند که کسی به dashboard نگاه کند. ارتباط این موضوع با <a href="../02-01-performance/sli-slo-sla.md">SLI / SLO / SLA</a> و <a href="../02-01-performance/measurement-tools.md">ابزارهای اندازه‌گیری</a> را ببینید.</li>
</ul>

<h2 dir="rtl" align="right">6) چطور سرویس‌ها را decouple کنیم؟</h2>

<p dir="rtl" align="right">این موضوع به معماری Microservice نزدیک است، ولی Microservice فقط یکی از راه‌های رسیدن به decoupling است و به‌تنهایی آن را تضمین نمی‌کند (می‌شود Microservice داشت و باز هم coupling شدید داشت). در حالت monolith عموما با مشکلاتی چون deploy و scale سخت و همینطور coupling زیاد طرف هستیم. به عنوان مثال در سیستمی مانند این service A -> service B -> service C مشکل coupling داریم که اگر مثلا سرویس C به مشکل بخورد، همه سیستم down خواهد شد. برای رفع این مشکلات از راه حل های decoupling استفاده می کنیم.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">این بخش در واقع <strong>دو مشکل جدا</strong> را با هم مطرح می‌کند:</p>

<ol dir="rtl">
  <li><strong>Deploy و scale مستقل:</strong> در monolith برای تغییر کوچک باید کل سیستم deploy شود و نمی‌شود فقط بخش پرفشار را scale کرد. راه‌حل معمول جدا کردن سرویس‌هاست.</li>
  <li><strong>زنجیره‌ی call همزمان:</strong> در <code>A -> B -> C</code>، اگر C کند یا خراب شود، B و A هم منتظر می‌مانند و تمام سیستم خراب می‌شود (cascading failure&rlm;). این مشکل ربطی به monolith یا Microservice بودن ندارد و فقط به همزمان بودن call‌ها مربوط است.</li>
</ol>

<p dir="rtl" align="right">راه‌حل‌های مشکل دوم:</p>

<ul dir="rtl">
  <li>ارتباط async با event یا صف، تا A منتظر C نماند</li>
  <li>timeout برای هر call</li>
  <li>circuit breaker (وقتی C خراب است، به‌جای انتظار، فوراً خطا بدهیم)</li>
  <li>bulkhead (منابع هر dependency را جدا کنیم تا یکی کل سیستم را نخواباند)</li>
  <li>retry با backoff</li>
</ul>

<p dir="rtl" align="right">Microservice هزینه هم دارد: latency شبکه، تراکنش‌های توزیع‌شده و پیچیدگی operation. برای همین «Microservice همیشه بهتر است» درست نیست.</p>

<h2 dir="rtl" align="right">7) چطور real-time scale کنیم؟</h2>

<p dir="rtl" align="right">سیستم های بزرگی همچون WhatsApp, Slack, Trading system, Uber دارای تعداد بسیار زیادی (از ده‌ها هزار تا میلیون‌ها) persistent connection هستند. لذا اینکه چطور آن ها را به صورت  real time مورد scale قرار دهیم، بسیار اهمیت خواهند داشت.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">چرا این مسئله سخت است؟ چون connection ماندگار یعنی <strong>state روی سرور</strong>:</p>

<ul dir="rtl">
  <li>باید پیام به سروری برسد که کاربر به آن وصل است</li>
  <li>ارسال یک پیام به هزاران نفر همزمان (<strong>fan-out</strong>)</li>
  <li>بعد از یک قطعی، همه‌ی کاربران همزمان reconnect می‌کنند و فشار ناگهانی ایجاد می‌شود</li>
  <li>هر connection مقداری حافظه مصرف می‌کند</li>
  <li>یک client کند می‌تواند صف ارسال را روی سرور پر کند</li>
</ul>

<p dir="rtl" align="right">تکنیک‌های رایج:</p>

<ul dir="rtl">
  <li>WebSocket، SSE یا long polling برای نگه داشتن connection</li>
  <li>Pub/Sub (مثل Redis Pub/Sub یا NATS) برای رساندن پیام به سرور درست</li>
  <li>جدا کردن لایه‌ی <strong>gateway</strong> (که فقط connection را نگه می‌دارد) از منطق کسب‌وکار</li>
</ul>

<p dir="rtl" align="right">در Go، goroutine کم‌هزینه است و برای نگه داشتن تعداد زیادی connection همزمان مناسب است. این یکی از دلایل محبوبیت Go برای این نوع سیستم‌هاست.</p>

<h2 dir="rtl" align="right">8) چقدر consistency داده‌ها مهم است؟</h2>

<p dir="rtl" align="right">در سیستم‌های توزیع‌شده همیشه tradeoff داریم که با CAP Theorem توضیح داده می‌شود. طبق CAP، وقتی partition شبکه رخ می‌دهد، باید بین Consistency و Availability یکی را انتخاب کنیم. چون partition در سیستم توزیع‌شده قابل‌اجتناب نیست، انتخاب عملی بین CP و AP است؛ و وقتی partition وجود ندارد، می‌توان هر دو را داشت. در این مورد هم از مواردی چون Strong Consistency (مانند تراکنش  بانکی) یا Eventual Consistency (لایک های اینستاگرام) استفاده می‌کنیم.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">مثال: دو data center نمی‌توانند با هم ارتباط بگیرند (partition) و کاربر یک تراکنش بانکی ثبت می‌کند.</p>

<ul dir="rtl">
  <li><strong>سیستم CP:</strong> درخواست را رد می‌کند و خطا برمی‌گرداند، تا داده در دو طرف ناسازگار نشود. Availability را فدا کرده است.</li>
  <li><strong>سیستم AP:</strong> درخواست را می‌پذیرد و بعداً دو طرف را با هم هماهنگ می‌کند. ممکن است برای مدتی داده‌ها متفاوت باشند. Consistency را فدا کرده است.</li>
</ul>

<p dir="rtl" align="right">برای تراکنش بانکی CP مناسب‌تر است و برای لایک اینستاگرام AP کافی است.</p>

<ul dir="rtl">
  <li>حرف C در CAP با حرف C در ACID <strong>یکی نیست</strong>. در CAP منظور این است که همه‌ی کلاینت‌ها همیشه آخرین مقدار را ببینند. در ACID منظور این است که تراکنش، قوانین و محدودیت‌های داده را نقض نکند.</li>
  <li>حتی وقتی partition نداریم، بین <strong>latency و consistency</strong> trade-off داریم (این چارچوب PACELC نام دارد). هرچه برای هماهنگی بین نسخه‌ها بیشتر صبر کنیم، پاسخ دیرتر می‌رسد. همین‌جا است که Cache با TTL، stale بودن داده را می‌پذیرد تا سریع‌تر باشد (<a href="../02-03-caching/ttl-strategy.md">TTL Strategy</a>&rlm;).</li>
  <li>در یک سیستم واحد معمولاً <strong>هر feature انتخاب جدا دارد</strong>: پرداخت Strong، فید و لایک Eventual&rlm;.</li>
</ul>

---

<h1 dir="rtl" align="right">چالش‌های تکمیلی</h1>

<p dir="rtl" align="right">این سه چالش در لیست اصلی نبودند، ولی در عمل بسیار مهم‌اند.</p>

<h2 dir="rtl" align="right">9) چطور سریع باشیم؟ (Latency)</h2>

<p dir="rtl" align="right">بخش بزرگی از همین مجموعه درباره‌ی این سؤال است، ولی در لیست اصلی جدا نیامده بود. کاربر فقط به «کار می‌کند یا نه» توجه نمی‌کند؛ سرعت پاسخ هم برایش مهم است. مهم این است که به‌جای میانگین، به Percentile (مثل P95 و P99) نگاه کنیم، چون کندترین درخواست‌ها بیشترین ناراحتی را ایجاد می‌کنند. (<a href="../02-01-performance/latency.md">Latency</a> و <a href="../02-01-performance/percentiles.md">Percentile-based Thinking</a>)</p>

<h2 dir="rtl" align="right">10) چطور در برابر خطا مقاوم باشیم؟ (Fault Tolerance)</h2>

<p dir="rtl" align="right">در مقیاس بزرگ، خطا استثنا نیست و همیشه رخ می‌دهد: یک سرویس کند می‌شود، یک connection قطع می‌شود، یک instance می‌افتد. سیستم باید طوری طراحی شود که این خطاها به کل سیستم سرایت نکنند:</p>

<ul dir="rtl">
  <li><strong>timeout</strong> برای هر call، تا منتظر بی‌نهایت نمانیم</li>
  <li><strong>retry با backoff</strong>، چون retry بی‌ملاحظه خودش می‌تواند فشار را چند برابر کند (retry storm)</li>
  <li><strong>circuit breaker</strong> و <strong>bulkhead</strong> (در مورد 6 توضیح داده شد)</li>
  <li><strong>idempotency</strong>، تا retry باعث اثر تکراری نشود (مثلاً دو بار کسر موجودی)</li>
</ul>

<p dir="rtl" align="right">مثالی از اهمیت timeout در <a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a> آمده است.</p>

<h2 dir="rtl" align="right">11) چطور داده را تقسیم کنیم؟ (Sharding / Partitioning)</h2>

<p dir="rtl" align="right">یک دیتابیس واحد، چه از نظر حجم داده و چه از نظر تعداد write در ثانیه، سقف دارد. راه‌حل این است که داده بین چند دیتابیس یا partition تقسیم شود. چالش‌ها:</p>

<ul dir="rtl">
  <li>انتخاب <strong>shard key</strong>. اگر کلید بد انتخاب شود، بعضی partitionها بسیار شلوغ‌تر از بقیه می‌شوند (hot partition&rlm;).</li>
  <li>query و تراکنشی که چند shard را درگیر کند، سخت و کند می‌شود.</li>
  <li>تغییر تعداد shard بعد از شروع کار (resharding) بسیار پرهزینه است.</li>
</ul>

---

<h2 dir="rtl" align="right">نقشه‌ی چالش‌ها در این مجموعه</h2>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">چالش</th>
      <th align="right">جایی که در این مجموعه پوشش داده شده</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">1) Scale</td>
      <td align="right">هنوز بخش مستقل ندارد</td>
    </tr>
    <tr>
      <td align="right">2) Availability</td>
      <td align="right"><a href="../02-01-performance/sli-slo-sla.md">SLI / SLO / SLA</a></td>
    </tr>
    <tr>
      <td align="right">3) Request از دست نرفتن</td>
      <td align="right">هنوز بخش مستقل ندارد</td>
    </tr>
    <tr>
      <td align="right">4) I/O Bottleneck</td>
      <td align="right"><a href="../02-03-caching/caching.md">Caching</a>، <a href="../02-04-batching-and-pagination/batching.md">Batching</a>، <a href="../02-04-batching-and-pagination/pagination.md">Pagination</a>، <a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a></td>
    </tr>
    <tr>
      <td align="right">5) Monitoring</td>
      <td align="right"><a href="../02-01-performance/measurement-tools.md">ابزارهای اندازه‌گیری</a>، <a href="../02-01-performance/sli-slo-sla.md">SLI / SLO / SLA</a></td>
    </tr>
    <tr>
      <td align="right">6) Decoupling</td>
      <td align="right">هنوز بخش مستقل ندارد</td>
    </tr>
    <tr>
      <td align="right">7) Real-time</td>
      <td align="right">هنوز بخش مستقل ندارد</td>
    </tr>
    <tr>
      <td align="right">8) Consistency</td>
      <td align="right"><a href="../02-03-caching/ttl-strategy.md">TTL Strategy</a> (از زاویه‌ی Cache)</td>
    </tr>
    <tr>
      <td align="right">9) Latency</td>
      <td align="right"><a href="../02-01-performance/latency.md">Latency</a>، <a href="../02-01-performance/percentiles.md">Percentiles</a></td>
    </tr>
    <tr>
      <td align="right">10) Fault Tolerance</td>
      <td align="right"><a href="../02-02-connection-pooling/connection-pooling.md">Connection Pooling</a> (بخش timeout)</td>
    </tr>
    <tr>
      <td align="right">11) Sharding</td>
      <td align="right">هنوز بخش مستقل ندارد</td>
    </tr>
  </tbody>
</table>

</div>
