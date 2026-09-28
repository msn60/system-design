# Release Readiness Checklist

قبل از release باید مشخص باشد: 

## Reliability

- SLO تعریف شده؟
- alert ها meaningful هستند؟
- rollback plan وجود دارد؟
- canary strategy آماده است؟

## Database

- query plan بررسی شده؟
- migration rollback-safe است؟
- index ها validate شده‌اند؟

## Performance

- p95/p99 acceptable است؟
- peak traffic تست شده؟
- capacity limit مشخص است؟

## Dependency Safety

- timeout تنظیم شده‌اند؟
- retry policy امن است؟
- circuit breaker وجود دارد؟
- fallback behavior تست شده؟

خلاصه مبحث Performance

Performance engineering در backend فقط درباره “کم کردن latency” نیست. بلکه درباره: 

- predictability
- stability
- scalability
- graceful degradation
- observability
- saturation control
- resilience under load

است. سیستم production-grade باید:

- تحت workload واقعی قابل پیش‌بینی بماند،
- tail latency کنترل‌شده داشته باشد،
- failure را isolate کند،
- و هنگام overload به‌صورت کنترل‌شده degrade شود، نه اینکه collapse کند.

از بخش به بعد در مورد Pattern ها و Anti pattern های performance صحبت می‌کنیم: 
