<div dir="rtl" align="right">

<h1 dir="rtl" align="right">Percentile-based Thinking</h1>

<p dir="rtl" align="right">هرگز با <strong>Average</strong>  کار نکنید. در خیلی از موارد، یک سری داده ( که ممکن است تعداد آن ها بسیار اندک می باشد) در میانگین تاثیر زیادی گذاشته و عدد میانگین را از عدد واقعی بسیار دور می‌کنند. از طرف دیگر، اگر بیشتر درخواست‌ها سریع باشند ولی درصد کمی خیلی کند باشند، average ممکن است هنوز قابل قبول دیده شود، در حالی که کاربران واقعی اذیت می‌شوند. برای همین در سیستم‌دیزاین و production monitoring از percentile استفاده می‌کنیم.  عموما به جای میانگین  از درصد های زیر برای performance استفاده می کنیم</p>

<table dir="rtl">
  <thead>
    <tr>
      <th align="right">Percentile</th>
      <th align="right">نوع</th>
      <th align="right">معنی</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td align="right">P50</td>
      <td align="right">تجربه کاربر معمولی</td>
      <td align="right">یعنی ۵۰٪ درخواست‌ها سریع‌تر از این مقدار جواب گرفته‌اند<br>P50 = 100ms یعنی نصف request ها زیر ۱۰۰ میلی‌ثانیه بوده‌اند.</td>
    </tr>
    <tr>
      <td align="right">P90</td>
      <td align="right">کاربران سنگین</td>
      <td align="right">یعنی 90 درصد درخواست‌ها سریع‌تر از این مقدار جواب گرفته‌اند</td>
    </tr>
    <tr>
      <td align="right">P95</td>
      <td align="right">اکثر کاربران</td>
      <td align="right">یعنی 95 درصد درخواست‌ها سریع‌تر از این مقدار جواب گرفته‌اند<br>این معمولاً برای SLO خیلی استفاده می‌شود.</td>
    </tr>
    <tr>
      <td align="right">P99</td>
      <td align="right">بدترین تجربه‌ها</td>
      <td align="right">یعنی 99 درصد درخواست‌ها سریع‌تر از این مقدار جواب گرفته‌اند</td>
    </tr>
  </tbody>
</table>

<h2 dir="rtl" align="right">چرا P95 و P99 مهم‌تر از Average هستند؟</h2>

<p dir="rtl" align="right">چون مشکل واقعی معمولاً در tail latency است. Tail latency یعنی کندترین بخش درخواست‌ها. تصور کنید endpoint زیر را داریم:</p>

</div>

```
P50: 80ms
P90: 200ms
P95: 600ms
P99: 4s
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">این یعنی اکثر کاربران تجربه خوبی دارند، اما بدترین ۱٪ تجربه بسیار بدی دارند.</p>

<p dir="rtl" align="right">حالا چرا این مهم است؟ چون در سیستم‌های مدرن، یک درخواست کاربر ممکن است به چند سرویس دیگر وابسته باشد.</p>

<p dir="rtl" align="right">مثلاً request باز کردن صفحه home در Uber&rlm;:</p>

</div>

```
User API
Pricing API
Driver Location API
Map API
Promotion API
Payment API
```

<div dir="rtl" align="right">

<p dir="rtl" align="right">اگر هر کدام از این سرویس‌ها P99 بدی داشته باشند، احتمال اینکه کل request کند شود بالا می‌رود.</p>

<p dir="rtl" align="right">به این می‌گویند:  Tail Latency Amplification یعنی کندی‌های کوچک در سرویس‌های مختلف جمع می‌شوند و تجربه نهایی کاربر را خراب می‌کنند.</p>

<p dir="rtl" align="right">استاندارد حرفه‌ای: SLO روی Mean تعریف نمی‌شود</p>

<ul dir="rtl">
  <li>SLO روی P95 یا P99 تعریف می‌شود</li>
  <li>نه روی Mean</li>
</ul>

<p dir="rtl" align="right">چون average نمی‌گوید چند درصد کاربران تجربه بد دارند. یک SLO حرفه‌ای معمولاً این شکلی است:</p>

</div>

```
For /payment/authorize:
P95 latency < 400ms
Error rate < 0.1%
Availability > 99.9%
```
