# Gaming API Requirements

Этот документ фиксирует, каких HTTP-запросов не хватает, чтобы страница `/gaming` снова работала так, как раньше: реальные зоны, реальные места, реальные конфигурации ПК, реальные тарифы и реальная занятость по времени.

## Цель

Старый сценарий `/gaming` требовал:

- несколько игровых зон
- список мест внутри каждой зоны
- конфигурацию каждого места
- тарифы для каждой зоны
- публичную занятость мест по выбранной дате/слотам

Сейчас backend уже умеет:

- `GET /api/v1/public/gaming`
- `POST /api/v1/zone`
- `POST /api/v1/service`
- `POST /api/v1/booking`

Но этого недостаточно для полного parity, потому что через HTTP пока нельзя:

- создать `zone_tags`
- создать `computer_configurations`
- создать `zone_places`
- получить публичную занятость мест по времени

## Requirement

Для полного восстановления старого поведения `/gaming` backend должен поддерживать:

1. `POST /api/v1/zone-tag`
2. `GET /api/v1/zone-tag`
3. `POST /api/v1/configuration`
4. `GET /api/v1/configuration`
5. `POST /api/v1/place`
6. `GET /api/v1/place?zone_id=<id>`
7. `GET /api/v1/public/gaming?date=YYYY-MM-DD`
8. Либо отдельный `GET /api/v1/public/gaming/availability?zone_id=<id>&date=YYYY-MM-DD`

Минимальный набор данных, который должен видеть frontend в `public/gaming`:

- `zones[]`
- `zones[].name`
- `zones[].description`
- `zones[].details_json.prefix`
- `zones[].details_json.title`
- `zones[].details_json.hexColor`
- `zones[].places[]`
- `zones[].places[].label`
- `zones[].places[].configuration_id`
- `zones[].places[].specs`
- `zones[].services[]`
- `zones[].services[].name`
- `zones[].services[].price`
- `zones[].services[].currency`
- `zones[].services[].duration`
- `zones[].availability[]` или эквивалентный отдельный endpoint

## Expected Public Payload

Пример того, что должен возвращать `GET /api/v1/public/gaming?date=2026-05-01`, чтобы страница выглядела как раньше:

```json
{
  "timestamp": "2026-05-01T10:00:00Z",
  "zones": [
    {
      "id": 1,
      "name": "VIP Arena",
      "zone_type": "game",
      "zone_tag_id": 2,
      "capacity": 6,
      "description": "Флагманская зона для соревновательной игры.",
      "is_active": true,
      "details_json": {
        "prefix": "VIP",
        "title": "ARENA",
        "hexColor": "#22d3ee"
      },
      "places": [
        {
          "id": 11,
          "label": "VIP 1",
          "configuration_id": 5,
          "sort_order": 1,
          "is_active": true,
          "specs": [
            { "title": "Видеокарта", "value": "RTX 4090" },
            { "title": "Процессор", "value": "Intel Core i9 14900K" },
            { "title": "Монитор", "value": "360Hz" }
          ]
        }
      ],
      "services": [
        {
          "id": 21,
          "name": "VIP час",
          "duration": 60,
          "price": "250",
          "currency": "RUB",
          "description": "RTX 4090 • 360Hz"
        }
      ],
      "availability": [
        {
          "place_id": 11,
          "date": "2026-05-01",
          "slots": [
            { "time": "10:00", "taken": false },
            { "time": "11:00", "taken": true },
            { "time": "12:00", "taken": true }
          ]
        }
      ]
    }
  ]
}
```

## Existing Requests

Ниже запросы, которые уже можно делать сейчас.

### 1. Создать игровую зону

Требуется admin-cookie.

