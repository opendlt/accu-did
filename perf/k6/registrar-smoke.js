import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
export let errorRate = new Rate('errors');

export let options = {
  duration: __ENV.DURATION || '60s',
  vus: parseInt(__ENV.VUS) || 5, // Lower VUs for registrar to avoid overwhelming
  rps: parseInt(__ENV.RPS) || 10, // Lower RPS for registrar operations
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests must complete within 2s (higher for registrar)
    http_req_failed: ['rate<0.2'], // Error rate must be less than 20% (higher tolerance)
    errors: ['rate<0.2'],
  },
};

const REGISTRAR_URL = __ENV.REGISTRAR_URL || 'http://127.0.0.1:8081';
const API_KEY = __ENV.API_KEY || '';

let testCounter = 0;

export default function () {
  testCounter++;
  const testDID = `did:acc:perf-test-${testCounter}-${__VU}`;

  // Test health check
  const healthResponse = http.get(`${REGISTRAR_URL}/healthz`);

  const healthSuccess = check(healthResponse, {
    'health check status is 200': (r) => r.status === 200,
    'health check response time < 100ms': (r) => r.timings.duration < 100,
  });

  if (!healthSuccess) {
    errorRate.add(1);
    return;
  }

  // Test DID creation
  const createPayload = {
    didDocument: {
      '@context': ['https://www.w3.org/ns/did/v1'],
      id: testDID,
      verificationMethod: [{
        id: `${testDID}#key-1`,
        type: 'Ed25519VerificationKey2020',
        controller: testDID,
        publicKeyMultibase: 'z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK'
      }],
      authentication: [`${testDID}#key-1`]
    }
  };

  const headers = {
    'Content-Type': 'application/json',
  };

  if (API_KEY) {
    headers['Authorization'] = `Bearer ${API_KEY}`;
  }

  const createResponse = http.post(
    `${REGISTRAR_URL}/register`,
    JSON.stringify(createPayload),
    { headers }
  );

  const createSuccess = check(createResponse, {
    'create status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'create response time < 2s': (r) => r.timings.duration < 2000,
    'create response has transaction info': (r) => {
      if (r.status >= 200 && r.status < 300) {
        try {
          const body = JSON.parse(r.body);
          return body.txids || body.accounts;
        } catch (e) {
          return false;
        }
      }
      return true; // Don't fail check for error responses
    },
  });

  if (!createSuccess) {
    errorRate.add(1);
  }

  // Test native deactivation (only if create succeeded)
  if (createResponse.status >= 200 && createResponse.status < 300) {
    const deactivatePayload = {
      did: testDID,
      deactivate: true
    };

    const deactivateResponse = http.post(
      `${REGISTRAR_URL}/native/deactivate`,
      JSON.stringify(deactivatePayload),
      { headers }
    );

    const deactivateSuccess = check(deactivateResponse, {
      'deactivate status is 200': (r) => r.status === 200,
      'deactivate response time < 2s': (r) => r.timings.duration < 2000,
      'deactivate response has tombstone info': (r) => {
        if (r.status === 200) {
          try {
            const body = JSON.parse(r.body);
            return body.didState && body.didState.action === 'deactivate';
          } catch (e) {
            return false;
          }
        }
        return true;
      },
    });

    if (!deactivateSuccess) {
      errorRate.add(1);
    }
  }

  sleep(2); // Longer sleep for registrar operations
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: '  ', enableColors: true }),
    'perf/registrar-smoke-results.json': JSON.stringify(data, null, 2),
  };
}

// Simple text summary since k6 doesn't export textSummary by default
function textSummary(data, options = {}) {
  const indent = options.indent || '';
  const colors = options.enableColors && typeof console !== 'undefined';

  const format = (value, unit = '') => `${value}${unit}`;
  const color = (text, color) => colors ? `\u001b[${color}m${text}\u001b[0m` : text;

  const metrics = data.metrics;

  let summary = `\n${indent}${color('=== Registrar Performance Summary ===', '1;36')}\n`;
  summary += `${indent}Duration: ${format(Math.round(data.state.testRunDurationMs / 1000), 's')}\n`;
  summary += `${indent}VUs: ${data.options.vus}\n`;

  if (metrics.http_reqs) {
    summary += `${indent}Total Requests: ${metrics.http_reqs.values.count}\n`;
    summary += `${indent}Request Rate: ${format(Math.round(metrics.http_reqs.values.rate * 100) / 100, '/s')}\n`;
  }

  if (metrics.http_req_duration) {
    summary += `${indent}Response Time - avg: ${format(Math.round(metrics.http_req_duration.values.avg), 'ms')}, `;
    summary += `p95: ${format(Math.round(metrics.http_req_duration.values['p(95)']), 'ms')}\n`;
  }

  if (metrics.http_req_failed) {
    const failRate = Math.round(metrics.http_req_failed.values.rate * 10000) / 100;
    const statusColor = failRate > 20 ? '1;31' : '1;32';
    summary += `${indent}${color(`Error Rate: ${format(failRate, '%')}`, statusColor)}\n`;
  }

  summary += `${indent}${color('=====================================', '1;36')}\n`;

  return summary;
}