const puppeteer = require('puppeteer');

(async () => {
    const browser = await puppeteer.launch({
        defaultViewport: { width: 1200, height: 900 },
        headless: false
    });
    const page = await browser.newPage();
    await page.goto('https://www.investing.com/dividends-calendar/');

    await page.click('#filterStateAnchor');
    await page.type('#searchText_dividends', 'Royal Dutch Shell');
    const el = await page.$$('tr');
    console.log(el);
    await browser.close();
})();
