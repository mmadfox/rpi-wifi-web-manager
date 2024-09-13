document.addEventListener("DOMContentLoaded", () => {
  main();
});

var loader = {};
var render = {};
const msg = {};

msg.show = function (text) {
  document.querySelector("#msg-text").innerText = text;
  document.querySelector("#msg").classList.remove("hidden");
};

msg.close = function () {
  document.querySelector("#msg").classList.add("hidden");
};

loader.show = function () {
  document.querySelector("#loader").classList.remove("hidden");
};

loader.hide = function () {
  document.querySelector("#loader").classList.add("hidden");
};

loader.noDataShow = function () {
  document.querySelector("#no-data").classList.remove("none");
};

loader.noDataHide = function () {
  document.querySelector("#no-data").classList.add("none");
};

loader.activeConnectionShow = function () {
  document.querySelector(".container").classList.add("container-top");
  document.querySelector("#active-connection").classList.remove("none");
};

loader.activeConnectionHide = function () {
  document.querySelector("#active-connection").classList.add("none");
  document.querySelector(".container").classList.remove("container-top");
};

var api = {};

api.connect = async function (params) {
  try {
    const url = host() + "/api/wifi-conn";
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    const response = await fetch(url, {
      method: "POST",
      body: JSON.stringify(params),
      headers: headers,
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      if ("error" in errorResponse) {
        throw new Error(errorResponse.error);
      } else {
        throw new Error(`Response status: ${response.status}`);
      }
    }
    return await response.json();
  } catch (e) {
    throw e;
  }
};

api.showIfaces = async function () {
  const url = host() + "/api/ifaces";
  try {
    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      if ("error" in errorResponse) {
        throw new Error(errorResponse.error);
      } else {
        throw new Error(`Response status: ${response.status}`);
      }
    }
    return await response.json();
  } catch (e) {
    throw e;
  }
};

