# Batching

Batching یعنی چند عملیات کوچک و مشابه را جمع کنیم و به‌جای اجرای جداگانه آن‌ها، در قالب یک عملیات بزرگ‌تر پردازش کنیم. هدف اصلی، کاهش هزینه‌های ثابت و تکرارشونده هر عملیات است؛ هزینه‌هایی مانند Network Round Trip، گرفتن Connection، Parse و Execute کردن Query، Transaction Commit، RPC Metadata، Serialization، System Call و Disk Flush.

بدون Batching:

```
Request 1 → Operation 1
Request 2 → Operation 2
Request 3 → Operation 3
Request 4 → Operation 4
```

با Batching:

```
Request 1
Request 2
Request 3
Request 4
→ One Batch
→ One Operation
```

Batching زمانی بیشترین ارزش را دارد که هر عملیات علاوه بر هزینه متناسب با حجم داده، یک Fixed Overhead نیز داشته باشد. با جمع‌کردن چند Item در یک Batch، این هزینه ثابت میان همه Itemها تقسیم می‌شود.

## Mental Model

فرض کنید باید صد بسته را از یک انبار به ساختمان دیگری منتقل کنیم. در روش اول، هر بار یک بسته برداشته، با خودرو به مقصد می‌رویم و دوباره برمی‌گردیم. در روش دوم، چندین بسته را داخل یک کامیون قرار داده و با یک بار رفت‌وآمد منتقل می‌کنیم.

```
Many Small Operations
→ Share One Expensive Fixed Cost
```

هرچند بارگیری، مدیریت و تخلیه یک Batch پیچیده‌تر است، هزینه رفت‌وآمد، راه‌اندازی و هماهنگی برای هر Item کاهش پیدا می‌کند.

## اثر Batching بر Performance

اثر اصلی Batching معمولاً افزایش Throughput و بهبود Resource Efficiency است:

```
Fewer Operations
→ Lower Fixed Overhead
→ More Work per Unit of Time
→ Higher Throughput
```

Batching می‌تواند Network Efficiency، Database Connection Usage، CPU Cost per Item، Serialization Overhead، Disk I/O Efficiency و Transaction Throughput را بهبود دهد. بااین‌حال، Itemها ممکن است برای تشکیل Batch منتظر بمانند؛ بنابراین Batching معمولاً میان Throughput و Latency یک Trade-off ایجاد می‌کند:

```
Throughput ↑
Resource Efficiency ↑
Per-Item Waiting Time May Increase
```

Batching بیشتر یک Throughput Optimization است و الزاماً Latency تک‌تک Itemها را کاهش نمی‌دهد.

## مثال Database و N+1 Query

فرض کنید باید اطلاعات هزار User از Database خوانده شود. روش ضعیف این است که برای هر User یک Query مستقل اجرا شود:

```
for _, id := range userIDs {
	user, err := repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	users = append(users, user)
}

Flow:
1000 User IDs → 1000 Separate SELECT Queries
```

برای هر Query هزینه‌های زیر تکرار می‌شوند:

```
Acquire Connection
→ Send Query
→ Network Round Trip
→ Parse and Execute
→ Receive Result
→ Release Connection
```

اگر هر Query به‌صورت Sequential حدود پنج میلی‌ثانیه طول بکشد:

```
1000 × 5ms ≈ 5 Seconds
```

با Batch Query می‌توان چند ID را در یک Query خواند:

```
SELECT user_id, name, email
FROM users
WHERE user_id IN (:ids);
```

مثلاً:

```
1000 User IDs
→ 10 Batches of 100 IDs
→ 10 Database Queries
```

اگر هر Batch حدود پانزده میلی‌ثانیه طول بکشد:

```
10 × 15ms ≈ 150ms
```

این اعداد صرفاً نمونه‌اند، اما نشان می‌دهند کاهش Round Trip و Connection Overhead می‌تواند اثر زیادی روی Performance داشته باشد.

## چرا یک Batch بسیار بزرگ مناسب نیست؟

ممکن است تصور شود تمام هزار ID باید در یک Query ارسال شوند:

```
1000 IDs → One Query
```

اما Batch بسیار بزرگ می‌تواند مشکلات زیر را ایجاد کند:

```
Large SQL Statement
Database Parameter Limit
High Memory Usage
Long Query Duration
Large Lock Window
Large Response
Timeout Risk
Poor Partial-Failure Handling
```

بنابراین هدف Batching ساخت بزرگ‌ترین Batch ممکن نیست؛ هدف پیدا کردن اندازه‌ای است که Fixed Overhead را کاهش دهد، بدون اینکه Latency، Memory Usage، Failure Scope یا Contention بیش از حد افزایش یابد.

