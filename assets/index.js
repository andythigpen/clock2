(() => {
  // js/carousel.ts
  async function advanceCarousel() {
    const el = document.getElementById("carousel");
    if (el === null) {
      return;
    }
    const response = await fetch("/components/weather-current");
    const html = await response.text();
    el.innerHTML = html;
  }
  setInterval(advanceCarousel, 15e3);

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
  setInterval(updateClock, 1e3);
})();
