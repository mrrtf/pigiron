[ordered.json](ordered.json) and [regular.json](regular.json) are two example request for the padlist endpoint.

```
curl -d '@ordered.json' -X POST http://localhost:8080/v2/padlist | jq .
```
