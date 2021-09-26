'use strict'

import { GetDataByCrieteries, POSTRequestWithParams } from "./api.js";
import { SafelyCreateElement } from "./utils.js";

const exchangeInfo = {
    'token': "",
    'point': {},
    "symbol": {},
}

export const SetUserToken = token => exchangeInfo.token = token;

const clearCurrencie = () => {
    // null active symbol
    exchangeInfo.symbol = {};

    // null currencie order book
    document.querySelector('.currencie-book').innerHTML = "";
}

const unsubscribeFromActiveCurrencie = async() => {
    if (Object.values(exchangeInfo.point).length === 0 || Object.values(exchangeInfo.symbol).length === 0) return;

    const resp = await POSTRequestWithParams("/api/exchange-unsubscribe", { "point": exchangeInfo.point.id, "symbol": exchangeInfo.symbol.symbol, "token": exchangeInfo.token })
    if (resp.code !== 200) {
        clearCurrencie();
        return document.querySelector('.exchange-error').textContent = "The symbol is not available. Try later, or change exchange symbol";
    }
}

export const RenderBookRows = (bids, asks) => {
    if (!bids || !asks) return document.querySelector('.exchange-error').textContent = "Something wrong. Please try later";

    // clear content
    document.querySelector('.currencie-book .bids_body').innerHTML = "";
    document.querySelector('.currencie-book .asks_body').innerHTML = "";

    // generate bids & asks
    const bidsMax = bids.reduce((acc, v) => acc += parseFloat(v[1]), 0.0);
    const asksMax = asks.reduce((acc, v) => acc += parseFloat(v[1]), 0.0);
    const maxTotal = bidsMax > asksMax ? bidsMax : asksMax;


    let bidsTotal = 0.0;
    document.querySelector('.currencie-book .bids_body').append(...bids.map(
        b => {
            const b0 = parseFloat(b[0])
            const b1 = parseFloat(b[1])
            bidsTotal += b1
            return SafelyCreateElement(
                "div", { className: "bids-row" }, [{
                    type: "span",
                    attrs: { textContent: 1 }
                }, {
                    type: "span",
                    attrs: { textContent: b1.toFixed(8) }
                }, {
                    type: "span",
                    attrs: { textContent: bidsTotal.toFixed(8) },
                }, {
                    type: "span",
                    attrs: { textContent: b0.toFixed(8) }
                }, {
                    type: "span",
                    attrs: { className: "bids-row-bg" },
                    cb: dom => dom.style.width = bidsTotal / maxTotal * 100 + "%"
                }]
            )
        }
    ));

    let asksTotal = 0.0;
    document.querySelector('.currencie-book .asks_body').append(...asks.map(
        a => {
            const a0 = parseFloat(a[0])
            const a1 = parseFloat(a[1])
            asksTotal += a1
            return SafelyCreateElement(
                "div", { className: "asks-row" }, [{
                    type: "span",
                    attrs: { textContent: 1 }
                }, {
                    type: "span",
                    attrs: { textContent: a1.toFixed(8) }
                }, {
                    type: "span",
                    attrs: { textContent: asksTotal.toFixed(8) },
                }, {
                    type: "span",
                    attrs: { textContent: a0.toFixed(8) }
                }, {
                    type: "span",
                    attrs: { className: "asks-row-bg" },
                    cb: dom => dom.style.width = asksTotal / maxTotal * 100 + "%"
                }]
            )
        }
    ))
}

