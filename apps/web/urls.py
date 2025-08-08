from django.urls import path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    path('estados/', views.state_list, name='state_list'),
    path('cidades/', views.city_list, name='city_list'),
    path('distritos/', views.district_list, name='district_list'),
    path('empresas/', views.company_list, name='company_list'),
]
