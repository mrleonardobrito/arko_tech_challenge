const StateView = {
  template: `
    <div class="row">
      <div class="col">
        <h2>Estados</h2>
        <data-table
          fetch-url="/api/states/"
          :server-side="true"
          :page-size="20"
          :columns="columns"
        ></data-table>
      </div>
    </div>
  `,
  data() {
    return {
      columns: [
        { key: "name", label: "Nome", sortable: true },
        { key: "acronym", label: "Sigla", sortable: true },
      ],
    };
  },
  components: { DataTable },
};