```
Too Small Batch
→ Too Many Round Trips

Too Large Batch
→ High Latency, Memory and Failure Cost

Balanced Batch
→ Better Throughput with Controlled Latency
```

## Size-Based Batching

در Size-based Batching، وقتی تعداد Itemها به حد مشخصی برسد، Batch پردازش می‌شود:

```
Batch Size = 100

Collect Items
→ Count Reaches 100
→ Flush Batch
```

این مدل در Traffic بالا خوب عمل می‌کند، زیرا Batch سریع پر می‌شود. در Traffic پایین ممکن است Batch هیچ‌گاه کامل نشود و Itemها مدت زیادی منتظر بمانند:

```
Only 20 Items Arrive
Batch Size = 100
→ Batch Never Fills
→ Requests Wait Indefinitely
```

به همین دلیل Size-based Batching معمولاً نباید به‌تنهایی استفاده شود.

## Time-Based Batching

در Time-based Batching، پس از گذشت مدت مشخصی Batch پردازش می‌شود، حتی اگر به اندازه کامل نرسیده باشد:

```
Maximum Wait = 10ms

Collect Items
→ 10ms Passes
→ Flush Current Batch
```

این مدل Maximum Waiting Time را محدود می‌کند، اما در Traffic پایین ممکن است Batchهای کوچک‌تری تولید کند.

## Size-or-Time Batching

مدل رایج Production ترکیب محدودیت تعداد و زمان است:

```
Flush When:
Batch Size Reaches 100
OR
10ms Passes
```

رفتار این مدل:

```
High Traffic → Batch Fills by Size
Low Traffic → Batch Flushes by Timeout
```

به این ترتیب، در Traffic بالا Throughput مناسب حفظ می‌شود و در Traffic پایین هیچ Itemی برای مدت نامحدود منتظر نمی‌ماند.

## Byte-Based Batching

تعداد Itemها همیشه معیار کافی نیست:

```
100 Items × 1KB = 100KB
100 Items × 5MB = 500MB
```

بهتر است Batcher علاوه بر تعداد و زمان، محدودیت حجم Payload نیز داشته باشد:

```
Flush When:
Item Count ≥ 100
OR
Payload Size ≥ 1MB
OR
Wait Time ≥ 10ms
```

این مدل مانع تشکیل Batchهایی می‌شود که تعداد Item مناسبی دارند، اما از نظر Memory یا Network Payload بیش از حد بزرگ‌اند.

## مثال Logging

اگر Application برای هر Log یک Write و Flush جدا انجام دهد:

```
Log Event → Serialize → Disk Write → Flush
```

در نرخ ده هزار Log در ثانیه، تعداد زیادی عملیات کوچک ایجاد می‌شود. با Batching می‌توان Logها را جمع و یک‌جا نوشت:

```
Collect 500 Log Events
→ Serialize Batch
→ Write Once
```

یا:

```
Collect Logs for 100ms
→ Flush Batch
```

این روش Throughput را افزایش می‌دهد، اما اگر Process پیش از Flush شدن Buffer Crash کند، Logهای داخل Memory ممکن است از دست بروند:

```
Larger Batch
→ Better Throughput
→ More Buffered Data at Risk
→ Higher Flush Latency
```

برای Logهای حیاتی ممکن است Durable Buffer یا Write-ahead Log لازم باشد.

## مثال Message Broker

Producer می‌تواند Messageها را به‌جای ارسال جداگانه، به‌صورت Batch ارسال کند:

بدون Batch:

```
Message 1 → Network Request
Message 2 → Network Request
Message 3 → Network Request
```

با Batch:

```
Message 1
Message 2
Message 3
→ One Network Request
```

در Message Brokerهایی مانند Kafka معمولاً دو تنظیم اصلی وجود دارد:

```
Maximum Batch Size
Maximum Linger Time
```

Producer تا رسیدن Batch به اندازه مشخص یا گذشت مدت کوتاهی صبر می‌کند. Batch بزرگ‌تر Throughput را افزایش می‌دهد، اما Linger Time می‌تواند Latency هر Message را بیشتر کند.

## مثال Bulk API

به‌جای ارسال چند HTTP Request مستقل:

```
POST /users/1/activate
POST /users/2/activate
POST /users/3/activate
```

می‌توان Bulk Endpoint طراحی کرد:

```
POST /users/bulk-activate

{
  "user_ids": [1, 2, 3]
}

Flow:
Many Client Requests
→ One HTTP Request
→ One Authentication and Validation Phase
→ Batch Processing
→ One Response
```

