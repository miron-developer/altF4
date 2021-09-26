'use strict'


/**
 * recursively create elements and return root
 * @param {string} type - dom element type
 * @param {object} attrs - dom element's attributes
 * @param {[]} children - nested dom elements' configs
 * @param {(dom)=>{}} cb - actions after create dom element
 * @returns 
 */
export const SafelyCreateElement = (type = "div", attrs = {}, children = [], cb) => {
    const DOMElem = document.createElement(type);
    Object.assign(DOMElem, attrs);
    DOMElem.append(...children.map(c => SafelyCreateElement(c.type, c.attrs, c.children, c.cb)))
    cb && cb(DOMElem);
    return DOMElem;
}