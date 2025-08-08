from django.shortcuts import render


def home(request):
    return render(request, 'home.html')


def state_list(request):
    return render(request, 'state/list.html')


def city_list(request):
    return render(request, 'city/list.html')


def district_list(request):
    return render(request, 'district/list.html')


def company_list(request):
    return render(request, 'company/list.html')