مزایا:

```
Fewer HTTP Round Trips
Lower Header Overhead
Lower Authentication Overhead
Better Database Batching
```

اما Bulk API نیازمند تعریف روشن Partial Failure، Retry، Idempotency و Per-item Result است.

## Partial Failure

فرض کنید Batch شامل صد Item است و پنج Item شکست می‌خورند. دو مدل اصلی وجود دارد.

### All-or-Nothing

```
One Item Fails
→ Entire Batch Rolls Back
```

این مدل برای عملیاتی مناسب است که باید Transactional و Atomic باشند. مزیت آن Consistency Semantics ساده است؛ اما یک Item نامعتبر می‌تواند کل Batch را Fail کند:

```
One Bad Item → Entire Batch Failure
```

Batchهای مالی یا عملیاتی که باید به‌صورت یک واحد Commit شوند، ممکن است به این مدل نیاز داشته باشند.

### Partial Success

```
95 Items Succeed
5 Items Fail
→ Return Per-Item Results
```

نمونه Response:

```
{
  "results": [
    {
      "id": 1,
      "status": "success"
    },
    {
      "id": 2,
      "status": "invalid"
    },
    {
      "id": 3,
      "status": "success"
    }
  ]
}
```

این مدل Progress و Throughput بهتری ایجاد می‌کند، اما Client باید نتیجه هر Item را جداگانه مدیریت کند.

## Error Semantics در Bulk API

HTTP Status به‌تنهایی نمی‌تواند همیشه نتیجه تمام Itemهای Batch را بیان کند. ممکن است Response با `200 OK` اعلام کند Batch پردازش شده، اما تعدادی از Itemها شکست خورده‌اند. در برخی APIها نیز می‌توان از `207 Multi-Status` استفاده کرد.

مهم‌تر از Status Code، قرارداد روشن API است:

```
Was the Batch Accepted?
Which Items Succeeded?
Which Items Failed?
Which Failures Are Retryable?
```

هر Item باید Identifier یا Correlation ID داشته باشد تا نتیجه آن به Request اصلی Map شود.

## Retry و Idempotency

Retry کل Batch می‌تواند Itemهای موفق را دوباره اجرا کند:

```
100 Items Sent
→ 90 Succeeded
→ Response Lost
→ Client Retries Entire Batch
→ 90 Successful Items May Run Again
```

بنابراین Batching تقریباً همیشه به Idempotency Strategy نیاز دارد:

```
Per-Item Idempotency Key
Batch Idempotency Key
Upsert
Deduplication Table
Unique Business Constraint
```

نمونه Request:

```
{
  "batch_id": "batch-2026-001",
  "items": [
    {
      "idempotency_key": "item-1001",
      "user_id": 123
    }
  ]
}
```

Batch ID برای تشخیص Retry کل Batch مفید است، اما اگر Partial Success مجاز باشد، معمولاً Per-item Idempotency نیز لازم خواهد بود.

## Batching و Transaction

قرار دادن همه Itemها در یک Transaction مزایایی دارد:

```
Atomicity
Fewer Commits
Lower Commit Overhead
```

اما Batch بزرگ‌تر Transaction را طولانی‌تر می‌کند:

```
Longer Transaction
More Locks
Higher Contention
Larger Rollback
More Expensive Retry
```

رابطه کلی:

```
More Items per Transaction
→ Fewer Commits
→ Longer Lock and Failure Window
```

بنابراین Batch Size باید با توجه به Commit Cost، Lock Duration، Contention و Retry Cost تعیین شود.

## Database Insert Batching

روش ضعیف:

```
INSERT INTO events (...) VALUES (...);
INSERT INTO events (...) VALUES (...);
INSERT INTO events (...) VALUES (...);
```

روش بهتر:

```
INSERT INTO events (event_id, event_type, payload)
VALUES
    (:id1, :type1, :payload1),
    (:id2, :type2, :payload2),
    (:id3, :type3, :payload3);
```

بسته به Database و Driver می‌توان از Multi-row Insert، Array Binding، Prepared Statement یا Bulk API اختصاصی Database استفاده کرد. در Go باید قابلیت‌های Driver و Database هدف بررسی شوند؛ یک روش واحد برای تمام Databaseها وجود ندارد.

## Batching در Queue Consumer

پردازش تک‌به‌تک Messageها:

```
Poll One Message
→ Process
→ Insert One Row
→ Commit
→ Ack
```

پردازش Batch:

