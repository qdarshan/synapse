chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: "saveSelectedText",
    title: "Save Selected Text",
    contexts: ["selection"],
  });
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
  const tabId = tab.id;
  const tabURL = tab.url;

  if (info.menuItemId === "saveSelectedText") {
    const selectedText = info.selectionText;

    chrome.scripting.executeScript(
      {
        target: { tabId },
        files: ["Readability.js"],
      },
      () => {
        chrome.scripting.executeScript({
          target: { tabId },
          func: (text, url) => {
            const docClone = document.cloneNode(true);
            const article = new Readability(docClone).parse();

            fetch("http://localhost:8080/api/notes", {
              method: "POST",
              body: JSON.stringify({
                url: url,
                title: article.title,
                content: text,
              }),
              headers: {
                "Content-type": "application/json; charset=UTF-8",
              },
            })
              .then((response) => response.json())
              .then((jsonResponse) => console.log(jsonResponse))
              .catch((error) => {
                console.error("Fetch error:", error.message);
                alert(`Failed to save article: ${error.message}`);
              });
          },
          args: [selectedText, tabURL],
        });
      },
    );
  }
});
