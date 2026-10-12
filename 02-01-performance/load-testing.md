<div dir="rtl" align="right">

<h1 dir="rtl" align="right">Load Testing از نگاه مهندسی</h1>

<p dir="rtl" align="right">هدف load test فقط “RPS بالا” نیست. هدف: شناخت behavior سیستم تحت workload واقعی است.</p>

<p dir="rtl" align="right">اشتباه رایج عموما این است که فقط steady traffic تست می‌شود. در حالی که production واقعی: burst و spike و noisy traffic و uneven distribution دارد.</p>

<h2 dir="rtl" align="right">سناریوهای ضروری</h2>

</div>

- steady state
- burst
- spike
- soak
- dependency degradation
- hot partition
- slow client
- partial outage
- retry storm simulation

<div dir="rtl" align="right">

<p dir="rtl" align="right">نکته بسیار مهم: Load Test بدون observability تقریباً بی‌ارزش است.</p>

<p dir="rtl" align="right">باید همزمان:  metrics ،traces ،logs ،DB stats ،GC behavior  ،queuemetrics جمع‌آوری شوند.</p>

</div>
