from django.db import DatabaseError as DjangoDatabaseError
from .models import City
from apps.api.errors import DatabaseError


def list_cities(page=1, page_size=10):
    try:
        queryset = City.objects.order_by('name')
        total = queryset.count()
        offset = (page - 1) * page_size
        results = queryset[offset:offset+page_size]
        return {
            'results': results,
            'total': total,
            'page': page,
            'page_size': page_size,
            'total_pages': (total + page_size - 1) // page_size
        }
    except DjangoDatabaseError as e:
        raise DatabaseError(f'Erro ao listar cidades: {str(e)}')
