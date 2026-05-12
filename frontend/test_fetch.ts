import { $fetch } from 'ofetch';

async function run() {
    try {
        const res = await $fetch('http://localhost:3000/api/v1/auth/session', {
            headers: {
                'Cookie': 'pg_session=123e4567-e89b-12d3-a456-426614174000'
            }
        });
        console.log(res);
    } catch (e) {
        console.error(e);
    }
}
run();
