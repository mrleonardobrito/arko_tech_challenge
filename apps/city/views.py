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

        cities = list_cities(page=page, page_size=page_size)
        serializer = self.get_serializer(cities['results'], many=True)

        return Response({
            'count': cities['total'],
            'next': f'/api/cities/?page={page + 1}&page_size={page_size}' if page < cities['total_pages'] else None,
            'previous': f'/api/cities/?page={page - 1}&page_size={page_size}' if page > 1 else None,
            'results': serializer.data
        })
