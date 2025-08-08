class API {
  static async fetchData(endpoint, page = 1) {
    try {
      const response = await fetch(`/api/${endpoint}/?page=${page}`);
      if (!response.ok) {
        throw new Error("Erro ao buscar dados");
      }
      return await response.json();
    } catch (error) {
      console.error("Erro:", error);
      throw error;
    }
  }

  static renderPagination(data, currentPage, entityName) {
    const paginationElement = document.getElementById("pagination");
    if (!paginationElement) return;

    let html = '<ul class="pagination justify-content-center">';

    // Botão anterior
    if (data.previous) {
      html += `
                <li class="page-item">
                    <a class="page-link" href="#" data-page="${
                      currentPage - 1
                    }" aria-label="Anterior">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>`;
    }

    // Números das páginas
    const totalPages = Math.ceil(data.count / data.results.length);
    for (let i = 1; i <= totalPages; i++) {
      html += `
                <li class="page-item ${i === currentPage ? "active" : ""}">
                    <a class="page-link" href="#" data-page="${i}">${i}</a>
                </li>`;
    }

    // Botão próximo
    if (data.next) {
      html += `
                <li class="page-item">
                    <a class="page-link" href="#" data-page="${
                      currentPage + 1
                    }" aria-label="Próximo">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>`;
    }

    html += "</ul>";
    paginationElement.innerHTML = html;

    // Adiciona eventos aos links de paginação
    paginationElement.querySelectorAll(".page-link").forEach((link) => {
      link.addEventListener("click", async (e) => {
        e.preventDefault();
        const page = parseInt(e.target.closest(".page-link").dataset.page);
        await loadData(page, entityName);
        history.pushState(null, "", `?page=${page}`);
      });
    });
  }
}
