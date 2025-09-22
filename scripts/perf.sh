#!/usr/bin/env bash
# scripts/perf.sh - Run performance tests using k6
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.."; pwd)"
K6_SCRIPT="${K6_SCRIPT:-$ROOT/perf/k6/resolve-smoke.js}"

echo "🚀 Running performance tests..."

# Check if k6 script exists, create a basic one if not
if [ ! -f "$K6_SCRIPT" ]; then
    echo "⚠️  k6 script not found at $K6_SCRIPT"
    echo "   Creating basic smoke test..."

    mkdir -p "$(dirname "$K6_SCRIPT")"
    cat > "$K6_SCRIPT" <<'EOF'
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 10 },  // Ramp up to 10 users
        { duration: '1m', target: 10 },   // Stay at 10 users
        { duration: '30s', target: 0 },   // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    },
};

export default function () {
    // Test resolver health endpoint
    let resolverHealth = http.get('http://127.0.0.1:8080/healthz');
    check(resolverHealth, {
        'resolver health status is 200': (r) => r.status === 200,
    });

    // Test registrar health endpoint
    let registrarHealth = http.get('http://127.0.0.1:8081/healthz');
    check(registrarHealth, {
        'registrar health status is 200': (r) => r.status === 200,
    });

    // Test DID resolution
    let didResolve = http.get('http://127.0.0.1:8080/resolve?did=did:acc:alice');
    check(didResolve, {
        'DID resolution status is 200 or 404': (r) => r.status === 200 || r.status === 404,
    });

    sleep(1);
}
EOF
fi

# Run k6 test
EXIT_CODE=0

if command -v k6 >/dev/null 2>&1; then
    echo "Running k6 locally..."
    k6 run "$K6_SCRIPT" || EXIT_CODE=$?
else
    echo "Running k6 in Docker..."
    docker run --rm -i --network host grafana/k6 run - < "$K6_SCRIPT" || EXIT_CODE=$?
fi

if [ $EXIT_CODE -eq 0 ]; then
    echo "[OK] Performance tests completed successfully"
else
    echo "[FAIL] Performance tests failed"
fi

exit $EXIT_CODE