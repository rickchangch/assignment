import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 100,
    iterations: 10000,
};

export default function () {
    let res = http.get('http://localhost:8888/isolation?action=UPSERT&level=RR');

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