api.switchIface = async function (iface) {
  try {
    const url = host() + "/api/ifaces/switch";
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    const response = await fetch(url, {
      method: "POST",
      body: JSON.stringify({ ifname: iface }),
      headers: headers,
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      if ("error" in errorResponse) {
        throw new Error(errorResponse.error);
      } else {
        throw new Error(`Response status: ${response.status}`);
      }
    }
    return await response.json();
  } catch (e) {
    throw e;
  }
};

api.disconnect = async function () {
  try {
    const url = host() + "/api/wifi-close";
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    const response = await fetch(url, {
      method: "POST",
      headers: headers,
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      if ("error" in errorResponse) {
        throw new Error(errorResponse.error);
      } else {
        throw new Error(`Response status: ${response.status}`);
      }
    }
    return await response.json();
  } catch (e) {
    throw e;
  }
};

api.fetchList = async function () {
  const url = host() + "/api/wifi-list";
  try {
    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      const errorResponse = await response.json();
      if ("error" in errorResponse) {
        throw new Error(errorResponse.error);
      } else {
        throw new Error(`Response status: ${response.status}`);
      }
    }
    return await response.json();
  } catch (e) {
    throw e;
  }
};

render.resetWiFiListContainer = function () {
  document.getElementById("wifi-list").innerHTML = "";
};

render.setActiveConnection = function (ssid) {
  document.getElementById("active-ssid").innerText = ssid;
  document.getElementById("btn-disconnect").setAttribute("data-ssid", ssid);
};

render.hideActiveConnection = function (ssid) {
  document.getElementById("active-ssid").innerText = "";
  document.getElementById("btn-disconnect").setAttribute("data-ssid", "");
};

render.showIfaces = function (ifaces) {
  const container = document.getElementById("ifaces");
  container.innerHTML = "";
  const template = document.getElementById("iface-item-template").content;
  ifaces.forEach((iface) => {
    const clone = document.importNode(template, true);
    const item = clone.querySelector(".iface");
    item.setAttribute("data-iface", iface.name);
    if (iface.default) {
      item.classList.add("active");
    }
    if (!iface.active) {
      item.classList.add("disabled");
    }
    item.innerText = iface.name;
    container.appendChild(clone);
  });
};

render.listenEvents = function () {
  document.querySelector("#wifi-list").addEventListener("click", (e) => {
    if (e.target && e.target.classList.contains("btn-connect")) {
      const wifiItem = e.target.closest(".wifi-item");
      const ssid = wifiItem.querySelector(".wifi-ssid").textContent;
      const form = document.querySelector("#wifi-connect-form");
      render.lockWiFiList = true;
      form.classList.remove("hidden");

      const onConnect = () => {
        const sp = document.querySelector("#save-password").value;
        const customEvent = new CustomEvent("connect", {
          detail: {
            ssid: form.querySelector("#ssid").value,
            password: form.querySelector("#password").value,
            savePoint: sp === "on",
          },
        });
        document.dispatchEvent(customEvent);
      };

      const onClose = () => {
        const customEvent = new CustomEvent("closeForm", {
          detail: { ssid: ssid },
        });
        document.dispatchEvent(customEvent);
        form.querySelector(".connect").removeEventListener("click", onConnect);
        form.querySelector(".close").removeEventListener("click", onClose);
        document.removeEventListener("closeConnectionForm", onClose);
        form.classList.add("hidden");
        form.querySelector("#password").value = "";
      };

      document.addEventListener("closeConnectionForm", (e) => {
        onClose();
      });

      form.querySelector("#ssid").value = ssid;
      form.querySelector(".connect").addEventListener("click", onConnect);
      form.querySelector(".close").addEventListener("click", onClose);
      form.classList.remove("none");
    }
  });
};

render.displayWifiList = function (wifiList) {
  const container = document.getElementById("wifi-list");
  const template = document.getElementById("wifi-item-template").content;
  wifiList.forEach((wifi) => {
    const clone = document.importNode(template, true);
    clone.querySelector(".wifi-ssid").textContent = wifi.ssid;
    clone.querySelector(".wifi-freq").textContent = wifi.freq;
    clone.querySelector(".wifi-signal").textContent = `Signal: ${wifi.signal}%`;
    clone.querySelector(".btn-connect").setAttribute("data-ssid", wifi.ssid);
    container.appendChild(clone);
  });
};

async function loadWiFiList() {
  try {
    render.resetWiFiListContainer();
    const wifiList = await api.fetchList();
    if (wifiList.active !== "None") {
      render.setActiveConnection(wifiList.active);
      loader.activeConnectionShow();
    } else {
      loader.activeConnectionHide();
    }
    if (wifiList.list.length > 0) {
      loader.noDataHide();
      render.displayWifiList(wifiList.list);
    } else {
      loader.noDataShow();
    }
  } catch (e) {
    msg.show(e.message);
  }
}

function host() {
  const hostname = window.location.hostname;
  const port = window.location.port;
  const proto = window.location.protocol;
  if (port !== undefined) {
    return proto + "//" + hostname + ":" + port;
  }
  return proto + "//" + hostname;
}

var controller = {
  lockWiFiList: false,
};

controller.updateAll = async function () {
  loader.show();
  // wifis
  await loadWiFiList();
  // ifaces
  await controller.updateIfaces();
  loader.hide();
};

controller.handleWifiDisconnect = function () {
  document
    .querySelector("#btn-disconnect")
    .addEventListener("click", async () => {
      loader.show();
      try {
        await api.disconnect();
        await loadWiFiList();
        await controller.updateIfaces();
      } catch (e) {
        msg.show(e.message);
      } finally {
        loader.hide();
      }
    });
};

controller.handleIfaceSwitch = function () {
  document
    .querySelector("#ifaces-container")
    .addEventListener("click", async (e) => {
      if (e.target && !e.target.classList.contains("iface")) {
        return;
      }
      if (
        e.target.classList.contains("active") ||
        e.target.classList.contains("disabled")
      ) {
        return;
      }
      loader.show();
      const iface = e.target.getAttribute("data-iface");
      try {
        await api.switchIface(iface);
        await controller.updateIfaces();
      } catch (e) {
        msg.show(e.message);
      } finally {
        loader.hide();
      }
    });
};

controller.handleWifiConnect = function () {
  document.addEventListener("closeForm", () => {
    controller.lockWiFiList = false;
  });

  document.addEventListener("connect", async (e) => {
    loader.show();
    try {
      const result = await api.connect(e.detail);
      document.dispatchEvent(
        new CustomEvent("closeConnectionForm", {
          detail: {
            status: result.ok,
          },
        })
      );
      await loadWiFiList();
      await controller.updateIfaces();
    } catch (e) {
      loader.hide();
      msg.show(e.message);
    } finally {
      loader.hide();
    }
  });
};

controller.updateIfaces = async function () {
  const result = await api.showIfaces();
  render.showIfaces(result.ifaces);
};

controller.listenEvents = function () {
  render.listenEvents();
};

async function main() {
  console.log("webnm running...");

  controller.handleIfaceSwitch();
  controller.handleWifiDisconnect();
  controller.handleWifiConnect();
  controller.listenEvents();
  controller.updateAll();

  setInterval(async () => {
    if (!controller.lockWiFiList) {
      controller.updateAll();
    }
  }, 30000);
}
