async function advanceCarousel() {
  const el = document.getElementById("carousel");
  if (el === null) {
    return;
  }
  const response = await fetch("/components/weather-current");
  const html = await response.text();
  el.innerHTML = html;
}

setInterval(advanceCarousel, 15000);
