from typing import Optional
from rest_framework.views import exception_handler
from rest_framework.response import Response
from rest_framework import status
from django.core.exceptions import ValidationError as DjangoValidationError
from django.db import IntegrityError
from django.http import Http404
from .messages import get_friendly_error_message


def custom_exception_handler(exc: Exception, context: dict) -> Optional[Response]:
    response = exception_handler(exc, context)
    error_type = exc.__class__.__name__

    if response is not None:
        details = response.data if hasattr(response, 'data') else None
        data = {
            'error': {
                'type': error_type,
                'message': get_friendly_error_message(error_type, details),
            }
        }
        response.data = data
        return response

    if isinstance(exc, DjangoValidationError):
        details = exc.message_dict if hasattr(
            exc, 'message_dict') else str(exc)
        data = {
            'error': {
                'type': 'ValidationError',
                'message': get_friendly_error_message('ValidationError', details),
                'details': details
            }
        }
        return Response(data, status=status.HTTP_400_BAD_REQUEST)

    if isinstance(exc, IntegrityError):
        data = {
            'error': {
                'type': 'IntegrityError',
                'message': get_friendly_error_message('IntegrityError'),
            }
        }
        return Response(data, status=status.HTTP_400_BAD_REQUEST)

    if isinstance(exc, Http404):
        data = {
            'error': {
                'type': 'NotFound',
                'message': get_friendly_error_message('NotFound'),
            }
        }
        return Response(data, status=status.HTTP_404_NOT_FOUND)

    data = {
        'error': {
            'type': error_type,
            'message': get_friendly_error_message('default'),
        }
    }
    return Response(data, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
