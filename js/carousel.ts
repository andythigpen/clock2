const components = [
  // "weather-current",
  "weather-forecast",
];

function* carouselGenerator() {
  let current = 0;
  while (true) {
    const reset = yield components[current];
    current += 1;
    if (current >= components.length || reset) {
      current = 0;
    }
  }
}

const carousel = carouselGenerator();

async function advanceCarousel() {
  const el = document.getElementById("carousel");
  if (el === null) {
    return;
  }
  const component = carousel.next().value;
  const response = await fetch(`/components/${component}`);
  const html = await response.text();
  el.classList.add("opacity-0");
  await new Promise((resolve) => setTimeout(resolve, 700));
  el.innerHTML = html;
  await new Promise((resolve) => setTimeout(resolve, 700));
  el.classList.remove("opacity-0");
}

setInterval(advanceCarousel, 15000);

globalThis.initCarousel = function () {
  if (globalThis.initCarousel.initialized) {
    return;
  }
  globalThis.initCarousel.initialized = true;
  advanceCarousel();
};
