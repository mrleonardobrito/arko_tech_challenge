from typing import Any, Dict


def get_friendly_error_message(error_type: str, details: Dict[str, Any] = None) -> str:
    messages = {
        'ValidationError': 'Os dados fornecidos são inválidos.',
        'PermissionDenied': 'Você não tem permissão para realizar esta ação.',
        'ObjectDoesNotExist': 'O recurso solicitado não foi encontrado.',
        'MultipleObjectsReturned': 'Foram encontrados múltiplos registros quando esperava-se apenas um.',
        'SuspiciousOperation': 'A operação solicitada foi considerada suspeita e foi bloqueada.',

        'DatabaseError': 'Ocorreu um erro ao acessar o banco de dados.',
        'IntegrityError': 'A operação viola a integridade do banco de dados.',
        'DataError': 'Os dados fornecidos são inválidos para o banco de dados.',
        'OperationalError': 'Ocorreu um erro operacional no banco de dados.',

        'Http404': 'O recurso solicitado não foi encontrado.',
        'BadRequest': 'A requisição é inválida ou mal formatada.',
        'PermissionDenied': 'Você não tem permissão para acessar este recurso.',

        'NotAuthenticated': 'Você precisa estar autenticado para realizar esta ação.',
        'AuthenticationFailed': 'Falha na autenticação. Verifique suas credenciais.',
        'NotFound': 'O recurso solicitado não foi encontrado.',
        'MethodNotAllowed': 'O método HTTP utilizado não é permitido para este recurso.',
        'NotAcceptable': 'O formato solicitado não está disponível.',
        'UnsupportedMediaType': 'O formato dos dados enviados não é suportado.',
        'Throttled': 'Você realizou muitas requisições. Tente novamente mais tarde.',
        'ParseError': 'Não foi possível interpretar os dados enviados.',

        'InvalidPageError': 'O número da página informado é inválido.',
        'InvalidPageSizeError': 'O tamanho da página informado é inválido.',
        'ResourceNotFoundError': 'O recurso solicitado não foi encontrado.',
        'ValidationError': 'Os dados fornecidos são inválidos.',
        'DatabaseError': 'Ocorreu um erro ao acessar o banco de dados.',

        'default': 'Ocorreu um erro inesperado. Por favor, tente novamente.'
    }

    if error_type == 'ValidationError' and details:
        field_errors = []
        for field, errors in details.items():
            if isinstance(errors, list):
                field_errors.append(f"{field}: {', '.join(errors)}")
            else:
                field_errors.append(f"{field}: {errors}")
        if field_errors:
            return f"Erros de validação: {'; '.join(field_errors)}"

    return messages.get(error_type, messages['default'])
