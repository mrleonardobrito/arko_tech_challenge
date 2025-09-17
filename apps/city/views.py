from rest_framework import viewsets
from rest_framework.response import Response
from drf_spectacular.utils import extend_schema, extend_schema_view, OpenApiParameter
from apps.api.errors import InvalidPageError, InvalidPageSizeError
from .models import City
from .services import list_cities
from .serializers import CitySerializer


@extend_schema_view(
    list=extend_schema(
        summary='Listar cidades',
        description='Retorna uma lista de cidades com paginação',
        tags=['Cidade'],
        parameters=[
            OpenApiParameter(
                name='page',
                type=int,
                description='Número da página (padrão: 1)',
                required=False
            ),
            OpenApiParameter(
                name='page_size',
                type=int,
                description='Quantidade de itens por página (padrão: 10, máximo: 100)',
                required=False
            ),
            OpenApiParameter(
                name='sort_by',
                type=str,
                description='Campo para ordenação (name, state)',
                required=False
            ),
            OpenApiParameter(
                name='sort_dir',
                type=str,
                description='Direção da ordenação (asc, desc)',
                required=False
            ),
            OpenApiParameter(
                name='query',
                type=str,
                description='Termo de busca (nome da cidade, estado)',
                required=False
            )
        ]
    )
)
class CityViewSet(viewsets.GenericViewSet):
    """
    ViewSet para listar cidades.

    A listagem é paginada com 10 itens por página.
    Use os parâmetros ?page= e ?page_size= para controlar a paginação.
    """
    serializer_class = CitySerializer

    def list(self, request):
        try:
            page = int(request.query_params.get('page', 1))
            if page < 1:
                raise InvalidPageError(
                    'O número da página deve ser maior que zero')
        except ValueError:
            raise InvalidPageError(
                'O número da página deve ser um número inteiro')

        try:
            page_size = int(request.query_params.get('page_size', 10))
            if page_size < 1 or page_size > 100:
                raise InvalidPageSizeError(
                    'O tamanho da página deve estar entre 1 e 100')
        except ValueError:
            raise InvalidPageSizeError(
                'O tamanho da página deve ser um número inteiro')

        sort_by = request.query_params.get('sort_by')
        sort_dir = request.query_params.get('sort_dir', 'asc')
        query = request.query_params.get('query')

        cities = list_cities(
            page=page,
            page_size=page_size,
            sort_by=sort_by,
            sort_dir=sort_dir,
            query=query,
        )
        serializer = self.get_serializer(cities['results'], many=True)

        def build_url(target_page: int | None):
            if not target_page:
                return None
            params = [f'page={target_page}', f'page_size={page_size}']
            if sort_by:
                params.append(f'sort_by={sort_by}')
            if sort_dir:
                params.append(f'sort_dir={sort_dir}')
            if query:
                params.append(f'query={query}')
            return f"/api/cities/?{'&'.join(params)}"

        return Response({
            'count': cities['total'],
            'next': build_url(page + 1 if page < cities['total_pages'] else None),
            'previous': build_url(page - 1 if page > 1 else None),
            'results': serializer.data
        })
