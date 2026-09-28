# Performance

how quickly and efficiently a system responds to user requests.

**Performance: product owner’s view**

Slow performance can lead to user frustration, abandonment, and lost revenue.

- **How quickly should the system respond to** a user action to feel instantaneous?
- What is the **expected number of users** the system should handle without slowing down?
- Are there specific critical user journeys that must have a high level of performance?

Performance یعنی: سیستم با چه سرعت، پایداری و بهره‌وری به درخواست‌های کاربر پاسخ می‌دهد، در شرایط مختلف و تحت بارهای متفاوت.

## آیا performance همان سرعت است؟

اما در عمل Performance فقط «سرعت» نیست. چهار لایه مهم دارد:

### 1.1 Perceived Performance (عملکرد ادراک‌شده)

چیزی که کاربر حس می‌کند، نه لزوماً آنچه واقعاً اتفاق افتاده

آیا سیستم «فوری» به نظر می‌رسد؟ آیا کاربر حس می‌کند سیستم زنده است؟ آیا در حین پردازش بازخورد می‌گیرد؟

 گاهی سیستم کند است، ولی چون Feedback مناسب دارد، کاربر آن را سریع تلقی می‌کند.

### 1.2 Actual Performance (عملکرد واقعی)

اعداد و واقعیت فنی:

- زمان پاسخ
- توان پردازش
- مصرف منابع
- تحمل بار همزمان

### 1.3 Consistency (پایداری عملکرد)

- آیا همیشه سریع است؟
- یا گاهی سریع و گاهی خیلی کند؟

ناپایداری از کندی خطرناک‌تر است.

### 1.4 Scalability Impact

- وقتی کاربران زیاد می‌شوند چه اتفاقی می‌افتد؟
- آیا تجربه افت می‌کند یا ثابت می‌ماند؟

## نقش Performance در User Journey

همه بخش‌ها نیاز به Performance یکسان ندارند.

### Critical User Journeys

مسیرهایی که اگر کند باشند، کاربر می‌رود:

- Login / Sign Up
- Search
- Checkout / Payment
- Navigation اولیه
- First Content Load\
 این مسیرها باید:
- سریع‌تر از میانگین باشند
- Animation + Skeleton + Feedback داشته باشند
- Error Rate نزدیک صفر داشته باشند

### Non-Critical Flows

مثل:

- Analytics
- History
- Settings
- Export ها

می‌توانند:

- Lazy Load شوند
- Async باشند
- Feedback طولانی‌تری داشته باشند

## Performance در System Design

در سیستم دیزاین، Performance فقط مسئله Backend نیست.

**Design Decisions** مؤثر بر **Performance:**

- حجم تصاویر
- Animation Duration
- Use of Skeletons
- Pagination vs Infinite Scroll
- Client-side vs Server-side rendering
- Preload / Prefetch
- Component Complexity

از دید بک‌اند، Performance یعنی: توانایی سیستم در پاسخ‌دهی سریع، پایدار و قابل پیش‌بینی به درخواست‌ها، تحت بارهای مختلف، با حداقل مصرف منابع و بدون افت کیفیت. Performance یعنی سیستم بتواند: درخواست‌ها را سریع، پایدار، قابل پیش‌بینی و با مصرف منطقی منابع پاسخ دهد؛ حتی وقتی بار سیستم تغییر می‌کند.

چهار ستون اصلی دارد:

- Latency – هر درخواست چقدر طول می‌کشد
- Throughput – چند درخواست در واحد زمان پردازش می‌شود
- Concurrency – چند درخواست همزمان بدون افت کیفیت
- Stability under Load – رفتار سیستم وقتی تحت فشار است

**در نظر داشته باشید که Performance ≠ Speed** 

یکی از خطاهای رایج بک‌اند: «اگه سریع باشه، Performance خوبه»

مثلاً ممکن است یک endpoint در تست local فقط در 20ms جواب بدهد، اما وقتی ۵۰۰ کاربر همزمان می‌آیند:

- latency برود روی ۳ ثانیه
- connection pool دیتابیس پر شود
- go routine ها زیاد شوند
- CPU بالا برود
- timeout زیاد شود
- retry storm ایجاد شود

اینجا endpoint در حالت عادی سریع بوده، ولی performance واقعی ندارد.

در واقع Performance ترکیب این‌هاست:

| عامل | توضیح |
|---|---|
| Latency | زمان پاسخ |
| Variance | نوسان زمان پاسخ |
| Saturation | پر شدن منابع |
| Backpressure | واکنش به فشار |
| Failure Behavior | رفتار در خرابی |

سیستمی که گاهی 50ms و گاهی 5s جواب می‌دهد، بدتر از سیستمی است که همیشه 300ms جواب می‌دهد.
