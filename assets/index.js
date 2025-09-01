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
    if (hour == 0) {
      return "12";
    }
    return hour.toString().padStart(2, "0");
  }
  function formattedMinute(date) {
    const min = date.getMinutes();
    return min.toString().padStart(2, "0");
  }
  var prevDate;
  var prevTime;
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
    const date = `${weekday} ${month} ${day}`;
    const time = `${hour}:${min}`;
    if (date !== prevDate) {
      elDate.innerHTML = date;
    }
    if (time !== prevTime) {
      elTime.innerHTML = time;
    }
    prevDate = date;
    prevTime = time;
  }
  globalThis.initClock = function() {
    if (globalThis.initClock.initialized) {
      return;
    }
    globalThis.initClock.initialized = true;
    updateClock();
    setInterval(updateClock, 1e3);
  };

  // js/home.ts
  var steps = {
    dawndusk1: {
      "--tw-gradient-from": "#734c67",
      "--tw-gradient-via": "#313862",
      "--tw-gradient-via-position": "10%",
      "--tw-gradient-to": "#011b32",
      "--tw-gradient-to-position": "30%"
    },
    dawndusk2: {
      "--tw-gradient-from": "#734c67",
      "--tw-gradient-via": "#313862",
      "--tw-gradient-via-position": "30%",
      "--tw-gradient-to": "#011b32",
      "--tw-gradient-to-position": "60%"
    },
    dawndusk3: {
      "--tw-gradient-from": "#734c67",
      "--tw-gradient-via": "#313862",
      "--tw-gradient-via-position": "60%",
      "--tw-gradient-to": "#011b32",
      "--tw-gradient-to-position": "100%"
    },
    day: {
      "--tw-gradient-from": "#4c6b73",
      "--tw-gradient-via": "#313862",
      "--tw-gradient-via-position": "60%",
      "--tw-gradient-to": "#011b32",
      "--tw-gradient-to-position": "100%"
    },
    night: {
      "--tw-gradient-from": "#150d0d",
      "--tw-gradient-via": "#000000",
      "--tw-gradient-via-position": "60%",
      "--tw-gradient-to": "#000000",
      "--tw-gradient-to-position": "100%"
    }
  };
  async function updateHomeBackground() {
    try {
      const el = document.getElementById("home");
      if (el === null) {
        return;
      }
      const response = await fetch("/api/sun");
      const body = await response.json();
      if (body.nextRising === void 0 || body.nextSetting === void 0) {
        return;
      }
      const nextRising = new Date(body.nextRising).getTime();
      const nextSetting = new Date(body.nextSetting).getTime();
      const now = (/* @__PURE__ */ new Date()).getTime();
      const hour = (/* @__PURE__ */ new Date()).getHours();
      const oneHourAfterRising = nextRising + 36e5;
      const twoHoursAfterRising = nextRising + 72e5;
      const threeHoursAfterRising = nextRising + 108e5;
      const oneHourBeforeSetting = nextSetting - 36e5;
      const twoHoursBeforeSetting = nextSetting - 72e5;
      const threeHoursBeforeSetting = nextSetting - 108e5;
      const twelveHoursFromNow = now + 432e5;
      let props;
      if (now >= nextRising && now < oneHourAfterRising || now <= nextSetting && now > oneHourBeforeSetting) {
        props = steps.dawndusk1;
      } else if (now >= oneHourAfterRising && now < twoHoursAfterRising || now <= oneHourBeforeSetting && now > twoHoursBeforeSetting) {
        props = steps.dawndusk2;
      } else if (now >= twoHoursAfterRising && now < threeHoursAfterRising || now <= twoHoursBeforeSetting && now > threeHoursBeforeSetting) {
        props = steps.dawndusk3;
      } else if (hour > 12 && (nextSetting >= twelveHoursFromNow || nextSetting <= now) || hour < 12 && nextRising > now && nextRising <= twelveHoursFromNow) {
        props = steps.night;
      } else {
        props = steps.day;
      }
      for (const [prop, val] of Object.entries(props)) {
        el.style.setProperty(prop, val);
      }
    } finally {
      setTimeout(updateHomeBackground, 6e4);
    }
  }
  globalThis.initHome = function() {
    if (globalThis.initHome.initialized) {
      return;
    }
    globalThis.initHome.initialized = true;
    updateHomeBackground();
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
