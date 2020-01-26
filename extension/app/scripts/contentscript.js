import React from "react";
import ReactDOM from "react-dom";
import OverlayWindow from "./overlayWindow";

console.log(`content script`);

const app = document.createElement("div");
app.id = "overlay-window-root";
document.body.appendChild(app);

const overlayWindowInstance = ReactDOM.render(<OverlayWindow />, app);

chrome.runtime.onMessage.addListener(function(message, sender, sendResponse) {
  console.log(message);
  overlayWindowInstance.addText(message);
  return true;
});
