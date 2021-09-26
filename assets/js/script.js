'use strict'

import { GetExPoints } from "./exchange.js";
import { CreateWSConnection } from "./ws.js";


const init = () => {
    CreateWSConnection();
    GetExPoints();
}

document.addEventListener("DOMContentLoaded", () => {
    init();
})