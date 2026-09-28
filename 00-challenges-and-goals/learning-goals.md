# هدف نهایی یادگیری System Design

 برای تبدیل شدن به یک مهندس قوی در System Design چه ذهنیت و چه دانش‌هایی لازم است

## Design Knowledge & Mindset

فقط دانستن ابزارها کافی نیست. مهم‌تر طرز فکر طراحی سیستم است. مهندس جونیور معمولاً می‌پرسد: از چه دیتابیسی استفاده کنم؟اما مهندس سینیور می‌پرسد:

- مشکل سیستم چیست؟
- بیشترین فشار کجاست؟
- Latency مهم‌تر است یا consistency؟

یعنی اول problem modeling بعد tool selection

یک مثال واقعی: فرض کن می‌خواهیم سیستم notification مثل Instagram طراحی کنیم. یک ذهنیت بد: بیاییم Kafka استفاده کنیم چون همه استفاده می‌کنند. یک ذهنیت درست:ما peak داریم؟ ordering مهم است؟ delivery guarantee لازم است؟ fan-out چقدر است؟ بعد ابزار انتخاب می‌شود.

## Decision Making & Technology Selection

در دنیای واقعی چندین راه حل همیشه وجود دارد. کار مهندس ارشد این است: Trade-offs را بفهمد. بهترین تصمیم برای context بگیرد.

مثال: ساخت یک chat system  انتخاب‌ها:

- REST API مزیت: ساده	معایب: realtime نیست
- WebSocket مزیت: realtime	معایب: مدیریت connection سخت
- MQTT مزیت: lightweight	معایب: infra خاص می‌خواهد

هیچ گزینه‌ای بهترین مطلق نیست.  همیشه: context decides

## Upgrade your career and responsibilities

سیستم دیزاین معمولاً نقطه تفاوت این‌هاست: 

- Mid-level engineer
- Senior engineer
- Staff engineer

فرق اصلی-  میدلول: feature implementation – سینیور: system ownership -  architecture decisions   - scalability planning

مثال: میدلول: feature: send notification سینیور فکر می‌کند:

- اگر 100M notification/day شد چه؟
- اگر queue پر شد چه؟
- اگر worker crash کرد چه؟