const subscribeToCurrencie = async(symbol) => {
    if (!symbol) return document.querySelector('.exchange-error').textContent = "Please, select exchange currencie";

    // unsubscribe active symbol
    unsubscribeFromActiveCurrencie();

    // clear
    clearCurrencie();

    // change active currencie
    document.querySelectorAll(".exchange-currencie").forEach(p => p.classList.contains("exchange-" + symbol.symbol) ? p.classList.add("active") : p.classList.remove("active"));
    exchangeInfo.symbol = symbol;

    // generate currencie info
    document.querySelector('.currencie-book').append(SafelyCreateElement(
        "div", { className: "book-header" }, [{
            type: "span",
            attrs: { className: "book-header_symbol", textContent: "Order book: " + symbol.symbol }
        }, {
            type: "span",
            attrs: { className: "book-header_baseAsset", textContent: "Base: " + symbol.baseAsset }
        }, {
            type: "span",
            attrs: { className: "book-header_quoteAsset", textContent: "Quote:" + symbol.quoteAsset }
        }]
    ))

    // get first view of order book
    const resp = await GetDataByCrieteries("/api/exchange-subscribe", { "point": exchangeInfo.point.id, "symbol": exchangeInfo.symbol.symbol, "token": exchangeInfo.token })
    if (resp.code !== 200) return document.querySelector('.exchange-error').textContent = "The symbol is not available. Try later, or change exchange symbol";

    // generate headers of each side
    document.querySelector('.currencie-book').append(SafelyCreateElement(
        "div", { className: "book-body" }, [{
            type: "div",
            attrs: { className: "book-body_bids" },
            children: [{
                type: "div",
                attrs: { className: "bids_header" },
                children: [{
                    type: "span",
                    attrs: { textContent: "count" }
                }, {
                    type: "span",
                    attrs: { textContent: "amount" }
                }, {
                    type: "span",
                    attrs: { textContent: "total" }
                }, {
                    type: "span",
                    attrs: { textContent: "price" }
                }]
            }, {
                type: "div",
                attrs: { className: "bids_body" },
            }]
        }, {
            type: "div",
            attrs: { className: "book-body_asks" },
            children: [{
                type: "div",
                attrs: { className: "asks_header" },
                children: [{
                    type: "span",
                    attrs: { textContent: "count" }
                }, {
                    type: "span",
                    attrs: { textContent: "amount" }
                }, {
                    type: "span",
                    attrs: { textContent: "total" }
                }, {
                    type: "span",
                    attrs: { textContent: "price" }
                }]
            }, {
                type: "div",
                attrs: { className: "asks_body" },
            }]
        }],
    ))

    RenderBookRows(resp.data.bids, resp.data.asks);
}

const getExCurrencies = async(point) => {
    if (!point) return document.querySelector('.exchange-error').textContent = "Please, select exchange point";

    // unsubscribe active symbol
    unsubscribeFromActiveCurrencie();

    // clear
    clearCurrencie();
    document.querySelector('.exchange-currencies').innerHTML = "";

    // change active point
    document.querySelectorAll(".exchange-point").forEach(p => p.classList.contains("exchange-" + point.name) ? p.classList.add("active") : p.classList.remove("active"));
    exchangeInfo.point = point;

    const resp = await GetDataByCrieteries("/api/exchange-currencies", { "point": exchangeInfo.point.id, "token": exchangeInfo.token })
    if (resp.code !== 200) return document.querySelector('.exchange-error').textContent = "The point-service is not available. Try later, or change exchange point";

    // generation currencies
    document.querySelector(".exchange-currencies").append(...resp.data.map(
        c =>
        SafelyCreateElement(
            "div", { className: "exchange-currencie exchange-" + c.symbol, textContent: c.baseAsset + "-" + c.quoteAsset }, [],
            dom => {
                // set data-id attr
                dom.dataset.id = c.symbol;

                // set event onclick to select currencie & get first view of order book
                dom.onclick = () => subscribeToCurrencie(c)
            }
        )
    ))
}

export const GetExPoints = async() => {
    const resp = await GetDataByCrieteries("/api/exchange-points", {});
    if (resp.code !== 200) return document.querySelector('.exchange-error').textContent = "The service is not available. Try later";

    // generation points
    document.querySelector(".exchange-points").append(...resp.data.map(
        p =>
        SafelyCreateElement(
            "div", { className: "exchange-point exchange-" + p.name, textContent: p.name }, [],
            dom => {
                // set data-id attr
                dom.dataset.id = p.id;

                // set event onclick to select point & get currencies
                dom.onclick = () => getExCurrencies(p);
            }
        )
    ))
}