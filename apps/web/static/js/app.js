const config = {
  setup() {
    const tabs = ref([
      { label: "Estados", path: "/estados" },
      { label: "Cidades", path: "/cidades" },
      { label: "Distritos", path: "/distritos" },
      { label: "Empresas", path: "/empresas" },
    ]);

    return { tabs };
  },
};

const app = createApp(config);
app.use(router);
app.mount("#app");
