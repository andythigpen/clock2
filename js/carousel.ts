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

// setInterval(advanceCarousel, 15000);

globalThis.initCarousel = function () {
  if (globalThis.initCarousel.initialized) {
    return;
  }
  globalThis.initCarousel.initialized = true;
  advanceCarousel();
};
