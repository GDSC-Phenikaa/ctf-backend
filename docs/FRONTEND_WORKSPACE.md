# Frontend Integration Guide for Pwnbox / Workspaces

This guide explains how to integrate the dynamic workspace feature into your frontend (React, Next.js, Vue, or plain HTML).

## Overview

The backend is completely self-contained. When you start a workspace, the backend spins up a Docker container and natively proxies the internal Desktop web server mapping. You don't need to configure WebSockets on the frontend heavily. The easiest way is to simply embed the desktop using an `<iframe>`.

---

## The Workflow

There are three main steps your frontend needs to manage:
1. **Starting the Workspace.**
2. **Displaying the VNC Player.**
3. **Stopping the Workspace.**

### 1. Starting the Workspace

When the user clicks "Start Pwnbox", make an authenticated `POST` request to the backend:

```javascript
async function startWorkspace() {
  const response = await fetch("https://api.yourdomain.com/workspace/start", {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${userToken}`, // Assumes standard auth setup
    }
  });

  const data = await response.json();
  if (response.ok) {
    // data.status === "running"
    // Now you can render the iframe!
    setWorkspaceActive(true);
  }
}
```

### 2. Displaying the VNC Player (The Easy Way)

Because the base image (`kasmweb/ubuntu-bionic-desktop`) comes with its own HTML client built-in, you **do not** need to install a library like `@novnc/novnc` directly in your package.json! 

You can simply mount an `<iframe>` pointing at the backend's native proxy endpoint. The backend will seamlessly tunnel the HTML assets and the WebSocket streams.

```jsx
// React Example

function PwnboxView() {
  // You must include the Authentication cookies if you rely on cookie auth,
  // or pass a token parameter if you use query param auth depending on your setup.
  
  // Note: The trailing slash is important so the browser requests relative assets correctly
  const proxyUrl = "https://api.yourdomain.com/workspace/proxy/";

  return (
    <div className="w-full h-[800px] border border-gray-700 rounded-lg overflow-hidden">
        <iframe 
          src={proxyUrl}
          className="w-full h-full border-0"
          title="User Pwnbox"
          allow="clipboard-read; clipboard-write"
        />
    </div>
  );
}
```

> **Important Note about Auth:** The `<iframe>` will issue a `GET` request to your backend. Browsers do not automatically attach `Authorization: Bearer xyz` headers to `<iframe>` requests. 
> - If your backend relies exclusively on **Cookies**, this will work out-of-the-box!
> - If your backend relies on **Bearer Headers**, you should modify the backend `middlewares.AuthMiddleware` to also check for a token in the URL query string (e.g. `?token=YOUR_TOKEN`), and have the `iframe` URL be `/workspace/proxy/?token=XYZ`.

### 3. Checking Workspace Status

If the user refreshes the page, you can check if they already have a workspace running:

```javascript
async function checkWorkspace() {
  const response = await fetch("https://api.yourdomain.com/workspace/status", {
      headers: { "Authorization": `Bearer ${userToken}` }
  });
  
  if (response.status === 200) {
      // The user has a workspace already! Auto-mount the UI.
      setWorkspaceActive(true);
  } else {
      setWorkspaceActive(false);
  }
}
```

### 4. Stopping the Workspace

When the user clicks "Stop Workspace" or "Destroy Pwnbox", call the stop endpoint.

```javascript
async function stopWorkspace() {
  await fetch("https://api.yourdomain.com/workspace/stop", {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${userToken}`
    }
  });

  // Unmount the iframe
  setWorkspaceActive(false);
}
```

---

## Advanced: Custom noVNC Client

If you *don't* want to use the iframe and you want to build a completely custom UI controls (custom full screen buttons, custom clipboard toolings, etc):

1. Install the library: `npm install @novnc/novnc`
2. Point the Web Socket driver to the exact proxy route.

```javascript
import RFB from '@novnc/novnc/core/rfb';

// Mount the canvas locally
const canvasContainer = document.getElementById("vnc-canvas");

// Connect RFB exactly to the proxied websockify endpoint.
// For Kasmweb images, the websocket usually runs on the /websockify path natively.
// Using "wss://" instead of "https://" forces the WS upgrade immediately
const wsUrl = "wss://api.yourdomain.com/workspace/proxy/websockify";

const rfb = new RFB(canvasContainer, wsUrl, {
    credentials: { password: "password" } // This connects to VNC_PW set in docker.go
});
```
