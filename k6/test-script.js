import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
const cacheHits = new Counter('cache_hits');
const cacheMisses = new Counter('cache_misses');

export let options = {
    scenarios: {
        defaultFunction: {
            exec: 'defaultFunction',
            executor: 'ramping-vus',
            startTime: '0s',
            stages: [
                { duration: '30s', target: 50 },
                { duration: '1m', target: 100 },
                { duration: '30s', target: 0 },
            ],
        },
        cacheHit: {
            exec: 'cacheHit',
            executor: 'ramping-vus',
            startTime: '2m5s',
            stages: [
                { duration: '30s', target: 50 },
                { duration: '1m', target: 100 },
                { duration: '30s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<500'],
        http_req_failed: ['rate<0.01'],
        cache_hits: ['count>0'],
        cache_misses: ['count>0'],
    },
};

export function cacheHit() {
    const userId = Math.floor(Math.random() * 20) + 1;
    try {
    const res = http.get(
        `http://localhost:8080/users/${userId}/recommendations?limit=10`
    );
    const responseBody = JSON.parse(res.body);
    if (responseBody.metadata.cache_hit) {
        cacheHits.add(1);
    } else {
        cacheMisses.add(1);
    }
    check(res, {
        '[cache hit] status is 200': (r) => r.status === 200 && JSON.parse(r.body).metadata.cache_hit === true,});
} catch (error) {
    cacheMisses.add(1);
}
    // check(res, {
    //     'status is 200': (r) => r.status === 200,
    //     'has recommendations': (r) => JSON.parse(r.body).metadata.cache_hit === true,
    // });
    sleep(0.1);
}
export function defaultFunction() {
    const userId = Math.floor(Math.random() * 20) + 1;
    const res = http.get(
        `http://localhost:8080/users/${userId}/recommendations?limit=10`
    );
    check(res, {
        'status is 200': (r) => r.status === 200,
        'has recommendations': (r) => JSON.parse(r.body).recommendations.length > 0,
    });
    sleep(0.1);
}
