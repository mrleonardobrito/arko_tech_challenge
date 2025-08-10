from rest_framework import serializers
from .models import State


class StateSerializer(serializers.ModelSerializer):
    """
    Serializador para o modelo State (Estado).

    Campos:
    - id: Identificador único do estado
    - name: Nome completo do estado
    - acronym: Sigla do estado (2 caracteres)
    """
    class Meta:
        model = State
        fields = ['id', 'name', 'acronym']
        swagger_schema_fields = {
            "title": "Estado",
            "description": "Representa um estado brasileiro com seu nome e sigla"
        }
