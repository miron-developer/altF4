import { RenderBookRows, SetUserToken } from "./exchange.js";

let wss = null;

const selectAct = (data) => {
    if (!data) return;

    // App ws msgs
    if (data.msgType === 1) return SetUserToken(data.body);
    if (data.msgType === 11) return RenderBookRows(data.body.b, data.body.a);
}

export const CloseWSConnection = () => wss && wss.close();

export const CreateWSConnection = () => {
    wss = new WebSocket('wss://' + window.location.host + '/ws/');

    wss.onopen = () => console.log("websocket connected");
    wss.onclose = (e) => {
        wss = null;
        console.log('closed', e);
    };

    wss.onerror = err => {
        wss = null;
        console.log('err', err);
    };

    wss.onmessage = e => {
        const data = JSON.parse(e.data);
        selectAct(data);
    }
}

export const SendWSMessage = (msgType = 1, receiver = 0, body) => {
    if (wss === null) return { 'err': 'do not sended' };
    wss.send(JSON.stringify({ "msgType": msgType, "receiver": "".concat(receiver), "body": body }));
}