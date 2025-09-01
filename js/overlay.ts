let brightness: number;

async function updateOverlay() {
  const params = new URLSearchParams();
  params.append("current", brightness.toString());
  try {
    const response = await fetch(`/api/display/brightness?${params}`);
    const body = await response.json();
    if (typeof body.brightness === "number") {
      brightness = body.brightness;
    }
    const opacity = 100 - brightness;
    const el = document.getElementById("overlay");
    if (el === null) {
      return;
    }
    el.style.opacity = `${opacity}%`;
  } finally {
    // reschedule forever
    setTimeout(updateOverlay, 1000);
  }
}

globalThis.initOverlay = function () {
  if (globalThis.initOverlay.initialized) {
    return;
  }
  brightness = 100;
  globalThis.initOverlay.initialized = true;
  updateOverlay();
};
