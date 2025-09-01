const locale = "en-US";

function shortMonth(date: Date): string {
  return new Intl.DateTimeFormat(locale, { month: "short" }).format(date);
}

function shortDay(date: Date): string {
  return new Intl.DateTimeFormat(locale, { weekday: "short" }).format(date);
}

function formattedDay(date: Date): string {
  return new Intl.DateTimeFormat(locale, { day: "2-digit" }).format(date);
}

function formattedHour(date: Date): string {
  const hour = date.getHours() % 12;
  if (hour == 0) {
    return "12";
  }
  return hour.toString().padStart(2, "0");
}

function formattedMinute(date: Date): string {
  const min = date.getMinutes();
  return min.toString().padStart(2, "0");
}

let prevDate: string;
let prevTime: string;

function updateClock() {
  const elDate = document.getElementById("date");
  const elTime = document.getElementById("time");
  if (elDate === null || elTime === null) {
    return;
  }
  const now = new Date();
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

globalThis.initClock = function () {
  if (globalThis.initClock.initialized) {
    return;
  }
  globalThis.initClock.initialized = true;
  updateClock();
  setInterval(updateClock, 1000);
};
