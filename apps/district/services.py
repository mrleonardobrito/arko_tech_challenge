from django.db import DatabaseError as DjangoDatabaseError
from django.db.models import Q
from .models import District
from apps.api.errors import DatabaseError


def list_districts(
    page: int = 1,
    page_size: int = 10,
    sort_by: str | None = None,
    sort_dir: str = "asc",
    query: str | None = None,
):
    try:
        queryset = District.objects.all()

        if query:
            ft = query.strip()
            if ft:
                queryset = queryset.filter(
                    Q(name__icontains=ft) | Q(city__name__icontains=ft) | Q(
                        city__state__name__icontains=ft) | Q(city__state__acronym__icontains=ft)
                )

        allowed_sorts = {
            "name": "name",
            "city": "city__name",
            "state": "city__state__name",
        }
        sort_field = allowed_sorts.get((sort_by or "").lower(), "name")
        sort_prefix = "-" if str(sort_dir).lower() == "desc" else ""
        queryset = queryset.order_by(f"{sort_prefix}{sort_field}")

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
        raise DatabaseError(f'Erro ao listar distritos: {str(e)}')
