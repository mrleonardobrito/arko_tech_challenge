from rest_framework import viewsets
from rest_framework.response import Response
from drf_spectacular.utils import extend_schema, extend_schema_view, OpenApiParameter
from apps.api.errors import InvalidPageError, InvalidPageSizeError
from .models import Company
from .services import list_companies
from .serializers import CompanySerializer


@extend_schema_view(
    list=extend_schema(
        summary='Listar empresas',
        description='Retorna uma lista de empresas com paginação',
        tags=['Empresa'],
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
class CompanyViewSet(viewsets.GenericViewSet):
    """
    ViewSet para listar empresas.

    A listagem é paginada com 10 itens por página.
    Use os parâmetros ?page= e ?page_size= para controlar a paginação.
    """
    serializer_class = CompanySerializer

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

        companies = list_companies(page=page, page_size=page_size)
        serializer = self.get_serializer(companies['results'], many=True)

        return Response({
            'count': companies['total'],
            'next': f'/api/companies/?page={page + 1}&page_size={page_size}' if page < companies['total_pages'] else None,
            'previous': f'/api/companies/?page={page - 1}&page_size={page_size}' if page > 1 else None,
            'results': serializer.data
        })
