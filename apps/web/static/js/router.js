const routes = [
  { path: "/", component: HomeView, label: "Home" },
  { path: "/estados", component: StateView, label: "Estados" },
  { path: "/cidades", component: CityView, label: "Cidades" },
  { path: "/distritos", component: DistrictView, label: "Distritos" },
  { path: "/empresas", component: CompanyView, label: "Empresas" },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});
