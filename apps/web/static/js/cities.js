async function loadData(page = 1) {
  try {
    const data = await API.fetchData("cities", page);
    renderTable(data.results);
    API.renderPagination(data, page, "cities");
  } catch (error) {
    console.error("Erro ao carregar cidades:", error);
  }
}

function renderTable(cities) {
  const tbody = document.querySelector("#dataTable tbody");
  if (!tbody) return;

  tbody.innerHTML =
    cities
      .map(
        (city) => `
        <tr>
            <td>${city.name || "-"}</td>
            <td>${city.state.name} (${city.state.acronym})</td>
        </tr>
    `
      )
      .join("") ||
    '<tr><td colspan="2" class="text-center">Nenhuma cidade encontrada.</td></tr>';
}

document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const page = parseInt(urlParams.get("page")) || 1;
  loadData(page);
});
