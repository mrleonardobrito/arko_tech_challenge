const DistrictView = {
  template: `
    <div class="row">
      <div class="col">
        <h2>Distritos</h2>
        <data-table
          fetch-url="/api/districts/"
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
          key: "city",
          displayKey: "city.name",
          label: "Cidade",
          sortable: true,
        },
        {
          key: "state",
          displayKey: "city.state.name",
          label: "Estado",
          sortable: true,
        },
      ],
    };
  },
  components: { DataTable },
};
