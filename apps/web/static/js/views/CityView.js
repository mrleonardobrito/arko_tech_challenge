const CityView = {
  template: `
    <div class="row">
      <div class="col">
        <h2>Cidades</h2>
        <data-table
          fetch-url="/api/cities/"
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
        {
          key: "state",
          displayKey: "state.name",
          label: "Estado",
          sortable: true,
        },
      ],
    };
  },
  components: { DataTable },
};
