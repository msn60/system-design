<div dir="rtl" align="right">

<h1 dir="rtl" align="right">Non functional Requirements (NFRs)</h1>

<p dir="rtl" align="right">یک سیستم دو بخش دارد:</p>

<p dir="rtl" align="right">بخش اول - Functional Requirements: سیستم چه کاری انجام می‌دهد؟ (مثلاً: کاربر بتواند پول انتقال دهد، کاربر بتواند عکس آپلود کند)</p>

<p dir="rtl" align="right">بخش دوم - Non-Functional Requirements: سیستم چگونه آن کار را انجام می‌دهد؟ (مثلاً: چقدر سریع؟ چقدر امن؟ چقدر پایدار؟)</p>

<p dir="rtl" align="right">در سیستم دیزاین، چالش اصلی همیشه NFR ها هستند. نوشتن کدی که یک ایمیل بفرستد ساده است، اما سیستمی که ۱ میلیارد ایمیل را بدون گم شدن و با سرعت بالا بفرستد، هنر مهندسی است.</p>

<h2 dir="rtl" align="right">چرا NFR ها در سیستم دیزاین حیاتی هستند؟</h2>

<p dir="rtl" align="right">بدون NFR، سیستم شما در دنیای واقعی شکست می‌خورد.</p>

<ul dir="rtl">
  <li>اگر سیستم شما Scalable نباشد -&gt; با آمدن ۱۰۰۰ کاربر جدید Crash می‌کند.</li>
  <li>اگر Available نباشد -&gt; وقتی کاربر نیاز دارد، سرویس قطع است (مثل قطع شدن درگاه بانکی در شب عید).</li>
  <li>اگر Consistent نباشد -&gt; کاربر پول واریز می‌کند ولی موجودی‌اش آپدیت نمی‌شود.</li>
</ul>

<p dir="rtl" align="right">در مصاحبه‌های سیستم دیزاین، هدف اصلی ارزیابی توانایی فرد در بالانس کردن این NFR ها (Trade-offs) است.</p>

<h2 dir="rtl" align="right">مهم‌ترین NFR ها (The Pillars)</h2>

<p dir="rtl" align="right">تقریباً تمام سیستم‌های توزیع شده حول این ۵ محور می‌چرخند (به نوعی 5 دسته بندی اصلی هستند):</p>

<h3 dir="rtl" align="right">مقیاس‌پذیری - Scalability</h3>

<p dir="rtl" align="right"> آیا سیستم می‌تواند با افزایش بار (ترافیک/داده)، رشد کند؟ NFR هایی که مربوط به رشد سیستم، افزایش کاربر، افزایش داده و مدیریت بار هستند.</p>

<ul dir="rtl">
  <li>Scalability</li>
  <li>Throughput</li>
  <li>Extensibility (قابلیت توسعه‌پذیری برای رشد آینده)</li>
</ul>

<h3 dir="rtl" align="right">در دسترس بودن - Availability</h3>

<p dir="rtl" align="right"> سیستم در چند درصد از مواقع سالم و آماده سرویس‌دهی است؟ (مثلاً ۴ تا ۹ یا ۹۹.۹۹٪)NFR هایی که هدف آن‌ها اطمینان از در دسترس بودن سرویس است.</p>

<ul dir="rtl">
  <li>Availability</li>
  <li>Fault Tolerance</li>
  <li>Resilience</li>
</ul>

<h3 dir="rtl" align="right">قابلیت اطمینان - Reliability</h3>

<p dir="rtl" align="right">آیا سیستم “درست” کار می‌کند و داده‌ای را از دست نمی‌دهد؟ (حتی اگر بخشی از آن خراب شود) .  NFR هایی که تضمین می‌کنند سیستم درست، بدون خطا و قابل اعتماد کار می‌کند</p>

<ul dir="rtl">
  <li>Reliability</li>
  <li>Durability</li>
  <li>Consistency</li>
</ul>

<h3 dir="rtl" align="right">کارایی - Performance</h3>

<p dir="rtl" align="right"> سرعت پاسخگویی (Latency) و ظرفیت پردازش (Throughput) چقدر است؟</p>

<p dir="rtl" align="right">هرچیزی که مربوط به سرعت، تأخیر و کارایی اجرای درخواست‌ هاست.</p>

<ul dir="rtl">
  <li>Performance</li>
  <li>Latency</li>
  <li>Throughput</li>
  <li>Efficiency</li>
</ul>

<h3 dir="rtl" align="right">نگهداری - Maintainability</h3>

<p dir="rtl" align="right">چقدر راحت می‌توان سیستم را تغییر داد، باگ‌ها را پیدا کرد (Observability) و توسعه داد؟ هر NFR مربوط به نگهداری، تغییر، عیب‌یابی و توسعه سیستم</p>

<ul dir="rtl">
  <li>Maintainability</li>
  <li>Observability</li>
  <li>Security</li>
  <li>Interoperability</li>
  <li>Compliance</li>
  <li>Usability</li>
</ul>

