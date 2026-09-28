# SLO / SLA Template 

## Scope

باید دقیق مشخص شود:

- چه endpoint هایی داخل scope هستند
- measurement از کجا انجام می‌شود
- چه trafficی excluded است

## SLI ها

معمولاً:

- availability
- latency
- error rate
- freshness
- saturation

## SLO ها

باید: 

- measurable
- realistic
- business-aligned

باشند.

## Error Budget

بخش بسیار مهمی که اغلب فراموش می‌شود. مثلاً: 99.9% availability یعنی حدود 43 دقیقه downtime در ماه

### Policy مهم

اگر error budget سریع مصرف شود:

- release freeze
- rollback
- incident review

باید فعال شوند.
