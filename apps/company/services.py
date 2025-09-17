from django.db import DatabaseError as DjangoDatabaseError
from django.db.models import Q
from .models import Company
from apps.api.errors import DatabaseError


def list_companies(
    page: int = 1,
    page_size: int = 10,
    sort_by: str | None = None,
    sort_dir: str = "asc",
    query: str | None = None,
):
    try:
        queryset = Company.objects.all()

        if query:
            ft = query.strip()
            if ft:
                queryset = queryset.filter(
                    Q(social_name__icontains=ft) | Q(cnpj__icontains=ft) | Q(
                        federative_entity__icontains=ft)
                )

        allowed_sorts = {
            "social_name": "social_name",
            "cnpj": "cnpj",
            "federative_entity": "federative_entity",
            "social_capital": "social_capital",
            "company_size": "company_size",
        }
        sort_field = allowed_sorts.get((sort_by or "").lower(), "social_name")
        sort_prefix = "-" if str(sort_dir).lower() == "desc" else ""
        queryset = queryset.order_by(f"{sort_prefix}{sort_field}")

        if page == 1:
            try:
                results = queryset[:page_size + 1]
                has_more = len(results) > page_size
                if has_more:
                    results = results[:page_size]
                    total = page_size * 100
                else:
                    total = len(results)
            except:
                total = queryset.count()
                offset = (page - 1) * page_size
                results = queryset[offset: offset + page_size]
        else:
            total = queryset.count()
            offset = (page - 1) * page_size
            results = queryset[offset: offset + page_size]

        return {
            'results': results,
            'total': total,
            'page': page,
            'page_size': page_size,
            'total_pages': (total + page_size - 1) // page_size
        }
    except DjangoDatabaseError as e:
        raise DatabaseError(f'Erro ao listar empresas: {str(e)}')
