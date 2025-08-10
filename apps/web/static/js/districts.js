async function loadData(page = 1) {
  try {
    const data = await API.fetchData("districts", page);
    renderTable(data.results);
    API.renderPagination(data, page, "districts");
  } catch (error) {
    console.error("Erro ao carregar distritos:", error);
  }
}

function renderTable(districts) {
  const tbody = document.querySelector("#dataTable tbody");
  if (!tbody) return;

  tbody.innerHTML =
    districts
      .map(
        (district) => `
        <tr>
            <td>${district.name || "-"}</td>
            <td>${district.city.name || "-"}</td>
            <td>${district.city.state.name || "-"} (${
          district.city.state.acronym || "-"
        })</td>
        </tr>
    `
      )
      .join("") ||
    '<tr><td colspan="3" class="text-center">Nenhum distrito encontrado.</td></tr>';
}

document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const page = parseInt(urlParams.get("page")) || 1;
  loadData(page);
});
