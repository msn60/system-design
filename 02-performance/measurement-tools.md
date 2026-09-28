# ابزارهای اندازه‌گیری متریک های Performance

برای اندازه‌گیری Latency، TTFB و Response Time، ابزارهای مختلفی در لایه‌های متفاوت سیستم استفاده می‌شوند.

## Open Telemetry

OpenTelemetry برای جمع‌آوری metrics، logs و traces به‌صورت استاندارد استفاده می‌شود. با Open Telemetry می‌توان مشاهده کرد که: 

- Request وارد کدام سرویس شده است.
- در هر سرویس چقدر زمان مصرف شده است.
- کدام DB query کند بوده است.
- کدام external call باعث delay شده است.
- کدام span بیشترین سهم را در latency داشته است.

## Prometheus Histogram Metrics

Prometheus Histogram برای اندازه‌گیری duration ها و محاسبه‌ی percentile ها بسیار کاربردی است. نمونه‌ی متریک‌های رایج:

- http_request_duration_seconds_bucket
- http_request_duration_seconds_sum
- http_request_duration_seconds_count

با استفاده از این متریک‌ها می‌توان P95، P99 و سایر percentile ها را محاسبه کرد. 

## Jaeger Tracing

Jaeger برای distributed tracing استفاده می‌شود و کمک می‌کند مسیر یک request در چند سرویس مختلف بررسی شود. مثال برای endpoint زیر: GET /feed ممکن است trace به شکل زیر باشد:

- API Gateway:        20ms
- Feed Service:      400ms
- User Service:       50ms
- Post Service:      120ms
- Ranking Service:   250ms
- Redis:              10ms
- PostgreSQL:         90ms

با این اطلاعات می‌توان مشخص کرد بیشترین زمان request در کدام سرویس یا dependency مصرف شده است.

##  API Gateway Metrics

در API Gateway ها و reverse proxy هایی مانند Nginx، Envoy و Traefik، متریک‌های زیر بسیار مهم هستند:

- request_duration
- upstream_connect_time
- upstream_first_byte_time
- upstream_response_time
- status_code
- retry_count
- active_connections
- pending_requests

این متریک‌ها کمک می‌کنند مشخص شود مشکل از gateway، شبکه، upstream service یا backend اصلی است. برای مثال در Nginx:

- upstream_connect_time
- upstream_header_time
- upstream_response_time
- request_time

تقریباً می‌توان آن‌ها را این‌گونه تفسیر کرد:

- upstream_connect_time   = زمان اتصال به upstream
- upstream_header_time    = زمان تا دریافت اولین header/byte از upstream
- upstream_response_time  = زمان کامل شدن پاسخ upstream
- request_time            = کل زمان request از دید Nginx

## Browser / Frontend Metrics

برای بررسی تجربه‌ی واقعی کاربر در browser، از API هایی مانند Navigation Timing API و Resource Timing API استفاده می‌شود. این ابزارها می‌توانند موارد زیر را اندازه‌گیری کنند:

- DNS Time
- TCP Connect Time
- TLS Handshake Time
- TTFB
- Content Download Time
- DOM Load Time
- Largest Contentful Paint

در طراحی سیستم‌های backend، همیشه لازم نیست وارد تمام جزئیات frontend metrics شویم، اما اگر سرویس web-facing باشد، فهم این متریک‌ها برای تحلیل تجربه‌ی کاربر بسیار مفید است.

| **دسته** | **پارامتر** | **چه چیزی را اندازه‌گیری می‌کند؟** | **کاربرد اصلی** |
|---|---|---|---|
| **Network Latency** | DNS Time | زمان resolve شدن دامنه | تحلیل client/network |
| **Network Latency** | TCP Connect Time | زمان برقراری connection | تحلیل connection overhead |
| **Network Latency** | TLS Handshake Time | زمان امن‌سازی ارتباط | تحلیل HTTPS overhead |
| **Network Latency** | RTT | زمان رفت‌وبرگشت شبکه | تحلیل کیفیت شبکه |
| **Application Latency** | Queue Time | انتظار قبل از پردازش | تشخیص saturation |
| **Application Latency** | Handler / Dispatch Time | زمان رسیدن request به handler | تحلیل framework/runtime |
| **Application Latency** | Business Logic Time | زمان اجرای منطق اصلی برنامه | تحلیل application bottleneck |
| **Internal Latency** | DB Latency | زمان اجرای query | تشخیص bottleneck دیتابیس |
| **Internal Latency** | Cache Latency | زمان پاسخ Redis/Memcached | تحلیل cache layer |
| **Internal Latency** | RPC Latency | زمان call به سرویس دیگر | تحلیل microservices |
| **Internal Latency** | Message Broker Latency | زمان publish/consume پیام | تحلیل async pipeline |
| **TTFB** | Time to First Byte | زمان تا دریافت اولین byte | وب، streaming، response بزرگ |
| **TTFB** | Upstream First Byte | زمان اولین پاسخ از upstream | تحلیل gateway/proxy |
| **TTFB** | Serialization Before First Write | زمان آماده‌سازی اولین خروجی | تحلیل response generation |
| **Response Time** | Total Response Time | زمان کامل شدن response | SLO اصلی API |
| **Response Time** | Download Time | زمان دریافت کامل body | فایل، payload بزرگ، شبکه کند |
| **Response Time** | TTLB | زمان دریافت آخرین byte | معادل عملی full response time |
| **Advanced** | P95 / P99 / P99.9 | تجربه‌ی کاربران کندتر | SLO و production monitoring |
| **Advanced** | Jitter | نوسان latency | real-time systems |
| **Advanced** | Queue Depth | عمق صف‌ها | تشخیص فشار و saturation |
| **Advanced** | Throughput + Latency | رفتار سیستم تحت بار | capacity planning |
