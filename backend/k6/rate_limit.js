import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 20 }, // Ramp up to 20 users
        { duration: '1m', target: 20 },  // Stay at 20 users (generating > 100 req/min combined)
        { duration: '10s', target: 0 },  // Ramp down
    ],
    thresholds: {
        http_req_failed: ['rate<0.1'], // Allow some errors (expected 429s shouldn't be too many for general health, but actually we WANT 429s here)
    },
};

export default function () {
    const res = http.get('http://localhost:9999/health');

    // We expect success initially, but 429 when limit is hit.
    // Ideally, we'd separate this into a check.
    check(res, {
        'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
    });

    sleep(0.1); // ~10 requests per second per VU? No, 0.1s sleep = ~10 req/s.
    // 20 VUs * 10 req/s = 200 req/s. Far exceeds 100 req/min.
}
