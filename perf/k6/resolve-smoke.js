import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
export let errorRate = new Rate('errors');

export let options = {
  duration: __ENV.DURATION || '60s',
  vus: parseInt(__ENV.VUS) || 10,
  rps: parseInt(__ENV.RPS) || 100,
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete within 500ms
    http_req_failed: ['rate<0.1'], // Error rate must be less than 10%
    errors: ['rate<0.1'],
  },
};

const RESOLVER_URL = __ENV.RESOLVER_URL || 'http://127.0.0.1:8080';

// Test DIDs
const TEST_DIDS = [
  'did:acc:alice',
  'did:acc:bob',
  'did:acc:company.example',
  'did:acc:test.user',
  'did:acc:beastmode.acme',
];

export default function () {
  // Pick random test DID
  const testDID = TEST_DIDS[Math.floor(Math.random() * TEST_DIDS.length)];
  const url = `${RESOLVER_URL}/resolve?did=${testDID}`;

  const response = http.get(url, {
    headers: {
      'Accept': 'application/did+json',
    },
  });

  // Check response
  const success = check(response, {
    'status is 200 or 404': (r) => r.status === 200 || r.status === 404,
    'response time < 500ms': (r) => r.timings.duration < 500,
    'has content-type header': (r) => r.headers['Content-Type'] !== undefined,
  });

  if (!success) {
    errorRate.add(1);
  }

  // Check for 410 deactivated responses
  if (response.status === 410) {
    check(response, {
      'deactivated response has proper content-type': (r) =>
        r.headers['Content-Type'].includes('application/json'),
      'deactivated response has error field': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.error === 'deactivated';
        } catch (e) {
          return false;
        }
      },
    });
  }

  // Check for valid 200 responses
  if (response.status === 200) {
    check(response, {
      'valid DID resolution result': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.didDocument && body.didDocument.id;
        } catch (e) {
          return false;
        }
      },
    });
  }

  sleep(1);
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: '  ', enableColors: true }),
    'perf/resolve-smoke-results.json': JSON.stringify(data, null, 2),
  };
}

// Simple text summary since k6 doesn't export textSummary by default
function textSummary(data, options = {}) {
  const indent = options.indent || '';
  const colors = options.enableColors && typeof console !== 'undefined';

  const format = (value, unit = '') => `${value}${unit}`;
  const color = (text, color) => colors ? `\u001b[${color}m${text}\u001b[0m` : text;

  const metrics = data.metrics;

  let summary = `\n${indent}${color('=== Performance Summary ===', '1;36')}\n`;
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
    const statusColor = failRate > 10 ? '1;31' : '1;32';
    summary += `${indent}${color(`Error Rate: ${format(failRate, '%')}`, statusColor)}\n`;
  }

  summary += `${indent}${color('================================', '1;36')}\n`;

  return summary;
}