# Throughput & Capacity

Throughput یعنی: سیستم در واحد زمان چند request می‌تواند پردازش کند. معمولاً با RPS یا QPS بیان می‌شود: 

- RPS = Requests Per Second
- QPS = Queries Per Second

مثلاً: Service can handle 2,000 RPS یعنی این سرویس می‌تواند در هر ثانیه ۲۰۰۰ درخواست را پردازش کند. اما نکته مهم: افزایش Throughput بدون کنترل Latency خطرناک است. مثلاً ممکن است سیستم بتواند ۵۰۰۰ RPS قبول کند، اما latency از 100ms برود روی 8s. این خوب نیست. در عمل باید بگویی: The service handles 5,000 RPS with P95 latency under 300ms. نه فقط: The service handles 5,000 RPS. چون throughput بدون latency ناقص است. 

## Throughput به چه چیزهایی وابسته است؟

در بک‌اند، throughput معمولاً به این‌ها وابسته است:

- CPU
- Memory
- Disk IO
- Network IO
- Database
- Cache
- Locks
- Connection Pool
- External Services

مثلاً در Go ممکن است کد application خیلی سبک باشد، اما bottleneck اصلی دیتابیس باشد. به عنوان مثال فرض کنید یک  endpoint برای دریافت سفارش های کاربران را دارید.  اگر هر request دو query به دیتابیس بزند و دیتابیس فقط بتواند ۱۰۰۰ query در ثانیه تحمل کند، حداکثر throughput تقریبی این endpoint می‌شود: 

```
1000 DB queries/sec ÷ 2 queries/request = 500 requests/sec
```

حتی اگر Go service بتواند از نظر CPU ده‌ها هزار RPS هندل کند، دیتابیس محدودت می‌کند.

نتیجه اینکه: **افزایش Throughput بدون کنترل Latency مساوی خواهد بود با Disaster**
