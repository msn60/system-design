# اصول کلیدی Backend Performance Engineering

## اصل 1 — Tail Latency مهم‌تر از Average است

در production: کاربران average را تجربه نمی‌کنند، بلکه worst-case نسبی را تجربه می‌کنند. به عنوان مثال در سیستمی که Average = 80ms و P99 = 4s می باشد، این سیستم از دید production مشکل دارد.

## اصل 2 — Latency باید Decompose شود

Latency یک عدد واحد نیست. هر request معمولاً شامل این بخش‌هاست: 

```
Total Latency  = queue_time + app_processing_time + db_time + external_dependency_time + serialization_time + network_transfer_time
```

Decomposition مهم است زیرا optimization بدون شناخت bottleneck معمولاً بی‌اثر است. 

## اصل 3 — Saturation دشمن پنهان سیستم است

سیستم‌ها معمولاً قبل از crash شدن degrade می‌شوند. عموما حالت های Saturation، وقتی resource تقریباً پر شده اتفاق می افتند، مانند: 

- CPU
- DB pool
- thread pool
- queue
- network bandwidth

نتیجه Saturation  می تواند مواردی چون زیر باشد: 

- latency spike
- timeout
- retry storm
- queue explosion

## اصل 4 — Back Pressure ضروری است

اگر سیستم نتواند فشار را کنترل کند: overload propagate می‌شود. موارد زیر مثال هایی از Backpressure هستند: 

- bounded queue
- concurrency limit
- request shedding
- adaptive throttling

## اصل 5 — Retry بدون کنترل خطرناک است

Retry می‌تواند outage را بدتر کند. به عنوان مثال Retry Storm می‌تواند: 

```
dependency slow → clients retry → traffic increases → dependency slower → more retries → collapse
```

قواعد حرفه‌ای Retry عبارتند از: 

- exponential backoff
- jitter
- retry budget
- idempotency
- retry only on transient failures

## اصل 6 — Cache باید Load را کاهش دهد، نه مشکل را پنهان کند

Cache اشتباه: stale inconsistency و cache stampede و invalidation chaos ایجاد می‌کند.

## اصل 7 — Capacity Margin لازم است

سیستم نباید دائماً نزدیک 100% utilization کار کند. زیرا در utilization بالا: queueing delay غیرخطی می‌شود. این یکی از مهم‌ترین مفاهیم performance engineering است. 
