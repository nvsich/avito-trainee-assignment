import http from 'k6/http';
import {check, sleep} from 'k6';
import {randomIntBetween} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const API_BASE = 'http://localhost:8080/api';
const MAX_USER_COUNT = 100;
const MIN_USER_COUNT = 20;
const RPS = 1000;
const TEST_DURATION = '30s';

export const options = {
    scenarios: {
        main: {
            executor: 'constant-arrival-rate',
            rate: RPS,
            timeUnit: '1s',
            duration: TEST_DURATION,
            preAllocatedVUs: MIN_USER_COUNT,
            maxVUs: MAX_USER_COUNT,
        }
    },
    thresholds: {
        http_req_failed: ['rate<0.001'],
        http_req_duration: ['avg<50'],
    },
};

const users = new Map();

const createUser = () => {
    const credentials = {
        username: `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
        password: Math.random().toString(36).substr(2, 12)
    };

    const response = http.post(
        `${API_BASE}/auth`,
        JSON.stringify(credentials),
        {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        }
    );

    check(response, {
        "status is 200": (r) => r.status === 200,
    });

    credentials.token = response.json().token;
    users.set(credentials.username, credentials);
};

const purchase = () => {
    const user = Array.from(users.values())[randomIntBetween(0, users.size - 1)];

    const response = http.get(
        `${API_BASE}/buy/cup`,
        {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${user.token}`
            }
        }
    );

    check(response, {
        "status is 200": (r) => r.status === 200,
    });
};

const transfer = () => {
    const sender = Array.from(users.values())[randomIntBetween(0, users.size - 1)];
    let recipient = Array.from(users.values())[randomIntBetween(0, users.size - 1)];

    while (recipient.username === sender.username) {
        recipient = Array.from(users.values())[randomIntBetween(0, users.size - 1)];
    }

    const response = http.post(
        `${API_BASE}/sendCoin`,
        JSON.stringify({
            toUser: recipient,
            amount: 1
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${sender.token}`
            }
        }
    );

    check(response, {
        "status is 200": (r) => r.status === 200,
    });
};

const getUserInfo = () => {
    const user = Array.from(users.values())[randomIntBetween(0, users.size - 1)];

    const response = http.get(
        `${API_BASE}/info`,
        {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${user.token}`
            }
        }
    );

    check(response, {
        "status is 200": (r) => r.status === 200,
    });
};

export default function () {
    createUser();
    sleep(randomIntBetween(0, 5));
    purchase();
    sleep(randomIntBetween(0, 5));
    transfer();
    sleep(randomIntBetween(0, 5));
    getUserInfo();
    sleep(randomIntBetween(0, 5));
}