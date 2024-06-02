import http from 'k6/http';
import { sleep, check } from 'k6';
import { Trend } from 'k6/metrics';

export let options = {
    tags: { name: 'demo-api: dice roll' , testid: "demo-api-load-test"},
    duration: 30,
    vu: 10,
};

export default function () {
    const res = http.get('http://demo-api:8080/dice/roll');
    check(res, { 'status was 200': (r) => r.status == 200 });
    // Additional metrics can be collected if needed
    // let trend = new Trend('response_time');
    // trend.add(res.timings.duration);
    sleep(1)
}
