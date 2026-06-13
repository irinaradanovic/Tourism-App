# Monitoring i tracing

Ovaj folder pokrece Grafana, Loki, Promtail i Tempo za observability demo.

## Pokretanje

Prvo pokrenuti monitoring stack, jer glavni `docker-compose.yml` koristi eksternu Docker mrezu `observability` koju ovaj stack kreira.

```bash
docker compose -f monitoring/docker-compose.yml up --build
```

Zatim u drugom terminalu pokrenuti aplikaciju:

```bash
docker compose up --build
```

Grafana je dostupna na:

```text
http://localhost:3000
```

Kredencijali:

```text
username: admin
password: admin
```

## Tracing demo za checkout

Trace je implementiran za tok:

```text
POST /api/purchase/checkout -> gateway -> purchase-service
```

`gateway` pravi server span za HTTP zahtev i client span za gRPC poziv. `purchase-service` preuzima trace context iz gRPC metadata i dodaje server span za `Checkout` i child span za `CheckoutCartAsync`.

Za test je potreban validan JWT token korisnika:

```bash
curl -X POST http://localhost/api/purchase/checkout \
  -H "Authorization: Bearer <validan_token>"
```

U Grafani otvoriti **Explore**, izabrati datasource **Tempo**, pa pretraziti trace-ove za servis `gateway` ili `purchase-service`. Ocekivani spanovi su:

```text
POST /api/purchase/checkout
grpc.Checkout purchase-service
PurchaseService.Checkout
PurchaseService.CheckoutCartAsync
```
