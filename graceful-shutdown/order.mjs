#!/usr/bin/env zx

const ticker = new Promise(resolve => {
    let i = 0;
    const interval = setInterval(() => {
        i++;
        process.stdout.write(`\r[client] order sent, (${i}s) waiting for response...`);
    }, 1000);

    resolve({stop: () => clearInterval(interval)});
});

const {stop} = await ticker;

try {
    const response = await fetch("http://localhost:8080/api/v1/orders", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            product_id: randRange(1000, 9999),
            price: randRange(10000, 60000).toString(),
            count: randRange(1, 5)
        })
    });
    const body = await response.json();
    console.log(`\n[client] response received:`);
    console.table(body);
} catch (err) {
    console.log(`\n[client] error while calling API:`, err);
} finally {
    stop()
}

function randRange(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}