```
Poll 100 Messages
→ Process
→ Batch Insert
→ Commit
→ Ack Messages
```

Batching می‌تواند Throughput Consumer را افزایش دهد، اما Failure Handling را پیچیده‌تر می‌کند. اگر Database Commit موفق شود ولی Ack شکست بخورد:

```
Database Commit Success
→ Ack Failure
→ Messages Redelivered
```

Consumer باید Idempotent باشد. همچنین یک Message خراب نباید کل Batch را وارد Retry بی‌نهایت کند. Message نامعتبر ممکن است نیاز به Dead-letter Queue داشته باشد.

## Queueing Delay و Latency

اگر Maximum Batch Wait برابر ده میلی‌ثانیه باشد، اولین Item ممکن است تقریباً ده میلی‌ثانیه منتظر بماند، در حالی که آخرین Item درست پیش از Flush وارد شده و تقریباً بدون انتظار پردازش می‌شود:

```
First Item Wait ≈ 10ms
Last Item Wait ≈ 0ms
Average Wait ≈ 5ms
```

برای Background Processing این Delay ممکن است قابل قبول باشد، اما در API با Latency Budget بسیار محدود، همین چند میلی‌ثانیه مهم است.

Total Latency هر Item تقریباً از این اجزا تشکیل می‌شود:

```
Queue Wait
+ Batch Formation Wait
+ Batch Processing
+ Result Delivery
```

## Tail Latency

Batchهای بزرگ می‌توانند P95 و P99 را بدتر کنند. دلایل اصلی:

```
Waiting for Batch Formation
Slowest Item Delays the Batch
Large Serialization
Long Database Query
Large Response
Retry of a Large Batch
```

اگر نتیجه تمام Itemها فقط پس از پایان کامل Batch برگردد:

```
Batch Latency ≈ Slowest Item Latency
```

## Head-of-Line Blocking

فرض کنید یک Batch صد Item دارد:

```
99 Items = 5ms
1 Item = 2s
```

اگر همه Callerها منتظر اتمام Batch باشند:

```
All 100 Items Wait 2 Seconds
```

این حالت Head-of-Line Blocking ایجاد می‌کند. راهکارهای ممکن:

```
Smaller Batches
Group Similar Work
Per-Item Timeout
Partial Result Delivery
Separate Slow Items
```

Itemهایی با Cost یا SLA بسیار متفاوت بهتر است در یک Batch مشترک قرار نگیرند.

## Homogeneous Batching

Batching زمانی مؤثرتر است که Itemها ویژگی‌های مشابه داشته باشند:

```
Same Operation
Similar Size
Similar Cost
Same Destination
Same Consistency Requirement
```

برای مثال، Batch کردن صد Event مشابه بهتر از ترکیب Readهای کوچک، Writeهای بزرگ، Remote Callهای کند و تراکنش‌های حیاتی در یک Batch است.

Grouping می‌تواند بر اساس این موارد انجام شود:

```
Operation Type
Tenant
Database Shard
Destination Service
Priority
Payload Size
```

## Batching در سیستم Sharded

در سیستم Sharded، یک Batch ورودی ممکن است به چند Batch مقصدی تقسیم شود:

```
Input Batch
→ Group by Shard
→ Batch for Shard A
→ Batch for Shard B
→ Batch for Shard C
```

مثلاً:

```
shard = hash(userID) % shardCount
```

اجرای یک Batch روی چند Shard ممکن است Scatter-Gather ایجاد کند، Partial Failure را پیچیده کند و Latency را به کندترین Shard وابسته سازد. بنابراین Batch باید تا حد امکان Destination-aware باشد.

## انتخاب Batch Size

Batch Size ثابت و عمومی برای همه سیستم‌ها وجود ندارد. اندازه مناسب به عوامل زیر بستگی دارد:

```
Item Size
Database Limits
Payload Size
Source Latency
Transaction Cost
Memory Budget
Latency SLO
Traffic Rate
Failure Cost
```

روش مناسب:

```
Start Conservative
→ Benchmark
→ Load Test
→ Observe Throughput and P99
→ Tune
```

می‌توان اندازه‌های مختلف را آزمایش کرد:

```
Batch Size = 1, 10, 50, 100, 500, 1000
```

Metricهای مقایسه:

```
Items per Second
Batch Latency
Per-Item Latency
Database CPU
Connection Pool Usage
Memory Consumption
Failure Rate
P95/P99
```

معمولاً بعد از نقطه‌ای مشخص، افزایش Batch Size بهبود Throughput کمی ایجاد می‌کند، اما Memory، Latency و Failure Scope را افزایش می‌دهد.

