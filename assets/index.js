(() => {
  // js/carousel.ts
  async function advanceCarousel() {
    const el = document.getElementById("carousel");
    if (el === null) {
      return;
    }
    const response = await fetch("/carousel");
    const html = await response.text();
    el.classList.add("opacity-0");
    await new Promise((resolve) => setTimeout(resolve, 700));
    el.innerHTML = html;
    await new Promise((resolve) => setTimeout(resolve, 700));
    el.classList.remove("opacity-0");
  }
  globalThis.initCarousel = function() {
    if (globalThis.initCarousel.initialized) {
      return;
    }
    globalThis.initCarousel.initialized = true;
    advanceCarousel();
    setInterval(advanceCarousel, 15e3);
  };

  // js/clock.ts
  var locale = "en-US";
  function shortMonth(date) {
    return new Intl.DateTimeFormat(locale, { month: "short" }).format(date);
  }
  function shortDay(date) {
    return new Intl.DateTimeFormat(locale, { weekday: "short" }).format(date);
  }
  function formattedDay(date) {
    return new Intl.DateTimeFormat(locale, { day: "2-digit" }).format(date);
  }
  function formattedHour(date) {
    const hour = date.getHours() % 12;
    return hour.toString().padStart(2, "0");
  }
  function formattedMinute(date) {
    const min = date.getMinutes();
    return min.toString().padStart(2, "0");
  }
  function updateClock() {
    const elDate = document.getElementById("date");
    const elTime = document.getElementById("time");
    if (elDate === null || elTime === null) {
      return;
    }
    const now = /* @__PURE__ */ new Date();
    const weekday = shortDay(now);
    const month = shortMonth(now);
    const day = formattedDay(now);
    const hour = formattedHour(now);
    const min = formattedMinute(now);
    elDate.innerHTML = `${weekday} ${month} ${day}`;
    elTime.innerHTML = `${hour}:${min}`;
  }
  globalThis.initClock = function() {
    if (globalThis.initClock.initialized) {
      return;
    }
    globalThis.initClock.initialized = true;
    updateClock();
    setInterval(updateClock, 1e3);
  };

  // js/overlay.ts
  var brightness;
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
      setTimeout(updateOverlay, 1e3);
    }
  }
  globalThis.initOverlay = function() {
    if (globalThis.initOverlay.initialized) {
      return;
    }
    brightness = 100;
    globalThis.initOverlay.initialized = true;
    updateOverlay();
  };
})();
