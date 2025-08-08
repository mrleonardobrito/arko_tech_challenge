from rest_framework.exceptions import APIException
from rest_framework import status


class BaseAPIException(APIException):
    """
    Classe base para exceções da API.
    """

    def __init__(self, message=None, code=None):
        super().__init__(detail=message, code=code)


class InvalidPageError(BaseAPIException):
    """
    Exceção para erros de paginação.
    """
    status_code = status.HTTP_400_BAD_REQUEST
    default_detail = 'O número da página deve ser um número inteiro positivo.'
    default_code = 'invalid_page'


class InvalidPageSizeError(BaseAPIException):
    """
    Exceção para erros de tamanho de página.
    """
    status_code = status.HTTP_400_BAD_REQUEST
    default_detail = 'O tamanho da página deve estar entre 1 e 100 registros.'
    default_code = 'invalid_page_size'


class ResourceNotFoundError(BaseAPIException):
    """
    Exceção para recursos não encontrados.
    """
    status_code = status.HTTP_404_NOT_FOUND
    default_detail = 'O recurso solicitado não foi encontrado no sistema.'
    default_code = 'not_found'


class ValidationError(BaseAPIException):
    """
    Exceção para erros de validação.
    """
    status_code = status.HTTP_400_BAD_REQUEST
    default_detail = 'Os dados fornecidos são inválidos. Por favor, verifique e tente novamente.'
    default_code = 'validation_error'


class DatabaseError(BaseAPIException):
    """
    Exceção para erros de banco de dados.
    """
    status_code = status.HTTP_500_INTERNAL_SERVER_ERROR
    default_detail = 'Ocorreu um erro ao acessar o banco de dados. Por favor, tente novamente mais tarde.'
    default_code = 'database_error'
