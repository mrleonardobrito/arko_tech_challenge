const DataTable = {
  name: "DataTable",
  props: {
    columns: { type: Array, default: () => [] },
    pageSize: { type: Number, default: 10 },
    fetchUrl: { type: String, required: true },
  },
  data() {
    return {
      page: 1,
      sortBy: null,
      sortDir: "asc",
      search: "",
      fetchedItems: [],
      isLoading: false,
      totalItems: 0,
      pageSizeOptions: [10, 20, 30, 50, 100],
      currentPageSize: this.pageSize,
    };
  },
  computed: {
    totalPages() {
      return Math.max(1, Math.ceil(this.totalItems / this.currentPageSize));
    },
    rows() {
      return this.fetchedItems;
    },
    paginationPages() {
      const total = this.totalPages;
      const current = this.page;
      const delta = 2;
      const range = [];
      const rangeWithDots = [];

      if (total <= 7) {
        for (let i = 1; i <= total; i++) {
          range.push(i);
        }
        return range;
      }
      const start = Math.max(1, current - delta);
      const end = Math.min(total, current + delta);
      for (let i = start; i <= end; i++) {
        range.push(i);
      }
      if (start > 1) {
        if (start > 2) {
          rangeWithDots.push(1, "...");
        } else {
          rangeWithDots.push(1);
        }
      }
      rangeWithDots.push(...range);
      if (end < total) {
        if (end < total - 1) {
          rangeWithDots.push("...", total);
        } else {
          rangeWithDots.push(total);
        }
      }

      return rangeWithDots;
    },
  },
  methods: {
    getNestedValue(obj, path) {
      return path.split(".").reduce((current, key) => {
        return current && current[key] !== undefined ? current[key] : "";
      }, obj);
    },
    onSearchInput() {
      this.page = 1;
      this.requestData();
    },
    onPageSizeChange() {
      this.page = 1;
      this.requestData();
    },
    async requestData() {
      this.isLoading = true;
      try {
        const url = new URL(this.fetchUrl, window.location.origin);
        url.searchParams.set("page", String(this.page));
        url.searchParams.set("page_size", String(this.currentPageSize));
        if (this.sortBy) {
          url.searchParams.set("sort_by", this.sortBy);
          url.searchParams.set("sort_dir", this.sortDir);
        }
        if (this.search) {
          url.searchParams.set("query", this.search);
        }

        const resp = await fetch(url.toString());
        if (!resp.ok) throw new Error("Falha ao buscar dados");
        const data = await resp.json();
        const results = Array.isArray(data)
          ? data
          : Array.isArray(data.results)
          ? data.results
          : [];
        this.fetchedItems = results;
        this.totalItems =
          typeof data.count === "number" ? data.count : results.length;
      } catch (e) {
        console.error(e);
      } finally {
        this.isLoading = false;
      }
    },
    changeSort(col) {
      if (!col.sortable) return;
      if (this.sortBy === col.key) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = col.key;
        this.sortDir = "asc";
      }
      this.page = 1;
      this.requestData();
    },
    goToPage(n) {
      if (n < 1 || n > this.totalPages) return;
      this.page = n;
      this.requestData();
    },
  },
  mounted() {
    if (!this.fetchUrl) return;
    this.requestData();
  },
  template: `
      <div>
        <div class="mb-2 d-flex justify-content-between align-items-center">
          <input v-model="search" @input="onSearchInput" class="form-control w-50" placeholder="Pesquisar...">
          <div class="d-flex align-items-center gap-2">
            <label for="pageSizeSelect" class="form-label mb-0">Itens por página:</label>
            <select 
              id="pageSizeSelect"
              v-model="currentPageSize" 
              @change="onPageSizeChange" 
              class="form-select w-auto"
            >
              <option v-for="size in pageSizeOptions" :key="size" :value="size" v-text="size">
              </option>
            </select>
            <div v-if="isLoading" class="text-muted">Carregando...</div>
          </div>
        </div>
  
        <table class="table table-striped table-hover">
          <thead>
            <tr>
              <th v-for="col in columns" :key="col.key" @click="changeSort(col)" style="cursor: pointer;">
                <span v-text="col.label"></span>
                <span v-if="col.sortable">
                  <small v-if="sortBy !== col.key">⇅</small>
                  <small v-else v-text="sortDir === 'asc' ? '▲' : '▼'"></small>
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in rows" :key="item.id" style="cursor:pointer;" tabindex="0">
              <td v-for="col in columns" :key="col.key">
                <slot :name="'cell-'+col.key" :item="item">
                  <span v-text="getNestedValue(item, col.displayKey || col.key)"></span>
                </slot>
              </td>
            </tr>
            <tr v-if="rows.length === 0">
              <td :colspan="columns.length" class="text-center">Nenhum registro</td>
            </tr>
          </tbody>
        </table>
  
        <nav aria-label="Paginação" class="d-flex justify-content-between align-items-center">
          <div>Mostrando página <span v-text="page"></span> / <span v-text="totalPages"></span></div>
          <ul class="pagination mb-0">
            <li class="page-item" :class="{ disabled: page===1 }">
              <button class="page-link" @click="goToPage(page-1)">Anterior</button>
            </li>
            <li v-for="p in paginationPages" :key="p" class="page-item" :class="{ active: p===page, disabled: p==='...' }">
              <button v-if="p !== '...'" class="page-link" @click="goToPage(p)" v-text="p"></button>
              <span v-else class="page-link" v-text="p"></span>
            </li>
            <li class="page-item" :class="{ disabled: page===totalPages }">
              <button class="page-link" @click="goToPage(page+1)">Próximo</button>
            </li>
          </ul>
        </nav>
      </div>
    `,
};
