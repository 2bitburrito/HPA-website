function initializeNavigation() {
  const nav = document.querySelector("#mainNavbar");
  const hamburger = document.querySelector(".hamburger");
  const navList = document.querySelector("nav ul");

  window.addEventListener("scroll", () => {
    if (window.scrollY > 50) {
      nav.classList.add("scrolled");
    } else {
      nav.classList.remove("scrolled");
    }
  });

  hamburger?.addEventListener("click", () => {
    navList.classList.toggle("active");
    hamburger.classList.toggle("active");
  });
}

// Initialize everything when DOM is ready
document.addEventListener("DOMContentLoaded", () => {
  initializeNavigation();
});