## Backpressure در Batcher

Batcher نباید Queue نامحدود داشته باشد:

```
Incoming Items
→ Unbounded In-Memory Queue
→ Memory Growth
→ OOM
```

طراحی مناسب:

```
Bounded Queue
→ Block, Reject or Shed When Full
```

هنگام پر شدن Queue باید Policy مشخص باشد:

```
Block Producer
Return Error
Drop Low-Priority Item
Write to Durable Queue
```

Bounded Queue باعث می‌شود ظرفیت محدود پردازش به Producer منتقل شود و سیستم به‌جای رشد نامحدود Memory، رفتار کنترل‌شده داشته باشد.

## پیاده‌سازی Size-or-Time Batcher در Go

ساختار ساده:

```
type Item struct {
	ID    string
	Value []byte
}

type BatchProcessor interface {
	Process(ctx context.Context, items []Item) error
}

type Batcher struct {
	maxItems int
	maxWait  time.Duration
	input    chan Item
	process  BatchProcessor
}
```

Constructor با Validation:

```
func NewBatcher(
	maxItems int,
	maxWait time.Duration,
	queueSize int,
	process BatchProcessor,
) (*Batcher, error) {
	if maxItems <= 0 {
		return nil, errors.New("max items must be positive")
	}
	if maxWait <= 0 {
		return nil, errors.New("max wait must be positive")
	}
	if queueSize <= 0 {
		return nil, errors.New("queue size must be positive")
	}
	if process == nil {
		return nil, errors.New("batch processor is required")
	}

	return &Batcher{
		maxItems: maxItems,
		maxWait:  maxWait,
		input:    make(chan Item, queueSize),
		process:  process,
	}, nil
}
```

Loop اصلی:

```
func (b *Batcher) Run(ctx context.Context) error {
	timer := time.NewTimer(b.maxWait)
	defer timer.Stop()

	batch := make([]Item, 0, b.maxItems)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}

		items := append([]Item(nil), batch...)
		batch = batch[:0]

		return b.process.Process(ctx, items)
	}

	for {
		select {
		case <-ctx.Done():
			if err := flush(); err != nil {
				return errors.Join(ctx.Err(), err)
			}
			return ctx.Err()

		case item, ok := <-b.input:
			if !ok {
				return flush()
			}

			batch = append(batch, item)
			if len(batch) >= b.maxItems {
				if err := flush(); err != nil {
					return err
				}

				resetTimer(timer, b.maxWait)
			}

		case <-timer.C:
			if err := flush(); err != nil {
				return err
			}

			timer.Reset(b.maxWait)
		}
	}
}
```

مدیریت امن Timer:

```
func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	timer.Reset(duration)
}
```

ارسال Item:

```
func (b *Batcher) Add(ctx context.Context, item Item) error {
	select {
	case b.input <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

Flow:
Add
→ Bounded Input Queue
→ Collect Items
→ Flush by Size or Time
→ Process Batch
```

این نمونه آموزشی است. در Production باید درباره Processing Timeout، Retry، Partial Failure، Per-item Result، Panic Recovery، Shutdown Semantics و Error Classification تصمیم‌گیری شود.

## Blocking Processor و Batch Workerها

در نمونه بالا، `Process` در همان Goroutine جمع‌کننده اجرا می‌شود. اگر پردازش Batch دو ثانیه طول بکشد، Collector در این مدت Item جدید دریافت نمی‌کند:

```
Flush Starts
→ Collector Stops Receiving
→ Input Queue Fills
```

این رفتار می‌تواند Backpressure طبیعی ایجاد کند و در برخی سیستم‌ها مطلوب باشد. اگر Throughput بیشتری لازم باشد، Collector می‌تواند Batchهای تشکیل‌شده را به تعداد محدودی Worker تحویل دهد:

```
Item Collector
→ Bounded Batch Queue
→ Fixed Batch Workers
```

تعداد Workerها باید محدود باشد. ایجاد یک Goroutine جدید برای هر Batch می‌تواند دوباره Concurrency را Unbounded کند و Database، Broker یا Remote Service را Saturate سازد.

## Result-Aware Batching

در بسیاری از سیستم‌ها هر Caller باید نتیجه Item خودش را دریافت کند. در این حالت Request می‌تواند Response Channel اختصاصی داشته باشد:

```
type Result struct {
	Value []byte
	Err   error
}

type BatchRequest struct {
	ID       string
	Response chan Result
}

Flow:
Caller
→ Submit Request
→ Wait on Personal Response Channel
→ Batcher Collects Requests
→ Execute One Batch Operation
→ Match Results by ID
→ Respond to Each Caller
```

