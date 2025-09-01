const steps = {
  dawndusk1: {
    "--tw-gradient-from": "#734c67",
    "--tw-gradient-via": "#313862",
    "--tw-gradient-via-position": "10%",
    "--tw-gradient-to": "#011b32",
    "--tw-gradient-to-position": "30%",
  },
  dawndusk2: {
    "--tw-gradient-from": "#734c67",
    "--tw-gradient-via": "#313862",
    "--tw-gradient-via-position": "30%",
    "--tw-gradient-to": "#011b32",
    "--tw-gradient-to-position": "60%",
  },
  dawndusk3: {
    "--tw-gradient-from": "#734c67",
    "--tw-gradient-via": "#313862",
    "--tw-gradient-via-position": "60%",
    "--tw-gradient-to": "#011b32",
    "--tw-gradient-to-position": "100%",
  },
  day: {
    "--tw-gradient-from": "#4c6b73",
    "--tw-gradient-via": "#313862",
    "--tw-gradient-via-position": "60%",
    "--tw-gradient-to": "#011b32",
    "--tw-gradient-to-position": "100%",
  },
  night: {
    "--tw-gradient-from": "#150d0d",
    "--tw-gradient-via": "#000000",
    "--tw-gradient-via-position": "60%",
    "--tw-gradient-to": "#000000",
    "--tw-gradient-to-position": "100%",
  },
};

async function updateHomeBackground() {
  try {
    const el = document.getElementById("home");
    if (el === null) {
      return;
    }

    const response = await fetch("/api/sun");
    const body = await response.json();
    if (body.nextRising === undefined || body.nextSetting === undefined) {
      return;
    }
    const nextRising = new Date(body.nextRising).getTime();
    const nextSetting = new Date(body.nextSetting).getTime();
    const now = new Date().getTime();
    const hour = new Date().getHours();
    const oneHourAfterRising = nextRising + 3600000;
    const twoHoursAfterRising = nextRising + 7200000;
    const threeHoursAfterRising = nextRising + 10800000;
    const oneHourBeforeSetting = nextSetting - 3600000;
    const twoHoursBeforeSetting = nextSetting - 7200000;
    const threeHoursBeforeSetting = nextSetting - 10800000;
    const twelveHoursFromNow = now + 43200000;
    let props: Record<string, string>;
    if (
      (now >= nextRising && now < oneHourAfterRising) ||
      (now <= nextSetting && now > oneHourBeforeSetting)
    ) {
      props = steps.dawndusk1;
    } else if (
      (now >= oneHourAfterRising && now < twoHoursAfterRising) ||
      (now <= oneHourBeforeSetting && now > twoHoursBeforeSetting)
    ) {
      props = steps.dawndusk2;
    } else if (
      (now >= twoHoursAfterRising && now < threeHoursAfterRising) ||
      (now <= twoHoursBeforeSetting && now > threeHoursBeforeSetting)
    ) {
      props = steps.dawndusk3;
    } else if (
      (hour > 12 &&
        (nextSetting >= twelveHoursFromNow || nextSetting <= now)) ||
      (hour < 12 && nextRising > now && nextRising <= twelveHoursFromNow)
    ) {
      props = steps.night;
    } else {
      props = steps.day;
    }
    for (const [prop, val] of Object.entries(props)) {
      el.style.setProperty(prop, val);
    }
  } finally {
    // reschedule forever
    setTimeout(updateHomeBackground, 60000);
  }
}

globalThis.initHome = function () {
  if (globalThis.initHome.initialized) {
    return;
  }
  globalThis.initHome.initialized = true;
  updateHomeBackground();
};
