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

// Send contact form to Lambda for Ses
async function initializeContactForm() {
  const form = document.getElementById("contact-form");

  async function handleSubmit(e) {
    e.preventDefault();

    const formData = {
      name: document.getElementById("form-name").value,
      email: document.getElementById("form-email").value,
      message: document.getElementById("form-message").value,
    };
    const jsonData = JSON.stringify(formData);
    try {
      const response = await fetch(
        "https://q9ut7p24g0.execute-api.ap-southeast-2.amazonaws.com/dev/contact",
        {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: jsonData,
        }
      );
      const res = await response.json();
      if (res.success) {
        alert("EMAIL SENT SUCCESSFUL");
      } else {
        alert("EMAIL SEND UNSUCCESSFUL");
      }
    } catch (e) {
      console.error(`Error sending form ${e}`);
    }
  }

  form.addEventListener("submit", handleSubmit);
}

// Initialize everything when DOM is ready
document.addEventListener("DOMContentLoaded", () => {
  initializeNavigation();
  initializeContactForm();
});