این Pattern در DataLoaderها و Request-scoped Batching رایج است.

## DataLoader Pattern

در GraphQL یا سیستم‌های Resolver-based ممکن است صد Resolver اطلاعات صد User را جداگانه درخواست کنند.

بدون DataLoader:

```
Resolver 1 → SELECT User 1
Resolver 2 → SELECT User 2
...
Resolver 100 → SELECT User 100
```

با DataLoader:

```
Resolver Requests
→ Collect IDs During a Short Window
→ One Batch Query
→ Map Results Back to Resolvers

Query:
SELECT user_id, name, email
FROM users
WHERE user_id IN (...);
```

DataLoader معمولاً دو Pattern را ترکیب می‌کند:

```
Batching
+ Request-Scoped Deduplication
```

اگر چند Resolver یک User را درخواست کنند، ID آن فقط یک بار در Batch قرار می‌گیرد.

## Ordering نتایج

Database تضمین نمی‌کند نتایج Query با ترتیب IDهای ورودی برگردند:

```
Input IDs = [7, 2, 9]
Database Result = [2, 7, 9]
```

بنابراین نتایج باید با Identifier Map شوند:

```
usersByID := make(map[int64]User, len(users))

for _, user := range users {
	usersByID[user.ID] = user
}

for _, id := range requestedIDs {
	user, ok := usersByID[id]
	// Build the output in the original request order.
}
```

اعتماد به ترتیب خروجی Database می‌تواند باعث برگرداندن نتیجه اشتباه به Caller شود.

## Timeout و Cancellation

در Batching معمولاً سه زمان متفاوت وجود دارد:

```
Maximum Batch Formation Wait
Batch Processing Timeout
Individual Caller Deadline
```

اگر Caller فقط پنج میلی‌ثانیه Deadline داشته باشد، اما Batcher تا ده میلی‌ثانیه برای تشکیل Batch صبر کند، Request پیش از Flush منقضی می‌شود. بنابراین Batching برای مسیرهایی با Deadline بسیار کوتاه ممکن است مناسب نباشد.

اگر Caller Cancel شود اما Item قبلاً داخل Queue قرار گرفته باشد، دو انتخاب وجود دارد:

```
Remove Cancelled Item Before Flush
OR
Process It Anyway and Discard the Result
```

حذف Item از Queue و Batch پیچیدگی بیشتری دارد؛ پردازش آن نیز ممکن است Work بیهوده ایجاد کند. انتخاب مناسب به هزینه عملیات، Deadline Semantics و طراحی Queue بستگی دارد.

## Failure Modeهای مهم

### Batch بیش از حد بزرگ

```
Memory Spike
Timeout
Long Transaction
Lock Contention
Large Rollback
High Retry Cost
```

### Batch بیش از حد کوچک

```
Little Throughput Gain
Too Many Round Trips
Batching Complexity Without Meaningful Benefit
```

### نبود Time-Based Flush

اگر Batch فقط براساس Size Flush شود، در Traffic پایین Itemها ممکن است برای همیشه منتظر بمانند.

### Poison Item

یک Item نامعتبر ممکن است کل Batch را Fail کند:

```
One Invalid Item
→ Entire Batch Retry
→ Same Invalid Item Fails Again
```

راهکارها:

```
Per-Item Validation
Partial Success
Split Failed Batch
Dead-Letter Queue
```

### Retry Storm

```
Large Batch Failure
→ Immediate Retry
→ Source Overload
```

راهکارها:

```
Exponential Backoff
Jitter
Retry Budget
Smaller Retry Batch
```

### Duplicate Processing

گم‌شدن Response یا Ack می‌تواند باعث اجرای دوباره Batch شود. Idempotency و Deduplication لازم است.

### Queue Growth

اگر Producer سریع‌تر از Processor باشد، Queue پر می‌شود. Queue باید Bounded و رفتار هنگام Capacity Exhaustion مشخص باشد.

### Shutdown Data Loss

اگر Process هنگام وجود Itemهای Pending خاموش شود، داده ممکن است از دست برود. Shutdown Flow مناسب:

```
Stop Accepting New Items
→ Flush Pending Batch
→ Wait for In-Flight Work
→ Exit
```

اگر Delivery Guarantee مهم باشد، Buffer صرفاً In-memory کافی نیست و باید Durable Queue یا Persistent Log استفاده شود.

## Observability

Metricهای مهم:

