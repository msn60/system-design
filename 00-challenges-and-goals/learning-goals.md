<div dir="rtl" align="right">

<h1 dir="rtl" align="right">هدف نهایی یادگیری System Design</h1>

<p dir="rtl" align="right">برای تبدیل شدن به یک مهندس قوی در System Design، چه ذهنیت و چه دانش‌هایی لازم است؟</p>

<h2 dir="rtl" align="right">Design Knowledge & Mindset</h2>

<p dir="rtl" align="right">فقط دانستن ابزارها کافی نیست. مهم‌تر طرز فکر طراحی سیستم است. مهندس جونیور معمولاً می‌پرسد: از چه دیتابیسی استفاده کنم؟ اما مهندس سینیور می‌پرسد:</p>

<ul dir="rtl">
  <li>مشکل سیستم چیست؟</li>
  <li>بیشترین فشار کجاست؟</li>
  <li>آیا Latency مهم‌تر است یا consistency؟</li>
</ul>

<p dir="rtl" align="right">یعنی اول <strong>problem modeling</strong> بعد tool selection</p>

<p dir="rtl" align="right"><strong>یک مثال واقعی</strong>: فرض کن می‌خواهیم سیستم notification مثل Instagram طراحی کنیم. یک ذهنیت بد: بیاییم Kafka استفاده کنیم چون همه استفاده می‌کنند. یک ذهنیت درست: ما peak داریم؟ ordering مهم است؟ delivery guarantee لازم است؟ fan-out چقدر است؟ بعد ابزار انتخاب می‌شود.</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">هر کدام از آن چهار سؤال، یک تصمیم طراحی را مشخص می‌کند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">سؤال</th>
      <th align="right">چه چیزی را تعیین می‌کند</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">ما peak داریم؟</td>
      <td align="right">آیا برای جذب ترافیک ناگهانی به صف نیاز داریم؟ (<a href="challenges.md">چالش 3</a>)</td>
    </tr>
    <tr>
      <td align="right">ordering مهم است؟</td>
      <td align="right">آیا ترتیب پیام‌های یک کاربر باید حفظ شود؟ اگر بله، معمولاً بر اساس یک کلید (مثل user id) تقسیم‌بندی می‌کنیم.</td>
    </tr>
    <tr>
      <td align="right">delivery guarantee لازم است؟</td>
      <td align="right">مثلاً «حداقل یک بار» (at-least-once) تحویل داده شود. در این حالت پیام ممکن است تکراری برسد و پردازش باید idempotent باشد.</td>
    </tr>
    <tr>
      <td align="right">fan-out چقدر است؟</td>
      <td align="right">یک رویداد برای چند نفر notification می‌سازد؟ پستی از حسابی با ده میلیون فالوور، با یک پست از حساب معمولی، طراحی کاملاً متفاوتی می‌خواهد.</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">نکته: مسئله این نیست که Kafka انتخاب بدی است. ممکن است بعد از جواب دادن به این سؤال‌ها هم Kafka انتخاب درست باشد. مشکل <strong>انتخاب بدون دلیل</strong> است.</p>

<p dir="rtl" align="right"><strong>قدم اول: عدد.</strong> قبل از انتخاب ابزار، مقیاس را به عدد تبدیل کنید. مثلاً اگر روزانه 100 میلیون notification داریم:</p>

<div dir="ltr" align="left">

```text
100,000,000 ÷ 86,400 ثانیه ≈ 1,157 notification در ثانیه (میانگین)
```

</div>

<p dir="rtl" align="right">این عدد میانگین بدی نیست، ولی ترافیک واقعی یکنواخت نیست. اگر فرض کنیم peak پنج برابر میانگین است (فرض تقریبی)، به حدود 5,800 در ثانیه می‌رسیم. طراحی باید برای peak باشد، نه میانگین.</p>

<p dir="rtl" align="right">یک فرایند ساده برای شروع هر طراحی:</p>

