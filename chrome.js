const puppeteer = require("puppeteer");

async function capture(url) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(url);
    page.on("console", msg => {
        console.log(msg.text());
    });
}

if (process.argv.length == 3) {
    console.log("Ctrl+C to exit...");
    capture(process.argv[2]);
} else {
    console.log(`Usage: node ${process.argv[1]} <url>`);
    process.exit(1);
}