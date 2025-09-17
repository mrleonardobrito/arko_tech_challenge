const CompanyView = {
  template: `
    <div class="row">
      <div class="col">
        <h2>Empresas</h2>
        <data-table
          fetch-url="/api/companies/"
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
        { key: "cnpj", label: "CNPJ", sortable: true },
        { key: "social_name", label: "Razão Social", sortable: true },
        {
          key: "juridical_nature",
          label: "Natureza Jurídica",
          sortable: false,
        },
        { key: "social_capital", label: "Capital Social", sortable: true },
        { key: "company_size", label: "Porte", sortable: true },
        { key: "federative_entity", label: "UF", sortable: true },
      ],
    };
  },
  components: { DataTable },
};
