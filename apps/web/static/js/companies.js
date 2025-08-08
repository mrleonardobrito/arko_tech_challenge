async function loadData(page = 1) {
  try {
    const data = await API.fetchData("companies", page);
    renderTable(data.results);
    API.renderPagination(data, page, "companies");
  } catch (error) {
    console.error("Erro ao carregar empresas:", error);
  }
}

function renderTable(companies) {
  const tbody = document.querySelector("#dataTable tbody");
  if (!tbody) return;

  tbody.innerHTML =
    companies
      .map(
        (company) => `
        <tr>
            <td>${company.cnpj}</td>
            <td>${company.social_name}</td>
            <td>${company.district.name}</td>
            <td>${company.city.name}</td>
            <td>${company.state.name} (${company.state.acronym})</td>
        </tr>
    `
      )
      .join("") ||
    '<tr><td colspan="5" class="text-center">Nenhuma empresa encontrada.</td></tr>';
}

document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const page = parseInt(urlParams.get("page")) || 1;
  loadData(page);
});
