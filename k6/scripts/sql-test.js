import sql from 'k6/x/sql';
import mysql from 'k6/x/sql/driver/mysql';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

// The database connection string
const dbStr = 'app:password@tcp(mysql:3306)/orders_db';

// Custom metrics to track SQL performance
const queryDuration = new Trend('sql_query_duration', true);
const queryErrorRate = new Rate('sql_query_error_rate');

export const options = {
    vus: 10,
    duration: '30s',
    thresholds: {
        'sql_query_duration': ['p(95)<500'], // 95% of queries must complete below 500ms
        'sql_query_error_rate': ['rate<0.01'], // Error rate must be < 1%
    },
};

// Initialize the database connection exactly once per VU
let db;

export function setup() {
    // We don't connect in setup because setup only runs once for the whole test,
    // and we want separate connections for each Virtual User (VU) to simulate real load.
}

export default function () {
    // Open connection on first iteration for each VU
    if (!db) {
        db = sql.open(mysql, dbStr);
    }

    // 1. Benchmark: Simple query returning 1 row (simulating a fast lookup)
    const startSimple = Date.now();
    let resultSimple;
    try {
        resultSimple = db.query('SELECT 1 AS val;');
        queryDuration.add(Date.now() - startSimple, { type: 'simple_select' });
        queryErrorRate.add(0);

        check(resultSimple, {
            'simple query successful': (r) => r.length === 1,
        });
    } catch (err) {
        console.error(`Simple query error: ${err}`);
        queryErrorRate.add(1);
    }

    // 2. Benchmark: Query the application's actual table (if exists) or simulate wait
    // Often in performance tests you might query actual business data.
    const startTable = Date.now();
    try {
        // Find existing orders limits to 10. Since this is a test, the table might be empty.
        const resultTable = db.query('SELECT * FROM orders LIMIT 10;');
        queryDuration.add(Date.now() - startTable, { type: 'table_select' });
        queryErrorRate.add(0);

        check(resultTable, {
            'table query successful': (r) => r !== null,
        });
    } catch (err) {
        // If the table doesn't exist, this will error, but we want to know that.
        console.error(`Table query error: ${err}`);
        queryErrorRate.add(1);
    }

    // A small sleep to simulate think time and avoid overloading the DB unrealistically fast
    // depending on the load model.
    sleep(0.5);
}

export function teardown() {
    // Teardown also runs only once.
    // Cleanup of the DB connections is handled by k6 at the end of VU execution.
    // However, explicit closing is good practice if possible,
    // but in xk6-sql current API structure, we typically open per VU and let it close on exit.
}
