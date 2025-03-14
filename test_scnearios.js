import http from 'k6/http';
import { sleep } from 'k6';
import { check } from 'k6';

export let options = {
  scenarios: {
    load_test: {
      executor: 'constant-arrival-rate',
      rate: 50, 
      duration: '2m', 
      preAllocatedVUs: 100, 
      maxVUs: 200, 
    },
    stress_test: {
      executor: 'ramping-arrival-rate',
      startRate: 10, 
      timeUnit: '1s',
      stages: [
        { target: 100, duration: '2m' }, 
        { target: 200, duration: '3m' }, 
        { target: 0, duration: '1m' }, 
      ],
      preAllocatedVUs: 200,
      maxVUs: 400,
    },
    spike_test: {
      executor: 'constant-arrival-rate',
      rate: 500, 
      duration: '30s', 
      preAllocatedVUs: 500,
      maxVUs: 1000,
    },
    soak_test: {
      executor: 'constant-arrival-rate',
      rate: 20, 
      duration: '10m', 
      preAllocatedVUs: 50,
      maxVUs: 100,
    },
  },
};

const API_URL = 'http://localhost:3002/api/version'; 

export default function () {
  const res = http.get(API_URL);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1); 
}