<ol dir="rtl">
  <li><strong>Requirements:</strong> سیستم دقیقاً چه کاری انجام می‌دهد؟ چه چیزهایی را انجام نمی‌دهد؟ (<a href="../01-non-functional-requirements/non-functional-requirements.md">Non-Functional Requirements</a>)</li>
  <li><strong>Estimation:</strong> تعداد کاربر، RPS، حجم داده، نسبت read به write</li>
  <li><strong>High-level design:</strong> اجزای اصلی و ارتباط بین آن‌ها</li>
  <li><strong>Deep dive:</strong> سخت‌ترین بخش (bottleneck) را عمیق‌تر بررسی کنیم</li>
  <li><strong>Trade-offs و خطاها:</strong> هر انتخاب چه هزینه‌ای دارد و سیستم چطور خراب می‌شود؟</li>
</ol>

<h2 dir="rtl" align="right">Decision Making & Technology Selection</h2>

<p dir="rtl" align="right">در دنیای واقعی چندین راه حل همیشه وجود دارد. کار مهندس ارشد این است: Trade-offs را بفهمد. بهترین تصمیم را برای context بگیرد.</p>

<p dir="rtl" align="right">مثال: ساخت یک chat system. انتخاب‌ها:</p>

<ul dir="rtl">
  <li>انتخاب REST API — مزیت: ساده؛ معایب: به‌صورت پیش‌فرض realtime نیست (client باید مدام سرور را poll کند)</li>
  <li>انتخاب WebSocket — مزیت: realtime؛ معایب: مدیریت connection سخت</li>
  <li>انتخاب MQTT — مزیت: lightweight؛ معایب: infra خاص می‌خواهد (یک broker)</li>
</ul>

<p dir="rtl" align="right">هیچ گزینه‌ای بهترین مطلق نیست. همیشه: context decides</p>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<p dir="rtl" align="right">دو گزینه‌ی دیگر هم معمولاً در این مقایسه می‌آیند:</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">گزینه</th>
      <th align="right">مزیت</th>
      <th align="right">عیب</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">Long polling</td>
      <td align="right">روی HTTP معمولی کار می‌کند و ساده‌تر از WebSocket است</td>
      <td align="right">هر پیام یک request جدید می‌خواهد و فشار روی سرور بیشتر است</td>
    </tr>
    <tr>
      <td align="right">SSE (Server-Sent Events)</td>
      <td align="right">ساده است و reconnect خودکار دارد</td>
      <td align="right">فقط از سرور به client است. برای ارسال پیام، باید از یک request معمولی استفاده کرد</td>
    </tr>
  </tbody>
</table>

<p dir="rtl" align="right">«Context» یعنی چه؟ این معیارها جواب را عوض می‌کنند:</p>

<ul dir="rtl">
  <li><strong>تعداد کاربر:</strong> یک chat داخلی برای 1,000 نفر با یک chat موبایل برای 10 میلیون نفر یک جواب ندارد. برای اولی polling ساده هم کافی است، برای دومی به WebSocket (یا MQTT) و یک لایه‌ی gateway نیاز داریم (<a href="challenges.md">چالش 7</a>).</li>
  <li><strong>نوع client:</strong> موبایل با شبکه ضعیف و مرورگر دسکتاپ نیازهای متفاوتی دارند.</li>
  <li><strong>مهارت تیم و هزینه‌ی نگه‌داری:</strong> گزینه‌ای که تیم نتواند آن را operate کند، گزینه‌ی خوبی نیست، حتی اگر روی کاغذ بهتر باشد.</li>
  <li><strong>برگشت‌پذیری:</strong> تغییر یک تصمیم چقدر سخت است؟ تغییر کتابخانه‌ی HTTP ساده است، ولی تغییر دیتابیس اصلی بسیار پرهزینه است. تصمیم‌های برگشت‌ناپذیر باید دقیق‌تر بررسی شوند.</li>
</ul>

<p dir="rtl" align="right">یک عادت کاربردی: هر تصمیم مهم را در یک یادداشت کوتاه بنویسید (به این روش <strong>ADR</strong>، یعنی Architecture Decision Record، می‌گویند). شامل: مسئله، گزینه‌ها، انتخاب، و دلیل. شش ماه بعد کمک می‌کند بفهمید چرا آن تصمیم گرفته شده بود.</p>

<h2 dir="rtl" align="right">Upgrade your career and responsibilities</h2>

<p dir="rtl" align="right">سیستم دیزاین معمولاً نقطه تفاوت این‌هاست:</p>

<ul dir="rtl">
  <li>Mid-level engineer</li>
  <li>Senior engineer</li>
  <li>Staff engineer</li>
