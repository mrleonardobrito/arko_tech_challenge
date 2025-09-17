const HomeView = {
  template: `
    <div class="text-center py-5">
      <h1 class="display-4">Bem-vindo ao Arko Tech Challenge</h1>
      <p class="lead">Explore os dados de estados, municípios, distritos e empresas.</p>
      
      <div class="row mt-5">
        <div class="col-md-3">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Estados</h5>
              <p class="card-text">Lista de estados brasileiros</p>
              <router-link to="/estados" class="btn btn-primary">Ver Estados</router-link>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Municípios</h5>
              <p class="card-text">Lista de municípios por estado</p>
              <router-link to="/cidades" class="btn btn-primary">Ver Municípios</router-link>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Distritos</h5>
              <p class="card-text">Lista de distritos por município</p>
              <router-link to="/distritos" class="btn btn-primary">Ver Distritos</router-link>
            </div>
          </div>
        </div>
        <div class="col-md-3">
          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Empresas</h5>
              <p class="card-text">Lista de empresas cadastradas</p>
              <router-link to="/empresas" class="btn btn-primary">Ver Empresas</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
};
