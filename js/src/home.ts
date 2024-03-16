interface Link {
  url: string;
  href: string;
  b64_code: string;
}

init();
async function init() {
  const existingLinks: Link[] = JSON.parse(localStorage.getItem("links") || "[]");
  const list = document.querySelector("#recent") as HTMLDListElement;
  hydrateForm(existingLinks);
  existingLinks.forEach((link) => {
    list.appendChild(createLink(link));
  });
}

function retrieveLinkTemplate() {
  const template = document.querySelector("#link-template") as HTMLTemplateElement;
  const el = template.content.cloneNode(true) as HTMLDivElement;
  el.id = "";
  return el;
}

function hydrateCopyBtn(template: HTMLDivElement, link: Link) {
  const btn = template.querySelector("button") as HTMLButtonElement;
  btn.addEventListener("click", () => {
    navigator.clipboard.writeText(window.location + link.b64_code).then(() => {
      const defaultSpan = btn.querySelector("#default-icon") as HTMLSpanElement;
      defaultSpan.classList.remove("inline-flex");
      defaultSpan.classList.add("hidden");
      const copiedSpan = btn.querySelector("#success-icon") as HTMLSpanElement;
      copiedSpan.classList.remove("hidden");
      copiedSpan.classList.add("inline-flex");
      setTimeout(() => {
        defaultSpan.classList.remove("hidden");
        defaultSpan.classList.add("inline-flex");
        copiedSpan.classList.remove("inline-flex");
        copiedSpan.classList.add("hidden");
      }, 1500);
    });
  });
  return btn;
}
function createLink(link: Link) {
  const template = retrieveLinkTemplate();
  hydrateCopyBtn(template, link);
  template.querySelector("a")!.setAttribute("href", link.href);
  template.querySelector("dd")!.textContent = link.url;
  template.querySelector("dt")!.textContent = link.href;
  return template;
}
function hydrateForm(existingLinks: Link[] = []) {
  const form = document.querySelector("form") as HTMLFormElement;
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = new FormData(form);
    const href = data.get("href") as string;
    const existingLink = existingLinks.find((link) => link.href === href);
    if (existingLink) {
      e.preventDefault();
      form.reset();
      tada(existingLink);
      return;
    }

    try {
      const response = await fetch("/", {
        method: "POST",
        body: data,
        headers: {
          "X-From-Js": "true",
        },
      });
      if (!response.ok) {
        throw new Error("Request error");
      }
      const { link: link_data, signed_token } = await response.json();
      existingLinks = [link_data, ...existingLinks];
      form.querySelector("input[name=csrf-token]")!.setAttribute("value", signed_token);
      localStorage.setItem("links", JSON.stringify(existingLinks));
      const list = document.querySelector("#recent") as HTMLDListElement;
      list.prepend(createLink(link_data));
      form.reset();
    } catch (e) {
      console.log(e);
      console.log("Fetch error");
    }
  });
}

function tada(link: Link) {
  const el = document.querySelector(`[href="${link.href}"]`) as HTMLAnchorElement;
  const parent = el.parentElement as HTMLDivElement;
  parent.classList.add("animate-tada");
  setTimeout(() => {
    parent.classList.remove("animate-tada");
  }, 1000);
}