```
batch_created_total
batch_processed_total
batch_failed_total
batch_items_total
batch_size
batch_payload_bytes
batch_wait_duration_seconds
batch_processing_duration_seconds
batch_queue_depth
batch_queue_capacity
batch_flush_by_size_total
batch_flush_by_timeout_total
batch_partial_failure_total
batch_retry_total
batch_dropped_items_total
```

Distributionهای مهم:

```
Batch Size Distribution
Batch Formation Wait Distribution
Batch Processing P95/P99
Per-Item End-to-End Latency
Queue Depth Over Time
Payload Size Distribution
```

Average Batch Size به‌تنهایی کافی نیست. ممکن است Average مناسب باشد، اما تعداد محدودی Item در Traffic پایین تا Maximum Wait کامل منتظر بمانند و P99 بدی ایجاد کنند.

Alertهای مهم:

```
Batch Queue Near Capacity
Batch Processing P99 Increases
Flush-by-Timeout Dominates Unexpectedly
Average Batch Size Drops
Partial Failure Rate Increases
Retry Rate Increases
Dropped Item Count Is Non-Zero
Shutdown Flush Fails
```

اگر بیشتر Batchها با Timeout و اندازه بسیار کوچک Flush شوند، ممکن است Traffic برای Batching کافی نباشد یا Configuration نامناسب باشد.

## Production Example: Analytics Events

فرض کنید سرویس در هر ثانیه پنجاه هزار Event تولید می‌کند.

بدون Batching:

```
50K Events/sec
→ 50K Network Sends
→ 50K Broker or Database Operations
```

Policy نمونه:

```
Max Items = 500
Max Payload = 1MB
Max Wait = 20ms

Flow:
Collect Events
→ Flush by Count, Bytes or Time
→ Send to Broker or Database
```

در Traffic بالا Batch با Size یا Byte Limit پر می‌شود. در Traffic پایین Max Wait اجازه نمی‌دهد Event بیش از بیست میلی‌ثانیه منتظر بماند. اگر Delivery مهم باشد، باید از Durable Queue یا Write-ahead Log استفاده شود، زیرا Buffer حافظه‌ای در Crash از بین می‌رود.

## Production Example: Notification Service

فرض کنید Preference هزار User باید خوانده شود.

روش ضعیف:

```
1000 Users
→ 1000 Preference Queries
```

روش بهتر:

```
Group User IDs by Database Shard
→ Batch Query per Shard
→ Map Results by User ID
→ Continue Notification Processing
```

اگر تعدادی User وجود نداشته باشند، نتیجه باید Per-user باشد و کل Batch نباید لزوماً Fail شود.

## Production Example: Payment Settlement

در Settlement، تراکنش‌ها معمولاً در Batchهای مشخص جمع می‌شوند. Batching تعداد Write و Commit را کاهش می‌دهد، اما Requirements سخت‌تری دارد:

```
Correct Ordering
No Duplicate Settlement
Auditable Batch ID
Durable Input
Reconciliation
Atomic File Generation
```

در این سیستم Throughput تنها معیار نیست. Batch باید قابل Audit، Resume و Retry امن باشد. Batch ID، Item ID، Reconciliation Record و Idempotency حیاتی‌اند.

## Review Scenario: Telemetry Service

فرض کنید:

```
Incoming Events = 100K/sec
Current Design = One Database Insert per Event
Database Safe Transactions = 5K/sec
Average Event Size = 2KB
Maximum Acceptable Added Latency = 50ms
Some Events May Be Invalid
```

### مشکل طراحی فعلی

سرویس صد هزار Transaction در ثانیه تولید می‌کند، در حالی که Database فقط پنج هزار Transaction در ثانیه ظرفیت امن دارد. حتی اگر Database از نظر Row Throughput توان کافی داشته باشد، هزینه Network Round Trip، Transaction و Commit باعث Saturation می‌شود.

### Strategy مناسب

```
Bounded Input Queue
→ Validate Events
→ Batch by Count, Bytes or Time
→ Bulk Insert
→ Commit
```

Policy اولیه:

```
Maximum Items = 500
Maximum Payload = 1MB
Maximum Wait = 20ms
```

چون هر Event حدود دو کیلوبایت است، محدودیت یک مگابایت تقریباً با پانصد Event هماهنگ است.

### دلیل انتخاب Max Wait بیست میلی‌ثانیه

حداکثر Latency اضافه‌شده پنجاه میلی‌ثانیه است. تمام این Budget نباید صرف انتظار تشکیل Batch شود، زیرا Queueing، Batch Processing و Database Latency نیز زمان مصرف می‌کنند. بنابراین بیست میلی‌ثانیه می‌تواند نقطه شروع محافظه‌کارانه‌ای باشد.

