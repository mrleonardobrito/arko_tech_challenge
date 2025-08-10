async function loadData(page = 1) {
  try {
    const data = await API.fetchData("states", page);
    renderTable(data.results);
    API.renderPagination(data, page, "states");
  } catch (error) {
    console.error("Erro ao carregar estados:", error);
  }
}

function renderTable(states) {
  const tbody = document.querySelector("#dataTable tbody");
  if (!tbody) return;

  tbody.innerHTML =
    states
      .map(
        (state) => `
        <tr>
            <td>${state.name || "-"}</td>
            <td>${state.acronym || "-"}</td>
        </tr>
    `
      )
      .join("") ||
    '<tr><td colspan="2" class="text-center">Nenhum estado encontrado.</td></tr>';
}

document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const page = parseInt(urlParams.get("page")) || 1;
  loadData(page);
});
