'use strict';

module.exports = {
    generateRandomId,
    generateOrderPayload,
    generateUpdatePayload,
};

function generateRandomId(userContext, events, done) {
    userContext.vars.orderId = Math.floor(Math.random() * 100) + 1;
    return done();
}

function generateOrderPayload(userContext, events, done) {
    userContext.vars.productName = `負荷テスト商品_${Date.now()}`;
    userContext.vars.quantity = Math.floor(Math.random() * 10) + 1;
    return done();
}

function generateUpdatePayload(userContext, events, done) {
    userContext.vars.quantity = Math.floor(Math.random() * 20) + 1;
    userContext.vars.note = `負荷テスト更新_${Date.now()}`;
    return done();
}
