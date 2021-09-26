'use strict'


const formDataToString = (data = new FormData()) => {
    let res = "";
    for (let [k, v] of data.entries())
        res += k + "=" + v + "&"
    return res.slice(0, -1)
}

// use fetching by both method
export const Fetching = async(action, data, method = "POST") => {
    if (action === undefined) return { err: "action undefined" };


    const fetchOption = { 'method': method };
    if (method === "GET") action += "?" + encodeURI(formDataToString(data));
    else fetchOption["body"] = data;

    return await fetch(action, fetchOption)
        .then(res => res.json())
        .catch(() => Object.assign({}, { 'code': 500, 'err': "500: internal client error" }));
}

// convert from object to FormData
const prepareDataToFetch = (datas = {}) => {
    const data = new FormData();
    for (let [k, v] of Object.entries(datas)) data.append(k, v);
    return data;
}

// get data by criteries & type
export const GetDataByCrieteries = async(datatype, criteries = {}) => await Fetching(datatype, prepareDataToFetch(criteries), 'GET');

// send post req to host with params
export const POSTRequestWithParams = async(to, params = {}) => await Fetching(to, prepareDataToFetch(params));