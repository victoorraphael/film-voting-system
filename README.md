#Film Voting Plataform

Rank system where the user can upvote your favorites films, downvote other ones and add films

## Techs
- Golang
- gRPC
- MongoDB

## Endpoints

### GET /film/

```bash
    curl https://client-victoorraphael.cloud.okteto.net/film/
```

### Response

```json
[
    {
        "id": "613e9850601c53d9720db872",
        "name": "lalaland",
        "upvotes": 2,
        "downvotes": 1,
        "score": 1
    },
    {
        "id": "613e98d2601c53d9720db873",
        "name": "ghost warrior",
        "upvotes": 3,
        "downvotes": 1,
        "score": 2
    },
    {
        "id": "613fb99b0a79e4c2a2fddef8",
        "name": "cavaleiro fantasma",
        "downvotes": 1,
        "score": -1
    },
    {
        "id": "614497d718292a2cc603888c",
        "name": "quarto de guerra",
        "upvotes": 2,
        "score": 2
    }
]
```

### GET /film/:id/

```bash
    curl https://client-victoorraphael.cloud.okteto.net/film/613e98d2601c53d9720db873/
```

### Response

```json
{
  "film": {
    "id": "613e98d2601c53d9720db873",
    "name": "ghost warrior",
    "upvotes": 3,
    "downvotes": 1,
    "score": 2
  }
}
```

### POST /film/

```bash
    curl -X POST https://client-victoorraphael.cloud.okteto.net/film/ --header "Content-Type: application/json"  --data '{"name": "Sherk"}'
```

### Response

```json
{
  "film": {
    "id": "6144a2ef2b198a1a84f2703b",
    "name": "Sherk"
  }
}
```

### POST /film/upvote/:id/

```bash
    curl -X POST https://client-victoorraphael.cloud.okteto.net/film/upvote/6144a2ef2b198a1a84f2703b/
```

### Response

```json
{
  "success": true
}
```

### POST /film/downvote/:id/

```bash
    curl -X POST https://client-victoorraphael.cloud.okteto.net/film/downvote/6144a2ef2b198a1a84f2703b/
```

### Response

```json
{
  "success": true
}
```

### DELETE /film/:id/

```bash
    curl -X DELETE https://client-victoorraphael.cloud.okteto.net/film/6144a2ef2b198a1a84f2703b/
```

### Response

```json
{
  "success": true
}
```
