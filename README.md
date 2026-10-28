# Interseguro Coding Challenge

Solución al reto técnico de Interseguro. Implementa dos APIs (Go + Node.js) que colaboran para factorizar/rotar matrices y calcular estadísticas.

## ⚠️ Decisión ante la ambigüedad del enunciado

El PDF menciona "rotación" en la sección de arquitectura (p.2) y "factorización QR" en funcionalidad (p.3). Ambas operaciones están implementadas:

| Endpoint Go                  | Operación                  | Envía a Node          |
|------------------------------|----------------------------|-----------------------|
| `POST /api/v1/qr`            | Factorización QR (Q, R)    | `[Q, R]`              |
| `POST /api/v1/rotate`        | Rotación 90° horaria       | `[rotated]`           |

## Arquitectura

Cliente → Go API (cálculo matricial) → Node API (estadísticas) → Cliente

## Endpoints

### Go API (puerto 8080)
- `GET  /health`
- `POST /api/v1/qr`     → `{ "matrix": [[1,2],[3,4]] }`
- `POST /api/v1/rotate` → `{ "matrix": [[1,2],[3,4]] }`

### Node API (puerto 3000)
- `GET  /health`
- `POST /api/v1/statistics` → `{ "matrices": [[[...]], [[...]]] }`

## Ejecución local

```bash
cp .env.example .env
docker compose up --build

## 🌐 Demo en vivo (Render)

| Servicio | URL |
| :--- | :--- |
| Frontend + Node API + GO API | https://interseguro-challenge-ww9k.onrender.com 