</ul>

<p dir="rtl" align="right">فرق اصلی:</p>

<ul dir="rtl">
  <li>میدلول: feature implementation</li>
  <li>سینیور: system ownership، architecture decisions، scalability planning</li>
</ul>

<p dir="rtl" align="right">مثال: میدلول: feature: send notification&rlm;. سینیور فکر می‌کند:</p>

<ul dir="rtl">
  <li>اگر 100M notification/day شد چه؟</li>
  <li>اگر queue پر شد چه؟</li>
  <li>اگر worker crash کرد چه؟</li>
</ul>

<h3 dir="rtl" align="right">نکات تکمیلی</h3>

<ul dir="rtl">
  <li><strong>تعریف سطح‌ها بین شرکت‌ها فرق می‌کند.</strong> این بخش یک الگوی کلی است، نه یک قانون.</li>
  <li><strong>Staff engineer</strong> در متن اصلی توضیح داده نشده بود. معمولاً فرق او با سینیور این است که تصمیم‌هایش فقط یک سرویس را درگیر نمی‌کند و روی چند تیم یا کل سیستم اثر می‌گذارد. مسئله‌ی او اغلب «کدام مسئله ارزش حل کردن دارد» و هماهنگ کردن چند تیم است، نه فقط نوشتن راه‌حل.</li>
  <li><strong>System ownership</strong> یعنی فقط نوشتن کد نیست: سیستم را مانیتور می‌کنی، وقتی خراب شد جوابگو هستی، هزینه‌اش را می‌دانی، و برای آینده‌اش برنامه داری.</li>
  <li>سه سؤال سینیور هر کدام به یکی از <a href="challenges.md">چالش‌های اصلی</a> برمی‌گردد:</li>
</ul>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">سؤال</th>
      <th align="right">چالش مرتبط</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">اگر 100M notification/day شد چه؟</td>
      <td align="right">چالش 1 (scale). این عدد حدود 1,157 در ثانیه در میانگین است، ولی باید برای peak بررسی شود.</td>
    </tr>
    <tr>
      <td align="right">اگر queue پر شد چه؟</td>
      <td align="right">چالش 3 (request از دست ندهیم): backpressure و load shedding</td>
    </tr>
    <tr>
      <td align="right">اگر worker crash کرد چه؟</td>
      <td align="right">چالش 10 (fault tolerance): پیام باید دوباره پردازش شود، و پردازش باید idempotent باشد</td>
    </tr>
  </tbody>
</table>

---

<h2 dir="rtl" align="right">بعد از این مجموعه باید بتوانی...</h2>

<p dir="rtl" align="right">یک چک‌لیست ساده برای سنجیدن پیشرفت:</p>

<ul dir="rtl">
  <li>مقیاس یک سیستم را به عدد تبدیل کنی (RPS، حجم داده، نسبت read به write)</li>
  <li>تشخیص بدهی bottleneck از CPU است یا I/O</li>
  <li>برای هر انتخاب فنی، حداقل دو گزینه و trade-off هر کدام را بگویی</li>
  <li>توضیح بدهی سیستم چطور خراب می‌شود و چطور متوجه خرابی‌اش می‌شوی</li>
  <li>طراحی را به زبان ساده برای یک هم‌تیمی یا مصاحبه‌کننده توضیح بدهی</li>
</ul>

<h2 dir="rtl" align="right">مسیر یادگیری در این مجموعه</h2>

<ol dir="rtl">
  <li><a href="challenges.md">چالش‌های اصلی System Design</a>: نقشه‌ی کلی مسائل</li>
  <li><a href="../01-non-functional-requirements/README.md">Non-Functional Requirements</a>: سیستم چه ویژگی‌هایی باید داشته باشد</li>
  <li><a href="../02-01-performance/README.md">Performance</a>: چطور سرعت را اندازه بگیریم و بفهمیم</li>
  <li><a href="../02-02-connection-pooling/README.md">Connection Pooling</a>: اولین الگوی بهبود Performance</li>
  <li><a href="../02-03-caching/README.md">Caching</a>: کاهش فشار روی دیتابیس</li>
  <li><a href="../02-04-batching-and-pagination/README.md">Batching و Pagination</a>: کم کردن تعداد و حجم عملیات</li>
</ol>

</div>
