# =====================================================================
# Smoke Test - Interseguro Coding Challenge
# ---------------------------------------------------------------------
# Valida end-to-end los 4 flujos criticos del sistema:
#   1. Health de ambas APIs (Go + Node)
#   2. Factorizacion QR de matriz cuadrada 2x2
#   3. Factorizacion QR de matriz ancha 2x3 (valida padding)
#   4. Rotacion 90 grados CW de matriz 2x3
#
# Uso:
#   cd C:\Users\USER\Desktop\interseguro-challenge
#   .\scripts\smoke-test.ps1
# =====================================================================

$ErrorActionPreference = "Stop"
$goApi   = "http://localhost:8080"
$nodeApi = "http://localhost:3000"

$passed = 0
$total  = 4

function Write-Step {
    param([string]$Message)
    Write-Host "`n[$($script:passed + 1)/$total] $Message" -ForegroundColor Cyan
}

function Write-Ok {
    param([string]$Message)
    Write-Host "  OK  $Message" -ForegroundColor Green
    $script:passed++
}

# ---------------------------------------------------------------------
# Test 1: Healthchecks
# ---------------------------------------------------------------------
Write-Step "Healthchecks de ambas APIs..."

try {
    $goHealth   = (Invoke-RestMethod -Uri "$goApi/health"   -TimeoutSec 5).status
    $nodeHealth = (Invoke-RestMethod -Uri "$nodeApi/health" -TimeoutSec 5).status
}
catch {
    Write-Host "  ERROR: no se pudo conectar a las APIs. Esta corriendo docker compose up?" -ForegroundColor Red
    throw
}

if ($goHealth -ne "ok")   { throw "Go API no healthy: $goHealth" }
if ($nodeHealth -ne "ok") { throw "Node API no healthy: $nodeHealth" }

Write-Ok "Go=ok, Node=ok"

# ---------------------------------------------------------------------
# Test 2: QR de matriz cuadrada 2x2
# ---------------------------------------------------------------------
Write-Step "Factorizacion QR de matriz 2x2..."

$body = '{"matrix":[[1,2],[3,4]]}'
$r = Invoke-RestMethod -Method Post -Uri "$goApi/api/v1/qr" `
    -ContentType "application/json" -Body $body -TimeoutSec 10

if ($r.operation -ne "QR_FACTORIZATION") {
    throw "operacion incorrecta: $($r.operation)"
}
if ($r.qr.q.rows -ne 2 -or $r.qr.q.columns -ne 2) {
    throw "Q incorrecta: $($r.qr.q.rows)x$($r.qr.q.columns)"
}
if ($r.qr.r.rows -ne 2 -or $r.qr.r.columns -ne 2) {
    throw "R incorrecta: $($r.qr.r.rows)x$($r.qr.r.columns)"
}
if ($r.statistics.count -ne 8) {
    throw "count incorrecto: $($r.statistics.count)"
}

Write-Ok "Q=2x2, R=2x2, count=$($r.statistics.count)"

# ---------------------------------------------------------------------
# Test 3: QR de matriz ancha 2x3 (valida el padding con ceros)
# ---------------------------------------------------------------------
Write-Step "Factorizacion QR de matriz ancha 2x3 (test del padding)..."

$body = '{"matrix":[[1,2,3],[4,5,6]]}'
$r = Invoke-RestMethod -Method Post -Uri "$goApi/api/v1/qr" `
    -ContentType "application/json" -Body $body -TimeoutSec 10

if ($r.qr.q.rows -ne 2 -or $r.qr.q.columns -ne 2) {
    throw "Q incorrecta para matriz ancha: $($r.qr.q.rows)x$($r.qr.q.columns)"
}
if ($r.qr.r.rows -ne 2 -or $r.qr.r.columns -ne 3) {
    throw "R incorrecta para matriz ancha: $($r.qr.r.rows)x$($r.qr.r.columns)"
}
if ($r.statistics.count -ne 10) {
    throw "count incorrecto: $($r.statistics.count)"
}

# Verificar que R es triangular superior (R[1][0] debe ser 0)
if ([Math]::Abs($r.qr.r.data[1][0]) -gt 1e-9) {
    throw "R no es triangular superior: R[1][0]=$($r.qr.r.data[1][0])"
}

Write-Ok "Q=2x2, R=2x3 triangular, count=$($r.statistics.count)"

# ---------------------------------------------------------------------
# Test 4: Rotacion 90 grados CW de matriz 2x3
# ---------------------------------------------------------------------
Write-Step "Rotacion 90 grados CW de matriz 2x3..."

$body = '{"matrix":[[1,2,3],[4,5,6]]}'
$r = Invoke-RestMethod -Method Post -Uri "$goApi/api/v1/rotate" `
    -ContentType "application/json" -Body $body -TimeoutSec 10

if ($r.operation -ne "ROTATE_90_CW") {
    throw "operacion incorrecta: $($r.operation)"
}
if ($r.rotated.rows -ne 3 -or $r.rotated.columns -ne 2) {
    throw "dimensiones incorrectas: $($r.rotated.rows)x$($r.rotated.columns)"
}

$expected = @(
    @(4, 1),
    @(5, 2),
    @(6, 3)
)

for ($i = 0; $i -lt 3; $i++) {
    for ($j = 0; $j -lt 2; $j++) {
        if ($r.rotated.data[$i][$j] -ne $expected[$i][$j]) {
            throw "rotacion incorrecta en [$i][$j]: esperado $($expected[$i][$j]), obtenido $($r.rotated.data[$i][$j])"
        }
    }
}

Write-Ok "rotated = [[4,1],[5,2],[6,3]]"

# ---------------------------------------------------------------------
# Resumen
# ---------------------------------------------------------------------
Write-Host "`n=====================================================" -ForegroundColor Green
Write-Host " TODOS LOS SMOKE TESTS PASARON ($passed/$total)" -ForegroundColor Green
Write-Host "=====================================================" -ForegroundColor Green