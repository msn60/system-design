# Load Testing از نگاه مهندسی

هدف load test فقط “RPS بالا” نیست. هدف: شناخت behavior سیستم تحت workload واقعی است.

 اشتباه رایج عموما این است که فقط steady traffic تست می‌شود. در حالی که production واقعی: burst و spike و noisy traffic و uneven distribution دارد.

## سناریوهای ضروری

- steady state
- burst
- spike
- soak
- dependency degradation
- hot partition
- slow client
- partial outage
- retry storm simulation

نکته بسیار مهم: Load Test بدون observability تقریباً بی‌ارزش است.

باید همزمان:  metrics ،traces ،logs ،DB stats ،GC behavior  ،queuemetrics جمع‌آوری شوند.
