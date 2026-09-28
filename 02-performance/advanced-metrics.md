# معیارهای پیشرفته Performance برای سیستم‌های بزرگ

در سیستم‌های production و scale بالا، فقط دانستن یک عدد برای latency یا response time کافی نیست. باید رفتار سیستم در درصدهای مختلف، زیر بار بالا و در شرایط ناپایدار بررسی شود. 

## Tail Latency Percentiles

همانگونه که پیشتر هم گفته شد، در سیستم‌های بزرگ، استفاده از میانگین یا average برای تحلیل performance کافی و حتی گمراه‌کننده است. به‌جای آن باید از percentile ها استفاده شود. 

در سیستم‌های بزرگ، حتی ۱٪ از درخواست‌ها می‌تواند میلیون‌ها request باشد. به همین دلیل، SLOها معمولاً روی P95 یا P99 تعریف می‌شوند، نه روی average.

نمونه‌ی SLO حرفه‌ای:

- P95 Response Time for read APIs should be under 300ms.
- P99 Response Time for payment authorization should be under 1s.

## Jitter

Jitter یعنی نوسان latency در طول زمان. برای مثال، سیستم زیر رفتار پایدارتر و قابل‌پیش‌بینی‌تری دارد: 

- 100ms
- 110ms
- 95ms
- 105ms

اما سیستم زیر با وجود اینکه ممکن است میانگین مشابهی داشته باشد، تجربه‌ی بدتری ایجاد می‌کند:

- 30ms
- 40ms
- 900ms
- 50ms
- 1200ms

Jitter مخصوصاً در سیستم‌های زیر اهمیت زیادی دارد:

- Real-time Messaging
- Voice / Video Call
- Online Gaming
- Trading Systems
- IoT
- Robotics

## Queue Depth Metrics

Queue Depth یکی از مهم‌ترین نشانه‌های فشار روی سیستم است. در بسیاری از موارد، قبل از اینکه سیستم error بدهد یا crash کند، ابتدا queueها شروع به رشد می‌کنند. متریک‌های مهم در این بخش عبارت‌اند از:

- Request Queue Length
- Worker Pool Queue Length
- Thread Pool Saturation
- Goroutine Count
- DB Connection Pool Wait Count
- Kafka Consumer Lag
- RabbitMQ Queue Depth

برای مثال، اگر RPS ثابت باشد اما goroutine count، memory usage و queue length به‌تدریج رشد کند، احتمالاً requestها در سیستم گیر کرده‌اند و concurrency در حال افزایش است.

## Throughput و ارتباط آن با Latency

Throughput به‌تنهایی برای ارزیابی performance کافی نیست. اینکه یک سیستم بتواند تعداد زیادی request دریافت کند، لزوماً به این معنا نیست که عملکرد خوبی دارد. عبارت ناقص: The service can handle 5000 RPS. عبارت دقیق‌تر: 

**The service can handle 5000 RPS with P95 Response Time under 300ms.**

زیرا ممکن است سیستم 5000 RPS دریافت کند، اما response time آن به ۱۰ ثانیه برسد. چنین سیستمی از نظر performance قابل قبول نیست. نمونه‌ی تحلیل ظرفیت:

- At 1000 RPS: P95 = 120ms
- At 3000 RPS: P95 = 250ms
- At 5000 RPS: P95 = 700ms
- At 7000 RPS: Error rate increases and queue depth grows

این نوع تحلیل نشان می‌دهد سیستم تا چه نقطه‌ای پایدار است و از چه نقطه‌ای به بعد وارد saturation می‌شود.