<h2 dir="rtl" align="right">اولویت‌بندی: با چه ترتیبی به آن‌ها نگاه کنیم؟</h2>

<p dir="rtl" align="right">این بخش مهم‌ترین قسمت برای یک مهندس سینیور است. اولویت‌بندی ثابت نیست و به Context بستگی دارد، اما یک فریم‌ورک استاندارد برای تحلیل وجود دارد:</p>

<h3 dir="rtl" align="right">گام اول: فهمیدن ماهیت بیزینس (The Trade-off Start)</h3>

<p dir="rtl" align="right">قبل از هر چیز بپرس: “هزینه خرابی یا اشتباه در این سیستم چقدر است؟”</p>

<p dir="rtl" align="right">اگر سیستم بانکی طراحی می‌کنی: Consistency &amp; Reliability اولویت ۱ هستند. (پول نباید گم شود).</p>

<p dir="rtl" align="right">اگر توییتر یا اینستاگرام طراحی می‌کنی: Availability &amp; Scalability اولویت ۱ هستند. (اگر یک لایک دیرتر ثبت شود فاجعه نمی‌شود، اما سیستم نباید بخوابد).</p>

<h3 dir="rtl" align="right">گام دوم: اعمال فریم‌ورک CAP (برای داده‌ها)</h3>

<p dir="rtl" align="right">در سیستم‌های توزیع شده، تو نمی‌توانی همزمان هم ۱۰۰٪ Consistency داشته باشی و هم ۱۰۰٪ Availability. باید انتخاب کنی:</p>

<ul dir="rtl">
  <li>مورد CP (Consistency/Partition): اولویت با صحت داده است (مثل سیستم‌های مالی).</li>
  <li>مورد AP (Availability/Partition): اولویت با بالا بودن سیستم است (مثل شبکه اجتماعی).</li>
</ul>

<h3 dir="rtl" align="right">گام سوم: تعیین سقف Performance</h3>

<p dir="rtl" align="right">حالا که فهمیدی سیستم باید چطور رفتار کند، می‌پرسی: “حداقل سرعت قابل قبول چقدر است؟” (همان بحث p95 و p99 که داشتیم).</p>

<p dir="rtl" align="right">خلاصه اولویت‌بندی در یک نگاه:</p>

<ul dir="rtl">
  <li>کدام مهم تر است؟ Consistency vs Availability: (انتخاب استراتژی اصلی بر اساس بیزینس).</li>
  <li>در نظر گرفتن Scalability: (تخمین تعداد کاربران و حجم داده - سیستم چقدر باید بزرگ شود؟).</li>
  <li>مورد Reliability: (چقدر تحمل خطا یا Fault Tolerance نیاز داریم؟).</li>
  <li>کارایی سیستم و Performance: (بهینه‌سازی برای سرعت و تجربه کاربری).</li>
</ul>

<p dir="rtl" align="right"><strong>مثال واقعی برای جمع‌بندی:</strong></p>

<p dir="rtl" align="right">فرض کن می‌خواهیم سیستم رزرو بلیط هواپیما طراحی کنیم.</p>

<ul dir="rtl">
  <li>اولویت ۱: Consistency. نباید یک صندلی به دو نفر فروخته شود. پس حاضریم سیستم کمی کندتر باشد یا اگر شبکه قطع شد اصلاً رزرو نگیرد، اما اشتباه نکند.</li>
  <li>اولویت ۲: Reliability. اگر کاربر پول پرداخت کرد، حتماً باید بلیت صادر شود. داده نباید گم شود.</li>
  <li>اولویت ۳: Scalability. در فصل مسافرت ترافیک ۱۰ برابر می‌شود.</li>
  <li>اولویت ۴: Performance. زمان جستجوی بلیت باید معقول باشد.</li>
</ul>

<p dir="rtl" align="right">یا مثلا برای سیستم “ارسال نوتیفیکیشن لایک اینستاگرام”، این اولویت‌بندی  تغییری می‌کند؟</p>

<ul dir="rtl">
  <li>اولویت 1: Availability اگر نوتیفیکیشن چند ثانیه دیر برسد → قابل قبول  اما اگر سیستم چند لحظه down شود → کاربران میلیاردی هیچ نوتیفیکیشنی دریافت نمی‌کنند و اگر اگر سیستم کامل قطع شود → engagement کاهش شدید پیدا می‌کند</li>
  <li>اولویت 2: Scalability سیستم نوتیفیکیشن اینستاگرام با چالش‌های عظیم scale مواجه است. اگر کافی نباشد: گلوگاه ایجاد  می شود، نوتیفیکیشن ها صف می‌خورند و latency بالا می‌رود.</li>
  <li>اولویت 3: performance بسیار مهم است اما <strong>critical path نیست</strong>.</li>
</ul>

<p dir="rtl" align="right">نکته بسیار مهم در سیستم دیزاین که همیشه باید از خودمان بپرسیم: اگر سیستم خراب شود، بدترین اتفاق چیست؟ همین سوال ساده کل اولویت NFR ها را تعیین می‌کند.</p>

</div>
