interface Link {
  url: string;
  href: string;
}
init();

async function init() {
  const existingLinks: Link[] = JSON.parse(localStorage.getItem("links") || "[]");
  const csrfToken = document.querySelector('[name="csrf-token"]')!.getAttribute("value");
  for (const link of existingLinks) {
    await htmx
      .ajax("POST", "/link", {
        values: { href: link.href, "csrf-token": csrfToken },
        target: "#recent",
        swap: "beforeend",
      })
      .then(() => {
        console.log("done", link);
      });
  }
}
