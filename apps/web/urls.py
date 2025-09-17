from django.urls import path, re_path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    # Catch-all pattern para servir a SPA, mas excluindo rotas da API e admin
    re_path(r'^(?!api/|admin/).*$', views.home, name='spa-fallback'),
]
