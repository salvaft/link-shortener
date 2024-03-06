interface Link {
  url: string;
  href: string;
}
const copyBtnTemplate = document.querySelector("#copy-btn") as HTMLTemplateElement;

const createCopyBtn = (link: Link) => {
  const template = copyBtnTemplate.content.cloneNode(true) as HTMLTemplateElement;
  const btn = template.querySelector("button") as HTMLButtonElement;
  btn.addEventListener("click", () => {
    navigator.clipboard.writeText(link.url).then(() => {
      const defaultSpan = btn.querySelector("#default-icon") as HTMLSpanElement;
      console.log(defaultSpan);
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
};

function createLink(link: Link) {
  const a = document.createElement("a");
  const div = document.createElement("div");
  const dt = document.createElement("dt");
  const dd = document.createElement("dd");
  a.classList.add("flex", "flex-col");
  div.classList.add("flex", "justify-between", "items-center", "pb-3");
  dt.classList.add("mb-1", "text-gray-500", "md:text-lg", "dark:text-gray-400");
  dd.classList.add("text-lg", "font-semibold");
  a.href = link.href;
  dd.textContent = link.url;
  dt.textContent = link.href;
  a.appendChild(dd);
  a.appendChild(dt);
  div.appendChild(a);
  div.appendChild(createCopyBtn(link));
  return div;
}
const list = document.querySelector("#recent") as HTMLDListElement;
const existingLinks: Link[] = JSON.parse(localStorage.getItem("links") || "[]");
existingLinks.forEach((link) => {
  list.appendChild(createLink(link));
});

const form = document.querySelector("form") as HTMLFormElement;
form.addEventListener("submit", async (e) => {
  e.preventDefault();
  if (
    existingLinks.some((link) => {
      return link.href === form.url.value;
    })
  ) {
    form.reset();
    // TODO: Show toaster with existing link
    return;
  }
  const data = new FormData(form);
  //TODO: Handle errors and error responses
  const response = await fetch("/", {
    method: "POST",
    body: data,
    headers: {
      "X-From-Js": "true",
    },
  });
  const link = await response.json();
  existingLinks.push(link);
  localStorage.setItem("links", JSON.stringify(existingLinks));

  list.appendChild(createLink(link));
  form.reset();
});
