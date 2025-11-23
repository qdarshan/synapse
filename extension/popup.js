document.addEventListener("DOMContentLoaded", async () => {
  try {
    const notes = await fetchNotes();
    console.log("Fetched notes:", notes);

    const notesContainer = document.getElementById("notesContainer");

    if (!notes || notes.length === 0) {
      notesContainer.innerHTML = "<p>No notes found.</p>";
      return;
    }

    notesContainer.innerHTML = "";
    notes.data.forEach((note) => {
      const div = document.createElement("div");
      div.className = "note";
      div.textContent = note.title;
      notesContainer.appendChild(div);
    });
  } catch (err) {
    console.error(err);
  }
});

async function fetchNotes() {
  const response = await fetch("http://localhost:8080/api/notes", {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });

  if (!response.ok) {
    throw new Error("Failed to fetch notes");
  }

  return await response.json();
}
