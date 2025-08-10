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

  if (!companies || companies.length === 0) {
    tbody.innerHTML =
      '<tr><td colspan="6" class="text-center">Nenhuma empresa encontrada.</td></tr>';
    return;
  }

  console.log("Renderizando empresas:", companies);

  tbody.innerHTML = companies
    .map(
      (company) => `
      <tr>
          <td>${company.cnpj || "-"}</td>
          <td>${company.social_name || "-"}</td>
          <td>${company.juridical_nature || "-"}</td>
          <td>${formatCurrency(company.social_capital) || "-"}</td>
          <td>${formatCompanySize(company.company_size) || "-"}</td>
          <td>${company.federative_entity || "-"}</td>
      </tr>
  `
    )
    .join("");
}

function formatCurrency(value) {
  if (!value) return "-";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(value);
}

function formatCompanySize(size) {
  const sizes = {
    "01": "Não Informado",
    "02": "Micro Empresa",
    "03": "Pequeno Porte",
    "04": "Médio Porte",
    "05": "Grande Porte",
  };
  return sizes[size] || size || "-";
}

document.addEventListener("DOMContentLoaded", () => {
  const urlParams = new URLSearchParams(window.location.search);
  const page = parseInt(urlParams.get("page")) || 1;
  const page_size = parseInt(urlParams.get("page_size")) || 20;
  loadData(page, page_size);
});
