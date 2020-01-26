console.log("background");

let sock;

function sendTextToCurrentTab(text) {
  /*  chrome.tabs.query({ active: true, currentWindow: true }, result => {
    const targetTab = result.shift();
    const message = text;
    chrome.tabs.sendMessage(targetTab.id, message);
  });

  const queryInfo = {
    url: "*localhost:3000*"
  };

  chrome.tabs.query(queryInfo, function(result) {
    console.log(result);
    const targetTab = result.shift();
    console.log(targetTab);
    const message = text;
    chrome.tabs.sendMessage(targetTab.id, message);
  }); */

  chrome.tabs.query({ currentWindow: true }, function(tabs) {
    tabs.forEach(function(tab) {
      console.log("Tab ID: ", tab.id);
      chrome.tabs.sendMessage(tab.id, text);
    });
  });

  return;
}

async function initSocket(wsUrl) {
  sock = new WebSocket(wsUrl);

  sock.addEventListener("open", function(e) {
    console.log("socked connected: " + wsUrl);
  });

  sock.addEventListener("close", function(e) {
    console.log("socket closed: " + wsUrl);
  });

  sock.addEventListener("message", function(e) {
    const data = JSON.parse(e.data);
    if (data.type !== "message") return;

    console.log(data.text);
    sendTextToCurrentTab(data.text);
    return;
  });
}

function closeSocket() {
  if (!!sock) {
    sock.close();
  }
}

chrome.runtime.onMessage.addListener(function(message, sender, sendResponse) {
  console.log(message);

  if (message.type == "start") {
    closeSocket();

    const wsUrl = message.data.wsUrl;
    initSocket(wsUrl);
  }

  return true;
});
