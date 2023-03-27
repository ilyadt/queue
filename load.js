import http from 'k6/http';
import { check } from 'k6';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        queueAdd: {
            executor: 'constant-vus',
            exec: 'queueAdd',
            vus: 70,
            duration: '30s',
        },
        consume: {
            executor: 'constant-vus',
            exec: 'consume',
            vus: 60,
            duration: '30s',
        },
    },
};

export function queueAdd() {
    const put = http.put('http://127.0.0.1:2802/color?v=white');
    check(put, { 'status was 200': (r) => {
            if (r.status !== 200) {
                console.log(r)
            }
            return r.status === 200
        }
    });
}

export function consume() {
    const getwait = http.get('http://127.0.0.1:2802/color?timeout=1');
    check(getwait, { 'status was 200': (r) => {
            if (r.status !== 200) {
                console.log(r)
            }
            return r.status === 200
    }});
}
