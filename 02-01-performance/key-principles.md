<div dir="rtl" align="right">

<h1 dir="rtl" align="right">اصول کلیدی Backend Performance Engineering</h1>

<h2 dir="rtl" align="right">اصل 1 — Tail Latency مهم‌تر از Average است</h2>

<p dir="rtl" align="right">در production: کاربران average را تجربه نمی‌کنند، بلکه worst-case نسبی را تجربه می‌کنند. به عنوان مثال در سیستمی که Average = 80ms و P99 = 4s می باشد، این سیستم از دید production مشکل دارد.</p>

<h2 dir="rtl" align="right">اصل 2 — Latency باید Decompose شود</h2>

<p dir="rtl" align="right">Latency یک عدد واحد نیست. هر request معمولاً شامل این بخش‌هاست:</p>

<div dir="ltr" align="left">

```
Total Latency  = queue_time + app_processing_time + db_time + external_dependency_time + serialization_time + network_transfer_time
```

</div>

<p dir="rtl" align="right">Decomposition مهم است زیرا optimization بدون شناخت bottleneck معمولاً بی‌اثر است.</p>

<h2 dir="rtl" align="right">اصل 3 — Saturation دشمن پنهان سیستم است</h2>

<p dir="rtl" align="right">سیستم‌ها معمولاً قبل از crash شدن degrade می‌شوند. عموما حالت های Saturation، وقتی resource تقریباً پر شده اتفاق می افتند، مانند:</p>

<ul dir="rtl">
  <li>CPU</li>
  <li>DB pool</li>
  <li>thread pool</li>
  <li>queue</li>
  <li>network bandwidth</li>
</ul>

<p dir="rtl" align="right">نتیجه Saturation  می تواند مواردی چون زیر باشد:</p>

<ul dir="rtl">
  <li>latency spike</li>
  <li>timeout</li>
  <li>retry storm</li>
  <li>queue explosion</li>
</ul>

<h2 dir="rtl" align="right">اصل 4 — Back Pressure ضروری است</h2>

<p dir="rtl" align="right">اگر سیستم نتواند فشار را کنترل کند: overload propagate می‌شود. موارد زیر مثال هایی از Backpressure هستند:</p>

<ul dir="rtl">
  <li>bounded queue</li>
  <li>concurrency limit</li>
  <li>request shedding</li>
  <li>adaptive throttling</li>
</ul>

<h2 dir="rtl" align="right">اصل 5 — Retry بدون کنترل خطرناک است</h2>

<p dir="rtl" align="right">Retry می‌تواند outage را بدتر کند. به عنوان مثال Retry Storm می‌تواند:</p>

<div dir="ltr" align="left">

```
dependency slow → clients retry → traffic increases → dependency slower → more retries → collapse
```

</div>

<p dir="rtl" align="right">قواعد حرفه‌ای Retry عبارتند از:</p>

<ul dir="rtl">
  <li>exponential backoff</li>
  <li>jitter</li>
  <li>retry budget</li>
  <li>idempotency</li>
  <li>retry only on transient failures</li>
</ul>

<h2 dir="rtl" align="right">اصل 6 — Cache باید Load را کاهش دهد، نه مشکل را پنهان کند</h2>

<p dir="rtl" align="right">Cache اشتباه: stale inconsistency و cache stampede و invalidation chaos ایجاد می‌کند.</p>

<h2 dir="rtl" align="right">اصل 7 — Capacity Margin لازم است</h2>

<p dir="rtl" align="right">سیستم نباید دائماً نزدیک 100% utilization کار کند. زیرا در utilization بالا: queueing delay غیرخطی می‌شود. این یکی از مهم‌ترین مفاهیم performance engineering است.</p>

</div>