```bash
curl -X POST "http://localhost/api/v1/zone" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<ADMIN_SESSION_COOKIE>" \
  -d '{
    "name": "VIP Arena",
    "type": "game",
    "zone_tag_id": 1,
    "capacity": 6,
    "description": "Флагманская зона для соревновательной игры.",
    "is_active": true,
    "details_json": "{\"prefix\":\"VIP\",\"title\":\"ARENA\",\"hexColor\":\"#22d3ee\"}"
  }'
```

### 2. Создать тариф для игровой зоны

Требуется admin-cookie.

```bash
curl -X POST "http://localhost/api/v1/service" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<ADMIN_SESSION_COOKIE>" \
  -d '{
    "name": "VIP час",
    "zone_id": 1,
    "duration": 60,
    "price": "250",
    "currency": "RUB",
    "description": "RTX 4090 • 360Hz"
  }'
```

### 3. Прочитать публичный каталог gaming

```bash
curl "http://localhost/api/v1/public/gaming"
```

### 4. Создать бронь

Требуется авторизованный пользователь.

```bash
curl -X POST "http://localhost/api/v1/booking" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<USER_SESSION_COOKIE>" \
  -d '{
    "zone_id": 1,
    "service_id": 21,
    "place_id": 11,
    "start_time": "2026-05-01T10:00:00Z",
    "end_time": "2026-05-01T11:00:00Z",
    "participants": 1,
    "contact_name": "Yaroslav",
    "contact_email": "yaroslav@example.com",
    "contact_phone": "+79990000000",
    "details_json": "{}"
  }'
```

## Missing Requests

Ниже запросы, которых сейчас не хватает в backend, но именно они нужны, чтобы вернуть старую функциональность без ручных SQL-вставок.

### 1. Создать zone tag

```bash
curl -X POST "http://localhost/api/v1/zone-tag" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<ADMIN_SESSION_COOKIE>" \
  -d '{
    "name": "vip-gaming"
  }'
```

### 2. Создать конфигурацию ПК

```bash
curl -X POST "http://localhost/api/v1/configuration" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<ADMIN_SESSION_COOKIE>" \
  -d '{
    "zone_tag_id": 1,
    "specs_json": [
      { "title": "Видеокарта", "value": "RTX 4090" },
      { "title": "Процессор", "value": "Intel Core i9 14900K" },
      { "title": "Оперативная память", "value": "64GB DDR5" },
      { "title": "Монитор", "value": "360Hz" },
      { "title": "Периферия", "value": "Logitech G Pro" }
    ]
  }'
```

### 3. Создать место в игровой зоне

```bash
curl -X POST "http://localhost/api/v1/place" \
  -H "Content-Type: application/json" \
  -H "Cookie: pg_session=<ADMIN_SESSION_COOKIE>" \
  -d '{
    "zone_id": 1,
    "label": "VIP 1",
    "configuration_id": 5,
    "sort_order": 1,
    "is_active": true
  }'
```

### 4. Получить занятость мест по дате

Вариант A, если availability включается прямо в `public/gaming`:

```bash
curl "http://localhost/api/v1/public/gaming?date=2026-05-01"
```

Вариант B, если availability выделяется отдельно:

```bash
curl "http://localhost/api/v1/public/gaming/availability?zone_id=1&date=2026-05-01"
```

## Full Example Flow

Ниже целевой сценарий, который должен давать тот же функциональный эффект, что и старый фронтовый mock:

1. Создать `zone_tag`
2. Создать `configuration`
3. Создать `zone`
4. Создать `place`
5. Создать `service`
6. Прочитать `GET /api/v1/public/gaming?date=...`
7. Нарисовать табы, характеристики, цены, места и занятые часы
8. Создать бронь через `POST /api/v1/booking`

## Important Note

Без endpoint'ов для `configuration`, `place` и `availability` страница `/gaming` не сможет быть полностью такой же, как раньше. Сейчас можно восстановить только:

- реальные зоны
- реальные описания
- реальные тарифы

Но нельзя полноценно восстановить:

- реальные игровые места
- реальные конфигурации по местам
- реальную занятость по временным слотам
