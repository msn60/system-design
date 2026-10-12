<div dir="rtl" align="right">

<h1 dir="rtl" align="right">Release Readiness Checklist</h1>

<p dir="rtl" align="right">قبل از release باید مشخص باشد:</p>

<h2 dir="rtl" align="right">Reliability</h2>

<ul dir="rtl">
  <li>SLO تعریف شده؟</li>
  <li>alert ها meaningful هستند؟</li>
  <li>rollback plan وجود دارد؟</li>
  <li>canary strategy آماده است؟</li>
</ul>

<h2 dir="rtl" align="right">Database</h2>

<ul dir="rtl">
  <li>query plan بررسی شده؟</li>
  <li>migration rollback-safe است؟</li>
  <li>index ها validate شده‌اند؟</li>
</ul>

<h2 dir="rtl" align="right">Performance</h2>

<ul dir="rtl">
  <li>p95/p99 acceptable است؟</li>
  <li>peak traffic تست شده؟</li>
  <li>capacity limit مشخص است؟</li>
</ul>

<h2 dir="rtl" align="right">Dependency Safety</h2>

<ul dir="rtl">
  <li>timeout تنظیم شده‌اند؟</li>
  <li>retry policy امن است؟</li>
  <li>circuit breaker وجود دارد؟</li>
  <li>fallback behavior تست شده؟</li>
</ul>

<p dir="rtl" align="right">خلاصه مبحث Performance</p>

<p dir="rtl" align="right">Performance engineering در backend فقط درباره “کم کردن latency” نیست. بلکه درباره:</p>

<ul dir="rtl">
  <li>predictability</li>
  <li>stability</li>
  <li>scalability</li>
  <li>graceful degradation</li>
  <li>observability</li>
  <li>saturation control</li>
  <li>resilience under load</li>
</ul>

<p dir="rtl" align="right">است. سیستم production-grade باید:</p>

<ul dir="rtl">
  <li>تحت workload واقعی قابل پیش‌بینی بماند،</li>
  <li>tail latency کنترل‌شده داشته باشد،</li>
  <li>failure را isolate کند،</li>
  <li>و هنگام overload به‌صورت کنترل‌شده degrade شود، نه اینکه collapse کند.</li>
</ul>

<p dir="rtl" align="right">از بخش به بعد در مورد Pattern ها و Anti pattern های performance صحبت می‌کنیم:</p>

</div>