### مدیریت Event نامعتبر

Validation ساده باید پیش از Bulk Insert انجام شود:

```
Valid Events → Batch Insert
Invalid Events → Reject or Dead-Letter
```

اگر فقط Database Constraint بتواند خطا را تشخیص دهد، Processor باید از Per-item Result، Split Batch یا Isolation Strategy پشتیبانی کند تا یک Event کل Batch را Poison نکند.

### گم‌شدن Ack یا Response

اگر Bulk Insert موفق شود ولی Ack از دست برود، Eventها ممکن است دوباره ارسال شوند. هر Event باید ID یکتا داشته باشد و Database Constraint یا Deduplication Strategy از Insert تکراری جلوگیری کند.

### Transaction Size

یک Transaction برای پانصد Event ممکن است مناسب باشد، اما باید با Load Test تأیید شود. Batch بزرگ‌تر Commit Overhead را کاهش می‌دهد، اما Rollback، Lock Duration و Retry Cost را افزایش می‌دهد.

### Metricهای مهم

```
Events per Second
Database Transactions per Second
Batch Size and Bytes
Batch Formation Wait
Bulk Insert P99
Queue Depth
Invalid Event Rate
Duplicate Rate
Retry Rate
Dropped Events
```

### نتیجه تقریبی

اگر Average Batch Size پانصد باشد:

```
100K Events/sec ÷ 500
≈ 200 Batch Transactions/sec
```

Transaction Count از حدود صد هزار به دویست Transaction در ثانیه کاهش پیدا می‌کند. این محاسبه تخمینی است و هزینه واقعی Bulk Insert باید با Benchmark و Load Test اندازه‌گیری شود.

## نکات Production-oriented

در طراحی Production باید پاسخ این پرسش‌ها مشخص باشد:

```
What Triggers a Flush?
What Is the Maximum Item Count?
What Is the Maximum Payload Size?
What Is the Maximum Wait?
Is the Input Queue Bounded?
What Happens When the Queue Is Full?
Are Partial Failures Allowed?
How Are Failed Items Reported?
How Are Retries Performed?
Are Operations Idempotent?
How Is Graceful Shutdown Handled?
How Are Slow or Invalid Items Isolated?
How Is Batch Ordering Preserved?
```

Batching بدون پاسخ روشن به این پرسش‌ها ممکن است Throughput را بهتر کند، اما Reliability، Tail Latency و Operational Complexity را بدتر سازد.

## نکات Interview-oriented

در مصاحبه، صرفاً گفتن «درخواست‌ها را Batch می‌کنیم» کافی نیست. باید مشخص شود Batching چه هزینه‌ای را کاهش می‌دهد و چه Trade-offهایی ایجاد می‌کند:

```
Batching Reduces:
Network Round Trips
Transaction Count
Serialization Overhead
Connection Usage
Per-Item Fixed Cost

Batching Adds:
Queueing Delay
Partial-Failure Complexity
Retry and Idempotency Challenges
Memory Buffering
Head-of-Line Blocking
Larger Failure Scope
```

همچنین باید مشخص شود Batch براساس چه شرطی Flush می‌شود:

```
Count
Payload Bytes
Time
or a Combination
```

و Source چگونه محافظت می‌شود:

```
Bounded Queue
Limited Batch Workers
Processing Timeout
Retry with Backoff
```

## جمع‌بندی و Mental Model نهایی

Batching چند عملیات کوچک و مشابه را جمع می‌کند تا هزینه‌های ثابت میان آن‌ها به اشتراک گذاشته شود:

```
Many Small Operations
→ One Controlled Batch
→ Fewer Round Trips and Commits
→ Higher Throughput
```

اما Batch بزرگ‌تر همیشه بهتر نیست:

```
Small Batch
→ Lower Waiting Time
→ Lower Efficiency

Large Batch
→ Higher Efficiency
→ More Latency, Memory and Failure Cost
```

مدل مناسب Production:

```
Collect Similar Items
→ Flush by Count, Bytes or Time
→ Process with Bounded Concurrency
→ Handle Partial Failure and Retry Safely
→ Preserve Idempotency and Ordering
→ Observe Batch Size, Queueing and P99
```

مدل ذهنی نهایی این است: **Batching هزینه ثابت چند عملیات را به اشتراک می‌گذارد و Throughput را افزایش می‌دهد، اما این بهبود را با Queueing Delay، Failure Scope بزرگ‌تر و پیچیدگی بیشتر در Retry، Idempotency و Partial Failure به دست می‌آورد.**
