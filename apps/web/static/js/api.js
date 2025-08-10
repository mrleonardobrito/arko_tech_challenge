class API {
  static async fetchData(endpoint, page = 1, page_size = 20) {
    try {
      const response = await fetch(
        `/api/${endpoint}/?page=${page}&page_size=${page_size}`
      );
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

    const totalPages = Math.ceil(data.count / data.results.length);
    const maxVisiblePages = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
    let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

    if (endPage - startPage + 1 < maxVisiblePages) {
      startPage = Math.max(1, endPage - maxVisiblePages + 1);
    }

    let html = '<ul class="pagination justify-content-center">';

    html += `
      <li class="page-item ${currentPage === 1 ? "disabled" : ""}">
        <a class="page-link" href="#" data-page="1" aria-label="Primeira">
          <span aria-hidden="true">&laquo;&laquo;</span>
        </a>
      </li>`;

    html += `
      <li class="page-item ${!data.previous ? "disabled" : ""}">
        <a class="page-link" href="#" data-page="${
          currentPage - 1
        }" aria-label="Anterior">
          <span aria-hidden="true">&laquo;</span>
        </a>
      </li>`;

    if (startPage > 1) {
      html += `
        <li class="page-item disabled">
          <span class="page-link">...</span>
        </li>`;
    }

    for (let i = startPage; i <= endPage; i++) {
      html += `
        <li class="page-item ${i === currentPage ? "active" : ""}">
          <a class="page-link" href="#" data-page="${i}">${i}</a>
        </li>`;
    }

    if (endPage < totalPages) {
      html += `
        <li class="page-item disabled">
          <span class="page-link">...</span>
        </li>`;
    }

    html += `
      <li class="page-item ${!data.next ? "disabled" : ""}">
        <a class="page-link" href="#" data-page="${
          currentPage + 1
        }" aria-label="Próximo">
          <span aria-hidden="true">&raquo;</span>
        </a>
      </li>`;

    html += `
      <li class="page-item ${currentPage === totalPages ? "disabled" : ""}">
        <a class="page-link" href="#" data-page="${totalPages}" aria-label="Última">
          <span aria-hidden="true">&raquo;&raquo;</span>
        </a>
      </li>`;

    html += "</ul>";

    html += `
      <div class="text-center mt-2">
        <small class="text-muted">
          Mostrando ${data.results.length} de ${data.count} registros
        </small>
      </div>`;

    paginationElement.innerHTML = html;

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
